# Load Generator

The load generator is a high-performance Go client designed to saturate the system and measure stable state performance.

## Design

- **Goroutine-Per-Worker**: Uses a configurable number of goroutines to maintain constant concurrency.
- **Blocking Stubs**: Uses blocking gRPC calls per goroutine for simple end-to-end latency measurement.
- **HdrHistogram**: Uses a High Dynamic Range Histogram to capture latency percentiles (p50, p95, p99) with high precision and low overhead.

## Metrics Reported

| Metric | Description |
| :--- | :--- |
| **ok_rps** | Successful requests per second. |
| **p50_us** | Median latency in microseconds. |
| **p99_us** | Tail latency (99th percentile) in microseconds. |
| **errors** | Number of failed requests (e.g., Deadline Exceeded). |

## Usage

The generator is typically run via the root `Makefile`:
```bash
make sweep REQUESTS=50000 CONCURRENCY="1 8 16"
```
