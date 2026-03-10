package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

func TestMainExecution_GRPC(t *testing.T) {
	// Let it run a single test and naturally exit
	os.Setenv("REQUESTS", "1")
	os.Setenv("CONCURRENCY", "1")
	os.Setenv("WARMUP", "0")
	os.Setenv("RUNS", "1")
	os.Setenv("TARGET_PORT", "0") // Will error out locally or connect, but will fail fast
	os.Setenv("MODE", "grpc")

	main()
}

func TestMainExecution_REST(t *testing.T) {
	os.Setenv("REQUESTS", "1")
	os.Setenv("CONCURRENCY", "1")
	os.Setenv("WARMUP", "0")
	os.Setenv("RUNS", "1")
	os.Setenv("TARGET_PORT", "0")
	os.Setenv("MODE", "rest")

	main()
}

func TestInitClients(t *testing.T) {
	_, httpC := initClients("rest", "127.0.0.1:8080", false, "", 2, 1000)
	if httpC == nil {
		t.Fatal("expected http client")
	}

	// Since we mock TLS by creating a dummy path that will fail if tried,
	// we will skip the exact grpc tls coverage unless necessary because dial options
	// block the client creation if certificates aren't present.
}

func TestExecuteRequestREST(t *testing.T) {
	// Start a dummy HTTP server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong (rest)"))
	}))
	defer srv.Close()

	// Parse host:port
	addr := srv.Listener.Addr().String()

	cfg := PhaseConfig{
		Mode:       "rest",
		Scheme:     "http",
		Target:     addr,
		N:          1,
		C:          1,
		DeadlineMs: 1000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := srv.Client()

	msg, err := executeRequest(ctx, nil, client, cfg, 0)
	if err != nil || msg != "pong (rest)" {
		t.Fatalf("expected ping success, got %v %s", err, msg)
	}
}

func TestWorkerAndPhase(t *testing.T) {
	// A simple wrapper mock where N=0 to just execute the outer structure
	cfg := PhaseConfig{
		Mode:       "rest",
		Scheme:     "http",
		Target:     "dummy",
		N:          0, // Skip real requests
		C:          1,
		DeadlineMs: 1000,
	}

	ok, errs := runPhase(nil, nil, cfg, nil)
	if ok != 0 || errs != 0 {
		t.Fatalf("expected 0, got %d %d", ok, errs)
	}

	// Now test with N=1 to see it hit an error
	cfg.N = 1
	client := &http.Client{Timeout: 50 * time.Millisecond}
	ok, errs = runPhase(nil, client, cfg, nil)
	if ok != 0 || errs != 1 {
		t.Fatalf("expected 1 error, got %d %d", ok, errs)
	}

	// Now with a legitimate server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong (rest)"))
	}))
	defer srv.Close()

	hist := hdrhistogram.New(1, 60_000_000, 3)

	cfg.Target = srv.Listener.Addr().String()
	cfg.C = 2
	cfg.N = 4 // 2 per worker
	cfg.PrintExample = true

	ok, errs = runPhase(nil, srv.Client(), cfg, hist)
	if ok != 4 || errs != 0 {
		t.Fatalf("expected 4 ok, 0 errs, got %d %d", ok, errs)
	}
}
