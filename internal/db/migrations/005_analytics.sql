CREATE TABLE analytics_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    audit_id    UUID REFERENCES audits(id) ON DELETE SET NULL,
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    event_type  TEXT NOT NULL,
    step        TEXT,
    metadata    JSONB NOT NULL DEFAULT '{}',
    session_id  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Event types:
-- tool_loaded | step_viewed | step_completed | step_abandoned
-- answer_given | live_roi_milestone | contact_submitted
-- audit_generated | pdf_downloaded | cta_clicked

CREATE INDEX idx_analytics_tenant ON analytics_events(tenant_id, event_type, created_at);
CREATE INDEX idx_analytics_audit ON analytics_events(audit_id);
