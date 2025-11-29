import { useMemo, useState } from 'react';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useUserAnalytics } from '@/features/analytics/hooks';
import { FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Select } from '@/components/Select';
import { ChartCard, ChartTitle, ChartArea, ChartLegend, LegendItem, Swatch } from '@/components/ChartFrame';
import { LineChart } from '@/components/LineChart';
import { Sparkline } from '@/components/Sparkline';
import { getCssVarValue } from '@/components/chartConfig';
import { RANGE_OPTIONS, computeRange } from '@/utils/dateRanges';

export function HomeInsightsPage() {
  const defaultRange = useMemo(() => RANGE_OPTIONS[0].value, []);
  const initialRange = useMemo(() => computeRange(defaultRange), [defaultRange]);
  const [selectedRange, setSelectedRange] = useState(defaultRange);
  const [start, setStart] = useState(initialRange.startDate);
  const [end, setEnd] = useState(initialRange.endDate);
  const { data, isLoading, error } = useUserAnalytics(start, end);

  const handleRangeChange = (value: string) => {
    setSelectedRange(value);
    const next = computeRange(value);
    setStart(next.startDate);
    setEnd(next.endDate);
  };

  return (
    <div style={{ display: 'grid', gap: 16 }}>
      <h1 style={{ margin: 0 }}>Overview</h1>

      <Card>
        <CardHeader style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 12, flexWrap: 'wrap' }}>
          <CardTitle style={{ margin: 0 }}>Performance Snapshot</CardTitle>
          <div style={{ minWidth: 220 }}>
            <FieldLabel style={{ fontSize: 12, textTransform: 'uppercase', color: 'var(--text-muted)' }}>Range</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative', width: '100%' }}>
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
        </CardHeader>
        <div style={{ display: 'grid', gap: 12, gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))' }}>
          <div style={{ border: '1px solid var(--border)', borderRadius: 'var(--radius)', padding: '12px 14px' }}>
            <small style={{ color: 'var(--text-muted)' }}>Total appointments</small>
            <div style={{ fontSize: 28, fontWeight: 700 }}>
              {isLoading ? '—' : data?.totalAppointments ?? 0}
            </div>
          </div>
          <div style={{ border: '1px solid var(--border)', borderRadius: 'var(--radius)', padding: '12px 14px' }}>
            <small style={{ color: 'var(--text-muted)' }}>Total bookings</small>
            <div style={{ fontSize: 28, fontWeight: 700 }}>
              {isLoading ? '—' : data?.totalBookings ?? 0}
            </div>
            {!!(data?.bookingsPerDay && data.bookingsPerDay.length) && (
              <div style={{ marginTop: 8 }}>
                <Sparkline series={data.bookingsPerDay} />
              </div>
            )}
          </div>
        </div>
        {error && <div style={{ color: 'var(--danger)' }}>Failed to load analytics.</div>}
      </Card>

      <ChartCard>
        <ChartTitle>Bookings per day</ChartTitle>
        <ChartArea>
          <LineChart
            series={data?.bookingsPerDay || []}
            secondary={data?.cancellationsPerDay || []}
            primaryLabel="Bookings"
            secondaryLabel="Cancellations"
          />
        </ChartArea>
        <ChartLegend>
          <LegendItem><Swatch color={getCssVarValue('--primary', '#146C43')} /> Bookings</LegendItem>
          <LegendItem><Swatch color={getCssVarValue('--secondary', '#2EB872')} /> Cancellations</LegendItem>
        </ChartLegend>
      </ChartCard>
    </div>
  );
}
