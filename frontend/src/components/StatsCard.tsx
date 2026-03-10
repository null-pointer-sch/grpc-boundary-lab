import React from 'react';
import { cn } from '../utils/cn';

interface StatsCardProps {
    title: string;
    value: string | number;
    subtitle?: string;
    icon?: React.ReactNode;
    trend?: {
        value: number;
        isPositive: boolean;
    };
    className?: string;
}

export const StatsCard: React.FC<StatsCardProps> = ({
    title,
    value,
    subtitle,
    icon,
    trend,
    className
}) => {
    return (
        <div className={cn(
            "group relative overflow-hidden p-8",
            className
        )}>
            {/* Dynamic Glass Gradient */}
            <div className="absolute inset-0 bg-gradient-to-br from-white/40 to-white/5 opacity-0 group-hover:opacity-100 transition-opacity duration-700" />
            
            <div className="relative flex flex-col items-start justify-between h-full gap-8">
                <div className="flex items-center justify-between w-full">
                    <p className="text-[10px] font-black uppercase tracking-[0.4em] text-redwood-text-muted/70">{title}</p>
                    {icon && (
                        <div className="rounded-[20px] bg-white/40 p-3 text-redwood-text-muted/80 shadow-lg border border-white/40 group-hover:scale-110 group-hover:bg-white/60 transition-all duration-500">
                            {React.cloneElement(icon as React.ReactElement<any>, { size: 20 })}
                        </div>
                    )}
                </div>

                <div className="flex flex-col gap-1">
                    <div className="flex items-baseline gap-3">
                        <h3 className="text-5xl font-black tracking-tighter text-redwood-text tabular-nums drop-shadow-sm">{value}</h3>
                        {subtitle && (
                            <span className="text-[10px] font-black text-redwood-text-muted/30 uppercase tracking-[0.3em] font-sans">{subtitle}</span>
                        )}
                    </div>
                </div>

                {trend && (
                    <div className="flex items-center gap-3 text-[10px] font-black uppercase tracking-[0.3em]">
                        <span className={cn(
                            "px-3 py-1.5 rounded-full shadow-lg backdrop-blur-md border border-white/40",
                            trend.isPositive ? "bg-redwood-pine/20 text-redwood-pine" : "bg-redwood-rose/20 text-redwood-rose"
                        )}>
                            {trend.value > 0 ? '↑' : '↓'} {Math.abs(trend.value)}%
                        </span>
                        <span className="text-redwood-text-muted/40">DELTA</span>
                    </div>
                )}
            </div>

            {/* Micro-texture Overlay */}
            <div className="absolute inset-0 pointer-events-none opacity-20 mix-blend-overlay bg-[url('https://www.transparenttextures.com/patterns/pinstriped-suit.png')]" />
        </div>
    );
};
