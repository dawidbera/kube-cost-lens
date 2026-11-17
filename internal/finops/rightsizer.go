// Package finops provides algorithms for rightsizing container compute allocations.
package finops

import (
	"math"
)

// Default sizing coefficients.
const (
	// DefaultHeadroomMargin is the standard safety buffer added to actual consumption (20%).
	DefaultHeadroomMargin = 1.20
	// DefaultLimitMultiplier is the ratio applied to recommended requests for limits (1.5x).
	DefaultLimitMultiplier = 1.50

	// MinRecommendedCPU is the floor for recommended CPU (10 millicores).
	MinRecommendedCPU = 0.010
	// MinRecommendedMemoryBytes is the floor for recommended Memory (32 MiB).
	MinRecommendedMemoryBytes = 32 * 1024 * 1024
)

// ContainerRecommendation encapsulates rightsizing suggestions for a single container.
type ContainerRecommendation struct {
	ContainerName string

	CurrentCPURequest float64
	CurrentCPULimit   float64
	CurrentCPUUsage   float64

	CurrentMemoryRequest int64
	CurrentMemoryLimit   int64
	CurrentMemoryUsage   int64

	RecommendedCPURequest float64
	RecommendedCPULimit   float64

	RecommendedMemoryRequest int64
	RecommendedMemoryLimit   int64

	CurrentMonthlyCost     float64
	RecommendedMonthlyCost float64
	MonthlySavings         float64
	SavingsPercentage      float64
}

// WorkloadRecommendation encapsulates complete rightsizing suggestions for a workload.
type WorkloadRecommendation struct {
	Namespace              string
	WorkloadName           string
	WorkloadKind           string
	Replicas               int32
	Containers             []ContainerRecommendation
	CurrentMonthlyCost     float64
	RecommendedMonthlyCost float64
	MonthlySavings         float64
	SavingsPercentage      float64
}

// Rightsizer computes container and workload rightsizing suggestions based on usage and headroom.
type Rightsizer struct {
	calculator      *CostCalculator
	headroomMargin  float64
	limitMultiplier float64
}

// NewRightsizer initializes a rightsizer with the given calculator and tuning parameters.
func NewRightsizer(calculator *CostCalculator, headroomMargin, limitMultiplier float64) *Rightsizer {
	if headroomMargin <= 0 {
		headroomMargin = DefaultHeadroomMargin
	}
	if limitMultiplier <= 0 {
		limitMultiplier = DefaultLimitMultiplier
	}
	return &Rightsizer{
		calculator:      calculator,
		headroomMargin:  headroomMargin,
		limitMultiplier: limitMultiplier,
	}
}

// RecommendContainer evaluates a single container and returns rightsizing recommendations.
func (r *Rightsizer) RecommendContainer(
	containerName string,
	cpuReq, cpuLimit, cpuUsage float64,
	memReq, memLimit, memUsage int64,
) ContainerRecommendation {
	// 1. Recommended CPU Request = Actual Usage * Headroom
	recCPUReq := cpuUsage * r.headroomMargin
	if recCPUReq < MinRecommendedCPU {
		recCPUReq = MinRecommendedCPU
	}

	// 2. Recommended CPU Limit = Recommended Request * Multiplier
	recCPULimit := recCPUReq * r.limitMultiplier

	// 3. Recommended Memory Request = Actual Usage * Headroom
	recMemReq := int64(math.Ceil(float64(memUsage) * r.headroomMargin))
	if recMemReq < MinRecommendedMemoryBytes {
		recMemReq = MinRecommendedMemoryBytes
	}

	// 4. Recommended Memory Limit = Recommended Request * Multiplier
	recMemLimit := int64(math.Ceil(float64(recMemReq) * r.limitMultiplier))

	// Cost projections
	currCost := r.calculator.MonthlyCost(cpuReq, memReq)
	recCost := r.calculator.MonthlyCost(recCPUReq, recMemReq)
	savings := currCost - recCost
	if savings < 0 {
		savings = 0
	}

	var savingsPct float64
	if currCost > 0 {
		savingsPct = (savings / currCost) * 100.0
	}

	return ContainerRecommendation{
		ContainerName:            containerName,
		CurrentCPURequest:        cpuReq,
		CurrentCPULimit:          cpuLimit,
		CurrentCPUUsage:          cpuUsage,
		CurrentMemoryRequest:     memReq,
		CurrentMemoryLimit:       memLimit,
		CurrentMemoryUsage:       memUsage,
		RecommendedCPURequest:    math.Round(recCPUReq*1000) / 1000,
		RecommendedCPULimit:      math.Round(recCPULimit*1000) / 1000,
		RecommendedMemoryRequest: recMemReq,
		RecommendedMemoryLimit:   recMemLimit,
		CurrentMonthlyCost:       math.Round(currCost*100) / 100,
		RecommendedMonthlyCost:   math.Round(recCost*100) / 100,
		MonthlySavings:           math.Round(savings*100) / 100,
		SavingsPercentage:        math.Round(savingsPct*10) / 10,
	}
}

// RecommendWorkload aggregates recommendations across all containers in a workload multiplied by replica count.
func (r *Rightsizer) RecommendWorkload(
	namespace, workloadName, workloadKind string,
	replicas int32,
	containers []ContainerRecommendation,
) WorkloadRecommendation {
	if replicas <= 0 {
		replicas = 1
	}

	var totalCurrCost, totalRecCost float64
	for _, c := range containers {
		totalCurrCost += c.CurrentMonthlyCost * float64(replicas)
		totalRecCost += c.RecommendedMonthlyCost * float64(replicas)
	}

	savings := totalCurrCost - totalRecCost
	if savings < 0 {
		savings = 0
	}

	var savingsPct float64
	if totalCurrCost > 0 {
		savingsPct = (savings / totalCurrCost) * 100.0
	}

	return WorkloadRecommendation{
		Namespace:              namespace,
		WorkloadName:           workloadName,
		WorkloadKind:           workloadKind,
		Replicas:               replicas,
		Containers:             containers,
		CurrentMonthlyCost:     math.Round(totalCurrCost*100) / 100,
		RecommendedMonthlyCost: math.Round(totalRecCost*100) / 100,
		MonthlySavings:         math.Round(savings*100) / 100,
		SavingsPercentage:      math.Round(savingsPct*10) / 10,
	}
}
