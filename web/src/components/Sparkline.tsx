import type { ResponsesTimeSeriesPoint } from '@appointment-master/api-client';
import { useMemo } from 'react';
import { Chart } from 'primereact/chart';
import { getCssVarValue } from './chartConfig';

type SparklineProps = {
  series: ResponsesTimeSeriesPoint[];
  colorVar?: string;
  height?: number;
};

export function Sparkline({ series, colorVar = '--primary', height = 60 }: SparklineProps) {
  const { labels, dataPoints } = useMemo(() => {
    const normalized = (series || [])
      .filter(point => typeof point?.count === 'number' && typeof point?.date === 'string')
      .map(point => ({ date: point.date as string, count: point.count as number }));
    return {
      labels: normalized.map(point => point.date),
      dataPoints: normalized.map(point => point.count),
    };
  }, [series]);

  if (!labels.length) {
    return (
      <div className="chart-empty" style={{ height }}>
        No data
      </div>
    );
  }

  const color = getCssVarValue(colorVar, '#146C43');

  return (
    <div style={{ width: '100%', height }}>
      <Chart
        type="line"
        data={{
          labels,
          datasets: [
            {
              data: dataPoints,
              borderColor: color,
              backgroundColor: color,
              borderWidth: 2,
              fill: false,
              tension: 0.35,
              pointRadius: 0,
            },
          ],
        }}
        options={{
          responsive: true,
          maintainAspectRatio: false,
          plugins: { legend: { display: false }, tooltip: { enabled: true } },
          scales: {
            x: { display: false },
            y: { display: false },
          },
        }}
      />
    </div>
  );
}
