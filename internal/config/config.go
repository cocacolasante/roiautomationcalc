package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	BaseURL  string
	Env      string
	LogLevel string
	AdminKey string

	DatabaseURL string

	RedisURL         string
	AsynqConcurrency int

	AnthropicAPIKey    string
	AnthropicModel     string
	AnthropicMaxTokens int

	EncryptionKey string
	JWTSecret     string

	PDFStoragePath           string
	PDFGenerationTimeoutSecs int

	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPass        string
	NotifyFromEmail string
	NotifyFromName  string

	DefaultDiscordWebhookURL string
	TwilioAccountSID         string
	TwilioAuthToken          string
	TwilioFromNumber         string

	N8NBaseURL string
	N8NAPIKey  string

	DefaultHourlyRate          float64
	DefaultWorkingHoursPerYear int
	DefaultAutomationSavingPct float64

	HighValueROIThreshold float64

	// Blueprint Command billing
	PortalsURL     string
	BPAInternalKey string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:     getEnv("PORT", "8093"),
		BaseURL:  getEnv("BASE_URL", "http://localhost:8093"),
		Env:      getEnv("ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		AdminKey: getEnv("BLUEPRINT_ADMIN_KEY", "dev-admin-key"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://blueprint:changeme@localhost:5432/blueprint_roi?sslmode=disable"),

		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		AsynqConcurrency: getEnvInt("ASYNQ_CONCURRENCY", 10),

		AnthropicAPIKey:    getEnv("ANTHROPIC_API_KEY", ""),
		AnthropicModel:     getEnv("ANTHROPIC_MODEL", "claude-sonnet-4-20250514"),
		AnthropicMaxTokens: getEnvInt("ANTHROPIC_AUDIT_MAX_TOKENS", 1500),

		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),

		PDFStoragePath:           getEnv("PDF_STORAGE_PATH", "/data/pdfs"),
		PDFGenerationTimeoutSecs: getEnvInt("PDF_GENERATION_TIMEOUT_SECONDS", 30),

		SMTPHost:        getEnv("SMTP_HOST", ""),
		SMTPPort:        getEnvInt("SMTP_PORT", 587),
		SMTPUser:        getEnv("SMTP_USER", ""),
		SMTPPass:        getEnv("SMTP_PASS", ""),
		NotifyFromEmail: getEnv("NOTIFICATION_FROM_EMAIL", "roi@blueprintautomation.tech"),
		NotifyFromName:  getEnv("NOTIFICATION_FROM_NAME", "Blueprint ROI"),

		DefaultDiscordWebhookURL: getEnv("DEFAULT_DISCORD_WEBHOOK_URL", ""),
		TwilioAccountSID:         getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:          getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber:         getEnv("TWILIO_FROM_NUMBER", ""),

		N8NBaseURL: getEnv("N8N_BASE_URL", "http://localhost:5678"),
		N8NAPIKey:  getEnv("N8N_API_KEY", ""),

		DefaultHourlyRate:          getEnvFloat("DEFAULT_HOURLY_RATE", 35.0),
		DefaultWorkingHoursPerYear: getEnvInt("DEFAULT_WORKING_HOURS_PER_YEAR", 2080),
		DefaultAutomationSavingPct: getEnvFloat("DEFAULT_AUTOMATION_SAVING_PCT", 0.75),

		HighValueROIThreshold: getEnvFloat("HIGH_VALUE_ROI_THRESHOLD", 50000),

		PortalsURL:     getEnv("PORTALS_URL", ""),
		BPAInternalKey: getEnv("BPA_INTERNAL_KEY", ""),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.AdminKey == "" && c.Env == "production" {
		return fmt.Errorf("BLUEPRINT_ADMIN_KEY is required in production")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
