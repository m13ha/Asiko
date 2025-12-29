-- 20251215090000_add_anti_scalping_unique_indexes.up.sql

CREATE UNIQUE INDEX IF NOT EXISTS uniq_bookings_active_email
    ON bookings (appointment_id, lower(email))
    WHERE email IS NOT NULL
      AND email <> ''
      AND status IN ('active', 'ongoing', 'pending', 'confirmed');

CREATE UNIQUE INDEX IF NOT EXISTS uniq_bookings_active_device
    ON bookings (appointment_id, device_id)
    WHERE device_id IS NOT NULL
      AND device_id <> ''
      AND status IN ('active', 'ongoing', 'pending', 'confirmed');

CREATE UNIQUE INDEX IF NOT EXISTS uniq_bookings_active_phone
    ON bookings (appointment_id, phone)
    WHERE phone IS NOT NULL
      AND phone <> ''
      AND status IN ('active', 'ongoing', 'pending', 'confirmed');
