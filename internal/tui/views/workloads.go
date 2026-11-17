// Package views provides rendering components for individual tabs and dialogs.
package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

// BuildWorkloadsColumns returns column definitions for the workload table.
func BuildWorkloadsColumns(width int) []table.Column {
	// Base widths
	nsWidth := 14
	nameWidth := 20
	kindWidth := 12
	replWidth := 6
	cpuWidth := 14
	ramWidth := 14
	costWidth := 12
	wasteWidth := 14

	// Dynamically expand name if terminal is wide
	extra := width - 110
	if extra > 0 {
		nameWidth += extra / 2
		nsWidth += extra / 4
	}

	return []table.Column{
		{Title: "NAMESPACE", Width: nsWidth},
		{Title: "WORKLOAD", Width: nameWidth},
		{Title: "KIND", Width: kindWidth},
		{Title: "REPL", Width: replWidth},
		{Title: "CPU REQ/ACT", Width: cpuWidth},
		{Title: "RAM REQ/ACT", Width: ramWidth},
		{Title: "MO. COST", Width: costWidth},
		{Title: "MO. WASTE", Width: wasteWidth},
	}
}

// BuildWorkloadsRows converts workload data into table rows, applying filter if present.
func BuildWorkloadsRows(workloads []k8s.WorkloadCostData, filter string) ([]table.Row, []k8s.WorkloadCostData) {
	var rows []table.Row
	var filtered []k8s.WorkloadCostData

	filter = strings.ToLower(strings.TrimSpace(filter))

	for _, w := range workloads {
		if filter != "" {
			nsMatch := strings.Contains(strings.ToLower(w.Namespace), filter)
			nameMatch := strings.Contains(strings.ToLower(w.Name), filter)
			kindMatch := strings.Contains(strings.ToLower(w.Kind), filter)
			if !nsMatch && !nameMatch && !kindMatch {
				continue
			}
		}

		filtered = append(filtered, w)

		cpuStr := fmt.Sprintf("%s / %s", k8s.FormatCPU(w.TotalCPURequest), k8s.FormatCPU(w.TotalCPUUsage))
		ramStr := fmt.Sprintf("%s / %s", k8s.FormatMemoryBytes(w.TotalMemoryRequest), k8s.FormatMemoryBytes(w.TotalMemoryUsage))
		costStr := fmt.Sprintf("$%.2f", w.MonthlyCost)

		wasteIndicator := ""
		if w.WastePercentage >= 70.0 {
			wasteIndicator = " ⚠"
		}
		wasteStr := fmt.Sprintf("$%.2f (%.0f%%)%s", w.MonthlyWaste, w.WastePercentage, wasteIndicator)

		row := table.Row{
			w.Namespace,
			w.Name,
			w.Kind,
			fmt.Sprintf("%d/%d", w.ReadyReplicas, w.Replicas),
			cpuStr,
			ramStr,
			costStr,
			wasteStr,
		}
		rows = append(rows, row)
	}

	return rows, filtered
}
