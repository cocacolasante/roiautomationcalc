package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/blueprintautomation/blueprint-roi/internal/ai"
	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	TaskGeneratePDF      = "pdf:generate"
	TaskSyncCRM          = "crm:sync"
	TaskSendNotification = "notification:send"
)

type Engine struct {
	db          *pgxpool.Pool
	calculator  *Calculator
	recommender *Recommender
	scorer      *Scorer
	analyzer    *ai.Analyzer
	queue       *asynq.Client
	log         *zap.Logger
}

func NewEngine(
	db *pgxpool.Pool,
	calculator *Calculator,
	recommender *Recommender,
	scorer *Scorer,
	analyzer *ai.Analyzer,
	queue *asynq.Client,
	log *zap.Logger,
) *Engine {
	return &Engine{db: db, calculator: calculator, recommender: recommender,
		scorer: scorer, analyzer: analyzer, queue: queue, log: log}
}

func (e *Engine) Generate(ctx context.Context, req *AuditRequest, t *tenant.Tenant, cfg *ToolConfig) (*Audit, error) {
	roiCalc := e.calculator.Calculate(req.Answers, cfg)
	leadScore, leadQuality := e.scorer.Score(req.Answers, roiCalc, t)

	// Build recommendations from cfg products
	var recommendations []*Recommendation
	for _, product := range cfg.Products {
		score := scoreProduct(product, req.Answers, roiCalc)
		if score < 0.3 {
			continue
		}
		monthlyROI := product.EstHoursSavedPerWeek * 52 / 12 * roiCalc.HourlyRate
		priority := "low"
		if score >= 0.7 && roiCalc.RecoverableCostPerYear >= 25000 {
			priority = "high"
		} else if score >= 0.5 {
			priority = "medium"
		}
		insight := product.ROIDescription
		for _, key := range product.RelevantKeys {
			if item := roiCalc.FindByKey(key); item != nil && item.HoursPerWeek > 0 {
				insight = fmt.Sprintf("Based on your %.1f hours/week on related tasks, this could save significant time.", item.HoursPerWeek)
				break
			}
		}
		recommendations = append(recommendations, &Recommendation{
			Product: product, RelevanceScore: score,
			EstimatedHoursSaved: product.EstHoursSavedPerWeek, EstimatedMonthlySaved: monthlyROI,
			PaybackPeriodMonths: product.TimeToROIMonths, SpecificInsight: insight, Priority: priority,
		})
	}

	analyzeReq := buildAnalyzeRequest(req, roiCalc, recommendations, t)
	analysis, err := e.analyzer.Analyze(ctx, analyzeReq)
	if err != nil {
		e.log.Warn("AI analysis failed, using fallback", zap.Error(err))
		analysis = fallbackAnalysis(req, roiCalc)
	}

	var findings []Finding
	for _, f := range analysis.Findings {
		findings = append(findings, Finding{
			Category: f.Category, Title: f.Title, Description: f.Description,
			Impact: f.Impact, HoursPerWeek: f.HoursPerWeek, AnnualCost: f.AnnualCost,
		})
	}

	now := time.Now()
	audit := &Audit{
		ID:                   uuid.New(),
		TenantID:             t.ID,
		BusinessName:         req.BusinessName,
		ContactName:          req.ContactName,
		ContactEmail:         req.ContactEmail,
		ContactPhone:         req.ContactPhone,
		TeamSize:             req.TeamSize,
		Industry:             req.Industry,
		Answers:              req.Answers,
		ROISummary:           roiCalc,
		Findings:             findings,
		ExecutiveSummary:     analysis.Executive,
		TopPriority:          analysis.TopPriority,
		Recommendations:      recommendations,
		LeadScore:            leadScore,
		LeadQuality:          leadQuality,
		EstimatedAnnualValue: roiCalc.RecoverableCostPerYear,
		Status:               "complete",
		SourceURL:            req.SourceURL,
		UTMSource:            req.UTMSource,
		UTMCampaign:          req.UTMCampaign,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := e.saveAudit(ctx, audit); err != nil {
		return nil, fmt.Errorf("save audit: %w", err)
	}
	if err := e.createLead(ctx, audit); err != nil {
		e.log.Warn("Failed to create lead", zap.Error(err))
	}

	e.enqueueTask(TaskGeneratePDF, audit.ID.String())
	e.enqueueTask(TaskSyncCRM, audit.ID.String())
	e.enqueueTask(TaskSendNotification, audit.ID.String())

	return audit, nil
}

func (e *Engine) saveAudit(ctx context.Context, a *Audit) error {
	answersJSON, _ := json.Marshal(a.Answers)
	roiJSON, _ := json.Marshal(a.ROISummary)
	findingsJSON, _ := json.Marshal(a.Findings)
	recsJSON, _ := json.Marshal(a.Recommendations)

	_, err := e.db.Exec(ctx, `
		INSERT INTO audits (
			id, tenant_id, business_name, contact_name, contact_email,
			contact_phone, team_size, industry, answers,
			roi_summary, findings, executive_summary, top_priority,
			recommendations, lead_score, lead_quality, estimated_annual_value,
			status, source_url, utm_source, utm_campaign, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23
		)`,
		a.ID, a.TenantID, a.BusinessName, a.ContactName, a.ContactEmail,
		a.ContactPhone, a.TeamSize, a.Industry, answersJSON,
		roiJSON, findingsJSON, a.ExecutiveSummary, a.TopPriority,
		recsJSON, a.LeadScore, a.LeadQuality, a.EstimatedAnnualValue,
		a.Status, a.SourceURL, a.UTMSource, a.UTMCampaign,
		a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (e *Engine) createLead(ctx context.Context, a *Audit) error {
	var topProducts []string
	for i, rec := range a.Recommendations {
		if i >= 3 {
			break
		}
		topProducts = append(topProducts, rec.Product.Key)
	}
	topProductsJSON, _ := json.Marshal(topProducts)
	_, err := e.db.Exec(ctx, `
		INSERT INTO leads (audit_id, tenant_id, name, email, phone, company,
			lead_score, estimated_roi, top_products, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		a.ID, a.TenantID, a.ContactName, a.ContactEmail, a.ContactPhone,
		a.BusinessName, a.LeadScore, a.EstimatedAnnualValue, topProductsJSON, "new",
	)
	return err
}

func (e *Engine) enqueueTask(taskType, id string) {
	task := asynq.NewTask(taskType, []byte(`{"id":"`+id+`"}`))
	if _, err := e.queue.Enqueue(task); err != nil {
		e.log.Warn("Failed to enqueue task", zap.String("type", taskType), zap.Error(err))
	}
}

func scoreProduct(product *Product, answers map[string]interface{}, roi *ROICalculation) float64 {
	score := 0.0
	for _, key := range product.RelevantKeys {
		answer, ok := answers[key]
		if !ok {
			continue
		}
		if item := roi.FindByKey(key); item != nil && item.HoursPerWeek >= product.MinWeeklyHours {
			score += 0.3
		}
		for reqKey, reqValue := range product.RequiredAnswers {
			if reqKey == key {
				strAnswer := fmt.Sprintf("%v", answer)
				if strAnswer == reqValue {
					score += 0.4
				}
			}
		}
	}
	if score > 1.0 {
		score = 1.0
	}
	return score
}

func buildAnalyzeRequest(req *AuditRequest, roi *ROICalculation, recs []*Recommendation, _ *tenant.Tenant) *ai.AnalyzeRequest {
	var wf strings.Builder
	for _, item := range roi.Breakdown {
		wf.WriteString(fmt.Sprintf("- %s: %.1f hours/week ($%.0f/year)\n",
			item.Description, item.HoursPerWeek, item.AnnualCost))
	}
	var pl strings.Builder
	for _, rec := range recs {
		pl.WriteString(fmt.Sprintf("- %s: %s (est. %.1f hours/week saved)\n",
			rec.Product.Name, rec.Product.Description, rec.EstimatedHoursSaved))
	}
	recoverablePct := 0.0
	if roi.TotalCostPerYear > 0 {
		recoverablePct = roi.RecoverableCostPerYear / roi.TotalCostPerYear * 100
	}
	teamSize := req.TeamSize
	if teamSize == 0 {
		if ts, ok := req.Answers["team_size"]; ok {
			if f, ok2 := ts.(float64); ok2 {
				teamSize = int(f)
			}
		}
	}
	return &ai.AnalyzeRequest{
		BusinessName: req.BusinessName, BusinessType: req.Industry,
		TeamSize: teamSize, Industry: req.Industry,
		WorkflowBreakdown: wf.String(), TotalHours: roi.TotalHoursPerWeek,
		AnnualCost: roi.TotalCostPerYear, RecoverableAnnual: roi.RecoverableCostPerYear,
		RecoverablePct: recoverablePct, ProductsList: pl.String(),
	}
}

func fallbackAnalysis(req *AuditRequest, roi *ROICalculation) *ai.AuditAnalysis {
	businessName := req.BusinessName
	if businessName == "" {
		businessName = "your business"
	}
	return &ai.AuditAnalysis{
		Executive: fmt.Sprintf(
			"Based on the workflow data for %s, your team is spending approximately %.1f hours per week on tasks that could be significantly reduced through automation. "+
				"At your current pace, this represents $%.0f per year in labor costs applied to manual processes. "+
				"With the right automation systems, you could recover approximately $%.0f of that annually — freeing your team to focus on higher-value work.",
			businessName, roi.TotalHoursPerWeek, roi.TotalCostPerYear, roi.RecoverableCostPerYear,
		),
		TopPriority: "Focus on automating your highest-volume repetitive tasks first for the fastest return on investment.",
	}
}
