package gateway_test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/gateway"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// mockBackend is a trivial BackendClient that prepends "mock-pong: ".
type mockBackend struct {
}

func (m *mockBackend) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "mock-pong: " + req.GetMessage()}, nil
}

func TestGatewayForwarding(t *testing.T) {
	// --- Start mock backend on bufconn ---
	backendLis := bufconn.Listen(bufSize)
	backendSrv := grpc.NewServer()

	type mockPingService struct {
		pb.UnimplementedPingServiceServer
	}
	pb.RegisterPingServiceServer(backendSrv, &mockPingService{})
	go func() {
		if err := backendSrv.Serve(backendLis); err != nil {
			log.Fatalf("backend serve: %v", err)
		}
	}()
	defer backendSrv.Stop()

	backendConn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return backendLis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial backend: %v", err)
	}
	defer backendConn.Close()

	// --- Start gateway on bufconn ---
	gatewayLis := bufconn.Listen(bufSize)
	gatewaySrv := grpc.NewServer()
	pb.RegisterPingServiceServer(gatewaySrv, &gateway.PingProxy{
		Backend: &mockBackend{},
	})
	go func() {
		if err := gatewaySrv.Serve(gatewayLis); err != nil {
			log.Fatalf("gateway serve: %v", err)
		}
	}()
	defer gatewaySrv.Stop()

	gwConn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return gatewayLis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial gateway: %v", err)
	}
	defer gwConn.Close()

	// --- Test ---
	client := pb.NewPingServiceClient(gwConn)
	resp, err := client.Ping(context.Background(), &pb.PingRequest{Message: "hello"})
	if err != nil {
		t.Fatalf("ping: %v", err)
	}
	if resp.GetMessage() != "mock-pong: hello" {
		t.Errorf("got %q, want %q", resp.GetMessage(), "mock-pong: hello")
	}
}
