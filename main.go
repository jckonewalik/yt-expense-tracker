package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jckonewalik/yt-expense-tracker/config"
	"github.com/jckonewalik/yt-expense-tracker/logger"
	"github.com/jckonewalik/yt-expense-tracker/services/auth"
	"github.com/jckonewalik/yt-expense-tracker/services/httputils"
	"github.com/jckonewalik/yt-expense-tracker/types"
	_ "github.com/lib/pq"
)

var validate *validator.Validate

func main() {

	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	route := http.NewServeMux()
	route.HandleFunc("GET /hello", auth.WithJWT(handleHello))
	route.HandleFunc("POST /signup", handleSignup)

	v1 := http.StripPrefix("/api/v1", route)

	err := http.ListenAndServe(":3000", v1)
	if err != nil {
		logger.Log.Error.Printf("error starting server: %v", err)
	}
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(types.UserName)
	httputils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Hello, %s", name))
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type KeycloakErrorResponse struct {
	Message string `json:"errorMessage"`
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	// decode payload
	var input types.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid input"))
		return
	}

	// validate payload
	err = validate.Struct(input)
	if err != nil {
		var payloadErrs = []error{}
		for _, err := range err.(validator.ValidationErrors) {
			payloadErrs = append(payloadErrs, fmt.Errorf("%s: %s", err.Field(), err.Tag()))
		}
		httputils.WriteErrors(w, http.StatusBadRequest, payloadErrs)
		return
	}

	if input.Password != input.PasswordConfirmation {
		httputils.WriteError(w, http.StatusBadRequest, fmt.Errorf("password confirmation doesn't match"))
		return
	}

	// check if user already exists

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", "ytexpensestracker", "admin@123", "ytexpensestracker"))
	if err != nil {
		logger.Log.Error.Printf("error getting db access. %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	err = db.Ping()
	if err != nil {
		logger.Log.Error.Printf("error pinging db. %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1 FROM users WHERE email = $1 OR login = $2", input.Email, input.Login)
	if err != nil {
		logger.Log.Error.Printf("error quering users table. %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	defer rows.Close()
	if rows.Next() {
		httputils.WriteError(w, http.StatusBadRequest, fmt.Errorf("there is already a user registered with this login or email"))
		return
	}

	// persist user in database
	_, err = db.Exec("INSERT INTO users (login, email, first_name, last_name) VALUES ($1, $2, $3, $4)", input.Login, input.Email, input.FirstName, input.LastName)
	if err != nil {
		logger.Log.Error.Printf("error inserting user. payload: %v. error: %v", input, err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}

	// create user in keycloak
	requestBody := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{"my-api"},
		"client_secret": []string{config.Env.AuthApiClientSecret},
	}
	dataReader := requestBody.Encode()
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/realms/yt-expense-tracker/protocol/openid-connect/token",
		strings.NewReader(dataReader))
	if err != nil {
		logger.Log.Error.Printf("error getting keycloak access token. %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		logger.Log.Error.Printf("error getting keycloak access token. %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		var resp KeycloakErrorResponse
		err := json.NewDecoder(res.Body).Decode(&resp)
		if err != nil {
			logger.Log.Error.Printf("error reading keycloak token response. %v", err)
			httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
			return
		}
		logger.Log.Error.Printf("error getting keycloak access token. %v", resp.Message)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}

	var token TokenResponse
	json.NewDecoder(res.Body).Decode(&token)

	userPayload := map[string]any{"email": input.Email, "username": input.Login, "credentials": []map[string]any{{"type": "password", "value": input.Password, "temporary": false}},
		"firstName": input.FirstName, "lastName": input.LastName, "emailVerified": false, "enabled": true, "requiredActions": []any{}, "groups": []any{}}
	jsonBody, err := json.Marshal(userPayload)
	if err != nil {
		logger.Log.Error.Printf("error encoding create user payload: %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	bodyReader := bytes.NewReader(jsonBody)
	req, err = http.NewRequest(http.MethodPost, "http://localhost:8080/admin/realms/yt-expense-tracker/users", bodyReader)
	if err != nil {
		logger.Log.Error.Printf("error creating keycloak users request: %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err = client.Do(req)
	if err != nil {
		logger.Log.Error.Printf("error requesting keycloak to create user: %v", err)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		var resp KeycloakErrorResponse
		err := json.NewDecoder(res.Body).Decode(&resp)
		if err != nil {
			logger.Log.Error.Printf("error reading keycloak users response. %v", err)
			httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
			return
		}
		logger.Log.Error.Printf("error creating keycloak user. %v", resp.Message)
		httputils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happend"))
		return
	}

	// return response
	httputils.WriteJSON(w, http.StatusCreated, nil)
}
