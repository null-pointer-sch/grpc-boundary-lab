import React from 'react';
import { RefreshCw, Activity, Zap, BarChart2, Shield, ShieldOff } from 'lucide-react';

import { useDashboard } from '../hooks/useDashboard';
import { ModeBadge } from '../components/ModeBadge';
import { StatsCard } from '../components/StatsCard';
import { LatencyChart } from '../components/LatencyChart';

export const Dashboard: React.FC = () => {
    const {
        mode, stats, pingResult,
        protocol, tlsEnabled,
        loading, fetchingStats, pinging, error,
        setProtocol, setTlsEnabled, refresh, handlePing,
    } = useDashboard();

    return (
        <div className="mx-auto max-w-5xl px-4 py-12 sm:px-6 lg:px-8">
            {/* Header */}
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

                <div className="flex flex-col items-end gap-3 min-w-[200px]">
                    {/* Protocol Toggle */}
                    <div className="flex gap-2 p-1 rounded-lg bg-zinc-900 border border-zinc-800">
                        <button
                            type="button"
                            onClick={(e) => { e.preventDefault(); setProtocol('grpc'); }}
                            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition-all ${protocol === 'grpc'
                                ? 'bg-indigo-500/20 text-indigo-400 border border-indigo-500/30'
                                : 'text-zinc-500 hover:text-zinc-300'
                                }`}
                        >
                            gRPC
                        </button>
                        <button
                            type="button"
                            onClick={(e) => { e.preventDefault(); setProtocol('rest'); }}
                            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition-all ${protocol === 'rest'
                                ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                                : 'text-zinc-500 hover:text-zinc-300'
                                }`}
                        >
                            REST
                        </button>
                    </div>

                    {/* TLS Toggle */}
                    <div className="flex gap-2 p-1 rounded-lg bg-zinc-900 border border-zinc-800">
                        <button
                            type="button"
                            onClick={(e) => { e.preventDefault(); setTlsEnabled(false); }}
                            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition-all ${!tlsEnabled
                                ? 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
                                : 'text-zinc-500 hover:text-zinc-300'
                                }`}
                        >
                            <ShieldOff size={12} />
                            Plain
                        </button>
                        <button
                            type="button"
                            onClick={(e) => { e.preventDefault(); setTlsEnabled(true); }}
                            disabled={mode !== null && !mode.tlsAvailable}
                            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition-all ${tlsEnabled
                                ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                                : 'text-zinc-500 hover:text-zinc-300'
                                } disabled:opacity-30 disabled:cursor-not-allowed`}
                            title={mode !== null && !mode.tlsAvailable ? 'TLS certs not available' : 'Enable mTLS'}
                        >
                            <Shield size={12} />
                            mTLS
                        </button>
                    </div>

                    <div className="flex items-center gap-3">
                        <ModeBadge mode={mode} protocol={protocol} tlsEnabled={tlsEnabled} />
                        <button
                            type="button"
                            onClick={(e) => { e.preventDefault(); refresh(); }}
                            className="group flex items-center gap-2 rounded-lg bg-zinc-800/80 px-3 py-1.5 text-xs font-medium text-zinc-300 transition-colors hover:bg-zinc-700 hover:text-white border border-zinc-700/50"
                            disabled={loading}
                        >
                            <RefreshCw size={14} className={loading ? 'animate-spin' : 'group-hover:rotate-180 transition-transform duration-500'} />
                            Refresh
                        </button>
                    </div>
                </div>
            </div>

            {/* Error */}
            {error && (
                <div className="mb-8 rounded-xl bg-rose-500/10 p-4 border border-rose-500/20 text-rose-400 font-medium flex items-center gap-3">
                    <Activity size={20} />
                    {error}
                </div>
            )}

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

            {/* Chart + Action Panel */}
            <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
                <div className="lg:col-span-2 rounded-2xl border border-zinc-800 bg-zinc-900/40 p-6 backdrop-blur-sm">
                    <div className="mb-6 flex items-center justify-between">
                        <div className="flex items-baseline gap-3">
                            <h2 className="text-lg font-semibold text-zinc-200">Latency Distribution</h2>
                            <span className="text-xs font-medium text-zinc-500 bg-zinc-800/50 px-2 py-0.5 rounded border border-zinc-700/30 shadow-sm flex items-center gap-1">
                                <span className="font-bold text-emerald-500/80">↓</span> Lower is better
                            </span>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="text-xs font-medium text-zinc-500 uppercase tracking-wider">
                                {protocol.toUpperCase()}
                            </span>
                            <span className={`inline-flex items-center gap-1 text-xs font-medium text-emerald-500 transition-opacity ${tlsEnabled ? 'opacity-100' : 'opacity-0'}`}>
                                <Shield size={10} /> TLS
                            </span>
                        </div>
                    </div>
                    <div className={`transition-opacity duration-300 ${fetchingStats ? 'opacity-50 grayscale' : 'opacity-100'}`}>
                        <LatencyChart stats={stats} />
                    </div>
                </div>

                {/* Action Panel */}
                <div className="flex flex-col gap-6 rounded-2xl border border-zinc-800 bg-zinc-900/40 p-6 backdrop-blur-sm">
                    <div>
                        <h2 className="mb-2 text-lg font-semibold text-zinc-200">Test Connection</h2>
                        <p className="mb-6 text-sm text-zinc-400">
                            Run a quick ping test through the gateway proxy to verify connectivity and measure baseline round-trip time.
                        </p>
                    </div>

                    <div className="min-w-[140px]">
                        <button
                            type="button"
                            onClick={(e) => { e.preventDefault(); handlePing(); }}
                            disabled={pinging || !mode}
                            className="flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-500 px-6 py-3.5 text-sm font-semibold text-white shadow-lg shadow-indigo-500/20 transition-all hover:bg-indigo-400 disabled:opacity-50 disabled:shadow-none"
                        >
                            <div className="flex w-[100px] items-center justify-center gap-2">
                                {pinging ? (
                                    <RefreshCw className="animate-spin" size={18} />
                                ) : (
                                    <Zap size={18} className="fill-white" />
                                )}
                                <span>{pinging ? 'Pinging...' : 'Run Ping'}</span>
                            </div>
                        </button>
                    </div>

                    {pingResult && (
                        <div className={`mt-4 rounded-xl border border-emerald-500/20 bg-emerald-500/10 p-4 text-center animate-in fade-in slide-in-from-bottom-2 duration-300 transition-opacity ${pinging ? 'opacity-50 grayscale' : 'opacity-100'}`}>
                            <span className="block text-2xl font-bold tracking-tight text-emerald-400 tabular-nums">
                                {pingResult.latencyMs} <span className="text-sm font-medium opacity-70">ms</span>
                            </span>
                            <span className="mt-1 block text-xs font-medium uppercase tracking-widest text-emerald-500/70">
                                Round Trip
                            </span>
                        </div>
                    )}

                    {/* Security Status */}
                    <div className={`mt-auto rounded-xl p-4 border transition-colors duration-300 ${tlsEnabled
                        ? 'border-emerald-500/20 bg-emerald-500/5'
                        : 'border-zinc-700/50 bg-zinc-800/30'
                        }`}>
                        <div className="flex items-center gap-2 mb-2">
                            {tlsEnabled ? (
                                <Shield size={16} className="text-emerald-500" />
                            ) : (
                                <ShieldOff size={16} className="text-zinc-500" />
                            )}
                            <span className={`text-xs font-semibold uppercase tracking-wider ${tlsEnabled ? 'text-emerald-400' : 'text-zinc-500'
                                }`}>
                                {tlsEnabled ? 'mTLS Active' : 'Plaintext'}
                            </span>
                        </div>
                        <p className="text-xs text-zinc-500 min-h-[3rem]">
                            {tlsEnabled
                                ? 'Traffic encrypted with local CA certificates. Gateway → Backend verified via mTLS.'
                                : mode?.tlsAvailable
                                    ? 'TLS available — toggle mTLS above to encrypt gateway traffic.'
                                    : 'No TLS certificates found. Run certs/gen.sh to generate.'}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};
