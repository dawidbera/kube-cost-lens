// Package tui provides terminal program lifecycle and initialization.
package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

// AppOptions contains startup configuration for the TUI engine.
type AppOptions struct {
	Provider        k8s.ClusterProvider
	PricingCard     finops.PricingCard
	NamespaceFilter string
	RefreshInterval time.Duration
}

// RunApp initializes Bubbletea with alternate screen buffer and begins event loop.
func RunApp(opts AppOptions) error {
	model := NewModel(opts.Provider, opts.PricingCard, opts.RefreshInterval)
	if opts.NamespaceFilter != "" {
		model.SelectedNS = opts.NamespaceFilter
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running terminal UI: %w", err)
	}

	return nil
}
