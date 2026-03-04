package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/null-pointer-sch/grpc-boundary-lab/gateway"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/tlsutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := envOrDefault("GATEWAY_PORT", "50052")
	restPort := envOrDefault("REST_PORT", "8080")
	backendHost := envOrDefault("BACKEND_HOST", "127.0.0.1")
	backendPort := envOrDefault("BACKEND_PORT", "50051")
	backendRestPort := envOrDefault("BACKEND_REST_PORT", "8081")

	protocol := os.Getenv("PROTOCOL")
	tlsEnabled := os.Getenv("TLS") == "1"
	certDir := envOrDefault("CERT_DIR", "/certs")

	backendAddr := fmt.Sprintf("%s:%s", backendHost, backendPort)
	backendAddrRest := fmt.Sprintf("http://%s:%s", backendHost, backendRestPort)
	if tlsEnabled {
		backendAddrRest = fmt.Sprintf("https://%s:%s", backendHost, backendRestPort)
	}

	// 1. Setup Outbound Clients to Backend
	var grpcDialOpts []grpc.DialOption
	var httpClient *http.Client

	if tlsEnabled {
		// Use verified TLS config explicitly trusted by our Local CA
		tlsConfig, err := tlsutil.LoadClientConfig(certDir + "/ca.crt")
		if err != nil {
			log.Fatalf("failed to load client CA cert: %v", err)
		}
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
			Timeout: 5 * time.Second,
		}
	} else {
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	backendConn, err := grpc.NewClient(backendAddr, grpcDialOpts...)
	if err != nil {
		log.Fatalf("failed to connect to backend at %s: %v", backendAddr, err)
	}
	defer backendConn.Close()

	grpcBackend := &gateway.GrpcBackendClient{Client: pb.NewPingServiceClient(backendConn)}
	restBackend := &gateway.RestBackendClient{TargetURL: backendAddrRest, HTTPClient: httpClient}

	// 2. Determine Primary Backend for gRPC proxy based on PROTOCOL env var defaults to gRPC if unset
	var primaryBackend gateway.BackendClient = grpcBackend
	if protocol == "rest" {
		primaryBackend = restBackend
	}

	// 3. Setup Gateway Inbound Listeners
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcServerOpts []grpc.ServerOption
	httpServer := &http.Server{
		Addr:    ":" + restPort,
		Handler: gateway.NewRESTServer(protocol, tlsEnabled, grpcBackend, restBackend),
	}

	if tlsEnabled {
		tlsConfig, err := tlsutil.LoadServerConfig(certDir+"/gateway.crt", certDir+"/gateway.key")
		if err != nil {
			log.Fatalf("failed to load Gateway TLS cert: %v", err)
		}
		grpcServerOpts = append(grpcServerOpts, grpc.Creds(credentials.NewTLS(tlsConfig)))
		httpServer.TLSConfig = tlsConfig
	}

	srv := grpc.NewServer(grpcServerOpts...)
	pb.RegisterPingServiceServer(srv, &gateway.PingProxy{
		Backend: primaryBackend,
	})
	reflection.Register(srv)

	go func() {
		fmt.Printf("gateway REST listening on :%s (TLS=%v)\n", restPort, tlsEnabled)
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
		fmt.Fprintln(os.Stderr, "shutting down gateway...")
		srv.GracefulStop()
		httpServer.Close()
		backendConn.Close()
	}()

	fmt.Printf("gateway gRPC listening on :%s (forwarding to %s, PROTOCOL=%s, TLS=%v)\n",
		port, backendAddr, protocol, tlsEnabled)
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
