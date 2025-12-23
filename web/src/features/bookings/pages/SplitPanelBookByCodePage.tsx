import { useEffect, useMemo, useState, useCallback } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { format, parseISO } from 'date-fns';
import FingerprintJS from '@sparkstone/fingerprintjs';
import type { EntitiesBooking } from '@appointment-master/api-client';

import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Field, FieldLabel, FieldRow, IconSlot } from '@/components/Field';
import { SuccessBurst } from '@/components/SuccessBurst';

import { useIsAuthed } from '@/stores/authStore';
import { useDeviceToken } from '@/features/auth/hooks';
import { useAppointmentByAppCode } from '@/features/appointments/hooks';
import {
  useAvailableSlotsByDay,
  useAvailableDates,
  useBookGuest,
  useBookRegistered,
} from '../hooks';

import { SplitPanelBooking } from '../components/SplitPanelBooking';
import { BookingForm, type BookingFormValues } from '../components/BookingForm';
import { BookingSummary } from '../components/BookingSummary';

// ============================================================================
// Types & Helpers
// ============================================================================

type BookingStep = 'code' | 'select' | 'details';

function normalizeDateOnly(value?: string | null): string | null {
  if (!value) return null;
  try {
    const isoDateString = value.trim().replace(' ', 'T');
    const date = parseISO(isoDateString.endsWith('Z') ? isoDateString : `${isoDateString}Z`);
    if (Number.isNaN(date.getTime())) return null;
    return format(date, 'yyyy-MM-dd');
  } catch {
    return null;
  }
}

function formatSlotTime(slot?: EntitiesBooking | null): string {
  if (!slot?.startTime || !slot?.endTime) return '';
  try {
    const start = new Date(slot.startTime);
    const end = new Date(slot.endTime);
    return `${format(start, 'p')} – ${format(end, 'p')}`;
  } catch {
    return '';
  }
}

// ============================================================================
// Main Component
// ============================================================================

/**
 * Split Panel Booking Page
 * 
 * A redesigned booking flow using the split panel calendar design:
 * 1. Code Entry - Enter appointment code
 * 2. Select - Split panel with calendar + time picker
 * 3. Details - Booking form / confirmation
 */
export function SplitPanelBookByCodePage() {
  const [search] = useSearchParams();
  const navigate = useNavigate();
  const isAuthed = useIsAuthed();

  // State
  const [step, setStep] = useState<BookingStep>('code');
  const [appCode, setAppCode] = useState(search.get('code') || '');
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedSlot, setSelectedSlot] = useState<EntitiesBooking | null>(null);
  const [registeredCount, setRegisteredCount] = useState(1);
  const [deviceToken, setDeviceToken] = useState<string | null>(null);
  const [showBurst, setShowBurst] = useState(false);

  // Queries
  const isValidCode = appCode.trim().length > 0;
  const appointmentDetails = useAppointmentByAppCode(appCode);
  const availableDatesQuery = useAvailableDates(isValidCode ? appCode : '');
  const daySlots = useAvailableSlotsByDay(isValidCode ? appCode : '', selectedDate || '');

  // Mutations
  const bookGuest = useBookGuest();
  const bookReg = useBookRegistered();
  const generateDeviceToken = useDeviceToken();

  // Derived data
  const availableDates = useMemo(() => {
    return availableDatesQuery.data || [];
  }, [availableDatesQuery.data]);

  const spotsRemaining = useMemo(() => {
    if (!selectedSlot) return 0;
    const capacity = selectedSlot.capacity ?? selectedSlot.attendeeCount ?? 1;
    const booked = selectedSlot.seatsBooked ?? 0;
    return Math.max(capacity - booked, 0);
  }, [selectedSlot]);

  const appointmentInfo = useMemo(() => ({
    title: appointmentDetails.data?.title,
    description: appointmentDetails.data?.description,
    hostName: appointmentDetails.data?.ownerName,
    duration: appointmentDetails.data?.bookingDuration,
    maxAttendees: appointmentDetails.data?.maxAttendees,
    type: appointmentDetails.data?.type,
    startDate: appointmentDetails.data?.startDate,
    endDate: appointmentDetails.data?.endDate,
  }), [appointmentDetails.data]);

  // Effects
  useEffect(() => {
    if (search.get('code')) {
      setAppCode(search.get('code')!);
      setStep('select');
    }
  }, [search]);

  useEffect(() => {
    setSelectedSlot(null);
    setRegisteredCount(1);
  }, [selectedDate]);

  useEffect(() => {
    if (selectedSlot && spotsRemaining > 0) {
      setRegisteredCount((prev) => Math.min(prev, spotsRemaining));
    }
  }, [selectedSlot, spotsRemaining]);

  // Device token generation for strict anti-scalping
  useEffect(() => {
    if (appointmentDetails.data?.antiScalpingLevel === 'strict') {
      FingerprintJS.load()
        .then((fp) => fp.get())
        .then((result) => {
          generateDeviceToken.mutate(
            { deviceId: result.visitorId },
            { onSuccess: (res) => res.device_token && setDeviceToken(res.device_token) }
          );
        })
        .catch(() => {
          const fallback = `web-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
          generateDeviceToken.mutate(
            { deviceId: fallback },
            { onSuccess: (res) => res.device_token && setDeviceToken(res.device_token) }
          );
        });
    } else if (appointmentDetails.data) {
      setDeviceToken(null);
    }
  }, [appointmentDetails.data]);

  // Handlers
  const handleCodeSubmit = useCallback(() => {
    if (!isValidCode) return;
    setSelectedDate(null);
    setSelectedSlot(null);
    setStep('select');
  }, [isValidCode]);

  const handleDateSelect = useCallback((date: string) => {
    setSelectedDate(date);
    setSelectedSlot(null);
  }, []);

  const handleSlotSelect = useCallback((slot: EntitiesBooking) => {
    setSelectedSlot(slot);
  }, []);

  const handleContinueToDetails = useCallback(() => {
    if (selectedDate && selectedSlot) {
      setStep('details');
    }
  }, [selectedDate, selectedSlot]);

  const handleBackToSelect = useCallback(() => {
    setStep('select');
  }, []);

  const normalizedDateIso = selectedDate ? `${selectedDate}T00:00:00Z` : '';

  const handleGuestSubmit = useCallback(
    (v: BookingFormValues) => {
      if (!selectedSlot) return;
      if (appointmentDetails.data?.antiScalpingLevel === 'strict' && !deviceToken) return;

      bookGuest.mutate(
        {
          appCode,
          date: normalizedDateIso,
          startTime: selectedSlot.startTime!,
          endTime: selectedSlot.endTime!,
          attendeeCount: v.attendeeCount,
          name: v.name,
          email: v.email,
          phone: v.phone,
          description: v.description,
          deviceToken: deviceToken || undefined,
        },
        {
          onSuccess: () => {
            setShowBurst(true);
            setTimeout(() => setShowBurst(false), 700);
          },
        }
      );
    },
    [appCode, normalizedDateIso, selectedSlot, deviceToken, appointmentDetails.data, bookGuest]
  );

  const handleRegisteredSubmit = useCallback(() => {
    if (!selectedSlot) return;
    if (appointmentDetails.data?.antiScalpingLevel === 'strict' && !deviceToken) return;

    bookReg.mutate(
      {
        appCode,
        date: normalizedDateIso,
        startTime: selectedSlot.startTime!,
        endTime: selectedSlot.endTime!,
        attendeeCount: registeredCount,
        deviceToken: deviceToken || undefined,
      },
      {
        onSuccess: () => {
          setShowBurst(true);
          setTimeout(() => setShowBurst(false), 700);
        },
      }
    );
  }, [appCode, normalizedDateIso, selectedSlot, registeredCount, deviceToken, appointmentDetails.data, bookReg]);

  // Render
  return (
    <div className="max-w-5xl mx-auto my-6 px-4">
      {/* Step 1: Code Entry */}
      {step === 'code' && (
        <Card className="max-w-lg mx-auto">
          <CardHeader>
            <CardTitle>Book an Appointment</CardTitle>
          </CardHeader>
          <div className="space-y-4">
            <Field>
              <FieldLabel>Appointment Code</FieldLabel>
              <p className="text-sm text-[var(--text-muted)] mb-2">
                Enter the code shared by your host to view available times.
              </p>
              <FieldRow>
                <div className="relative flex-1">
                  <IconSlot>
                    <i className="pi pi-hashtag" aria-hidden="true" />
                  </IconSlot>
                  <Input
                    value={appCode}
                    onChange={(e) => setAppCode(e.target.value)}
                    placeholder="AP-XXXXX"
                    className="pl-9"
                    onKeyDown={(e) => e.key === 'Enter' && handleCodeSubmit()}
                  />
                </div>
              </FieldRow>
            </Field>
            <Button
              variant="primary"
              size="lg"
              disabled={!isValidCode || appointmentDetails.isFetching}
              onClick={handleCodeSubmit}
              className="w-full"
            >
              {appointmentDetails.isFetching ? 'Loading...' : 'View Available Times'}
            </Button>
          </div>
        </Card>
      )}

      {/* Step 2: Split Panel Selection */}
      {step === 'select' && (
        <div className="space-y-4">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="sm" onClick={() => setStep('code')}>
              <i className="pi pi-arrow-left mr-2" /> Back
            </Button>
            <span className="text-sm text-[var(--text-muted)]">Code: {appCode}</span>
          </div>

          {availableDatesQuery.isLoading || appointmentDetails.isFetching ? (
            <Card className="py-12 text-center">
              <div className="flex items-center justify-center gap-3">
                <div className="w-5 h-5 border-2 border-[var(--primary)] border-t-transparent rounded-full animate-spin" />
                <span className="text-[var(--text-muted)]">Loading appointment details...</span>
              </div>
            </Card>
          ) : availableDatesQuery.error ? (
            <Card className="py-8 text-center">
              <p className="text-[var(--danger)]">Failed to load availability. Please try again.</p>
              <Button variant="outline" className="mt-4" onClick={() => setStep('code')}>
                Try Different Code
              </Button>
            </Card>
          ) : (
            <SplitPanelBooking
              appointment={appointmentInfo}
              availableDates={availableDates}
              slots={daySlots.data?.items || []}
              selectedDate={selectedDate}
              selectedSlot={selectedSlot}
              isLoadingSlots={daySlots.isFetching}
              onDateSelect={handleDateSelect}
              onSlotSelect={handleSlotSelect}
              onContinue={handleContinueToDetails}
            />
          )}
        </div>
      )}

      {/* Step 3: Booking Details */}
      {step === 'details' && (
        <div className="max-w-2xl mx-auto space-y-4">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="sm" onClick={handleBackToSelect}>
              <i className="pi pi-arrow-left mr-2" /> Back
            </Button>
            <span className="text-sm text-[var(--text-muted)]">Change date or time</span>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Complete Your Booking</CardTitle>
            </CardHeader>

            <BookingSummary
              appCode={appCode}
              date={selectedDate || ''}
              startTime={selectedSlot?.startTime || ''}
              endTime={selectedSlot?.endTime || ''}
              attendeeCount={isAuthed ? registeredCount : undefined}
            />

            <div className="mt-4">
              {isAuthed ? (
                <div className="space-y-4">
                  <Field>
                    <FieldLabel>Number of Attendees</FieldLabel>
                    <FieldRow>
                      <div className="relative max-w-[160px]">
                        <IconSlot>
                          <i className="pi pi-users" aria-hidden="true" />
                        </IconSlot>
                        <Input
                          type="number"
                          min={1}
                          max={spotsRemaining || 1}
                          value={String(registeredCount)}
                          onChange={(e) =>
                            setRegisteredCount(
                              Math.max(1, Math.min(Number(e.target.value) || 1, spotsRemaining || 1))
                            )
                          }
                          className="pl-9"
                        />
                      </div>
                    </FieldRow>
                    <small className="text-xs text-[var(--text-muted)]">
                      {spotsRemaining === 1 ? '1 spot left' : `${spotsRemaining} spots left`}
                    </small>
                  </Field>
                  <Button
                    variant="primary"
                    size="lg"
                    className="w-full"
                    onClick={handleRegisteredSubmit}
                    disabled={bookReg.isPending || spotsRemaining < 1}
                  >
                    {bookReg.isPending ? 'Confirming...' : 'Confirm Booking'}
                  </Button>
                </div>
              ) : (
                <BookingForm
                  onSubmit={handleGuestSubmit}
                  pending={bookGuest.isPending}
                  maxAttendees={spotsRemaining || undefined}
                />
              )}
            </div>
          </Card>
        </div>
      )}

      {/* Success animation */}
      <SuccessBurst show={showBurst} />

      {/* Sign up prompt for guests */}
      {!isAuthed && step !== 'code' && (
        <Card className="max-w-lg mx-auto mt-6">
          <CardHeader>
            <CardTitle>Want to create appointments?</CardTitle>
          </CardHeader>
          <p className="text-sm text-[var(--text-muted)] mb-3">
            Register an account to set up your own appointment codes and manage bookings.
          </p>
          <Button variant="outline" onClick={() => navigate('/signup')}>
            Create an Account
          </Button>
        </Card>
      )}
    </div>
  );
}

export default SplitPanelBookByCodePage;
