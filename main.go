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

	v1 := http.StripPrefix("/api/v1", route)

	http.ListenAndServe(":3000", v1)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(types.UserName)
	httputils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Hello, %s", name))
}
