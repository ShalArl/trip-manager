-- Add image_key column to locations table
ALTER TABLE locations
    ADD COLUMN IF NOT EXISTS image_key VARCHAR(500);