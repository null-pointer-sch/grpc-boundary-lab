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
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/tlsutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	tlsEnabled := envOrDefault("TLS", "0") == "1"
	certDir := envOrDefault("CERT_DIR", "/certs")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcOpts []grpc.ServerOption
	restMux := http.NewServeMux()
	restSrvHandler := &restPingServer{addr: ":" + restPort}
	restMux.HandleFunc("/api/ping", restSrvHandler.handlePing)

	httpServer := &http.Server{
		Addr:    ":" + restPort,
		Handler: restMux,
	}

	if tlsEnabled {
		tlsConfig, err := tlsutil.LoadServerConfig(certDir+"/backend.crt", certDir+"/backend.key")
		if err != nil {
			log.Fatalf("failed to load TLS cert: %v", err)
		}
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(tlsConfig)))
		httpServer.TLSConfig = tlsConfig
	}

	srv := grpc.NewServer(grpcOpts...)
	pb.RegisterPingServiceServer(srv, &pingServer{})
	reflection.Register(srv)

	go func() {
		fmt.Printf("backend REST listening on :%s (TLS=%v)\n", restPort, tlsEnabled)
		var err error
		if tlsEnabled {
			err = httpServer.ListenAndServeTLS("", "")
		} else {
			err = httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
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

	fmt.Printf("backend listening on :%s (TLS=%v)\n", port, tlsEnabled)
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
