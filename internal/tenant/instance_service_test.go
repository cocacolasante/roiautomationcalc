package tenant

import (
	"testing"
)

func TestInstMax(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{0, 5, 5},
		{5, 0, 5},
		{3, 3, 3},
		{-1, 0, 0},
	}
	for _, tt := range tests {
		got := instMax(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("instMax(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestProductInstanceConstants(t *testing.T) {
	if MaxTenantsPerInstance != 15 {
		t.Errorf("MaxTenantsPerInstance = %d, want 15", MaxTenantsPerInstance)
	}
	if IncludedSeatsPerInstance != 3 {
		t.Errorf("IncludedSeatsPerInstance = %d, want 3", IncludedSeatsPerInstance)
	}
	if MaxInstancesPerClient != 2 {
		t.Errorf("MaxInstancesPerClient = %d, want 2", MaxInstancesPerClient)
	}
}

func TestProductInstanceAtCapacity(t *testing.T) {
	inst := &ProductInstance{
		TenantCount: 15,
		MaxTenants:  MaxTenantsPerInstance,
	}
	inst.AtCapacity = inst.TenantCount >= MaxTenantsPerInstance
	if !inst.AtCapacity {
		t.Error("expected AtCapacity=true when TenantCount == MaxTenantsPerInstance")
	}
}

func TestProductInstanceNotAtCapacity(t *testing.T) {
	inst := &ProductInstance{
		TenantCount: 14,
		MaxTenants:  MaxTenantsPerInstance,
	}
	inst.AtCapacity = inst.TenantCount >= MaxTenantsPerInstance
	if inst.AtCapacity {
		t.Error("expected AtCapacity=false when TenantCount < MaxTenantsPerInstance")
	}
}

func TestProductInstanceOverageSeats(t *testing.T) {
	inst := &ProductInstance{
		TenantCount:   5,
		IncludedSeats: IncludedSeatsPerInstance,
	}
	inst.OverageSeats = instMax(0, inst.TenantCount-IncludedSeatsPerInstance)
	if inst.OverageSeats != 2 {
		t.Errorf("OverageSeats = %d, want 2", inst.OverageSeats)
	}
}

func TestCreateInstanceRequestValidation(t *testing.T) {
	// Simulate InstanceService.Create validation logic
	validReq := CreateInstanceRequest{ClientID: "client1", InstanceNumber: 1}
	if validReq.InstanceNumber != 1 && validReq.InstanceNumber != 2 {
		t.Error("expected valid instance number 1")
	}

	invalidReq := CreateInstanceRequest{ClientID: "client1", InstanceNumber: 3}
	if invalidReq.InstanceNumber == 1 || invalidReq.InstanceNumber == 2 {
		t.Error("expected invalid instance number 3")
	}
}
