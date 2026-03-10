import React from 'react';
import { 
  RefreshCw, 
  Activity, 
  Zap, 
  BarChart2, 
  Shield, 
  ShieldOff, 
  Settings, 
  ShieldCheck
} from 'lucide-react';

import { useDashboard } from '../hooks/useDashboard';
import { ModeBadge } from '../components/ModeBadge';
import { StatsCard } from '../components/StatsCard';
import { LatencyChart } from '../components/LatencyChart';
import { cn } from '../utils/cn';

const DashboardStatsGrid: React.FC<{ stats: any }> = ({ stats }) => (
    <div className="mb-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <StatsCard
            title="Throughput"
            value={stats ? Math.round(stats.rps).toLocaleString() : '--'}
            subtitle="RPS"
            icon={<Zap className="text-primary-red" />}
        />
        <StatsCard
            title="P50 Latency"
            value={stats ? stats.p50.toFixed(2) : '--'}
            subtitle="MS"
            icon={<BarChart2 className="text-viz-teal" />}
        />
        <StatsCard
            title="P95 Latency"
            value={stats ? stats.p95.toFixed(2) : '--'}
            subtitle="MS"
            icon={<BarChart2 className="text-viz-slate" />}
        />
        <StatsCard
            title="P99 Latency"
            value={stats ? stats.p99.toFixed(2) : '--'}
            subtitle="MS"
            icon={<BarChart2 className="text-primary-red" />}
        />
    </div>
);

export const Dashboard: React.FC = () => {
    const {
        mode, stats, pingResult,
        protocol, tlsEnabled,
        loading, fetchingStats, pinging, error,
        setProtocol, setTlsEnabled, refresh, handlePing,
    } = useDashboard();

    return (
        <div className="p-8 max-w-[1600px] mx-auto">
            {/* Context Header */}
            <div className="mb-8 flex flex-col md:flex-row md:items-end justify-between gap-6">
                <div>
                    <h1 className="text-[32px] font-bold text-text-main leading-tight tracking-tight">Boundary Lab Performance</h1>
                    <p className="text-[16px] text-text-sub mt-1">Real-time observability across secure network boundaries</p>
                </div>
                
                <div className="flex items-center gap-3">
                    <button 
                        onClick={() => refresh()}
                        className="rw-button-tertiary"
                        disabled={loading}
                    >
                        <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
                        Refresh Data
                    </button>
                    <button className="rw-button-primary">
                        <Zap size={16} />
                        Benchmark All
                    </button>
                    <button className="p-2 hover:bg-neutral-100 rounded-md transition-colors">
                        <Settings size={20} className="text-text-sub" />
                    </button>
                </div>
            </div>

            {/* Config & Status Bar */}
            <div className="rw-card p-4 mb-8 flex flex-wrap items-center justify-between gap-6 bg-neutral-100/50">
                <div className="flex items-center gap-8">
                   <div className="flex flex-col gap-1">
                      <span className="text-[10px] font-bold text-text-sub uppercase tracking-wider">Protocol</span>
                      <div className="flex gap-1 p-1 bg-white border border-neutral-200 rounded-md">
                        <button 
                            onClick={() => setProtocol('grpc')}
                            className={cn("px-4 py-1.5 text-xs font-bold rounded transition-all", protocol === 'grpc' ? "bg-primary-red text-white" : "hover:bg-neutral-100")}
                        >gRPC</button>
                        <button 
                            onClick={() => setProtocol('rest')}
                            className={cn("px-4 py-1.5 text-xs font-bold rounded transition-all", protocol === 'rest' ? "bg-primary-red text-white" : "hover:bg-neutral-100")}
                        >REST</button>
                      </div>
                   </div>

                   <div className="flex flex-col gap-1">
                      <span className="text-[10px] font-bold text-text-sub uppercase tracking-wider">Security Layer</span>
                      <div className="flex gap-1 p-1 bg-white border border-neutral-200 rounded-md">
                        <button 
                            onClick={() => setTlsEnabled(false)}
                            className={cn("flex items-center gap-2 px-4 py-1.5 text-xs font-bold rounded transition-all", tlsEnabled ? "hover:bg-neutral-100" : "bg-primary-red text-white")}
                        ><ShieldOff size={14} /> Plain</button>
                        <button 
                            onClick={() => setTlsEnabled(true)}
                            disabled={mode !== null && !mode.tlsAvailable}
                            className={cn("flex items-center gap-2 px-4 py-1.5 text-xs font-bold rounded transition-all", tlsEnabled ? "bg-primary-red text-white" : "hover:bg-neutral-100")}
                        ><Shield size={14} /> mTLS</button>
                      </div>
                   </div>
                </div>

                <div className="flex items-center gap-4">
                   <div className="flex flex-col items-end gap-1">
                      <span className="text-[10px] font-bold text-text-sub uppercase tracking-wider">Current Status</span>
                      <ModeBadge mode={mode} protocol={protocol} tlsEnabled={tlsEnabled} />
                   </div>
                </div>
            </div>

            {/* Error Overlay */}
            {error && (
                <div className="mb-8 rw-badge rw-badge-error w-full py-4 px-6 rounded-lg normal-case tracking-normal text-sm flex items-center gap-3">
                    <Activity size={20} />
                    <span>{error}</span>
                </div>
            )}

            {/* Stats Grid */}
            <DashboardStatsGrid stats={stats} />

            {/* Main Area: Chart + Connection Test */}
            <div className="grid grid-cols-1 gap-8 lg:grid-cols-3 items-start">
                <div className="lg:col-span-2 rw-card">
                    <div className="p-6 border-b border-neutral-100 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                           <h2 className="text-[16px] font-bold text-text-main">Latency Perspective</h2>
                           <span className="rw-badge rw-badge-success text-[10px]">Real-time</span>
                        </div>
                        <span className="text-[12px] text-text-sub">Protocol: <span className="uppercase font-bold text-text-main">{protocol}</span></span>
                    </div>
                    <div className={cn("p-6 transition-all duration-500", fetchingStats ? 'opacity-40 grayscale-[50%]' : 'opacity-100')}>
                        <LatencyChart stats={stats} />
                    </div>
                </div>

                {/* Connection Utility */}
                <div className="flex flex-col gap-6">
                    <div className="rw-card p-6 flex-1">
                        <h2 className="text-[16px] font-bold mb-4 border-b border-neutral-100 pb-3 text-text-main">Test Connection</h2>
                        <p className="text-[14px] text-text-sub mb-6 leading-relaxed">
                            Simulate high-frequency binary probes across the repository boundary to verify tunnel stability.
                        </p>
                        
                        <div className="flex flex-col gap-4">
                            <button
                                onClick={() => handlePing()}
                                disabled={pinging || !mode}
                                className="rw-button-primary w-full h-12 text-md"
                            >
                                {pinging ? (
                                    <>
                                        <RefreshCw size={20} className="mr-2 inline-block animate-spin" />
                                        Pinging...
                                    </>
                                ) : (
                                    <>
                                        <Activity size={20} className="mr-2 inline-block" />
                                        Initialize Probe
                                    </>
                                )}
                            </button>
                            
                            {pingResult && (
                                <div className="mt-4 p-6 bg-neutral-100 rounded-lg border border-neutral-200">
                                   <div className="flex items-baseline justify-center gap-2">
                                      <span className="text-[42px] font-bold text-viz-teal tracking-tighter tabular-nums">{pingResult.latencyMs}</span>
                                      <span className="text-[14px] font-bold text-text-sub">ms</span>
                                   </div>
                                   <div className="flex items-center justify-center gap-2 mt-1">
                                      <div className="w-1.5 h-1.5 rounded-full bg-viz-teal animate-pulse" />
                                      <span className="text-[12px] font-medium text-text-sub uppercase tracking-widest">Calculated Round Trip</span>
                                   </div>
                                </div>
                            )}
                        </div>
                    </div>

                    <div className={cn(
                        "rw-card p-6 border-t-4 transition-all duration-300", 
                        tlsEnabled ? "border-t-success bg-white" : "border-t-viz-amber bg-neutral-100/30"
                    )}>
                        <div className="flex items-center gap-3 mb-3">
                            <div className={cn("p-1.5 rounded-md", tlsEnabled ? "bg-success/10 text-success" : "bg-viz-amber/10 text-viz-amber")}>
                                {tlsEnabled ? <ShieldCheck size={18} /> : <ShieldOff size={18} />}
                            </div>
                            <span className="text-[14px] font-bold uppercase tracking-wide text-text-main">
                                {tlsEnabled ? 'Secure Identity' : 'Public Access'}
                            </span>
                        </div>
                        <p className="text-[12px] text-text-sub leading-loose">
                            {tlsEnabled 
                             ? 'Encryption verified with OCI standard mutual TLS handshake. Security context is strictly enforced.' 
                             : 'Plaintext communication optimized for laboratory observation. Not recommended for production environments.'}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};
