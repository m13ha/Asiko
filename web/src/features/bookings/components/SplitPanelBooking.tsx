import { memo, useMemo, useState, useCallback, useEffect } from 'react';
import type { EntitiesBooking } from '@appointment-master/api-client';
import { format, parseISO } from 'date-fns';
import { Calendar } from 'primereact/calendar';
import { Button } from '@/components/Button';
import { Card } from '@/components/Card';
import { AppointmentDetailsPanel, type AppointmentInfo } from './AppointmentDetailsPanel';
import { InlineTimePicker } from './InlineTimePicker';
import { ChevronLeft, ChevronRight, Calendar as CalendarIcon } from 'lucide-react';

// ============================================================================
// Types
// ============================================================================

export interface SplitPanelBookingProps {
  /** Appointment details for the left panel */
  appointment: AppointmentInfo;
  /** Available dates in ISO format (yyyy-MM-dd) */
  availableDates: string[];
  /** Time slots for the selected date */
  slots: EntitiesBooking[];
  /** Currently selected date */
  selectedDate: string | null;
  /** Currently selected time slot */
  selectedSlot: EntitiesBooking | null;
  /** Loading state for slots */
  isLoadingSlots?: boolean;
  /** Callback when date is selected */
  onDateSelect: (date: string) => void;
  /** Callback when slot is selected */
  onSlotSelect: (slot: EntitiesBooking) => void;
  /** Callback to proceed to booking details */
  onContinue: () => void;
  /** Optional class name */
  className?: string;
}

// ============================================================================
// Sub-components
// ============================================================================

interface DateTemplateEvent {
  day: number;
  month: number;
  year: number;
  otherMonth: boolean;
  today: boolean;
  selectable: boolean;
}

interface CalendarSectionProps {
  availableDates: string[];
  selectedDate: string | null;
  onDateSelect: (date: string) => void;
}

/** Calendar section with date selection - memoized */
const CalendarSection = memo(function CalendarSection({
  availableDates,
  selectedDate,
  onDateSelect,
}: CalendarSectionProps) {
  const availableSet = useMemo(() => new Set(availableDates), [availableDates]);
  
  const selected = useMemo(
    () => (selectedDate ? parseISO(selectedDate) : null),
    [selectedDate]
  );
  
  const defaultMonth = useMemo(() => {
    if (selected) return selected;
    if (availableDates.length > 0) return parseISO(availableDates[0]);
    return new Date();
  }, [selected, availableDates]);

  const handleChange = useCallback(
    (event: { value?: Date | Date[] | null }) => {
      const value = (Array.isArray(event.value) ? event.value[0] : event.value) as Date | null;
      if (!value) return;
      const iso = format(value, 'yyyy-MM-dd');
      if (availableSet.has(iso)) {
        onDateSelect(iso);
      }
    },
    [availableSet, onDateSelect]
  );

  const dateTemplate = useCallback(
    (event: DateTemplateEvent) => {
      const current = new Date(event.year, event.month, event.day);
      const iso = format(current, 'yyyy-MM-dd');
      const isAvailable = availableSet.has(iso);
      const isSelected = selected ? iso === format(selected, 'yyyy-MM-dd') : false;

      const baseClasses =
        'inline-flex items-center justify-center w-9 h-9 rounded-lg text-sm font-medium transition-all duration-200';
      const availabilityClasses = isAvailable
        ? 'text-[var(--text)] cursor-pointer hover:bg-[color-mix(in_oklab,var(--primary)_15%,transparent)]'
        : 'text-[var(--text-muted)] opacity-40 line-through cursor-not-allowed';
      const selectedClasses = isSelected
        ? 'bg-[var(--primary)] text-[var(--primary-contrast)] shadow-[0_0_0_3px_color-mix(in_oklab,var(--primary)_30%,transparent)]'
        : '';
      const otherMonthClasses = event.otherMonth ? 'opacity-30' : '';

      return (
        <span className={`${baseClasses} ${availabilityClasses} ${selectedClasses} ${otherMonthClasses}`}>
          {event.day}
        </span>
      );
    },
    [availableSet, selected]
  );

  return (
    <div className="rounded-2xl border border-[var(--border)] bg-[var(--bg-elevated)] p-4">
      <div className="flex items-center gap-2 mb-3">
        <CalendarIcon className="w-4 h-4 text-[var(--primary)]" />
        <h3 className="text-sm font-semibold text-[var(--text)] m-0">Select a Date</h3>
      </div>
      <Calendar
        value={selected ?? defaultMonth ?? null}
        onChange={handleChange}
        inline
        numberOfMonths={1}
        dateTemplate={dateTemplate}
        showWeek={false}
        className="w-full split-panel-calendar"
        prevIcon={<ChevronLeft className="w-4 h-4" />}
        nextIcon={<ChevronRight className="w-4 h-4" />}
      />
      {/* Custom styles for the calendar */}
      <style>{`
        .split-panel-calendar .p-datepicker {
          width: 100% !important;
          border: none !important;
          box-shadow: none !important;
          background: transparent !important;
          color: var(--text) !important;
        }
        .split-panel-calendar .p-datepicker-header {
          background: transparent !important;
          border-bottom: 1px solid var(--border) !important;
          padding-bottom: 0.75rem !important;
          margin-bottom: 0.75rem !important;
        }
        .split-panel-calendar .p-datepicker-title {
          font-weight: 600 !important;
          color: var(--text) !important;
        }
        .split-panel-calendar .p-datepicker-prev,
        .split-panel-calendar .p-datepicker-next {
          color: var(--text-muted) !important;
          transition: color 0.2s !important;
        }
        .split-panel-calendar .p-datepicker-prev:hover,
        .split-panel-calendar .p-datepicker-next:hover {
          color: var(--primary) !important;
          background: transparent !important;
        }
        .split-panel-calendar .p-datepicker table td {
          padding: 0.15rem !important;
        }
        .split-panel-calendar .p-datepicker table th {
          color: var(--text-muted) !important;
          font-weight: 500 !important;
          font-size: 0.75rem !important;
        }
      `}</style>
    </div>
  );
});

interface SelectedDateDisplayProps {
  selectedDate: string | null;
}

/** Display the currently selected date */
const SelectedDateDisplay = memo(function SelectedDateDisplay({
  selectedDate,
}: SelectedDateDisplayProps) {
  const formattedDate = useMemo(() => {
    if (!selectedDate) return null;
    try {
      const date = parseISO(selectedDate);
      return format(date, 'EEEE, MMMM d, yyyy');
    } catch {
      return selectedDate;
    }
  }, [selectedDate]);

  if (!formattedDate) return null;

  return (
    <div className="flex items-center gap-2 px-4 py-3 rounded-xl bg-[color-mix(in_oklab,var(--primary)_10%,var(--bg))] border border-[color-mix(in_oklab,var(--primary)_25%,var(--border))]">
      <CalendarIcon className="w-4 h-4 text-[var(--primary)]" />
      <span className="text-sm font-semibold text-[var(--text)]">{formattedDate}</span>
    </div>
  );
});

// ============================================================================
// Main Component
// ============================================================================

/**
 * Split Panel Booking Layout
 * 
 * A desktop-optimized booking interface with:
 * - Left panel: Appointment details (fixed)
 * - Right panel: Calendar + time slot picker
 * 
 * Collapses to single column on mobile.
 */
export const SplitPanelBooking = memo(function SplitPanelBooking({
  appointment,
  availableDates,
  slots,
  selectedDate,
  selectedSlot,
  isLoadingSlots = false,
  onDateSelect,
  onSlotSelect,
  onContinue,
  className = '',
}: SplitPanelBookingProps) {
  const hasSelection = Boolean(selectedDate && selectedSlot);

  return (
    <div className={`grid gap-6 lg:grid-cols-[320px_1fr] items-start ${className}`}>
      {/* Left Panel - Appointment Details */}
      <div className="lg:sticky lg:top-4">
        <AppointmentDetailsPanel appointment={appointment} />
      </div>

      {/* Right Panel - Date & Time Selection */}
      <div className="space-y-5">
        {/* Calendar */}
        <CalendarSection
          availableDates={availableDates}
          selectedDate={selectedDate}
          onDateSelect={onDateSelect}
        />

        {/* Selected Date Display */}
        {selectedDate && <SelectedDateDisplay selectedDate={selectedDate} />}

        {/* Time Slots */}
        {selectedDate && (
          <Card className="!p-4">
            {isLoadingSlots ? (
              <div className="flex items-center justify-center py-8 gap-3">
                <div className="w-5 h-5 border-2 border-[var(--primary)] border-t-transparent rounded-full animate-spin" />
                <span className="text-sm text-[var(--text-muted)]">Loading available times...</span>
              </div>
            ) : (
              <InlineTimePicker
                slots={slots}
                selectedSlot={selectedSlot}
                onSelect={onSlotSelect}
              />
            )}
          </Card>
        )}

        {/* Continue Button */}
        <div className="flex justify-end pt-2">
          <Button
            variant="primary"
            size="lg"
            disabled={!hasSelection}
            onClick={onContinue}
            className="min-w-[200px]"
          >
            {hasSelection ? 'Continue to Details' : 'Select date & time'}
          </Button>
        </div>
      </div>
    </div>
  );
});

export default SplitPanelBooking;
