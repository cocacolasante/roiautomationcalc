package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blueprintautomation/blueprint-roi/internal/audit"
)

type WebhookAdapter struct {
	webhookURL string
	client     *http.Client
}

func NewWebhookAdapter(webhookURL string) *WebhookAdapter {
	return &WebhookAdapter{webhookURL: webhookURL, client: &http.Client{}}
}

func (w *WebhookAdapter) CreateContact(ctx context.Context, lead *audit.Lead, auditData *audit.Audit) (string, error) {
	payload := map[string]interface{}{
		"event":     "new_lead",
		"lead":      lead,
		"audit":     auditData,
		"roi":       auditData.ROISummary,
		"leadScore": lead.LeadScore,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", w.webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("webhook returned %d", resp.StatusCode)
	}
	return "webhook_sent", nil
}

func (w *WebhookAdapter) UpdateContact(_ context.Context, _ string, _ *audit.Lead) error {
	return nil
}
