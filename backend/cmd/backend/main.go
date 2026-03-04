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
	port := envutil.GetOrDefault("BACKEND_PORT", "50051")
	portTLS := envutil.GetOrDefault("BACKEND_PORT_TLS", "50151")
	restPort := envutil.GetOrDefault("BACKEND_REST_PORT", "8081")
	restPortTLS := envutil.GetOrDefault("BACKEND_REST_PORT_TLS", "8181")
	certDir := envutil.GetOrDefault("CERT_DIR", "/certs")

	restMux := http.NewServeMux()
	restHandler := &service.RESTPingHandler{}
	restMux.HandleFunc("/api/ping", restHandler.HandlePing)

	// ── 1. Always start plaintext gRPC + REST ────────────────────────────

	plainLis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("backend: failed to listen on :%s: %v", port, err)
	}

	plainGRPC := grpc.NewServer()
	pb.RegisterPingServiceServer(plainGRPC, &service.PingServer{})
	reflection.Register(plainGRPC)

	plainHTTP := &http.Server{Addr: ":" + restPort, Handler: restMux}

	go func() {
		fmt.Printf("backend gRPC listening on :%s (plain)\n", port)
		if err := plainGRPC.Serve(plainLis); err != nil {
			log.Fatalf("backend: failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		fmt.Printf("backend REST listening on :%s (plain)\n", restPort)
		if err := plainHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("backend: failed to serve REST: %v", err)
		}
	}()

	// ── 2. Optionally start TLS gRPC + REST (if certs present) ───────────

	tlsConfig, loadErr := tlsutil.LoadServerConfig(certDir+"/backend.crt", certDir+"/backend.key")
	if loadErr == nil {
		tlsLis, err := net.Listen("tcp", ":"+portTLS)
		if err != nil {
			log.Fatalf("backend: failed to listen on TLS port :%s: %v", portTLS, err)
		}

		tlsGRPC := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
		pb.RegisterPingServiceServer(tlsGRPC, &service.PingServer{})
		reflection.Register(tlsGRPC)

		tlsHTTP := &http.Server{Addr: ":" + restPortTLS, Handler: restMux, TLSConfig: tlsConfig}

		go func() {
			fmt.Printf("backend gRPC listening on :%s (TLS)\n", portTLS)
			if err := tlsGRPC.Serve(tlsLis); err != nil {
				log.Fatalf("backend: failed to serve TLS gRPC: %v", err)
			}
		}()

		go func() {
			fmt.Printf("backend REST listening on :%s (TLS)\n", restPortTLS)
			if err := tlsHTTP.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatalf("backend: failed to serve TLS REST: %v", err)
			}
		}()

		log.Printf("backend TLS listeners active (certs from %s)", certDir)
	} else {
		log.Printf("backend TLS not available: %v", loadErr)
	}

	// ── 3. Graceful shutdown ─────────────────────────────────────────────

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Fprintln(os.Stderr, "shutting down backend...")
	plainGRPC.GracefulStop()
	plainHTTP.Close()
}
