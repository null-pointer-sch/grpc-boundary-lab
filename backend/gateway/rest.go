package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

type RESTServer struct {
	Backend     pb.PingServiceClient
	BackendREST string
	httpClient  *http.Client
}

func NewRESTServer(backend pb.PingServiceClient, backendRestAddr string) *RESTServer {
	return &RESTServer{
		Backend:     backend,
		BackendREST: backendRestAddr,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
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

func (s *RESTServer) handleMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	target := r.URL.Query().Get("target")
	if target == "" {
		target = "grpc"
	}

	data := map[string]interface{}{
		"protocol": target,
		"tls":      false,
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
	target := r.URL.Query().Get("target")

	var msg string

	if target == "rest" {
		url := fmt.Sprintf("http://%s/api/ping?message=ping%%20from%%20frontend", s.BackendREST)
		resp, err := s.httpClient.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var res map[string]interface{}
		json.Unmarshal(body, &res)
		if m, ok := res["message"].(string); ok {
			msg = m
		} else {
			msg = "pong (rest)"
		}
	} else {
		// Call the gRPC backend
		req := &pb.PingRequest{Message: "ping from frontend"}
		resp, err := s.Backend.Ping(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		msg = resp.GetMessage()
	}

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

	target := r.URL.Query().Get("target")

	// Provide some dummy baseline data simulating the loadgen results for now
	data := map[string]interface{}{
		"protocol": target,
		"tls":      false,
	}

	if target == "rest" {
		data["rps"] = 7455.45
		data["p50"] = 3.74
		data["p95"] = 7.25
		data["p99"] = 9.15
	} else {
		data["rps"] = 22348.19
		data["p50"] = 1.30
		data["p95"] = 2.18
		data["p99"] = 2.74
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
