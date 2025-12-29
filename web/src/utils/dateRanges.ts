export const RANGE_OPTIONS = [
  { value: '7d', label: 'Last 7 days', durationMs: 7 * 24 * 60 * 60 * 1000 },
  { value: '14d', label: 'Last 14 days', durationMs: 14 * 24 * 60 * 60 * 1000 },
  { value: '30d', label: 'Last 30 days', durationMs: 30 * 24 * 60 * 60 * 1000 },
  { value: '60d', label: 'Last 60 days', durationMs: 60 * 24 * 60 * 60 * 1000 },
  { value: '90d', label: 'Last 90 days', durationMs: 90 * 24 * 60 * 60 * 1000 },
  { value: '180d', label: 'Last 6 months', durationMs: 180 * 24 * 60 * 60 * 1000 },
];

function formatDate(d: Date) {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

export function computeRange(value: string) {
  const fallback = RANGE_OPTIONS[0];
  const match = RANGE_OPTIONS.find((option) => option.value === value) ?? fallback;
  const now = new Date();
  const dayMs = 24 * 60 * 60 * 1000;
  const days = Math.max(1, Math.round(match.durationMs / dayMs));
  const start = new Date(now.getTime());
  start.setDate(now.getDate() - (days - 1));
  return {
    startDate: formatDate(start),
    endDate: formatDate(now),
  };
}
