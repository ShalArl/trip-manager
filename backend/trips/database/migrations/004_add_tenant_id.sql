ALTER TABLE trips
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

ALTER TABLE transports
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

ALTER TABLE accommodations
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

CREATE INDEX IF NOT EXISTS idx_trips_tenant_id ON trips(tenant_id);
CREATE INDEX IF NOT EXISTS idx_trips_tenant_user ON trips(tenant_id, user_id);

-- Row-Level Security
ALTER TABLE trips ENABLE ROW LEVEL SECURITY;
ALTER TABLE transports ENABLE ROW LEVEL SECURITY;
ALTER TABLE accommodations ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_trips ON trips
    USING (tenant_id = current_setting('app.tenant_id', true));

CREATE POLICY tenant_isolation_transports ON transports
    USING (tenant_id = current_setting('app.tenant_id', true));

CREATE POLICY tenant_isolation_accommodations ON accommodations
    USING (tenant_id = current_setting('app.tenant_id', true));