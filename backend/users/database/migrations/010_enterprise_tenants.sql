CREATE TABLE IF NOT EXISTS enterprise_tenants (
                                                  tenant_id  VARCHAR(255) PRIMARY KEY,
                                                  db_url     TEXT NOT NULL,
                                                  created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE enterprise_tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE enterprise_tenants FORCE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'enterprise_tenants' AND policyname = 'enterprise_tenants_platform_admin') THEN
        CREATE POLICY enterprise_tenants_platform_admin ON enterprise_tenants
            USING (current_setting('app.tenant_id', true) IN ('default', ''));
    END IF;
END $$;