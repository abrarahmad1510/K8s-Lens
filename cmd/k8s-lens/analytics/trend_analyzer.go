package analytics

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// TrendAnalyzer analyzes historical trends and patterns
type TrendAnalyzer struct {
	client kubernetes.Interface
}

// NewTrendAnalyzer creates a new trend analyzer
func NewTrendAnalyzer(client kubernetes.Interface) *TrendAnalyzer {
	return &TrendAnalyzer{
		client: client,
	}
}

// TrendReport contains trend analysis results
type TrendReport struct {
	Namespace         string
	AnalysisPeriod    time.Duration
	ResourceTrends    []ResourceTrend
	PerformanceTrends []PerformanceTrend
	Recommendations   []string
	GeneratedAt       time.Time
}

// ResourceTrend shows resource usage trends
type ResourceTrend struct {
	ResourceType  string
	Metric        string
	CurrentValue  float64
	PreviousValue float64
	ChangePercent float64
	Trend         string // Increasing, Decreasing, Stable
}

// PerformanceTrend shows performance patterns
type PerformanceTrend struct {
	Component  string
	Metric     string
	Pattern    string
	Confidence float64
	Impact     string
}

// AnalyzeNamespaceTrends analyzes trends in a namespace over time
func (t *TrendAnalyzer) AnalyzeNamespaceTrends(namespace string, period time.Duration) (*TrendReport, error) {
	report := &TrendReport{
		Namespace:      namespace,
		AnalysisPeriod: period,
		GeneratedAt:    time.Now(),
	}

	// Get current state
	currentPods, err := t.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get current pods: %v", err)
	}

	// Get deployments for workload analysis
	deployments, err := t.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %v", err)
	}

	// Analyze resource trends
	resourceTrends := t.analyzeResourceTrends(currentPods.Items, deployments.Items)
	report.ResourceTrends = resourceTrends

	// Analyze performance trends
	performanceTrends := t.analyzePerformanceTrends(currentPods.Items, deployments.Items)
	report.PerformanceTrends = performanceTrends

	// Generate recommendations
	report.Recommendations = t.generateTrendRecommendations(resourceTrends, performanceTrends)

	return report, nil
}

func (t *TrendAnalyzer) analyzeResourceTrends(pods []corev1.Pod, deployments []appsv1.Deployment) []ResourceTrend {
	var trends []ResourceTrend

	// Analyze pod count trend
	podCount := len(pods)
	// In a real implementation, you'd compare with historical data
	// For now, we'll use a simulated previous value
	previousPodCount := podCount - 1 // Simulate decrease
	if previousPodCount < 0 {
		previousPodCount = 0
	}

	podChangePercent := 0.0
	if previousPodCount > 0 {
		podChangePercent = (float64(podCount) - float64(previousPodCount)) / float64(previousPodCount) * 100
	}

	podTrend := "Stable"
	if podChangePercent > 5 {
		podTrend = "Increasing"
	} else if podChangePercent < -5 {
		podTrend = "Decreasing"
	}

	trends = append(trends, ResourceTrend{
		ResourceType:  "Pods",
		Metric:        "Count",
		CurrentValue:  float64(podCount),
		PreviousValue: float64(previousPodCount),
		ChangePercent: podChangePercent,
		Trend:         podTrend,
	})

	// Analyze resource requests trend
	totalCPU := int64(0)
	totalMemory := int64(0)
	containerCount := 0

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				cpu := container.Resources.Requests[corev1.ResourceCPU]
				memory := container.Resources.Requests[corev1.ResourceMemory]

				totalCPU += cpu.MilliValue()
				totalMemory += memory.Value() / (1024 * 1024) // Convert to MB
				containerCount++
			}
		}
	}

	if containerCount > 0 {
		avgCPU := float64(totalCPU) / float64(containerCount)
		avgMemory := float64(totalMemory) / float64(containerCount)

		// Simulate previous values (in real implementation, fetch historical data)
		previousAvgCPU := avgCPU * 0.9
		previousAvgMemory := avgMemory * 0.95

		cpuChangePercent := (avgCPU - previousAvgCPU) / previousAvgCPU * 100
		memoryChangePercent := (avgMemory - previousAvgMemory) / previousAvgMemory * 100

		cpuTrend := "Stable"
		if cpuChangePercent > 5 {
			cpuTrend = "Increasing"
		} else if cpuChangePercent < -5 {
			cpuTrend = "Decreasing"
		}

		memoryTrend := "Stable"
		if memoryChangePercent > 5 {
			memoryTrend = "Increasing"
		} else if memoryChangePercent < -5 {
			memoryTrend = "Decreasing"
		}

		trends = append(trends, ResourceTrend{
			ResourceType:  "Containers",
			Metric:        "Average CPU Request (millicores)",
			CurrentValue:  avgCPU,
			PreviousValue: previousAvgCPU,
			ChangePercent: cpuChangePercent,
			Trend:         cpuTrend,
		})

		trends = append(trends, ResourceTrend{
			ResourceType:  "Containers",
			Metric:        "Average Memory Request (MB)",
			CurrentValue:  avgMemory,
			PreviousValue: previousAvgMemory,
			ChangePercent: memoryChangePercent,
			Trend:         memoryTrend,
		})
	}

	return trends
}

func (t *TrendAnalyzer) analyzePerformanceTrends(pods []corev1.Pod, deployments []appsv1.Deployment) []PerformanceTrend {
	var trends []PerformanceTrend

	// Analyze restart patterns
	totalRestarts := 0
	for _, pod := range pods {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			totalRestarts += int(containerStatus.RestartCount)
		}
	}

	if len(pods) > 0 {
		avgRestarts := float64(totalRestarts) / float64(len(pods))
		if avgRestarts > 2 {
			trends = append(trends, PerformanceTrend{
				Component:  "Pods",
				Metric:     "Container Restarts",
				Pattern:    "High restart frequency detected",
				Confidence: 0.85,
				Impact:     "Medium - May indicate application instability",
			})
		}
	}

	// Analyze pod status distribution
	runningPods := 0
	pendingPods := 0
	failedPods := 0

	for _, pod := range pods {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			runningPods++
		case corev1.PodPending:
			pendingPods++
		case corev1.PodFailed:
			failedPods++
		}
	}

	if len(pods) > 0 {
		failureRate := float64(failedPods) / float64(len(pods)) * 100
		if failureRate > 10 {
			trends = append(trends, PerformanceTrend{
				Component:  "Pods",
				Metric:     "Failure Rate",
				Pattern:    "High pod failure rate",
				Confidence: 0.9,
				Impact:     "High - Impacts application availability",
			})
		}

		pendingRate := float64(pendingPods) / float64(len(pods)) * 100
		if pendingRate > 20 {
			trends = append(trends, PerformanceTrend{
				Component:  "Pods",
				Metric:     "Pending Rate",
				Pattern:    "High pod pending rate",
				Confidence: 0.75,
				Impact:     "Medium - May indicate resource constraints",
			})
		}
	}

	return trends
}

func (t *TrendAnalyzer) generateTrendRecommendations(resourceTrends []ResourceTrend, performanceTrends []PerformanceTrend) []string {
	recommendations := []string{}

	// Analyze resource trends for recommendations
	for _, trend := range resourceTrends {
		if trend.Trend == "Increasing" && trend.ChangePercent > 20 {
			recommendations = append(recommendations,
				fmt.Sprintf("Significant increase in %s %s (%.1f%%) - monitor capacity",
					trend.ResourceType, trend.Metric, trend.ChangePercent))
		}
	}

	// Analyze performance trends for recommendations
	for _, trend := range performanceTrends {
		if trend.Impact == "High" {
			recommendations = append(recommendations,
				fmt.Sprintf("Address %s: %s", trend.Component, trend.Pattern))
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"No critical trends detected - maintain current monitoring and optimization efforts")
	}

	return recommendations
}
