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

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/envutil"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/gateway"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/tlsutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := envutil.GetOrDefault("GATEWAY_PORT", "50052")
	restPort := envutil.GetOrDefault("REST_PORT", "8080")
	backendHost := envutil.GetOrDefault("BACKEND_HOST", "127.0.0.1")
	backendPort := envutil.GetOrDefault("BACKEND_PORT", "50051")
	backendPortTLS := envutil.GetOrDefault("BACKEND_PORT_TLS", "50151")
	backendRestPort := envutil.GetOrDefault("BACKEND_REST_PORT", "8081")
	backendRestPortTLS := envutil.GetOrDefault("BACKEND_REST_PORT_TLS", "8181")

	protocol := os.Getenv("PROTOCOL")
	certDir := envutil.GetOrDefault("CERT_DIR", "/certs")

	backendAddr := fmt.Sprintf("%s:%s", backendHost, backendPort)
	backendAddrTLS := fmt.Sprintf("%s:%s", backendHost, backendPortTLS)
	backendAddrRest := fmt.Sprintf("http://%s:%s", backendHost, backendRestPort)
	backendAddrRestTLS := fmt.Sprintf("https://%s:%s", backendHost, backendRestPortTLS)

	// ── 1. Always create plaintext clients ───────────────────────────────

	plainConn, err := grpc.NewClient(backendAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect (plain) to backend at %s: %v", backendAddr, err)
	}
	defer plainConn.Close()

	grpcPlain := &gateway.GrpcBackendClient{Client: pb.NewPingServiceClient(plainConn)}
	restPlain := &gateway.RestBackendClient{
		TargetURL:  backendAddrRest,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}

	// ── 2. Optionally create TLS clients (if certs are present) ──────────

	var grpcTLS gateway.BackendClient
	var restTLS gateway.BackendClient

	clientTLS, loadErr := tlsutil.LoadClientConfig(certDir + "/ca.crt")
	if loadErr == nil {
		tlsConn, dialErr := grpc.NewClient(backendAddrTLS,
			grpc.WithTransportCredentials(credentials.NewTLS(clientTLS)))
		if dialErr == nil {
			defer tlsConn.Close()
			grpcTLS = &gateway.GrpcBackendClient{Client: pb.NewPingServiceClient(tlsConn)}
		} else {
			log.Printf("warning: TLS gRPC dial failed (TLS toggle disabled for gRPC): %v", dialErr)
		}

		restTLS = &gateway.RestBackendClient{
			TargetURL: backendAddrRestTLS,
			HTTPClient: &http.Client{
				Transport: &http.Transport{TLSClientConfig: clientTLS},
				Timeout:   5 * time.Second,
			},
		}
		log.Printf("TLS clients initialised (certs from %s)", certDir)
	} else {
		log.Printf("TLS clients not available (certs not found at %s): %v", certDir, loadErr)
	}

	// ── 3. Setup Gateway Inbound Listeners ───────────────────────────────

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcServerOpts []grpc.ServerOption

	restServer := gateway.NewRESTServer(protocol, grpcPlain, restPlain, grpcTLS, restTLS)
	httpServer := &http.Server{
		Addr:    ":" + restPort,
		Handler: restServer,
	}

	// Determine primary backend for the gRPC proxy
	var primaryBackend gateway.BackendClient = grpcPlain
	if protocol == "rest" {
		primaryBackend = restPlain
	}

	srv := grpc.NewServer(grpcServerOpts...)
	pb.RegisterPingServiceServer(srv, &gateway.PingProxy{
		Backend: primaryBackend,
	})
	reflection.Register(srv)

	go func() {
		fmt.Printf("gateway REST listening on :%s (TLS-toggle=%v)\n", restPort, grpcTLS != nil)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	}()

	fmt.Printf("gateway gRPC listening on :%s (forwarding to %s, PROTOCOL=%s)\n",
		port, backendAddr, protocol)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
