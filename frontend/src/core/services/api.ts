import axios from 'axios';
import type { ModeInfo, PingResponse, BenchmarkStats } from '../types/api';

// Create an Axios instance with base configuration
export const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const api = {
    // Get current mode (REST/gRPC + TLS)
    getMode: async (target: 'grpc' | 'rest' = 'grpc', tls: boolean = false): Promise<ModeInfo> => {
        const params = new URLSearchParams({ target, tls: String(tls) });
        const { data } = await apiClient.get<ModeInfo>(`/mode?${params.toString()}`);
        return data;
    },

    // Ping endpoint to test latency
    ping: async (target: 'grpc' | 'rest' = 'grpc', tls: boolean = false): Promise<PingResponse> => {
        const params = new URLSearchParams({ target, tls: String(tls) });
        const { data } = await apiClient.get<PingResponse>(`/ping?${params.toString()}`);
        return data;
    },

    // Get latest benchmark run statistics
    getLatestBench: async (target: 'grpc' | 'rest' = 'grpc', tls: boolean = false): Promise<BenchmarkStats> => {
        const params = new URLSearchParams({ target, tls: String(tls) });
        const { data } = await apiClient.get<BenchmarkStats>(`/bench/latest?${params.toString()}`);
        return data;
    },
};
