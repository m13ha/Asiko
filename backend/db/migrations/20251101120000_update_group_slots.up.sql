-- 20251101120000_update_group_slots.up.sql

ALTER TABLE bookings
    ADD COLUMN IF NOT EXISTS is_slot BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS capacity INTEGER,
    ADD COLUMN IF NOT EXISTS seats_booked INTEGER NOT NULL DEFAULT 0;

-- Backfill slot metadata for pre-generated availability rows
UPDATE bookings b
SET
    is_slot = true,
    capacity = COALESCE(a.max_attendees, 1),
    seats_booked = 0
FROM appointments a
WHERE b.appointment_id = a.id
  AND b.available = true;

-- Ensure existing reservation rows have sensible defaults
UPDATE bookings
SET
    is_slot = false,
    capacity = COALESCE(capacity, GREATEST(attendee_count, 1)),
    seats_booked = 0
WHERE available = false;

-- Enforce non-null capacity going forward
ALTER TABLE bookings
    ALTER COLUMN capacity SET DEFAULT 1,
    ALTER COLUMN capacity SET NOT NULL;
