package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/envutil"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/service"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/tlsutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := envutil.LoadConfig()

	restMux := http.NewServeMux()
	restHandler := &service.RESTPingHandler{}
	restMux.HandleFunc("/api/ping", restHandler.HandlePing)

	// ── 1. Always start plaintext gRPC + REST ────────────────────────────

	plainLis, err := net.Listen("tcp", ":"+cfg.BackendPort)
	if err != nil {
		log.Fatalf("backend: failed to listen on :%s: %v", cfg.BackendPort, err)
	}

	plainGRPC := grpc.NewServer()
	pb.RegisterPingServiceServer(plainGRPC, &service.PingServer{})
	reflection.Register(plainGRPC)

	plainHTTP := &http.Server{Addr: ":" + cfg.BackendRESTPort, Handler: restMux}

	go func() {
		fmt.Printf("backend gRPC listening on :%s (plain)\n", cfg.BackendPort)
		if err := plainGRPC.Serve(plainLis); err != nil {
			log.Printf("backend: failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		fmt.Printf("backend REST listening on :%s (plain)\n", cfg.BackendRESTPort)
		if err := plainHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("backend: failed to serve REST: %v", err)
		}
	}()

	// ── 2. Optionally start TLS gRPC + REST (if certs present) ───────────
	tlsGRPC, tlsHTTP := startTLSServers(cfg, restMux)

	// ── 3. Graceful shutdown ─────────────────────────────────────────────

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Fprintln(os.Stderr, "shutting down backend...")
	plainGRPC.GracefulStop()
	plainHTTP.Close()
	if tlsGRPC != nil {
		tlsGRPC.GracefulStop()
	}
	if tlsHTTP != nil {
		tlsHTTP.Close()
	}
}

func startTLSServers(cfg *envutil.Config, restMux *http.ServeMux) (*grpc.Server, *http.Server) {
	tlsConfig, loadErr := tlsutil.LoadServerConfig(cfg.CertDir+"/backend.crt", cfg.CertDir+"/backend.key")
	if loadErr != nil {
		log.Printf("backend TLS not available: %v", loadErr)
		return nil, nil
	}

	tlsLis, err := net.Listen("tcp", ":"+cfg.BackendPortTLS)
	if err != nil {
		log.Fatalf("backend: failed to listen on TLS port :%s: %v", cfg.BackendPortTLS, err)
	}

	tlsGRPC := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	pb.RegisterPingServiceServer(tlsGRPC, &service.PingServer{})
	reflection.Register(tlsGRPC)

	tlsHTTP := &http.Server{Addr: ":" + cfg.BackendRESTPortTLS, Handler: restMux, TLSConfig: tlsConfig}

	go func() {
		fmt.Printf("backend gRPC listening on :%s (TLS)\n", cfg.BackendPortTLS)
		if err := tlsGRPC.Serve(tlsLis); err != nil {
			log.Printf("backend: failed to serve TLS gRPC: %v", err)
		}
	}()

	go func() {
		fmt.Printf("backend REST listening on :%s (TLS)\n", cfg.BackendRESTPortTLS)
		if err := tlsHTTP.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Printf("backend: failed to serve TLS REST: %v", err)
		}
	}()

	log.Printf("backend TLS listeners active (certs from %s)", cfg.CertDir)
	return tlsGRPC, tlsHTTP
}
