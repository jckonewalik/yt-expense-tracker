package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jckonewalik/yt-expense-tracker/services/auth"
	"github.com/jckonewalik/yt-expense-tracker/services/httputils"
	"github.com/jckonewalik/yt-expense-tracker/types"
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

	http.ListenAndServe(":3000", v1)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(types.UserName)
	httputils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Hello, %s", name))
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	// decode payload
	var input types.SignUpInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid input"))
		return
	}

	// validate payload
	if err := validate.Struct(input); err != nil {
		var payloadErrors = []error{}
		for _, err := range err.(validator.ValidationErrors) {
			payloadErrors = append(payloadErrors, fmt.Errorf("%s: %s", err.Field(), err.Tag()))
		}

		httputils.WriteErrors(w, http.StatusBadRequest, payloadErrors)
		return
	}

	if input.Password != input.PasswordConfirmation {
		httputils.WriteError(w, http.StatusBadRequest, fmt.Errorf("password confirmation doesn't match"))
		return
	}

	// persist user into database

	// create user in keycloak

	// return response
	httputils.WriteJSON(w, http.StatusCreated, nil)
}
