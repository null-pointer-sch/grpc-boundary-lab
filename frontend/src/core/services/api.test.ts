import { describe, it, expect, vi, beforeEach } from 'vitest';
import { api, apiClient } from './api';
import type { ModeInfo, PingResponse, BenchmarkStats } from '../types/api';

describe('api service', () => {
    beforeEach(() => {
        vi.restoreAllMocks();
    });

    it('getMode requests the correct endpoint', async () => {
        const mockData: ModeInfo = { protocol: 'grpc', tls: false, tlsAvailable: true };
        vi.spyOn(apiClient, 'get').mockResolvedValueOnce({ data: mockData });

        const result = await api.getMode('grpc', false);
        expect(apiClient.get).toHaveBeenCalledWith('/mode?target=grpc&tls=false');
        expect(result).toEqual(mockData);
    });

    it('ping requests the correct endpoint', async () => {
        const mockData: PingResponse = { message: 'pong', latencyMs: 5 };
        vi.spyOn(apiClient, 'get').mockResolvedValueOnce({ data: mockData });

        const result = await api.ping('rest', true);
        expect(apiClient.get).toHaveBeenCalledWith('/ping?target=rest&tls=true');
        expect(result).toEqual(mockData);
    });

    it('getLatestBench requests the correct endpoint', async () => {
        const mockData: BenchmarkStats = { protocol: 'grpc', tls: true, rps: 100, p50: 1, p95: 2, p99: 3 };
        vi.spyOn(apiClient, 'get').mockResolvedValueOnce({ data: mockData });

        const result = await api.getLatestBench('grpc', true);
        expect(apiClient.get).toHaveBeenCalledWith('/bench/latest?target=grpc&tls=true');
        expect(result).toEqual(mockData);
    });
});
