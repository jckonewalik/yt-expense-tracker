package httputils

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		json.NewEncoder(w).Encode(body)
	}
}
