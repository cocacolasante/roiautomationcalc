package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type EventType string

const (
	EventToolLoaded       EventType = "tool_loaded"
	EventStepViewed       EventType = "step_viewed"
	EventStepCompleted    EventType = "step_completed"
	EventStepAbandoned    EventType = "step_abandoned"
	EventAnswerGiven      EventType = "answer_given"
	EventROIMilestone     EventType = "live_roi_milestone"
	EventContactSubmitted EventType = "contact_submitted"
	EventAuditGenerated   EventType = "audit_generated"
	EventPDFDownloaded    EventType = "pdf_downloaded"
	EventCTAClicked       EventType = "cta_clicked"
)

type Event struct {
	ID        uuid.UUID              `json:"id"`
	AuditID   *uuid.UUID             `json:"auditId,omitempty"`
	TenantID  uuid.UUID              `json:"tenantId"`
	EventType EventType              `json:"eventType"`
	Step      string                 `json:"step,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	SessionID string                 `json:"sessionId,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
}

type FunnelData struct {
	Loads     int `json:"loads"`
	Started   int `json:"started"`
	Completed int `json:"completed"`
}

type TierBreakdown struct {
	Hot  int `json:"hot"`
	Warm int `json:"warm"`
	Cold int `json:"cold"`
}

type FunnelStats struct {
	TotalLeads      int           `json:"totalLeads"`
	AvgLeadScore    float64       `json:"avgLeadScore"`
	AvgRoiValue     float64       `json:"avgRoiValue"`
	Funnel          FunnelData    `json:"funnel"`
	TierBreakdown   TierBreakdown `json:"tierBreakdown"`
	PDFsDownloaded  int           `json:"pdfsDownloaded"`
	CTAsClicked     int           `json:"ctasClicked"`
	CompletionRate  float64       `json:"completionRate"`
	LeadCaptureRate float64       `json:"leadCaptureRate"`
}

type Tracker struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewTracker(db *pgxpool.Pool, log *zap.Logger) *Tracker {
	return &Tracker{db: db, log: log}
}

func (t *Tracker) Track(ctx context.Context, event *Event) {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	event.CreatedAt = time.Now()

	go func() {
		_, err := t.db.Exec(context.Background(), `
			INSERT INTO analytics_events (id, audit_id, tenant_id, event_type, step, metadata, session_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			event.ID, event.AuditID, event.TenantID, event.EventType,
			event.Step, event.Metadata, event.SessionID, event.CreatedAt,
		)
		if err != nil {
			t.log.Warn("Failed to track event", zap.Error(err))
		}
	}()
}

func (t *Tracker) GetFunnelStats(ctx context.Context, tenantID uuid.UUID, days int) (*FunnelStats, error) {
	since := time.Now().AddDate(0, 0, -days)
	stats := &FunnelStats{}

	// Funnel events from analytics_events
	rows, err := t.db.Query(ctx, `
		SELECT event_type, COUNT(*) as count
		FROM analytics_events
		WHERE tenant_id = $1 AND created_at >= $2
		GROUP BY event_type`, tenantID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auditsCompleted, contactsProvided int
	for rows.Next() {
		var eventType string
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			continue
		}
		switch EventType(eventType) {
		case EventToolLoaded:
			stats.Funnel.Loads = count
		case EventStepViewed:
			stats.Funnel.Started = count
		case EventAuditGenerated:
			auditsCompleted = count
			stats.Funnel.Completed = count
		case EventContactSubmitted:
			contactsProvided = count
		case EventPDFDownloaded:
			stats.PDFsDownloaded = count
		case EventCTAClicked:
			stats.CTAsClicked = count
		}
	}

	if stats.Funnel.Loads > 0 {
		stats.CompletionRate = float64(auditsCompleted) / float64(stats.Funnel.Loads)
	}
	if auditsCompleted > 0 {
		stats.LeadCaptureRate = float64(contactsProvided) / float64(auditsCompleted)
	}

	// Lead counts, avg score, avg ROI, and tier breakdown from leads + audits tables
	t.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(AVG(l.lead_score), 0),
			COALESCE(AVG(l.estimated_roi), 0),
			COUNT(*) FILTER (WHERE LOWER(a.lead_quality) = 'hot'),
			COUNT(*) FILTER (WHERE LOWER(a.lead_quality) = 'warm'),
			COUNT(*) FILTER (WHERE LOWER(a.lead_quality) = 'cold')
		FROM leads l
		JOIN audits a ON l.audit_id = a.id
		WHERE l.tenant_id = $1 AND l.created_at >= $2`,
		tenantID, since,
	).Scan(
		&stats.TotalLeads,
		&stats.AvgLeadScore,
		&stats.AvgRoiValue,
		&stats.TierBreakdown.Hot,
		&stats.TierBreakdown.Warm,
		&stats.TierBreakdown.Cold,
	)

	return stats, nil
}
