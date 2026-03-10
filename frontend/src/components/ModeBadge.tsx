import React from 'react';
import { Shield, ShieldAlert } from 'lucide-react';
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
            <div className={cn("rw-badge rw-badge-info animate-pulse", className)}>
                Initializing...
            </div>
        );
    }

    return (
        <div className={cn("flex items-center gap-2", className)}>
            <div className={cn(
                "rw-badge",
                protocol === 'grpc' ? "rw-badge-info" : "rw-badge-warning"
            )}>
                {protocol}
            </div>

            <div className={cn(
                "rw-badge font-bold",
                tlsEnabled ? "rw-badge-success" : "rw-badge-error"
            )}>
                {tlsEnabled ? (
                    <><Shield size={12} className="mr-1.5" /> Secured</>
                ) : (
                    <><ShieldAlert size={12} className="mr-1.5" /> Unsecured</>
                )}
            </div>
        </div>
    );
};
