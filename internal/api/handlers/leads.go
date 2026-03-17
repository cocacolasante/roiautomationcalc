package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	auditpkg "github.com/blueprintautomation/blueprint-roi/internal/audit"
)

type LeadsHandler struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewLeadsHandler(db *pgxpool.Pool, log *zap.Logger) *LeadsHandler {
	return &LeadsHandler{db: db, log: log}
}

func (h *LeadsHandler) ListLeads(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "id")
	ctx := r.Context()

	id, err := uuid.Parse(tenantID)
	if err != nil {
		http.Error(w, `{"error":"invalid tenant id"}`, http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT l.id, l.audit_id, l.name, l.email, l.phone, l.company,
		       l.lead_score, l.estimated_roi, l.top_products, l.status, l.created_at,
		       COALESCE(a.executive_summary,''), COALESCE(a.top_priority,'')
		FROM leads l
		JOIN audits a ON l.audit_id = a.id
		WHERE l.tenant_id = $1
		ORDER BY l.created_at DESC LIMIT 100`, id)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leads []map[string]interface{}
	for rows.Next() {
		var lead auditpkg.Lead
		var topProductsJSON []byte
		var execSummary, topPriority string
		if err := rows.Scan(
			&lead.ID, &lead.AuditID, &lead.Name, &lead.Email, &lead.Phone, &lead.Company,
			&lead.LeadScore, &lead.EstimatedROI, &topProductsJSON, &lead.Status, &lead.CreatedAt,
			&execSummary, &topPriority,
		); err != nil {
			continue
		}
		json.Unmarshal(topProductsJSON, &lead.TopProducts)
		leads = append(leads, map[string]interface{}{
			"lead": lead, "execSummary": execSummary, "topPriority": topPriority,
		})
	}
	if leads == nil {
		leads = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"leads": leads})
}

func (h *LeadsHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	leadID := chi.URLParam(r, "leadId")
	ctx := r.Context()

	id, err := uuid.Parse(leadID)
	if err != nil {
		http.Error(w, `{"error":"invalid lead id"}`, http.StatusBadRequest)
		return
	}

	var lead auditpkg.Lead
	var topProductsJSON []byte
	err = h.db.QueryRow(ctx, `
		SELECT id, audit_id, tenant_id, name, email, phone, company,
		       lead_score, estimated_roi, top_products, status, created_at
		FROM leads WHERE id = $1`, id).Scan(
		&lead.ID, &lead.AuditID, &lead.TenantID, &lead.Name, &lead.Email,
		&lead.Phone, &lead.Company, &lead.LeadScore, &lead.EstimatedROI,
		&topProductsJSON, &lead.Status, &lead.CreatedAt,
	)
	if err != nil {
		http.Error(w, `{"error":"lead not found"}`, http.StatusNotFound)
		return
	}
	json.Unmarshal(topProductsJSON, &lead.TopProducts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}
