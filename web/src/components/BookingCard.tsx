import { useState } from 'react';
import { Calendar, Clock, ChevronRight, Copy, Check } from 'lucide-react';
import { format } from 'date-fns';
import { Link } from 'react-router-dom';
import { Card } from './Card';
import { Badge } from './Badge';

interface BookingCardProps {
  booking: {
    id?: string;
    bookingCode?: string;
    appCode?: string;
    date?: string;
    startTime?: string;
    endTime?: string;
    status?: string;
    email?: string;
    name?: string;
    seatsBooked?: number;
  };
  showActions?: boolean;
}

function formatTime(timeStr?: string) {
  if (!timeStr) return 'Time TBD';
  try {
    const date = new Date(timeStr);
    return format(date, 'p');
  } catch {
    return timeStr;
  }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return 'Date TBD';
  try {
    const date = new Date(dateStr);
    return format(date, 'EEE, MMM d, yyyy');
  } catch {
    return dateStr;
  }
}

function getStatusTone(status?: string): 'default' | 'success' | 'warning' | 'danger' | 'info' | 'muted' {
  const s = (status || '').toLowerCase();
  if (['active', 'confirmed'].includes(s)) return 'success';
  if (['pending', 'ongoing'].includes(s)) return 'warning';
  if (['cancelled', 'canceled', 'rejected'].includes(s)) return 'danger';
  if (['expired'].includes(s)) return 'muted';
  return 'default';
}

export function BookingCard({ booking, showActions = false }: BookingCardProps) {
  const [copiedCode, setCopiedCode] = useState<string | null>(null);
  const tone = getStatusTone(booking.status);

  const copyToClipboard = (code: string) => {
    navigator.clipboard.writeText(code);
    setCopiedCode(code);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  const cardBody = (
    <Card className="hover:shadow-md transition-shadow duration-300 overflow-hidden border-l-4" style={{ borderLeftColor: `var(--${tone === 'default' ? 'border' : tone})` }}>
      <div className="flex items-start justify-between gap-2 mb-3">
        <div className="min-w-0">
          <h2 className="text-sm sm:text-base font-semibold text-[var(--text)] truncate">
            Booking for {booking.appCode || '—'}
          </h2>
          <div className="text-xs text-[var(--text-muted)] truncate">
            {formatDate(booking.date)} • {formatTime(booking.startTime)} – {formatTime(booking.endTime)}
          </div>
        </div>
        {showActions && (
          <div className="shrink-0 inline-flex items-center gap-1 text-xs font-semibold text-[var(--primary)]">
            <span className="hidden sm:inline">Details</span>
            <ChevronRight className="w-4 h-4" />
          </div>
        )}
      </div>

      {booking.bookingCode && (
        <div className="flex items-center gap-2 text-xs text-[var(--text-muted)] mb-3">
          <button
            onClick={(event) => {
              event.preventDefault();
              event.stopPropagation();
              copyToClipboard(booking.bookingCode!);
            }}
            className="inline-flex items-center gap-2 px-2.5 py-1.5 rounded-md text-xs font-semibold border transition-colors"
            style={{
              background: copiedCode === booking.bookingCode
                ? 'var(--success-faded, rgba(34, 197, 94, 0.1))'
                : 'var(--primary-faded, rgba(59, 130, 246, 0.1))',
              borderColor: copiedCode === booking.bookingCode
                ? 'var(--success)'
                : 'var(--primary)',
              color: copiedCode === booking.bookingCode ? 'var(--success)' : 'var(--primary)',
            }}
          >
            <span className="font-mono tracking-wide truncate">{booking.bookingCode}</span>
            {copiedCode === booking.bookingCode ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
            <span className="hidden sm:inline">{copiedCode === booking.bookingCode ? 'Copied' : 'Copy'}</span>
          </button>
        </div>
      )}

      <div className="flex flex-wrap items-center gap-2 mt-auto">
        <Badge tone={tone}>
          {booking.status || '—'}
        </Badge>
        {booking.seatsBooked ? (
          <Badge tone="info">
            {booking.seatsBooked} seat{booking.seatsBooked !== 1 ? 's' : ''}
          </Badge>
        ) : null}
      </div>
    </Card>
  );

  if (booking.bookingCode) {
    return (
      <Link to={`/bookings/${booking.bookingCode}`} className="block no-underline">
        {cardBody}
      </Link>
    );
  }

  return cardBody;
}
