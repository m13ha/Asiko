import { useLocation, useParams } from 'react-router-dom';
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
  const fromState = loc.state?.appointment;
  const { data: my } = useMyAppointments();
  const match = my?.items?.find((i: any) => i.id === id);
  const appt = fromState || match;
  const { data: users, isLoading, error } = useAppointmentUsers(id);
  const reject = useRejectBooking(appt?.appCode || id);

  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <Card>
        <CardHeader>
          <CardTitle>Appointment Details</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 8 }}>
          <div><small>ID:</small> <strong>{id}</strong></div>
          {appt?.title && <div><small>Title:</small> <strong>{appt.title}</strong></div>}
          {appt?.appCode && (
            <div>
              <small>Code:</small> <strong>{appt.appCode}</strong> <CopyButton value={appt.appCode} ariaLabel="Copy appointment code" />
            </div>
          )}
          {appt && (
            <div>
              <small>When:</small> {appt.startDate} {appt.startTime} → {appt.endDate} {appt.endTime}
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
