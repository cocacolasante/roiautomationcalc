package crm

import (
	"fmt"

	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

func NewAdapterForTenant(t *tenant.Tenant) (CRMAdapter, error) {
	switch t.CRMAdapter {
	case "hubspot":
		return NewHubSpotAdapter(t.CRMConfig)
	case "gohighlevel", "ghl":
		return NewGHLAdapter(t.CRMConfig)
	case "webhook":
		return NewWebhookAdapter(t.EventWebhookURL), nil
	case "":
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown CRM adapter: %s", t.CRMAdapter)
	}
}
