package views

import (
	"strings"
	"testing"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
)

func TestBuildColumnsAndRows(t *testing.T) {
	// Workloads
	wCols := BuildWorkloadsColumns(140)
	if len(wCols) == 0 {
		t.Fatal("expected non-empty workload columns")
	}

	workloads := []k8s.WorkloadCostData{
		{
			Namespace:          "default",
			Name:               "payment-api",
			Kind:               "Deployment",
			Replicas:           3,
			ReadyReplicas:      3,
			TotalCPURequest:    3.0,
			TotalCPUUsage:      0.5,
			TotalMemoryRequest: 6 * 1024 * 1024 * 1024,
			TotalMemoryUsage:   1 * 1024 * 1024 * 1024,
			MonthlyCost:        120.0,
			MonthlyWaste:       90.0,
			WastePercentage:    75.0,
			EfficiencyScore:    25.0,
		},
	}

	wRows, filtered := BuildWorkloadsRows(workloads, "payment")
	if len(wRows) != 1 || len(filtered) != 1 {
		t.Fatalf("expected 1 row, got %d", len(wRows))
	}
	if wRows[0][0] != "default" || wRows[0][1] != "payment-api" {
		t.Errorf("unexpected row data: %+v", wRows[0])
	}

	// Filter out
	_, empty := BuildWorkloadsRows(workloads, "non-existent")
	if len(empty) != 0 {
		t.Errorf("expected 0 rows for non-existent filter, got %d", len(empty))
	}

	// Namespaces
	nsCols := BuildNamespacesColumns(120)
	if len(nsCols) == 0 {
		t.Fatal("expected namespace columns")
	}
	nsRows := BuildNamespacesRows([]k8s.NamespaceCostData{
		{
			Namespace:       "default",
			WorkloadCount:   1,
			PodCount:        3,
			MonthlyCost:     120.0,
			MonthlyWaste:    90.0,
			WastePercentage: 75.0,
		},
	}, "")
	if len(nsRows) != 1 {
		t.Fatalf("expected 1 namespace row, got %d", len(nsRows))
	}

	// Nodes
	nodeCols := BuildNodesColumns(120)
	if len(nodeCols) == 0 {
		t.Fatal("expected node columns")
	}
	nodeRows := BuildNodesRows([]k8s.NodeCostData{
		{
			Name:                    "node-1",
			Status:                  "Ready",
			Roles:                   "worker",
			MonthlyCost:             250.0,
			CPUAllocationPercent:    80.0,
			MemoryAllocationPercent: 75.0,
		},
	}, "")
	if len(nodeRows) != 1 {
		t.Fatalf("expected 1 node row, got %d", len(nodeRows))
	}
}

func TestModalRenderers(t *testing.T) {
	calc := finops.NewCostCalculator(finops.DefaultPricingCard())
	rightsizer := finops.NewRightsizer(calc, 1.20, 1.50)

	workload := &k8s.WorkloadCostData{
		Namespace:     "default",
		Name:          "payment-api",
		Kind:          "Deployment",
		Replicas:      3,
		ReadyReplicas: 3,
		Containers: []k8s.ContainerCostData{
			{
				Name:            "payment-api",
				CPURequest:      1.0,
				CPULimit:        2.0,
				CPUUsage:        0.2,
				MemoryRequest:   2 * 1024 * 1024 * 1024,
				MemoryLimit:     4 * 1024 * 1024 * 1024,
				MemoryUsage:     512 * 1024 * 1024,
				MonthlyCost:     40.0,
				MonthlyWaste:    30.0,
				WastePercentage: 75.0,
			},
		},
		MonthlyCost:     120.0,
		MonthlyWaste:    90.0,
		WastePercentage: 75.0,
		EfficiencyScore: 25.0,
	}

	details := RenderContainerDetailsModal(workload, rightsizer, 100)
	if !strings.Contains(details, "payment-api") || !strings.Contains(details, "RIGHTSITING RECOMMENDATIONS") {
		t.Errorf("unexpected container details output: %s", details)
	}

	patchModal := RenderPatchModal("payment-api", "default", "yaml: true", 100)
	if !strings.Contains(patchModal, "payment-api") || !strings.Contains(patchModal, "yaml: true") {
		t.Errorf("unexpected patch modal output: %s", patchModal)
	}
}
