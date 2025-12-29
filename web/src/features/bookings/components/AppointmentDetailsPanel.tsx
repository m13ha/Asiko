import { memo, useMemo } from 'react';
import { Badge } from '@/components/Badge';
import { Clock, Users, FileText, User } from 'lucide-react';

export interface AppointmentInfo {
  title?: string;
  description?: string;
  hostName?: string;
  duration?: number;
  maxAttendees?: number;
  type?: string;
  startDate?: string;
  endDate?: string;
}

export interface AppointmentDetailsPanelProps {
  appointment: AppointmentInfo;
  className?: string;
}

/** Memoized detail row to prevent unnecessary re-renders */
const DetailRow = memo(function DetailRow({
  icon: Icon,
  label,
  value,
  accentColor = 'var(--primary)',
}: {
  icon: React.ComponentType<{ className?: string; style?: React.CSSProperties }>;
  label: string;
  value: string | number;
  accentColor?: string;
}) {
  return (
    <div className="flex items-center gap-3">
      <div
        className="p-2 rounded-lg"
        style={{ background: `color-mix(in oklab, ${accentColor} 15%, transparent)` }}
      >
        <Icon className="w-4 h-4" style={{ color: accentColor }} />
      </div>
      <div className="min-w-0 flex-1">
        <p className="text-xs text-[var(--text-muted)] font-medium m-0">{label}</p>
        <p className="text-sm font-semibold text-[var(--text)] m-0 truncate">{value}</p>
      </div>
    </div>
  );
});

/** Format the appointment type for display */
function formatAppointmentType(type?: string): string {
  switch (type?.toLowerCase()) {
    case 'single':
      return '1:1 Meeting';
    case 'group':
      return 'Group Session';
    case 'party':
      return 'Shared Event';
    default:
      return 'Appointment';
  }
}

/** Format date range for display */
function formatDateRange(startDate?: string, endDate?: string): string {
  if (!startDate) return 'Flexible dates';
  const start = new Date(startDate);
  const formatter = new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric' });
  
  if (!endDate || startDate === endDate) {
    return formatter.format(start);
  }
  
  const end = new Date(endDate);
  return `${formatter.format(start)} – ${formatter.format(end)}`;
}

/**
 * Left panel component displaying appointment details.
 * Designed for the split panel booking layout.
 */
export const AppointmentDetailsPanel = memo(function AppointmentDetailsPanel({
  appointment,
  className = '',
}: AppointmentDetailsPanelProps) {
  const formattedType = useMemo(() => formatAppointmentType(appointment.type), [appointment.type]);
  const formattedDates = useMemo(
    () => formatDateRange(appointment.startDate, appointment.endDate),
    [appointment.startDate, appointment.endDate]
  );

  return (
    <div
      className={`flex flex-col gap-5 p-5 rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] ${className}`}
    >
      {/* Header with title and type badge */}
      <div className="space-y-2">
        <Badge tone="info" className="mb-1">
          {formattedType}
        </Badge>
        <h2 className="text-xl font-bold text-[var(--text)] m-0 leading-tight break-words whitespace-normal">
          {appointment.title || 'Appointment'}
        </h2>
        {appointment.hostName && (
          <p className="text-sm text-[var(--text-muted)] m-0 flex items-center gap-1.5">
            <User className="w-3.5 h-3.5" />
            Hosted by {appointment.hostName}
          </p>
        )}
      </div>

      {/* Divider */}
      <div className="h-px bg-[var(--border)]" />

      {/* Details grid */}
      <div className="grid gap-4">
        {appointment.duration && (
          <DetailRow
            icon={Clock}
            label="Duration"
            value={`${appointment.duration} minutes`}
            accentColor="var(--primary)"
          />
        )}

        {appointment.maxAttendees && appointment.maxAttendees > 1 && (
          <DetailRow
            icon={Users}
            label="Capacity"
            value={`${appointment.maxAttendees} attendees`}
            accentColor="var(--success)"
          />
        )}

        {formattedDates !== 'Flexible dates' && (
          <DetailRow
            icon={FileText}
            label="Available"
            value={formattedDates}
            accentColor="var(--warning)"
          />
        )}
      </div>

      {/* Description if present */}
      {appointment.description && (
        <>
          <div className="h-px bg-[var(--border)]" />
          <div>
            <p className="text-xs text-[var(--text-muted)] font-medium mb-1.5">About</p>
            <p className="text-sm text-[var(--text)] m-0 leading-relaxed">
              {appointment.description}
            </p>
          </div>
        </>
      )}
    </div>
  );
});

export default AppointmentDetailsPanel;
