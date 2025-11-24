import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Input } from '@/components/Input';
import { Textarea } from '@/components/Textarea';
import { Button } from '@/components/Button';
import { Spinner } from '@/components/Spinner';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Badge } from '@/components/Badge';

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
    <form onSubmit={submit} className="form-grid booking-form">
      <div className="form-grid-two">
        <Field>
          <FieldLabel>Name</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><i className="pi pi-user" aria-hidden="true" /></IconSlot>
              <Input {...register('name')} autoComplete="name" style={{ paddingLeft: 36 }} />
            </div>
          </FieldRow>
          {errors.name && <small className="field-error">{errors.name.message}</small>}
        </Field>
        <Field>
          <FieldLabel>
            Contact <Badge tone="muted">email or phone</Badge>
          </FieldLabel>
          <FieldRow style={{ gap: 8, flexWrap: 'wrap' }}>
            <div style={{ position: 'relative', flex: '1 1 180px' }}>
              <IconSlot><i className="pi pi-envelope" aria-hidden="true" /></IconSlot>
              <Input type="email" placeholder="Email" {...register('email')} autoComplete="email" style={{ paddingLeft: 36 }} />
            </div>
            <div style={{ position: 'relative', flex: '1 1 140px' }}>
              <IconSlot><i className="pi pi-phone" aria-hidden="true" /></IconSlot>
              <Input placeholder="Phone" {...register('phone')} autoComplete="tel" style={{ paddingLeft: 36 }} />
            </div>
          </FieldRow>
          {errors.email && <small className="field-error">{errors.email.message}</small>}
        </Field>
      </div>

      <div className="form-grid-two">
        <Field>
          <FieldLabel>Attendees</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative', maxWidth: 220 }}>
              <IconSlot><i className="pi pi-users" aria-hidden="true" /></IconSlot>
              <Input
                type="number"
                min={1}
                max={maxAttendees}
                {...register('attendeeCount', { valueAsNumber: true })}
                style={{ paddingLeft: 36 }}
              />
            </div>
          </FieldRow>
          {errors.attendeeCount && <small className="field-error">{errors.attendeeCount.message}</small>}
          {!errors.attendeeCount && maxAttendees && <small className="field-hint">{maxAttendees} spots remain for this slot.</small>}
        </Field>

        <Field>
          <FieldLabel>Notes (optional)</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><i className="pi pi-file" aria-hidden="true" /></IconSlot>
              <Textarea {...register('description')} rows={3} style={{ paddingLeft: 36 }} placeholder="Anything the host should know?" />
            </div>
          </FieldRow>
        </Field>
      </div>

      <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Button variant="primary" disabled={pending} size="lg">
          {pending ? (<><Spinner /> Confirming...</>) : 'Confirm booking'}
        </Button>
      </div>
    </form>
  );
}
