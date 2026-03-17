package crm

import (
	"context"

	"github.com/blueprintautomation/blueprint-roi/internal/audit"
)

// CRMAdapter is the interface all CRM integrations implement.
type CRMAdapter interface {
	CreateContact(ctx context.Context, lead *audit.Lead, auditData *audit.Audit) (string, error)
	UpdateContact(ctx context.Context, contactID string, lead *audit.Lead) error
}
