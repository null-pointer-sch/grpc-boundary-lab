import React from 'react';
import { Shield, ShieldAlert, Cpu } from 'lucide-react';
import { cn } from '../utils/cn';
import type { ModeInfo } from '../core/types/api';

interface ModeBadgeProps {
    mode: ModeInfo | null;
    protocol: 'grpc' | 'rest';
    tlsEnabled: boolean;
    className?: string;
}

export const ModeBadge: React.FC<ModeBadgeProps> = ({ mode, protocol, tlsEnabled, className }) => {
    if (!mode) {
        return (
            <div className={cn("inline-flex items-center gap-2 px-3 py-1 rounded-full bg-zinc-800 text-zinc-400 text-sm font-medium animate-pulse", className)}>
                <Cpu size={16} />
                Connecting...
            </div>
        );
    }

    const isGrpc = protocol === 'grpc';

    return (
        <div className={cn("flex items-center gap-3", className)}>
            <div className={cn(
                "inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-semibold border",
                isGrpc
                    ? "bg-indigo-500/10 text-indigo-400 border-indigo-500/20"
                    : "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
            )}>
                <span className="relative flex h-2 w-2 mr-1">
                    <span className={cn(
                        "animate-ping absolute inline-flex h-full w-full rounded-full opacity-75",
                        isGrpc ? "bg-indigo-400" : "bg-emerald-400"
                    )}></span>
                    <span className={cn(
                        "relative inline-flex rounded-full h-2 w-2",
                        isGrpc ? "bg-indigo-500" : "bg-emerald-500"
                    )}></span>
                </span>
                {protocol.toUpperCase()}
            </div>

            <div className={cn(
                "inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium border",
                tlsEnabled
                    ? "bg-zinc-800/50 text-zinc-300 border-zinc-700/50"
                    : "bg-amber-500/10 text-amber-500 border-amber-500/20"
            )}>
                {tlsEnabled ? <Shield size={14} className="text-emerald-500" /> : <ShieldAlert size={14} />}
                {tlsEnabled ? 'TLS Enabled' : 'No TLS'}
            </div>
        </div>
    );
};
