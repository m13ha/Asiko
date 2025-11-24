import { useEffect, useMemo, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import styled from 'styled-components';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { differenceInCalendarDays } from 'date-fns';
import { Steps } from 'primereact/steps';
import type { StepsSelectEvent } from 'primereact/steps';
import type { MenuItem } from 'primereact/menuitem';
import { Input } from '@/components/Input';
import { Textarea } from '@/components/Textarea';
import { Select } from '@/components/Select';
import { Button } from '@/components/Button';
import { Spinner } from '@/components/Spinner';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import * as API from '@appointment-master/api-client';

const parseDateOnly = (value?: string) => {
  if (!value) return null;
  const d = new Date(`${value}T00:00:00Z`);
  return Number.isNaN(d.getTime()) ? null : d;
};

const normalizeTimeString = (time: string) => {
  if (!time.includes(':')) return `${time}:00`;
  const parts = time.split(':');
  if (parts.length === 2) return `${parts[0]}:${parts[1]}:00`;
  return `${parts[0]}:${parts[1]}:${parts[2] || '00'}`;
};

const parseTimeOnly = (value?: string) => {
  if (!value) return null;
  const t = new Date(`1970-01-01T${normalizeTimeString(value)}Z`);
  return Number.isNaN(t.getTime()) ? null : t;
};

const schema = z
  .object({
    title: z.string().min(1, 'Please provide a title'),
    description: z.string().optional(),
    type: z.enum(
      [API.EntitiesAppointmentType.Single, API.EntitiesAppointmentType.Group, API.EntitiesAppointmentType.Party],
      { required_error: 'Select an appointment type' }
    ),
    bookingDuration: z.coerce.number().min(5, 'Duration must be at least 5 minutes'),
    startDate: z.string().min(1, 'Select a start date'),
    endDate: z.string().min(1, 'Select an end date'),
    startTime: z.string().min(1, 'Select a start time'),
    endTime: z.string().min(1, 'Select an end time'),
    maxAttendees: z.coerce.number().min(1, 'Capacity must be at least 1').optional(),
    antiScalpingLevel: z
      .nativeEnum(API.EntitiesAntiScalpingLevel)
      .default(API.EntitiesAntiScalpingLevel.ScalpingStandard),
  })
  .superRefine((value, ctx) => {
    if (value.type !== API.EntitiesAppointmentType.Single && !value.maxAttendees) {
      ctx.addIssue({
        path: ['maxAttendees'],
        code: z.ZodIssueCode.custom,
        message: 'Capacity is required for group or party appointments',
      });
    }

    const startDate = parseDateOnly(value.startDate);
    const endDate = parseDateOnly(value.endDate);
    if (startDate && endDate && endDate < startDate) {
      ctx.addIssue({
        path: ['endDate'],
        code: z.ZodIssueCode.custom,
        message: 'End date cannot be before start date',
      });
    }

    const startTime = parseTimeOnly(value.startTime);
    const endTime = parseTimeOnly(value.endTime);
    if (startTime && endTime) {
      if (endTime <= startTime) {
        ctx.addIssue({
          path: ['endTime'],
          code: z.ZodIssueCode.custom,
          message: 'End time must be after start time',
        });
      } else {
        const windowMinutes = (endTime.getTime() - startTime.getTime()) / 60000;
        if (value.bookingDuration > windowMinutes) {
          ctx.addIssue({
            path: ['bookingDuration'],
            code: z.ZodIssueCode.custom,
            message: 'Booking duration exceeds the daily time window',
          });
        }
      }
    }
  });

export type AppointmentFormValues = z.infer<typeof schema>;

const steps = [
  { key: 'basics', label: 'Basics', fields: ['title', 'type', 'description'] },
  { key: 'schedule', label: 'Schedule', fields: ['startDate', 'endDate', 'startTime', 'endTime', 'bookingDuration'] },
  { key: 'capacity', label: 'Capacity & Review', fields: ['maxAttendees', 'antiScalpingLevel'] },
] as const;

const appointmentTypeOptions = [
  { label: 'Single (1:1)', value: API.EntitiesAppointmentType.Single },
  { label: 'Group (per-slot capacity)', value: API.EntitiesAppointmentType.Group },
  { label: 'Party (shared capacity)', value: API.EntitiesAppointmentType.Party },
];

const antiScalpingOptions = [
  { label: 'None – open access', value: API.EntitiesAntiScalpingLevel.ScalpingNone },
  { label: 'Standard – device checks', value: API.EntitiesAntiScalpingLevel.ScalpingStandard },
  { label: 'Strict – owner approval', value: API.EntitiesAntiScalpingLevel.ScalpingStrict },
];

export function AppointmentForm({ onSubmit, pending }: { onSubmit: (v: AppointmentFormValues) => void; pending?: boolean }) {
  const [currentStep, setCurrentStep] = useState(0);
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
    trigger,
  } = useForm<AppointmentFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      type: API.EntitiesAppointmentType.Single,
      bookingDuration: 30,
      antiScalpingLevel: API.EntitiesAntiScalpingLevel.ScalpingStandard,
    },
  });

  const type = watch('type');
  const startDate = watch('startDate');
  const endDate = watch('endDate');
  const startTime = watch('startTime');
  const endTime = watch('endTime');
  const bookingDuration = watch('bookingDuration');
  const maxAttendees = watch('maxAttendees');

  useEffect(() => {
    if (startDate && !endDate) {
      setValue('endDate', startDate);
    }
  }, [startDate, endDate, setValue]);

  useEffect(() => {
    if (type === API.EntitiesAppointmentType.Single && maxAttendees) {
      setValue('maxAttendees', undefined);
    }
  }, [type, maxAttendees, setValue]);

  const timezone = useMemo(() => Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC', []);

  const summary = useMemo(() => {
    if (!startDate || !endDate || !startTime || !endTime || !bookingDuration) {
      return { days: null, windowMinutes: null, slotsPerDay: null };
    }

    const days = differenceInCalendarDays(new Date(`${endDate}T00:00:00Z`), new Date(`${startDate}T00:00:00Z`)) + 1;
    const startDateTime = new Date(`${startDate}T${startTime}`);
    const endDateTime = new Date(`${startDate}T${endTime}`);
    const windowMinutes = Math.max(0, (endDateTime.getTime() - startDateTime.getTime()) / 60000);
    const slotsPerDay = windowMinutes > 0 ? Math.floor(windowMinutes / bookingDuration) : 0;

    return { days, windowMinutes, slotsPerDay };
  }, [startDate, endDate, startTime, endTime, bookingDuration]);

  const stepItems = useMemo<MenuItem[]>(
    () =>
      steps.map((step, index) => ({
        id: step.key,
        label: step.label,
        disabled: index > currentStep + 1,
      })),
    [currentStep]
  );

  const handleStepSelect = (event: StepsSelectEvent) => {
    if (event.index <= currentStep) {
      setCurrentStep(event.index);
    }
  };

  const handleNext = async () => {
    const fields = [...steps[currentStep].fields] as (keyof AppointmentFormValues)[];
    const ok = await trigger(fields, { shouldFocus: true });
    if (!ok) return;
    setCurrentStep((step) => Math.min(step + 1, steps.length - 1));
  };

  const handleBack = () => setCurrentStep((step) => Math.max(step - 1, 0));

  const isLastStep = currentStep === steps.length - 1;

  return (
    <StyledForm onSubmit={handleSubmit(onSubmit)}>
      <StepperWrapper>
        <Steps
          model={stepItems}
          activeIndex={currentStep}
          onSelect={handleStepSelect}
          readOnly={false}
          className="appointment-steps"
        />
      </StepperWrapper>

      {currentStep === 0 && (
        <FormSection>
        <SectionHeader>
          <SectionTitle>Basics</SectionTitle>
          <Helper>Tell attendees what this booking is about.</Helper>
        </SectionHeader>
        <GridTwo>
          <Field>
            <FieldLabel>Title</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-list" aria-hidden="true" /></IconSlot>
                <Input {...register('title')} style={{ paddingLeft: 36 }} placeholder="Consultation, Office Hours..." />
              </div>
            </FieldRow>
            {errors.title && <ErrorText>{errors.title.message}</ErrorText>}
          </Field>
          <Field>
            <FieldLabel>Type</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-list" aria-hidden="true" /></IconSlot>
                <Controller
                  name="type"
                  control={control}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      onChange={(event) => field.onChange(event.value)}
                      onBlur={field.onBlur}
                      name={field.name}
                      options={appointmentTypeOptions}
                      optionLabel="label"
                      optionValue="value"
                      pt={{ input: { style: { paddingLeft: '36px' } } }}
                    />
                  )}
                />
              </div>
            </FieldRow>
          </Field>
        </GridTwo>
        <Field>
          <FieldLabel>Description <small>(optional)</small></FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><i className="pi pi-file" aria-hidden="true" /></IconSlot>
              <Textarea {...register('description')} style={{ paddingLeft: 36 }} placeholder="Share agenda, location, prep info..." />
            </div>
          </FieldRow>
        </Field>
        </FormSection>
      )}

      {currentStep === 1 && (
        <FormSection>
        <SectionHeader>
          <SectionTitle>Schedule</SectionTitle>
          <Helper>Times are saved in {timezone}. Attendees convert automatically.</Helper>
        </SectionHeader>
        <GridTwo>
          <Field>
            <FieldLabel>Start date</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-calendar" aria-hidden="true" /></IconSlot>
                <Input type="date" {...register('startDate')} style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
            {errors.startDate && <ErrorText>{errors.startDate.message}</ErrorText>}
          </Field>
          <Field>
            <FieldLabel>End date</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-calendar" aria-hidden="true" /></IconSlot>
                <Input type="date" {...register('endDate')} style={{ paddingLeft: 36 }} disabled={!startDate} min={startDate} />
              </div>
            </FieldRow>
            {errors.endDate && <ErrorText>{errors.endDate.message}</ErrorText>}
          </Field>
        </GridTwo>
        <GridTwo>
          <Field>
            <FieldLabel>Start time</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-clock" aria-hidden="true" /></IconSlot>
                <Input type="time" {...register('startTime')} style={{ paddingLeft: 36 }} />
              </div>
            </FieldRow>
            {errors.startTime && <ErrorText>{errors.startTime.message}</ErrorText>}
          </Field>
          <Field>
            <FieldLabel>End time</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-clock" aria-hidden="true" /></IconSlot>
                <Input type="time" {...register('endTime')} style={{ paddingLeft: 36 }} disabled={!startTime} min={startTime} />
              </div>
            </FieldRow>
            {errors.endTime && <ErrorText>{errors.endTime.message}</ErrorText>}
          </Field>
        </GridTwo>
        <Field>
          <FieldLabel>Booking duration (minutes)</FieldLabel>
          <FieldRow>
            <div style={{ position: 'relative' }}>
              <IconSlot><i className="pi pi-stopwatch" aria-hidden="true" /></IconSlot>
              <Input type="number" min={5} step={5} {...register('bookingDuration', { valueAsNumber: true })} style={{ paddingLeft: 36 }} />
            </div>
          </FieldRow>
          <Helper>Slot length determines how many bookings fit inside your daily window.</Helper>
          {errors.bookingDuration && <ErrorText>{errors.bookingDuration.message}</ErrorText>}
        </Field>
        </FormSection>
      )}

      {currentStep === 2 && (
        <FormSection>
        <SectionHeader>
          <SectionTitle>Capacity & Rules</SectionTitle>
          <Helper>Keep control over how people book your time.</Helper>
        </SectionHeader>
        <GridTwo>
          {(type === API.EntitiesAppointmentType.Group || type === API.EntitiesAppointmentType.Party) && (
            <Field>
              <FieldLabel>Max attendees per slot</FieldLabel>
              <FieldRow>
                <div style={{ position: 'relative' }}>
                  <IconSlot><i className="pi pi-users" aria-hidden="true" /></IconSlot>
                  <Input type="number" min={1} {...register('maxAttendees', { valueAsNumber: true })} style={{ paddingLeft: 36 }} placeholder="e.g. 5" />
                </div>
              </FieldRow>
              <Helper>Only shown to invitees if capacity is limited.</Helper>
              {errors.maxAttendees && <ErrorText>{errors.maxAttendees.message}</ErrorText>}
            </Field>
          )}
          <Field>
            <FieldLabel>Anti-scalping level</FieldLabel>
            <FieldRow>
              <div style={{ position: 'relative' }}>
                <IconSlot><i className="pi pi-shield" aria-hidden="true" /></IconSlot>
                <Controller
                  name="antiScalpingLevel"
                  control={control}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      onChange={(event) => field.onChange(event.value)}
                      onBlur={field.onBlur}
                      name={field.name}
                      options={antiScalpingOptions}
                      optionLabel="label"
                      optionValue="value"
                      pt={{ input: { style: { paddingLeft: '36px' } } }}
                    />
                  )}
                />
              </div>
            </FieldRow>
          </Field>
        </GridTwo>
        <SummaryCard>
          <SummaryHeader>
            <i className="pi pi-info-circle" aria-hidden="true" />
            <strong>Appointment Summary</strong>
          </SummaryHeader>
          <SummaryRow>
            <span>Date range</span>
            <strong>{startDate && endDate ? `${startDate} → ${endDate}` : 'Pending selection'}</strong>
          </SummaryRow>
          <SummaryRow>
            <span>Window per day</span>
            <strong>
              {summary.windowMinutes ? `${summary.windowMinutes} mins` : 'Choose start/end times'}
            </strong>
          </SummaryRow>
          <SummaryRow>
            <span>Slots per day</span>
            <strong>{summary.slotsPerDay ? summary.slotsPerDay : '-'}</strong>
          </SummaryRow>
          <SummaryRow>
            <span>Total days</span>
            <strong>{summary.days ? summary.days : '-'}</strong>
          </SummaryRow>
          <SummaryRow>
            <span>Guest experience</span>
            <strong>{type === API.EntitiesAppointmentType.Single ? '1:1 bookings' : 'Shared slots'}</strong>
          </SummaryRow>
        </SummaryCard>
      </FormSection>
      )}

      <StickyActions>
        {currentStep > 0 && (
          <Button type="button" onClick={handleBack}>
            Back
          </Button>
        )}
        {!isLastStep && (
          <Button type="button" variant="primary" onClick={handleNext}>
            Next
          </Button>
        )}
        {isLastStep && (
          <Button variant="primary" disabled={pending} style={{ minWidth: 200 }}>
            {pending ? (
              <>
                <Spinner /> Saving…
              </>
            ) : (
              'Create appointment'
            )}
          </Button>
        )}
      </StickyActions>
    </StyledForm>
  );
}

const StyledForm = styled.form`
  display: grid;
  gap: 20px;
`;

const FormSection = styled.section`
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg-elevated);
  padding: 18px;
  display: grid;
  gap: 20px;
`;

const SectionHeader = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

const SectionTitle = styled.h3`
  margin: 0;
  font-size: 16px;
`;

const Helper = styled.p`
  margin: 0;
  font-size: 12px;
  color: var(--text-muted);
`;

const GridTwo = styled.div`
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
`;

const ErrorText = styled.small`
  color: var(--danger);
`;

const SummaryCard = styled.div`
  border: 1px dashed color-mix(in oklab, var(--primary) 30%, var(--border));
  border-radius: var(--radius);
  padding: 12px;
  background: color-mix(in oklab, var(--primary) 4%, var(--bg));
  display: grid;
  gap: 8px;
`;

const SummaryHeader = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text);
  font-size: 14px;
`;

const SummaryRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  color: var(--text);
`;

const StickyActions = styled.div`
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  flex-wrap: wrap;
  margin-top: 12px;
`;

const StepperWrapper = styled.div`
  padding: 8px 0 4px;
`;
