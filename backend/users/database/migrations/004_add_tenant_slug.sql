ALTER TABLE tenants
    ADD COLUMN IF NOT EXISTS slug VARCHAR(255) UNIQUE;

UPDATE tenants
SET slug = LOWER(REGEXP_REPLACE(name, '[^a-zA-Z0-9]+', '-', 'g'))
WHERE slug IS NULL;