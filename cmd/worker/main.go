package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	auditpkg "github.com/blueprintautomation/blueprint-roi/internal/audit"
	"github.com/blueprintautomation/blueprint-roi/internal/config"
	"github.com/blueprintautomation/blueprint-roi/internal/db"
	"github.com/blueprintautomation/blueprint-roi/internal/notification"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	var log *zap.Logger
	if cfg.Env == "development" {
		log, _ = zap.NewDevelopment()
	} else {
		log, _ = zap.NewProduction()
	}
	defer log.Sync()

	ctx := context.Background()

	pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("postgres connect failed", zap.Error(err))
	}
	defer pool.Close()

	redisOpts, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		log.Fatal("parse redis URI", zap.Error(err))
	}

	srv := asynq.NewServer(redisOpts, asynq.Config{
		Concurrency: cfg.AsynqConcurrency,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	mux := asynq.NewServeMux()

	mux.HandleFunc(auditpkg.TaskSendNotification, func(ctx context.Context, t *asynq.Task) error {
		var payload struct{ ID string }
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		var name, company, tenantName string
		var roi float64
		var score int
		var discordWebhook string

		pool.QueryRow(ctx, `
			SELECT a.contact_name, a.business_name, a.estimated_annual_value,
			       a.lead_score, tn.company_name, COALESCE(tn.notify_discord,'')
			FROM audits a
			JOIN tenants tn ON a.tenant_id = tn.id
			WHERE a.id = $1`, payload.ID).Scan(
			&name, &company, &roi, &score, &tenantName, &discordWebhook,
		)

		webhookURL := discordWebhook
		if webhookURL == "" {
			webhookURL = cfg.DefaultDiscordWebhookURL
		}

		if webhookURL != "" {
			notifier := notification.NewDiscordNotifier(webhookURL)
			if err := notifier.SendAuditComplete(ctx, name, company, roi, score, tenantName); err != nil {
				log.Warn("Discord notification failed", zap.Error(err))
			}
		}
		return nil
	})

	mux.HandleFunc(auditpkg.TaskGeneratePDF, func(ctx context.Context, t *asynq.Task) error {
		var payload struct{ ID string }
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		log.Info("PDF generation queued", zap.String("auditId", payload.ID))
		return nil
	})

	mux.HandleFunc(auditpkg.TaskSyncCRM, func(ctx context.Context, t *asynq.Task) error {
		var payload struct{ ID string }
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		log.Info("CRM sync queued", zap.String("auditId", payload.ID))
		return nil
	})

	log.Info("Blueprint ROI Worker starting", zap.Int("concurrency", cfg.AsynqConcurrency))
	if err := srv.Run(mux); err != nil {
		log.Fatal("worker error", zap.Error(err))
	}
}
