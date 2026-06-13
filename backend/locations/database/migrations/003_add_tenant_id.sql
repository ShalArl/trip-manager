ALTER TABLE locations
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

ALTER TABLE location_images
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

CREATE INDEX IF NOT EXISTS idx_locations_tenant_id ON locations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_locations_tenant_user ON locations(tenant_id, user_id);

ALTER TABLE locations ENABLE ROW LEVEL SECURITY;
ALTER TABLE location_images ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_locations ON locations
    USING (tenant_id = current_setting('app.tenant_id', true));

CREATE POLICY tenant_isolation_location_images ON location_images
    USING (tenant_id = current_setting('app.tenant_id', true));