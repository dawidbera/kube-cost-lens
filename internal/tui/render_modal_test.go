package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

func TestRenderAllModals(t *testing.T) {
	mock := k8s.NewMockClusterProvider()
	pricing := finops.DefaultPricingCard()
	m := NewModel(mock, pricing, 5*time.Second)
	m.Width = 100
	m.Height = 30

	// 1. Details modal
	m.ActiveModal = ModalContainerDetails
	m.Viewport.SetContent("container details text")
	outDetails := m.View()
	if !strings.Contains(outDetails, "CONTAINER DRILLDOWN") {
		t.Error("expected container drilldown header in modal view")
	}

	// 2. Patch modal
	m.ActiveModal = ModalPatchPreview
	m.Viewport.SetContent("patch yaml content")
	outPatch := m.View()
	if !strings.Contains(outPatch, "STRATEGIC MERGE PATCH") {
		t.Error("expected patch header in modal view")
	}

	// 3. Help modal
	m.ActiveModal = ModalHelp
	outHelp := m.View()
	if !strings.Contains(outHelp, "KEYBOARD NAVIGATION") {
		t.Error("expected help header in modal view")
	}
}
