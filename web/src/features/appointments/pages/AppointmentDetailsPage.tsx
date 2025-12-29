import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { useState } from 'react';
import { format } from 'date-fns';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { useAppointmentByAppCode, useAppointmentUsers, useDeleteAppointment, useMyAppointments } from '../hooks';
import { useConfirmBooking, useRejectBooking } from '@/features/bookings/hooks';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';
import { ListItem } from '@/components/ListItem';
import { CopyButton } from '@/components/CopyButton';
import { PaginatedGrid } from '@/components/PaginatedGrid';
import { Pagination } from '@/components/Pagination';
import { usePagination } from '@/hooks/usePagination';
import { Dialog } from 'primereact/dialog';

export function AppointmentDetailsPage() {
  const { id = '' } = useParams();
  const loc = useLocation() as any;
  const navigate = useNavigate();
  const fromState = loc.state?.appointment;
  const { data: my } = useMyAppointments();
  const match = my?.items?.find((i: any) => i.id === id);
  const appt = fromState || match;
  
  const pagination = usePagination(1, 10);
  const { data: users, isLoading, error } = useAppointmentUsers(appt?.appCode, pagination.params, { 
    enabled: !!appt?.appCode 
  });
  const deleteAppointment = useDeleteAppointment();
  const reject = useRejectBooking(appt?.appCode || id);
  const confirm = useConfirmBooking(appt?.appCode || id);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const bookingLink = typeof window !== 'undefined' && appt?.appCode ? `${window.location.origin}/book-by-code?code=${appt.appCode}` : '';
  const hasBookings = Boolean((users?.total ?? users?.items?.length ?? 0) > 0);
  const visibleUsers = (users?.items ?? []).filter((u: any) => !['active', 'cancelled', 'canceled', 'rejected'].includes(String(u?.status || '').toLowerCase()));
  const sortedUsers = visibleUsers.slice().sort((a: any, b: any) => {
    const aDate = parseDate(a?.date)?.getTime() ?? 0;
    const bDate = parseDate(b?.date)?.getTime() ?? 0;
    if (aDate !== bDate) return aDate - bDate;
    const aStart = parseDate(a?.startTime)?.getTime() ?? 0;
    const bStart = parseDate(b?.startTime)?.getTime() ?? 0;
    return aStart - bStart;
  });

  const schedule = formatSchedule(appt);
  const details = [
    { label: 'Type', value: appt?.type ? capitalize(appt.type) : null },
    {
      label: 'Capacity',
      value: appt?.type === 'party'
        ? appt?.maxAttendees ? `${appt.maxAttendees} total` : 'Unlimited'
        : appt?.maxAttendees ? `${appt.maxAttendees} per slot` : '1 per slot',
    },
    { label: 'Slot length', value: appt?.bookingDuration ? `${appt.bookingDuration} mins` : null },
    {
      label: 'Status',
      value: appt?.status
        ? String(appt.status).toLowerCase() === 'expired'
          ? 'Completed'
          : capitalize(appt.status)
        : 'Active',
    },
  ].filter((item) => item.value);

  if (!appt) {
    return (
      <Card className="p-6">
        <CardHeader>
          <CardTitle>Appointment not found</CardTitle>
        </CardHeader>
        <div className="text-sm text-[var(--text-muted)]">We couldn&apos;t locate this appointment.</div>
        <div className="mt-4">
          <Button variant="outline" onClick={() => navigate('/appointments')}>Back to appointments</Button>
        </div>
      </Card>
    );
  }

  const handleDelete = () => {
    if (!appt?.id) return;
    setConfirmOpen(true);
  };

  const handleConfirmDelete = () => {
    if (!appt?.id) return;
    deleteAppointment.mutate(appt.id, {
      onSuccess: () => {
        setConfirmOpen(false);
        navigate('/appointments');
      },
    });
  };

  return (
    <div className="grid gap-6">
      <div>
        <Button variant="ghost" onClick={() => navigate(-1)}>
          Back
        </Button>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Appointment Details</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 16 }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>Title</span>
            <strong style={{ fontSize: 20, wordBreak: 'break-word', overflowWrap: 'anywhere' }}>{appt?.title || 'Untitled appointment'}</strong>
          </div>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
            {!hasBookings && (
              <Button
                variant="primary"
                onClick={() => navigate(`/appointments/${appt.id}/edit`, { state: { appointment: appt } })}
              >
                Edit appointment
              </Button>
            )}
            <Button variant="outline" onClick={handleDelete} disabled={deleteAppointment.isPending}>
              Delete appointment
            </Button>
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
      <Dialog
        header="Delete appointment?"
        visible={confirmOpen}
        onHide={() => setConfirmOpen(false)}
        className="w-full max-w-md"
        contentClassName="!p-0"
      >
        <div className="grid gap-4 p-4">
          <div className="text-sm text-[var(--text-muted)]">
            This will delete the appointment and cancel any active bookings. This action cannot be undone.
          </div>
          <div className="flex flex-wrap justify-end gap-2">
            <Button variant="ghost" onClick={() => setConfirmOpen(false)} disabled={deleteAppointment.isPending}>
              Keep appointment
            </Button>
            <Button variant="primary" onClick={handleConfirmDelete} disabled={deleteAppointment.isPending}>
              {deleteAppointment.isPending ? 'Deleting…' : 'Delete'}
            </Button>
          </div>
        </div>
      </Dialog>

      <Card>
        <CardHeader>
          <CardTitle>Registered Users/Bookings</CardTitle>
        </CardHeader>
        
        {/* Desktop Table View */}
        <div className="hidden lg:block">
          {isLoading ? (
            <div className="flex justify-center items-center py-8">Loading...</div>
          ) : error ? (
            <div className="text-red-600 py-4">Failed to load users.</div>
          ) : !users?.items?.length ? (
            <EmptyState>
              <EmptyTitle>No bookings yet</EmptyTitle>
              <EmptyDescription>Share the appointment code to start receiving bookings.</EmptyDescription>
            </EmptyState>
          ) : (
            <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-1)] overflow-hidden">
              <table className="w-full text-left text-sm">
                <thead className="bg-[var(--bg-muted)] text-[var(--text-muted)]">
                  <tr>
                    <th className="px-4 py-3 font-semibold">Name/Contact</th>
                    <th className="px-4 py-3 font-semibold">Code</th>
                    <th className="px-4 py-3 font-semibold">Status</th>
                    <th className="px-4 py-3 font-semibold">Date</th>
                    <th className="px-4 py-3 font-semibold">Time</th>
                    <th className="px-4 py-3 font-semibold">Seats</th>
                    <th className="px-4 py-3 font-semibold">Action</th>
                  </tr>
                </thead>
                <tbody>
                  {sortedUsers.map((u: any) => {
                    const status = String(u?.status || '').toLowerCase();
                    const canConfirm = status === 'pending';
                    const canReject = ['active', 'pending', 'confirmed', 'ongoing'].includes(status);
                    return (
                    <tr key={u.id} className="border-t border-[var(--border)] hover:bg-[color-mix(in_oklab,var(--primary)_6%,transparent)]">
                      <td className="px-4 py-3 font-medium">{u.name || u.email || u.phone}</td>
                      <td className="px-4 py-3 font-mono text-xs text-[var(--text-muted)]">{u.bookingCode}</td>
                      <td className="px-4 py-3 text-xs uppercase tracking-wide">{u.status}</td>
                      <td className="px-4 py-3">{formatDate(u.date)}</td>
                      <td className="px-4 py-3">{formatTime(u.startTime)} – {formatTime(u.endTime)}</td>
                      <td className="px-4 py-3">{u.seatsBooked ?? 1}</td>
                      <td className="px-4 py-3">
                        <div className="flex flex-wrap gap-2">
                          {canConfirm && (
                            <Button variant="primary" size="sm" onClick={() => confirm.mutate(u.bookingCode)}>
                              Confirm
                            </Button>
                          )}
                          {canReject && (
                            <Button variant="ghost" size="sm" onClick={() => reject.mutate(u.bookingCode)}>
                              Reject
                            </Button>
                          )}
                        </div>
                      </td>
                    </tr>
                  )})}
                </tbody>
              </table>
            </div>
          )}
          {(users?.totalPages ?? 0) > 1 && (
            <div className="mt-4">
              <Pagination
                currentPage={pagination.page}
                totalPages={users?.totalPages ?? 1}
                onPageChange={pagination.updatePage}
              />
            </div>
          )}
        </div>

        {/* Mobile Grid View */}
        <div className="lg:hidden">
          <PaginatedGrid
            data={users ? { ...users, items: visibleUsers } : users}
            isLoading={isLoading}
            error={error}
            onPageChange={pagination.updatePage}
            layout="list"
            renderItem={(u: any) => {
              const status = String(u?.status || '').toLowerCase();
              const canConfirm = status === 'pending';
              const canReject = ['active', 'pending', 'confirmed', 'ongoing'].includes(status);
              return (
              <ListItem key={u.id}>
                <div>
                  <div style={{ fontWeight: 600 }}>{u.name || u.email || u.phone}</div>
                  <div style={{ fontSize: 12, opacity: 0.8 }}>
                    {formatDate(u.date)} • {formatTime(u.startTime)} – {formatTime(u.endTime)}
                  </div>
                  <div style={{ fontSize: 12, opacity: 0.8 }}>Code: {u.bookingCode} • {u.status}</div>
                </div>
                <div className="flex gap-2">
                  {canConfirm && (
                    <Button variant="primary" size="sm" onClick={() => confirm.mutate(u.bookingCode)}>
                      Confirm
                    </Button>
                  )}
                  {canReject && (
                    <Button variant="ghost" size="sm" onClick={() => reject.mutate(u.bookingCode)}>
                      Reject
                    </Button>
                  )}
                </div>
              </ListItem>
            )}}
            emptyState={
              <EmptyState>
                <EmptyTitle>No bookings yet</EmptyTitle>
                <EmptyDescription>Share the appointment code to start receiving bookings.</EmptyDescription>
              </EmptyState>
            }
          />
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

function formatDate(value?: string) {
  const d = parseDate(value);
  if (!d) return '—';
  return format(d, 'EEE, MMM d, yyyy');
}

function formatTime(value?: string) {
  const d = parseDate(value);
  if (!d) return '—';
  return format(d, 'p');
}

function capitalize(value?: string) {
  if (!value) return '';
  return value.charAt(0).toUpperCase() + value.slice(1);
}
