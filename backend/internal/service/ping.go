// Package service contains the core business logic for the backend.
package service

import (
	"context"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

// PingServer implements the gRPC PingService.
type PingServer struct {
	pb.UnimplementedPingServiceServer
}

// Ping replies with "pong: " followed by the incoming message.
func (s *PingServer) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong: " + req.GetMessage()}, nil
}
