import type { ResponsesTimeSeriesPoint } from '@appointment-master/api-client';

function pathFromSeries(series: ResponsesTimeSeriesPoint[], w: number, h: number, pad = 2) {
  const pts = (series || []).filter(s => typeof s.count === 'number');
  if (pts.length === 0) return '';
  const xs = pts.map((_, i) => i);
  const ys = pts.map(p => p.count as number);
  const xMin = 0, xMax = Math.max(pts.length - 1, 1);
  const yMin = Math.min(...ys);
  const yMax = Math.max(...ys);
  const xScale = (i: number) => pad + (i - xMin) / (xMax - xMin || 1) * (w - pad * 2);
  const yScale = (v: number) => pad + (1 - (v - yMin) / (yMax - yMin || 1)) * (h - pad * 2);
  const d = xs.map((i, idx) => `${idx === 0 ? 'M' : 'L'}${xScale(i)},${yScale(ys[idx])}`).join(' ');
  return d;
}

export function Sparkline({ series, colorVar = '--primary', width = 140, height = 36, title = 'sparkline' }: { series: ResponsesTimeSeriesPoint[]; colorVar?: string; width?: number; height?: number; title?: string; }) {
  const d = pathFromSeries(series || [], width, height, 2);
  const color = getComputedStyleColor(colorVar);
  return (
    <svg role="img" viewBox={`0 0 ${width} ${height}`} width={width} height={height} aria-label={title} preserveAspectRatio="none">
      <title>{title}</title>
      <path d={d} fill="none" stroke={color} strokeWidth={2} strokeLinejoin="round" strokeLinecap="round" />
    </svg>
  );
}

function getComputedStyleColor(varName: string) {
  if (typeof window === 'undefined') return '#146C43';
  const c = getComputedStyle(document.documentElement).getPropertyValue(varName).trim();
  return c || '#146C43';
}

