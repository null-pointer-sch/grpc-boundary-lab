import React from 'react';
import { Shield, ShieldAlert, Cpu } from 'lucide-react';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import type { ModeInfo } from '../types/api';

/** Utility to merge tailwind classes */
export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

interface ModeBadgeProps {
    mode: ModeInfo | null;
    className?: string;
}

export const ModeBadge: React.FC<ModeBadgeProps> = ({ mode, className }) => {
    if (!mode) {
        return (
            <div className={cn("inline-flex items-center gap-2 px-3 py-1 rounded-full bg-zinc-800 text-zinc-400 text-sm font-medium animate-pulse", className)}>
                <Cpu size={16} />
                Connecting...
            </div>
        );
    }

    const isGrpc = mode.protocol === 'grpc';

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
                {mode.protocol.toUpperCase()}
            </div>

            <div className={cn(
                "inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium border",
                mode.tls
                    ? "bg-zinc-800/50 text-zinc-300 border-zinc-700/50"
                    : "bg-amber-500/10 text-amber-500 border-amber-500/20"
            )}>
                {mode.tls ? <Shield size={14} className="text-emerald-500" /> : <ShieldAlert size={14} />}
                {mode.tls ? 'TLS Enabled' : 'No TLS'}
            </div>
        </div>
    );
};
