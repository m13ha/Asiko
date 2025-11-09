import type { ResponsesTimeSeriesPoint } from '@appointment-master/api-client';

function pathFromSeries(series: ResponsesTimeSeriesPoint[], w: number, h: number, pad = 12) {
  const pts = (series || []).filter(s => typeof s.count === 'number' && s.date);
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

export function LineChart({
  series,
  secondary,
  title = 'Chart',
  desc,
}: {
  series: ResponsesTimeSeriesPoint[];
  secondary?: ResponsesTimeSeriesPoint[];
  title?: string;
  desc?: string;
}) {
  const w = 640, h = 220;
  const d1 = pathFromSeries(series || [], w, h);
  const d2 = secondary ? pathFromSeries(secondary || [], w, h) : '';
  const has = !!d1;
  return (
    <svg role="img" viewBox={`0 0 ${w} ${h}`} width="100%" height="100%" aria-label={title} preserveAspectRatio="none">
      <title>{title}</title>
      {desc ? <desc>{desc}</desc> : null}
      {/* grid */}
      <defs>
        <pattern id="grid" width="32" height="24" patternUnits="userSpaceOnUse">
          <rect width="100%" height="100%" fill="transparent" />
          <path d={`M32 0 V24 M0 24 H32`} stroke="var(--border)" strokeWidth="1" opacity="0.6" />
        </pattern>
      </defs>
      <rect x="0" y="0" width="100%" height="100%" fill="url(#grid)" />
      {has ? (
        <>
          {d2 && <path d={d2} fill="none" stroke={getComputedStyleColor('--secondary')} strokeWidth={2} strokeLinejoin="round" strokeLinecap="round" />}
          <path d={d1} fill="none" stroke={getComputedStyleColor('--primary')} strokeWidth={2.5} strokeLinejoin="round" strokeLinecap="round" />
        </>
      ) : (
        <text x="50%" y="50%" dominantBaseline="middle" textAnchor="middle" fill="var(--text-muted)" fontSize="12">No data</text>
      )}
    </svg>
  );
}

function getComputedStyleColor(varName: string) {
  if (typeof window === 'undefined') return '#146C43';
  const c = getComputedStyle(document.documentElement).getPropertyValue(varName).trim();
  return c || '#146C43';
}

