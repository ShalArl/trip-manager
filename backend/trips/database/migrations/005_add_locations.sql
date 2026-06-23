CREATE TABLE IF NOT EXISTS locations
(
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id           UUID         NOT NULL REFERENCES trips (id) ON DELETE CASCADE,
    user_id           UUID         NOT NULL,
    user_name         VARCHAR(255) NOT NULL,
    user_email        VARCHAR(255) NOT NULL,
    user_avatar_key   VARCHAR(500),
    name              VARCHAR(255) NOT NULL,
    city              VARCHAR(255) NOT NULL,
    country           VARCHAR(255) NOT NULL,
    country_code      VARCHAR(10),
    short_description VARCHAR(500) NOT NULL,
    date_from         DATE         NOT NULL,
    date_to           DATE         NOT NULL,
    latitude          DOUBLE PRECISION,
    longitude         DOUBLE PRECISION,
    notes             TEXT,
    sequence          INTEGER,
    tenant_id         VARCHAR(255) NOT NULL DEFAULT 'default',
    created_at        TIMESTAMP        DEFAULT NOW(),
    updated_at        TIMESTAMP        DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS location_images
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    location_id UUID         NOT NULL REFERENCES locations (id) ON DELETE CASCADE,
    image_key   VARCHAR(500) NOT NULL,
    sequence    INTEGER      NOT NULL DEFAULT 0,
    tenant_id   VARCHAR(255) NOT NULL DEFAULT 'default',
    created_at  TIMESTAMP             DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_locations_trip_id ON locations (trip_id);
CREATE INDEX IF NOT EXISTS idx_locations_tenant_id ON locations (tenant_id);
CREATE INDEX IF NOT EXISTS idx_locations_tenant_user ON locations (tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_location_images_location_id ON location_images (location_id);


ALTER TABLE locations ENABLE ROW LEVEL SECURITY;
ALTER TABLE location_images ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE tablename = 'locations' AND policyname = 'tenant_isolation_locations'
    ) THEN
        CREATE POLICY tenant_isolation_locations ON locations
            USING (tenant_id = current_setting('app.tenant_id', true));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE tablename = 'location_images' AND policyname = 'tenant_isolation_location_images'
    ) THEN
        CREATE POLICY tenant_isolation_location_images ON location_images
            USING (tenant_id = current_setting('app.tenant_id', true));
    END IF;
END $$;