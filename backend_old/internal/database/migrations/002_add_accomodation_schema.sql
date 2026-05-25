CREATE TABLE IF NOT EXISTS accommodations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id         UUID NOT NULL,
    user_id         UUID NOT NULL,
    location_id     UUID NOT NULL,
    name            VARCHAR(255) NOT NULL,
    address         VARCHAR(255),
    check_in        TIMESTAMP,
    check_out       TIMESTAMP,
    price_per_night NUMERIC(10, 2),
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (trip_id)     REFERENCES trips(id)     ON DELETE CASCADE,
    FOREIGN KEY (user_id)     REFERENCES users(id)     ON DELETE CASCADE,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_accommodations_trip_id     ON accommodations(trip_id);
CREATE INDEX IF NOT EXISTS idx_accommodations_user_id     ON accommodations(user_id);
CREATE INDEX IF NOT EXISTS idx_accommodations_location_id ON accommodations(location_id);