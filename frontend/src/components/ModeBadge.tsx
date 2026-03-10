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
            <div className={cn("inline-flex items-center gap-3 px-5 py-2.5 rounded-2xl bg-white/40 backdrop-blur-md text-redwood-text-muted text-[10px] font-black uppercase tracking-[0.3em] animate-pulse border border-white/40 shadow-xl", className)}>
                <Cpu size={14} className="animate-spin-slow" />
                INITIALIZING...
            </div>
        );
    }

    const isGrpc = protocol === 'grpc';

    return (
        <div className={cn("flex items-center gap-4", className)}>
            <div className={cn(
                "inline-flex items-center gap-3 px-5 py-2.5 rounded-2xl text-[10px] font-black uppercase tracking-[0.3em] border shadow-2xl transition-all duration-500 hover:scale-105",
                isGrpc
                    ? "bg-redwood-pine/80 backdrop-blur-md text-white border-redwood-pine/30"
                    : "bg-redwood-teal/80 backdrop-blur-md text-white border-redwood-teal/30"
            )}>
                <div className="relative flex h-2.5 w-2.5">
                    <div className="animate-ping absolute inline-flex h-full w-full rounded-full bg-white opacity-40"></div>
                    <div className="relative inline-flex rounded-full h-2.5 w-2.5 bg-white shadow-sm"></div>
                </div>
                {protocol}
            </div>

            <div className={cn(
                "inline-flex items-center gap-3 px-5 py-2.5 rounded-2xl text-[10px] font-black uppercase tracking-[0.3em] border shadow-2xl transition-all duration-500 hover:scale-105",
                tlsEnabled
                    ? "bg-white/90 backdrop-blur-md text-redwood-pine border-white/50"
                    : "bg-redwood-plum/80 backdrop-blur-md text-white border-redwood-plum/30"
            )}>
                {tlsEnabled ? <Shield size={14} className="text-redwood-pine" /> : <ShieldAlert size={14} />}
                {tlsEnabled ? 'SECURED' : 'PLAIN'}
            </div>
        </div>
    );
};
