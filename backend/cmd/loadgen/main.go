package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/envutil"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/tlsutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := envutil.GetOrDefault("TARGET_HOST", "127.0.0.1")
	port := envutil.GetOrDefault("TARGET_PORT", "50052")
	mode := strings.ToLower(envutil.GetOrDefault("MODE", "grpc"))
	tlsEnabled := envutil.GetOrDefault("TLS", "0") == "1"
	certDir := envutil.GetOrDefault("CERT_DIR", "/certs")
	n := envutil.GetInt("REQUESTS", 100)
	c := envutil.GetInt("CONCURRENCY", 1)
	warmup := envutil.GetInt("WARMUP", 2000)
	deadlineMs := envutil.GetInt64("DEADLINE_MS", 20000)
	runs := envutil.GetInt("RUNS", 1)

	addr := fmt.Sprintf("%s:%s", host, port)

	var grpcClient pb.PingServiceClient
	var httpClient *http.Client
	var scheme string = "http"

	if tlsEnabled {
		scheme = "https"
	}

	if mode == "grpc" {
		var grpcDialOpts []grpc.DialOption
		if tlsEnabled {
			tlsConfig, err := tlsutil.LoadClientConfig(certDir + "/ca.crt")
			if err != nil {
				log.Fatalf("failed to load client CA cert: %v", err)
			}
			grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		} else {
			grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		conn, err := grpc.NewClient(addr, grpcDialOpts...)
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		defer conn.Close()
		grpcClient = pb.NewPingServiceClient(conn)
	} else {
		var transport *http.Transport
		if tlsEnabled {
			tlsConfig, err := tlsutil.LoadClientConfig(certDir + "/ca.crt")
			if err != nil {
				log.Fatalf("failed to load client CA cert: %v", err)
			}
			transport = &http.Transport{
				MaxIdleConns:        1000,
				MaxIdleConnsPerHost: c,
				IdleConnTimeout:     90 * time.Second,
				TLSClientConfig:     tlsConfig,
			}
		} else {
			transport = &http.Transport{
				MaxIdleConns:        1000,
				MaxIdleConnsPerHost: c,
				IdleConnTimeout:     90 * time.Second,
			}
		}
		httpClient = &http.Client{
			Transport: transport,
			Timeout:   time.Duration(deadlineMs) * time.Millisecond,
		}
	}

	// Warmup
	if warmup > 0 {
		fmt.Printf("warmup: %d requests overhead (%s, tls=%v) with concurrency=%d\n", warmup, mode, tlsEnabled, c)
		runPhase(grpcClient, httpClient, mode, scheme, addr, warmup, c, deadlineMs, false, nil)
	}

	fmt.Println("run,mode,attempted,ok,errors,concurrency,seconds,ok_rps,p50_us,p95_us,p99_us,max_us")

	var minRps, maxRps, sumRps float64
	minRps = math.Inf(1)

	for r := 1; r <= runs; r++ {
		hist := hdrhistogram.New(1, 60_000_000, 3) // micros
		printExample := r == 1

		t0 := time.Now()
		ok, errors := runPhase(grpcClient, httpClient, mode, scheme, addr, n, c, deadlineMs, printExample, hist)
		elapsed := time.Since(t0).Seconds()

		okRps := float64(ok) / elapsed

		fmt.Printf("%d,%s,%d,%d,%d,%d,%.3f,%.2f,%d,%d,%d,%d\n",
			r, mode, n, ok, errors, c, elapsed, okRps,
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

func runPhase(grpcClient pb.PingServiceClient, httpClient *http.Client, mode, scheme, target string, n, c int, deadlineMs int64, printExample bool, hist *hdrhistogram.Histogram) (okCount, errCount int64) {
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
				var err error
				var val string

				if mode == "grpc" {
					resp, rErr := grpcClient.Ping(ctx, &pb.PingRequest{Message: fmt.Sprintf("hi %d", i)})
					err = rErr
					if err == nil {
						val = resp.GetMessage()
					}
				} else {
					url := fmt.Sprintf("%s://%s/api/ping?message=hi%%20%d", scheme, target, i)
					req, rErr := http.NewRequestWithContext(ctx, "GET", url, nil)
					if rErr == nil {
						resp, dErr := httpClient.Do(req)
						err = dErr
						if err == nil {
							if resp.StatusCode != 200 {
								err = fmt.Errorf("status: %d", resp.StatusCode)
							}
							resp.Body.Close()
							val = "pong (rest)"
						}
					} else {
						err = rErr
					}
				}

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
						fmt.Printf("example response: %s\n", val)
					}
				}
			}
		}()
	}

	wg.Wait()
	return ok.Load(), errors.Load()
}
