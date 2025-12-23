import { useState } from 'react';
import { Calendar, Clock, Users, ChevronRight, Copy, Check } from 'lucide-react';
import { Link } from 'react-router-dom';
import { format } from 'date-fns';

function formatLabel(value?: string) {
  if (!value) return '';
  return value.charAt(0).toUpperCase() + value.slice(1);
}

function toDate(value?: string) {
  if (!value) return null;
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? null : d;
}

function formatDateRange(startDate?: string, endDate?: string) {
  const start = toDate(startDate);
  const end = toDate(endDate);
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
  const start = toDate(startTime);
  const end = toDate(endTime);
  if (start && end) return `${format(start, 'p')} – ${format(end, 'p')}`;
  if (start) return format(start, 'p');
  if (end) return format(end, 'p');
  return 'Time TBD';
}

function getStatusClasses(status?: string) {
  const value = (status || '').toLowerCase();
  if (['active', 'ongoing'].includes(value)) return { accent: 'var(--success)' };
  if (['pending', 'draft'].includes(value)) return { accent: 'var(--warning)' };
  if (['completed'].includes(value)) return { accent: 'var(--primary)' };
  if (['canceled', 'cancelled', 'expired'].includes(value)) return { accent: 'var(--danger)' };
  return { accent: 'var(--text-muted)' };
}

export function AppointmentCard({ item }: { item: any }) {
  const [copiedCode, setCopiedCode] = useState(false);
  const dateLabel = formatDateRange(item.startDate, item.endDate);
  const timeLabel = formatTimeRange(item.startTime, item.endTime);
  const statusClasses = getStatusClasses(item.status);

  const copyToClipboard = (code: string) => {
    navigator.clipboard.writeText(code);
    setCopiedCode(true);
    setTimeout(() => setCopiedCode(false), 2000);
  };

  const statusStyle = {
    background: `color-mix(in oklab, ${statusClasses.accent} 20%, transparent)`,
    color: statusClasses.accent,
    borderColor: `color-mix(in oklab, ${statusClasses.accent} 55%, transparent)`,
  };

  const cardBody = (
    <div className="rounded-xl border border-[var(--border)] bg-[var(--bg-elevated)] text-[var(--text)] shadow-[var(--elev-1)] hover:shadow-[var(--elev-2)] transition-all duration-300 overflow-hidden" style={{ borderLeftWidth: 4, borderLeftColor: statusClasses.accent }}>
      <div className="p-4 sm:p-5 space-y-3">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <h2 className="text-sm sm:text-base font-semibold text-[var(--text)] truncate">
              {item.title || 'Untitled appointment'}
            </h2>
            <div className="text-xs text-[var(--text-muted)] truncate">
              {dateLabel} • {timeLabel}
            </div>
          </div>
          <div className="shrink-0 inline-flex items-center gap-1 text-xs font-semibold text-[var(--primary)]">
            <span className="hidden sm:inline">Manage</span>
            <ChevronRight className="w-4 h-4" />
          </div>
        </div>

        {item.appCode && (
          <div className="flex items-center gap-2 text-xs text-[var(--text-muted)]">
            <button
              onClick={(event) => {
                event.preventDefault();
                event.stopPropagation();
                copyToClipboard(item.appCode);
              }}
              className="inline-flex items-center gap-2 px-2.5 py-1.5 rounded-md text-xs font-semibold border"
              style={{
                background: copiedCode
                  ? 'color-mix(in oklab, var(--success) 15%, transparent)'
                  : 'color-mix(in oklab, var(--primary) 10%, transparent)',
                borderColor: copiedCode
                  ? 'color-mix(in oklab, var(--success) 45%, transparent)'
                  : 'color-mix(in oklab, var(--primary) 35%, transparent)',
                color: copiedCode ? 'var(--success)' : 'var(--primary)',
              }}
            >
              <span className="font-mono tracking-wide truncate">{item.appCode}</span>
              {copiedCode ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
              <span className="hidden sm:inline">{copiedCode ? 'Copied' : 'Copy'}</span>
            </button>
          </div>
        )}

        <div className="flex flex-wrap items-center gap-2">
          <span className="px-2 py-0.5 rounded-full text-xs font-semibold border" style={statusStyle}>
            {formatLabel(item.status)}
          </span>
          <span className="px-2 py-0.5 rounded-full text-xs font-semibold border" style={{ background: 'color-mix(in oklab, var(--primary) 12%, transparent)', color: 'var(--primary)', borderColor: 'color-mix(in oklab, var(--primary) 45%, transparent)' }}>
            {formatLabel(String(item.type))}
          </span>
          {item.maxAttendees ? (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold border whitespace-nowrap" style={{ background: 'color-mix(in oklab, var(--warning) 12%, transparent)', borderColor: 'color-mix(in oklab, var(--warning) 35%, transparent)' }}>
              <Users className="w-3.5 h-3.5" />
              {item.maxAttendees} slots
            </span>
          ) : null}
        </div>
      </div>
    </div>
  );

  if (item?.id) {
    return (
      <Link to={`/appointments/${item.id}`} state={{ appointment: item }} className="block">
        {cardBody}
      </Link>
    );
  }

  return cardBody;
}
