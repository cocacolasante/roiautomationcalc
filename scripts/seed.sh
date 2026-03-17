#!/usr/bin/env bash
set -euo pipefail

DB_URL="${DATABASE_URL:-postgres://blueprint:blueprint@localhost:5432/blueprint_roi?sslmode=disable}"

echo "→ Seeding database..."

psql "$DB_URL" <<'SQL'
-- Default tenant
INSERT INTO tenants (id, name, slug, admin_email, active, crm_adapter, created_at)
VALUES (
  'aaaaaaaa-0000-0000-0000-000000000001',
  'Blueprint Demo',
  'default',
  'admin@blueprintautomation.io',
  true,
  'webhook',
  NOW()
) ON CONFLICT (slug) DO NOTHING;

-- Sample questions
INSERT INTO questions (id, tenant_id, category, key, question, hint, input_type, unit, min_value, max_value, automation_pct, is_required, sort_order)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-0000-0000-0000-000000000001', 'email',      'email_hours',       'How many hours/week does your team spend on email triage & follow-up?', 'Include reading, sorting, and writing routine emails', 'hours_slider', 'hrs/wk', 0, 40, 0.70, true,  1),
  (gen_random_uuid(), 'aaaaaaaa-0000-0000-0000-000000000001', 'leads',      'lead_entry_hours',  'How many hours/week is spent manually entering or qualifying leads?',    'CRM data entry, lead scoring, initial outreach',         'hours_slider', 'hrs/wk', 0, 40, 0.80, true,  2),
  (gen_random_uuid(), 'aaaaaaaa-0000-0000-0000-000000000001', 'billing',    'invoice_hours',     'How many hours/week does invoicing and payment follow-up take?',         'Creating invoices, sending reminders, reconciling payments','hours_slider','hrs/wk',0,20,0.75, false, 3),
  (gen_random_uuid(), 'aaaaaaaa-0000-0000-0000-000000000001', 'scheduling', 'scheduling_hours',  'How many hours/week are spent on scheduling meetings or appointments?',  'Back-and-forth emails, calendar management',             'hours_slider', 'hrs/wk', 0, 20, 0.85, false, 4),
  (gen_random_uuid(), 'aaaaaaaa-0000-0000-0000-000000000001', 'reporting',  'reporting_hours',   'How many hours/week do you spend pulling reports or compiling data?',    'Weekly reports, dashboards, spreadsheet updates',        'hours_slider', 'hrs/wk', 0, 20, 0.65, false, 5),
  (gen_random_uuid(), 'aaaaaaaa-0000-0000-0000-000000000001', 'support',    'support_hours',     'How many hours/week does your team handle customer support queries?',    'FAQs, ticket responses, status updates',                 'hours_slider', 'hrs/wk', 0, 40, 0.60, false, 6)
ON CONFLICT DO NOTHING;

-- Tool config
INSERT INTO tool_config (tenant_id, config, updated_at)
VALUES (
  'aaaaaaaa-0000-0000-0000-000000000001',
  '{
    "toolTitle": "Free Automation Audit",
    "toolSubtitle": "Discover exactly where your business is losing time — and how to get it back.",
    "welcomeMessage": "Answer 10 quick questions and get a personalized report showing your automation ROI potential.",
    "companyName": "Blueprint Automation",
    "ctaButtonText": "Book a Free Strategy Call",
    "ctaUrl": "https://calendly.com/blueprintautomation/strategy",
    "ctaDescription": "Let us build a custom automation roadmap for your business.",
    "defaultHourlyRate": 75,
    "primaryColor": "#2563eb"
  }',
  NOW()
) ON CONFLICT (tenant_id) DO UPDATE SET config = EXCLUDED.config, updated_at = NOW();
SQL

echo "✓ Seed complete."
