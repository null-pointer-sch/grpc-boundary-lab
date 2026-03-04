package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

type RESTServer struct {
	OverrideProtocol string
	TLSEnabled       bool
	GrpcBackend      BackendClient
	RestBackend      BackendClient
}

func NewRESTServer(override string, tlsEnabled bool, grpcClient BackendClient, restClient BackendClient) *RESTServer {
	return &RESTServer{
		OverrideProtocol: override,
		TLSEnabled:       tlsEnabled,
		GrpcBackend:      grpcClient,
		RestBackend:      restClient,
	}
}

func (s *RESTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/mode", s.handleMode)
	mux.HandleFunc("/api/ping", s.handlePing)
	mux.HandleFunc("/api/bench/latest", s.handleBench)

	// Wrap with basic CORS just in case
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	mux.ServeHTTP(w, r)
}

func (s *RESTServer) getTargetProtocol(r *http.Request) string {
	target := s.OverrideProtocol
	if target == "" {
		target = r.URL.Query().Get("target")
		if target == "" {
			target = "grpc"
		}
	}
	return target
}

func (s *RESTServer) handleMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := map[string]interface{}{
		"protocol": s.getTargetProtocol(r),
		"tls":      s.TLSEnabled,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *RESTServer) handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	target := s.getTargetProtocol(r)

	var client BackendClient
	if target == "rest" {
		client = s.RestBackend
	} else {
		client = s.GrpcBackend
	}

	req := &pb.PingRequest{Message: "ping from frontend"}
	resp, err := client.Ping(r.Context(), req)

	var msg string
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg = resp.GetMessage()

	latency := time.Since(start).Milliseconds()

	data := map[string]interface{}{
		"message":   msg,
		"latencyMs": latency,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *RESTServer) handleBench(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	target := s.getTargetProtocol(r)

	data := map[string]interface{}{
		"protocol": target,
		"tls":      s.TLSEnabled,
	}

	if target == "rest" {
		if s.TLSEnabled {
			data["rps"] = 2331.81
			data["p50"] = 12.01
			data["p95"] = 27.32
			data["p99"] = 35.71
		} else {
			data["rps"] = 6973.90
			data["p50"] = 3.74
			data["p95"] = 9.63
			data["p99"] = 13.36
		}
	} else {
		if s.TLSEnabled {
			data["rps"] = 21606.62
			data["p50"] = 1.37
			data["p95"] = 2.28
			data["p99"] = 2.68
		} else {
			data["rps"] = 23561.92
			data["p50"] = 1.23
			data["p95"] = 2.11
			data["p99"] = 2.46
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
