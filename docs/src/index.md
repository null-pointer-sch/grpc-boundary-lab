# gRPC Boundary Lab

This project explores the performance impact of introducing a gateway boundary
in front of a gRPC backend service.

The setup:

- Backend (port 50051 for gRPC, 8081 for REST)
- Gateway (port 50052 for gRPC, 8080 for REST → forwards to backend)
- Frontend (React Dashboard UI on port 80)
- Load Generator with configurable:
  - protocol (gRPC vs REST)
  - concurrency
  - warmup
  - deadlines
  - repeated runs
  - percentile reporting (p50 / p95 / p99)

The goal is to quantify:

- Throughput degradation
- Tail latency amplification
- Saturation behavior under load
- Baseline comparison between HTTP/2 multiplexing (gRPC) vs HTTP/1.1 Keep-Alive (REST)

---

## Quick Start

The recommended way to run the entire stack (Backend, Gateway, and Dashboard) locally is via Docker Compose:

```bash
make all
```

Once running, open your browser to `http://localhost/` to view the interactive performance dashboard. You can toggle between tracking `gRPC` proxy latency and standard `REST` proxies.

### Manual CLI Testing

If you prefer to run isolated services without Docker:

```bash
# Terminal A
make backend

# Terminal B
make gateway
```

### Run load against backend

```bash
make loadgen-backend
```

### Run full sweep

```bash
make sweep RUNS=5
```
