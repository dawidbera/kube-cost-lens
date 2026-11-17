package finops

import (
	"testing"
)

func TestRightsizerRecommendations(t *testing.T) {
	pricing := DefaultPricingCard()
	calc := NewCostCalculator(pricing)
	rightsizer := NewRightsizer(calc, 1.20, 1.50)

	// Container with 2.0 cores requested, but only using 0.5 cores
	// 4 GiB requested, but only using 1 GiB
	reqMem := int64(4 * 1024 * 1024 * 1024)
	usedMem := int64(1 * 1024 * 1024 * 1024)

	rec := rightsizer.RecommendContainer(
		"app",
		2.0, 3.0, 0.5,
		reqMem, reqMem*2, usedMem,
	)

	// Expected Recommended CPU = 0.5 * 1.20 = 0.60
	// Expected Recommended CPU Limit = 0.60 * 1.50 = 0.90
	if rec.RecommendedCPURequest != 0.60 {
		t.Errorf("expected recommended CPU 0.60, got %f", rec.RecommendedCPURequest)
	}
	if rec.RecommendedCPULimit != 0.90 {
		t.Errorf("expected recommended CPU limit 0.90, got %f", rec.RecommendedCPULimit)
	}

	// Savings should be positive
	if rec.MonthlySavings <= 0 {
		t.Errorf("expected positive monthly savings, got %f", rec.MonthlySavings)
	}
	if rec.SavingsPercentage <= 0 || rec.SavingsPercentage > 100 {
		t.Errorf("expected valid savings percentage, got %f", rec.SavingsPercentage)
	}

	// Workload aggregation
	wRec := rightsizer.RecommendWorkload("default", "web", "Deployment", 3, []ContainerRecommendation{rec})
	if wRec.Replicas != 3 {
		t.Errorf("expected 3 replicas, got %d", wRec.Replicas)
	}
	if wRec.MonthlySavings <= rec.MonthlySavings {
		t.Errorf("aggregated savings should be higher than single container savings")
	}
}

func TestRightsizerFloors(t *testing.T) {
	pricing := DefaultPricingCard()
	calc := NewCostCalculator(pricing)
	rightsizer := NewRightsizer(calc, 1.20, 1.50)

	// Near-zero consumption
	rec := rightsizer.RecommendContainer(
		"idle",
		1.0, 1.0, 0.0001,
		1024*1024*1024, 1024*1024*1024, 1024,
	)

	if rec.RecommendedCPURequest < MinRecommendedCPU {
		t.Errorf("recommended CPU %f should be >= floor %f", rec.RecommendedCPURequest, MinRecommendedCPU)
	}
	if rec.RecommendedMemoryRequest < MinRecommendedMemoryBytes {
		t.Errorf("recommended Memory %d should be >= floor %d", rec.RecommendedMemoryRequest, MinRecommendedMemoryBytes)
	}
}
