package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/blueprintautomation/blueprint-roi/internal/audit"
	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

type ConfigHandler struct {
	tenantSvc *tenant.Service
	db        *pgxpool.Pool
	log       *zap.Logger
}

func NewConfigHandler(tenantSvc *tenant.Service, db *pgxpool.Pool, log *zap.Logger) *ConfigHandler {
	return &ConfigHandler{tenantSvc: tenantSvc, db: db, log: log}
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
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

	questions, _ := h.GetQuestions(ctx, t.ID)
	products, _ := h.GetProducts(ctx, t.ID)

	cfg := &audit.ToolConfig{
		TenantID:            t.ID.String(),
		CompanyName:         t.CompanyName,
		LogoURL:             t.LogoURL,
		PrimaryColor:        t.PrimaryColor,
		AccentColor:         t.AccentColor,
		ToolTitle:           t.ToolTitle,
		ToolSubtitle:        t.ToolSubtitle,
		WelcomeMessage:      t.WelcomeMessage,
		CTAText:             t.CTAText,
		CTAUrl:              t.CTAUrl,
		ThankyouMessage:     t.ThankyouMessage,
		DefaultHourlyRate:   t.DefaultHourlyRate,
		IndustryHourlyRate:  t.IndustryHourlyRates,
		WorkingHoursPerYear: t.WorkingHoursPerYear,
		AutomationSavingPct: t.AutomationSavingPct,
		Questions:           questions,
		Products:            products,
		RequireContact:      t.RequireContact,
		CTASendReport:       t.CTASendReport,
		GoogleTagID:         t.GoogleTagID,
		MetaPixelID:         t.MetaPixelID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (h *ConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "tenantId"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	var req tenant.UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	t, err := h.tenantSvc.UpdateConfig(r.Context(), id, &req)
	if err != nil {
		h.log.Error("Failed to update config", zap.Error(err))
		http.Error(w, `{"error":"failed to update config"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *ConfigHandler) GetQuestions(ctx context.Context, tenantID uuid.UUID) ([]*audit.Question, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, tenant_id, key, category, question, COALESCE(hint,''), input_type, options,
		       COALESCE(min_value,0), COALESCE(max_value,40), COALESCE(default_value,0),
		       unit, hours_per_unit, is_automatic, automation_pct, relevant_products,
		       sort_order, is_active, is_required
		FROM questions WHERE tenant_id = $1 AND is_active = true ORDER BY sort_order`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*audit.Question
	for rows.Next() {
		q := &audit.Question{}
		var optionsJSON, productsJSON []byte
		if err := rows.Scan(
			&q.ID, &q.TenantID, &q.Key, &q.Category, &q.Question, &q.Hint,
			&q.InputType, &optionsJSON, &q.MinValue, &q.MaxValue, &q.DefaultValue,
			&q.Unit, &q.HoursPerUnit, &q.IsAutomatic, &q.AutomationPct, &productsJSON,
			&q.SortOrder, &q.IsActive, &q.IsRequired,
		); err != nil {
			continue
		}
		json.Unmarshal(optionsJSON, &q.Options)
		json.Unmarshal(productsJSON, &q.RelevantProducts)
		questions = append(questions, q)
	}
	return questions, nil
}

func (h *ConfigHandler) GetProducts(ctx context.Context, tenantID uuid.UUID) ([]*audit.Product, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, tenant_id, key, name, description, category,
		       COALESCE(setup_fee,0), COALESCE(monthly_fee,0), COALESCE(roi_description,''),
		       min_weekly_hours, relevant_keys, required_answers,
		       est_hours_saved_pw, time_to_roi_months, sort_order, is_active
		FROM products WHERE tenant_id = $1 AND is_active = true ORDER BY sort_order`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*audit.Product
	for rows.Next() {
		p := &audit.Product{}
		var keysJSON, answersJSON []byte
		if err := rows.Scan(
			&p.ID, &p.TenantID, &p.Key, &p.Name, &p.Description, &p.Category,
			&p.SetupFee, &p.MonthlyFee, &p.ROIDescription,
			&p.MinWeeklyHours, &keysJSON, &answersJSON,
			&p.EstHoursSavedPerWeek, &p.TimeToROIMonths, &p.SortOrder, &p.IsActive,
		); err != nil {
			continue
		}
		json.Unmarshal(keysJSON, &p.RelevantKeys)
		json.Unmarshal(answersJSON, &p.RequiredAnswers)
		products = append(products, p)
	}
	return products, nil
}
