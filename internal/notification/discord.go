package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{webhookURL: webhookURL, client: &http.Client{}}
}

type discordEmbed struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Color       int                  `json:"color"`
	Fields      []discordEmbedField  `json:"fields"`
}

type discordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func (d *DiscordNotifier) SendAuditComplete(ctx context.Context, name, company string, roiAmount float64, score int, tenantName string) error {
	if d.webhookURL == "" {
		return nil
	}
	color := 0x7289DA
	emoji := "📊"
	if roiAmount >= 50000 {
		color = 0xFF6B35
		emoji = "🔥"
	}
	embed := discordEmbed{
		Title:       fmt.Sprintf("%s New Audit Completed", emoji),
		Description: fmt.Sprintf("**%s** from **%s** just completed an automation audit", name, company),
		Color:       color,
		Fields: []discordEmbedField{
			{Name: "Annual Opportunity", Value: fmt.Sprintf("$%.0f", roiAmount), Inline: true},
			{Name: "Lead Score", Value: fmt.Sprintf("%d/100", score), Inline: true},
			{Name: "Tool", Value: tenantName, Inline: true},
		},
	}
	payload := map[string]interface{}{"embeds": []discordEmbed{embed}}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", d.webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook returned %d", resp.StatusCode)
	}
	return nil
}
