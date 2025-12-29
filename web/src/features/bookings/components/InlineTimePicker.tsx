import { memo, useMemo, useCallback } from 'react';
import type { EntitiesBooking } from '@appointment-master/api-client';
import * as API from '@appointment-master/api-client';
import { Badge } from '@/components/Badge';

export interface InlineTimePickerProps {
  slots: EntitiesBooking[];
  selectedSlot: EntitiesBooking | null;
  onSelect: (slot: EntitiesBooking) => void;
  className?: string;
  appointmentType?: API.EntitiesAppointmentType;
}

/** Calculate remaining spots for a slot */
function remainingSpots(slot: EntitiesBooking): number {
  const capacity = slot.capacity ?? slot.attendeeCount ?? 1;
  const booked = slot.seatsBooked ?? 0;
  return Math.max(capacity - booked, 0);
}

/** Format time for display */
function formatTime(isoTime: string): string {
  try {
    const date = new Date(isoTime);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  } catch {
    return isoTime;
  }
}

/** Generate unique key for slot */
function slotKey(slot: EntitiesBooking): string {
  return `${slot.id ?? 'slot'}-${slot.startTime}-${slot.endTime}`;
}

interface TimeSlotButtonProps {
  slot: EntitiesBooking;
  isSelected: boolean;
  spots: number;
  appointmentType?: API.EntitiesAppointmentType;
  onSelect: (slot: EntitiesBooking) => void;
}

/** Individual time slot button - memoized for performance */
const TimeSlotButton = memo(function TimeSlotButton({
  slot,
  isSelected,
  spots,
  appointmentType,
  onSelect,
}: TimeSlotButtonProps) {
  const timeLabel = useMemo(() => formatTime(slot.startTime!), [slot.startTime]);
  const capacity = slot.capacity ?? slot.attendeeCount ?? spots;
  const isLowSpots = appointmentType !== API.EntitiesAppointmentType.Single && capacity > 1 && spots <= 2;

  const handleClick = useCallback(() => {
    onSelect(slot);
  }, [onSelect, slot]);

  return (
    <button
      type="button"
      onClick={handleClick}
      className={[
        'w-full px-4 py-3 rounded-xl border text-left transition-all duration-200',
        'hover:border-[color-mix(in_oklab,var(--primary)_50%,var(--border))]',
        'hover:shadow-[0_0_0_3px_color-mix(in_oklab,var(--primary)_12%,transparent)]',
        'focus:outline-none focus:ring-2 focus:ring-[var(--primary)] focus:ring-offset-2 focus:ring-offset-[var(--bg)]',
        isSelected
          ? 'border-[var(--primary)] bg-[color-mix(in_oklab,var(--primary)_12%,var(--bg-elevated))] shadow-[0_0_0_3px_color-mix(in_oklab,var(--primary)_20%,transparent)]'
          : 'border-[var(--border)] bg-[var(--bg-elevated)]',
      ].join(' ')}
      aria-pressed={isSelected}
    >
      <div className="flex items-center gap-3">
        <div>
          <div className="font-semibold text-[var(--text)]">{timeLabel}</div>
          <div className="text-xs text-[var(--text-muted)]">
            {spots === 1 ? '1 spot' : `${spots} spots`}
          </div>
        </div>
        <Badge tone={isLowSpots ? 'danger' : 'success'}>
          {isLowSpots ? 'Filling' : 'Open'}
        </Badge>
      </div>
    </button>
  );
});

/**
 * Horizontal scrolling time slot picker.
 * Displays available time slots as selectable pills/buttons.
 */
export const InlineTimePicker = memo(function InlineTimePicker({
  slots,
  selectedSlot,
  onSelect,
  className = '',
  appointmentType,
}: InlineTimePickerProps) {
  // Filter to only available slots with remaining capacity
  const availableSlots = useMemo(
    () => (slots || []).filter((s) => s.available !== false && remainingSpots(s) > 0),
    [slots]
  );

  const selectedKey = selectedSlot ? slotKey(selectedSlot) : null;

  if (!availableSlots.length) {
    return (
      <div className={`text-center py-6 text-[var(--text-muted)] ${className}`}>
        <p className="m-0 text-sm">No time slots available for this day.</p>
        <p className="m-0 text-xs mt-1">Try selecting a different date.</p>
      </div>
    );
  }

  return (
    <div className={`space-y-3 ${className}`}>
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-[var(--text)] m-0">Available Times</h3>
        <span className="text-xs text-[var(--text-muted)]">
          {availableSlots.length} slot{availableSlots.length !== 1 ? 's' : ''}
        </span>
      </div>
      
      {/* Grid container */}
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
        {availableSlots.map((slot) => (
          <TimeSlotButton
            key={slotKey(slot)}
            slot={slot}
            isSelected={slotKey(slot) === selectedKey}
            spots={remainingSpots(slot)}
            appointmentType={appointmentType}
            onSelect={onSelect}
          />
        ))}
      </div>
    </div>
  );
});

export default InlineTimePicker;
