import React, { useState, useEffect, useCallback } from 'react';
import { RefreshCw, Activity, Zap, BarChart2 } from 'lucide-react';
import { api } from '../api/client';
import type { ModeInfo, BenchmarkStats, PingResponse } from '../types/api';

import { ModeBadge } from '../components/ModeBadge';
import { StatsCard } from '../components/StatsCard';
import { LatencyChart } from '../components/LatencyChart';

export const Dashboard: React.FC = () => {
    const [mode, setMode] = useState<ModeInfo | null>(null);
    const [stats, setStats] = useState<BenchmarkStats | null>(null);
    const [pingResult, setPingResult] = useState<PingResponse | null>(null);
    const [protocol, setProtocol] = useState<'grpc' | 'rest'>('grpc');

    const [loadingMode, setLoadingMode] = useState(true);
    const [loadingStats, setLoadingStats] = useState(true);
    const [pinging, setPinging] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchMode = useCallback(async (target: 'grpc' | 'rest') => {
        try {
            setLoadingMode(true);
            const data = await api.getMode(target);
            setMode(data);
            setError(null);
        } catch (err: any) {
            console.error('Failed to fetch mode:', err);
            setError('Could not connect to gateway');
        } finally {
            setLoadingMode(false);
        }
    }, []);

    const fetchStats = useCallback(async (target: 'grpc' | 'rest') => {
        try {
            setLoadingStats(true);
            const data = await api.getLatestBench(target);
            setStats(data);
            setError(null);
        } catch (err: any) {
            console.error('Failed to fetch stats:', err);
            // It's okay if stats don't exist yet, so we don't necessarily error out the whole page
        } finally {
            setLoadingStats(false);
        }
    }, []);

    const handlePing = async () => {
        try {
            setPinging(true);
            setPingResult(null);
            const start = performance.now();
            const data = await api.ping(protocol);
            // Ensure we have at least frontend measured latency if backend didn't supply it
            const latencyMs = data.latencyMs || parseFloat((performance.now() - start).toFixed(2));
            setPingResult({ ...data, latencyMs });
            setError(null);
        } catch (err: any) {
            console.error('Ping failed:', err);
            setError('Ping failed');
        } finally {
            setPinging(false);
        }
    };

    useEffect(() => {
        fetchMode(protocol);
        fetchStats(protocol);

        // Optional polling for mode / stats
        const interval = setInterval(() => {
            fetchMode(protocol);
            fetchStats(protocol);
        }, 10000); // 10s refresh

        return () => clearInterval(interval);
    }, [fetchMode, fetchStats, protocol]);

    return (
        <div className="mx-auto max-w-5xl px-4 py-12 sm:px-6 lg:px-8">
            {/* Header Section */}
            <div className="mb-10 flex flex-col items-start justify-between gap-6 sm:flex-row sm:items-center">
                <div>
                    <h1 className="text-3xl font-extrabold tracking-tight text-zinc-50 flex items-center gap-3 mb-2">
                        <Activity className="h-8 w-8 text-indigo-500" />
                        Backend Lab Metrics
                    </h1>
                    <p className="text-zinc-400">
                        Real-time performance benchmark across gRPC / REST boundaries
                    </p>
                </div>

                <div className="flex flex-col items-end gap-3">
                    <div className="flex gap-2 p-1 rounded-lg bg-zinc-900 border border-zinc-800">
                        <button
                            onClick={() => setProtocol('grpc')}
                            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition-all ${protocol === 'grpc'
                                    ? 'bg-indigo-500/20 text-indigo-400 border border-indigo-500/30'
                                    : 'text-zinc-500 hover:text-zinc-300'
                                }`}
                        >
                            gRPC
                        </button>
                        <button
                            onClick={() => setProtocol('rest')}
                            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition-all ${protocol === 'rest'
                                    ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                                    : 'text-zinc-500 hover:text-zinc-300'
                                }`}
                        >
                            REST
                        </button>
                    </div>

                    <div className="flex items-center gap-3">
                        <ModeBadge mode={mode} />
                        <button
                            onClick={() => { fetchMode(protocol); fetchStats(protocol); }}
                            className="group flex items-center gap-2 rounded-lg bg-zinc-800/80 px-3 py-1.5 text-xs font-medium text-zinc-300 transition-colors hover:bg-zinc-700 hover:text-white border border-zinc-700/50"
                            disabled={loadingMode || loadingStats}
                        >
                            <RefreshCw size={14} className={loadingMode || loadingStats ? 'animate-spin' : 'group-hover:rotate-180 transition-transform duration-500'} />
                            Refresh
                        </button>
                    </div>
                </div>
            </div>

            {error ? (
                <div className="mb-8 rounded-xl bg-rose-500/10 p-4 border border-rose-500/20 text-rose-400 font-medium flex items-center gap-3">
                    <Activity size={20} />
                    {error}
                </div>
            ) : null}

            {/* Stats Grid */}
            <div className="mb-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
                <StatsCard
                    title="Throughput (RPS)"
                    value={stats ? Math.round(stats.rps).toLocaleString() : '--'}
                    subtitle="req/s"
                    icon={<Zap size={24} className="text-amber-400" />}
                />
                <StatsCard
                    title="p50 Latency"
                    value={stats ? stats.p50.toFixed(2) : '--'}
                    subtitle="ms"
                    icon={<BarChart2 size={24} className="text-emerald-400" />}
                />
                <StatsCard
                    title="p95 Latency"
                    value={stats ? stats.p95.toFixed(2) : '--'}
                    subtitle="ms"
                    icon={<BarChart2 size={24} className="text-amber-400" />}
                />
                <StatsCard
                    title="p99 Latency"
                    value={stats ? stats.p99.toFixed(2) : '--'}
                    subtitle="ms"
                    icon={<BarChart2 size={24} className="text-rose-400" />}
                />
            </div>

            {/* Content Section: Chart and Action Panel */}
            <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
                {/* Latency Chart */}
                <div className="lg:col-span-2 rounded-2xl border border-zinc-800 bg-zinc-900/40 p-6 backdrop-blur-sm">
                    <div className="mb-4 flex items-center justify-between">
                        <h2 className="text-lg font-semibold text-zinc-200">Latency Distribution</h2>
                        {stats && (
                            <span className="text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                Protocol: {stats.protocol}
                            </span>
                        )}
                    </div>
                    <LatencyChart stats={stats} />
                </div>

                {/* Action Panel */}
                <div className="flex flex-col gap-6 rounded-2xl border border-zinc-800 bg-zinc-900/40 p-6 backdrop-blur-sm">
                    <div>
                        <h2 className="mb-2 text-lg font-semibold text-zinc-200">Test Connection</h2>
                        <p className="mb-6 text-sm text-zinc-400">
                            Run a quick ping test through the gateway proxy to verify connectivity and measure baseline round-trip time.
                        </p>
                    </div>

                    <button
                        onClick={handlePing}
                        disabled={pinging || !mode}
                        className="flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-500 px-6 py-3.5 text-sm font-semibold text-white shadow-lg shadow-indigo-500/20 transition-all hover:bg-indigo-400 disabled:opacity-50 disabled:shadow-none"
                    >
                        {pinging ? (
                            <RefreshCw className="animate-spin" size={18} />
                        ) : (
                            <Zap size={18} className="fill-white" />
                        )}
                        {pinging ? 'Pinging...' : 'Run Ping'}
                    </button>

                    {pingResult && (
                        <div className="mt-4 rounded-xl border border-emerald-500/20 bg-emerald-500/10 p-4 text-center animate-in fade-in slide-in-from-bottom-2 duration-300">
                            <span className="block text-2xl font-bold tracking-tight text-emerald-400">
                                {pingResult.latencyMs} <span className="text-sm font-medium opacity-70">ms</span>
                            </span>
                            <span className="mt-1 block text-xs font-medium uppercase tracking-widest text-emerald-500/70">
                                Round Trip
                            </span>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};
