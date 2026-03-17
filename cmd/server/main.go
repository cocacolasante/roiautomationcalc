package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/blueprintautomation/blueprint-roi/internal/ai"
	"github.com/blueprintautomation/blueprint-roi/internal/analytics"
	"github.com/blueprintautomation/blueprint-roi/internal/api"
	auditpkg "github.com/blueprintautomation/blueprint-roi/internal/audit"
	"github.com/blueprintautomation/blueprint-roi/internal/config"
	"github.com/blueprintautomation/blueprint-roi/internal/db"
	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
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
		log.Fatal("parse redis URI for asynq", zap.Error(err))
	}
	asynqClient := asynq.NewClient(redisOpts)
	defer asynqClient.Close()

	analyzer := ai.NewAnalyzer(cfg.AnthropicAPIKey, cfg.AnthropicModel, cfg.AnthropicMaxTokens)
	calculator := auditpkg.NewCalculator()
	recommender := auditpkg.NewRecommender()
	scorer := auditpkg.NewScorer()
	engine := auditpkg.NewEngine(pool, calculator, recommender, scorer, analyzer, asynqClient, log)

	tenantSvc := tenant.NewService(pool)
	tracker := analytics.NewTracker(pool, log)

	router := api.NewRouter(&api.RouterConfig{
		AdminKey:  cfg.AdminKey,
		Engine:    engine,
		TenantSvc: tenantSvc,
		DB:        pool,
		Tracker:   tracker,
		Log:       log,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("Blueprint ROI API starting", zap.String("port", cfg.Port), zap.String("env", cfg.Env))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
