import { Card } from '@/components/Card';
import { Badge } from '@/components/Badge';
import { Link } from 'react-router-dom';
import { CopyButton } from '@/components/CopyButton';
import { format } from 'date-fns';
import './appointmentCard.css';

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

function statusTone(status?: string) {
  const value = (status || '').toLowerCase();
  if (['active', 'open'].includes(value)) return 'success';
  if (['pending', 'draft'].includes(value)) return 'warning';
  if (['canceled', 'cancelled', 'expired'].includes(value)) return 'danger';
  return 'muted';
}

function typeTone(type?: string) {
  const value = (type || '').toLowerCase();
  if (value.includes('group')) return 'primary';
  if (value.includes('party')) return 'info';
  return 'secondary';
}

export function AppointmentCard({ item }: { item: any }) {
  const dateLabel = formatDateRange(item.startDate, item.endDate);
  const timeLabel = formatTimeRange(item.startTime, item.endTime);

  return (
    <Card className="appointment-card">
      <div className="appointment-card__top">
        <div className="appointment-card__code">
          <div className="code-badge">
            <span className="code-badge__label">Code</span>
            <span className="code-badge__value">{item.appCode || '—'}</span>
          </div>
          {item.appCode && <CopyButton value={item.appCode} ariaLabel="Copy appointment code" />}
        </div>
        <div className="appointment-card__chips">
          {item.status && <Badge tone={statusTone(item.status)}>{formatLabel(item.status)}</Badge>}
          {item.type && <Badge tone={typeTone(String(item.type))}>{formatLabel(String(item.type))}</Badge>}
        </div>
      </div>

      <div className="appointment-card__body">
        <div className="appointment-card__title">{item.title || 'Untitled appointment'}</div>
        {item.description && <div className="appointment-card__desc">{item.description}</div>}

        <div className="appointment-card__meta">
          <div className="meta-row">
            <span className="meta-icon" aria-hidden="true">📅</span>
            <div>
              <div className="meta-label">Dates</div>
              <div className="meta-value">{dateLabel}</div>
            </div>
          </div>
          <div className="meta-row">
            <span className="meta-icon" aria-hidden="true">⏰</span>
            <div>
              <div className="meta-label">Time</div>
              <div className="meta-value">{timeLabel}</div>
            </div>
          </div>
          {item.maxAttendees && (
            <div className="meta-row">
              <span className="meta-icon" aria-hidden="true">👥</span>
              <div>
                <div className="meta-label">Capacity</div>
                <div className="meta-value">{item.maxAttendees} slots</div>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="appointment-card__footer">
        <Link to={`/appointments/${item.id}`} state={{ appointment: item }} className="link-button">
          View details
        </Link>
      </div>
    </Card>
  );
}
