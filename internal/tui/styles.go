// Package tui provides Lipgloss visual styling and color themes.
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Terminal color palette definitions.
const (
	ColorPrimary   = lipgloss.Color("#7D56F4") // Deep Purple
	ColorSecondary = lipgloss.Color("#04B575") // Emerald Green
	ColorWarning   = lipgloss.Color("#FFB86C") // Warm Amber
	ColorDanger    = lipgloss.Color("#FF5555") // Crimson Red
	ColorMuted     = lipgloss.Color("#6272A4") // Dracula Slate
	ColorBgModal   = lipgloss.Color("#21222C") // Deep Charcoal
	ColorWhite     = lipgloss.Color("#F8F8F2") // Off-white
	ColorGray      = lipgloss.Color("#44475A") // Dark Slate Gray
)

// Lipgloss styles for structural components.
var (
	// AppHeaderStyle formats the top title bar.
	AppHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorPrimary).
			Padding(0, 1)

	// ContextTagStyle formats the active cluster context name.
	ContextTagStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorMuted).
			Padding(0, 1)

	// CostTagStyle formats the aggregated monthly cluster burn rate.
	CostTagStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Padding(0, 1)

	// TabActiveStyle highlights the currently selected view tab.
	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorPrimary).
			Padding(0, 1)

	// TabInactiveStyle styles unselected view tabs.
	TabInactiveStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 1)

	// FooterStyle styles the bottom helper bar.
	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorGray).
			Padding(0, 1)

	// StatusMessageStyle styles transient status announcements.
	StatusMessageStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Padding(0, 1)

	// ModalBoxStyle styles floating drilldown and patch modal containers.
	ModalBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Background(ColorBgModal).
			Padding(1, 2)

	// ModalTitleStyle formats modal headings.
	ModalTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	// FilterPromptStyle formats the live search input prompt.
	FilterPromptStyle = lipgloss.NewStyle().
				Foreground(ColorWarning).
				Bold(true)

	// WasteHighStyle formats critical waste alerts (>70%).
	WasteHighStyle = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	// WasteMediumStyle formats moderate waste warnings (40-70%).
	WasteMediumStyle = lipgloss.NewStyle().
				Foreground(ColorWarning)

	// WasteLowStyle formats healthy/low waste metrics (<40%).
	WasteLowStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)
)

// TableDefaultStyles returns configured Lipgloss styles for bubbles table components.
func TableDefaultStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ColorMuted).
		BorderBottom(true).
		Bold(true).
		Foreground(ColorWhite)
	s.Selected = s.Selected.
		Foreground(ColorWhite).
		Background(ColorPrimary).
		Bold(true)
	return s
}
