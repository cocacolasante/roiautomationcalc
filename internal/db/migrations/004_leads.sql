CREATE TABLE leads (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    audit_id        UUID NOT NULL REFERENCES audits(id) ON DELETE CASCADE,
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name            TEXT,
    email           TEXT,
    phone           TEXT,
    company         TEXT,
    lead_score      INTEGER,
    estimated_roi   NUMERIC,
    top_products    JSONB NOT NULL DEFAULT '[]',
    crm_contact_id  TEXT,
    status          TEXT NOT NULL DEFAULT 'new',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leads_tenant ON leads(tenant_id);
CREATE INDEX idx_leads_email ON leads(tenant_id, email);
CREATE INDEX idx_leads_score ON leads(tenant_id, lead_score DESC);
