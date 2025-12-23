import { Button } from '@/components/Button';
import { PaginatedList } from '@/components/PaginatedList';
import { usePagination } from '@/hooks/usePagination';
import { useMarkAllRead, useNotifications } from '../hooks';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';
import { Bell, CheckCheck, Clock } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export function NotificationsPage() {
  const pagination = usePagination(1, 10);
  const { data, isLoading, error } = useNotifications(pagination.params);
  const markAll = useMarkAllRead();

  return (
    <div className="mx-auto py-8 px-4 max-w-4xl">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-6 gap-4">
        <div>
          <h1 className="text-2xl font-bold text-[var(--text)] m-0">Notifications</h1>
          <p className="text-[var(--text-muted)] mt-1">Stay updated with your latest activity</p>
        </div>
        <Button 
          onClick={() => markAll.mutate()} 
          disabled={markAll.isPending || (data?.items?.every(n => n.is_read) ?? true)}
          variant="outline"
          size="sm"
          className="flex items-center gap-2"
        >
          <CheckCheck size={14} />
          Mark all read
        </Button>
      </div>

      <PaginatedList
        data={data}
        isLoading={isLoading}
        error={error}
        onPageChange={pagination.updatePage}
        renderItem={(notification: any) => (
          <div 
            key={notification.id} 
            className={`
              relative p-4 transition-colors border-b border-[var(--border)] last:border-b-0
              ${notification.is_read 
                ? 'bg-[var(--bg-elevated)]' 
                : 'bg-[color-mix(in_oklab,var(--primary)_3%,var(--bg-elevated))] border-l-4 border-l-[var(--primary)]'}
              hover:bg-[color-mix(in_oklab,var(--primary)_5%,var(--bg-elevated))]
            `}
          >
            <div className="flex gap-4">
              <div className={`p-2 rounded-full h-fit ${notification.is_read ? 'bg-[var(--bg)] text-[var(--text-muted)]' : 'bg-[color-mix(in_oklab,var(--primary)_10%,transparent)] text-[var(--primary)]'}`}>
                <Bell size={18} />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between gap-2 mb-1">
                  <div className="flex items-center min-w-0">
                    {!notification.is_read && (
                      <div className="w-2 h-2 rounded-full bg-[var(--primary)] mr-2 shrink-0" />
                    )}
                    <h3 className={`font-semibold truncate text-sm sm:text-base ${notification.is_read ? 'text-[var(--text)]' : 'text-[var(--primary)]'}`}>
                      {notification.title || 'Notification'}
                    </h3>
                  </div>
                  <span className="text-[10px] sm:text-xs text-[var(--text-muted)] flex items-center gap-1 shrink-0">
                    <Clock size={12} />
                    {notification.created_at ? formatDistanceToNow(new Date(notification.created_at), { addSuffix: true }) : ''}
                  </span>
                </div>
                <p className="text-[var(--text-muted)] text-sm leading-relaxed line-clamp-2">
                  {notification.message || notification.body || ''}
                </p>
              </div>
            </div>
          </div>
        )}
        emptyState={
          <EmptyState>
            <div className="p-4 bg-[var(--bg-elevated)] rounded-full mb-4">
              <Bell size={32} className="text-[var(--text-muted)]" />
            </div>
            <EmptyTitle>All caught up!</EmptyTitle>
            <EmptyDescription>No new notifications at the moment. We'll alert you when something happens.</EmptyDescription>
          </EmptyState>
        }
        itemsClassName="flex flex-col rounded-xl border border-[var(--border)] overflow-hidden bg-[var(--bg-elevated)] shadow-[var(--elev-1)]"
      />
    </div>
  );
}
