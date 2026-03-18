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
		       COALESCE(a.executive_summary,''), COALESCE(a.top_priority,''),
		       COALESCE(a.lead_quality,''), a.answers,
		       COALESCE(t.default_hourly_rate, 35), COALESCE(t.automation_saving_pct, 0.75)
		FROM leads l
		JOIN audits a ON l.audit_id = a.id
		JOIN tenants t ON l.tenant_id = t.id
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
		var topProductsJSON, answersJSON []byte
		var execSummary, topPriority, leadQuality string
		var tenantHourlyRate, tenantAutoPct float64
		if err := rows.Scan(
			&lead.ID, &lead.AuditID, &lead.Name, &lead.Email, &lead.Phone, &lead.Company,
			&lead.LeadScore, &lead.EstimatedROI, &topProductsJSON, &lead.Status, &lead.CreatedAt,
			&execSummary, &topPriority, &leadQuality, &answersJSON,
			&tenantHourlyRate, &tenantAutoPct,
		); err != nil {
			continue
		}
		json.Unmarshal(topProductsJSON, &lead.TopProducts)

		// Recompute ROI on the fly if stored value is zero but automation answers exist
		if lead.EstimatedROI == 0 && answersJSON != nil {
			var answers map[string]interface{}
			if json.Unmarshal(answersJSON, &answers) == nil {
				if _, hasAreas := answers["automation_areas"]; hasAreas {
					cfg := &auditpkg.ToolConfig{
						DefaultHourlyRate:   tenantHourlyRate,
						AutomationSavingPct: tenantAutoPct,
					}
					recomputed := auditpkg.NewCalculator().Calculate(answers, cfg)
					if recomputed != nil && recomputed.RecoverableCostPerYear > 0 {
						lead.EstimatedROI = recomputed.RecoverableCostPerYear
						updatedROI, _ := json.Marshal(recomputed)
						h.db.Exec(ctx,
							`UPDATE audits SET roi_summary=$1, estimated_annual_value=$2 WHERE id=$3`,
							updatedROI, recomputed.RecoverableCostPerYear, lead.AuditID)
						h.db.Exec(ctx,
							`UPDATE leads SET estimated_roi=$1 WHERE id=$2`,
							recomputed.RecoverableCostPerYear, lead.ID)
					}
				}
			}
		}

		leads = append(leads, map[string]interface{}{
			"id":             lead.ID,
			"auditId":        lead.AuditID,
			"contactName":    lead.Name,
			"contactEmail":   lead.Email,
			"contactPhone":   lead.Phone,
			"businessName":   lead.Company,
			"leadScore":      lead.LeadScore,
			"leadTier":       leadQuality,
			"roiAnnualValue": lead.EstimatedROI,
			"status":         lead.Status,
			"createdAt":      lead.CreatedAt,
			"execSummary":    execSummary,
			"topPriority":    topPriority,
		})
	}
	if leads == nil {
		leads = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"leads": leads})
}

func (h *LeadsHandler) UpdateLeadStatus(w http.ResponseWriter, r *http.Request) {
	leadID := chi.URLParam(r, "leadId")
	ctx := r.Context()

	id, err := uuid.Parse(leadID)
	if err != nil {
		http.Error(w, `{"error":"invalid lead id"}`, http.StatusBadRequest)
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	valid := map[string]bool{
		"new": true, "contacted": true, "qualified": true,
		"proposal": true, "closed_won": true, "closed_lost": true,
	}
	if !valid[body.Status] {
		http.Error(w, `{"error":"invalid status"}`, http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec(ctx, `UPDATE leads SET status=$1 WHERE id=$2`, body.Status, id)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

func toFloatIface(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
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
	var topProductsJSON, roiJSON, findingsJSON, answersJSON []byte
	var leadQuality, execSummary, topPriority string
	var tenantHourlyRate float64
	var tenantAutoPct float64
	err = h.db.QueryRow(ctx, `
		SELECT l.id, l.audit_id, l.tenant_id, l.name, l.email, l.phone, l.company,
		       l.lead_score, l.estimated_roi, l.top_products, l.status, l.created_at,
		       COALESCE(a.lead_quality,''), COALESCE(a.executive_summary,''),
		       COALESCE(a.top_priority,''), a.roi_summary, a.findings, a.answers,
		       COALESCE(t.default_hourly_rate, 35), COALESCE(t.automation_saving_pct, 0.75)
		FROM leads l
		JOIN audits a ON l.audit_id = a.id
		JOIN tenants t ON l.tenant_id = t.id
		WHERE l.id = $1`, id).Scan(
		&lead.ID, &lead.AuditID, &lead.TenantID, &lead.Name, &lead.Email,
		&lead.Phone, &lead.Company, &lead.LeadScore, &lead.EstimatedROI,
		&topProductsJSON, &lead.Status, &lead.CreatedAt,
		&leadQuality, &execSummary, &topPriority, &roiJSON, &findingsJSON, &answersJSON,
		&tenantHourlyRate, &tenantAutoPct,
	)
	if err != nil {
		http.Error(w, `{"error":"lead not found"}`, http.StatusNotFound)
		return
	}

	json.Unmarshal(topProductsJSON, &lead.TopProducts)

	var roiSummary interface{}
	var findings interface{}
	var answers map[string]interface{}
	json.Unmarshal(roiJSON, &roiSummary)
	json.Unmarshal(findingsJSON, &findings)
	json.Unmarshal(answersJSON, &answers)

	// If stored ROI is zero but answers contain automation selections, recompute on the fly
	if roiStored, ok := roiSummary.(map[string]interface{}); ok {
		stored := toFloatIface(roiStored["recoverableCostPerYear"])
		if stored == 0 && answers != nil {
			if _, hasAreas := answers["automation_areas"]; hasAreas {
				cfg := &auditpkg.ToolConfig{
					DefaultHourlyRate:   tenantHourlyRate,
					AutomationSavingPct: tenantAutoPct,
				}
				calc := auditpkg.NewCalculator()
				recomputed := calc.Calculate(answers, cfg)
				if recomputed != nil && recomputed.RecoverableCostPerYear > 0 {
					roiSummary = recomputed
					lead.EstimatedROI = recomputed.RecoverableCostPerYear
					// Persist the corrected values so future fetches are fast
					updatedROI, _ := json.Marshal(recomputed)
					h.db.Exec(ctx,
						`UPDATE audits SET roi_summary=$1, estimated_annual_value=$2 WHERE id=$3`,
						updatedROI, recomputed.RecoverableCostPerYear, lead.AuditID)
					h.db.Exec(ctx,
						`UPDATE leads SET estimated_roi=$1 WHERE id=$2`,
						recomputed.RecoverableCostPerYear, lead.ID)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":             lead.ID,
		"auditId":        lead.AuditID,
		"contactName":    lead.Name,
		"contactEmail":   lead.Email,
		"contactPhone":   lead.Phone,
		"businessName":   lead.Company,
		"leadScore":      lead.LeadScore,
		"leadTier":       leadQuality,
		"roiAnnualValue": lead.EstimatedROI,
		"status":         lead.Status,
		"createdAt":      lead.CreatedAt,
		"audit": map[string]interface{}{
			"roiSummary":       roiSummary,
			"findings":         findings,
			"executiveSummary": execSummary,
			"topPriority":      topPriority,
		},
		"answers": answers,
	})
}
