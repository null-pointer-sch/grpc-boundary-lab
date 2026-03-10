#!/usr/bin/env bash
set -e

# Automatically run the 4 benchmarks for the Gateway
export WARMUP=500
export REQUESTS=5000
export CERT_DIR=../certs

echo "Building binaries..."
make build >/dev/null

run_benchmark() {
    local proto=$1
    local tls_val=$2
    local label=$3

    echo ""
    echo "============================================================"
    echo " Benchmarking: Gateway (PROTOCOL=$proto, TLS=$tls_val)"
    echo "============================================================"

    # Start backend
    TLS=$tls_val ./bin/backend >/dev/null 2>&1 &
    BACKEND_PID=$!

    # Start gateway
    PROTOCOL=$proto TLS=$tls_val ./bin/gateway >/dev/null 2>&1 &
    GATEWAY_PID=$!

    sleep 1 # wait for servers to bind

    # Run loadgen using Makefile alias
    make --no-print-directory bench-$label

    # Kill servers
    kill $GATEWAY_PID
    kill $BACKEND_PID
    wait $GATEWAY_PID 2>/dev/null || true
    wait $BACKEND_PID 2>/dev/null || true
    
    return 0
}

# Run the 4 combinations
run_benchmark rest 0 rest
run_benchmark rest 1 rest-tls
run_benchmark grpc 0 grpc
run_benchmark grpc 1 grpc-tls

echo ""
echo "Done."
exit 0
