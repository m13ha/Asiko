import { useNavigate, useParams } from 'react-router-dom';
import { format } from 'date-fns';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useBookingByCode, useCancelBooking, useUpdateBooking, useAvailableDates, useAvailableSlotsByDay } from '../hooks';
import { useAppointmentByAppCode } from '@/features/appointments/hooks';
import { useEffect, useMemo, useState } from 'react';
import { Button } from '@/components/Button';
import { CopyButton } from '@/components/CopyButton';
import { AvailabilityCalendar } from '../components/AvailabilityCalendar';
import { InlineTimePicker } from '../components/InlineTimePicker';
import type { EntitiesBooking } from '@appointment-master/api-client';
import * as API from '@appointment-master/api-client';

export function BookingManagePage() {
  const { bookingCode = '' } = useParams();
  const navigate = useNavigate();
  const { data, isLoading, error, refetch } = useBookingByCode(bookingCode);
  const [selectedDate, setSelectedDate] = useState('');
  const [selectedSlot, setSelectedSlot] = useState<EntitiesBooking | null>(null);
  const [attendeeCount, setAttendeeCount] = useState(1);
  const update = useUpdateBooking(bookingCode);
  const cancel = useCancelBooking(bookingCode);
  const appCode = data?.appCode || '';
  const appointmentDetails = useAppointmentByAppCode(appCode);
  const availableDatesQuery = useAvailableDates(appCode);
  const daySlots = useAvailableSlotsByDay(appCode, selectedDate);
  const schedule = formatSchedule(data);
  const details = [
    { label: 'Status', value: data?.status ? capitalize(data.status) : null },
    { label: 'Appointment code', value: data?.appCode || null },
    { label: 'Seats', value: typeof data?.seatsBooked === 'number' ? String(data.seatsBooked) : null },
  ].filter((item) => item.value);

  const bookingDate = useMemo(() => toIsoDay(data?.date), [data?.date]);
  const availableDates = availableDatesQuery.data || [];
  const appointmentType = appointmentDetails.data?.type;
  const statusValue = String(data?.status || '').toLowerCase();
  const canReschedule = appointmentType !== API.EntitiesAppointmentType.Party && statusValue !== 'ongoing';
  const canEditAttendees = appointmentType === API.EntitiesAppointmentType.Group;

  useEffect(() => {
    if (!availableDates.length) return;
    if (selectedDate) return;
    if (bookingDate && availableDates.includes(bookingDate)) {
      setSelectedDate(bookingDate);
      return;
    }
    setSelectedDate(availableDates[0]);
  }, [availableDates, bookingDate, selectedDate]);

  useEffect(() => {
    const initialCount = data?.attendeeCount ?? data?.seatsBooked ?? 1;
    if (initialCount && initialCount > 0) {
      setAttendeeCount(initialCount);
    }
  }, [data?.attendeeCount, data?.seatsBooked]);

  useEffect(() => {
    if (!selectedSlot || !selectedDate) return;
    const slotDate = toIsoDay(selectedSlot.date);
    if (slotDate && slotDate !== selectedDate) {
      setSelectedSlot(null);
    }
  }, [selectedDate, selectedSlot]);

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

            {!canReschedule && (
              <div className="rounded-lg border border-[var(--border)] bg-[var(--bg-elevated)] p-3 text-sm text-[var(--text-muted)]">
                {statusValue === 'ongoing'
                  ? 'Ongoing bookings cannot be rescheduled.'
                  : 'Party appointments can’t be rescheduled to a different slot.'}
              </div>
            )}

            {canReschedule && (
              <div className="grid gap-4">
                {canEditAttendees && (
                  <div className="flex flex-wrap items-center gap-3">
                    <span className="text-sm text-[var(--text-muted)]">Attendees</span>
                    <div className="flex items-center gap-2">
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => setAttendeeCount((count) => Math.max(1, count - 1))}
                        disabled={attendeeCount <= 1}
                      >
                        -
                      </Button>
                      <span className="min-w-[32px] text-center text-sm font-semibold">{attendeeCount}</span>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => setAttendeeCount((count) => count + 1)}
                      >
                        +
                      </Button>
                    </div>
                  </div>
                )}
                <AvailabilityCalendar
                  availableDates={availableDates}
                  selectedDate={selectedDate}
                  onSelect={(value) => setSelectedDate(value)}
                />
                <Card className="!p-4">
                  {daySlots.isFetching ? (
                    <div className="flex items-center justify-center py-6 gap-3">
                      <div className="w-5 h-5 border-2 border-[var(--primary)] border-t-transparent rounded-full animate-spin" />
                      <span className="text-sm text-[var(--text-muted)]">Loading available times...</span>
                    </div>
                  ) : (
                    <InlineTimePicker
                      slots={daySlots.data?.items || []}
                      selectedSlot={selectedSlot}
                      onSelect={setSelectedSlot}
                      appointmentType={appointmentType}
                    />
                  )}
                </Card>
              </div>
            )}

            <div className="flex flex-wrap gap-2">
              <Button
                variant="primary"
                onClick={() => {
                  if (!selectedSlot) return;
                  const normalizedDateIso = selectedDate ? `${selectedDate}T00:00:00Z` : '';
                  update.mutate(
                    {
                      appCode,
                      date: normalizedDateIso,
                      startTime: selectedSlot.startTime || '',
                      endTime: selectedSlot.endTime || '',
                      attendeeCount: canEditAttendees ? attendeeCount : 1,
                      name: data?.name || undefined,
                      email: data?.email || undefined,
                      phone: data?.phone || undefined,
                    },
                    {
                      onSuccess: () => {
                        window.setTimeout(() => {
                          setSelectedSlot(null);
                          setSelectedDate('');
                          setAttendeeCount(1);
                          refetch();
                        }, 400);
                      },
                    }
                  );
                }}
                disabled={!canReschedule || !selectedDate || !selectedSlot || attendeeCount < 1 || update.isPending}
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

function toIsoDay(value?: string) {
  if (!value) return '';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return '';
  return format(d, 'yyyy-MM-dd');
}

function capitalize(value?: string) {
  if (!value) return '';
  return value.charAt(0).toUpperCase() + value.slice(1);
}
