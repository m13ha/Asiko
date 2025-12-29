import { useMemo, useState } from 'react';
import { Card, CardTitle } from '@/components/Card';
import { useUserAnalytics } from '@/features/analytics/hooks';
import { FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Select } from '@/components/Select';
import { ChartCard, ChartTitle, ChartArea, ChartLegend, LegendItem, Swatch } from '@/components/ChartFrame';
import { LineChart } from '@/components/LineChart';
import { getCssVarValue } from '@/components/chartConfig';
import { RANGE_OPTIONS, computeRange } from '@/utils/dateRanges';

type KpiCardProps = {
  label: string;
  value: string | number;
  hint?: string;
  sparkline?: React.ReactNode;
  className?: string;
};

function KpiCard({ label, value, hint, sparkline, className = '' }: KpiCardProps) {
  return (
    <Card className={`p-4 min-w-0 overflow-hidden ${className}`}>
      <div className="text-xs uppercase tracking-wide text-[var(--text-muted)]">{label}</div>
      <div className="mt-2 text-2xl font-semibold">{value}</div>
      {hint ? <div className="mt-1 text-sm text-[var(--text-muted)] break-words">{hint}</div> : null}
      {sparkline ? <div className="mt-3 overflow-hidden">{sparkline}</div> : null}
    </Card>
  );
}

function formatPercent(value: number) {
  if (!Number.isFinite(value)) return '—';
  return `${Math.round(value * 10) / 10}%`;
}

export function HomeInsightsPage() {
  const defaultRange = useMemo(() => RANGE_OPTIONS[0].value, []);
  const initialRange = useMemo(() => computeRange(defaultRange), [defaultRange]);
  const [selectedRange, setSelectedRange] = useState(defaultRange);
  const [start, setStart] = useState(initialRange.startDate);
  const [end, setEnd] = useState(initialRange.endDate);
  const { data, isLoading, error } = useUserAnalytics(start, end);

  const rangeLabel = useMemo(
    () => RANGE_OPTIONS.find(option => option.value === selectedRange)?.label ?? 'Range',
    [selectedRange]
  );

  const totalAppointments = isLoading ? '—' : (data?.totalAppointments ?? 0);
  const totalBookings = isLoading ? '—' : (data?.totalBookings ?? 0);

  const peak = useMemo(() => {
    const bookings = data?.bookingsPerDay ?? [];
    let best: { date: string; count: number } | null = null;
    for (const point of bookings) {
      if (typeof point?.date !== 'string' || typeof point?.count !== 'number') continue;
      if (!best || point.count > best.count) best = { date: point.date, count: point.count };
    }
    return best;
  }, [data?.bookingsPerDay]);

  const handleRangeChange = (value: string) => {
    setSelectedRange(value);
    const next = computeRange(value);
    setStart(next.startDate);
    setEnd(next.endDate);
  };

  return (
    <div className="grid gap-4 overflow-x-hidden">
      <div className="w-full flex items-start justify-between gap-3 flex-wrap">
        <div>
          <h1 className="m-0 text-2xl font-semibold">Overview</h1>
          <div className="mt-1 text-sm text-[var(--text-muted)]">{rangeLabel}</div>
        </div>

        <div className="min-w-0">
          <FieldLabel className="text-xs uppercase tracking-wide text-[var(--text-muted)] hidden md:flex">Range</FieldLabel>
          <FieldRow>
            <div className="relative w-full">
              <IconSlot>
                <i className="pi pi-clock" aria-hidden="true" />
              </IconSlot>
              <Select
                value={selectedRange}
                options={RANGE_OPTIONS}
                optionLabel="label"
                optionValue="value"
                onChange={(event) => handleRangeChange(event.value)}
                pt={{ input: { style: { paddingLeft: '36px' } } }}
              />
            </div>
          </FieldRow>
        </div>
      </div>


      <div className="grid gap-4">
        <div className="lg:col-span-2 min-w-0">
          <ChartCard>
            <ChartTitle>Bookings per day</ChartTitle>
            <ChartArea>
              <LineChart
                series={data?.bookingsPerDay || []}
                secondary={data?.cancellationsPerDay || []}
                primaryLabel="Bookings"
                secondaryLabel="Cancellations"
                rangeStart={start}
                rangeEnd={end}
              />
            </ChartArea>
            <ChartLegend>
              <LegendItem><Swatch color={getCssVarValue('--primary', '#146C43')} /> Bookings</LegendItem>
              <LegendItem><Swatch color={getCssVarValue('--secondary', '#2EB872')} /> Cancellations</LegendItem>
            </ChartLegend>
          </ChartCard>
        </div>
      </div>

      <section aria-label="Key metrics">
        <div className="flex flex-wrap gap-3 ">
          <KpiCard className="flex-1 min-w-[160px] sm:min-w-[220px]" label="Total appointments" value={totalAppointments} />
          <KpiCard
            className="flex-1 min-w-[160px] sm:min-w-[220px]"
            label="Total bookings"
            value={totalBookings}
          />
          <KpiCard
            className="flex-1 min-w-[160px] sm:min-w-[220px]"
            label="Cancellation rate"
            value={isLoading ? '—' : formatPercent(data?.cancellationRate ?? 0)}
            hint={isLoading ? undefined : `${data?.totalCancellations ?? 0} cancellations`}
          />
          <KpiCard
            className="flex-1 min-w-[160px] sm:min-w-[220px]"
            label="Avg bookings / day"
            value={isLoading ? '—' : (Math.round((data?.avgBookingsPerDay ?? 0) * 10) / 10)}
          />
        </div>
      </section>

      {error && <div className="text-[var(--danger)] text-sm">Failed to load analytics.</div>}
    </div>
  );
}
