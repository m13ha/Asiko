-- 000001_create_initial_tables.up.sql

CREATE TYPE appointment_type AS ENUM ('single', 'group', 'party');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    phone_number TEXT UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    booking_duration INTEGER NOT NULL,
    max_attendees INTEGER DEFAULT 1,
    type appointment_type NOT NULL DEFAULT 'single',
    owner_id UUID NOT NULL REFERENCES users(id),
    app_code TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    description TEXT,
    attendees_booked INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON appointments(deleted_at);

CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id),
    app_code TEXT NOT NULL,
    user_id UUID REFERENCES users(id),
    name TEXT,
    email TEXT,
    phone TEXT,
    date TIMESTAMPTZ NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    available BOOLEAN NOT NULL DEFAULT true,
    attendee_count INTEGER DEFAULT 1,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    booking_code TEXT UNIQUE NOT NULL,
    notification_status TEXT DEFAULT '',
    notification_channel TEXT DEFAULT '',
    status TEXT DEFAULT 'active',
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_bookings_deleted_at ON bookings(deleted_at);
