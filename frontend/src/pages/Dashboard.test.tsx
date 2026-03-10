import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Dashboard } from './Dashboard';
import { useDashboard } from '../hooks/useDashboard';

vi.mock('../hooks/useDashboard');

describe('Dashboard Component', () => {
    it('renders loading state', () => {
        vi.mocked(useDashboard).mockReturnValue({
            loading: true,
            error: null,
            mode: null,
            stats: null,
            pingResult: null,
            protocol: 'grpc',
            tlsEnabled: false,
            fetchingStats: false,
            pinging: false,
            setProtocol: vi.fn(),
            setTlsEnabled: vi.fn(),
            refresh: vi.fn(),
            handlePing: vi.fn(),
        });

        render(<Dashboard />);
        expect(screen.getByText(/Initializing.../i)).toBeInTheDocument();
    });

    it('renders error state', () => {
        vi.mocked(useDashboard).mockReturnValue({
            loading: false,
            error: 'Failed to connect',
            mode: null,
            stats: null,
            pingResult: null,
            protocol: 'grpc',
            tlsEnabled: false,
            fetchingStats: false,
            pinging: false,
            setProtocol: vi.fn(),
            setTlsEnabled: vi.fn(),
            refresh: vi.fn(),
            handlePing: vi.fn(),
        });

        render(<Dashboard />);
        expect(screen.getByText(/Failed to connect/)).toBeInTheDocument();
    });

    it('renders success state and handles interactions', () => {
        const setProtocol = vi.fn();
        const setTlsEnabled = vi.fn();
        const refresh = vi.fn();
        const handlePing = vi.fn();

        vi.mocked(useDashboard).mockReturnValue({
            loading: false,
            error: null,
            mode: { protocol: 'grpc', tls: false, tlsAvailable: true },
            stats: { protocol: 'grpc', tls: false, rps: 1000, p50: 1, p95: 2, p99: 3 },
            pingResult: { message: 'pong', latencyMs: 42 },
            protocol: 'grpc',
            tlsEnabled: false,
            fetchingStats: false,
            pinging: false,
            setProtocol,
            setTlsEnabled,
            refresh,
            handlePing,
        });

        render(<Dashboard />);
        
        // Protocol buttons
        fireEvent.click(screen.getByText('REST'));
        expect(setProtocol).toHaveBeenCalledWith('rest');

        fireEvent.click(screen.getByText('gRPC'));
        expect(setProtocol).toHaveBeenCalledWith('grpc');

        // Security buttons
        fireEvent.click(screen.getByText(/Plain/i));
        expect(setTlsEnabled).toHaveBeenCalledWith(false);

        fireEvent.click(screen.getByText(/mTLS/i));
        expect(setTlsEnabled).toHaveBeenCalledWith(true);

        // Actions
        fireEvent.click(screen.getByText('Refresh Data'));
        expect(refresh).toHaveBeenCalled();

        fireEvent.click(screen.getByText('Initialize Probe'));
        expect(handlePing).toHaveBeenCalled();
    });
});
