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

	grpcClient, httpClient := initClients(mode, addr, tlsEnabled, certDir, c, deadlineMs)

	scheme := "http"
	if tlsEnabled {
		scheme = "https"
	}

	cfg := PhaseConfig{Mode: mode, Scheme: scheme, Target: addr, C: c, DeadlineMs: deadlineMs}

	// Warmup
	if warmup > 0 {
		fmt.Printf("warmup: %d requests overhead (%s, tls=%v) with concurrency=%d\n", warmup, mode, tlsEnabled, c)
		wCfg := cfg
		wCfg.N = warmup
		runPhase(grpcClient, httpClient, wCfg, nil)
	}

	fmt.Println("run,mode,attempted,ok,errors,concurrency,seconds,ok_rps,p50_us,p95_us,p99_us,max_us")

	var minRps, maxRps, sumRps float64
	minRps = math.Inf(1)

	cfg.N = n

	for r := 1; r <= runs; r++ {
		executeRun(r, cfg, grpcClient, httpClient, &minRps, &maxRps, &sumRps)
	}

	if runs > 1 {
		fmt.Printf("ok_rps summary: avg=%.2f min=%.2f max=%.2f\n",
			sumRps/float64(runs), minRps, maxRps)
	}
}

func executeRun(r int, cfg PhaseConfig, grpcClient pb.PingServiceClient, httpClient *http.Client, minRps, maxRps, sumRps *float64) {
	hist := hdrhistogram.New(1, 60_000_000, 3) // micros
	cfg.PrintExample = (r == 1)

	t0 := time.Now()
	ok, errors := runPhase(grpcClient, httpClient, cfg, hist)
	elapsed := time.Since(t0).Seconds()

	okRps := float64(ok) / elapsed

	fmt.Printf("%d,%s,%d,%d,%d,%d,%.3f,%.2f,%d,%d,%d,%d\n",
		r, cfg.Mode, cfg.N, ok, errors, cfg.C, elapsed, okRps,
		hist.ValueAtQuantile(50.0),
		hist.ValueAtQuantile(95.0),
		hist.ValueAtQuantile(99.0),
		hist.Max(),
	)

	if okRps < *minRps {
		*minRps = okRps
	}
	if okRps > *maxRps {
		*maxRps = okRps
	}
	*sumRps += okRps
}

type PhaseConfig struct {
	Mode         string
	Scheme       string
	Target       string
	N            int
	C            int
	DeadlineMs   int64
	PrintExample bool
}

func initClients(mode, addr string, tlsEnabled bool, certDir string, c int, deadlineMs int64) (pb.PingServiceClient, *http.Client) {
	var grpcClient pb.PingServiceClient
	var httpClient *http.Client

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
		// Notice: conn is deliberately not closed here for the lifespan of loadgen
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
	return grpcClient, httpClient
}

type WorkerState struct {
	ok            *atomic.Int64
	errors        *atomic.Int64
	printedErrors *atomic.Int32
	wg            *sync.WaitGroup
	hist          *hdrhistogram.Histogram
}

func runPhase(grpcClient pb.PingServiceClient, httpClient *http.Client, cfg PhaseConfig, hist *hdrhistogram.Histogram) (int64, int64) {
	var wg sync.WaitGroup
	var ok, errors atomic.Int64
	var printedErrors atomic.Int32

	ws := &WorkerState{
		ok:            &ok,
		errors:        &errors,
		printedErrors: &printedErrors,
		wg:            &wg,
		hist:          hist,
	}

	for worker := 0; worker < cfg.C; worker++ {
		wg.Add(1)
		go runWorker(worker, grpcClient, httpClient, cfg, ws)
	}

	wg.Wait()
	return ok.Load(), errors.Load()
}

func runWorker(workerID int, grpcClient pb.PingServiceClient, httpClient *http.Client, cfg PhaseConfig, ws *WorkerState) {
	defer ws.wg.Done()

	base := cfg.N / cfg.C
	extra := cfg.N % cfg.C
	myN := base
	if workerID < extra {
		myN++
	}
	startIndex := workerID*base + min(workerID, extra)

	for j := 0; j < myN; j++ {
		i := startIndex + j
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DeadlineMs)*time.Millisecond)

		startNs := time.Now()
		val, err := executeRequest(ctx, grpcClient, httpClient, cfg, i)
		cancel()

		if err != nil {
			ws.errors.Add(1)
			if ws.printedErrors.Add(1) <= 3 {
				fmt.Printf("error example: %v\n", err)
			}
			continue
		}

		ws.ok.Add(1)
		if ws.hist != nil {
			durUs := time.Since(startNs).Microseconds()
			if durUs < 1 {
				durUs = 1
			}
			ws.hist.RecordValue(durUs)

			if cfg.PrintExample && workerID == 0 && j == 0 {
				fmt.Printf("example response: %s\n", val)
			}
		}
	}
}

func executeRequest(ctx context.Context, grpcClient pb.PingServiceClient, httpClient *http.Client, cfg PhaseConfig, i int) (string, error) {
	if cfg.Mode == "grpc" {
		resp, err := grpcClient.Ping(ctx, &pb.PingRequest{Message: fmt.Sprintf("hi %d", i)})
		if err != nil {
			return "", err
		}
		return resp.GetMessage(), nil
	}

	url := fmt.Sprintf("%s://%s/api/ping?message=hi%%20%d", cfg.Scheme, cfg.Target, i)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status: %d", resp.StatusCode)
	}
	return "pong (rest)", nil
}
