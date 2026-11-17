// Package k8s defines cluster models, interfaces, and Kubernetes API integrations.
package k8s

import (
	"context"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
)

// ContainerCostData holds resource allocation and cost metrics for a single container.
type ContainerCostData struct {
	Name string

	CPURequest float64 // in cores (e.g., 0.5 = 500m)
	CPULimit   float64 // in cores
	CPUUsage   float64 // in cores

	MemoryRequest int64 // in bytes
	MemoryLimit   int64 // in bytes
	MemoryUsage   int64 // in bytes

	MonthlyCost     float64
	MonthlyWaste    float64
	WastePercentage float64
}

// WorkloadCostData represents a high-level Kubernetes controller (Deployment, StatefulSet, DaemonSet).
type WorkloadCostData struct {
	Namespace     string
	Name          string
	Kind          string
	Replicas      int32
	ReadyReplicas int32

	Containers []ContainerCostData

	TotalCPURequest float64
	TotalCPUUsage   float64

	TotalMemoryRequest int64
	TotalMemoryUsage   int64

	MonthlyCost     float64
	MonthlyWaste    float64
	WastePercentage float64
	EfficiencyScore float64
}

// NamespaceCostData represents aggregate resource and financial totals for a namespace.
type NamespaceCostData struct {
	Namespace     string
	WorkloadCount int
	PodCount      int

	TotalCPURequest float64
	TotalCPUUsage   float64

	TotalMemoryRequest int64
	TotalMemoryUsage   int64

	MonthlyCost     float64
	MonthlyWaste    float64
	WastePercentage float64
}

// NodeCostData represents node compute capacity, allocation, and overcommit stats.
type NodeCostData struct {
	Name   string
	Status string
	Roles  string

	CPUCapacity    float64
	CPUAllocatable float64
	CPURequests    float64
	CPUUsage       float64

	MemoryCapacity    int64
	MemoryAllocatable int64
	MemoryRequests    int64
	MemoryUsage       int64

	MonthlyCost            float64
	CPUAllocationPercent   float64
	MemoryAllocationPercent float64
}

// ClusterProvider defines the interface for acquiring cluster inventory and usage metrics.
type ClusterProvider interface {
	// FetchClusterData queries workloads, namespaces, and nodes, applying cost calculations.
	FetchClusterData(ctx context.Context, namespaceFilter string, calc *finops.CostCalculator) ([]WorkloadCostData, []NamespaceCostData, []NodeCostData, error)
	// GetContextName returns the current active cluster/kubeconfig context name.
	GetContextName() string
	// Close terminates any active informers or connections.
	Close() error
}
