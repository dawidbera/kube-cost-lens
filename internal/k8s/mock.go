// Package k8s provides synthetic mock implementations for offline demo mode.
package k8s

import (
	"context"
	"math"
	"sort"
	"strings"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
)

// MockClusterProvider implements ClusterProvider for demo and offline test environments.
type MockClusterProvider struct {
	contextName string
}

// NewMockClusterProvider constructs a new mock provider.
func NewMockClusterProvider() *MockClusterProvider {
	return &MockClusterProvider{
		contextName: "demo-cluster-us-east-1",
	}
}

// GetContextName returns the simulated cluster context identifier.
func (m *MockClusterProvider) GetContextName() string {
	return m.contextName
}

// Close simulates releasing cluster resources.
func (m *MockClusterProvider) Close() error {
	return nil
}

// FetchClusterData generates synthetic workloads, aggregates namespace stats, and builds node inventory.
func (m *MockClusterProvider) FetchClusterData(
	_ context.Context,
	namespaceFilter string,
	calc *finops.CostCalculator,
) ([]WorkloadCostData, []NamespaceCostData, []NodeCostData, error) {
	if calc == nil {
		calc = finops.NewCostCalculator(finops.DefaultPricingCard())
	}

	rawWorkloads := []struct {
		namespace string
		name      string
		kind      string
		replicas  int32
		ready     int32
		cName     string
		cpuReq    float64
		cpuLimit  float64
		cpuUsed   float64
		memReqGB  float64
		memLimGB  float64
		memUsedGB float64
	}{
		{
			namespace: "default",
			name:      "payment-api",
			kind:      "Deployment",
			replicas:  3,
			ready:     3,
			cName:     "payment-api",
			cpuReq:    3.0,
			cpuLimit:  4.5,
			cpuUsed:   0.45,
			memReqGB:  6.0,
			memLimGB:  9.0,
			memUsedGB: 1.20,
		},
		{
			namespace: "monitoring",
			name:      "prometheus",
			kind:      "StatefulSet",
			replicas:  1,
			ready:     1,
			cName:     "prometheus-server",
			cpuReq:    2.0,
			cpuLimit:  3.0,
			cpuUsed:   1.80,
			memReqGB:  8.0,
			memLimGB:  12.0,
			memUsedGB: 7.10,
		},
		{
			namespace: "ingress-gw",
			name:      "envoy-proxy",
			kind:      "Deployment",
			replicas:  2,
			ready:     2,
			cName:     "envoy",
			cpuReq:    1.0,
			cpuLimit:  2.0,
			cpuUsed:   0.85,
			memReqGB:  2.0,
			memLimGB:  3.0,
			memUsedGB: 1.40,
		},
		{
			namespace: "analytics",
			name:      "clickhouse",
			kind:      "StatefulSet",
			replicas:  2,
			ready:     2,
			cName:     "clickhouse-server",
			cpuReq:    8.0,
			cpuLimit:  12.0,
			cpuUsed:   1.50,
			memReqGB:  32.0,
			memLimGB:  48.0,
			memUsedGB: 8.50,
		},
		{
			namespace: "auth",
			name:      "keycloak",
			kind:      "Deployment",
			replicas:  2,
			ready:     2,
			cName:     "keycloak",
			cpuReq:    2.0,
			cpuLimit:  3.0,
			cpuUsed:   0.60,
			memReqGB:  4.0,
			memLimGB:  6.0,
			memUsedGB: 1.80,
		},
		{
			namespace: "kube-system",
			name:      "cilium",
			kind:      "DaemonSet",
			replicas:  3,
			ready:     3,
			cName:     "cilium-agent",
			cpuReq:    0.50,
			cpuLimit:  1.0,
			cpuUsed:   0.35,
			memReqGB:  0.50,
			memLimGB:  1.0,
			memUsedGB: 0.38,
		},
		{
			namespace: "frontend",
			name:      "web-client",
			kind:      "Deployment",
			replicas:  4,
			ready:     4,
			cName:     "nginx-static",
			cpuReq:    1.0,
			cpuLimit:  1.5,
			cpuUsed:   0.20,
			memReqGB:  1.0,
			memLimGB:  2.0,
			memUsedGB: 0.30,
		},
		{
			namespace: "batch-jobs",
			name:      "report-worker",
			kind:      "Deployment",
			replicas:  1,
			ready:     1,
			cName:     "worker",
			cpuReq:    4.0,
			cpuLimit:  6.0,
			cpuUsed:   0.30,
			memReqGB:  16.0,
			memLimGB:  24.0,
			memUsedGB: 2.10,
		},
	}

	var workloads []WorkloadCostData
	nsMap := make(map[string]*NamespaceCostData)

	for _, raw := range rawWorkloads {
		if namespaceFilter != "" && !strings.EqualFold(raw.namespace, namespaceFilter) {
			continue
		}

		reqBytes := int64(raw.memReqGB * finops.BytesPerGigabyte)
		limBytes := int64(raw.memLimGB * finops.BytesPerGigabyte)
		usedBytes := int64(raw.memUsedGB * finops.BytesPerGigabyte)

		cWaste := calc.CalculateWaste(raw.cpuReq, raw.cpuUsed, reqBytes, usedBytes)

		container := ContainerCostData{
			Name:            raw.cName,
			CPURequest:      raw.cpuReq,
			CPULimit:        raw.cpuLimit,
			CPUUsage:        raw.cpuUsed,
			MemoryRequest:   reqBytes,
			MemoryLimit:     limBytes,
			MemoryUsage:     usedBytes,
			MonthlyCost:     cWaste.TotalCost,
			MonthlyWaste:    cWaste.TotalWaste,
			WastePercentage: cWaste.WastePercentage,
		}

		repl := float64(raw.replicas)
		totReqBytes := reqBytes * int64(raw.replicas)
		totUsedBytes := usedBytes * int64(raw.replicas)
		totCPUReq := raw.cpuReq * repl
		totCPUUsed := raw.cpuUsed * repl

		wWaste := calc.CalculateWaste(totCPUReq, totCPUUsed, totReqBytes, totUsedBytes)

		w := WorkloadCostData{
			Namespace:          raw.namespace,
			Name:               raw.name,
			Kind:               raw.kind,
			Replicas:           raw.replicas,
			ReadyReplicas:      raw.ready,
			Containers:         []ContainerCostData{container},
			TotalCPURequest:    totCPUReq,
			TotalCPUUsage:      totCPUUsed,
			TotalMemoryRequest: totReqBytes,
			TotalMemoryUsage:   totUsedBytes,
			MonthlyCost:        wWaste.TotalCost,
			MonthlyWaste:       wWaste.TotalWaste,
			WastePercentage:    wWaste.WastePercentage,
			EfficiencyScore:    wWaste.EfficiencyScore,
		}
		workloads = append(workloads, w)

		// Aggregate namespace stats
		ns, exists := nsMap[raw.namespace]
		if !exists {
			ns = &NamespaceCostData{
				Namespace: raw.namespace,
			}
			nsMap[raw.namespace] = ns
		}
		ns.WorkloadCount++
		ns.PodCount += int(raw.replicas)
		ns.TotalCPURequest += totCPUReq
		ns.TotalCPUUsage += totCPUUsed
		ns.TotalMemoryRequest += totReqBytes
		ns.TotalMemoryUsage += totUsedBytes
		ns.MonthlyCost += wWaste.TotalCost
		ns.MonthlyWaste += wWaste.TotalWaste
	}

	var namespaces []NamespaceCostData
	for _, ns := range nsMap {
		if ns.MonthlyCost > 0 {
			ns.WastePercentage = math.Round((ns.MonthlyWaste/ns.MonthlyCost)*1000) / 10
		}
		ns.MonthlyCost = math.Round(ns.MonthlyCost*100) / 100
		ns.MonthlyWaste = math.Round(ns.MonthlyWaste*100) / 100
		namespaces = append(namespaces, *ns)
	}

	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].MonthlyWaste > namespaces[j].MonthlyWaste
	})

	nodes := []NodeCostData{
		{
			Name:                    "node-pool-1-a1",
			Status:                  "Ready",
			Roles:                   "worker",
			CPUCapacity:             8.0,
			CPUAllocatable:          7.8,
			CPURequests:             6.5,
			CPUUsage:                2.4,
			MemoryCapacity:          32 * int64(finops.BytesPerGigabyte),
			MemoryAllocatable:       31 * int64(finops.BytesPerGigabyte),
			MemoryRequests:          26 * int64(finops.BytesPerGigabyte),
			MemoryUsage:             12 * int64(finops.BytesPerGigabyte),
			MonthlyCost:             calc.MonthlyCost(7.8, 31*int64(finops.BytesPerGigabyte)),
			CPUAllocationPercent:    math.Round((6.5/7.8)*1000) / 10,
			MemoryAllocationPercent: math.Round((26.0/31.0)*1000) / 10,
		},
		{
			Name:                    "node-pool-1-b2",
			Status:                  "Ready",
			Roles:                   "worker",
			CPUCapacity:             8.0,
			CPUAllocatable:          7.8,
			CPURequests:             5.8,
			CPUUsage:                3.1,
			MemoryCapacity:          32 * int64(finops.BytesPerGigabyte),
			MemoryAllocatable:       31 * int64(finops.BytesPerGigabyte),
			MemoryRequests:          24 * int64(finops.BytesPerGigabyte),
			MemoryUsage:             15 * int64(finops.BytesPerGigabyte),
			MonthlyCost:             calc.MonthlyCost(7.8, 31*int64(finops.BytesPerGigabyte)),
			CPUAllocationPercent:    math.Round((5.8/7.8)*1000) / 10,
			MemoryAllocationPercent: math.Round((24.0/31.0)*1000) / 10,
		},
		{
			Name:                    "node-pool-spot-c3",
			Status:                  "Ready",
			Roles:                   "worker,spot",
			CPUCapacity:             16.0,
			CPUAllocatable:          15.6,
			CPURequests:             14.2,
			CPUUsage:                5.2,
			MemoryCapacity:          64 * int64(finops.BytesPerGigabyte),
			MemoryAllocatable:       62 * int64(finops.BytesPerGigabyte),
			MemoryRequests:          52 * int64(finops.BytesPerGigabyte),
			MemoryUsage:             22 * int64(finops.BytesPerGigabyte),
			MonthlyCost:             calc.MonthlyCost(15.6, 62*int64(finops.BytesPerGigabyte)),
			CPUAllocationPercent:    math.Round((14.2/15.6)*1000) / 10,
			MemoryAllocationPercent: math.Round((52.0/62.0)*1000) / 10,
		},
	}

	return workloads, namespaces, nodes, nil
}
