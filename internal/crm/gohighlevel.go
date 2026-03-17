package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blueprintautomation/blueprint-roi/internal/audit"
)

type GHLAdapter struct {
	apiKey     string
	locationID string
	client     *http.Client
}

type ghlConfig struct {
	APIKey     string `json:"apiKey"`
	LocationID string `json:"locationId"`
}

func NewGHLAdapter(configJSON []byte) (*GHLAdapter, error) {
	var cfg ghlConfig
	if err := json.Unmarshal(configJSON, &cfg); err != nil {
		return nil, err
	}
	return &GHLAdapter{apiKey: cfg.APIKey, locationID: cfg.LocationID, client: &http.Client{}}, nil
}

func (g *GHLAdapter) CreateContact(ctx context.Context, lead *audit.Lead, _ *audit.Audit) (string, error) {
	payload := map[string]interface{}{
		"firstName":   lead.Name,
		"email":       lead.Email,
		"phone":       lead.Phone,
		"companyName": lead.Company,
		"locationId":  g.locationID,
		"tags":        []string{"roi-audit", fmt.Sprintf("score-%d", lead.LeadScore)},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://services.leadconnectorhq.com/contacts/", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Version", "2021-07-28")
	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("GHL API returned %d", resp.StatusCode)
	}
	var result struct {
		Contact struct{ ID string `json:"id"` } `json:"contact"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Contact.ID, nil
}

func (g *GHLAdapter) UpdateContact(_ context.Context, _ string, _ *audit.Lead) error {
	return nil
}
