package tui

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

func TestTUIModelLifecycle(t *testing.T) {
	mock := k8s.NewMockClusterProvider()
	pricing := finops.DefaultPricingCard()
	m := NewModel(mock, pricing, 5*time.Second)

	// 1. Terminal resize
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)

	if m.Width != 120 || m.Height != 40 {
		t.Errorf("expected size 120x40, got %dx%d", m.Width, m.Height)
	}

	// 2. Cluster data arrival
	workloads, namespaces, nodes, err := mock.FetchClusterData(context.Background(), "", m.Calculator)
	if err != nil {
		t.Fatalf("failed fetching mock data: %v", err)
	}

	updated, _ = m.Update(ClusterDataUpdatedMsg{
		Workloads:  workloads,
		Namespaces: namespaces,
		Nodes:      nodes,
		Err:        nil,
	})
	m = updated.(Model)

	if len(m.Workloads) == 0 {
		t.Fatal("expected populated workloads")
	}
	if len(m.Table.Rows()) == 0 {
		t.Fatal("expected table rows populated")
	}

	// 3. Tab cycling
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.ActiveView != ViewNamespaces {
		t.Errorf("expected ViewNamespaces, got %v", m.ActiveView)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.ActiveView != ViewNodes {
		t.Errorf("expected ViewNodes, got %v", m.ActiveView)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.ActiveView != ViewWorkloads {
		t.Errorf("expected back to ViewWorkloads, got %v", m.ActiveView)
	}

	// 4. Modal inspection
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.ActiveModal != ModalContainerDetails {
		t.Errorf("expected ModalContainerDetails, got %v", m.ActiveModal)
	}

	// 5. Open patch preview from container modal
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = updated.(Model)
	if m.ActiveModal != ModalPatchPreview {
		t.Errorf("expected ModalPatchPreview, got %v", m.ActiveModal)
	}
	if m.GeneratedPatch == "" {
		t.Error("expected generated patch YAML to be non-empty")
	}

	// 6. Close modal
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if m.ActiveModal != ModalNone {
		t.Errorf("expected ModalNone, got %v", m.ActiveModal)
	}

	// 7. View output check
	viewStr := m.View()
	if !strings.Contains(viewStr, "kube-cost-lens") {
		t.Error("rendered view does not contain title")
	}
	if !strings.Contains(viewStr, "payment-api") {
		t.Error("rendered view does not contain workload row")
	}
}

func TestTUIFiltering(t *testing.T) {
	mock := k8s.NewMockClusterProvider()
	pricing := finops.DefaultPricingCard()
	m := NewModel(mock, pricing, 5*time.Second)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)

	workloads, namespaces, nodes, _ := mock.FetchClusterData(context.Background(), "", m.Calculator)
	updated, _ = m.Update(ClusterDataUpdatedMsg{
		Workloads:  workloads,
		Namespaces: namespaces,
		Nodes:      nodes,
	})
	m = updated.(Model)

	initialRows := len(m.Table.Rows())

	// Enter filter mode
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = updated.(Model)
	if !m.IsFiltering {
		t.Error("expected IsFiltering true")
	}

	// Type filter for "prometheus"
	for _, r := range "prometheus" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = updated.(Model)
	}

	if len(m.Table.Rows()) >= initialRows {
		t.Errorf("filter should reduce rows, initial %d, current %d", initialRows, len(m.Table.Rows()))
	}
	if len(m.Table.Rows()) != 1 {
		t.Errorf("expected exactly 1 row for prometheus, got %d", len(m.Table.Rows()))
	}

	// Escape clears filter
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if len(m.Table.Rows()) != initialRows {
		t.Errorf("expected rows restored to %d, got %d", initialRows, len(m.Table.Rows()))
	}
}
