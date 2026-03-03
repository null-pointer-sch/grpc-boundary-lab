import React from 'react';
import { cn } from './ModeBadge';

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
            "relative overflow-hidden rounded-2xl border border-zinc-800 bg-zinc-900/50 p-6 backdrop-blur-sm transition-all hover:bg-zinc-900/80",
            className
        )}>
            {/* Decorative gradient blob */}
            <div className="absolute -right-6 -top-6 h-24 w-24 rounded-full bg-indigo-500/10 blur-2xl" />

            <div className="relative flex items-center justify-between">
                <div>
                    <p className="text-sm font-medium text-zinc-400">{title}</p>
                    <div className="mt-2 flex items-baseline gap-2">
                        <h3 className="text-3xl font-bold tracking-tight text-zinc-100">{value}</h3>
                        {subtitle && (
                            <span className="text-sm text-zinc-500">{subtitle}</span>
                        )}
                    </div>
                </div>

                {icon && (
                    <div className="rounded-xl bg-zinc-800/80 p-3 text-zinc-400 ring-1 ring-inset ring-zinc-700/50">
                        {icon}
                    </div>
                )}
            </div>

            {trend && (
                <div className="mt-4 flex items-center gap-2 text-sm">
                    <span className={cn(
                        "font-medium",
                        trend.isPositive ? "text-emerald-400" : "text-rose-400"
                    )}>
                        {trend.value > 0 ? '+' : ''}{trend.value}%
                    </span>
                    <span className="text-zinc-500">vs last run</span>
                </div>
            )}
        </div>
    );
};
