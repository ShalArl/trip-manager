CREATE TABLE IF NOT EXISTS locations (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id           UUID NOT NULL,
    user_id           UUID NOT NULL,
    user_name         VARCHAR(255) NOT NULL,
    user_email        VARCHAR(255) NOT NULL,
    user_avatar_key   VARCHAR(500),
    name              VARCHAR(255) NOT NULL,
    city              VARCHAR(255) NOT NULL,
    country           VARCHAR(255) NOT NULL,
    short_description VARCHAR(500) NOT NULL,
    date_from         DATE NOT NULL,
    date_to           DATE NOT NULL,
    latitude          DOUBLE PRECISION,
    longitude         DOUBLE PRECISION,
    notes             TEXT,
    sequence          INTEGER,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_locations_trip_id ON locations(trip_id);

CREATE TABLE IF NOT EXISTS location_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    image_key   VARCHAR(500) NOT NULL,
    sequence    INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_location_images_location_id ON location_images(location_id);