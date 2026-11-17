package finops

import (
	"math"
	"testing"
)

func TestCalculatorMonthlyCosts(t *testing.T) {
	pricing := PricingCard{
		Provider:            "Test",
		Currency:            "USD",
		CPURatePerHour:      0.04,
		MemoryRatePerGBHour: 0.005,
		HoursPerMonth:       700,
	}
	calc := NewCostCalculator(pricing)

	// 2 cores * 0.04 * 700 = $56.00
	cpuCost := calc.MonthlyCPUCost(2.0)
	if math.Abs(cpuCost-56.0) > 0.001 {
		t.Errorf("expected CPU cost 56.00, got %f", cpuCost)
	}

	// 4 GiB = 4 * 1024^3 bytes * 0.005 * 700 = $14.00
	fourGiB := int64(4 * 1024 * 1024 * 1024)
	memCost := calc.MonthlyMemoryCost(fourGiB)
	if math.Abs(memCost-14.0) > 0.001 {
		t.Errorf("expected Memory cost 14.00, got %f", memCost)
	}

	total := calc.MonthlyCost(2.0, fourGiB)
	if math.Abs(total-70.0) > 0.001 {
		t.Errorf("expected Total cost 70.00, got %f", total)
	}

	// Zero and negative inputs
	if calc.MonthlyCPUCost(0) != 0 || calc.MonthlyCPUCost(-1) != 0 {
		t.Errorf("zero/negative cpu cost should be 0")
	}
	if calc.MonthlyMemoryCost(0) != 0 || calc.MonthlyMemoryCost(-100) != 0 {
		t.Errorf("zero/negative mem cost should be 0")
	}
}

func TestCalculateWasteBreakdown(t *testing.T) {
	pricing := PricingCard{
		Provider:            "Test",
		Currency:            "USD",
		CPURatePerHour:      0.05,
		MemoryRatePerGBHour: 0.01,
		HoursPerMonth:       100,
	}
	calc := NewCostCalculator(pricing)

	// Requested: 2.0 cores, 4 GiB
	// Total cost = (2 * 0.05 * 100) + (4 * 0.01 * 100) = 10 + 4 = $14.00
	// Used: 0.5 cores, 1 GiB
	// CPU waste = (2 - 0.5) * 0.05 * 100 = 1.5 * 5 = $7.50
	// Mem waste = (4 - 1) * 0.01 * 100 = 3 * 1 = $3.00
	// Total waste = $10.50
	// Waste pct = (10.50 / 14.00) * 100 = 75%
	// Efficiency = 25%

	reqBytes := int64(4 * 1024 * 1024 * 1024)
	usedBytes := int64(1 * 1024 * 1024 * 1024)

	breakdown := calc.CalculateWaste(2.0, 0.5, reqBytes, usedBytes)

	if math.Abs(breakdown.TotalCost-14.0) > 0.01 {
		t.Errorf("expected total cost 14.00, got %f", breakdown.TotalCost)
	}
	if math.Abs(breakdown.TotalWaste-10.50) > 0.01 {
		t.Errorf("expected total waste 10.50, got %f", breakdown.TotalWaste)
	}
	if math.Abs(breakdown.CPUWaste-7.50) > 0.01 {
		t.Errorf("expected cpu waste 7.50, got %f", breakdown.CPUWaste)
	}
	if math.Abs(breakdown.MemoryWaste-3.00) > 0.01 {
		t.Errorf("expected mem waste 3.00, got %f", breakdown.MemoryWaste)
	}
	if math.Abs(breakdown.WastePercentage-75.0) > 0.1 {
		t.Errorf("expected waste percentage 75%%, got %f", breakdown.WastePercentage)
	}
	if math.Abs(breakdown.EfficiencyScore-25.0) > 0.1 {
		t.Errorf("expected efficiency score 25%%, got %f", breakdown.EfficiencyScore)
	}
}
