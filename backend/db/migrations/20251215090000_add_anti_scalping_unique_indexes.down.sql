-- 20251215090000_add_anti_scalping_unique_indexes.down.sql

DROP INDEX IF EXISTS uniq_bookings_active_email;
DROP INDEX IF EXISTS uniq_bookings_active_device;
DROP INDEX IF EXISTS uniq_bookings_active_phone;
