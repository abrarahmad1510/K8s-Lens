package machinelearning

import (
	"context"
	"fmt"
	"math"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// AnomalyDetector identifies unusual patterns in cluster behavior
type AnomalyDetector struct {
	client kubernetes.Interface
}

// NewAnomalyDetector creates a new anomaly detector
func NewAnomalyDetector(client kubernetes.Interface) *AnomalyDetector {
	return &AnomalyDetector{
		client: client,
	}
}

// AnomalyReport contains detected anomalies
type AnomalyReport struct {
	Namespace       string
	TotalPods       int
	Anomalies       []Anomaly
	Score           int // 0-100, higher is worse
	Recommendations []string
	Timestamp       time.Time
}

// Anomaly represents a detected anomaly
type Anomaly struct {
	Type       string
	Severity   string // Low, Medium, High, Critical
	Resource   string
	Message    string
	Confidence float64
	Timestamp  time.Time
}

// DetectNamespaceAnomalies analyzes a namespace for unusual patterns
func (a *AnomalyDetector) DetectNamespaceAnomalies(namespace string) (*AnomalyReport, error) {
	report := &AnomalyReport{
		Namespace: namespace,
		Timestamp: time.Now(),
	}

	// Get all pods in the namespace
	pods, err := a.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	report.TotalPods = len(pods.Items)

	// Analyze each pod for anomalies
	for _, pod := range pods.Items {
		podAnomalies := a.analyzePodAnomalies(&pod)
		report.Anomalies = append(report.Anomalies, podAnomalies...)
	}

	// Analyze namespace-level anomalies
	nsAnomalies := a.analyzeNamespaceLevelAnomalies(pods.Items)
	report.Anomalies = append(report.Anomalies, nsAnomalies...)

	// Calculate overall anomaly score
	report.Score = a.calculateAnomalyScore(report.Anomalies)
	report.Recommendations = a.generateRecommendations(report.Anomalies)

	return report, nil
}

func (a *AnomalyDetector) analyzePodAnomalies(pod *corev1.Pod) []Anomaly {
	var anomalies []Anomaly

	// Check for restart anomalies
	if a.detectRestartAnomaly(pod) {
		anomalies = append(anomalies, Anomaly{
			Type:       "RestartPattern",
			Severity:   "High",
			Resource:   pod.Name,
			Message:    "Unusual pod restart pattern detected",
			Confidence: 0.85,
			Timestamp:  time.Now(),
		})
	}

	// Check for resource anomalies
	resourceAnomalies := a.detectResourceAnomalies(pod)
	anomalies = append(anomalies, resourceAnomalies...)

	// Check for status anomalies
	statusAnomalies := a.detectStatusAnomalies(pod)
	anomalies = append(anomalies, statusAnomalies...)

	return anomalies
}

func (a *AnomalyDetector) detectRestartAnomaly(pod *corev1.Pod) bool {
	totalRestarts := 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		totalRestarts += int(containerStatus.RestartCount)
	}

	// If a pod has restarted more than 10 times, it's anomalous
	return totalRestarts > 10
}

func (a *AnomalyDetector) detectResourceAnomalies(pod *corev1.Pod) []Anomaly {
	var anomalies []Anomaly

	for _, container := range pod.Spec.Containers {
		// Check for missing resource requests
		if container.Resources.Requests == nil {
			anomalies = append(anomalies, Anomaly{
				Type:       "MissingResourceRequests",
				Severity:   "Medium",
				Resource:   fmt.Sprintf("%s/%s", pod.Name, container.Name),
				Message:    "Container missing resource requests",
				Confidence: 1.0,
				Timestamp:  time.Now(),
			})
		}

		// Check for unbalanced CPU/Memory ratios
		if container.Resources.Requests != nil {
			cpu := container.Resources.Requests[corev1.ResourceCPU]
			memory := container.Resources.Requests[corev1.ResourceMemory]

			if !cpu.IsZero() && !memory.IsZero() {
				cpuMilli := cpu.MilliValue()
				memoryMB := memory.Value() / (1024 * 1024)

				// Typical ratio: 1 CPU core per 4GB RAM
				if cpuMilli > 0 && memoryMB > 0 {
					ratio := float64(memoryMB) / float64(cpuMilli)
					if ratio < 500 || ratio > 8000 { // Outside typical range
						anomalies = append(anomalies, Anomaly{
							Type:       "UnbalancedResources",
							Severity:   "Low",
							Resource:   fmt.Sprintf("%s/%s", pod.Name, container.Name),
							Message:    fmt.Sprintf("Unusual CPU/Memory ratio: %.2f MB per CPU core", ratio),
							Confidence: 0.75,
							Timestamp:  time.Now(),
						})
					}
				}
			}
		}
	}

	return anomalies
}

func (a *AnomalyDetector) detectStatusAnomalies(pod *corev1.Pod) []Anomaly {
	var anomalies []Anomaly

	// Check for long-running pending pods
	if pod.Status.Phase == corev1.PodPending {
		duration := time.Since(pod.CreationTimestamp.Time)
		if duration > 10*time.Minute {
			anomalies = append(anomalies, Anomaly{
				Type:       "LongPending",
				Severity:   "High",
				Resource:   pod.Name,
				Message:    fmt.Sprintf("Pod has been pending for %v", duration),
				Confidence: 0.9,
				Timestamp:  time.Now(),
			})
		}
	}

	return anomalies
}

func (a *AnomalyDetector) analyzeNamespaceLevelAnomalies(pods []corev1.Pod) []Anomaly {
	var anomalies []Anomaly

	// Check namespace resource distribution
	totalCPU := int64(0)
	totalMemory := int64(0)

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				cpu := container.Resources.Requests[corev1.ResourceCPU]
				memory := container.Resources.Requests[corev1.ResourceMemory]

				totalCPU += cpu.MilliValue()
				totalMemory += memory.Value() / (1024 * 1024) // Convert to MB
			}
		}
	}

	// Check for resource concentration anomalies
	if len(pods) > 0 {
		avgCPUPerPod := float64(totalCPU) / float64(len(pods))
		if avgCPUPerPod > 4000 { // More than 4 CPUs per pod on average
			anomalies = append(anomalies, Anomaly{
				Type:       "HighResourceConcentration",
				Severity:   "Medium",
				Resource:   "Namespace",
				Message:    fmt.Sprintf("High CPU concentration: %.2f millicores per pod average", avgCPUPerPod),
				Confidence: 0.8,
				Timestamp:  time.Now(),
			})
		}
	}

	return anomalies
}

func (a *AnomalyDetector) calculateAnomalyScore(anomalies []Anomaly) int {
	score := 0
	for _, anomaly := range anomalies {
		switch anomaly.Severity {
		case "Critical":
			score += 10
		case "High":
			score += 7
		case "Medium":
			score += 4
		case "Low":
			score += 1
		}
	}

	return int(math.Min(100, float64(score)))
}

func (a *AnomalyDetector) generateRecommendations(anomalies []Anomaly) []string {
	recommendations := []string{}

	hasRestartAnomalies := false
	hasResourceAnomalies := false

	for _, anomaly := range anomalies {
		switch anomaly.Type {
		case "RestartPattern":
			hasRestartAnomalies = true
		case "MissingResourceRequests", "UnbalancedResources":
			hasResourceAnomalies = true
		}
	}

	if hasRestartAnomalies {
		recommendations = append(recommendations,
			"Investigate pod restart patterns - check application logs and resource limits")
	}

	if hasResourceAnomalies {
		recommendations = append(recommendations,
			"Review and optimize resource requests and limits for better scheduling")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "No critical issues detected - maintain current monitoring")
	}

	return recommendations
}
