import { useEffect, useMemo, useState } from 'react';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useAvailableSlots, useAvailableSlotsByDay, useBookGuest, useBookRegistered } from '../hooks';
import { SlotPicker } from '../components/SlotPicker';
import { BookingForm, BookingFormValues } from '../components/BookingForm';
import { BookingSummary } from '../components/BookingSummary';
import { useAuth } from '@/features/auth/AuthProvider';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { SuccessBurst } from '@/components/SuccessBurst';
import type { EntitiesBooking, EntitiesAntiScalpingLevel } from '@appointment-master/api-client';
import { AvailabilityCalendar } from '../components/AvailabilityCalendar';
import { format, parseISO } from 'date-fns';
import { Badge } from '@/components/Badge';
import { useAppointmentByAppCode } from '@/features/appointments/hooks';
import { useDeviceToken } from '@/features/auth/hooks';
import FingerprintJS from '@sparkstone/fingerprintjs';
import './bookingSteps.css';

function normalizeDateOnly(value?: string | null) {
  if (!value) return null;
  try {
    // By replacing space with 'T' and appending 'Z', we are treating the date as UTC.
    // This prevents timezone shifts from accidentally changing the date.
    const isoDateString = value.trim().replace(' ', 'T');
    const date = parseISO(isoDateString.endsWith('Z') ? isoDateString : `${isoDatein-memory-contextg}Z`);
    if (Number.isNaN(date.getTime())) return null;
    return format(date, 'yyyy-MM-dd');
  } catch {
    return null;
  }
}

export function BookByCodePage() {
  const [step, setStep] = useState<1 | 2 | 3 | 4>(1);
  const [lastSuccessfulStep, setLastSuccessfulStep] = useState<1 | 2 | 3 | 4>(1);
  const [search] = useSearchParams();
  const prefill = search.get('code') || '';
  const [appCode, setAppCode] = useState(prefill);
  const [date, setDate] = useState('');
  const [slot, setSlot] = useState<EntitiesBooking | null>(null);
  const [registeredCount, setRegisteredCount] = useState(1);
  const [deviceToken, setDeviceToken] = useState<string | null>(null);
  const isValidCode = appCode.trim().length > 0;
  const { isAuthed } = useAuth();
  const navigate = useNavigate();

  // Fetch appointment details to check anti-scalping level
  const appointmentDetails = useAppointmentByAppCode(appCode);
  const allSlots = useAvailableSlots(isValidCode ? appCode : '');
  const slots = useAvailableSlotsByDay(isValidCode ? appCode : '', date);
  const bookGuest = useBookGuest();
  const bookReg = useBookRegistered();
  const generateDeviceToken = useDeviceToken();
  const [showBurst, setShowBurst] = useState(false);

  const proceedToSummary = () => setStep(4);

  const spotsRemaining = slot ? Math.max((slot.capacity ?? slot.attendeeCount ?? 1) - (slot.seatsBooked ?? 0), 0) : 0;
  const normalizedDateIso = date ? `${date}T00:00:00Z` : '';

  useEffect(() => {
    setSlot(null);
    setRegisteredCount(1);
  }, [date]);

  useEffect(() => {
    if (allSlots.error) {
      setStep(1);
      setLastSuccessfulStep(1);
    }
  }, [allSlots.error]);

  useEffect(() => {
    if (slots.error && date) {
      setStep(2);
      setLastSuccessfulStep(2);
    }
  }, [slots.error, date]);

  useEffect(() => {
    if (slot) {
      setRegisteredCount((prev) => {
        const next = spotsRemaining > 0 ? spotsRemaining : 1;
        return Math.min(prev, next);
      });
    }
  }, [slot, spotsRemaining]);

  const availableDates = useMemo(() => {
    if (!allSlots.data?.items?.length) return [];
    const unique = new Set<string>();
    for (const item of allSlots.data.items) {
      const normalized = normalizeDateOnly(item.date);
      if (normalized) {
        unique.add(normalized);
      }
    }
    return Array.from(unique).sort();
  }, [allSlots.data?.items]);

  useEffect(() => {
    if (prefill) {
      setAppCode(prefill);
      setStep(2);
    }
  }, [prefill]);

  // Generate device token if the appointment requires strict anti-scalping
  useEffect(() => {
    if (appointmentDetails.data && appointmentDetails.data.antiScalpingLevel === 'strict') {
      // Generate a unique device ID using FingerprintJS
      FingerprintJS.load()
        .then(fp => fp.get())
        .then(result => {
          const uniqueDeviceId = result.visitorId; // Use the visitorId as the device ID

          generateDeviceToken.mutate(
            { deviceId: uniqueDeviceId },
            {
              onSuccess: (response) => {
                if (response.device_token) {
                  setDeviceToken(response.device_token);
                }
              },
              onError: (error) => {
                console.error('Failed to generate device token:', error);
                // Don't prevent booking, but notify user
              }
            }
          );
        })
        .catch(error => {
          console.error('Failed to generate device fingerprint:', error);
          // Fallback to a random device ID if fingerprinting fails
          const fallbackDeviceId = `web-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

          generateDeviceToken.mutate(
            { deviceId: fallbackDeviceId },
            {
              onSuccess: (response) => {
                if (response.device_token) {
                  setDeviceToken(response.device_token);
                }
              },
              onError: (error) => {
                console.error('Failed to generate fallback device token:', error);
              }
            }
          );
        });
    } else if (appointmentDetails.data && appointmentDetails.data.antiScalpingLevel !== 'strict') {
      // Clear device token if anti-scalping level is not strict
      setDeviceToken(null);
    }
  }, [appointmentDetails.data]);

  const goToCalendarStep = () => {
    if (!isValidCode) return;
    setLastSuccessfulStep(1);
    setDate('');
    setSlot(null);
    setRegisteredCount(1);
    setStep(2);
  };

  const onSubmitGuest = (v: BookingFormValues) => {
    if (!slot) return;

    // Check if appointment requires device token but it's not available
    if (appointmentDetails.data?.antiScalpingLevel === 'strict' && !deviceToken) {
      console.error('Device token is required but not available');
      return;
    }

    bookGuest.mutate(
      {
        appCode,
        date: normalizedDateIso,
        startTime: slot.startTime!,
        endTime: slot.endTime!,
        attendeeCount: v.attendeeCount,
        name: v.name,
        email: v.email,
        phone: v.phone,
        description: v.description,
        deviceToken: deviceToken || undefined,
      },
      {
        onSuccess: () => {
          setLastSuccessfulStep(4);
          setShowBurst(true);
          setTimeout(() => setShowBurst(false), 700);
        },
        onError: () => setStep(lastSuccessfulStep),
      }
    );
  };

  const onSubmitRegistered = () => {
    if (!slot) return;

    // Check if appointment requires device token but it's not available
    if (appointmentDetails.data?.antiScalpingLevel === 'strict' && !deviceToken) {
      console.error('Device token is required but not available');
      return;
    }

    bookReg.mutate(
      {
        appCode,
        date: normalizedDateIso,
        startTime: slot.startTime!,
        endTime: slot.endTime!,
        attendeeCount: registeredCount,
        deviceToken: deviceToken || undefined,
      },
      {
        onSuccess: () => {
          setLastSuccessfulStep(4);
          setShowBurst(true);
          setTimeout(() => setShowBurst(false), 700);
        },
        onError: () => setStep(lastSuccessfulStep),
      }
    );
  };

  return (
    <div style={{ maxWidth: 960, margin: '24px auto', display: 'grid', gap: 16, padding: '0 12px' }}>
      <Card>
        <CardHeader>
          <CardTitle>Book an appointment</CardTitle>
        </CardHeader>

        <ol className="stepper" aria-label="Booking steps">
          <li className={step >= 1 ? 'stepper__item is-active' : 'stepper__item'}>
            <span className="stepper__dot" aria-hidden="true" />
            <span className="stepper__label">Code</span>
          </li>
          <li className={step >= 2 ? 'stepper__item is-active' : 'stepper__item'}>
            <span className="stepper__dot" aria-hidden="true" />
            <span className="stepper__label">Day</span>
          </li>
          <li className={step >= 3 ? 'stepper__item is-active' : 'stepper__item'}>
            <span className="stepper__dot" aria-hidden="true" />
            <span className="stepper__label">Time</span>
          </li>
          <li className={step >= 4 ? 'stepper__item is-active' : 'stepper__item'}>
            <span className="stepper__dot" aria-hidden="true" />
            <span className="stepper__label">Review</span>
          </li>
        </ol>

        <div className="booking-step-grid">
          <div className="booking-step-card">
            {step === 1 && (
              <div style={{ display: 'grid', gap: 12 }}>
                <Field>
                  <FieldLabel>Appointment code</FieldLabel>
                  <small style={{ color: 'var(--text-muted)' }}>Enter the code shared by the host.</small>
                  <FieldRow>
                    <div style={{ position: 'relative', flex: 1 }}>
                      <IconSlot><i className="pi pi-hashtag" aria-hidden="true" /></IconSlot>
                      <Input value={appCode} onChange={(e) => setAppCode(e.target.value)} placeholder="AP-XXXXX" style={{ paddingLeft: 36 }} />
                    </div>
                    <Button
                      variant="primary"
                      disabled={!isValidCode || appointmentDetails.isFetching}
                      onClick={goToCalendarStep}
                      size="lg"
                    >
                      {appointmentDetails.isFetching ? 'Loading...' : 'Continue'}
                    </Button>
                  </FieldRow>
                </Field>
                {(appointmentDetails.isFetching || generateDeviceToken.isPending) && (
                  <div style={{ textAlign: 'center', padding: '20px 0' }}>
                    <div>Checking appointment details...</div>
                  </div>
                )}
              </div>
            )}

            {step === 2 && (
              <div style={{ display: 'grid', gap: 12 }}>
                <Field>
                  <FieldLabel>Pick a day</FieldLabel>
                  <small style={{ color: 'var(--text-muted)' }}>Available days are highlighted. Tap to see times.</small>
                </Field>
                {allSlots.isFetching && <div>Loading availability...</div>}
                {allSlots.error && <div style={{ color: 'var(--danger)' }}>Failed to load availability.</div>}
                {(appointmentDetails.isFetching || generateDeviceToken.isPending) && (
                  <div style={{ textAlign: 'center', padding: '20px 0' }}>
                    <div>Preparing booking options...</div>
                  </div>
                )}
                {!allSlots.isFetching && !allSlots.error && !availableDates.length && (
                  <div style={{ color: 'var(--text-muted)' }}>No open days for this code.</div>
                )}
                {availableDates.length > 0 && !appointmentDetails.isFetching && !generateDeviceToken.isPending && (
                  <AvailabilityCalendar
                    availableDates={availableDates}
                    selectedDate={date}
                    onSelect={(selected) => {
                      setDate(selected);
                      setSlot(null);
                      setLastSuccessfulStep(2);
                      setStep(3);
                    }}
                  />
                )}
              </div>
            )}

            {step === 3 && (
              <div style={{ display: 'grid', gap: 12 }}>
                <Field>
                  <FieldLabel>Selected day</FieldLabel>
                  <div className="selected-day-card">
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <Badge tone="primary">Day</Badge>
                      <div>{date ? formatSelectedDay(date) : 'Choose a day'}</div>
                    </div>
                    <Button variant="ghost" onClick={() => setStep(2)}>Change</Button>
                  </div>
                </Field>
                <Field>
                  <FieldLabel>Time</FieldLabel>
                  <div className="selected-slot-card">
                    {slot ? formatSlotSummary(slot) : 'No time selected'}
                  </div>
                  {slots.isFetching && <div>Loading slots...</div>}
                  {slots.error && <div style={{ color: 'var(--danger)' }}>Failed to load slots.</div>}
                  {(appointmentDetails.isFetching || generateDeviceToken.isPending) && (
                    <div style={{ textAlign: 'center', padding: '20px 0' }}>
                      <div>Preparing booking options...</div>
                    </div>
                  )}
                  {slots.data && !appointmentDetails.isFetching && !generateDeviceToken.isPending && (
                    <SlotPicker
                      slots={slots.data.items || []}
                      selected={slot}
                      onSelect={(s) => {
                        setSlot(s);
                        setLastSuccessfulStep(3);
                      }}
                    />
                  )}
                </Field>
                <div>
                  <Button variant="primary" disabled={!slot} onClick={proceedToSummary} size="lg" fullWidth>Review details</Button>
                </div>
              </div>
            )}

            {step === 4 && (
              <div style={{ display: 'grid', gap: 16 }}>
                <BookingSummary appCode={appCode} date={date} startTime={slot?.startTime || ''} endTime={slot?.endTime || ''} attendeeCount={isAuthed ? registeredCount : undefined} />
                {isAuthed ? (
                  <div style={{ display: 'grid', gap: 12 }}>
                    <Field>
                      <FieldLabel>Attendees</FieldLabel>
                      <FieldRow>
                        <div style={{ position: 'relative' }}>
                          <IconSlot><i className="pi pi-users" aria-hidden="true" /></IconSlot>
                          <Input
                            type="number"
                            min={1}
                            max={spotsRemaining || 1}
                            value={String(registeredCount)}
                            onChange={(e) => setRegisteredCount(Math.max(1, Math.min(Number(e.target.value) || 1, spotsRemaining || 1)))}
                            style={{ paddingLeft: 36 }}
                          />
                        </div>
                      </FieldRow>
                      <small style={{ color: 'var(--text-muted)' }}>{spotsRemaining === 1 ? '1 spot left' : `${spotsRemaining} spots left`}</small>
                    </Field>
                    <Button variant="primary" onClick={onSubmitRegistered} disabled={bookReg.isPending || spotsRemaining < 1} size="lg" fullWidth>Confirm booking</Button>
                  </div>
                ) : (
                  <BookingForm onSubmit={onSubmitGuest} pending={bookGuest.isPending} maxAttendees={spotsRemaining || undefined} />
                )}
              </div>
            )}
          </div>
        </div>
      </Card>
      <SuccessBurst show={showBurst} />
      {!isAuthed && (
        <Card>
          <CardHeader>
            <CardTitle>Want to create appointments?</CardTitle>
          </CardHeader>
          <div style={{ display: 'grid', gap: 12 }}>
            <p style={{ margin: 0, color: 'var(--text)' }}>
              Register an account to set up your own appointment codes, manage bookings, and access analytics.
            </p>
            <Button variant="primary" onClick={() => navigate('/signup')}>
              Create an account
            </Button>
          </div>
        </Card>
      )}
    </div>
  );
}

function formatSelectedDay(dateStr: string) {
  try {
    const d = new Date(`${dateStr}T00:00:00`);
    return format(d, 'EEE, MMM d, yyyy');
  } catch {
    return dateStr;
  }
}

function formatSlotSummary(slot?: EntitiesBooking | null) {
  if (!slot) return 'No time selected';
  const start = new Date(slot.startTime as string);
  const end = new Date(slot.endTime as string);
  return `${format(start, 'p')} – ${format(end, 'p')}`;
}


