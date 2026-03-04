package gateway

import (
	"context"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PingProxy implements PingServiceServer by forwarding calls to a backend client.
type PingProxy struct {
	pb.UnimplementedPingServiceServer
	Backend pb.PingServiceClient
}

// Ping forwards the request to the backend and returns its response.
func (p *PingProxy) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	resp, err := p.Backend.Ping(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "backend call failed: %v", err)
	}
	return resp, nil
}
