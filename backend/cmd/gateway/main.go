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
	cfg := envutil.LoadConfig()

	protocol := os.Getenv("PROTOCOL")

	backendAddr := fmt.Sprintf("%s:%s", cfg.BackendHost, cfg.BackendPort)
	backendAddrTLS := fmt.Sprintf("%s:%s", cfg.BackendHost, cfg.BackendPortTLS)
	backendAddrRest := fmt.Sprintf("http://%s:%s", cfg.BackendHost, cfg.BackendRESTPort)
	backendAddrRestTLS := fmt.Sprintf("https://%s:%s", cfg.BackendHost, cfg.BackendRESTPortTLS)

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

	clientTLS, loadErr := tlsutil.LoadClientConfig(cfg.CertDir + "/ca.crt")
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
		log.Printf("TLS clients initialised (certs from %s)", cfg.CertDir)
	} else {
		log.Printf("TLS clients not available (certs not found at %s): %v", cfg.CertDir, loadErr)
	}

	// ── 3. Setup Gateway Inbound Listeners ───────────────────────────────

	lis, err := net.Listen("tcp", ":"+cfg.GatewayPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcServerOpts []grpc.ServerOption

	restServer := gateway.NewRESTServer(protocol, grpcPlain, restPlain, grpcTLS, restTLS)
	httpServer := &http.Server{
		Addr:    ":" + cfg.GatewayRESTPort,
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
		fmt.Printf("gateway REST listening on :%s (TLS-toggle=%v)\n", cfg.GatewayRESTPort, grpcTLS != nil)
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
		cfg.GatewayPort, backendAddr, protocol)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
