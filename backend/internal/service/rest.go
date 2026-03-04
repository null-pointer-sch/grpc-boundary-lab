package service

import (
	"encoding/json"
	"net/http"
)

// RESTPingHandler serves the backend's REST ping endpoint.
type RESTPingHandler struct{}

// HandlePing responds with a JSON pong message.
func (h *RESTPingHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg := r.URL.Query().Get("message")
	if msg == "" {
		msg = "ping from frontend"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "pong: " + msg,
	})
}
