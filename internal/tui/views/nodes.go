// Package views provides node inventory and allocation table formatting.
package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

// BuildNodesColumns returns column definitions for node view.
func BuildNodesColumns(width int) []table.Column {
	nameWidth := 24
	extra := width - 110
	if extra > 0 {
		nameWidth += extra / 2
	}

	return []table.Column{
		{Title: "NODE NAME", Width: nameWidth},
		{Title: "STATUS", Width: 10},
		{Title: "ROLES", Width: 14},
		{Title: "CPU (REQ/ALLOC)", Width: 18},
		{Title: "CPU ALLOC %", Width: 14},
		{Title: "RAM (REQ/ALLOC)", Width: 18},
		{Title: "RAM ALLOC %", Width: 14},
		{Title: "EST. MO. COST", Width: 14},
	}
}

// BuildNodesRows converts node metrics into table rows.
func BuildNodesRows(nodes []k8s.NodeCostData, filter string) []table.Row {
	var rows []table.Row
	filter = strings.ToLower(strings.TrimSpace(filter))

	for _, n := range nodes {
		if filter != "" && !strings.Contains(strings.ToLower(n.Name), filter) && !strings.Contains(strings.ToLower(n.Roles), filter) {
			continue
		}

		cpuStr := fmt.Sprintf("%s / %s", k8s.FormatCPU(n.CPURequests), k8s.FormatCPU(n.CPUAllocatable))
		ramStr := fmt.Sprintf("%s / %s", k8s.FormatMemoryBytes(n.MemoryRequests), k8s.FormatMemoryBytes(n.MemoryAllocatable))
		costStr := fmt.Sprintf("$%.2f", n.MonthlyCost)

		rows = append(rows, table.Row{
			n.Name,
			n.Status,
			n.Roles,
			cpuStr,
			fmt.Sprintf("%.1f%%", n.CPUAllocationPercent),
			ramStr,
			fmt.Sprintf("%.1f%%", n.MemoryAllocationPercent),
			costStr,
		})
	}

	return rows
}
