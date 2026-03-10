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
import type { BenchmarkStats } from '../core/types/api';

interface LatencyChartProps {
    stats: BenchmarkStats | null;
}

const COLORS = {
    p50: '#00758F', 
    p95: '#7F2257', 
    p99: '#C74634', 
} as const;

const EMPTY_DATA = [
    { name: 'P50', value: 0, color: COLORS.p50 },
    { name: 'P95', value: 0, color: COLORS.p95 },
    { name: 'P99', value: 0, color: COLORS.p99 },
];

export const LatencyChart: React.FC<LatencyChartProps> = ({ stats }) => {
    const data = useMemo(() => {
        if (!stats) return EMPTY_DATA;
        return [
            { name: 'P50', value: stats.p50, color: COLORS.p50 },
            { name: 'P95', value: stats.p95, color: COLORS.p95 },
            { name: 'P99', value: stats.p99, color: COLORS.p99 },
        ];
    }, [stats]);

    return (
        <div className="h-96 w-full pt-10">
            <ResponsiveContainer width="100%" height="100%">
                <BarChart
                    data={data}
                    margin={{ top: 20, right: 30, left: 10, bottom: 40 }}
                    barGap={20}
                >
                    <defs>
                        {Object.entries(COLORS).map(([key, color]) => (
                            <linearGradient key={key} id={`brand-art-grad-${key}`} x1="0" y1="0" x2="0" y2="1">
                                <stop offset="0%" stopColor={color} stopOpacity={0.9} />
                                <stop offset="100%" stopColor={color} stopOpacity={0.3} />
                            </linearGradient>
                        ))}
                        {/* Glass Overlay Gradient */}
                        <linearGradient id="glassGradient" x1="0" y1="0" x2="1" y2="0">
                            <stop offset="0%" stopColor="white" stopOpacity={0.2} />
                            <stop offset="50%" stopColor="white" stopOpacity={0.05} />
                            <stop offset="100%" stopColor="white" stopOpacity={0.2} />
                        </linearGradient>
                    </defs>
                    
                    <CartesianGrid strokeDasharray="10 10" stroke="rgba(49, 45, 42, 0.08)" vertical={false} />
                    
                    <XAxis
                        dataKey="name"
                        stroke="#312D2A"
                        fontSize={12}
                        fontWeight={900}
                        tickLine={false}
                        axisLine={false}
                        dy={25}
                        letterSpacing="0.4em"
                    />
                    <YAxis
                        stroke="#645F5B"
                        fontSize={11}
                        fontWeight={700}
                        tickLine={false}
                        axisLine={false}
                        tickFormatter={(val) => `${val}ms`}
                        width={75}
                        dx={-10}
                    />
                    
                    <ReferenceLine y={0} stroke="rgba(49, 45, 42, 0.1)" />

                    <Bar
                        dataKey="value"
                        radius={[16, 16, 0, 0]}
                        maxBarSize={80}
                        isAnimationActive={true}
                        animationDuration={2000}
                        animationBegin={300}
                    >
                        {data.map((entry) => (
                            <Cell
                                key={entry.name}
                                fill={`url(#brand-art-grad-${entry.name.toLowerCase()})`}
                                stroke={entry.color}
                                strokeWidth={3}
                                strokeOpacity={0.2}
                            />
                        ))}
                    </Bar>
                </BarChart>
            </ResponsiveContainer>
        </div>
    );
};
