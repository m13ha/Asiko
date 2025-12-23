import { useNavigate, useParams } from 'react-router-dom';
import { format } from 'date-fns';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useBookingByCode, useCancelBooking, useUpdateBooking } from '../hooks';
import { Input } from '@/components/Input';
import { useState } from 'react';
import { Button } from '@/components/Button';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { CopyButton } from '@/components/CopyButton';

export function BookingManagePage() {
  const { bookingCode = '' } = useParams();
  const navigate = useNavigate();
  const { data, isLoading, error, refetch } = useBookingByCode(bookingCode);
  const [date, setDate] = useState('');
  const [startTime, setStart] = useState('');
  const [endTime, setEnd] = useState('');
  const update = useUpdateBooking(bookingCode);
  const cancel = useCancelBooking(bookingCode);
  const schedule = formatSchedule(data);
  const details = [
    { label: 'Status', value: data?.status ? capitalize(data.status) : null },
    { label: 'Appointment code', value: data?.appCode || null },
    { label: 'Seats', value: typeof data?.seatsBooked === 'number' ? String(data.seatsBooked) : null },
  ].filter((item) => item.value);

  return (
    <div style={{ display: 'grid', gap: 16 }}>
      <div>
        <Button variant="ghost" onClick={() => navigate(-1)}>
          Back
        </Button>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Booking Details</CardTitle>
        </CardHeader>
        {isLoading && <div>Loading...</div>}
        {error && <div style={{ color: 'var(--danger)' }}>Failed to load booking.</div>}
        {data && (
          <div style={{ display: 'grid', gap: 16 }}>
            {data?.bookingCode && (
              <div style={{
                border: '1px solid var(--border)',
                borderRadius: 'var(--radius)',
                padding: '12px 14px',
                display: 'flex',
                flexWrap: 'wrap',
                gap: 12,
                alignItems: 'center',
                justifyContent: 'space-between',
              }}>
                <div>
                  <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>Booking code</span>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 18, fontWeight: 700 }}>
                    <code style={{ fontSize: 18 }}>{data.bookingCode}</code>
                    <CopyButton value={data?.bookingCode || ''} ariaLabel="Copy booking code" />
                  </div>
                </div>
              </div>
            )}

            {schedule && (
              <div style={{ background: 'color-mix(in oklab, var(--primary) 6%, transparent)', borderRadius: 'var(--radius)', padding: '12px 14px' }}>
                <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>When</span>
                <div style={{ fontWeight: 600 }}>{schedule}</div>
              </div>
            )}

            {details.length > 0 && (
              <div style={{ display: 'grid', gap: 12, gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))' }}>
                {details.map((item) => (
                  <div key={item.label} style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                    <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>{item.label}</span>
                    <strong>{item.value}</strong>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </Card>

      {data?.status !== 'expired' && (
        <Card>
          <CardHeader>
            <CardTitle>Reschedule</CardTitle>
          </CardHeader>
          <div className="grid gap-4">
            <div className="rounded-lg border border-[var(--border)] bg-[var(--bg-subtle)] p-3 text-xs text-[var(--text-muted)]">
              Pick a new date and time to move this booking. We’ll keep your original booking code.
            </div>

            <div className="grid gap-3 md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)]">
              <Field>
                <FieldLabel>New date</FieldLabel>
                <FieldRow>
                  <div className="relative">
                    <IconSlot><i className="pi pi-calendar" aria-hidden="true" /></IconSlot>
                    <Input type="date" value={date} onChange={(e) => setDate(e.target.value)} className="pl-9" />
                  </div>
                </FieldRow>
              </Field>
              <Field>
                <FieldLabel>Start time</FieldLabel>
                <FieldRow>
                  <div className="relative">
                    <IconSlot><i className="pi pi-clock" aria-hidden="true" /></IconSlot>
                    <Input type="time" value={startTime} onChange={(e) => setStart(e.target.value)} className="pl-9" />
                  </div>
                </FieldRow>
              </Field>
              <Field>
                <FieldLabel>End time</FieldLabel>
                <FieldRow>
                  <div className="relative">
                    <IconSlot><i className="pi pi-clock" aria-hidden="true" /></IconSlot>
                    <Input type="time" value={endTime} onChange={(e) => setEnd(e.target.value)} className="pl-9" />
                  </div>
                </FieldRow>
              </Field>
            </div>

            <div className="flex flex-wrap gap-2">
              <Button
                variant="primary"
                onClick={() => update.mutate({ appCode: data?.appCode || '', date, startTime, endTime }, { onSuccess: () => refetch() })}
                disabled={!date || !startTime || !endTime || update.isPending}
              >
                {update.isPending ? 'Updating…' : 'Update booking'}
              </Button>
              <Button onClick={() => cancel.mutate(undefined, { onSuccess: () => refetch() })} disabled={cancel.isPending}>
                {cancel.isPending ? 'Cancelling…' : 'Cancel booking'}
              </Button>
            </div>
          </div>
        </Card>
      )}
    </div>
  );
}

function formatSchedule(booking: any) {
  if (!booking) return null;
  const date = parseDate(booking.date);
  const start = parseDate(booking.startTime);
  const end = parseDate(booking.endTime);
  if (!date) return null;

  const dayPart = format(date, 'EEE, MMM d, yyyy');
  let timePart = '';
  if (start && end) {
    timePart = `${format(start, 'p')} – ${format(end, 'p')}`;
  } else if (start) {
    timePart = format(start, 'p');
  }
  const tz = Intl.DateTimeFormat().resolvedOptions().timeZone || 'local time';
  return [dayPart, timePart].filter(Boolean).join(' • ') + ` (${tz})`;
}

function parseDate(value?: string) {
  if (!value) return null;
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? null : d;
}

function capitalize(value?: string) {
  if (!value) return '';
  return value.charAt(0).toUpperCase() + value.slice(1);
}
