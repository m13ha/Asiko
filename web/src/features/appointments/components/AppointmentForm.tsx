import { useEffect, useMemo, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';

import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { addDays, differenceInCalendarDays, differenceInMinutes, format, isBefore, startOfDay } from 'date-fns';
import { Stepper } from '@/components/Stepper';
import { Input } from '@/components/Input';
import { Textarea } from '@/components/Textarea';
import { Select } from '@/components/Select';
import { Button } from '@/components/Button';
import { Spinner } from '@/components/Spinner';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { TimePicker } from '@/components/TimePicker';
import { DatePicker } from '@/components/DatePicker';
import * as API from '@appointment-master/api-client';
const SectionTitle = ({ children }: { children: React.ReactNode }) => (
  <h3 className="m-0 text-base font-semibold text-[var(--text)]">{children}</h3>
);

const Helper = ({ children, className = '' }: { children: React.ReactNode; className?: string }) => (
  <p className={`m-0 text-xs text-[var(--text-muted)] ${className}`}>{children}</p>
);

const FormSection = ({ children }: { children: React.ReactNode }) => (
  <section className="p-5 sm:p-8 bg-[var(--bg-elevated)] border-y sm:border border-[var(--border)] sm:rounded-2xl shadow-[var(--elev-1)]">
    {children}
  </section>
);

const SummaryBox = ({ children }: { children: React.ReactNode }) => (
  <div className="p-4 border border-dashed border-[color-mix(in_oklab,var(--primary)_30%,var(--border))] rounded-xl bg-[color-mix(in_oklab,var(--primary)_4%,var(--bg-elevated))] grid gap-3">
    {children}
  </div>
);

const parseDateOnly = (value?: string) => {
  if (!value) return null;
  const d = new Date(`${value}T00:00:00`);
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
  const t = new Date(`1970-01-01T${normalizeTimeString(value)}`);
  return Number.isNaN(t.getTime()) ? null : t;
};

const combineDateTime = (date: Date, time: Date) =>
  new Date(date.getFullYear(), date.getMonth(), date.getDate(), time.getHours(), time.getMinutes(), time.getSeconds(), time.getMilliseconds());

const getDailyWindowMinutes = (startTime: Date, endTime: Date) => {
  const startMinutes = startTime.getHours() * 60 + startTime.getMinutes();
  const endMinutes = endTime.getHours() * 60 + endTime.getMinutes();
  if (endMinutes > startMinutes) {
    return endMinutes - startMinutes;
  }
  return 24 * 60 - startMinutes + endMinutes;
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
    const today = startOfDay(new Date());
    if (startDate && isBefore(startDate, today)) {
      ctx.addIssue({
        path: ['startDate'],
        code: z.ZodIssueCode.custom,
        message: 'Start date cannot be in the past',
      });
    }
    if (startDate && endDate && endDate < startDate) {
      ctx.addIssue({
        path: ['endDate'],
        code: z.ZodIssueCode.custom,
        message: 'End date cannot be before start date',
      });
    }

    const startTime = parseTimeOnly(value.startTime);
    const endTime = parseTimeOnly(value.endTime);
    if (startDate && endDate && startTime && endTime) {
      const startDateTime = combineDateTime(startDate, startTime);
      const endDateTime = combineDateTime(endDate, endTime);
      const now = new Date();
      if (startDateTime < now) {
        ctx.addIssue({
          path: ['startTime'],
          code: z.ZodIssueCode.custom,
          message: 'Start time cannot be in the past',
        });
      }

      if (endDateTime <= startDateTime) {
        ctx.addIssue({
          path: ['endTime'],
          code: z.ZodIssueCode.custom,
          message: 'End time must be after start time',
        });
      }

      const dailyWindowMinutes = getDailyWindowMinutes(startTime, endTime);
      if (value.type !== API.EntitiesAppointmentType.Party && value.bookingDuration > dailyWindowMinutes) {
        ctx.addIssue({
          path: ['bookingDuration'],
          code: z.ZodIssueCode.custom,
          message: 'Booking duration exceeds the daily time window',
        });
      }

      if (value.type === API.EntitiesAppointmentType.Party) {
        if (differenceInCalendarDays(endDate, startDate) > 1) {
          ctx.addIssue({
            path: ['endDate'],
            code: z.ZodIssueCode.custom,
            message: 'Party appointments can only span one overnight window',
          });
        }

        const fullWindowMinutes = differenceInMinutes(endDateTime, startDateTime);
        if (fullWindowMinutes > 24 * 60) {
          ctx.addIssue({
            path: ['endDate'],
            code: z.ZodIssueCode.custom,
            message: 'Party appointments cannot exceed 24 hours',
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

export function AppointmentForm({
  onSubmit,
  pending,
  initialValues,
  submitLabel = 'Create appointment',
}: {
  onSubmit: (v: AppointmentFormValues) => void;
  pending?: boolean;
  initialValues?: Partial<AppointmentFormValues>;
  submitLabel?: string;
}) {
  const [currentStep, setCurrentStep] = useState(0);
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
    trigger,
    reset,
  } = useForm<AppointmentFormValues>({
    resolver: zodResolver(schema),
    shouldUnregister: false,
    defaultValues: {
      type: API.EntitiesAppointmentType.Single,
      bookingDuration: 30,
      antiScalpingLevel: API.EntitiesAntiScalpingLevel.ScalpingStandard,
      ...initialValues,
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
    if (!initialValues) return;
    reset({
      type: API.EntitiesAppointmentType.Single,
      bookingDuration: 30,
      antiScalpingLevel: API.EntitiesAntiScalpingLevel.ScalpingStandard,
      ...initialValues,
    });
  }, [initialValues, reset]);

  useEffect(() => {
    if (!startDate) return;
    if (type !== API.EntitiesAppointmentType.Party) {
      if (!endDate) {
        setValue('endDate', startDate);
      }
      return;
    }

    const start = parseDateOnly(startDate);
    if (!start) return;
    const startClock = parseTimeOnly(startTime);
    const endClock = parseTimeOnly(endTime);
    const shouldOvernight = startClock && endClock ? getDailyWindowMinutes(startClock, endClock) > 0 && endClock <= startClock : false;
    const targetEnd = shouldOvernight ? addDays(start, 1) : start;
    const formatted = format(targetEnd, 'yyyy-MM-dd');
    if (endDate !== formatted) {
      setValue('endDate', formatted);
    }
  }, [startDate, endDate, startTime, endTime, type, setValue]);

  useEffect(() => {
    if (type === API.EntitiesAppointmentType.Single && maxAttendees) {
      setValue('maxAttendees', undefined);
    }
  }, [type, maxAttendees, setValue]);

  useEffect(() => {
    if (type !== API.EntitiesAppointmentType.Party) return;
    const start = parseDateOnly(startDate);
    const end = parseDateOnly(endDate);
    const startClock = parseTimeOnly(startTime);
    const endClock = parseTimeOnly(endTime);
    if (!start || !end || !startClock || !endClock) return;
    const windowMinutes = differenceInMinutes(combineDateTime(end, endClock), combineDateTime(start, startClock));
    if (windowMinutes > 0 && bookingDuration !== windowMinutes) {
      setValue('bookingDuration', windowMinutes);
    }
  }, [type, startDate, endDate, startTime, endTime, bookingDuration, setValue]);

  const timezone = useMemo(() => Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC', []);
  const labelClass = 'text-sm font-semibold text-[var(--text)]';
  const inputClass = 'pl-10 py-2.5 text-sm w-full';
  const textareaClass = 'w-full min-h-[120px] border p-3 text-sm';
  
  const timeValue = (value?: string) => parseTimeOnly(value);
  const toTimeString = (value: Date | null) => (value ? format(value, 'HH:mm') : '');
  const sameDay = Boolean(startDate && endDate && startDate === endDate);

  const summary = useMemo(() => {
    if (!startDate || !endDate || !startTime || !endTime || !bookingDuration) {
      return { days: null, windowMinutes: null, slotsPerDay: null };
    }

    const start = parseDateOnly(startDate);
    const end = parseDateOnly(endDate);
    const startClock = parseTimeOnly(startTime);
    const endClock = parseTimeOnly(endTime);
    if (!start || !end || !startClock || !endClock) {
      return { days: null, windowMinutes: null, slotsPerDay: null };
    }

    const days = differenceInCalendarDays(end, start) + 1;
    const dailyWindowMinutes = getDailyWindowMinutes(startClock, endClock);
    const slotsPerDay = dailyWindowMinutes > 0 ? Math.floor(dailyWindowMinutes / bookingDuration) : 0;

    return { days, windowMinutes: dailyWindowMinutes, slotsPerDay };
  }, [startDate, endDate, startTime, endTime, bookingDuration]);


  const handleStepSelect = (index: number) => {
    if (index <= currentStep) {
      setCurrentStep(index);
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
    <div className="py-4 sm:py-8 px-4 sm:px-4">
      <form onSubmit={handleSubmit(onSubmit)} className="grid gap-6 sm:gap-8">
      <div className="mb-4 sm:mb-8">
        <Stepper
          steps={steps}
          activeStep={currentStep}
          onStepClick={handleStepSelect}
        />
      </div>

        <div className={currentStep === 0 ? 'block' : 'hidden'} aria-hidden={currentStep !== 0}>
          <FormSection>
            <div className="flex flex-col gap-1 mb-6">
              <SectionTitle>Basics</SectionTitle>
              <Helper>Tell attendees what this booking is about.</Helper>
            </div>
            
            <div className="grid gap-6 md:grid-cols-2 mb-6">
              <Field>
              <FieldLabel className={labelClass}>Title</FieldLabel>
              <FieldRow>
                <IconSlot className="left-3"><i className="pi pi-bookmark" aria-hidden="true" /></IconSlot>
                <Input {...register('title')} className={inputClass} placeholder="Consultation, Office Hours..." />
              </FieldRow>
              {errors.title && <small className="text-red-500">{errors.title.message}</small>}
            </Field>
            
            <Field>
              <FieldLabel className={labelClass}>Type</FieldLabel>
              <FieldRow>
                <IconSlot className="left-3"><i className="pi pi-tag" aria-hidden="true" /></IconSlot>
                <div className="w-full">
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
                        className="w-full"
                      />
                    )}
                  />
                </div>
              </FieldRow>
            </Field>
          </div>

          <Field>
            <FieldLabel className={labelClass}>Description <small className="text-[var(--text-muted)] font-medium">(optional)</small></FieldLabel>
            <FieldRow>
              <Textarea {...register('description')} className={textareaClass} placeholder="Share agenda, location, prep info..." />
            </FieldRow>
          </Field>
          </FormSection>
        </div>

        <div className={currentStep === 1 ? 'block' : 'hidden'} aria-hidden={currentStep !== 1}>
          <FormSection>
            <div className="flex flex-col gap-1 mb-6">
              <SectionTitle>Schedule</SectionTitle>
              <Helper>Times are saved in {timezone}. Attendees convert automatically.</Helper>
            </div>
            
            <div className="grid gap-6 md:grid-cols-2 mb-6">
              <Field>
              <FieldLabel className={labelClass}>Start date</FieldLabel>
              <FieldRow>
                <div className="w-full">
                  <Controller
                    name="startDate"
                    control={control}
                    render={({ field }) => (
                      <DatePicker
                        value={field.value}
                        onChange={(date) => field.onChange(date ? format(date, 'yyyy-MM-dd') : '')}
                        placeholder="Select start date"
                        minDate={new Date()}
                        className="w-full"
                      />
                    )}
                  />
                </div>
              </FieldRow>
              {errors.startDate && <small className="text-red-500">{errors.startDate.message}</small>}
            </Field>
            
            <Field>
              <FieldLabel className={labelClass}>End date</FieldLabel>
              <FieldRow>
                <div className="w-full">
                  <Controller
                    name="endDate"
                    control={control}
                    render={({ field }) => (
                        <DatePicker
                          value={field.value}
                          onChange={(date) => field.onChange(date ? format(date, 'yyyy-MM-dd') : '')}
                          placeholder="Select end date"
                          disabled={!startDate}
                          minDate={startDate}
                          maxDate={
                            type === API.EntitiesAppointmentType.Party && startDate
                              ? addDays(new Date(`${startDate}T00:00:00`), 1)
                              : undefined
                          }
                          className="w-full"
                        />
                      )}
                    />
                </div>
              </FieldRow>
              {errors.endDate && <small className="text-red-500">{errors.endDate.message}</small>}
            </Field>
          </div>

            <div className="grid gap-6 md:grid-cols-2 mb-6">
              <Field>
                <FieldLabel className={labelClass}>Start time</FieldLabel>
                <FieldRow>
                  <IconSlot className="left-3"><i className="pi pi-clock" aria-hidden="true" /></IconSlot>
                  <div className="w-full">
                    <Controller
                      name="startTime"
                      control={control}
                      render={({ field }) => (
                        <TimePicker
                          value={timeValue(field.value)}
                          onChange={(value) => field.onChange(toTimeString(value))}
                          placeholder="Select start time"
                          className="w-full"
                        />
                      )}
                    />
                  </div>
                </FieldRow>
                {errors.startTime && <small className="text-red-500">{errors.startTime.message}</small>}
              </Field>
              
              <Field>
                <FieldLabel className={labelClass}>End time</FieldLabel>
                <FieldRow>
                  <IconSlot className="left-3"><i className="pi pi-clock" aria-hidden="true" /></IconSlot>
                  <div className="w-full">
                    <Controller
                      name="endTime"
                      control={control}
                      render={({ field }) => (
                      <TimePicker
                          value={timeValue(field.value)}
                          onChange={(value) => field.onChange(toTimeString(value))}
                          placeholder="Select end time"
                          disabled={!startTime}
                          minTime={sameDay ? timeValue(startTime) : undefined}
                          className="w-full"
                        />
                      )}
                    />
                  </div>
                </FieldRow>
                {errors.endTime && <small className="text-red-500">{errors.endTime.message}</small>}
              </Field>
            </div>

          {type !== API.EntitiesAppointmentType.Party && (
            <Field>
              <FieldLabel className={labelClass}>Booking duration (minutes)</FieldLabel>
              <FieldRow>
                <IconSlot className="left-3"><i className="pi pi-stopwatch" aria-hidden="true" /></IconSlot>
                <Input type="number" min={5} step={5} {...register('bookingDuration', { valueAsNumber: true })} className={inputClass} />
              </FieldRow>
              <Helper className="mt-1">Slot length determines how many bookings fit inside your daily window.</Helper>
              {errors.bookingDuration && <small className="text-red-500">{errors.bookingDuration.message}</small>}
            </Field>
          )}
          </FormSection>
        </div>

        <div className={currentStep === 2 ? 'block' : 'hidden'} aria-hidden={currentStep !== 2}>
          <FormSection>
            <div className="flex flex-col gap-1 mb-6">
              <SectionTitle>Capacity & Rules</SectionTitle>
              <Helper>Keep control over how people book your time.</Helper>
            </div>
            
            <div className="grid gap-6 md:grid-cols-2 mb-8">
              {(type === API.EntitiesAppointmentType.Group || type === API.EntitiesAppointmentType.Party) && (
                <Field>
                  <FieldLabel className={labelClass}>
                    {type === API.EntitiesAppointmentType.Party ? 'Total attendees for the party' : 'Max attendees per slot'}
                  </FieldLabel>
                  <FieldRow>
                    <IconSlot className="left-3"><i className="pi pi-users" aria-hidden="true" /></IconSlot>
                    <Input type="number" min={1} {...register('maxAttendees', { valueAsNumber: true })} className={inputClass} placeholder="e.g. 5" />
                  </FieldRow>
                  <Helper className="mt-1">Only shown to invitees if capacity is limited.</Helper>
                  {errors.maxAttendees && <small className="text-red-500">{errors.maxAttendees.message}</small>}
                </Field>
              )}
              
              <Field>
                <FieldLabel className={labelClass}>Anti-scalping level</FieldLabel>
                <FieldRow>
                  <IconSlot className="left-3"><i className="pi pi-shield" aria-hidden="true" /></IconSlot>
                  <div className="w-full">
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
                          className="w-full"
                        />
                      )}
                    />
                  </div>
                </FieldRow>
              </Field>
            </div>

            <SummaryBox>
              <div className="flex items-center gap-2 text-sm font-semibold mb-1 text-[var(--primary)]">
                <i className="pi pi-info-circle" aria-hidden="true" />
                Appointment Summary
              </div>
              <div className="flex justify-between items-center text-sm">
                <span className="text-[var(--text-muted)]">Date range</span>
                <strong className="text-[var(--text)]">{startDate && endDate ? `${startDate} → ${endDate}` : 'Pending selection'}</strong>
              </div>
              <div className="flex justify-between items-center text-sm">
                <span className="text-[var(--text-muted)]">Window per day</span>
                <strong className="text-[var(--text)]">
                  {summary.windowMinutes ? `${summary.windowMinutes} mins` : 'Choose start/end times'}
                </strong>
              </div>
              {type !== API.EntitiesAppointmentType.Party && (
                <div className="flex justify-between items-center text-sm">
                  <span className="text-[var(--text-muted)]">Slots per day</span>
                  <strong className="text-[var(--text)]">{summary.slotsPerDay ? summary.slotsPerDay : '-'}</strong>
                </div>
              )}
              <div className="flex justify-between items-center text-sm">
                <span className="text-[var(--text-muted)]">Total days</span>
                <strong className="text-[var(--text)]">{summary.days ? summary.days : '-'}</strong>
              </div>
              <div className="flex justify-between items-center text-sm">
                <span className="text-[var(--text-muted)]">Guest experience</span>
                <strong className="text-[var(--text)]">
                  {type === API.EntitiesAppointmentType.Single
                    ? '1:1 bookings'
                    : type === API.EntitiesAppointmentType.Party
                      ? 'Single shared slot'
                      : 'Shared slots'}
                </strong>
              </div>
            </SummaryBox>
          </FormSection>
        </div>

        <div className="flex flex-col-reverse sm:flex-row gap-3 justify-end mt-4">
          {currentStep > 0 && (
            <Button type="button" onClick={handleBack} variant="outline" className="w-full sm:w-auto">
              Back
            </Button>
          )}
          {!isLastStep && (
            <Button type="button" variant="primary" onClick={handleNext} className="w-full sm:w-auto">
              Next
            </Button>
          )}
          {isLastStep && (
            <Button variant="primary" disabled={pending} className="w-full sm:w-auto px-10">
              {pending ? (
                <>
                  <Spinner size="sm" /> Saving…
                </>
              ) : (
                submitLabel
              )}
            </Button>
          )}
        </div>
      </form>
    </div>
  );
}
