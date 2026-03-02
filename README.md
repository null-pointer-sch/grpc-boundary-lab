# grpc-boundary-lab

A small lab to measure the performance cost of introducing a gRPC gateway boundary in front of a backend service.

It compares two request paths:

- Direct: Client → Backend
- With gateway: Client → Gateway → Backend

The goal is to quantify:

- Throughput degradation
- Tail-latency amplification (p95 / p99)
- Saturation behavior under load

Live docs:
https://AndySchubert.github.io/grpc-boundary-lab/

---

## What's in this repo

- cmd/backend/ – gRPC backend service
- cmd/gateway/ – gRPC gateway that forwards to backend
- cmd/loadgen/ – load generator with percentile latency tracking (HdrHistogram)
- gateway/ – gateway proxy logic (testable)
- internal/proto/ – protobuf definitions and generated code
- docs/ + mkdocs.yml – documentation site
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

Build everything:

    make build

Run tests:

    make test

Start services (two terminals):

Terminal A:

    make backend

Terminal B:

    make gateway

---

## Run a load sweep

    make sweep REQUESTS=50000 CONCURRENCY="1 16 64"

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

---

## Documentation

Online:
https://AndySchubert.github.io/grpc-boundary-lab/

Local:

    make docs

Then open the local MkDocs URL printed in the terminal.

---

## Make targets

- make build
- make test
- make vet
- make backend
- make gateway
- make sweep
- make docs

---

## Status

- Gateway forwarding via gRPC client
- Automated load generator with percentile latency tracking
- Integration tests and CI
- MkDocs-based documentation site

---

## License

MIT
