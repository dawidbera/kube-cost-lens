package k8s

import (
	"testing"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestBuildWorkloadData(t *testing.T) {
	calc := finops.NewCostCalculator(finops.DefaultPricingCard())

	containers := []corev1.Container{
		{
			Name: "app",
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("1Gi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				},
			},
		},
	}

	w := buildWorkloadData("prod", "order-svc", "Deployment", 3, 3, containers, calc)

	if w.Namespace != "prod" || w.Name != "order-svc" || w.Kind != "Deployment" {
		t.Errorf("unexpected workload metadata: %+v", w)
	}
	if w.Replicas != 3 || w.ReadyReplicas != 3 {
		t.Errorf("expected 3 replicas, got %d", w.Replicas)
	}
	if w.TotalCPURequest != 1.5 { // 0.5 * 3
		t.Errorf("expected total CPU request 1.5, got %f", w.TotalCPURequest)
	}
	if len(w.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(w.Containers))
	}
	if w.MonthlyCost <= 0 || w.MonthlyWaste <= 0 {
		t.Errorf("expected positive monthly cost and waste: %+v", w)
	}
}
