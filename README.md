# grpc-boundary-lab

A small lab to measure the performance cost of introducing a gRPC gateway boundary in front of a backend service.

It compares two request paths:

- Direct: Client → Backend
- With gateway: Client → Gateway → Backend

The goal is to quantify:

- Throughput degradation
- Tail-latency amplification (p95 / p99)
- Saturation behavior under load
- **Protocol Overhead:** Direct comparisons between `gRPC` and `REST` proxies.



## What's in this repo

- `cmd/backend/` – gRPC & REST backend service
- `cmd/gateway/` – Gateway proxy that translates or forwards requests
- `cmd/loadgen/` – Load generator with percentile latency tracking (HdrHistogram) capable of sweeping `gRPC` and `REST`.
- `frontend/` – React/Vite Dashboard UI to visualize metrics directly.
- `gateway/` – Gateway proxy routing logic.
- `internal/proto/` – Protobuf definitions and generated code.
- `docker-compose.yml` – Local network orchestration.
- `docs/` + `mkdocs.yml` – Documentation site.
- Makefile – one-command workflows

---

## Prerequisites

- Go 1.26
- make

Optional:
- Docker (for controlled benchmark environments)
- Python (for docs tooling)

---

## Quick start

The fastest way to test the lab locally is using our master `Makefile` orchestrator. It uses Docker Compose to start the Backend, the Gateway, and a React-based Dashboard UI:

```bash
make all
```

Once running, navigate to `http://localhost/` in your browser to access the interactive dashboard. You can toggle between tracking `gRPC` proxy latency and standard `REST` proxies.

Alternatively, you can run isolated load tests via the CLI:

### CLI Sweeps

```bash
make sweep REQUESTS=50000 CONCURRENCY="1 16 64"
```

Suggested quick iteration:

    REQUESTS=20000 CONCURRENCY="1 2 4 8 16 32 64"

Adjust concurrency levels to observe saturation and latency amplification.

---

## What to look for

When introducing a gateway hop, expect overhead from:

- Additional scheduling and queuing
- Serialization/deserialization
- Extra transport hop
- Goroutine contention

Key indicators:

- Where p99 latency rises sharply
- Where throughput plateaus
- How the gateway shifts the latency knee
- The performance differences between `gRPC` Multiplexing and traditional HTTP Keep-Alive connection limits.

---

## Documentation

Online:
https://null-pointer-sch.github.io/grpc-boundary-lab/

Local:

    make docs

Then open the local MkDocs URL printed in the terminal.

---



## Status

- Gateway forwarding via gRPC client
- Automated load generator with percentile latency tracking
- Integration tests and CI
- MkDocs-based documentation site

---

## License

MIT
