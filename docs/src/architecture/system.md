# System Overview

The system consists of three components:

## Backend

- **Port**: 50051
- Implements `PingService`
- Handles unary gRPC requests

## Gateway

- **Port**: 50052
- Accepts the same gRPC API
- Forwards requests to backend via a gRPC client
- Introduces an additional network + serialisation boundary

## Load Generator

- Concurrent gRPC client using goroutines
- Supports configurable concurrency
- Reports latency percentiles via HdrHistogram
