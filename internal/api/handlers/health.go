package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

// HealthService wraps the Health handler with optional dependencies.
type HealthService struct {
	instanceSvc *tenant.InstanceService
}

// NewHealthService creates a new HealthService.
func NewHealthService(instanceSvc *tenant.InstanceService) *HealthService {
	return &HealthService{instanceSvc: instanceSvc}
}

// Health handles GET /api/health (method form, used when instance capacity is needed).
func (h *HealthService) Health(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status":    "ok",
		"service":   "blueprint-roi",
		"timestamp": time.Now().UTC(),
	}

	if h.instanceSvc != nil {
		instances, err := h.instanceSvc.ListAll(r.Context())
		if err == nil {
			totalTenants := 0
			totalMax := 0
			for _, inst := range instances {
				totalTenants += inst.TenantCount
				totalMax += inst.MaxTenants
			}
			resp["instance_capacity"] = map[string]interface{}{
				"total_instances": len(instances),
				"total_tenants":   totalTenants,
				"total_max":       totalMax,
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Health is the standalone function kept for backward compatibility.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"service":   "blueprint-roi",
		"timestamp": time.Now().UTC(),
	})
}
