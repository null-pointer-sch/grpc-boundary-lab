import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useDashboard } from './useDashboard';
import { api } from '../core/services/api';
import type { ModeInfo, BenchmarkStats, PingResponse } from '../core/types/api';

vi.mock('../core/services/api', () => ({
    api: {
        getMode: vi.fn(),
        getLatestBench: vi.fn(),
        ping: vi.fn(),
    }
}));

describe('useDashboard hook', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('initializes and fetches data on mount', async () => {
        const mockMode: ModeInfo = { protocol: 'grpc', tls: false, tlsAvailable: true };
        const mockStats: BenchmarkStats = { protocol: 'grpc', tls: false, rps: 100, p50: 1, p95: 2, p99: 3 };
        
        vi.mocked(api.getMode).mockResolvedValueOnce(mockMode);
        vi.mocked(api.getLatestBench).mockResolvedValueOnce(mockStats);

        const { result } = renderHook(() => useDashboard());

        // Initial state
        expect(result.current.loading).toBe(true);

        // Wait for state updates
        await act(async () => {
            await new Promise(resolve => setTimeout(resolve, 0));
        });

        expect(api.getMode).toHaveBeenCalledWith('grpc', false);
        expect(api.getLatestBench).toHaveBeenCalledWith('grpc', false);
        expect(result.current.loading).toBe(false);
        expect(result.current.mode).toEqual(mockMode);
        expect(result.current.stats).toEqual(mockStats);
        expect(result.current.error).toBeNull();
    });

    it('handles fetch errors on mount', async () => {
        vi.mocked(api.getMode).mockRejectedValueOnce(new Error('Network Error'));

        const { result } = renderHook(() => useDashboard());

        await act(async () => {
            await new Promise(resolve => setTimeout(resolve, 0));
        });

        expect(result.current.loading).toBe(false);
        expect(result.current.error).toBe('Could not connect to gateway');
    });

    it('handles ping correctly', async () => {
        const mockPing: PingResponse = { message: 'pong', latencyMs: 5 };
        vi.mocked(api.ping).mockResolvedValueOnce(mockPing);
        
        // Mock the initial fetch so it behaves normally
        vi.mocked(api.getMode).mockResolvedValueOnce({ protocol: 'grpc', tls: false, tlsAvailable: true });
        vi.mocked(api.getLatestBench).mockResolvedValueOnce({ protocol: 'grpc', tls: false, rps: 100, p50: 1, p95: 2, p99: 3 });

        const { result } = renderHook(() => useDashboard());

        await act(async () => {
            await new Promise(resolve => setTimeout(resolve, 0));
            await result.current.handlePing();
        });

        expect(api.ping).toHaveBeenCalledWith('grpc', false);
        expect(result.current.pingResult).toEqual(mockPing);
        expect(result.current.pinging).toBe(false);
    });

    it('updates protocol triggers fetch immediately', async () => {
        vi.mocked(api.getMode).mockResolvedValue({ protocol: 'rest', tls: false, tlsAvailable: true });
        vi.mocked(api.getLatestBench).mockResolvedValue({ protocol: 'rest', tls: false, rps: 100, p50: 1, p95: 2, p99: 3 });

        const { result } = renderHook(() => useDashboard());

        await act(async () => {
            result.current.setProtocol('rest');
        });

        // The hook triggers fetchAll(p, tls, false)
        expect(api.getMode).toHaveBeenCalledWith('rest', false);
        expect(result.current.protocol).toBe('rest');
    });

    it('updates tls triggers fetch immediately', async () => {
        vi.mocked(api.getMode).mockResolvedValue({ protocol: 'grpc', tls: true, tlsAvailable: true });
        vi.mocked(api.getLatestBench).mockResolvedValue({ protocol: 'grpc', tls: true, rps: 100, p50: 1, p95: 2, p99: 3 });

        const { result } = renderHook(() => useDashboard());

        await act(async () => {
            result.current.setTlsEnabled(true);
        });

        expect(api.getMode).toHaveBeenCalledWith('grpc', true);
        expect(result.current.tlsEnabled).toBe(true);
    });
});
