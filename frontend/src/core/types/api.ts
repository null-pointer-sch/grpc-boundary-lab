export interface ModeInfo {
    protocol: 'rest' | 'grpc';
    tls: boolean;
    tlsAvailable: boolean;
}

export interface PingResponse {
    message: string;
    latencyMs: number;
}

export interface BenchmarkStats {
    protocol: string;
    tls: boolean;
    rps: number;
    p50: number;
    p95: number;
    p99: number;
    timestamp?: string; // Optional if you decide to track when it was recorded
}
