package service

import (
	"net/http"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/core"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/httputil"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

// RESTPingHandler serves the backend's REST ping endpoint.
type RESTPingHandler struct{}

// HandlePing responds with a JSON pong message.
func (h *RESTPingHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.WriteErrorMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	msg := r.URL.Query().Get("message")
	resp := &pb.PingResponse{
		Message: core.GeneratePong(msg),
	}

	httputil.WriteJSON(w, resp)
}
