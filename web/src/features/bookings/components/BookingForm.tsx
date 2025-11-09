import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/Input';
import { Textarea } from '@/components/Textarea';
import { Button } from '@/components/Button';
import { Spinner } from '@/components/Spinner';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { User, Mail, Phone, Users, FileText } from 'lucide-react';

const schema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.string().email().optional(),
  phone: z.string().optional(),
  attendeeCount: z.coerce.number().min(1).default(1),
  description: z.string().optional(),
}).refine((v) => !!v.email || !!v.phone, { message: 'Email or phone is required', path: ['email'] });

export type BookingFormValues = z.infer<typeof schema>;

export function BookingForm({ onSubmit, pending, maxAttendees }: { onSubmit: (v: BookingFormValues) => void; pending?: boolean; maxAttendees?: number }) {
  const { register, handleSubmit, formState: { errors }, setError, clearErrors } = useForm<BookingFormValues>({ resolver: zodResolver(schema), defaultValues: { attendeeCount: 1 } });

  const submit = handleSubmit((values) => {
    if (maxAttendees && values.attendeeCount > maxAttendees) {
      setError('attendeeCount', { type: 'validate', message: `Only ${maxAttendees} spot${maxAttendees === 1 ? '' : 's'} remain for this slot` });
      return;
    }
    clearErrors('attendeeCount');
    onSubmit(values);
  });

  return (
    <form onSubmit={submit} className="form-grid">
      <Field>
        <FieldLabel>Name</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><User size={16} /></IconSlot>
            <Input {...register('name')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
        {errors.name && <small style={{ color: 'var(--danger)' }}>{errors.name.message}</small>}
      </Field>
      <Field>
        <FieldLabel>Email</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><Mail size={16} /></IconSlot>
            <Input type="email" {...register('email')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
      </Field>
      <Field>
        <FieldLabel>Phone</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><Phone size={16} /></IconSlot>
            <Input {...register('phone')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
      </Field>
      <Field>
        <FieldLabel>Attendees</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><Users size={16} /></IconSlot>
            <Input
              type="number"
              min={1}
              max={maxAttendees}
              {...register('attendeeCount', { valueAsNumber: true })}
              style={{ paddingLeft: 36 }}
            />
          </div>
        </FieldRow>
        {errors.attendeeCount && <small style={{ color: 'var(--danger)' }}>{errors.attendeeCount.message}</small>}
      </Field>
      <Field>
        <FieldLabel>Notes (optional)</FieldLabel>
        <FieldRow>
          <div style={{ position: 'relative' }}>
            <IconSlot><FileText size={16} /></IconSlot>
            <Textarea {...register('description')} style={{ paddingLeft: 36 }} />
          </div>
        </FieldRow>
      </Field>
      <div>
        <Button variant="primary" disabled={pending}>
          {pending ? (<><Spinner /> Confirming...</>) : 'Confirm booking'}
        </Button>
      </div>
    </form>
  );
}
