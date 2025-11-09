import { Card } from '@/components/Card';

export function BookingSummary({ appCode, date, startTime, endTime, attendeeCount }: { appCode: string; date: string; startTime: string; endTime: string; attendeeCount?: number; }) {
  return (
    <Card>
      <div style={{ display: 'grid', gap: 6 }}>
        <div><small>Code:</small> <strong>{appCode}</strong></div>
        <div><small>Date:</small> <strong>{date}</strong></div>
        <div><small>Time:</small> <strong>{startTime} - {endTime}</strong></div>
        {attendeeCount && <div><small>Attendees:</small> <strong>{attendeeCount}</strong></div>}
      </div>
    </Card>
  );
}

