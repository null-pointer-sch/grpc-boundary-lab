package main

import (
	"encoding/json"
	"net/http"
)

type restPingServer struct {
	addr string
}

func (s *restPingServer) handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg := r.URL.Query().Get("message")
	if msg == "" {
		msg = "ping from frontend"
	}

	data := map[string]interface{}{
		"message": "pong: " + msg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
