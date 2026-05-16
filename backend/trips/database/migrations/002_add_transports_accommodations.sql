-- 002_add_transports_accommodations.sql

CREATE TABLE IF NOT EXISTS transports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id         UUID NOT NULL,
    type            VARCHAR(50) NOT NULL CHECK (type IN ('flight', 'train', 'bus', 'car', 'ferry', 'other')),
    departure_place VARCHAR(255) NOT NULL,
    arrival_place   VARCHAR(255) NOT NULL,
    departure_time  TIMESTAMP,
    arrival_time    TIMESTAMP,
    booking_ref     VARCHAR(255),
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transports_trip_id ON transports(trip_id);

CREATE TABLE IF NOT EXISTS accommodations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id     UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    address     VARCHAR(500),
    check_in    TIMESTAMP,
    check_out   TIMESTAMP,
    booking_ref VARCHAR(255),
    notes       TEXT,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accommodations_trip_id ON accommodations(trip_id);