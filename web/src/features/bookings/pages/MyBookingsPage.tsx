import { useMyBookings } from '../hooks';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { CopyButton } from '@/components/CopyButton';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';

export function MyBookingsPage() {
  const { data, isLoading, error } = useMyBookings();
  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <h1 style={{ margin: 0 }}>My Bookings</h1>
      {isLoading && <div>Loading...</div>}
      {error && <div style={{ color: 'var(--danger)' }}>Failed to load bookings.</div>}
      <div style={{ display: 'grid', gap: 10, gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))' }}>
        {data?.items?.length ? (
          data.items.map((b: any) => (
            <Card key={b.id}>
              <CardHeader>
                <CardTitle>{b.date} {b.startTime} - {b.endTime}</CardTitle>
              </CardHeader>
              <div style={{ display: 'grid', gap: 6 }}>
                <div>
                  <small>Code:</small> <strong>{b.bookingCode}</strong> <CopyButton value={b.bookingCode} ariaLabel="Copy booking code" />
                </div>
                <div>
                  <small>Appointment:</small> <strong>{b.appCode}</strong> {b.appCode && <CopyButton value={b.appCode} ariaLabel="Copy appointment code" />}
                </div>
                <div><small>Status:</small> <strong>{b.status}</strong></div>
              </div>
            </Card>
          ))
        ) : (
          <EmptyState>
            <EmptyTitle>No bookings yet</EmptyTitle>
            <EmptyDescription>When you book a slot, it will appear here.</EmptyDescription>
          </EmptyState>
        )}
      </div>
    </div>
  );
}
