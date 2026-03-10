import React from 'react';
import { cn } from '../utils/cn';

interface StatsCardProps {
    title: string;
    value: string | number;
    subtitle?: string;
    icon?: React.ReactNode;
    className?: string;
}

export const StatsCard: React.FC<StatsCardProps> = ({
    title,
    value,
    subtitle,
    icon,
    className
}) => {
    return (
        <div className={cn(
            "rw-card p-6 group",
            className
        )}>
            <div className="flex items-center justify-between mb-4">
                <span className="text-[12px] font-bold text-text-sub uppercase tracking-wider">{title}</span>
                {icon && (
                    <div className="p-2 bg-neutral-100 rounded-md text-text-sub group-hover:bg-neutral-200 transition-colors">
                        {React.cloneElement(icon as React.ReactElement<any>, { size: 18 })}
                    </div>
                )}
            </div>

            <div className="flex items-baseline gap-2">
                <h3 className="text-[28px] font-bold text-text-main tabular-nums tracking-tight">{value}</h3>
                {subtitle && (
                    <span className="text-[12px] font-medium text-text-sub uppercase">{subtitle}</span>
                )}
            </div>
            
            {/* Subtle bottom accent for high-density feel */}
            <div className="mt-4 h-1 w-full bg-neutral-100 rounded-full overflow-hidden">
                <div className="h-full bg-primary-red/10 w-1/3" />
            </div>
        </div>
    );
};
