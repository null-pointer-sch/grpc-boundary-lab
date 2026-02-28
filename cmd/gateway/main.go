package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndySchubert/grpc-boundary-lab/gateway"
	pb "github.com/AndySchubert/grpc-boundary-lab/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := envOrDefault("GATEWAY_PORT", "50052")
	backendHost := envOrDefault("BACKEND_HOST", "127.0.0.1")
	backendPort := envOrDefault("BACKEND_PORT", "50051")

	backendAddr := fmt.Sprintf("%s:%s", backendHost, backendPort)

	backendConn, err := grpc.NewClient(
		backendAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to backend at %s: %v", backendAddr, err)
	}
	defer backendConn.Close()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterPingServiceServer(srv, &gateway.PingProxy{
		Backend: pb.NewPingServiceClient(backendConn),
	})

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Fprintln(os.Stderr, "shutting down gateway...")
		srv.GracefulStop()
		backendConn.Close()
	}()

	fmt.Printf("gateway listening on :%s (forwarding to %s)\n", port, backendAddr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
