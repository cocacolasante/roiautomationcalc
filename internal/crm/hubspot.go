package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blueprintautomation/blueprint-roi/internal/audit"
)

type HubSpotAdapter struct {
	apiKey string
	client *http.Client
}

type hubSpotConfig struct {
	APIKey string `json:"apiKey"`
}

func NewHubSpotAdapter(configJSON []byte) (*HubSpotAdapter, error) {
	var cfg hubSpotConfig
	if err := json.Unmarshal(configJSON, &cfg); err != nil {
		return nil, err
	}
	return &HubSpotAdapter{apiKey: cfg.APIKey, client: &http.Client{}}, nil
}

func (h *HubSpotAdapter) CreateContact(ctx context.Context, lead *audit.Lead, _ *audit.Audit) (string, error) {
	payload := map[string]interface{}{
		"properties": map[string]string{
			"firstname":           lead.Name,
			"email":               lead.Email,
			"phone":               lead.Phone,
			"company":             lead.Company,
			"roi_estimate":        fmt.Sprintf("%.0f", lead.EstimatedROI),
			"lead_score":          fmt.Sprintf("%d", lead.LeadScore),
			"automation_audit_id": lead.AuditID.String(),
		},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.hubapi.com/crm/v3/objects/contacts", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("hubspot API returned %d", resp.StatusCode)
	}
	var result struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.ID, nil
}

func (h *HubSpotAdapter) UpdateContact(_ context.Context, _ string, _ *audit.Lead) error {
	return nil
}
