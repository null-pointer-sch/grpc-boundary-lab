package gateway_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/gateway"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockBackendClient mocks the gateway.Pinger interface
type mockBackendClient struct {
	Response string
	Err      error
}

func (m *mockBackendClient) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return &pb.PingResponse{Message: m.Response}, nil
}

func TestRESTServer_Mode(t *testing.T) {
	grpcClient := &mockBackendClient{Response: "grpc"}
	restClient := &mockBackendClient{Response: "rest"}

	server := gateway.NewRESTServer("", grpcClient, restClient, nil, nil)

	// Test target fallback to grpc
	req := httptest.NewRequest(http.MethodGet, "/api/mode", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var payload map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	assert.Equal(t, "grpc", payload["protocol"])
	assert.Equal(t, false, payload["tls"])

	// Test explicit target rest
	req = httptest.NewRequest(http.MethodGet, "/api/mode?target=rest", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	assert.Equal(t, "rest", payload["protocol"])
}

func TestRESTServer_Ping(t *testing.T) {
	grpcClient := &mockBackendClient{Response: "grpc-pong"}
	server := gateway.NewRESTServer("", grpcClient, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var payload map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	assert.Equal(t, "grpc-pong", payload["message"])
}

func TestRESTServer_Bench(t *testing.T) {
	server := gateway.NewRESTServer("", nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/bench/latest?target=grpc", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var payload gateway.BenchmarkData
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	assert.Equal(t, "grpc", payload.Protocol)
	assert.False(t, payload.TLS)
}
