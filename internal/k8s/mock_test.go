package k8s

import (
	"context"
	"testing"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
)

func TestMockClusterProvider(t *testing.T) {
	provider := NewMockClusterProvider()
	if provider.GetContextName() == "" {
		t.Error("expected non-empty context name")
	}

	calc := finops.NewCostCalculator(finops.DefaultPricingCard())
	ctx := context.Background()

	// 1. Fetch all workloads
	workloads, namespaces, nodes, err := provider.FetchClusterData(ctx, "", calc)
	if err != nil {
		t.Fatalf("unexpected error fetching cluster data: %v", err)
	}

	if len(workloads) == 0 {
		t.Fatal("expected synthetic workloads, got 0")
	}
	if len(namespaces) == 0 {
		t.Fatal("expected namespaces, got 0")
	}
	if len(nodes) == 0 {
		t.Fatal("expected nodes, got 0")
	}

	// Verify workload calculations
	foundPayment := false
	for _, w := range workloads {
		if w.Name == "payment-api" {
			foundPayment = true
			if w.Replicas != 3 {
				t.Errorf("expected 3 replicas for payment-api, got %d", w.Replicas)
			}
			if w.MonthlyWaste <= 0 {
				t.Errorf("expected positive monthly waste for payment-api, got %f", w.MonthlyWaste)
			}
			if w.WastePercentage < 50.0 {
				t.Errorf("payment-api should have high waste percentage, got %f", w.WastePercentage)
			}
		}
	}
	if !foundPayment {
		t.Error("payment-api workload not found in mock data")
	}

	// 2. Fetch with namespace filter
	filteredWorkloads, _, _, err := provider.FetchClusterData(ctx, "monitoring", calc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, w := range filteredWorkloads {
		if w.Namespace != "monitoring" {
			t.Errorf("expected workload in namespace monitoring, got %s", w.Namespace)
		}
	}

	if err := provider.Close(); err != nil {
		t.Errorf("error closing provider: %v", err)
	}
}

func TestQuantityParsing(t *testing.T) {
	cpu, err := ParseCPUQuantity("500m")
	if err != nil || cpu != 0.5 {
		t.Errorf("expected 0.5 cores, got %f (err: %v)", cpu, err)
	}

	cpu2, err := ParseCPUQuantity("2")
	if err != nil || cpu2 != 2.0 {
		t.Errorf("expected 2.0 cores, got %f (err: %v)", cpu2, err)
	}

	mem, err := ParseMemoryQuantity("512Mi")
	if err != nil || mem != 512*1024*1024 {
		t.Errorf("expected 512MiB in bytes, got %d (err: %v)", mem, err)
	}

	mem2, err := ParseMemoryQuantity("2Gi")
	if err != nil || mem2 != 2*1024*1024*1024 {
		t.Errorf("expected 2GiB in bytes, got %d (err: %v)", mem2, err)
	}

	// Format helpers
	if FormatCPU(0.5) != "500m" {
		t.Errorf("expected 500m, got %s", FormatCPU(0.5))
	}
	if FormatCPU(2.0) != "2.00" {
		t.Errorf("expected 2.00, got %s", FormatCPU(2.0))
	}

	if FormatMemoryBytes(512*1024*1024) != "512Mi" {
		t.Errorf("expected 512Mi, got %s", FormatMemoryBytes(512*1024*1024))
	}
	if FormatMemoryBytes(2*1024*1024*1024) != "2Gi" {
		t.Errorf("expected 2Gi, got %s", FormatMemoryBytes(2*1024*1024*1024))
	}
}
