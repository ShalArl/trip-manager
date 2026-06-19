CREATE TABLE IF NOT EXISTS tenants (
                                       id          VARCHAR(255) PRIMARY KEY,
                                       name        VARCHAR(255) NOT NULL,
                                       tier        VARCHAR(50)  NOT NULL DEFAULT 'free' CHECK (tier IN ('free', 'standard', 'enterprise')),
                                       status      VARCHAR(50)  NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'cancelled')),
                                       branding    JSONB        NOT NULL DEFAULT '{}',
                                       created_at  TIMESTAMP    DEFAULT NOW(),
                                       updated_at  TIMESTAMP    DEFAULT NOW()
);

INSERT INTO tenants (id, name, tier) VALUES ('default', 'Default', 'free')
ON CONFLICT DO NOTHING;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255) NOT NULL DEFAULT 'default'
        REFERENCES tenants(id),
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'tenant_member'
        CHECK (role IN ('platform_admin', 'tenant_owner', 'tenant_admin', 'tenant_member'));

CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE tablename = 'tenants' AND policyname = 'tenant_isolation_tenants'
    ) THEN
        CREATE POLICY tenant_isolation_tenants ON tenants
            FOR SELECT
            USING (id = current_setting('app.tenant_id', true));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE tablename = 'tenants' AND policyname = 'tenant_isolation_tenants'
    ) THEN
        CREATE POLICY tenant_isolation_tenants ON tenants
            USING (id = current_setting('app.tenant_id', true));
    END IF;
END $$;