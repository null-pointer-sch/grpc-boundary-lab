package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMainExecution(t *testing.T) {
	// Set ports to 0 to dynamically assign and avoid conflicts
	os.Setenv("BACKEND_PORT", "0")
	os.Setenv("BACKEND_PORT_TLS", "0")
	os.Setenv("BACKEND_REST_PORT", "0")
	os.Setenv("BACKEND_REST_PORT_TLS", "0")
	os.Setenv("CERT_DIR", "/tmp/nonexistent_certs_to_skip_tls")

	// Run main asynchronously
	go func() {
		main()
	}()

	// Give the server time to start up
	time.Sleep(200 * time.Millisecond)

	// Send interrupt to trigger GracefulStop
	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err == nil {
		process.Signal(syscall.SIGINT)
	}

	// Give it time to shut down
	time.Sleep(100 * time.Millisecond)
}
