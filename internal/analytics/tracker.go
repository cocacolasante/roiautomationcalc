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

type FunnelStats struct {
	ToolLoads        int     `json:"toolLoads"`
	StepsStarted     int     `json:"stepsStarted"`
	AuditsCompleted  int     `json:"auditsCompleted"`
	ContactsProvided int     `json:"contactsProvided"`
	PDFsDownloaded   int     `json:"pdfsDownloaded"`
	CTAsClicked      int     `json:"ctasClicked"`
	CompletionRate   float64 `json:"completionRate"`
	LeadCaptureRate  float64 `json:"leadCaptureRate"`
	AvgROICalculated float64 `json:"avgRoiCalculated"`
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

	rows, err := t.db.Query(ctx, `
		SELECT event_type, COUNT(*) as count
		FROM analytics_events
		WHERE tenant_id = $1 AND created_at >= $2
		GROUP BY event_type`, tenantID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventType string
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			continue
		}
		switch EventType(eventType) {
		case EventToolLoaded:
			stats.ToolLoads = count
		case EventStepViewed:
			stats.StepsStarted = count
		case EventAuditGenerated:
			stats.AuditsCompleted = count
		case EventContactSubmitted:
			stats.ContactsProvided = count
		case EventPDFDownloaded:
			stats.PDFsDownloaded = count
		case EventCTAClicked:
			stats.CTAsClicked = count
		}
	}

	if stats.ToolLoads > 0 {
		stats.CompletionRate = float64(stats.AuditsCompleted) / float64(stats.ToolLoads)
	}
	if stats.AuditsCompleted > 0 {
		stats.LeadCaptureRate = float64(stats.ContactsProvided) / float64(stats.AuditsCompleted)
	}

	t.db.QueryRow(ctx, `
		SELECT COALESCE(AVG(estimated_annual_value), 0)
		FROM audits WHERE tenant_id = $1 AND created_at >= $2 AND status = 'complete'`,
		tenantID, since,
	).Scan(&stats.AvgROICalculated)

	return stats, nil
}
