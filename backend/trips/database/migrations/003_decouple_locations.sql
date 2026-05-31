-- 003_decouple_locations.sql
-- Entkopplung: location_id Fremdschlüssel entfernen, Ortsdaten direkt einbetten

-- TRANSPORTS
ALTER TABLE transports
    DROP COLUMN IF EXISTS from_location_id,
    DROP COLUMN IF EXISTS to_location_id,
    ADD COLUMN IF NOT EXISTS from_name         VARCHAR(255),
    ADD COLUMN IF NOT EXISTS from_city         VARCHAR(255),
    ADD COLUMN IF NOT EXISTS from_country      VARCHAR(255),
    ADD COLUMN IF NOT EXISTS from_country_code VARCHAR(255),
    ADD COLUMN IF NOT EXISTS from_lat          DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS from_lng          DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS to_name           VARCHAR(255),
    ADD COLUMN IF NOT EXISTS to_city           VARCHAR(255),
    ADD COLUMN IF NOT EXISTS to_country        VARCHAR(255),
    ADD COLUMN IF NOT EXISTS to_country_code   VARCHAR(255),
    ADD COLUMN IF NOT EXISTS to_lat            DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS to_lng            DOUBLE PRECISION;

-- ACCOMMODATIONS
ALTER TABLE accommodations
    DROP COLUMN IF EXISTS location_id,
    ADD COLUMN IF NOT EXISTS location_name         VARCHAR(255),
    ADD COLUMN IF NOT EXISTS location_city         VARCHAR(255),
    ADD COLUMN IF NOT EXISTS location_country      VARCHAR(255),
    ADD COLUMN IF NOT EXISTS location_country_code VARCHAR(255),
    ADD COLUMN IF NOT EXISTS location_lat          DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS location_lng          DOUBLE PRECISION;