package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/AndySchubert/grpc-boundary-lab/internal/proto"
	"github.com/HdrHistogram/hdrhistogram-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := envOrDefault("TARGET_HOST", "127.0.0.1")
	port := envOrDefault("TARGET_PORT", "50052")
	n := envInt("REQUESTS", 100)
	c := envInt("CONCURRENCY", 1)
	warmup := envInt("WARMUP", 2000)
	deadlineMs := envInt64("DEADLINE_MS", 20000)
	runs := envInt("RUNS", 1)

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPingServiceClient(conn)

	// Warmup
	if warmup > 0 {
		fmt.Printf("warmup: %d requests with concurrency=%d\n", warmup, c)
		runPhase(client, warmup, c, deadlineMs, false, nil)
	}

	fmt.Println("run,attempted,ok,errors,concurrency,seconds,ok_rps,p50_us,p95_us,p99_us,max_us")

	var minRps, maxRps, sumRps float64
	minRps = math.Inf(1)

	for r := 1; r <= runs; r++ {
		hist := hdrhistogram.New(1, 60_000_000, 3) // micros
		printExample := r == 1

		t0 := time.Now()
		ok, errors := runPhase(client, n, c, deadlineMs, printExample, hist)
		elapsed := time.Since(t0).Seconds()

		okRps := float64(ok) / elapsed

		fmt.Printf("%d,%d,%d,%d,%d,%.3f,%.2f,%d,%d,%d,%d\n",
			r, n, ok, errors, c, elapsed, okRps,
			hist.ValueAtQuantile(50.0),
			hist.ValueAtQuantile(95.0),
			hist.ValueAtQuantile(99.0),
			hist.Max(),
		)

		if okRps < minRps {
			minRps = okRps
		}
		if okRps > maxRps {
			maxRps = okRps
		}
		sumRps += okRps
	}

	if runs > 1 {
		fmt.Printf("ok_rps summary: avg=%.2f min=%.2f max=%.2f\n",
			sumRps/float64(runs), minRps, maxRps)
	}
}

func runPhase(client pb.PingServiceClient, n, c int, deadlineMs int64, printExample bool, hist *hdrhistogram.Histogram) (okCount, errCount int64) {
	var wg sync.WaitGroup
	var ok, errors atomic.Int64
	var printedErrors atomic.Int32

	for worker := 0; worker < c; worker++ {
		wg.Add(1)
		workerID := worker
		base := n / c
		extra := n % c
		myN := base
		if workerID < extra {
			myN++
		}
		startIndex := workerID*base + min(workerID, extra)

		go func() {
			defer wg.Done()
			for j := 0; j < myN; j++ {
				i := startIndex + j
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(deadlineMs)*time.Millisecond)

				startNs := time.Now()
				resp, err := client.Ping(ctx, &pb.PingRequest{Message: fmt.Sprintf("hi %d", i)})
				cancel()

				if err != nil {
					errors.Add(1)
					if printedErrors.Add(1) <= 3 {
						fmt.Printf("error example: %v\n", err)
					}
					continue
				}

				ok.Add(1)
				if hist != nil {
					durUs := time.Since(startNs).Microseconds()
					if durUs < 1 {
						durUs = 1
					}
					hist.RecordValue(durUs)

					if printExample && workerID == 0 && j == 0 {
						fmt.Printf("example response: %s\n", resp.GetMessage())
					}
				}
			}
		}()
	}

	wg.Wait()
	return ok.Load(), errors.Load()
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid %s: %v", key, err)
		}
		return n
	}
	return fallback
}

func envInt64(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Fatalf("invalid %s: %v", key, err)
		}
		return n
	}
	return fallback
}
