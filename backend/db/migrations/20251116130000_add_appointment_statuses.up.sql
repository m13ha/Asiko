-- 000004_add_appointment_statuses.up.sql

CREATE TYPE appointment_status AS ENUM ('pending', 'ongoing', 'completed', 'canceled', 'expired');

ALTER TABLE appointments
    ADD COLUMN status appointment_status NOT NULL DEFAULT 'pending';

UPDATE appointments
SET status = 'pending'
WHERE status IS NULL;
