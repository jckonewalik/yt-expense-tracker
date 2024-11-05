package httputils

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string][]string{"errors": []string{err.Error()}})
}

func WriteErrors(w http.ResponseWriter, status int, errors []error) {
	var sErrors []string
	for _, err := range errors {
		sErrors = append(sErrors, err.Error())
	}

	WriteJSON(w, status, map[string][]string{"errors": sErrors})
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		json.NewEncoder(w).Encode(body)
	}
}
