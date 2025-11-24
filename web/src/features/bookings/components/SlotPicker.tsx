import type { EntitiesBooking } from '@appointment-master/api-client';
import { useMemo, useState } from 'react';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';
import { Dialog } from 'primereact/dialog';
import { Button } from '@/components/Button';
import { Badge } from '@/components/Badge';

function remainingSpots(slot: EntitiesBooking): number {
  const capacity = slot.capacity ?? slot.attendeeCount ?? 1;
  const booked = slot.seatsBooked ?? 0;
  const remaining = capacity - booked;
  return remaining > 0 ? remaining : 0;
}

type SlotPickerProps = {
  slots: EntitiesBooking[];
  selected?: EntitiesBooking | null;
  onSelect: (slot: EntitiesBooking) => void;
};

type SlotOption = {
  label: string;
  value: string;
  slot: EntitiesBooking;
  spots: number;
};

function optionKey(slot: EntitiesBooking) {
  return `${slot.id ?? 'slot'}-${slot.startTime}-${slot.endTime}`;
}

export function SlotPicker({ slots, selected, onSelect }: SlotPickerProps) {
  const available = useMemo(() => (slots || []).filter((s) => s.available !== false && remainingSpots(s) > 0), [slots]);
  const [open, setOpen] = useState(false);

  if (!available.length) {
    return (
      <EmptyState>
        <EmptyTitle>No slots available</EmptyTitle>
        <EmptyDescription>Try another date or check back later.</EmptyDescription>
      </EmptyState>
    );
  }

  const options: SlotOption[] = available.map(slot => ({
    label: formatTimeRange(slot),
    value: optionKey(slot),
    slot,
    spots: remainingSpots(slot),
  }));

  const selectedKey = selected ? optionKey(selected) : null;

  const handleSelect = (slot: EntitiesBooking) => {
    onSelect(slot);
    setOpen(false);
  };

  return (
    <>
      <Button variant="secondary" onClick={() => setOpen(true)} size="lg" fullWidth>
        {selected ? `Change time • ${formatTimeRange(selected)}` : 'Select a time'}
      </Button>
      <Dialog
        header="Pick a time"
        visible={open}
        onHide={() => setOpen(false)}
        className="slot-dialog"
        contentClassName="slot-dialog-content"
      >
        <div className="slot-dialog-list">
          {options.map(option => {
            const active = option.value === selectedKey;
            const lowSpots = option.spots <= 2;
            return (
              <button
                key={option.value}
                type="button"
                className={['slot-dialog-item', active ? 'is-active' : ''].filter(Boolean).join(' ')}
                onClick={() => handleSelect(option.slot)}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 12, justifyContent: 'space-between', width: '100%' }}>
                  <div style={{ textAlign: 'left' }}>
                    <div style={{ fontWeight: 600 }}>{option.label}</div>
                    <small style={{ color: 'var(--text-muted)' }}>{option.spots === 1 ? '1 spot left' : `${option.spots} spots left`}</small>
                  </div>
                  <Badge tone={lowSpots ? 'danger' : 'primary'}>{lowSpots ? 'Filling' : 'Available'}</Badge>
                </div>
              </button>
            );
          })}
        </div>
      </Dialog>
    </>
  );
}

function formatTimeRange(slot: EntitiesBooking) {
  const start = new Date(slot.startTime as string);
  const end = new Date(slot.endTime as string);
  const startLabel = start.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  const endLabel = end.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  return `${startLabel} – ${endLabel}`;
}
