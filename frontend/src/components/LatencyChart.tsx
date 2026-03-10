import React, { useMemo } from 'react';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    ResponsiveContainer,
    Cell,
} from 'recharts';
import type { BenchmarkStats } from '../core/types/api';

interface LatencyChartProps {
    stats: BenchmarkStats | null;
}

const COLORS = {
    p50: '#00758F', // Viz Teal
    p95: '#4A6273', // Viz Slate
    p99: '#C74634', // Oracle Red
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
            { name: 'P50', value: stats.p50, color: COLORS.p50, label: 'Median' },
            { name: 'P95', value: stats.p95, color: COLORS.p95, label: '95th Percentile' },
            { name: 'P99', value: stats.p99, color: COLORS.p99, label: '99th Percentile' },
        ];
    }, [stats]);

    return (
        <div className="h-80 w-full">
            <ResponsiveContainer width="100%" height="100%">
                <BarChart
                    data={data}
                    margin={{ top: 10, right: 10, left: 0, bottom: 20 }}
                    barSize={60}
                >
                    <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#EBE9E7" />
                    
                    <XAxis
                        dataKey="name"
                        axisLine={false}
                        tickLine={false}
                        tick={{ fill: '#645F5B', fontSize: 12, fontWeight: 600 }}
                        dy={10}
                    />
                    <YAxis
                        axisLine={false}
                        tickLine={false}
                        tick={{ fill: '#645F5B', fontSize: 11 }}
                        tickFormatter={(val) => `${val}ms`}
                        width={45}
                    />
                    
                    <Bar
                        dataKey="value"
                        radius={[4, 4, 0, 0]}
                        isAnimationActive={true}
                    >
                        {data.map((entry, index) => (
                            <Cell 
                                key={`cell-${index}`} 
                                fill={entry.color} 
                                fillOpacity={0.8}
                                className="hover:fill-opacity-100 transition-all duration-300"
                            />
                        ))}
                    </Bar>
                </BarChart>
            </ResponsiveContainer>
            
            {/* Custom Legend */}
            <div className="flex justify-center gap-8 mt-6">
                {data.map((entry) => (
                    <div key={entry.name} className="flex items-center gap-2">
                        <div className="w-3 h-3 rounded-sm" style={{ backgroundColor: entry.color }} />
                        <span className="text-[11px] font-bold text-text-sub uppercase tracking-wider">{entry.label}</span>
                    </div>
                ))}
            </div>
        </div>
    );
};
