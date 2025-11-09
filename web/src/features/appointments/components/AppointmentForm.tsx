import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/Input';
import { Textarea } from '@/components/Textarea';
import { Select } from '@/components/Select';
import { Button } from '@/components/Button';
import { Spinner } from '@/components/Spinner';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Type as TypeIcon, FileText, Timer, Calendar, Clock, Users } from 'lucide-react';
import * as API from '@appointment-master/api-client';

const schema = z.object({
  title: z.string().min(1),
  description: z.string().optional(),
  type: z.enum([API.EntitiesAppointmentType.Single, API.EntitiesAppointmentType.Group, API.EntitiesAppointmentType.Party]),
  bookingDuration: z.coerce.number().min(5),
  startDate: z.string().min(1),
  endDate: z.string().min(1),
  startTime: z.string().min(1),
  endTime: z.string().min(1),
  maxAttendees: z.coerce.number().optional(),
});

export type AppointmentFormValues = z.infer<typeof schema>;

export function AppointmentForm({ onSubmit, pending }: { onSubmit: (v: AppointmentFormValues) => void; pending?: boolean }) {
  const { register, handleSubmit, formState: { errors } } = useForm<AppointmentFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      type: API.EntitiesAppointmentType.Single,
      bookingDuration: 30,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="form-grid">
      <Field>
        <FieldLabel>Title</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><TypeIcon size={16} /></IconSlot>
            <Input {...register('title')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
        {errors.title && <small style={{ color: 'var(--danger)' }}>{errors.title.message}</small>}
      </Field>
      <Field>
        <FieldLabel>Description</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><FileText size={16} /></IconSlot>
            <Textarea {...register('description')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
      </Field>
      <Field>
        <FieldLabel>Type</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><TypeIcon size={16} /></IconSlot>
            <Select {...register('type')} style={{ paddingLeft: 36 }}>
              <option value={API.EntitiesAppointmentType.Single}>Single</option>
              <option value={API.EntitiesAppointmentType.Group}>Group</option>
              <option value={API.EntitiesAppointmentType.Party}>Party</option>
            </Select>
          </div>
        </FieldRow>
      </Field>
      <Field>
        <FieldLabel>Booking duration (minutes)</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><Timer size={16} /></IconSlot>
            <Input type="number" min={5} step={5} {...register('bookingDuration')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
      </Field>
      <div
        style={{
          display: 'grid',
          gap: 12,
          gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
        }}
      >
        <Field>
          <FieldLabel>Start date</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><Calendar size={16} /></IconSlot>
              <Input type="date" {...register('startDate')} style={{ paddingLeft: 36 }} />
            </div>
          </FieldRow>
        </Field>
        <Field>
          <FieldLabel>End date</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><Calendar size={16} /></IconSlot>
              <Input type="date" {...register('endDate')} style={{ paddingLeft: 36 }} />
            </div>
          </FieldRow>
        </Field>
      </div>
      <div
        style={{
          display: 'grid',
          gap: 12,
          gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
        }}
      >
        <Field>
          <FieldLabel>Start time</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><Clock size={16} /></IconSlot>
              <Input type="time" {...register('startTime')} style={{ paddingLeft: 36 }} />
            </div>
          </FieldRow>
        </Field>
        <Field>
          <FieldLabel>End time</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><Clock size={16} /></IconSlot>
              <Input type="time" {...register('endTime')} style={{ paddingLeft: 36 }} />
            </div>
          </FieldRow>
        </Field>
      </div>
      <Field>
        <FieldLabel>Max attendees (optional)</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><Users size={16} /></IconSlot>
            <Input type="number" min={1} {...register('maxAttendees')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
      </Field>
      <div>
        <Button variant="primary" disabled={pending}>
          {pending ? (<><Spinner /> Creating...</>) : 'Create appointment'}
        </Button>
      </div>
    </form>
  );
}
