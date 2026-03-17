CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE tenants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key         TEXT NOT NULL UNIQUE,
    slug            TEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,

    company_name    TEXT NOT NULL DEFAULT '',
    logo_url        TEXT,
    primary_color   TEXT NOT NULL DEFAULT '#6C63FF',
    accent_color    TEXT NOT NULL DEFAULT '#F18F01',
    domain_name     TEXT,

    tool_title      TEXT NOT NULL DEFAULT 'Free Automation Audit',
    tool_subtitle   TEXT,
    welcome_message TEXT,
    cta_text        TEXT NOT NULL DEFAULT 'Book Your Free Consultation',
    cta_url         TEXT,
    thankyou_message TEXT,

    default_hourly_rate   NUMERIC(8,2) NOT NULL DEFAULT 35.00,
    working_hours_per_year INTEGER NOT NULL DEFAULT 2080,
    automation_saving_pct  NUMERIC(4,3) NOT NULL DEFAULT 0.75,
    industry_hourly_rates  JSONB NOT NULL DEFAULT '{}',

    require_contact     BOOLEAN NOT NULL DEFAULT true,
    cta_send_report     BOOLEAN NOT NULL DEFAULT true,
    send_notifications  BOOLEAN NOT NULL DEFAULT true,
    notify_email        TEXT,
    notify_discord      TEXT,
    notify_slack        TEXT,

    crm_adapter     TEXT,
    crm_config      BYTEA,

    google_tag_id   TEXT,
    meta_pixel_id   TEXT,

    event_webhook_url TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
