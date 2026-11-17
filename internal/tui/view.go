// Package tui provides declarative string-based terminal rendering.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the complete ANSI terminal screen buffer.
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Initializing kube-cost-lens..."
	}

	var totalClusterCost, totalClusterWaste float64
	for _, ns := range m.Namespaces {
		totalClusterCost += ns.MonthlyCost
		totalClusterWaste += ns.MonthlyWaste
	}

	// 1. Header Bar
	headerTitle := AppHeaderStyle.Render("kube-cost-lens v1.0.0")
	contextTag := ContextTagStyle.Render(fmt.Sprintf("Context: %s", m.K8sProvider.GetContextName()))
	costTag := CostTagStyle.Render(fmt.Sprintf("Monthly: $%.2f (Waste: $%.2f)", totalClusterCost, totalClusterWaste))
	loadingTag := ""
	if m.IsLoading {
		loadingTag = TabInactiveStyle.Render("⟳ Loading...")
	}

	topBar := lipgloss.JoinHorizontal(lipgloss.Center, headerTitle, " ", contextTag, " ", costTag, " ", loadingTag)

	// 2. Tabs Bar
	tabWorkloads := TabInactiveStyle.Render("[1] Workloads")
	tabNamespaces := TabInactiveStyle.Render("[2] Namespaces")
	tabNodes := TabInactiveStyle.Render("[3] Nodes")

	switch m.ActiveView {
	case ViewWorkloads:
		tabWorkloads = TabActiveStyle.Render("[1] Workloads (Active)")
	case ViewNamespaces:
		tabNamespaces = TabActiveStyle.Render("[2] Namespaces (Active)")
	case ViewNodes:
		tabNodes = TabActiveStyle.Render("[3] Nodes (Active)")
	}

	filterBar := ""
	if m.IsFiltering {
		filterBar = m.FilterInput.View()
	} else if m.FilterQuery != "" {
		filterBar = FilterPromptStyle.Render(fmt.Sprintf("Filter: %q", m.FilterQuery))
	} else {
		filterBar = TabInactiveStyle.Render("[/] Filter")
	}

	tabsBar := lipgloss.JoinHorizontal(lipgloss.Center, tabWorkloads, "  ", tabNamespaces, "  ", tabNodes, "   |   ", filterBar)

	// 3. Main Center Area (Table or Modal)
	var mainArea string
	if m.ActiveModal != ModalNone {
		mainArea = m.renderModal()
	} else {
		mainArea = m.Table.View()
	}

	// 4. Footer Bar
	var footerInfo string
	if m.StatusMessage != "" {
		footerInfo = StatusMessageStyle.Render(m.StatusMessage)
	} else if m.ErrorMessage != "" {
		footerInfo = WasteHighStyle.Render("Error: " + m.ErrorMessage)
	} else if m.ActiveView == ViewWorkloads && m.SelectedWorkload != nil {
		w := m.SelectedWorkload
		footerInfo = fmt.Sprintf("Selected: %s/%s (%s) | Waste: $%.2f (%.0f%%) | Rec. Savings Available: [p]",
			w.Namespace, w.Name, w.Kind, w.MonthlyWaste, w.WastePercentage)
	} else {
		footerInfo = fmt.Sprintf("Showing %d items | Tab: switch view | [/]: filter", len(m.Table.Rows()))
	}

	keyHelp := FooterStyle.Render(
		"[Enter] Inspect | [p] Patch | [/] Filter | [Tab] Switch View | [r] Refresh | [?] Help | [q] Quit",
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		tabsBar,
		strings.Repeat("─", m.Width),
		mainArea,
		strings.Repeat("─", m.Width),
		footerInfo,
		keyHelp,
	)
}

// renderModal builds the overlay view for active dialogs.
func (m Model) renderModal() string {
	var title string
	var body string

	switch m.ActiveModal {
	case ModalContainerDetails:
		title = ModalTitleStyle.Render("CONTAINER DRILLDOWN & SIZING RECOMMENDATIONS")
		body = m.Viewport.View()

	case ModalPatchPreview:
		title = ModalTitleStyle.Render("STRATEGIC MERGE PATCH MANIFEST PREVIEW")
		body = m.Viewport.View()

	case ModalHelp:
		title = ModalTitleStyle.Render("KUBE-COST-LENS KEYBOARD NAVIGATION & HELP")
		body = `
Navigation & Controls:
  ↑ / k, ↓ / j     Move cursor up/down in table
  Tab              Cycle forward to next view tab
  Shift+Tab        Cycle backward to previous view tab
  1, 2, 3          Directly jump to Workloads, Namespaces, or Nodes view
  Enter            Open container drilldown modal for selected workload
  p                Open 1-click Rightsizing Patch modal
  /                Open search filter input (ESC to cancel)
  r                Force immediate cluster refresh
  ?                Toggle this help screen
  q, Ctrl+C        Exit application

Modal Actions:
  c                Copy generated patch YAML to clipboard
  w                Save patch YAML to local disk (patch-<name>.yaml)
  Esc              Close modal overlay and return to table

FinOps Metrics:
  CPU Rate:        $0.0416/vCPU-hour (~$30/vCPU-month)
  RAM Rate:        $0.0055/GB-hour   (~$4/GB-month)
  Idle Waste:      max(0, Requested - Used) * Unit Rate * 730h
  Recommended Req: Actual Usage * 1.20 (20% safety headroom)
  Recommended Lim: Recommended Request * 1.50
`
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, body)
	modalWidth := m.Width - 8
	if modalWidth < 40 {
		modalWidth = 40
	}
	return ModalBoxStyle.Width(modalWidth).Render(content)
}
