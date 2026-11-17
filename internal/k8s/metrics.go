// Package k8s provides helper functions for parsing Kubernetes resource metrics.
package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
)

// ParseCPUQuantity converts Kubernetes CPU string quantity (e.g. "500m", "2") into float64 cores.
func ParseCPUQuantity(val string) (float64, error) {
	q, err := resource.ParseQuantity(val)
	if err != nil {
		return 0, fmt.Errorf("failed parsing CPU quantity %q: %w", val, err)
	}
	return float64(q.MilliValue()) / 1000.0, nil
}

// ParseMemoryQuantity converts Kubernetes memory string quantity (e.g. "512Mi", "4Gi") into bytes.
func ParseMemoryQuantity(val string) (int64, error) {
	q, err := resource.ParseQuantity(val)
	if err != nil {
		return 0, fmt.Errorf("failed parsing Memory quantity %q: %w", val, err)
	}
	return q.Value(), nil
}

// FormatCPU cores as human readable millicores or full cores string.
func FormatCPU(cores float64) string {
	if cores < 1.0 {
		return fmt.Sprintf("%dm", int64(cores*1000))
	}
	return fmt.Sprintf("%.2f", cores)
}

// FormatMemoryBytes formats bytes into human-readable Mi or Gi for Kubernetes manifests and UI.
func FormatMemoryBytes(bytes int64) string {
	const (
		MiB = 1024 * 1024
		GiB = 1024 * 1024 * 1024
	)
	if bytes >= GiB && bytes%GiB == 0 {
		return fmt.Sprintf("%dGi", bytes/GiB)
	}
	return fmt.Sprintf("%dMi", bytes/MiB)
}
