import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { format } from 'date-fns';
import { useMyBookings } from '../hooks';
import { BookingCard } from '@/components/BookingCard';
import { PaginatedGrid } from '@/components/PaginatedGrid';
import { Pagination } from '@/components/Pagination';
import { usePagination } from '@/hooks/usePagination';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Button } from '@/components/Button';
import { StatusFilter } from '@/components/StatusFilter';

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Ongoing', value: 'ongoing' },
  { label: 'Expired', value: 'expired' },
  { label: 'Pending', value: 'pending' },
  { label: 'Confirmed', value: 'confirmed' },
  { label: 'Cancelled', value: 'cancelled' },
  { label: 'Canceled', value: 'canceled' },
  { label: 'Rejected', value: 'rejected' },
];

const sortOptions = [
  { label: 'Recently updated', value: 'updated_desc' as const },
  { label: 'Recently created', value: 'created_desc' as const },
  { label: 'Booking date (soonest)', value: 'date_asc' as const },
  { label: 'Name (A–Z)', value: 'name_asc' as const },
];

type SortValue = (typeof sortOptions)[number]['value'];

function safeDate(value?: string | Date): number {
  if (!value) return 0;
  const d = new Date(value);
  const t = d.getTime();
  return Number.isNaN(t) ? 0 : t;
}

function formatDate(value?: string | Date) {
  if (!value) return 'Date TBD';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);
  return format(d, 'EEE, MMM d, yyyy');
}

function formatTime(value?: string | Date) {
  if (!value) return 'Time TBD';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);
  return format(d, 'p');
}

export function MyBookingsPage() {
  const [selectedStatuses, setSelectedStatuses] = useState<string[]>([]);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState<SortValue>('updated_desc');
  const pagination = usePagination(1, 10);
  const navigate = useNavigate();
  
  const { data, isLoading, error } = useMyBookings({ 
    statuses: selectedStatuses,
    ...pagination.params 
  });

  const filtered = useMemo(() => {
    const items = (data?.items ?? []) as any[];
    const needle = query.trim().toLowerCase();
    const searched = needle
      ? items.filter((item) => {
          const name = String(item?.name ?? '').toLowerCase();
          const email = String(item?.email ?? '').toLowerCase();
          const code = String(item?.bookingCode ?? '').toLowerCase();
          const appCode = String(item?.appCode ?? '').toLowerCase();
          return name.includes(needle) || email.includes(needle) || code.includes(needle) || appCode.includes(needle);
        })
      : items;

    const sorted = searched.slice().sort((a, b) => {
      switch (sort) {
        case 'created_desc':
          return safeDate(b?.createdAt) - safeDate(a?.createdAt);
        case 'date_asc': {
          const diff = safeDate(a?.date) - safeDate(b?.date);
          return diff !== 0 ? diff : safeDate(b?.updatedAt) - safeDate(a?.updatedAt);
        }
        case 'name_asc':
          return String(a?.name ?? '').localeCompare(String(b?.name ?? ''));
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
    setSort('updated_desc');
    pagination.updatePage(1);
  };
  
  return (
    <div className="grid gap-6">
      <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] p-4 sm:p-6 min-w-0">
        <div className="flex justify-between items-start flex-wrap gap-3">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-[var(--text)] m-0">My Bookings</h1>
            <div className="mt-1 text-sm text-[var(--text-muted)]">
              Manage your appointment bookings and schedules.
            </div>
          </div>
        </div>

        <div className="mt-5 grid gap-3">
          <div className="grid grid-cols-1 md:grid-cols-[minmax(0,1fr)_240px] gap-3 items-start">
            <div className="min-w-0">
              <div className="text-xs uppercase tracking-wide text-[var(--text-muted)] mb-1">Search</div>
              <Input
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search by code, name, or email…"
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
              <EmptyTitle>No bookings yet</EmptyTitle>
              <EmptyDescription>When you book a slot, it will appear here.</EmptyDescription>
            </EmptyState>
          </div>
        ) : (
          <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] overflow-hidden">
            <table className="w-full text-left text-sm">
              <thead className="bg-[var(--bg-muted)] text-[var(--text-muted)]">
                <tr>
                  <th className="px-4 py-3 font-semibold">Booking</th>
                  <th className="px-4 py-3 font-semibold">Appointment</th>
                  <th className="px-4 py-3 font-semibold">Status</th>
                  <th className="px-4 py-3 font-semibold">Date</th>
                  <th className="px-4 py-3 font-semibold">Time</th>
                  <th className="px-4 py-3 font-semibold">Seats</th>
                  <th className="px-4 py-3 font-semibold">Action</th>
                </tr>
              </thead>
              <tbody>
                {derivedData.items.map((booking: any) => (
                  <tr
                    key={booking.id}
                    className="border-t border-[var(--border)] hover:bg-[color-mix(in_oklab,var(--primary)_6%,transparent)] cursor-pointer"
                    onClick={() => booking.bookingCode && navigate(`/bookings/${booking.bookingCode}`)}
                  >
                    <td className="px-4 py-3 font-mono text-xs text-[var(--text-muted)]">{booking.bookingCode || '—'}</td>
                    <td className="px-4 py-3 font-semibold text-[var(--text)]">{booking.appCode || '—'}</td>
                    <td className="px-4 py-3 text-xs uppercase tracking-wide">{booking.status || '—'}</td>
                    <td className="px-4 py-3">{formatDate(booking.date)}</td>
                    <td className="px-4 py-3">{formatTime(booking.startTime)} – {formatTime(booking.endTime)}</td>
                    <td className="px-4 py-3">{booking.seatsBooked ?? '—'}</td>
                    <td className="px-4 py-3 text-[var(--primary)] font-semibold">Manage</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {derivedData?.totalPages && derivedData.totalPages > 1 && (
          <Pagination
            currentPage={pagination.page}
            totalPages={derivedData.totalPages ?? 1}
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
          renderItem={(booking: any) => (
            <BookingCard 
              key={booking.id} 
              booking={booking}
              showActions={true}
            />
          )}
          emptyState={
            <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] p-12 text-center">
              <EmptyState>
                <EmptyTitle>No bookings yet</EmptyTitle>
                <EmptyDescription>When you book a slot, it will appear here.</EmptyDescription>
              </EmptyState>
            </div>
          }
        />
      </div>
    </div>
  );
}
