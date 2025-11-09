import { useMemo, useState } from 'react';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useUserAnalytics } from '../hooks';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Select } from '@/components/Select';
import { Clock } from 'lucide-react';
import { ChartCard, ChartTitle, ChartArea, ChartLegend, LegendItem, Swatch } from '@/components/ChartFrame';
import { LineChart } from '@/components/LineChart';
import { Sparkline } from '@/components/Sparkline';

function formatDate(d: Date) {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

type RangeOption = {
  value: string;
  label: string;
  durationMs: number;
};

const RANGE_OPTIONS: RangeOption[] = [
  { value: '5m', label: 'Last 5 minutes', durationMs: 5 * 60 * 1000 },
  { value: '1h', label: 'Last 1 hour', durationMs: 60 * 60 * 1000 },
  { value: '5h', label: 'Last 5 hours', durationMs: 5 * 60 * 60 * 1000 },
  { value: '1d', label: 'Last day', durationMs: 24 * 60 * 60 * 1000 },
  { value: '3d', label: 'Last 3 days', durationMs: 3 * 24 * 60 * 60 * 1000 },
  { value: '7d', label: 'Last 7 days', durationMs: 7 * 24 * 60 * 60 * 1000 },
  { value: '14d', label: 'Last 14 days', durationMs: 14 * 24 * 60 * 60 * 1000 },
  { value: '30d', label: 'Last 30 days', durationMs: 30 * 24 * 60 * 60 * 1000 },
];

function computeRange(value: string) {
  const fallback = RANGE_OPTIONS[RANGE_OPTIONS.length - 1];
  const match = RANGE_OPTIONS.find((option) => option.value === value) ?? fallback;
  const now = new Date();
  const start = new Date(now.getTime() - match.durationMs);
  return {
    startDate: formatDate(start),
    endDate: formatDate(now),
  };
}

export function AnalyticsPage() {
  const defaultRange = useMemo(() => RANGE_OPTIONS[RANGE_OPTIONS.length - 1].value, []);
  const initialRange = useMemo(() => computeRange(defaultRange), [defaultRange]);
  const [selectedRange, setSelectedRange] = useState(defaultRange);
  const [start, setStart] = useState(initialRange.startDate);
  const [end, setEnd] = useState(initialRange.endDate);
  const { data, isLoading, error } = useUserAnalytics(start, end);
  const getVar = (name: string) => (typeof window !== 'undefined' ? getComputedStyle(document.documentElement).getPropertyValue(name).trim() : '#146C43');

  const handleRangeChange = (value: string) => {
    setSelectedRange(value);
    const next = computeRange(value);
    setStart(next.startDate);
    setEnd(next.endDate);
  };

  return (
    <div style={{ display: 'grid', gap: 16 }}>
      <h1 style={{ margin: 0 }}>Analytics</h1>
      <Card>
        <CardHeader>
          <CardTitle>Time Range</CardTitle>
        </CardHeader>
        <div
          style={{
            display: 'grid',
            gap: 12,
            gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
          }}
        >
          <Field>
            <FieldLabel>Range</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative', width: '100%' }}>
                <IconSlot>
                  <Clock size={16} />
                </IconSlot>
                <Select value={selectedRange} onChange={(event) => handleRangeChange(event.target.value)} style={{ paddingLeft: 36 }}>
                  {RANGE_OPTIONS.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </Select>
              </div>
            </FieldRow>
          </Field>
        </div>
      </Card>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: 12 }}>
        <Card>
          <CardHeader>
            <CardTitle>Total Appointments</CardTitle>
          </CardHeader>
          <div style={{ fontSize: 28, fontWeight: 700 }}>
            {isLoading ? '—' : data?.totalAppointments ?? 0}
          </div>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Total Bookings</CardTitle>
          </CardHeader>
          <div style={{ display: 'grid', gap: 6 }}>
            <div style={{ fontSize: 28, fontWeight: 700 }}>
              {isLoading ? '—' : data?.totalBookings ?? 0}
            </div>
            {!!(data?.bookingsPerDay && data.bookingsPerDay.length) && (
              <div style={{ opacity: 0.9 }}>
                <Sparkline series={data.bookingsPerDay} title="Bookings sparkline" />
              </div>
            )}
          </div>
        </Card>
      </div>

      {error && <div style={{ color: 'var(--danger)' }}>Failed to load analytics.</div>}

      {/* Bookings trend */}
      <ChartCard>
        <ChartTitle>Bookings per day</ChartTitle>
        <ChartArea>
          <LineChart
            series={data?.bookingsPerDay || []}
            secondary={data?.cancellationsPerDay || []}
            title="Bookings per day"
            desc="Primary line shows bookings; secondary shows cancellations"
          />
        </ChartArea>
        <ChartLegend>
          <LegendItem><Swatch color={getVar('--primary')} /> Bookings</LegendItem>
          <LegendItem><Swatch color={getVar('--secondary')} /> Cancellations</LegendItem>
        </ChartLegend>
      </ChartCard>
    </div>
  );
}
