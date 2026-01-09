package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")

	res := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}

func ErrorJSON(w http.ResponseWriter, statusCode int, message string, errs any) {
	w.Header().Set("Content-Type", "application/json")
	res := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Errors  any    `json:"errors,omitempty"`
	}{
		Success: false,
		Message: message,
		Errors:  errs,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)

}
