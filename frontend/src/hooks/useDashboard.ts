import { useState, useEffect, useCallback, useRef } from 'react';
import { api } from '../core/services/api';
import type { ModeInfo, BenchmarkStats, PingResponse } from '../core/types/api';

interface DashboardState {
    mode: ModeInfo | null;
    stats: BenchmarkStats | null;
    pingResult: PingResponse | null;
    protocol: 'grpc' | 'rest';
    tlsEnabled: boolean;
    loading: boolean;
    fetchingStats: boolean;
    pinging: boolean;
    error: string | null;
}

interface DashboardActions {
    setProtocol: (p: 'grpc' | 'rest') => void;
    setTlsEnabled: (v: boolean) => void;
    refresh: () => void;
    handlePing: () => Promise<void>;
}

export function useDashboard(): DashboardState & DashboardActions {
    const [mode, setMode] = useState<ModeInfo | null>(null);
    const [stats, setStats] = useState<BenchmarkStats | null>(null);
    const [pingResult, setPingResult] = useState<PingResponse | null>(null);

    // Store active toggles in refs so we can use them instantly for fetches without waiting for re-renders.
    const protocolRef = useRef<'grpc' | 'rest'>('grpc');
    const tlsRef = useRef<boolean>(false);

    // We still keep them in React state simply to drive the UI highlights on the buttons.
    const [protocol, _setProtocol] = useState<'grpc' | 'rest'>('grpc');
    const [tlsEnabled, _setTlsEnabled] = useState(false);

    const [loading, setLoading] = useState(true);
    const [fetchingStats, setFetchingStats] = useState(false);
    const [pinging, setPinging] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchIdRef = useRef(0);
    const initialLoadDone = useRef(false);

    const fetchAll = useCallback(async (target: 'grpc' | 'rest', tls: boolean, showSpinner: boolean) => {
        const id = ++fetchIdRef.current;
        if (showSpinner) setLoading(true);
        setFetchingStats(true);

        try {
            const [modeData, statsData] = await Promise.all([
                api.getMode(target, tls),
                api.getLatestBench(target, tls),
            ]);
            if (fetchIdRef.current !== id) return;
            setMode(modeData);
            setStats(statsData);
            setError(null);
        } catch (err: any) {
            if (fetchIdRef.current !== id) return;
            console.error('Fetch failed:', err);
            if (!initialLoadDone.current) setError('Could not connect to gateway');
        } finally {
            if (fetchIdRef.current === id) {
                setLoading(false);
                setFetchingStats(false);
                initialLoadDone.current = true;
            }
        }
    }, []);

    // Directly bind the toggle actions to initiate the fetch immediately, bypassing useEffect delays.
    const setProtocol = useCallback((p: 'grpc' | 'rest') => {
        if (p === protocolRef.current) return;
        protocolRef.current = p;
        _setProtocol(p);
        fetchAll(p, tlsRef.current, false);
    }, [fetchAll]);

    const setTlsEnabled = useCallback((tls: boolean) => {
        if (tls === tlsRef.current) return;
        tlsRef.current = tls;
        _setTlsEnabled(tls);
        fetchAll(protocolRef.current, tls, false);
    }, [fetchAll]);

    const refresh = useCallback(() => {
        fetchAll(protocolRef.current, tlsRef.current, true);
    }, [fetchAll]);

    const handlePing = useCallback(async () => {
        try {
            setPinging(true);
            const start = performance.now();
            const data = await api.ping(protocolRef.current, tlsRef.current);
            const latencyMs = data.latencyMs || Number.parseFloat((performance.now() - start).toFixed(2));
            setPingResult({ ...data, latencyMs });
            setError(null);
        } catch (err: any) {
            console.error('Ping failed:', err);
            setError('Ping failed');
        } finally {
            setPinging(false);
        }
    }, []);

    useEffect(() => {
        // Only run on mount and interval. Toggles are handled explicitly above.
        fetchAll(protocolRef.current, tlsRef.current, !initialLoadDone.current);

        const interval = setInterval(() => {
            fetchAll(protocolRef.current, tlsRef.current, false);
        }, 10000);

        return () => clearInterval(interval);
    }, [fetchAll]);

    return {
        mode, stats, pingResult,
        protocol, tlsEnabled,
        loading, fetchingStats, pinging, error,
        setProtocol, setTlsEnabled, refresh, handlePing,
    };
}
