package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMainExecution(t *testing.T) {
	os.Setenv("GATEWAY_PORT", "0")
	os.Setenv("REST_PORT", "0")
	os.Setenv("BACKEND_HOST", "127.0.0.1")
	// Try a bad/empty backend port so client connection fails fast or tries to connect but is ignored
	os.Setenv("BACKEND_PORT", "9999")
	os.Setenv("CERT_DIR", "/tmp/nonexistent")

	go func() {
		main()
	}()

	time.Sleep(200 * time.Millisecond)

	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err == nil {
		process.Signal(syscall.SIGINT)
	}

	time.Sleep(100 * time.Millisecond)
}

func TestMainExecution_TLS(t *testing.T) {
	os.Setenv("GATEWAY_PORT", "0")
	os.Setenv("REST_PORT", "0")
	os.Setenv("BACKEND_HOST", "127.0.0.1")
	os.Setenv("BACKEND_PORT", "9999")
	os.Setenv("CERT_DIR", "/tmp/certs")

	go func() {
		main()
	}()

	time.Sleep(200 * time.Millisecond)

	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err == nil {
		process.Signal(syscall.SIGINT)
	}

	time.Sleep(100 * time.Millisecond)
}
