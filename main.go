package main

import (
	"fmt"
	"net/http"

	"github.com/jckonewalik/yt-expense-tracker/services/auth"
	"github.com/jckonewalik/yt-expense-tracker/services/httputils"
	"github.com/jckonewalik/yt-expense-tracker/types"
)

func main() {

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

	// validate payload

	// check if user already exists

	// persist user in database

	// create user in keycloak

	// return response
}
