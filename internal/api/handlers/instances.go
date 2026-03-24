package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

// InstancesHandler handles product instance management endpoints.
type InstancesHandler struct {
	instanceSvc *tenant.InstanceService
}

// NewInstancesHandler creates a new InstancesHandler.
func NewInstancesHandler(instanceSvc *tenant.InstanceService) *InstancesHandler {
	return &InstancesHandler{instanceSvc: instanceSvc}
}

// Create handles POST /api/admin/instances.
func (h *InstancesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req tenant.CreateInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.ClientID == "" {
		http.Error(w, `{"error":"client_id is required"}`, http.StatusBadRequest)
		return
	}
	if req.InstanceNumber == 0 {
		req.InstanceNumber = 2
	}
	inst, err := h.instanceSvc.Create(r.Context(), req)
	if err != nil {
		if err == tenant.ErrInstanceLimitReached {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if err == tenant.ErrInvalidInstanceNumber {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inst)
}

// List handles GET /api/admin/instances.
func (h *InstancesHandler) List(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID != "" {
		instances, err := h.instanceSvc.ListByClientID(r.Context(), clientID)
		if err != nil {
			http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
			return
		}
		total := 0
		for _, inst := range instances {
			total += inst.TenantCount
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tenant.InstancesResponse{
			Instances:    instances,
			TotalTenants: total,
			TotalMax:     len(instances) * tenant.MaxTenantsPerInstance,
		})
		return
	}
	instances, err := h.instanceSvc.ListAll(r.Context())
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	total := 0
	for _, inst := range instances {
		total += inst.TenantCount
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant.InstancesResponse{
		Instances:    instances,
		TotalTenants: total,
		TotalMax:     len(instances) * tenant.MaxTenantsPerInstance,
	})
}

// Get handles GET /api/admin/instances/{id}.
func (h *InstancesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	inst, err := h.instanceSvc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"instance not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inst)
}
