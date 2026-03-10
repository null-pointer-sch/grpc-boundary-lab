// Package service contains the core business logic for the backend.
package service

import (
	"context"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/core"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

// PingServer implements the gRPC PingService.
type PingServer struct {
	pb.UnimplementedPingServiceServer
}

// Ping replies with a pong message using the core business logic.
func (s *PingServer) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: core.GeneratePong(req.GetMessage())}, nil
}
