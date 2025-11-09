import { useEffect, useMemo, useState } from 'react';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { useAvailableSlots, useAvailableSlotsByDay, useBookGuest, useBookRegistered } from '../hooks';
import { SlotPicker } from '../components/SlotPicker';
import { BookingForm, BookingFormValues } from '../components/BookingForm';
import { BookingSummary } from '../components/BookingSummary';
import { useAuth } from '@/features/auth/AuthProvider';
import { useNavigate } from 'react-router-dom';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { Hash, Calendar, Users } from 'lucide-react';
import { SuccessBurst } from '@/components/SuccessBurst';
import type { EntitiesBooking } from '@appointment-master/api-client';
import { AvailabilityCalendar } from '../components/AvailabilityCalendar';

export function BookByCodePage() {
  const [step, setStep] = useState<1 | 2 | 3 | 4>(1);
  const [lastSuccessfulStep, setLastSuccessfulStep] = useState<1 | 2 | 3 | 4>(1);
  const [appCode, setAppCode] = useState('');
  const [date, setDate] = useState('');
  const [slot, setSlot] = useState<EntitiesBooking | null>(null);
  const [registeredCount, setRegisteredCount] = useState(1);
  const isValidCode = appCode.trim().length > 0;
  const { isAuthed } = useAuth();
  const navigate = useNavigate();

  const allSlots = useAvailableSlots(isValidCode ? appCode : '');
  const slots = useAvailableSlotsByDay(isValidCode ? appCode : '', date);
  const bookGuest = useBookGuest();
  const bookReg = useBookRegistered();
  const [showBurst, setShowBurst] = useState(false);

  const proceedToSummary = () => setStep(4);

  const spotsRemaining = slot ? Math.max((slot.capacity ?? slot.attendeeCount ?? 1) - (slot.seatsBooked ?? 0), 0) : 0;
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
      if (item.date) {
        const justDate = item.date.split('T')[0];
        unique.add(justDate);
      }
    }
    return Array.from(unique).sort();
  }, [allSlots.data?.items]);

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
    bookGuest.mutate(
      { appCode, date, startTime: slot.startTime!, endTime: slot.endTime!, attendeeCount: v.attendeeCount, name: v.name, email: v.email, phone: v.phone, description: v.description },
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
    bookReg.mutate(
      { appCode, date, startTime: slot.startTime!, endTime: slot.endTime!, attendeeCount: registeredCount },
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
    <div style={{ maxWidth: 720, margin: '24px auto', display: 'grid', gap: 16 }}>
      <Card>
        <CardHeader>
          <CardTitle>Book by Code</CardTitle>
        </CardHeader>

        {step === 1 && (
          <div style={{ display: 'grid', gap: 12 }}>
            <Field>
              <FieldLabel>Enter appointment code</FieldLabel>
              <FieldRow>
                <div style={{ position: 'relative' }}>
                  <IconSlot><Hash size={16} /></IconSlot>
                  <Input value={appCode} onChange={(e) => setAppCode(e.target.value)} placeholder="AP-XXXXX" style={{ paddingLeft: 36 }} />
                </div>
                <Button variant="primary" disabled={!isValidCode} onClick={goToCalendarStep}>Find slots</Button>
              </FieldRow>
            </Field>
          </div>
        )}

        {step === 2 && (
          <div style={{ display: 'grid', gap: 12 }}>
            <Field>
              <FieldLabel>Pick a day</FieldLabel>
              <small style={{ color: 'var(--text-muted)' }}>Available days are highlighted. Select one to view open times.</small>
            </Field>
            {allSlots.isFetching && <div>Loading availability...</div>}
            {allSlots.error && <div style={{ color: 'var(--danger)' }}>Failed to load availability.</div>}
            {!allSlots.isFetching && !allSlots.error && !availableDates.length && (
              <div style={{ color: 'var(--text-muted)' }}>No open days found for this appointment code.</div>
            )}
            {availableDates.length > 0 && (
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
              <FieldLabel>
                Available slots <small style={{ color: 'var(--text-muted)' }}>({date || 'Select a day'})</small>
              </FieldLabel>
              <FieldRow>
                <div style={{ position: 'relative' }}>
                  <IconSlot><Calendar size={16} /></IconSlot>
                  <Input
                    type="date"
                    value={date}
                    onChange={(e) => {
                      setDate(e.target.value);
                      setSlot(null);
                      setLastSuccessfulStep(2);
                      setStep(3);
                    }}
                    style={{ paddingLeft: 36 }}
                  />
                </div>
                <Button variant="ghost" onClick={() => setStep(2)}>Change day</Button>
              </FieldRow>
            </Field>
            {slots.isFetching && <div>Loading slots...</div>}
            {slots.error && <div style={{ color: 'var(--danger)' }}>Failed to load slots.</div>}
            {slots.data && (
              <SlotPicker
                slots={slots.data.items || []}
                selected={slot}
                onSelect={(s) => {
                  setSlot(s);
                  setLastSuccessfulStep(3);
                }}
              />
            )}
            <div>
              <Button variant="primary" disabled={!slot} onClick={proceedToSummary}>Continue</Button>
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
                      <IconSlot><Users size={16} /></IconSlot>
                      <Input
                        type="number"
                        min={1}
                        max={spotsRemaining || 1}
                        value={registeredCount}
                        onChange={(e) => setRegisteredCount(Math.max(1, Math.min(Number(e.target.value) || 1, spotsRemaining || 1)))}
                        style={{ paddingLeft: 36 }}
                      />
                    </div>
                  </FieldRow>
                  <small style={{ color: 'var(--text-muted)' }}>{spotsRemaining === 1 ? '1 spot left' : `${spotsRemaining} spots left`}</small>
                </Field>
                <Button variant="primary" onClick={onSubmitRegistered} disabled={bookReg.isPending || spotsRemaining < 1}>Confirm booking</Button>
              </div>
            ) : (
              <BookingForm onSubmit={onSubmitGuest} pending={bookGuest.isPending} maxAttendees={spotsRemaining || undefined} />
            )}
          </div>
        )}
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
