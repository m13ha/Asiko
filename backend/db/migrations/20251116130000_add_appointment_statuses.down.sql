-- 000004_add_appointment_statuses.down.sql

ALTER TABLE appointments
    DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS appointment_status;
