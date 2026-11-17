// Package tui defines the Elm Architecture application model and message channels.
package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

// ViewMode defines the primary active view tab.
type ViewMode int

const (
	// ViewWorkloads displays controller workloads (Deployments, StatefulSets, DaemonSets).
	ViewWorkloads ViewMode = iota
	// ViewNamespaces displays cost aggregate summaries by namespace.
	ViewNamespaces
	// ViewNodes displays cluster node capacities and allocation rates.
	ViewNodes
)

// ModalMode defines which modal dialog is currently presented to the operator.
type ModalMode int

const (
	// ModalNone indicates no modal is active.
	ModalNone ModalMode = iota
	// ModalContainerDetails displays container-level drilldown.
	ModalContainerDetails
	// ModalPatchPreview displays generated YAML rightsizing patch.
	ModalPatchPreview
	// ModalHelp displays keyboard controls and legend.
	ModalHelp
)

// TickMsg is delivered on periodic refresh intervals.
type TickMsg time.Time

// ClusterDataUpdatedMsg delivers freshly queried cluster inventory and calculations.
type ClusterDataUpdatedMsg struct {
	Workloads  []k8s.WorkloadCostData
	Namespaces []k8s.NamespaceCostData
	Nodes      []k8s.NodeCostData
	Err        error
}

// ClipboardCopiedMsg reports the result of copying patch data to the system clipboard.
type ClipboardCopiedMsg struct {
	Success bool
	Message string
}

// PatchSavedMsg reports the result of exporting patch manifests to disk.
type PatchSavedMsg struct {
	Path string
	Err  error
}

// Model maintains centralized UI state following The Elm Architecture.
type Model struct {
	// Terminal Dimensions
	Width  int
	Height int

	// Navigation & Mode
	ActiveView  ViewMode
	ActiveModal ModalMode
	KeyMap      KeyMap

	// Filtering & Search
	SelectedNS  string
	FilterQuery string
	IsFiltering bool
	FilterInput textinput.Model

	// Status & Loading
	IsLoading     bool
	ErrorMessage  string
	StatusMessage string

	// UI Components
	Table    table.Model
	Viewport viewport.Model

	// Domain & Cluster Data
	Workloads         []k8s.WorkloadCostData
	FilteredWorkloads []k8s.WorkloadCostData
	Namespaces        []k8s.NamespaceCostData
	Nodes             []k8s.NodeCostData
	PricingCard       finops.PricingCard
	Calculator        *finops.CostCalculator
	Rightsizer        *finops.Rightsizer

	// Selected Context
	SelectedWorkload *k8s.WorkloadCostData
	GeneratedPatch   string

	// Provider & Execution
	K8sProvider     k8s.ClusterProvider
	RefreshInterval time.Duration
}

// NewModel constructs an initialized application model with default components.
func NewModel(provider k8s.ClusterProvider, pricing finops.PricingCard, refreshInterval time.Duration) Model {
	if refreshInterval <= 0 {
		refreshInterval = 10 * time.Second
	}

	calc := finops.NewCostCalculator(pricing)
	sizer := finops.NewRightsizer(calc, finops.DefaultHeadroomMargin, finops.DefaultLimitMultiplier)

	ti := textinput.New()
	ti.Placeholder = "type to filter workloads or namespaces..."
	ti.Prompt = "Filter: "
	ti.PromptStyle = FilterPromptStyle
	ti.CharLimit = 64

	t := table.New(
		table.WithColumns([]table.Column{}),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	t.SetStyles(TableDefaultStyles())

	vp := viewport.New(80, 20)

	m := Model{
		ActiveView:      ViewWorkloads,
		ActiveModal:     ModalNone,
		KeyMap:          DefaultKeyMap(),
		FilterInput:     ti,
		IsLoading:       true,
		Table:           t,
		Viewport:        vp,
		PricingCard:     pricing,
		Calculator:      calc,
		Rightsizer:      sizer,
		K8sProvider:     provider,
		RefreshInterval: refreshInterval,
	}

	return m
}

// Init starts initial data fetch and ticker command.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		FetchClusterDataCmd(m.K8sProvider, m.SelectedNS, m.Calculator),
		TickCmd(m.RefreshInterval),
	)
}

// TickCmd schedules the next periodic refresh tick.
func TickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// FetchClusterDataCmd dispatches asynchronous cluster queries.
func FetchClusterDataCmd(provider k8s.ClusterProvider, nsFilter string, calc *finops.CostCalculator) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		workloads, namespaces, nodes, err := provider.FetchClusterData(ctx, nsFilter, calc)
		return ClusterDataUpdatedMsg{
			Workloads:  workloads,
			Namespaces: namespaces,
			Nodes:      nodes,
			Err:        err,
		}
	}
}
