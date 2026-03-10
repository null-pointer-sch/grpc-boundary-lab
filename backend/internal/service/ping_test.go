package service_test

import (
	"context"
	"testing"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingServer_Ping(t *testing.T) {
	s := &service.PingServer{}

	req := &pb.PingRequest{Message: "test"}
	res, err := s.Ping(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "pong: test", res.GetMessage())
}
