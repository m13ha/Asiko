-- Create a new enum type for the anti-scalping levels
CREATE TYPE anti_scalping_level AS ENUM ('none', 'standard', 'strict');

-- Add the anti_scalping_level column to the appointments table
ALTER TABLE appointments ADD COLUMN anti_scalping_level anti_scalping_level NOT NULL DEFAULT 'none';

-- Add device_id column to the bookings table
ALTER TABLE bookings ADD COLUMN device_id TEXT;

-- Add an index for faster lookups on booking restrictions
CREATE INDEX IF NOT EXISTS idx_bookings_restriction_lookup ON bookings(appointment_id, email, device_id);
