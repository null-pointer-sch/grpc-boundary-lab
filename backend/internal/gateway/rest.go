package gateway

import (
	"net/http"
	"time"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/httputil"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

// RESTServer exposes the gateway's HTTP API.
// It holds both plaintext and (optional) TLS backend clients,
// allowing the frontend to toggle TLS per-request via ?tls=true.
type RESTServer struct {
	overrideProtocol string

	// Plaintext clients (always available).
	grpcBackend BackendClient
	restBackend BackendClient

	// TLS clients (nil when certs are not present).
	grpcBackendTLS BackendClient
	restBackendTLS BackendClient

	// Whether TLS clients were successfully initialised.
	TLSAvailable bool

	stats *StatsProvider
	mux   *http.ServeMux
}

// NewRESTServer creates a RESTServer and registers all routes once.
func NewRESTServer(override string, grpcClient, restClient, grpcTLS, restTLS BackendClient) *RESTServer {
	s := &RESTServer{
		overrideProtocol: override,
		grpcBackend:      grpcClient,
		restBackend:      restClient,
		grpcBackendTLS:   grpcTLS,
		restBackendTLS:   restTLS,
		TLSAvailable:     grpcTLS != nil,
		stats:            NewStatsProvider(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/mode", s.handleMode)
	mux.HandleFunc("/api/ping", s.handlePing)
	mux.HandleFunc("/api/bench/latest", s.handleBench)
	s.mux = mux

	return s
}

// ServeHTTP dispatches to the pre-built mux.
func (s *RESTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	s.mux.ServeHTTP(w, r)
}

func (s *RESTServer) targetProtocol(r *http.Request) string {
	if s.overrideProtocol != "" {
		return s.overrideProtocol
	}
	if t := r.URL.Query().Get("target"); t != "" {
		return t
	}
	return "grpc"
}

func (s *RESTServer) wantTLS(r *http.Request) bool {
	return r.URL.Query().Get("tls") == "true"
}

func (s *RESTServer) pickClient(r *http.Request) (client BackendClient, activeTLS bool) {
	target := s.targetProtocol(r)
	useTLS := s.wantTLS(r) && s.TLSAvailable

	if target == "rest" {
		if useTLS && s.restBackendTLS != nil {
			return s.restBackendTLS, true
		}
		return s.restBackend, false
	}
	if useTLS && s.grpcBackendTLS != nil {
		return s.grpcBackendTLS, true
	}
	return s.grpcBackend, false
}

// ── Handlers ─────────────────────────────────────────────────

func (s *RESTServer) handleMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.WriteErrorMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	_, activeTLS := s.pickClient(r)

	httputil.WriteJSON(w, map[string]any{
		"protocol":     s.targetProtocol(r),
		"tls":          activeTLS,
		"tlsAvailable": s.TLSAvailable,
	})
}

func (s *RESTServer) handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.WriteErrorMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	start := time.Now()
	client, _ := s.pickClient(r)

	resp, err := client.Ping(r.Context(), &pb.PingRequest{Message: "ping from frontend"})
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.WriteJSON(w, map[string]any{
		"message":   resp.GetMessage(),
		"latencyMs": time.Since(start).Milliseconds(),
	})
}

func (s *RESTServer) handleBench(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.WriteErrorMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	target := s.targetProtocol(r)
	_, activeTLS := s.pickClient(r)

	data, ok := s.stats.GetStats(target, activeTLS)
	if !ok {
		httputil.WriteErrorMessage(w, http.StatusNotFound, "Stats not found")
		return
	}

	httputil.WriteJSON(w, data)
}
