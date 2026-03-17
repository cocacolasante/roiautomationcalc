package ai

import (
	"context"
	"encoding/json"
	"fmt"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Analyzer struct {
	client    *anthropic.Client
	model     string
	maxTokens int64
}

func NewAnalyzer(apiKey, model string, maxTokens int) *Analyzer {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Analyzer{client: client, model: model, maxTokens: int64(maxTokens)}
}

func (a *Analyzer) Analyze(ctx context.Context, req *AnalyzeRequest) (*AuditAnalysis, error) {
	prompt := BuildAnalyzePrompt(req)

	msg, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(a.model)),
		MaxTokens: anthropic.F(a.maxTokens),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("claude API call: %w", err)
	}

	if len(msg.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	var textContent string
	for _, block := range msg.Content {
		if block.Type == "text" {
			textContent = block.Text
			break
		}
	}
	if textContent == "" {
		return nil, fmt.Errorf("no text content in Claude response")
	}

	jsonStr := extractJSON(textContent)
	var analysis AuditAnalysis
	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return nil, fmt.Errorf("parse Claude response: %w", err)
	}
	return &analysis, nil
}

func extractJSON(s string) string {
	start, end := 0, len(s)
	for i := 0; i < len(s); i++ {
		if s[i] == '{' {
			start = i
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '}' {
			end = i + 1
			break
		}
	}
	if start >= end {
		return s
	}
	return s[start:end]
}
