import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { format } from 'date-fns';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { useAppointmentUsers, useMyAppointments } from '../hooks';
import { useRejectBooking } from '@/features/bookings/hooks';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';
import { ListItem } from '@/components/ListItem';
import { CopyButton } from '@/components/CopyButton';

export function AppointmentDetailsPage() {
  const { id = '' } = useParams();
  const loc = useLocation() as any;
  const navigate = useNavigate();
  const fromState = loc.state?.appointment;
  const { data: my } = useMyAppointments();
  const match = my?.items?.find((i: any) => i.id === id);
  const appt = fromState || match;
  const { data: users, isLoading, error } = useAppointmentUsers(id);
  const reject = useRejectBooking(appt?.appCode || id);
  const bookingLink = typeof window !== 'undefined' && appt?.appCode ? `${window.location.origin}/book-by-code?code=${appt.appCode}` : '';

  const schedule = formatSchedule(appt);
  const details = [
    { label: 'Type', value: appt?.type ? capitalize(appt.type) : null },
    { label: 'Capacity', value: appt?.maxAttendees ? `${appt.maxAttendees} per slot` : '1 per slot' },
    { label: 'Slot length', value: appt?.bookingDuration ? `${appt.bookingDuration} mins` : null },
    { label: 'Status', value: appt?.status ? capitalize(appt.status) : 'Active' },
  ].filter((item) => item.value);

  return (
    <div style={{ display: 'grid', gap: 16 }}>
      <Card>
        <CardHeader>
          <CardTitle>Appointment Details</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 16 }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>Title</span>
            <strong style={{ fontSize: 20 }}>{appt?.title || 'Untitled appointment'}</strong>
          </div>
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
          {appt?.appCode && (
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
                <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>Appointment code</span>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 18, fontWeight: 700 }}>
                  <code style={{ fontSize: 18 }}>{appt.appCode}</code>
                  <CopyButton value={appt.appCode} ariaLabel="Copy appointment code" />
                </div>
              </div>
              <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
                <Button variant="primary" onClick={() => navigate(`/book-by-code?code=${appt.appCode}`)}>
                  Open booking flow
                </Button>
                {bookingLink && (
                  <CopyButton value={bookingLink} ariaLabel="Copy booking link" label="Copy link" />
                )}
              </div>
            </div>
          )}
        </div>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Registered Users/Bookings</CardTitle>
        </CardHeader>
        {isLoading && <div>Loading users...</div>}
        {error && <div style={{ color: 'var(--danger)' }}>Failed to load users.</div>}
        <div style={{ display: 'grid', gap: 8 }}>
          {users?.items?.map((u: any) => (
            <ListItem key={u.id}>
              <div>
                <div style={{ fontWeight:600 }}>{u.name || u.email || u.phone}</div>
                <div style={{ fontSize:12, opacity:0.8 }}>Code: {u.bookingCode} • {u.status}</div>
              </div>
              <Button variant="ghost" onClick={() => reject.mutate(u.bookingCode)}>Reject</Button>
            </ListItem>
          ))}
          {!users?.items?.length && (
            <EmptyState>
              <EmptyTitle>No bookings yet</EmptyTitle>
              <EmptyDescription>Share the appointment code to start receiving bookings.</EmptyDescription>
            </EmptyState>
          )}
        </div>
      </Card>
    </div>
  );
}

function formatSchedule(appt: any) {
  if (!appt) return null;
  const startDate = parseDate(appt.startDate);
  const endDate = parseDate(appt.endDate);
  const startTime = parseDate(appt.startTime);
  const endTime = parseDate(appt.endTime);
  if (!startDate || !endDate) return null;

  const sameDay = startDate.toDateString() === endDate.toDateString();
  const dayPart = sameDay
    ? format(startDate, 'EEE, MMM d')
    : `${format(startDate, 'EEE, MMM d')} → ${format(endDate, 'EEE, MMM d')}`;
  let timePart = '';
  if (startTime && endTime) {
    const sameTimeDay = startTime.toDateString() === endTime.toDateString();
    timePart = sameTimeDay
      ? `${format(startTime, 'p')} – ${format(endTime, 'p')}`
      : `${format(startTime, 'EEE, MMM d p')} → ${format(endTime, 'EEE, MMM d p')}`;
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
