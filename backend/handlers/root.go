package handlers

import (
	"encoding/json"
	"net/http"
)

// RootResponse defines the JSON structure for /
type RootResponse struct {
	Message string `json:"message"`
}

// RootHandler returns a simple greeting
func RootHandler(w http.ResponseWriter, r *http.Request) {
	response := RootResponse{
		Message: "Hello from Four in a Row Server! (Test)",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
