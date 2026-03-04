package gateway

import (
	"context"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

// BackendClient abstract the upstream connection, whether it resolves via REST or gRPC
type BackendClient interface {
	Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error)
}

// PingProxy implements PingServiceServer by forwarding calls to a backend client.
type PingProxy struct {
	pb.UnimplementedPingServiceServer
	Backend BackendClient
}

// Ping forwards the request to the upstream backend.
func (p *PingProxy) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return p.Backend.Ping(ctx, req)
}
