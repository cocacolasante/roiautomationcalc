package tenant

import "time"

const (
	MaxTenantsPerInstance    = 15
	IncludedSeatsPerInstance = 3
	MaxInstancesPerClient    = 2
)

type ProductInstance struct {
	ID                string    `json:"instanceId"`
	ClientID          string    `json:"clientId,omitempty"`
	InstanceNumber    int       `json:"instanceNumber"`
	PortalsInstanceID string    `json:"portalsInstanceId,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	TenantCount       int       `json:"tenantCount"`
	MaxTenants        int       `json:"maxTenants"`
	IncludedSeats     int       `json:"includedSeats"`
	OverageSeats      int       `json:"overageSeats"`
	AtCapacity        bool      `json:"atCapacity"`
}

type CreateInstanceRequest struct {
	ClientID       string `json:"client_id"`
	InstanceNumber int    `json:"instance_number"`
}

type InstancesResponse struct {
	Instances    []*ProductInstance `json:"instances"`
	TotalTenants int                `json:"total_tenants"`
	TotalMax     int                `json:"total_max"`
}
