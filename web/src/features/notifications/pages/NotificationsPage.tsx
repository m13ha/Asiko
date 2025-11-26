import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { PaginatedList } from '@/components/PaginatedList';
import { usePagination } from '@/hooks/usePagination';
import { useMarkAllRead, useNotifications } from '../hooks';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';

export function NotificationsPage() {
  const pagination = usePagination(1, 10);
  const { data, isLoading, error } = useNotifications(pagination.params);
  const markAll = useMarkAllRead();

  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0 }}>Notifications</h1>
        <Button onClick={() => markAll.mutate()} disabled={markAll.isPending}>Mark all read</Button>
      </div>
      <PaginatedList
        data={data}
        isLoading={isLoading}
        error={error}
        onPageChange={pagination.updatePage}
        renderItem={(notification: any) => (
          <Card key={notification.id}>
            <CardHeader>
              <CardTitle style={{ fontSize: 16 }}>{notification.title || 'Notification'}</CardTitle>
            </CardHeader>
            <div style={{ color: 'var(--text-muted)' }}>{notification.message || notification.body || ''}</div>
          </Card>
        )}
        emptyState={
          <EmptyState>
            <EmptyTitle>No notifications yet</EmptyTitle>
            <EmptyDescription>We'll let you know when there's activity.</EmptyDescription>
          </EmptyState>
        }
        itemsClassName="grid gap-2"
      />
    </div>
  );
}