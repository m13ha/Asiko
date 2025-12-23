import type { ResponsesTimeSeriesPoint } from '@appointment-master/api-client';
import { Chart } from 'primereact/chart';
import { useMemo } from 'react';
import { getCssVarValue } from './chartConfig';

type SeriesPoint = {
  date: string;
  count: number;
};

type LineChartProps = {
  series: ResponsesTimeSeriesPoint[];
  secondary?: ResponsesTimeSeriesPoint[];
  primaryLabel?: string;
  secondaryLabel?: string;
  height?: number;
};

function normalizeSeries(points: ResponsesTimeSeriesPoint[] = []): SeriesPoint[] {
  return (points || [])
    .filter((point): point is ResponsesTimeSeriesPoint & { date: string; count: number } => {
      return typeof point?.date === 'string' && typeof point?.count === 'number';
    })
    .map(point => ({ date: point.date as string, count: point.count as number }));
}

export function LineChart({
  series,
  secondary,
  primaryLabel = 'Primary',
  secondaryLabel = 'Secondary',
  height = 220,
}: LineChartProps) {
  const { labels, datasets } = useMemo(() => {
    const primary = normalizeSeries(series);
    const fallback = normalizeSeries(secondary);
    const base = primary.length ? primary : fallback;
    const labelSet = base.map(point => point.date);
    const primaryMap = new Map(primary.map(point => [point.date, point.count]));
    const secondaryMap = new Map(normalizeSeries(secondary).map(point => [point.date, point.count]));

    const primaryColor = getCssVarValue('--primary', '#146C43');
    const secondaryColor = getCssVarValue('--secondary', '#2EB872');

    const datasetPrimary = labelSet.map(label => primaryMap.get(label) ?? 0);
    const datasetSecondary = labelSet.map(label => secondaryMap.get(label) ?? 0);

    const ds = [];

    if (primaryMap.size) {
      ds.push({
        label: primaryLabel,
        data: datasetPrimary,
        borderColor: primaryColor,
        backgroundColor: primaryColor,
        fill: false,
        tension: 0.35,
        pointRadius: 2,
      });
    }

    if (secondaryMap.size) {
      ds.push({
        label: secondaryLabel,
        data: datasetSecondary,
        borderColor: secondaryColor,
        backgroundColor: secondaryColor,
        borderDash: [6, 6],
        fill: false,
        tension: 0.35,
        pointRadius: 2,
      });
    }

    return { labels: labelSet, datasets: ds };
  }, [series, secondary, primaryLabel, secondaryLabel]);

  if (!labels.length) {
    return (
      <div className="grid place-items-center text-[var(--text-muted)] border border-dashed border-[var(--border)] rounded-lg text-sm" style={{ height }}>
        No data
      </div>
    );
  }

  return (
    <div style={{ width: '100%', height }}>
      <Chart
        type="line"
        data={{ labels, datasets }}
        options={{
          maintainAspectRatio: false,
          responsive: true,
          plugins: {
            legend: { display: datasets.length > 1, position: 'bottom' },
            tooltip: { intersect: false, mode: 'index' },
          },
          scales: {
            x: {
              grid: { drawOnChartArea: false },
              ticks: { autoSkip: true, maxTicksLimit: 6 },
            },
            y: {
              beginAtZero: true,
              ticks: { precision: 0 },
              grid: { color: 'rgba(148, 163, 184, 0.25)' },
            },
          },
          elements: { point: { radius: 2, hoverRadius: 4 } },
        }}
      />
    </div>
  );
}
