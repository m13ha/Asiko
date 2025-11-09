import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { useMarkAllRead, useNotifications } from '../hooks';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';

export function NotificationsPage() {
  const { data, isLoading, error } = useNotifications();
  const markAll = useMarkAllRead();
  const items = data?.items || [];

  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0 }}>Notifications</h1>
        <Button onClick={() => markAll.mutate()} disabled={markAll.isPending}>Mark all read</Button>
      </div>
      {isLoading && <div>Loading...</div>}
      {error && <div style={{ color: 'var(--danger)' }}>Failed to load notifications.</div>}
      <div style={{ display: 'grid', gap: 8 }}>
        {items.length === 0 && (
          <EmptyState>
            <EmptyTitle>No notifications yet</EmptyTitle>
            <EmptyDescription>We’ll let you know when there’s activity.</EmptyDescription>
          </EmptyState>
        )}
        {items.map((n: any, i: number) => (
          <Card key={n.id || i}>
            <CardHeader>
              <CardTitle style={{ fontSize: 16 }}>{n.title || 'Notification'}</CardTitle>
            </CardHeader>
            <div style={{ color: 'var(--text-muted)' }}>{n.message || n.body || ''}</div>
          </Card>
        ))}
      </div>
    </div>
  );
}
