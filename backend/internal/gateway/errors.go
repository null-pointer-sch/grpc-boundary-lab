package gateway

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized JSON error structure.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// WriteError sends a standardized JSON error response.
func WriteError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(code),
		Message: err.Error(),
	})
}

// WriteErrorMessage sends a standardized JSON error response with a custom message.
func WriteErrorMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}
