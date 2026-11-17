// Package k8s contains helper functions to normalize workload and container specifications.
package k8s

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
)

// buildWorkloadData constructs a WorkloadCostData from container specs and replicas.
func buildWorkloadData(
	namespace, name, kind string,
	replicas, readyReplicas int32,
	containers []corev1.Container,
	calc *finops.CostCalculator,
) WorkloadCostData {
	if replicas <= 0 {
		replicas = 1
	}

	var containerData []ContainerCostData
	var totReqCPU, totLimitCPU, totUsedCPU float64
	var totReqMem, totLimitMem, totUsedMem int64

	for _, c := range containers {
		cpuReq := float64(c.Resources.Requests.Cpu().MilliValue()) / 1000.0
		cpuLim := float64(c.Resources.Limits.Cpu().MilliValue()) / 1000.0
		memReq := c.Resources.Requests.Memory().Value()
		memLim := c.Resources.Limits.Memory().Value()

		// Fallback defaults if not set in manifest
		if cpuReq <= 0 {
			cpuReq = 0.100 // 100m default
		}
		if memReq <= 0 {
			memReq = 128 * 1024 * 1024 // 128 MiB default
		}

		// Sensible utilization estimate (e.g., 30-50% of request)
		cpuUsed := cpuReq * 0.35
		memUsed := int64(float64(memReq) * 0.40)

		cWaste := calc.CalculateWaste(cpuReq, cpuUsed, memReq, memUsed)

		cData := ContainerCostData{
			Name:            c.Name,
			CPURequest:      cpuReq,
			CPULimit:        cpuLim,
			CPUUsage:        cpuUsed,
			MemoryRequest:   memReq,
			MemoryLimit:     memLim,
			MemoryUsage:     memUsed,
			MonthlyCost:     cWaste.TotalCost,
			MonthlyWaste:    cWaste.TotalWaste,
			WastePercentage: cWaste.WastePercentage,
		}

		containerData = append(containerData, cData)

		totReqCPU += cpuReq
		totLimitCPU += cpuLim
		totUsedCPU += cpuUsed
		totReqMem += memReq
		totLimitMem += memLim
		totUsedMem += memUsed
	}

	repl := float64(replicas)
	totalCPUReq := totReqCPU * repl
	totalCPUUsed := totUsedCPU * repl
	totalMemReq := totReqMem * int64(replicas)
	totalMemUsed := totUsedMem * int64(replicas)

	wWaste := calc.CalculateWaste(totalCPUReq, totalCPUUsed, totalMemReq, totalMemUsed)

	return WorkloadCostData{
		Namespace:          namespace,
		Name:               name,
		Kind:               kind,
		Replicas:           replicas,
		ReadyReplicas:      readyReplicas,
		Containers:         containerData,
		TotalCPURequest:    totalCPUReq,
		TotalCPUUsage:      totalCPUUsed,
		TotalMemoryRequest: totalMemReq,
		TotalMemoryUsage:   totalMemUsed,
		MonthlyCost:        wWaste.TotalCost,
		MonthlyWaste:       wWaste.TotalWaste,
		WastePercentage:    wWaste.WastePercentage,
		EfficiencyScore:    wWaste.EfficiencyScore,
	}
}
