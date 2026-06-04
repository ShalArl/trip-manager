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
        CHECK (role IN ('platform_admin', 'tenant_admin', 'tenant_member'));

CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);

-- RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_users ON users
    USING (tenant_id = current_setting('app.tenant_id', true));

CREATE POLICY tenant_isolation_tenants ON tenants
    USING (id = current_setting('app.tenant_id', true));