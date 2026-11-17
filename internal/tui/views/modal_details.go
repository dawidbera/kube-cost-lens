// Package views provides container-level inspection modal formatting.
package views

import (
	"fmt"
	"strings"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

// RenderContainerDetailsModal formats full container-level breakdown for a selected workload.
func RenderContainerDetailsModal(
	workload *k8s.WorkloadCostData,
	rightsizer *finops.Rightsizer,
	width int,
) string {
	if workload == nil {
		return "No workload selected."
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Workload: %s/%s (%s)\n", workload.Namespace, workload.Name, workload.Kind))
	sb.WriteString(fmt.Sprintf("Replicas: %d ready / %d total\n", workload.ReadyReplicas, workload.Replicas))
	sb.WriteString(fmt.Sprintf("Financials: Monthly Cost: $%.2f | Waste: $%.2f (%.1f%%) | Efficiency: %.1f%%\n\n",
		workload.MonthlyCost, workload.MonthlyWaste, workload.WastePercentage, workload.EfficiencyScore))

	sb.WriteString("--------------------------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("%-20s %-16s %-16s %-14s %-14s\n", "CONTAINER", "CPU (REQ/ACT)", "RAM (REQ/ACT)", "MO. COST", "MO. WASTE"))
	sb.WriteString("--------------------------------------------------------------------------------\n")

	var containerRecs []finops.ContainerRecommendation

	for _, c := range workload.Containers {
		cpuStr := fmt.Sprintf("%s / %s", k8s.FormatCPU(c.CPURequest), k8s.FormatCPU(c.CPUUsage))
		ramStr := fmt.Sprintf("%s / %s", k8s.FormatMemoryBytes(c.MemoryRequest), k8s.FormatMemoryBytes(c.MemoryUsage))
		costStr := fmt.Sprintf("$%.2f", c.MonthlyCost)
		wasteStr := fmt.Sprintf("$%.2f (%.0f%%)", c.MonthlyWaste, c.WastePercentage)

		sb.WriteString(fmt.Sprintf("%-20s %-16s %-16s %-14s %-14s\n",
			c.Name, cpuStr, ramStr, costStr, wasteStr))

		rec := rightsizer.RecommendContainer(c.Name, c.CPURequest, c.CPULimit, c.CPUUsage, c.MemoryRequest, c.MemoryLimit, c.MemoryUsage)
		containerRecs = append(containerRecs, rec)
	}

	wRec := rightsizer.RecommendWorkload(workload.Namespace, workload.Name, workload.Kind, workload.Replicas, containerRecs)

	sb.WriteString("\n=== RIGHTSITING RECOMMENDATIONS ===\n")
	for _, rec := range containerRecs {
		sb.WriteString(fmt.Sprintf("\nContainer [%s]:\n", rec.ContainerName))
		sb.WriteString(fmt.Sprintf("  • CPU Request: %s -> %s (Limit: %s)\n",
			k8s.FormatCPU(rec.CurrentCPURequest),
			k8s.FormatCPU(rec.RecommendedCPURequest),
			k8s.FormatCPU(rec.RecommendedCPULimit)))
		sb.WriteString(fmt.Sprintf("  • RAM Request: %s -> %s (Limit: %s)\n",
			k8s.FormatMemoryBytes(rec.CurrentMemoryRequest),
			k8s.FormatMemoryBytes(rec.RecommendedMemoryRequest),
			k8s.FormatMemoryBytes(rec.RecommendedMemoryLimit)))
		sb.WriteString(fmt.Sprintf("  • Monthly Savings: $%.2f (-%.1f%%)\n",
			rec.MonthlySavings, rec.SavingsPercentage))
	}

	sb.WriteString(fmt.Sprintf("\nWorkload Projected Monthly Savings: $%.2f (-%.1f%%)\n",
		wRec.MonthlySavings, wRec.SavingsPercentage))
	sb.WriteString("\n[p] Generate Patch | [Esc] Close Modal\n")

	return sb.String()
}
