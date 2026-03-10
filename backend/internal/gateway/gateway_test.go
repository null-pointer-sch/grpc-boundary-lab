package gateway_test

import (
	"context"
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/gateway"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingProxy_Ping(t *testing.T) {
	mockBackend := &mockBackendClient{Response: "proxied-pong"}
	proxy := &gateway.PingProxy{Backend: mockBackend}

	req := &pb.PingRequest{Message: "test"}
	res, err := proxy.Ping(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "proxied-pong", res.Message)
}
