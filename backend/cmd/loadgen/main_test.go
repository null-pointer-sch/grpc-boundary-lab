package main

import (
	"os"
	"testing"
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
