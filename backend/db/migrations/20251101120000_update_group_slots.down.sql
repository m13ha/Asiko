-- 20251101120000_update_group_slots.down.sql

ALTER TABLE bookings
    DROP COLUMN IF EXISTS seats_booked,
    DROP COLUMN IF EXISTS capacity,
    DROP COLUMN IF EXISTS is_slot;
