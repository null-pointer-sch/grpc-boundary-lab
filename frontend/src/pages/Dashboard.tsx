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
        <div className="mx-auto max-w-6xl px-4 pt-2 sm:pt-4 pb-20 sm:px-6 lg:px-8">
            {/* Artistic Header Section */}
            <div className="mb-16 flex flex-col items-start justify-between gap-10 md:flex-row md:items-center p-8 rounded-[40px] bg-white/20 backdrop-blur-2xl border border-white/30 shadow-2xl">
                <div className="flex items-center gap-6">
                    <div className="rounded-3xl bg-redwood-red p-5 shadow-2xl shadow-redwood-red/40 transform -rotate-6 hover:rotate-0 transition-all duration-700 hover:scale-110">
                        <Activity className="h-12 w-12 text-white" />
                    </div>
                    <div>
                        <h1 className="text-5xl font-black tracking-tighter text-redwood-text mb-2 uppercase drop-shadow-sm">
                            Realms of <span className="text-redwood-red">Data</span>
                        </h1>
                        <p className="text-xs font-black text-redwood-text-muted/80 uppercase tracking-[0.3em] ml-1">
                            Oracle Redwood Artistic Performance Lab
                        </p>
                    </div>
                </div>

                <div className="flex flex-col items-end gap-5 min-w-[320px]">
                    <div className="flex items-center gap-4">
                        {/* Protocol Control */}
                        <div className="redwood-segmented-control">
                            <button
                                type="button"
                                onClick={(e) => { e.preventDefault(); setProtocol('grpc'); }}
                                className={`redwood-segmented-item ${protocol === 'grpc' ? 'redwood-segmented-item-active' : 'redwood-segmented-item-inactive'}`}
                            >
                                gRPC
                            </button>
                            <button
                                type="button"
                                onClick={(e) => { e.preventDefault(); setProtocol('rest'); }}
                                className={`redwood-segmented-item ${protocol === 'rest' ? 'redwood-segmented-item-active' : 'redwood-segmented-item-inactive'}`}
                            >
                                REST
                            </button>
                        </div>

                        {/* Security Control */}
                        <div className="redwood-segmented-control">
                            <button
                                type="button"
                                onClick={(e) => { e.preventDefault(); setTlsEnabled(false); }}
                                className={`redwood-segmented-item flex items-center gap-2 ${!tlsEnabled ? 'redwood-segmented-item-active' : 'redwood-segmented-item-inactive'}`}
                            >
                                <ShieldOff size={12} />
                                PLAIN
                            </button>
                            <button
                                type="button"
                                onClick={(e) => { e.preventDefault(); setTlsEnabled(true); }}
                                disabled={mode !== null && !mode.tlsAvailable}
                                className={`redwood-segmented-item flex items-center gap-2 ${tlsEnabled ? 'redwood-segmented-item-active' : 'redwood-segmented-item-inactive'} disabled:opacity-30`}
                            >
                                <Shield size={12} />
                                mTLS
                            </button>
                        </div>
                    </div>

                    <div className="flex items-center gap-4">
                        <ModeBadge mode={mode} protocol={protocol} tlsEnabled={tlsEnabled} className="redwood-glass-badge" />
                        <button
                            type="button"
                            onClick={(e) => { e.preventDefault(); refresh(); }}
                            className="redwood-button-ghost"
                            disabled={loading}
                        >
                            <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
                            REFRESH
                        </button>
                    </div>
                </div>
            </div>

            {/* Error Overlay */}
            {error && (
                <div className="mb-12 rounded-3xl bg-redwood-rose/20 backdrop-blur-xl p-6 border border-redwood-rose/30 text-redwood-rose font-black uppercase tracking-widest flex items-center gap-5 shadow-2xl animate-pulse">
                    <Activity size={28} />
                    <span className="text-lg">{error}</span>
                </div>
            )}

            {/* Glass Stats Grid */}
            <div className="mb-12 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-4">
                <StatsCard
                    title="THROUGHPUT"
                    value={stats ? Math.round(stats.rps).toLocaleString() : '--'}
                    subtitle="RPS"
                    icon={<Zap size={24} className="text-redwood-red" />}
                    className="redwood-card border-b-8 border-b-redwood-red/30"
                />
                <StatsCard
                    title="P50 LATENCY"
                    value={stats ? stats.p50.toFixed(2) : '--'}
                    subtitle="MS"
                    icon={<BarChart2 size={24} className="text-redwood-teal" />}
                    className="redwood-card border-b-8 border-b-redwood-teal/30"
                />
                <StatsCard
                    title="P95 LATENCY"
                    value={stats ? stats.p95.toFixed(2) : '--'}
                    subtitle="MS"
                    icon={<BarChart2 size={24} className="text-redwood-plum" />}
                    className="redwood-card border-b-8 border-b-redwood-plum/30"
                />
                <StatsCard
                    title="P99 LATENCY"
                    value={stats ? stats.p99.toFixed(2) : '--'}
                    subtitle="MS"
                    icon={<BarChart2 size={24} className="text-redwood-rose" />}
                    className="redwood-card border-b-8 border-b-redwood-rose/30"
                />
            </div>

            {/* Secondary Content Layer */}
            <div className="grid grid-cols-1 gap-10 lg:grid-cols-3">
                <div className="lg:col-span-2 redwood-card p-10">
                    <div className="mb-10 flex items-center justify-between">
                        <div className="flex items-center gap-6">
                            <h2 className="text-2xl font-black text-redwood-text uppercase tracking-tighter">Performance Realm</h2>
                            <div className="flex items-center gap-2 px-6 py-2 bg-redwood-pine/10 rounded-full border border-redwood-pine/20 backdrop-blur-md">
                                <Activity size={14} className="text-redwood-pine" />
                                <span className="text-[10px] font-black text-redwood-pine uppercase tracking-[0.3em]">
                                    Live Stream
                                </span>
                            </div>
                        </div>
                        <div className="text-[10px] font-black text-redwood-text-muted/60 uppercase tracking-[0.4em]">
                            MODE: {protocol}
                        </div>
                    </div>
                    <div className={`transition-all duration-1000 ${fetchingStats ? 'opacity-20 blur-md scale-[0.95]' : 'opacity-100'}`}>
                        <LatencyChart stats={stats} />
                    </div>
                </div>

                {/* Connection Portal */}
                <div className="flex flex-col gap-10 redwood-card p-10 border-t-8 border-t-redwood-red">
                    <div>
                        <h2 className="mb-6 text-2xl font-black text-redwood-text uppercase tracking-tighter border-b border-redwood-text/10 pb-4">Connection Portal</h2>
                        <p className="text-xs font-black text-redwood-text-muted/60 leading-relaxed uppercase tracking-widest">
                            Initiate a high-speed probe across the repository boundaries. Visualizing real-time round-trip latency.
                        </p>
                    </div>

                    <button
                        type="button"
                        onClick={(e) => { e.preventDefault(); handlePing(); }}
                        disabled={pinging || !mode}
                        className="redwood-button-primary group"
                    >
                        <div className="flex items-center justify-center gap-4">
                            {pinging ? (
                                <RefreshCw className="animate-spin" size={24} />
                            ) : (
                                <Zap size={24} className="fill-white group-hover:scale-125 transition-transform" />
                            )}
                            <span className="text-lg">{pinging ? 'PROBING...' : 'INITIALIZE PING'}</span>
                        </div>
                    </button>

                    {pingResult && (
                        <div className={`rounded-[32px] border-2 border-redwood-teal/30 bg-redwood-teal/10 backdrop-blur-2xl p-10 text-center transition-all duration-700 ${pinging ? 'opacity-20 scale-90 translate-y-4' : 'opacity-100 shadow-3xl shadow-redwood-teal/10'}`}>
                            <span className="block text-6xl font-black tracking-tighter text-redwood-teal tabular-nums mb-3 drop-shadow-md">
                                {pingResult.latencyMs} <span className="text-xl opacity-40">MS</span>
                            </span>
                            <div className="flex items-center justify-center gap-2">
                                <div className="h-1 w-1 rounded-full bg-redwood-teal animate-ping" />
                                <span className="text-[10px] font-black uppercase tracking-[0.5em] text-redwood-teal/60">
                                    ROUND TRIP
                                </span>
                            </div>
                        </div>
                    )}

                    <div className={`mt-auto rounded-3xl p-6 border-2 transition-all duration-500 ${tlsEnabled
                        ? 'border-redwood-pine/30 bg-redwood-pine/10 shadow-2xl'
                        : 'border-white/20 bg-white/10 shadow-lg'
                        }`}>
                        <div className="flex items-center gap-3 mb-4">
                            <div className={`p-2 rounded-xl backdrop-blur-md ${tlsEnabled ? 'bg-redwood-pine text-white' : 'bg-white/20 text-redwood-text-muted'}`}>
                                {tlsEnabled ? <Shield size={20} /> : <ShieldOff size={20} />}
                            </div>
                            <span className={`text-[11px] font-black uppercase tracking-[0.3em] ${tlsEnabled ? 'text-redwood-pine' : 'text-redwood-text-muted'}`}>
                                {tlsEnabled ? 'IDENTITY SECURED' : 'OPEN TUNNEL'}
                            </span>
                        </div>
                        <p className="text-[10px] font-black text-redwood-text-muted/60 leading-loose uppercase tracking-[0.15em]">
                            {tlsEnabled
                                ? 'Transmission encrypted via mutual TLS certificates. Cryptographic proof of identity established.'
                                : 'Payload transmitted in plaintext. Security layer bypassed for maximum performance analysis.'}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};
