import React, { useMemo } from 'react';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    ResponsiveContainer,
    Cell,
    ReferenceLine,
} from 'recharts';
import type { BenchmarkStats } from '../types/api';

interface LatencyChartProps {
    stats: BenchmarkStats | null;
}

const COLORS = {
    p50: '#10b981',
    p95: '#f59e0b',
    p99: '#ef4444',
} as const;

const EMPTY_DATA = [
    { name: 'p50', value: 0, color: COLORS.p50 },
    { name: 'p95', value: 0, color: COLORS.p95 },
    { name: 'p99', value: 0, color: COLORS.p99 },
];

export const LatencyChart: React.FC<LatencyChartProps> = ({ stats }) => {
    const data = useMemo(() => {
        if (!stats) return EMPTY_DATA;
        return [
            { name: 'p50', value: stats.p50, color: COLORS.p50 },
            { name: 'p95', value: stats.p95, color: COLORS.p95 },
            { name: 'p99', value: stats.p99, color: COLORS.p99 },
        ];
    }, [stats]);

    return (
        <div className="h-72 w-full pt-4">
            <ResponsiveContainer width="100%" height="100%">
                <BarChart
                    data={data}
                    margin={{ top: 10, right: 10, left: -10, bottom: 0 }}
                    barGap={8}
                >
                    <defs>
                        {Object.entries(COLORS).map(([key, color]) => (
                            <linearGradient key={key} id={`grad-${key}`} x1="0" y1="0" x2="0" y2="1">
                                <stop offset="0%" stopColor={color} stopOpacity={0.9} />
                                <stop offset="100%" stopColor={color} stopOpacity={0.4} />
                            </linearGradient>
                        ))}
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="#27272a" vertical={false} />
                    <XAxis
                        dataKey="name"
                        stroke="#71717a"
                        fontSize={12}
                        fontWeight={600}
                        tickLine={false}
                        axisLine={false}
                        dy={10}
                    />
                    <YAxis
                        stroke="#71717a"
                        fontSize={11}
                        tickLine={false}
                        axisLine={false}
                        tickFormatter={(val) => `${val}ms`}
                        width={55}
                    />
                    <ReferenceLine y={0} stroke="#3f3f46" />

                    <Bar
                        dataKey="value"
                        radius={[6, 6, 0, 0]}
                        maxBarSize={56}
                        isAnimationActive={false}
                    >
                        {data.map((entry) => (
                            <Cell
                                key={entry.name}
                                fill={`url(#grad-${entry.name})`}
                                stroke={entry.color}
                                strokeWidth={1}
                                strokeOpacity={0.3}
                            />
                        ))}
                    </Bar>
                </BarChart>
            </ResponsiveContainer>
        </div>
    );
};
