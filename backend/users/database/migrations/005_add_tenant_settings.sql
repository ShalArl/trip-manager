ALTER TABLE tenants
    ADD COLUMN IF NOT EXISTS settings JSONB NOT NULL DEFAULT '{}';

-- Default-Limits für bestehende Tenants setzen
UPDATE tenants
SET settings = jsonb_build_object('maxActiveTrips', 3)
WHERE tier = 'free' AND NOT settings ? 'maxActiveTrips';

UPDATE tenants
SET settings = jsonb_build_object('maxActiveTrips', 0)
WHERE tier != 'free' AND NOT settings ? 'maxActiveTrips';