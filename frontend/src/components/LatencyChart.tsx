import React, { useMemo } from 'react';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    ResponsiveContainer,
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
    { name: '', P50: 0, P95: 0, P99: 0 },
];

export const LatencyChart: React.FC<LatencyChartProps> = ({ stats }) => {
    const data = useMemo(() => {
        if (!stats) return EMPTY_DATA;
        return [
            { name: '', P50: stats.p50, P95: stats.p95, P99: stats.p99 }
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
                    
                    <Bar dataKey="P50" fill={COLORS.p50} radius={[4, 4, 0, 0]} isAnimationActive={true} barSize={60} />
                    <Bar dataKey="P95" fill={COLORS.p95} radius={[4, 4, 0, 0]} isAnimationActive={true} barSize={60} />
                    <Bar dataKey="P99" fill={COLORS.p99} radius={[4, 4, 0, 0]} isAnimationActive={true} barSize={60} />
                </BarChart>
            </ResponsiveContainer>
            
            {/* Custom Legend */}
            <div className="flex justify-center gap-8 mt-6">
                <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-sm" style={{ backgroundColor: COLORS.p50 }} />
                    <span className="text-[11px] font-bold text-text-sub uppercase tracking-wider">Median</span>
                </div>
                <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-sm" style={{ backgroundColor: COLORS.p95 }} />
                    <span className="text-[11px] font-bold text-text-sub uppercase tracking-wider">95th Percentile</span>
                </div>
                <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-sm" style={{ backgroundColor: COLORS.p99 }} />
                    <span className="text-[11px] font-bold text-text-sub uppercase tracking-wider">99th Percentile</span>
                </div>
            </div>
        </div>
    );
};
