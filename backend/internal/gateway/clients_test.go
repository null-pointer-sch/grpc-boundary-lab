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
	"google.golang.org/grpc"
)

type mockPingClient struct {
	Response string
	Err      error
}

func (m *mockPingClient) Ping(ctx context.Context, in *pb.PingRequest, opts ...grpc.CallOption) (*pb.PingResponse, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return &pb.PingResponse{Message: m.Response}, nil
}

func TestGrpcBackendClient_Ping(t *testing.T) {
	client := &gateway.GrpcBackendClient{
		Client: &mockPingClient{Response: "grpc-pong"},
	}

	req := &pb.PingRequest{Message: "hello"}
	res, err := client.Ping(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "grpc-pong", res.Message)
}

func TestRestBackendClient_Ping(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/ping", r.URL.Path)
		assert.Equal(t, "hello from test", r.URL.Query().Get("message"))

		resp := map[string]string{"message": "rest-pong"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := &gateway.RestBackendClient{
		TargetURL:  ts.URL,
		HTTPClient: ts.Client(),
	}

	req := &pb.PingRequest{Message: "hello from test"}
	res, err := client.Ping(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "rest-pong", res.Message)
}

func TestRestBackendClient_Ping_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := &gateway.RestBackendClient{
		TargetURL:  ts.URL,
		HTTPClient: ts.Client(),
	}

	req := &pb.PingRequest{Message: "hello from test"}
	_, err := client.Ping(context.Background(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
