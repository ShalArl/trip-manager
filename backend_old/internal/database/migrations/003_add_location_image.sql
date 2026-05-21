-- Create location_images table for multiple images per location
CREATE TABLE IF NOT EXISTS location_images (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    image_key  VARCHAR(500) NOT NULL,
    sequence   INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_location_images_location_id ON location_images(location_id);