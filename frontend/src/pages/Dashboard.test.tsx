import { render, screen } from '@testing-library/react';
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
        expect(screen.getByTestId('dashboard-skeleton')).toBeInTheDocument();
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
});
