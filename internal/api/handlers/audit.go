package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	auditpkg "github.com/blueprintautomation/blueprint-roi/internal/audit"
	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

type AuditHandler struct {
	engine    *auditpkg.Engine
	tenantSvc *tenant.Service
	configH   *ConfigHandler
	db        *pgxpool.Pool
	log       *zap.Logger
}

func NewAuditHandler(engine *auditpkg.Engine, tenantSvc *tenant.Service, configH *ConfigHandler, db *pgxpool.Pool, log *zap.Logger) *AuditHandler {
	return &AuditHandler{engine: engine, tenantSvc: tenantSvc, configH: configH, db: db, log: log}
}

func (h *AuditHandler) Submit(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
	ctx := r.Context()

	var t *tenant.Tenant
	var err error
	if parsedID, parseErr := uuid.Parse(tenantID); parseErr == nil {
		t, err = h.tenantSvc.GetByID(ctx, parsedID)
	} else {
		t, err = h.tenantSvc.GetBySlug(ctx, tenantID)
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
		return
	}

	var req auditpkg.AuditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	req.TenantID = t.ID

	questions, _ := h.configH.GetQuestions(ctx, t.ID)
	products, _ := h.configH.GetProducts(ctx, t.ID)

	cfg := &auditpkg.ToolConfig{
		TenantID:            t.ID.String(),
		DefaultHourlyRate:   t.DefaultHourlyRate,
		AutomationSavingPct: t.AutomationSavingPct,
		WorkingHoursPerYear: t.WorkingHoursPerYear,
		Questions:           questions,
		Products:            products,
	}

	result, err := h.engine.Generate(ctx, &req, t, cfg)
	if err != nil {
		h.log.Error("Audit generation failed", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"audit generation failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *AuditHandler) GetAudit(w http.ResponseWriter, r *http.Request) {
	auditID := chi.URLParam(r, "auditId")
	ctx := r.Context()

	id, err := uuid.Parse(auditID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"invalid audit id"}`, http.StatusBadRequest)
		return
	}

	var a auditpkg.Audit
	var roiJSON, findingsJSON, recsJSON, answersJSON []byte
	err = h.db.QueryRow(ctx, `
		SELECT id, tenant_id, business_name, contact_name, contact_email,
		       answers, roi_summary, findings, COALESCE(executive_summary,''), COALESCE(top_priority,''),
		       recommendations, lead_score, COALESCE(lead_quality,''), estimated_annual_value,
		       status, created_at
		FROM audits WHERE id = $1`, id).Scan(
		&a.ID, &a.TenantID, &a.BusinessName, &a.ContactName, &a.ContactEmail,
		&answersJSON, &roiJSON, &findingsJSON, &a.ExecutiveSummary, &a.TopPriority,
		&recsJSON, &a.LeadScore, &a.LeadQuality, &a.EstimatedAnnualValue,
		&a.Status, &a.CreatedAt,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"audit not found"}`, http.StatusNotFound)
		return
	}

	json.Unmarshal(roiJSON, &a.ROISummary)
	json.Unmarshal(findingsJSON, &a.Findings)
	json.Unmarshal(recsJSON, &a.Recommendations)
	json.Unmarshal(answersJSON, &a.Answers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}
