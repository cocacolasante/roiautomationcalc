package tenant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Provisioner handles automated tenant setup including default questions and tool config.
type Provisioner struct {
	db *pgxpool.Pool
}

func NewProvisioner(db *pgxpool.Pool) *Provisioner {
	return &Provisioner{db: db}
}

// Provision creates a tenant and seeds it with default questions and tool config.
func (p *Provisioner) Provision(ctx context.Context, req CreateTenantRequest) (*Tenant, error) {
	svc := NewService(p.db)

	t, err := svc.Create(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("create tenant: %w", err)
	}

	if err := p.seedDefaultQuestions(ctx, t.ID); err != nil {
		return nil, fmt.Errorf("seed questions: %w", err)
	}

	if err := p.seedToolConfig(ctx, t); err != nil {
		return nil, fmt.Errorf("seed tool config: %w", err)
	}

	return t, nil
}

func (p *Provisioner) seedDefaultQuestions(ctx context.Context, tenantID uuid.UUID) error {
	type q struct {
		category      string
		key           string
		question      string
		hint          string
		inputType     string
		unit          string
		automationPct float64
		isRequired    bool
		sortOrder     int
	}

	defaults := []q{
		{"email", "email_hours", "How many hours/week does your team spend on email triage & follow-up?", "Include reading, sorting, and writing routine emails", "hours_slider", "hrs/wk", 0.70, true, 1},
		{"leads", "lead_entry_hours", "How many hours/week is spent manually entering or qualifying leads?", "CRM data entry, lead scoring, initial outreach", "hours_slider", "hrs/wk", 0.80, true, 2},
		{"billing", "invoice_hours", "How many hours/week does invoicing and payment follow-up take?", "Creating invoices, sending reminders, reconciling payments", "hours_slider", "hrs/wk", 0.75, false, 3},
		{"scheduling", "scheduling_hours", "How many hours/week are spent on scheduling meetings or appointments?", "Back-and-forth emails, calendar management", "hours_slider", "hrs/wk", 0.85, false, 4},
		{"reporting", "reporting_hours", "How many hours/week do you spend pulling reports or compiling data?", "Weekly reports, dashboards, spreadsheet updates", "hours_slider", "hrs/wk", 0.65, false, 5},
		{"support", "support_hours", "How many hours/week does your team handle customer support queries?", "FAQs, ticket responses, status updates", "hours_slider", "hrs/wk", 0.60, false, 6},
		{"social", "social_hours", "How many hours/week do you spend on social media posting and scheduling?", "Content creation excluded — just scheduling and posting", "hours_slider", "hrs/wk", 0.55, false, 7},
	}

	for _, dq := range defaults {
		_, err := p.db.Exec(ctx, `
			INSERT INTO questions (id, tenant_id, category, key, question, hint, input_type, unit, min_value, max_value, automation_pct, is_required, sort_order)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, 0, 40, $8, $9, $10)
			ON CONFLICT (tenant_id, key) DO NOTHING`,
			tenantID, dq.category, dq.key, dq.question, dq.hint, dq.inputType, dq.unit, dq.automationPct, dq.isRequired, dq.sortOrder,
		)
		if err != nil {
			return fmt.Errorf("insert question %s: %w", dq.key, err)
		}
	}
	return nil
}

func (p *Provisioner) seedToolConfig(ctx context.Context, t *Tenant) error {
	ctaText := t.CTAText
	if ctaText == "" {
		ctaText = "Book a Free Strategy Call"
	}
	title := t.ToolTitle
	if title == "" {
		title = "Free Automation Audit"
	}
	_, err := p.db.Exec(ctx, `
		INSERT INTO tool_config (tenant_id, config, updated_at)
		VALUES ($1, $2::jsonb, NOW())
		ON CONFLICT (tenant_id) DO UPDATE SET config = EXCLUDED.config, updated_at = NOW()`,
		t.ID,
		fmt.Sprintf(`{"toolTitle":%q,"companyName":%q,"ctaButtonText":%q,"ctaUrl":%q,"defaultHourlyRate":75,"primaryColor":"#2563eb"}`,
			title, t.CompanyName, ctaText, t.CTAUrl),
	)
	return err
}
