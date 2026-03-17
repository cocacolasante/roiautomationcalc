package tenant

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, req *CreateTenantRequest) (*Tenant, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("generate api key: %w", err)
	}

	slug := strings.ToLower(strings.ReplaceAll(req.Slug, " ", "-"))
	primaryColor := coalesce(req.PrimaryColor, "#6C63FF")
	accentColor := coalesce(req.AccentColor, "#F18F01")
	toolTitle := coalesce(req.ToolTitle, "Free Automation Audit")
	ctaText := coalesce(req.CTAText, "Book Your Free Consultation")
	hourlyRate := req.DefaultHourlyRate
	if hourlyRate == 0 {
		hourlyRate = 35.0
	}
	workingHours := req.WorkingHoursPerYear
	if workingHours == 0 {
		workingHours = 2080
	}
	automationPct := req.AutomationSavingPct
	if automationPct == 0 {
		automationPct = 0.75
	}

	industryRates := "{}"
	if len(req.IndustryHourlyRates) > 0 {
		b, _ := json.Marshal(req.IndustryHourlyRates)
		industryRates = string(b)
	}

	managedBy := req.ManagedBy
	if managedBy == "" {
		managedBy = "bpa"
	}

	var id uuid.UUID
	err = s.db.QueryRow(ctx, `
		INSERT INTO tenants (
			api_key, slug, name, company_name, primary_color, accent_color,
			tool_title, tool_subtitle, welcome_message, cta_text, cta_url,
			thankyou_message, default_hourly_rate, working_hours_per_year,
			automation_saving_pct, industry_hourly_rates, require_contact,
			cta_send_report, send_notifications, notify_email, notify_discord,
			event_webhook_url, created_at, updated_at,
			partner_id, managed_by, client_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
			$25, $26, $27
		) RETURNING id`,
		apiKey, slug, req.Name, coalesce(req.CompanyName, req.Name),
		primaryColor, accentColor, toolTitle,
		req.ToolSubtitle, req.WelcomeMessage, ctaText, req.CTAUrl,
		req.ThankyouMessage, hourlyRate, workingHours,
		automationPct, industryRates,
		req.RequireContact, req.CTASendReport, req.SendNotifications,
		req.NotifyEmail, req.NotifyDiscord, req.EventWebhookURL,
		time.Now(), time.Now(),
		req.PartnerID, managedBy, req.ClientID,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert tenant: %w", err)
	}

	return s.GetByID(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	t := &Tenant{}
	var industryRatesJSON []byte

	err := s.db.QueryRow(ctx, `SELECT
		id, api_key, slug, name, is_active,
		company_name, COALESCE(logo_url,''), primary_color, accent_color, COALESCE(domain_name,''),
		tool_title, COALESCE(tool_subtitle,''), COALESCE(welcome_message,''), cta_text,
		COALESCE(cta_url,''), COALESCE(thankyou_message,''),
		default_hourly_rate, working_hours_per_year, automation_saving_pct,
		industry_hourly_rates,
		require_contact, cta_send_report, send_notifications,
		COALESCE(notify_email,''), COALESCE(notify_discord,''), COALESCE(notify_slack,''),
		COALESCE(crm_adapter,''), COALESCE(google_tag_id,''), COALESCE(meta_pixel_id,''),
		COALESCE(event_webhook_url,''),
		created_at, updated_at,
		partner_id, managed_by, client_id
		FROM tenants WHERE id = $1`, id).Scan(
		&t.ID, &t.APIKey, &t.Slug, &t.Name, &t.IsActive,
		&t.CompanyName, &t.LogoURL, &t.PrimaryColor, &t.AccentColor, &t.DomainName,
		&t.ToolTitle, &t.ToolSubtitle, &t.WelcomeMessage, &t.CTAText,
		&t.CTAUrl, &t.ThankyouMessage,
		&t.DefaultHourlyRate, &t.WorkingHoursPerYear, &t.AutomationSavingPct,
		&industryRatesJSON,
		&t.RequireContact, &t.CTASendReport, &t.SendNotifications,
		&t.NotifyEmail, &t.NotifyDiscord, &t.NotifySlack,
		&t.CRMAdapter, &t.GoogleTagID, &t.MetaPixelID,
		&t.EventWebhookURL,
		&t.CreatedAt, &t.UpdatedAt,
		&t.PartnerID, &t.ManagedBy, &t.ClientID,
	)
	if err != nil {
		return nil, fmt.Errorf("get tenant: %w", err)
	}
	if len(industryRatesJSON) > 0 {
		json.Unmarshal(industryRatesJSON, &t.IndustryHourlyRates)
	}
	return t, nil
}

// ListByPartner returns all tenants for a given partner ID.
func (s *Service) ListByPartner(ctx context.Context, partnerID string) ([]*Tenant, error) {
	rows, err := s.db.Query(ctx, `SELECT id FROM tenants WHERE partner_id = $1 ORDER BY created_at DESC`, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*Tenant
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			continue
		}
		t, err := s.GetByID(ctx, id)
		if err != nil {
			continue
		}
		tenants = append(tenants, t)
	}
	if tenants == nil {
		tenants = []*Tenant{}
	}
	return tenants, nil
}

func (s *Service) UpdateConfig(ctx context.Context, id uuid.UUID, req *UpdateConfigRequest) (*Tenant, error) {
	industryRates := "{}"
	if len(req.IndustryHourlyRates) > 0 {
		b, _ := json.Marshal(req.IndustryHourlyRates)
		industryRates = string(b)
	}
	_, err := s.db.Exec(ctx, `
		UPDATE tenants SET
			company_name = $1, logo_url = $2, primary_color = $3, accent_color = $4,
			tool_title = $5, tool_subtitle = $6, welcome_message = $7,
			cta_text = $8, cta_url = $9, thankyou_message = $10,
			default_hourly_rate = $11, working_hours_per_year = $12,
			automation_saving_pct = $13, industry_hourly_rates = $14,
			require_contact = $15, cta_send_report = $16,
			notify_email = $17, notify_discord = $18,
			google_tag_id = $19, meta_pixel_id = $20,
			updated_at = NOW()
		WHERE id = $21`,
		req.CompanyName, req.LogoURL, req.PrimaryColor, req.AccentColor,
		req.ToolTitle, req.ToolSubtitle, req.WelcomeMessage,
		req.CTAText, req.CTAUrl, req.ThankyouMessage,
		req.DefaultHourlyRate, req.WorkingHoursPerYear,
		req.AutomationSavingPct, industryRates,
		req.RequireContact, req.CTASendReport,
		req.NotifyEmail, req.NotifyDiscord,
		req.GoogleTagID, req.MetaPixelID,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("update config: %w", err)
	}
	return s.GetByID(ctx, id)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `SELECT id FROM tenants WHERE slug = $1 AND is_active = true`, slug).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("get tenant by slug: %w", err)
	}
	return s.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Tenant, error) {
	rows, err := s.db.Query(ctx, `SELECT id, slug, name, company_name, primary_color, is_active, created_at FROM tenants ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*Tenant
	for rows.Next() {
		t := &Tenant{}
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.CompanyName, &t.PrimaryColor, &t.IsActive, &t.CreatedAt); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	if tenants == nil {
		tenants = []*Tenant{}
	}
	return tenants, nil
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "broi_" + hex.EncodeToString(b), nil
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
