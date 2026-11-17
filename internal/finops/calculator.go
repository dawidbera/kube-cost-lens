// Package finops provides financial formulas and waste aggregation metrics.
package finops

import (
	"math"
)

// BytesPerGigabyte defines the binary constant (1024^3) for converting bytes to GiB.
const BytesPerGigabyte = 1024.0 * 1024.0 * 1024.0

// CostCalculator computes infrastructure costs and idle waste metrics.
type CostCalculator struct {
	pricing PricingCard
}

// NewCostCalculator constructs a calculator instance bound to a specific pricing card.
func NewCostCalculator(pricing PricingCard) *CostCalculator {
	return &CostCalculator{
		pricing: pricing,
	}
}

// MonthlyCPUCost computes projected monthly cost for a given number of CPU cores.
func (c *CostCalculator) MonthlyCPUCost(cores float64) float64 {
	if cores <= 0 {
		return 0
	}
	return cores * c.pricing.CPURatePerHour * c.pricing.HoursPerMonth
}

// MonthlyMemoryCost computes projected monthly cost for a given memory quantity in bytes.
func (c *CostCalculator) MonthlyMemoryCost(memoryBytes int64) float64 {
	if memoryBytes <= 0 {
		return 0
	}
	gb := float64(memoryBytes) / BytesPerGigabyte
	return gb * c.pricing.MemoryRatePerGBHour * c.pricing.HoursPerMonth
}

// MonthlyCost calculates the combined monthly cost of CPU and RAM allocations.
func (c *CostCalculator) MonthlyCost(cores float64, memoryBytes int64) float64 {
	return c.MonthlyCPUCost(cores) + c.MonthlyMemoryCost(memoryBytes)
}

// MonthlyCPUWaste calculates the dollar cost of unutilized allocated CPU cores.
// Formula: max(0, requested - used) * unit_rate * hours_per_month.
func (c *CostCalculator) MonthlyCPUWaste(requestedCores, usedCores float64) float64 {
	wasteCores := requestedCores - usedCores
	if wasteCores <= 0 {
		return 0
	}
	return wasteCores * c.pricing.CPURatePerHour * c.pricing.HoursPerMonth
}

// MonthlyMemoryWaste calculates the dollar cost of unutilized allocated memory.
func (c *CostCalculator) MonthlyMemoryWaste(requestedBytes, usedBytes int64) float64 {
	wasteBytes := requestedBytes - usedBytes
	if wasteBytes <= 0 {
		return 0
	}
	wasteGB := float64(wasteBytes) / BytesPerGigabyte
	return wasteGB * c.pricing.MemoryRatePerGBHour * c.pricing.HoursPerMonth
}

// WasteBreakdown contains detailed cost and waste analysis values.
type WasteBreakdown struct {
	// TotalCost is the total monthly spend based on requested resources.
	TotalCost float64
	// TotalWaste is the total monthly spend lost to idle allocations.
	TotalWaste float64
	// CPUWaste is the monthly cost lost to idle CPU.
	CPUWaste float64
	// MemoryWaste is the monthly cost lost to idle RAM.
	MemoryWaste float64
	// WastePercentage is the percentage of requested spend that is unutilized (0-100%).
	WastePercentage float64
	// EfficiencyScore is the percentage of requested spend actively utilized (0-100%).
	EfficiencyScore float64
}

// CalculateWaste calculates complete monthly waste breakdown for requested vs used compute resources.
func (c *CostCalculator) CalculateWaste(requestedCores, usedCores float64, requestedBytes, usedBytes int64) WasteBreakdown {
	totalCost := c.MonthlyCost(requestedCores, requestedBytes)
	cpuWaste := c.MonthlyCPUWaste(requestedCores, usedCores)
	memWaste := c.MonthlyMemoryWaste(requestedBytes, usedBytes)
	totalWaste := cpuWaste + memWaste

	var wastePct float64
	if totalCost > 0 {
		wastePct = (totalWaste / totalCost) * 100.0
		if wastePct > 100.0 {
			wastePct = 100.0
		}
	}

	efficiency := 100.0 - wastePct
	if efficiency < 0 {
		efficiency = 0
	}

	return WasteBreakdown{
		TotalCost:       math.Round(totalCost*100) / 100,
		TotalWaste:      math.Round(totalWaste*100) / 100,
		CPUWaste:        math.Round(cpuWaste*100) / 100,
		MemoryWaste:     math.Round(memWaste*100) / 100,
		WastePercentage: math.Round(wastePct*10) / 10,
		EfficiencyScore: math.Round(efficiency*10) / 10,
	}
}
