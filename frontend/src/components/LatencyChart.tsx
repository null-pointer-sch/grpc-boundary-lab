import React from 'react';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    Cell
} from 'recharts';
import type { BenchmarkStats } from '../types/api';

interface LatencyChartProps {
    stats: BenchmarkStats | null;
}

export const LatencyChart: React.FC<LatencyChartProps> = ({ stats }) => {
    if (!stats) {
        return (
            <div className="flex h-64 w-full items-center justify-center rounded-2xl border border-zinc-800 border-dashed bg-zinc-900/30">
                <p className="text-sm text-zinc-500">No benchmark data available</p>
            </div>
        );
    }

    const data = [
        { name: 'p50', value: stats.p50, color: '#10b981' }, // Emerald
        { name: 'p95', value: stats.p95, color: '#f59e0b' }, // Amber
        { name: 'p99', value: stats.p99, color: '#ef4444' }, // Rose
    ];

    const CustomTooltip = ({ active, payload, label }: any) => {
        if (active && payload && payload.length) {
            return (
                <div className="rounded-xl border border-zinc-800 bg-zinc-950/90 p-3 shadow-xl backdrop-blur-md">
                    <p className="mb-1 text-sm font-medium text-zinc-300">{label} Latency</p>
                    <p className="text-lg font-bold" style={{ color: payload[0].payload.color }}>
                        {payload[0].value.toFixed(2)} ms
                    </p>
                </div>
            );
        }
        return null;
    };

    return (
        <div className="h-72 w-full pt-4">
            <ResponsiveContainer width="100%" height="100%">
                <BarChart
                    data={data}
                    margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
                >
                    <CartesianGrid strokeDasharray="3 3" stroke="#27272a" vertical={false} />
                    <XAxis
                        dataKey="name"
                        stroke="#71717a"
                        fontSize={12}
                        tickLine={false}
                        axisLine={false}
                        dy={10}
                    />
                    <YAxis
                        stroke="#71717a"
                        fontSize={12}
                        tickLine={false}
                        axisLine={false}
                        tickFormatter={(val) => `${val}ms`}
                    />
                    <Tooltip content={<CustomTooltip />} cursor={{ fill: '#27272a', opacity: 0.4 }} />
                    <Bar dataKey="value" radius={[4, 4, 0, 0]} maxBarSize={60}>
                        {data.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={entry.color} />
                        ))}
                    </Bar>
                </BarChart>
            </ResponsiveContainer>
        </div>
    );
};
