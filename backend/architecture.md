# Backend Architecture

The backend of the `grpc-boundary-lab` is designed as a modular, gRPC-first service with a flexible gateway to explore communication boundaries and performance trade-offs.

## Core Components

### 1. Service Layer (`internal/service`)
- **PingService (`ping.go`)**: Implements the core gRPC logic defined in the `.proto` files.
- **REST Handler (`rest.go`)**: Provides a native HTTP/REST implementation of the same functionality for comparison.

### 2. Protocol Definitions (`internal/proto`)
- Uses Protocol Buffers to define the `PingService` interface. This serves as the single source of truth for the primary communication contract.

### 3. Backend Server (`cmd/backend`)
The backend is a multi-protocol listener that can serve requests over:
- **Plaintext gRPC** (default: 50051)
- **TLS gRPC** (default: 50151)
- **Plaintext REST** (default: 8081)
- **TLS REST** (default: 8181)

It loads certificates from the directory specified by `CERT_DIR` to enable TLS.

### 4. Gateway (`cmd/gateway`)
The gateway acts as an intelligent intermediary. It is the primary entry point for the frontend and benchmarking tools.
- **Protocol Switching**: Can forward requests to the backend using either gRPC or REST, controlled by environment variables or query parameters.
- **Security Toggle**: Supports switching between plaintext and TLS communication with the backend service.
- **Dual Interface**: Exposes its own REST and gRPC endpoints to the outside world.

## Communication Flows

1. **Frontend to Backend**: 
   `Frontend (Browser) --[REST]--> Gateway --[gRPC or REST]--> Backend`
2. **Benchmark to Backend**:
   `Load Generator --[gRPC]--> Gateway --[gRPC or REST]--> Backend`

## Key Patterns
- **Internal Package**: All core logic and utilities are kept in the `internal/` directory to prevent external imports and maintain strict encapsulation.
- **Environment Driven**: Configuration is primarily managed via environment variables (e.g., `PROTOCOL`, `BACKEND_HOST`, `CERT_DIR`).
