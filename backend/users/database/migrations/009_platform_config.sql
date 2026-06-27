CREATE TABLE IF NOT EXISTS platform_config (
                                               key     VARCHAR(255) PRIMARY KEY,
                                               value   JSONB        NOT NULL,
                                               updated_at TIMESTAMP DEFAULT NOW()
);

-- Default-Werte
INSERT INTO platform_config (key, value) VALUES
                                             ('pricing_free',       '{"basePrice": 0, "freeApiCalls": 0, "pricePerCall": 0}'),
                                             ('pricing_standard',   '{"basePrice": 29, "freeApiCalls": 10000, "pricePerCall": 0.001}'),
                                             ('pricing_enterprise', '{"basePrice": 99, "freeApiCalls": 100000, "pricePerCall": 0.0005}')
ON CONFLICT (key) DO NOTHING;

ALTER TABLE platform_config ENABLE ROW LEVEL SECURITY;

CREATE POLICY platform_config_read ON platform_config
    USING (true);

CREATE POLICY platform_config_write ON platform_config
    FOR ALL
    WITH CHECK (current_setting('app.tenant_id', true) IN ('default', ''));