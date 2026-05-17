CREATE TABLE IF NOT EXISTS transports (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id          UUID NOT NULL,
    user_id          UUID NOT NULL,
    user_name        VARCHAR(255) NOT NULL,
    user_email       VARCHAR(255) NOT NULL,
    from_location_id UUID NOT NULL,
    to_location_id   UUID NOT NULL,
    departure_time   TIMESTAMP,
    arrival_time     TIMESTAMP,
    type             VARCHAR(50) NOT NULL CHECK (type IN ('flight', 'train', 'car', 'bus')),
    notes            TEXT,
    created_at       TIMESTAMP DEFAULT NOW(),
    updated_at       TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transports_trip_id ON transports(trip_id);

CREATE TABLE IF NOT EXISTS accommodations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id         UUID NOT NULL,
    user_id         UUID NOT NULL,
    user_name       VARCHAR(255) NOT NULL,
    user_email      VARCHAR(255) NOT NULL,
    location_id     UUID NOT NULL,
    name            VARCHAR(255) NOT NULL,
    address         VARCHAR(500),
    check_in        TIMESTAMP,
    check_out       TIMESTAMP,
    price_per_night DOUBLE PRECISION,
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accommodations_trip_id ON accommodations(trip_id);