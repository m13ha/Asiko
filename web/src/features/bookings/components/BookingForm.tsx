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
    <form onSubmit={submit} className="grid gap-3 p-4 border border-[var(--border)] rounded-xl bg-[var(--bg-elevated)]">
      <div className="grid gap-4 md:grid-cols-2">
        <Field>
          <FieldLabel>Name</FieldLabel>
          <FieldRow>
            <div className="relative w-full">
              <IconSlot><i className="pi pi-user" aria-hidden="true" /></IconSlot>
              <Input {...register('name')} autoComplete="name" className="pl-9" placeholder="John Doe" />
            </div>
          </FieldRow>
          {errors.name && <small className="text-xs text-[var(--danger)]">{errors.name.message}</small>}
        </Field>

        <Field>
          <FieldLabel>Email</FieldLabel>
          <FieldRow>
            <div className="relative w-full">
              <IconSlot><i className="pi pi-envelope" aria-hidden="true" /></IconSlot>
              <Input type="email" placeholder="john@example.com" {...register('email')} autoComplete="email" className="pl-9" />
            </div>
          </FieldRow>
          {errors.email && <small className="text-xs text-[var(--danger)]">{errors.email.message}</small>}
        </Field>

        <Field>
          <FieldLabel>Phone Number</FieldLabel>
          <FieldRow>
            <div className="relative w-full">
              <IconSlot><i className="pi pi-phone" aria-hidden="true" /></IconSlot>
              <Input placeholder="+1 (555) 000-0000" {...register('phone')} autoComplete="tel" className="pl-9" />
            </div>
          </FieldRow>
          {errors.phone && <small className="text-xs text-[var(--danger)]">{errors.phone.message}</small>}
        </Field>

        <Field>
          <FieldLabel>Attendees</FieldLabel>
          <FieldRow>
            <div className="relative w-full">
              <IconSlot><i className="pi pi-users" aria-hidden="true" /></IconSlot>
              <Input
                type="number"
                min={1}
                max={maxAttendees}
                {...register('attendeeCount', { valueAsNumber: true })}
                className="pl-9"
              />
            </div>
          </FieldRow>
          {errors.attendeeCount && <small className="text-xs text-[var(--danger)]">{errors.attendeeCount.message}</small>}
          {!errors.attendeeCount && maxAttendees && <small className="text-xs text-[var(--text-muted)]">{maxAttendees} spots remain for this slot.</small>}
        </Field>
      </div>

      <Field>
        <FieldLabel>Notes (optional)</FieldLabel>
        <FieldRow>
          <div className="relative w-full">
            <Textarea {...register('description')} rows={3} className="w-full p-2" placeholder="Anything the host should know?" />
          </div>
        </FieldRow>
      </Field>

      <div className="flex justify-end w-full">
        <Button variant="primary" disabled={pending} size="lg">
          {pending ? (<><Spinner /> Confirming...</>) : 'Confirm booking'}
        </Button>
      </div>
    </form>
  );
}
