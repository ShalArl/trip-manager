ALTER TABLE trips
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

ALTER TABLE transports
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

ALTER TABLE accommodations
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default';

CREATE INDEX IF NOT EXISTS idx_trips_tenant_id ON trips(tenant_id);
CREATE INDEX IF NOT EXISTS idx_trips_tenant_user ON trips(tenant_id, user_id);

ALTER TABLE trips ENABLE ROW LEVEL SECURITY;
ALTER TABLE transports ENABLE ROW LEVEL SECURITY;
ALTER TABLE accommodations ENABLE ROW LEVEL SECURITY;

ALTER TABLE trips FORCE ROW LEVEL SECURITY;
ALTER TABLE transports FORCE ROW LEVEL SECURITY;
ALTER TABLE accommodations FORCE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE tablename = 'trips' AND policyname = 'tenant_isolation_trips'
    ) THEN
        CREATE POLICY tenant_isolation_trips ON trips
            USING (tenant_id = current_setting('app.tenant_id', true));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE tablename = 'transports' AND policyname = 'tenant_isolation_transports'
    ) THEN
        CREATE POLICY tenant_isolation_transports ON transports
            USING (tenant_id = current_setting('app.tenant_id', true));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE tablename = 'accommodations' AND policyname = 'tenant_isolation_accommodations'
    ) THEN
        CREATE POLICY tenant_isolation_accommodations ON accommodations
            USING (tenant_id = current_setting('app.tenant_id', true));
    END IF;
END $$;