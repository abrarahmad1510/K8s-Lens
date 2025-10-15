package optimization

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ResourceOptimizer provides resource optimization recommendations
type ResourceOptimizer struct {
	client kubernetes.Interface
}

// NewResourceOptimizer creates a new ResourceOptimizer
func NewResourceOptimizer(client kubernetes.Interface) *ResourceOptimizer {
	return &ResourceOptimizer{
		client: client,
	}
}

// OptimizationReport contains resource optimization recommendations
type OptimizationReport struct {
	Namespace     string
	TotalPods     int
	AnalyzedPods  int
	Optimizations []Optimization
	CostSavings   CostSavings
	Summary       OptimizationSummary
}

// Optimization represents a single optimization recommendation
type Optimization struct {
	PodName       string
	ContainerName string
	Type          string
	Current       ResourceValues
	Recommended   ResourceValues
	Savings       CostSavings
	Confidence    int
	Description   string
}

// ResourceValues represents CPU and Memory values
type ResourceValues struct {
	CPU    string
	Memory string
}

// CostSavings represents estimated cost savings
type CostSavings struct {
	MonthlySavings float64
	PercentSavings float64
	Reason         string
}

// OptimizationSummary provides an overall summary
type OptimizationSummary struct {
	TotalMonthlySavings float64
	TotalOptimizations  int
	OverallConfidence   int
	RiskLevel           string
}

// AnalyzeNamespace analyzes resource usage in a namespace
func (r *ResourceOptimizer) AnalyzeNamespace(namespace string) (*OptimizationReport, error) {
	// Get all pods in the namespace
	pods, err := r.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods in namespace %s: %v", namespace, err)
	}

	report := &OptimizationReport{
		Namespace:    namespace,
		TotalPods:    len(pods.Items),
		AnalyzedPods: 0,
	}

	totalMonthlySavings := 0.0
	optimizationCount := 0
	totalConfidence := 0

	for _, pod := range pods.Items {
		podOptimizations := r.analyzePodResources(&pod)
		report.Optimizations = append(report.Optimizations, podOptimizations...)

		for _, opt := range podOptimizations {
			totalMonthlySavings += opt.Savings.MonthlySavings
			totalConfidence += opt.Confidence
			optimizationCount++
		}
		report.AnalyzedPods++
	}

	// Calculate summary
	if optimizationCount > 0 {
		report.Summary.TotalMonthlySavings = totalMonthlySavings
		report.Summary.TotalOptimizations = optimizationCount
		report.Summary.OverallConfidence = totalConfidence / optimizationCount

		// Determine risk level
		if report.Summary.OverallConfidence >= 80 {
			report.Summary.RiskLevel = "Low"
		} else if report.Summary.OverallConfidence >= 60 {
			report.Summary.RiskLevel = "Medium"
		} else {
			report.Summary.RiskLevel = "High"
		}
	}

	return report, nil
}

func (r *ResourceOptimizer) analyzePodResources(pod *corev1.Pod) []Optimization {
	var optimizations []Optimization

	for _, container := range pod.Spec.Containers {
		// Analyze requests vs potential optimizations
		if container.Resources.Requests != nil {
			cpuRequest := container.Resources.Requests[corev1.ResourceCPU]
			memoryRequest := container.Resources.Requests[corev1.ResourceMemory]

			// Check for over-provisioned CPU
			if !cpuRequest.IsZero() {
				currentCPU := cpuRequest.String()
				recommendedCPU := r.calculateRecommendedCPU(cpuRequest)

				if recommendedCPU != currentCPU {
					optimizations = append(optimizations, Optimization{
						PodName:       pod.Name,
						ContainerName: container.Name,
						Type:          "CPU Right-Sizing",
						Current:       ResourceValues{CPU: currentCPU},
						Recommended:   ResourceValues{CPU: recommendedCPU},
						Savings: CostSavings{
							MonthlySavings: r.calculateCPUSavings(cpuRequest, recommendedCPU),
							PercentSavings: 25.0,
							Reason:         "CPU is over-provisioned based on usage patterns",
						},
						Confidence:  75,
						Description: "Reduce CPU requests to match actual usage patterns",
					})
				}
			}

			// Check for over-provisioned Memory
			if !memoryRequest.IsZero() {
				currentMemory := memoryRequest.String()
				recommendedMemory := r.calculateRecommendedMemory(memoryRequest)

				if recommendedMemory != currentMemory {
					optimizations = append(optimizations, Optimization{
						PodName:       pod.Name,
						ContainerName: container.Name,
						Type:          "Memory Right-Sizing",
						Current:       ResourceValues{Memory: currentMemory},
						Recommended:   ResourceValues{Memory: recommendedMemory},
						Savings: CostSavings{
							MonthlySavings: r.calculateMemorySavings(memoryRequest, recommendedMemory),
							PercentSavings: 30.0,
							Reason:         "Memory is over-provisioned based on usage patterns",
						},
						Confidence:  80,
						Description: "Reduce memory requests to match actual usage patterns",
					})
				}
			}
		}

		// Check for missing limits
		if container.Resources.Limits == nil || len(container.Resources.Limits) == 0 {
			optimizations = append(optimizations, Optimization{
				PodName:       pod.Name,
				ContainerName: container.Name,
				Type:          "Missing Resource Limits",
				Current:       ResourceValues{CPU: "Not set", Memory: "Not set"},
				Recommended:   ResourceValues{CPU: "500m", Memory: "512Mi"},
				Savings: CostSavings{
					MonthlySavings: 0, // No direct savings, but prevents cost spikes
					PercentSavings: 0,
					Reason:         "Prevents cost spikes from resource exhaustion",
				},
				Confidence:  95,
				Description: "Add resource limits to prevent runaway resource consumption",
			})
		}
	}

	return optimizations
}

func (r *ResourceOptimizer) calculateRecommendedCPU(currentCPU resource.Quantity) string {
	// Simplified calculation - in real implementation, this would use metrics
	// For demonstration, we're recommending a fixed value
	return "250m"
}

func (r *ResourceOptimizer) calculateRecommendedMemory(currentMemory resource.Quantity) string {
	// Simplified calculation - in real implementation, this would use metrics
	// For demonstration, we're recommending a fixed value
	return "256Mi"
}

func (r *ResourceOptimizer) calculateCPUSavings(current resource.Quantity, recommended string) float64 {
	// Simplified calculation - real implementation would use cloud pricing
	// Convert both to milliCPU for comparison
	currentMilli := current.MilliValue()

	// Parse recommended (assuming format like "250m")
	var recommendedMilli int64
	if recommended == "250m" {
		recommendedMilli = 250
	} else {
		recommendedMilli = 500 // default assumption
	}

	// Calculate savings based on difference
	savings := float64(currentMilli-recommendedMilli) * 0.01 // $0.01 per milliCPU per month
	if savings < 0 {
		return 0
	}
	return savings
}

func (r *ResourceOptimizer) calculateMemorySavings(current resource.Quantity, recommended string) float64 {
	// Simplified calculation - real implementation would use cloud pricing
	// Convert both to bytes for comparison
	currentBytes := current.Value()

	// Parse recommended (assuming format like "256Mi")
	var recommendedBytes int64
	if recommended == "256Mi" {
		recommendedBytes = 256 * 1024 * 1024 // 256 MiB in bytes
	} else {
		recommendedBytes = 512 * 1024 * 1024 // default assumption
	}

	// Calculate savings based on difference
	savings := float64(currentBytes-recommendedBytes) * 0.000000001 // $0.001 per MB per month
	if savings < 0 {
		return 0
	}
	return savings
}
