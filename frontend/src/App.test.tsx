import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import App from './App';

// Mock the dashboard hook to provide stable data for the UI to render
vi.mock('./hooks/useDashboard', () => ({
    useDashboard: () => ({
        mode: { protocol: 'grpc', tls: false, tlsAvailable: true },
        stats: { protocol: 'grpc', tls: false, rps: 23561, p50: 1.23, p95: 2.11, p99: 2.46 },
        pingResult: { message: 'pong', latencyMs: 1.5 },
        protocol: 'grpc',
        tlsEnabled: false,
        loading: false,
        fetchingStats: false,
        pinging: false,
        error: null,
        setProtocol: vi.fn(),
        setTlsEnabled: vi.fn(),
        refresh: vi.fn(),
        handlePing: vi.fn()
    })
}));

describe('App Integration', () => {
    it('renders without crashing and displays core components', () => {
        render(<App />);
        
        // Check AppShell exists
        expect(screen.getByText(/Enterprise Performance Labs/i)).toBeInTheDocument();

        // Check TopBar exists
        expect(screen.getByPlaceholderText(/Search/i)).toBeInTheDocument();

        // Stats Cards are visible
        expect(screen.getByText(/23,561/)).toBeInTheDocument(); // formatted RPS
        expect(screen.getByText(/1.23/)).toBeInTheDocument(); // P50
        
    });
});
