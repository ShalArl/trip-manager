DO $$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'users_app') THEN
            EXECUTE format('CREATE USER users_app WITH PASSWORD %L', '{{APP_DB_PASSWORD}}');
        END IF;
    END $$;

GRANT CONNECT ON DATABASE users TO users_app;
GRANT USAGE ON SCHEMA public TO users_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO users_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO users_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO users_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO users_app;

ALTER TABLE users FORCE ROW LEVEL SECURITY;
ALTER TABLE tenants FORCE ROW LEVEL SECURITY;
ALTER TABLE tenant_invitations FORCE ROW LEVEL SECURITY;
ALTER TABLE advertisers FORCE ROW LEVEL SECURITY;
ALTER TABLE advertiser_tenants FORCE ROW LEVEL SECURITY;

-- Policies
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'users' AND policyname = 'tenant_isolation_users') THEN
        CREATE POLICY tenant_isolation_users ON users
            USING (tenant_id = current_setting('app.tenant_id', true))
            WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'tenant_invitations' AND policyname = 'tenant_isolation_invitations') THEN
        CREATE POLICY tenant_isolation_invitations ON tenant_invitations
            USING (tenant_id = current_setting('app.tenant_id', true))
            WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'advertisers' AND policyname = 'tenant_isolation_advertisers') THEN
        CREATE POLICY tenant_isolation_advertisers ON advertisers
            USING (true) WITH CHECK (true);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'advertiser_tenants' AND policyname = 'tenant_isolation_advertiser_tenants') THEN
        CREATE POLICY tenant_isolation_advertiser_tenants ON advertiser_tenants
            USING (true) WITH CHECK (true);
    END IF;
END $$;

-- Platform Admin Policy für tenants (darf alle sehen)
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'tenants' AND policyname = 'platform_admin_tenants') THEN
        CREATE POLICY platform_admin_tenants ON tenants
            USING (current_setting('app.tenant_id', true) IN ('default', ''));
    END IF;
END $$;

-- Gleiches für users (Platform Admin darf alle User sehen)
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'users' AND policyname = 'platform_admin_users') THEN
        CREATE POLICY platform_admin_users ON users
            USING (current_setting('app.tenant_id', true) IN ('default', ''));
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'tenants' AND policyname = 'tenant_create') THEN
        CREATE POLICY tenant_create ON tenants FOR INSERT WITH CHECK (true);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'tenants' AND policyname = 'tenant_isolation_tenants') THEN
        CREATE POLICY tenant_isolation_tenants ON tenants
            USING (id = current_setting('app.tenant_id', true))
            WITH CHECK (id = current_setting('app.tenant_id', true));
    END IF;
END $$;