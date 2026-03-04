package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type pingServer struct {
	pb.UnimplementedPingServiceServer
}

func (s *pingServer) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong: " + req.GetMessage()}, nil
}

func main() {
	port := envOrDefault("BACKEND_PORT", "50051")
	restPort := envOrDefault("BACKEND_REST_PORT", "8081")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterPingServiceServer(srv, &pingServer{})
	reflection.Register(srv)

	restMux := http.NewServeMux()
	restSrv := &restPingServer{addr: ":" + restPort}
	restMux.HandleFunc("/api/ping", restSrv.handlePing)

	httpServer := &http.Server{
		Addr:    ":" + restPort,
		Handler: restMux,
	}

	go func() {
		fmt.Printf("backend REST listening on :%s\n", restPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve REST: %v", err)
		}
	}()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Fprintln(os.Stderr, "shutting down backend...")
		srv.GracefulStop()
		httpServer.Close()
	}()

	fmt.Printf("backend listening on :%s\n", port)
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
