package httputil

import (
	"encoding/json"
	"net/http"
)

// WriteJSON is a helper that sets Content-Type and encodes v as JSON.
func WriteJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// If encoding fails, we can't do much on this writer, 
		// but standardizing the helper is the first step.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
