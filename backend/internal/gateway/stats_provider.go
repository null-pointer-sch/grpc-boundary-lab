package gateway

import "sync"

// BenchmarkData represents the metrics for a protocol/security combination.
type BenchmarkData struct {
	Protocol string  `json:"protocol"`
	TLS      bool    `json:"tls"`
	RPS      float64 `json:"rps"`
	P50      float64 `json:"p50"`
	P95      float64 `json:"p95"`
	P99      float64 `json:"p99"`
}

// StatsProvider manages access to benchmark statistics.
type StatsProvider struct {
	mu    sync.RWMutex
	stats map[string]BenchmarkData
}

// NewStatsProvider creates a provider with default (hardcoded) benchmark data.
func NewStatsProvider() *StatsProvider {
	p := &StatsProvider{
		stats: make(map[string]BenchmarkData),
	}
	p.initDefaults()
	return p
}

func (p *StatsProvider) initDefaults() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// REST stats
	p.stats["rest-false"] = BenchmarkData{Protocol: "rest", TLS: false, RPS: 6973.90, P50: 3.74, P95: 9.63, P99: 13.36}
	p.stats["rest-true"] = BenchmarkData{Protocol: "rest", TLS: true, RPS: 2331.81, P50: 12.01, P95: 27.32, P99: 35.71}

	// gRPC stats
	p.stats["grpc-false"] = BenchmarkData{Protocol: "grpc", TLS: false, RPS: 23561.92, P50: 1.23, P95: 2.11, P99: 2.46}
	p.stats["grpc-true"] = BenchmarkData{Protocol: "grpc", TLS: true, RPS: 21606.62, P50: 1.37, P95: 2.28, P99: 2.68}
}

// GetStats returns the benchmark data for the given protocol and TLS state.
func (p *StatsProvider) GetStats(protocol string, tls bool) (BenchmarkData, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	key := protocol
	if tls {
		key += "-true"
	} else {
		key += "-false"
	}

	data, ok := p.stats[key]
	return data, ok
}
