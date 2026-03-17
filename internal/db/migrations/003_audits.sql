CREATE TABLE audits (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    business_name   TEXT,
    contact_name    TEXT,
    contact_email   TEXT,
    contact_phone   TEXT,
    team_size       INTEGER,
    industry        TEXT,

    answers         JSONB NOT NULL DEFAULT '{}',

    roi_summary     JSONB NOT NULL DEFAULT '{}',
    findings        JSONB NOT NULL DEFAULT '[]',
    executive_summary TEXT,
    top_priority    TEXT,
    recommendations JSONB NOT NULL DEFAULT '[]',

    lead_score      INTEGER,
    lead_quality    TEXT,
    estimated_annual_value NUMERIC,

    status          TEXT NOT NULL DEFAULT 'generating',

    pdf_path        TEXT,
    pdf_generated_at TIMESTAMPTZ,

    source_url      TEXT,
    utm_source      TEXT,
    utm_campaign    TEXT,

    crm_contact_id  TEXT,
    crm_synced      BOOLEAN NOT NULL DEFAULT false,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audits_tenant ON audits(tenant_id);
CREATE INDEX idx_audits_email ON audits(tenant_id, contact_email);
CREATE INDEX idx_audits_date ON audits(tenant_id, created_at DESC);
