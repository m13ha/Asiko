-- Revert all changes from the up migration
ALTER TABLE appointments DROP COLUMN anti_scalping_level;

ALTER TABLE bookings DROP COLUMN device_id;

DROP INDEX IF EXISTS idx_bookings_restriction_lookup;

DROP TYPE IF EXISTS anti_scalping_level;
