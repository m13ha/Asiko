import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import * as API from '@appointment-master/api-client';
import { useMyAppointments } from '../hooks';
import { Button } from '@/components/Button';
import { AppointmentCard } from '../components/AppointmentCard';
import { EmptyState, EmptyTitle, EmptyDescription, EmptyAction } from '@/components/EmptyState';
import { PaginatedGrid } from '@/components/PaginatedGrid';
import { Pagination } from '@/components/Pagination';
import { usePagination } from '@/hooks/usePagination';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { format } from 'date-fns';
import { StatusFilter } from '@/components/StatusFilter';

const statusOptions = [
  { label: 'Pending', value: API.EntitiesAppointmentStatus.AppointmentStatusPending },
  { label: 'Ongoing', value: API.EntitiesAppointmentStatus.AppointmentStatusOngoing },
  { label: 'Completed', value: API.EntitiesAppointmentStatus.AppointmentStatusCompleted },
  { label: 'Canceled', value: API.EntitiesAppointmentStatus.AppointmentStatusCanceled },
];

const sortOptions = [
  { label: 'Recently updated', value: 'updated_desc' as const },
  { label: 'Recently created', value: 'created_desc' as const },
  { label: 'Start date (soonest)', value: 'start_asc' as const },
  { label: 'Title (A–Z)', value: 'title_asc' as const },
];

type SortValue = (typeof sortOptions)[number]['value'];

function safeDate(value?: string): number {
  if (!value) return 0;
  const d = new Date(value);
  const t = d.getTime();
  return Number.isNaN(t) ? 0 : t;
}

function formatDateRange(startDate?: string, endDate?: string) {
  const start = safeDate(startDate) ? new Date(startDate as string) : null;
  const end = safeDate(endDate) ? new Date(endDate as string) : null;
  if (start && end) {
    const sameDay = start.toDateString() === end.toDateString();
    return sameDay
      ? format(start, 'EEE, MMM d, yyyy')
      : `${format(start, 'MMM d, yyyy')} → ${format(end, 'MMM d, yyyy')}`;
  }
  if (start) return format(start, 'EEE, MMM d, yyyy');
  if (end) return format(end, 'EEE, MMM d, yyyy');
  return 'Date TBD';
}

function formatTimeRange(startTime?: string, endTime?: string) {
  const start = safeDate(startTime) ? new Date(startTime as string) : null;
  const end = safeDate(endTime) ? new Date(endTime as string) : null;
  if (start && end) return `${format(start, 'p')} – ${format(end, 'p')}`;
  if (start) return format(start, 'p');
  if (end) return format(end, 'p');
  return 'Time TBD';
}

export function MyAppointmentsPage() {
  const [selectedStatuses, setSelectedStatuses] = useState<API.EntitiesAppointmentStatus[]>([]);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState<SortValue>('start_asc');
  const pagination = usePagination(1, 10);
  const { data, isLoading, error } = useMyAppointments({ 
    statuses: selectedStatuses,
    ...pagination.params 
  });
  const navigate = useNavigate();

  const filtered = useMemo(() => {
    const items = (data?.items ?? []) as any[];
    const needle = query.trim().toLowerCase();
    const searched = needle
      ? items.filter((item) => {
          const title = String(item?.title ?? '').toLowerCase();
          const code = String(item?.appCode ?? '').toLowerCase();
          return title.includes(needle) || code.includes(needle);
        })
      : items;

    const sorted = searched.slice().sort((a, b) => {
      switch (sort) {
        case 'created_desc':
          return safeDate(b?.createdAt) - safeDate(a?.createdAt);
        case 'start_asc': {
          const diff = safeDate(a?.startDate) - safeDate(b?.startDate);
          return diff !== 0 ? diff : safeDate(b?.updatedAt) - safeDate(a?.updatedAt);
        }
        case 'title_asc':
          return String(a?.title ?? '').localeCompare(String(b?.title ?? ''));
        case 'updated_desc':
        default:
          return safeDate(b?.updatedAt) - safeDate(a?.updatedAt);
      }
    });

    return sorted;
  }, [data?.items, query, sort]);

  const derivedData = useMemo(() => {
    if (!data) return data;
    return { ...data, items: filtered };
  }, [data, filtered]);

  const clearFilters = () => {
    setSelectedStatuses([]);
    setQuery('');
    setSort('start_asc');
    pagination.updatePage(1);
  };

  return (
    <div className="grid gap-6">
      <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] p-4 sm:p-6 min-w-0">
        <div className="flex justify-between items-start flex-wrap gap-3">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-[var(--text)] m-0">My Appointments</h1>
            <div className="mt-1 text-sm text-[var(--text-muted)]">
              Browse, filter, and manage all your published schedules.
            </div>
          </div>
          <Button variant="primary" onClick={() => navigate('/appointments/new')}>
            Create appointment
          </Button>
        </div>

        <div className="mt-5 grid gap-3">
          <div className="grid grid-cols-1 md:grid-cols-[minmax(0,1fr)_240px] gap-3 items-start">
            <div className="min-w-0">
              <div className="text-xs uppercase tracking-wide text-[var(--text-muted)] mb-1">Search</div>
              <Input
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search by title or code…"
              />
            </div>
            <div>
              <div className="text-xs uppercase tracking-wide text-[var(--text-muted)] mb-1">Sort</div>
              <Select
                value={sort}
                options={sortOptions}
                optionLabel="label"
                optionValue="value"
                onChange={(e) => setSort(e.value as SortValue)}
              />
            </div>
          </div>

          <div className="flex items-center justify-between gap-3 flex-wrap">
            <StatusFilter
              options={statusOptions}
              selected={selectedStatuses}
              onChange={(next) => {
                setSelectedStatuses(next);
                pagination.updatePage(1);
              }}
            />

            {(selectedStatuses.length > 0 || query.trim() || sort !== 'updated_desc') && (
              <Button variant="ghost" size="sm" onClick={clearFilters}>
                Clear filters
              </Button>
            )}
          </div>

          <div className="text-sm text-[var(--text-muted)]">
            Showing {filtered.length}{typeof data?.total === 'number' ? ` of ${data.total}` : ''}.
          </div>
        </div>
      </div>
      
      <div className="hidden lg:block">
        {isLoading ? (
          <div className="flex justify-center items-center py-8">Loading...</div>
        ) : error ? (
          <div className="text-red-600 py-4">Failed to load data. Please try again.</div>
        ) : !derivedData?.items?.length ? (
          <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] p-12 text-center">
            <EmptyState>
              <EmptyTitle>No appointments yet</EmptyTitle>
              <EmptyDescription>Create your first appointment and share the code.</EmptyDescription>
              <EmptyAction>
                <Button variant="primary" onClick={() => navigate('/appointments/new')}>
                  Create appointment
                </Button>
              </EmptyAction>
            </EmptyState>
          </div>
        ) : (
          <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] overflow-hidden">
            <table className="w-full text-left text-sm">
              <thead className="bg-[var(--bg-muted)] text-[var(--text-muted)]">
                <tr>
                  <th className="px-4 py-3 font-semibold">Title</th>
                  <th className="px-4 py-3 font-semibold">Code</th>
                  <th className="px-4 py-3 font-semibold">Status</th>
                  <th className="px-4 py-3 font-semibold">Date</th>
                  <th className="px-4 py-3 font-semibold">Time</th>
                  <th className="px-4 py-3 font-semibold">Type</th>
                  <th className="px-4 py-3 font-semibold">Action</th>
                </tr>
              </thead>
              <tbody>
                {derivedData.items.map((item: any) => (
                  <tr
                    key={item.id}
                    className="border-t border-[var(--border)] hover:bg-[color-mix(in_oklab,var(--primary)_6%,transparent)] cursor-pointer"
                    onClick={() => navigate(`/appointments/${item.id}`, { state: { appointment: item } })}
                  >
                    <td className="px-4 py-3 font-semibold text-[var(--text)] whitespace-normal break-words">{item.title || 'Untitled appointment'}</td>
                    <td className="px-4 py-3 font-mono text-xs text-[var(--text-muted)]">{item.appCode || '—'}</td>
                    <td className="px-4 py-3 text-xs uppercase tracking-wide">{item.status || '—'}</td>
                    <td className="px-4 py-3">{formatDateRange(item.startDate, item.endDate)}</td>
                    <td className="px-4 py-3">{formatTimeRange(item.startTime, item.endTime)}</td>
                    <td className="px-4 py-3">{String(item.type || '—')}</td>
                    <td className="px-4 py-3 text-[var(--primary)] font-semibold">Manage</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {(derivedData?.totalPages ?? 0) > 1 && (
          <Pagination
            currentPage={pagination.page}
            totalPages={derivedData?.totalPages ?? 1}
            onPageChange={pagination.updatePage}
          />
        )}
      </div>

      <div className="lg:hidden">
        <PaginatedGrid
          data={derivedData}
          isLoading={isLoading}
          error={error}
          onPageChange={pagination.updatePage}
          layout="list"
          renderItem={(item: any) => <AppointmentCard key={item.id} item={item} />}
          emptyState={
            <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] p-12 text-center">
              <EmptyState>
                <EmptyTitle>No appointments yet</EmptyTitle>
                <EmptyDescription>Create your first appointment and share the code.</EmptyDescription>
                <EmptyAction>
                  <Button variant="primary" onClick={() => navigate('/appointments/new')}>
                    Create appointment
                  </Button>
                </EmptyAction>
              </EmptyState>
            </div>
          }
        />
      </div>
    </div>
  );
}
