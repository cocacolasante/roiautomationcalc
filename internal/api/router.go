package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/blueprintautomation/blueprint-roi/internal/analytics"
	auditpkg "github.com/blueprintautomation/blueprint-roi/internal/audit"
	"github.com/blueprintautomation/blueprint-roi/internal/api/handlers"
	"github.com/blueprintautomation/blueprint-roi/internal/api/middleware"
	"github.com/blueprintautomation/blueprint-roi/internal/billing"
	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

type RouterConfig struct {
	AdminKey      string
	Engine        *auditpkg.Engine
	TenantSvc     *tenant.Service
	DB            *pgxpool.Pool
	Tracker       *analytics.Tracker
	Log           *zap.Logger
	BillingClient *billing.TenantLimitClient
	InstanceSvc   *tenant.InstanceService
}

func NewRouter(cfg *RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Blueprint-Admin-Key"},
		MaxAge:         300,
	}))

	configH := handlers.NewConfigHandler(cfg.TenantSvc, cfg.DB, cfg.Log)
	auditH := handlers.NewAuditHandler(cfg.Engine, cfg.TenantSvc, configH, cfg.DB, cfg.Log)
	leadsH := handlers.NewLeadsHandler(cfg.DB, cfg.Log)
	tenantsH := handlers.NewTenantsHandler(cfg.TenantSvc, cfg.Log, cfg.BillingClient).WithDB(cfg.DB)
	analyticsH := handlers.NewAnalyticsHandler(cfg.Tracker, cfg.Log)
	healthSvc := handlers.NewHealthService(cfg.InstanceSvc)

	if cfg.InstanceSvc != nil {
		tenantsH = tenantsH.WithInstanceService(cfg.InstanceSvc)
	}

	instancesH := handlers.NewInstancesHandler(cfg.InstanceSvc)

	r.Get("/api/health", healthSvc.Health)

	embedFS := http.FileServer(http.Dir("./embed"))
	r.Handle("/embed/*", http.StripPrefix("/embed", embedFS))

	// Admin dashboard SPA — served at /admin/
	// Any sub-path falls back to index.html so React Router handles routing.
	adminFS := http.Dir("./admin-ui")
	r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
	})
	r.Handle("/admin/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/admin")
		f, err := adminFS.Open(path)
		if err != nil {
			// SPA fallback — let React Router handle the route
			http.ServeFile(w, r, "./admin-ui/index.html")
			return
		}
		f.Close()
		http.StripPrefix("/admin", http.FileServer(adminFS)).ServeHTTP(w, r)
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.RateLimit(60, time.Minute))
		r.Get("/config/{tenantId}", configH.GetConfig)
		r.Post("/audit/{tenantId}", auditH.Submit)
		r.Get("/audit/{auditId}", auditH.GetAudit)
		r.Post("/analytics", analyticsH.TrackEvent)
	})

	r.Route("/api/admin", func(r chi.Router) {
		r.Use(middleware.AdminAuth(cfg.AdminKey))
		r.Use(middleware.PartnerScopeMiddleware)
		r.Get("/verify", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok":true}`))
		})
		r.Route("/tenants", func(r chi.Router) {
			r.Get("/", tenantsH.List)
			r.Post("/", tenantsH.Create)
			r.Get("/{id}", tenantsH.Get)
			r.Delete("/{id}", tenantsH.Delete)
			r.Post("/{id}/suspend", tenantsH.Suspend)
			r.Post("/{id}/unsuspend", tenantsH.Unsuspend)
			r.Get("/{tenantId}/config", configH.GetConfig)
			r.Put("/{tenantId}/config", configH.UpdateConfig)
			r.Get("/{id}/leads", leadsH.ListLeads)
			r.Get("/{id}/leads/{leadId}", leadsH.GetLead)
			r.Patch("/{id}/leads/{leadId}/status", leadsH.UpdateLeadStatus)
			r.Get("/{id}/analytics", analyticsH.GetStats)
			r.Get("/{id}/stats", tenantsH.GetStats)
		})

		// Instance management
		r.Get("/instances", instancesH.List)
		r.Post("/instances", instancesH.Create)
		r.Get("/instances/{id}", instancesH.Get)
	})

	return r
}
