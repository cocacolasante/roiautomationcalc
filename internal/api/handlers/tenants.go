package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/blueprintautomation/blueprint-roi/internal/billing"
	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

type TenantsHandler struct {
	tenantSvc *tenant.Service
	log       *zap.Logger
	db        *pgxpool.Pool
	billing   *billing.TenantLimitClient
}

func NewTenantsHandler(tenantSvc *tenant.Service, log *zap.Logger, billingClient *billing.TenantLimitClient) *TenantsHandler {
	return &TenantsHandler{tenantSvc: tenantSvc, log: log, billing: billingClient}
}

// WithDB sets the DB pool on the handler (used for stats queries).
func (h *TenantsHandler) WithDB(db *pgxpool.Pool) *TenantsHandler {
	h.db = db
	return h
}

func (h *TenantsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req tenant.CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if !h.billing.CheckTenantAllowed(r.Context(), req.PortalsInstanceID) {
		http.Error(w, `{"error":"tenant limit reached"}`, http.StatusPaymentRequired)
		return
	}

	t, err := h.tenantSvc.Create(r.Context(), &req)
	if err != nil {
		h.log.Error("Failed to create tenant", zap.Error(err))
		http.Error(w, `{"error":"failed to create tenant"}`, http.StatusInternalServerError)
		return
	}

	go h.billing.IncrementTenantCount(context.Background(), req.PortalsInstanceID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TenantsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	t, err := h.tenantSvc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TenantsHandler) List(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.tenantSvc.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"tenants": tenants})
}

// GetStats returns aggregated metrics for a tenant (Blueprint Command reseller endpoint).
func (h *TenantsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	if _, err := h.tenantSvc.GetByID(r.Context(), id); err != nil {
		http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
		return
	}

	ctx := r.Context()

	var auditsCompleted int
	_ = h.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM audits WHERE tenant_id = $1 AND status = 'complete' AND created_at >= date_trunc('month', NOW())`,
		id,
	).Scan(&auditsCompleted)

	var leadsCaptured int
	_ = h.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM leads WHERE tenant_id = $1 AND created_at >= date_trunc('month', NOW())`,
		id,
	).Scan(&leadsCaptured)

	var avgROI float64
	_ = h.db.QueryRow(ctx,
		`SELECT COALESCE(AVG((roi_summary->>'annual_savings')::numeric), 0) FROM audits WHERE tenant_id = $1 AND status = 'complete' AND created_at >= date_trunc('month', NOW())`,
		id,
	).Scan(&avgROI)

	var totalAudits int
	_ = h.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM audits WHERE tenant_id = $1 AND created_at >= date_trunc('month', NOW())`,
		id,
	).Scan(&totalAudits)

	var completionRate float64
	if totalAudits > 0 {
		completionRate = float64(auditsCompleted) / float64(totalAudits) * 100
	}

	stats := map[string]interface{}{
		"audits_completed":       auditsCompleted,
		"leads_captured":         leadsCaptured,
		"avg_roi_calculated_usd": avgROI,
		"completion_rate_pct":    completionRate,
	}
	summary := fmt.Sprintf(
		"%d ROI audits completed this month, generating %d qualified leads.",
		auditsCompleted, leadsCaptured,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tenant_id":   id,
		"product_key": "roiautomationcalc",
		"stats":       stats,
		"summary":     summary,
	})
}
