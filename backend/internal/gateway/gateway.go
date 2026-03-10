package gateway

import (
	"context"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
)

// Pinger abstract the upstream connection, whether it resolves via REST or gRPC
type Pinger interface {
	Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error)
}

// PingProxy implements PingServiceServer by forwarding calls to a backend client.
type PingProxy struct {
	pb.UnimplementedPingServiceServer
	Backend Pinger
}

// Ping forwards the request to the upstream backend.
func (p *PingProxy) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return p.Backend.Ping(ctx, req)
}
