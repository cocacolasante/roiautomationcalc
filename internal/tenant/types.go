package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID       uuid.UUID `json:"id" db:"id"`
	APIKey   string    `json:"apiKey" db:"api_key"`
	Slug     string    `json:"slug" db:"slug"`
	Name     string    `json:"name" db:"name"`
	IsActive bool      `json:"isActive" db:"is_active"`

	CompanyName  string `json:"companyName" db:"company_name"`
	LogoURL      string `json:"logoUrl" db:"logo_url"`
	PrimaryColor string `json:"primaryColor" db:"primary_color"`
	AccentColor  string `json:"accentColor" db:"accent_color"`
	DomainName   string `json:"domainName" db:"domain_name"`

	ToolTitle       string `json:"toolTitle" db:"tool_title"`
	ToolSubtitle    string `json:"toolSubtitle" db:"tool_subtitle"`
	WelcomeMessage  string `json:"welcomeMessage" db:"welcome_message"`
	CTAText         string `json:"ctaText" db:"cta_text"`
	CTAUrl          string `json:"ctaUrl" db:"cta_url"`
	ThankyouMessage string `json:"thankyouMessage" db:"thankyou_message"`

	DefaultHourlyRate   float64            `json:"defaultHourlyRate" db:"default_hourly_rate"`
	WorkingHoursPerYear int                `json:"workingHoursPerYear" db:"working_hours_per_year"`
	AutomationSavingPct float64            `json:"automationSavingPct" db:"automation_saving_pct"`
	IndustryHourlyRates map[string]float64 `json:"industryHourlyRates" db:"industry_hourly_rates"`

	RequireContact    bool   `json:"requireContact" db:"require_contact"`
	CTASendReport     bool   `json:"ctaSendReport" db:"cta_send_report"`
	SendNotifications bool   `json:"sendNotifications" db:"send_notifications"`
	NotifyEmail       string `json:"notifyEmail" db:"notify_email"`
	NotifyDiscord     string `json:"notifyDiscord" db:"notify_discord"`
	NotifySlack       string `json:"notifySlack" db:"notify_slack"`

	CRMAdapter string `json:"crmAdapter" db:"crm_adapter"`
	CRMConfig  []byte `json:"-" db:"crm_config"`

	GoogleTagID string `json:"googleTagId" db:"google_tag_id"`
	MetaPixelID string `json:"metaPixelId" db:"meta_pixel_id"`

	EventWebhookURL string `json:"eventWebhookUrl" db:"event_webhook_url"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`

	// Blueprint Command reseller support
	PartnerID *string `json:"partner_id,omitempty" db:"partner_id"`
	ManagedBy string  `json:"managed_by" db:"managed_by"`
	ClientID  *string `json:"client_id,omitempty" db:"client_id"`
}

type UpdateConfigRequest struct {
	CompanyName         string             `json:"companyName"`
	LogoURL             string             `json:"logoUrl"`
	PrimaryColor        string             `json:"primaryColor"`
	AccentColor         string             `json:"accentColor"`
	ToolTitle           string             `json:"toolTitle"`
	ToolSubtitle        string             `json:"toolSubtitle"`
	WelcomeMessage      string             `json:"welcomeMessage"`
	CTAText             string             `json:"ctaText"`
	CTAUrl              string             `json:"ctaUrl"`
	ThankyouMessage     string             `json:"thankyouMessage"`
	DefaultHourlyRate   float64            `json:"defaultHourlyRate"`
	WorkingHoursPerYear int                `json:"workingHoursPerYear"`
	AutomationSavingPct float64            `json:"automationSavingPct"`
	IndustryHourlyRates map[string]float64 `json:"industryHourlyRates"`
	RequireContact      bool               `json:"requireContact"`
	CTASendReport       bool               `json:"ctaSendReport"`
	NotifyEmail         string             `json:"notifyEmail"`
	NotifyDiscord       string             `json:"notifyDiscord"`
	GoogleTagID         string             `json:"googleTagId"`
	MetaPixelID         string             `json:"metaPixelId"`
}

type CreateTenantRequest struct {
	Name                string             `json:"name" validate:"required"`
	Slug                string             `json:"slug" validate:"required"`
	CompanyName         string             `json:"companyName"`
	PrimaryColor        string             `json:"primaryColor"`
	AccentColor         string             `json:"accentColor"`
	ToolTitle           string             `json:"toolTitle"`
	ToolSubtitle        string             `json:"toolSubtitle"`
	WelcomeMessage      string             `json:"welcomeMessage"`
	CTAText             string             `json:"ctaText"`
	CTAUrl              string             `json:"ctaUrl"`
	ThankyouMessage     string             `json:"thankyouMessage"`
	DefaultHourlyRate   float64            `json:"defaultHourlyRate"`
	IndustryHourlyRates map[string]float64 `json:"industryHourlyRates"`
	WorkingHoursPerYear int                `json:"workingHoursPerYear"`
	AutomationSavingPct float64            `json:"automationSavingPct"`
	RequireContact      bool               `json:"requireContact"`
	CTASendReport       bool               `json:"ctaSendReport"`
	SendNotifications   bool               `json:"sendNotifications"`
	NotifyEmail         string             `json:"notifyEmail"`
	NotifyDiscord       string             `json:"notifyDiscord"`
	EventWebhookURL     string             `json:"eventWebhookUrl"`
	PartnerID           *string            `json:"partner_id,omitempty"`
	ManagedBy           string             `json:"managed_by,omitempty"`
	ClientID            *string            `json:"client_id,omitempty"`
}
