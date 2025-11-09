import { useParams } from 'react-router-dom';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useBookingByCode, useCancelBooking, useUpdateBooking } from '../hooks';
import { Input } from '@/components/Input';
import { useState } from 'react';
import { Button } from '@/components/Button';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Calendar, Clock } from 'lucide-react';
import { CopyButton } from '@/components/CopyButton';

export function BookingManagePage() {
  const { bookingCode = '' } = useParams();
  const { data, isLoading, error, refetch } = useBookingByCode(bookingCode);
  const [date, setDate] = useState('');
  const [startTime, setStart] = useState('');
  const [endTime, setEnd] = useState('');
  const update = useUpdateBooking(bookingCode);
  const cancel = useCancelBooking(bookingCode);

  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <Card>
        <CardHeader>
          <CardTitle>Manage Booking</CardTitle>
        </CardHeader>
        {isLoading && <div>Loading...</div>}
        {error && <div style={{ color: 'var(--danger)' }}>Failed to load booking.</div>}
        {data && (
          <div style={{ display: 'grid', gap: 8 }}>
            <div>
              <small>Booking Code:</small> <strong>{data.bookingCode}</strong> <CopyButton value={data?.bookingCode || ''} ariaLabel="Copy booking code" />
            </div>
            <div><small>Date:</small> <strong>{data.date}</strong></div>
            <div><small>Time:</small> <strong>{data.startTime} - {data.endTime}</strong></div>
          </div>
        )}
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Reschedule</CardTitle>
        </CardHeader>
        <div style={{ display: 'grid', gap: 8 }}>
          <Field>
            <FieldLabel>New date</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><Calendar size={16} /></IconSlot>
                <Input type="date" value={date} onChange={(e) => setDate(e.target.value)} style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
          </Field>
          <div
            style={{
              display: 'grid',
              gap: 8,
              gridTemplateColumns: 'repeat(auto-fit, minmax(160px, 1fr))',
            }}
          >
            <Field>
              <FieldLabel>Start time</FieldLabel>
              <FieldRow>
                <div style={{ position: 'relative' }}>
                  <IconSlot><Clock size={16} /></IconSlot>
                  <Input type="time" value={startTime} onChange={(e) => setStart(e.target.value)} style={{ paddingLeft: 36 }} />
                </div>
              </FieldRow>
            </Field>
            <Field>
              <FieldLabel>End time</FieldLabel>
              <FieldRow>
                <div style={{ position: 'relative' }}>
                  <IconSlot><Clock size={16} /></IconSlot>
                  <Input type="time" value={endTime} onChange={(e) => setEnd(e.target.value)} style={{ paddingLeft: 36 }} />
                </div>
              </FieldRow>
            </Field>
          </div>
          <div>
            <Button
              variant="primary"
              onClick={() => update.mutate({ appCode: data?.appCode || '', date, startTime, endTime }, { onSuccess: () => refetch() })}
              disabled={!date || !startTime || !endTime || update.isPending}
            >
              {update.isPending ? 'Updating…' : 'Update booking'}
            </Button>
          </div>
          <div>
            <Button onClick={() => cancel.mutate(undefined, { onSuccess: () => refetch() })} disabled={cancel.isPending}>
              {cancel.isPending ? 'Cancelling…' : 'Cancel booking'}
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}
