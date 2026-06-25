CREATE TABLE IF NOT EXISTS advertisers (
                                           id          VARCHAR(255) PRIMARY KEY,
                                           firebase_uid VARCHAR(255) UNIQUE NOT NULL,
                                           email       VARCHAR(255) UNIQUE NOT NULL,
                                           name        VARCHAR(255) NOT NULL,
                                           created_at  TIMESTAMP DEFAULT NOW(),
                                           updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS advertiser_tenants (
                                                  advertiser_id VARCHAR(255) NOT NULL REFERENCES advertisers(id) ON DELETE CASCADE,
                                                  tenant_id     VARCHAR(255) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
                                                  created_at    TIMESTAMP DEFAULT NOW(),
                                                  PRIMARY KEY (advertiser_id, tenant_id)
);

CREATE INDEX IF NOT EXISTS idx_advertiser_tenants_tenant ON advertiser_tenants(tenant_id);
CREATE INDEX IF NOT EXISTS idx_advertiser_tenants_advertiser ON advertiser_tenants(advertiser_id);