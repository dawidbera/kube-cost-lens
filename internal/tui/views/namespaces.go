// Package views provides namespace summary table formatting.
package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

// BuildNamespacesColumns returns column definitions for the namespace summary table.
func BuildNamespacesColumns(width int) []table.Column {
	nsWidth := 20
	extra := width - 100
	if extra > 0 {
		nsWidth += extra / 2
	}

	return []table.Column{
		{Title: "NAMESPACE", Width: nsWidth},
		{Title: "WORKLOADS", Width: 12},
		{Title: "PODS", Width: 8},
		{Title: "CPU REQ/ACT", Width: 16},
		{Title: "RAM REQ/ACT", Width: 16},
		{Title: "MO. COST", Width: 14},
		{Title: "MO. WASTE", Width: 16},
	}
}

// BuildNamespacesRows converts namespace aggregates into table rows.
func BuildNamespacesRows(namespaces []k8s.NamespaceCostData, filter string) []table.Row {
	var rows []table.Row
	filter = strings.ToLower(strings.TrimSpace(filter))

	for _, ns := range namespaces {
		if filter != "" && !strings.Contains(strings.ToLower(ns.Namespace), filter) {
			continue
		}

		cpuStr := fmt.Sprintf("%s / %s", k8s.FormatCPU(ns.TotalCPURequest), k8s.FormatCPU(ns.TotalCPUUsage))
		ramStr := fmt.Sprintf("%s / %s", k8s.FormatMemoryBytes(ns.TotalMemoryRequest), k8s.FormatMemoryBytes(ns.TotalMemoryUsage))
		costStr := fmt.Sprintf("$%.2f", ns.MonthlyCost)
		wasteStr := fmt.Sprintf("$%.2f (%.0f%%)", ns.MonthlyWaste, ns.WastePercentage)

		rows = append(rows, table.Row{
			ns.Namespace,
			fmt.Sprintf("%d", ns.WorkloadCount),
			fmt.Sprintf("%d", ns.PodCount),
			cpuStr,
			ramStr,
			costStr,
			wasteStr,
		})
	}

	return rows
}
