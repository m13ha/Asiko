import type { EntitiesBooking } from '@appointment-master/api-client';
import { Button } from '@/components/Button';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';

function remainingSpots(slot: EntitiesBooking): number {
  const capacity = slot.capacity ?? slot.attendeeCount ?? 1;
  const booked = slot.seatsBooked ?? 0;
  const remaining = capacity - booked;
  return remaining > 0 ? remaining : 0;
}

export function SlotPicker({ slots, selected, onSelect }: { slots: EntitiesBooking[]; selected?: EntitiesBooking | null; onSelect: (s: EntitiesBooking) => void; }) {
  const available = (slots || []).filter((s) => s.available !== false && remainingSpots(s) > 0);
  if (!available.length) {
    return (
      <EmptyState>
        <EmptyTitle>No slots available</EmptyTitle>
        <EmptyDescription>Try another date or check back later.</EmptyDescription>
      </EmptyState>
    );
  }

  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(180px, 1fr))', gap: 8 }}>
      {available.map((slot) => {
        const key = `${slot.id ?? slot.startTime}-${slot.startTime}-${slot.endTime}`;
        const isSel = selected?.id ? selected.id === slot.id : (selected?.startTime === slot.startTime && selected?.endTime === slot.endTime);
        const spots = remainingSpots(slot);
        const label = `${slot.startTime} - ${slot.endTime}`;
        const badge = spots === 1 ? '1 spot left' : `${spots} spots left`;
        return (
          <Button key={key} variant={isSel ? 'primary' : undefined} onClick={() => onSelect(slot)} disabled={spots <= 0}>
            <span style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start', lineHeight: 1.2 }}>
              <span>{label}</span>
              <small style={{ opacity: 0.7 }}>{badge}</small>
            </span>
          </Button>
        );
      })}
    </div>
  );
}
