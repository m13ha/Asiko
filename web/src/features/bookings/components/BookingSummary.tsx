import { format, parseISO } from 'date-fns';
import { Card } from '@/components/Card';

export function BookingSummary({ appCode, date, startTime, endTime, attendeeCount }: { appCode: string; date: string; startTime: string; endTime: string; attendeeCount?: number; }) {
  const formattedDate = date ? format(parseISO(date), 'EEEE, MMMM d, yyyy') : date;
  
  const formatTime = (iso: string) => {
    try {
      return format(parseISO(iso), 'h:mm a');
    } catch {
      return iso;
    }
  };

  const timeRange = startTime && endTime 
    ? `${formatTime(startTime)} - ${formatTime(endTime)}`
    : '';

  return (
    <Card className="!p-4">
      <div className="space-y-3">
        <div className="flex justify-between items-center border-b border-[var(--border)] pb-2">
            <span className="text-[var(--text-muted)] text-sm">Appointment Code</span>
            <span className="font-mono font-medium text-[var(--primary)]">{appCode}</span>
        </div>
        <div className="flex justify-between items-center">
            <span className="text-[var(--text-muted)] text-sm">Date</span>
            <span className="font-medium text-right">{formattedDate}</span>
        </div>
        <div className="flex justify-between items-center">
            <span className="text-[var(--text-muted)] text-sm">Time</span>
            <span className="font-medium text-right">{timeRange}</span>
        </div>
        {attendeeCount && (
            <div className="flex justify-between items-center pt-2 border-t border-[var(--border)]">
                <span className="text-[var(--text-muted)] text-sm">Attendees</span>
                <span className="font-medium">{attendeeCount}</span>
            </div>
        )}
      </div>
    </Card>
  );
}

