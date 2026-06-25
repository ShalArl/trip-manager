CREATE TABLE IF NOT EXISTS tenant_invitations (
                                                  id          VARCHAR(255) PRIMARY KEY,
                                                  tenant_id   VARCHAR(255) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
                                                  email       VARCHAR(255) NOT NULL,
                                                  role        VARCHAR(50)  NOT NULL DEFAULT 'tenant_member',
                                                  token       VARCHAR(255) NOT NULL UNIQUE,
                                                  created_by  VARCHAR(255) NOT NULL,
                                                  expires_at  TIMESTAMP    NOT NULL,
                                                  accepted_at TIMESTAMP,
                                                  created_at  TIMESTAMP    DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_invitations_token ON tenant_invitations(token);
CREATE INDEX IF NOT EXISTS idx_invitations_tenant ON tenant_invitations(tenant_id);