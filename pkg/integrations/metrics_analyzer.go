package integrations

import (
	"fmt"
	"time"

	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"k8s.io/client-go/kubernetes"
)

// MetricsAnalyzer combines Kubernetes and Prometheus data for enhanced analysis
type MetricsAnalyzer struct {
	k8sClient  kubernetes.Interface
	promClient *PrometheusClient
}

// NewMetricsAnalyzer creates a new metrics analyzer
func NewMetricsAnalyzer(k8sClient kubernetes.Interface, prometheusURL string) *MetricsAnalyzer {
	promClient := NewPrometheusClient(prometheusURL)
	return &MetricsAnalyzer{
		k8sClient:  k8sClient,
		promClient: promClient,
	}
}

// EnhancedPodReport combines diagnostic and metrics data
type EnhancedPodReport struct {
	PodReport       *diagnostics.PodReport
	PodMetrics      *PodMetrics
	Recommendations []string
	HealthScore     int
}

// AnalyzePodWithMetrics enhances pod analysis with metrics
func (m *MetricsAnalyzer) AnalyzePodWithMetrics(podName, namespace string) (*EnhancedPodReport, error) {
	// Get standard pod analysis
	podAnalyzer := diagnostics.NewPodAnalyzer(m.k8sClient, namespace)
	podReport, err := podAnalyzer.Analyze(podName)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze pod: %v", err)
	}

	// Get metrics
	metrics, err := m.promClient.GetPodMetrics(podName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %v", err)
	}

	report := &EnhancedPodReport{
		PodReport:  podReport,
		PodMetrics: metrics,
	}

	m.generateRecommendations(report)
	report.HealthScore = m.calculateHealthScore(report)
	return report, nil
}

// AnalyzeNodeWithMetrics enhances node analysis with metrics
func (m *MetricsAnalyzer) AnalyzeNodeWithMetrics(nodeName string) (*NodeMetrics, error) {
	metrics, err := m.promClient.GetNodeMetrics(nodeName)
	if err != nil {
		return &NodeMetrics{
			NodeName:  nodeName,
			Timestamp: time.Now(),
			Error:     fmt.Sprintf("Failed to get node metrics: %v", err),
		}, fmt.Errorf("failed to get node metrics: %v", err)
	}
	return metrics, nil
}

// AnalyzeClusterWithMetrics provides cluster-level metrics analysis
func (m *MetricsAnalyzer) AnalyzeClusterWithMetrics() (*ClusterMetrics, error) {
	metrics, err := m.promClient.GetClusterMetrics()
	if err != nil {
		return &ClusterMetrics{
			Timestamp: time.Now(),
			Error:     fmt.Sprintf("Failed to get cluster metrics: %v", err),
		}, fmt.Errorf("failed to get cluster metrics: %v", err)
	}
	return metrics, nil
}

func (m *MetricsAnalyzer) generateRecommendations(report *EnhancedPodReport) {
	var recommendations []string

	// Only generate metric-based recommendations if we have metrics
	if report.PodMetrics.Error == "" {
		// Check CPU usage
		if report.PodMetrics.CPUUsage > 0.8 {
			recommendations = append(recommendations,
				"High CPU usage detected - consider increasing CPU limits or optimizing application")
		} else if report.PodMetrics.CPUUsage < 0.1 {
			recommendations = append(recommendations,
				"Low CPU usage - consider reducing CPU requests to improve node utilization")
		}

		// Check memory usage
		if report.PodMetrics.MemoryUsage > 1024*1024*1024 {
			recommendations = append(recommendations,
				"High memory usage detected - monitor for memory leaks and consider increasing memory limits")
		}

		// Check network usage
		if report.PodMetrics.NetworkRx > 1000000 {
			recommendations = append(recommendations,
				"High network receive traffic - ensure network policies and limits are appropriate")
		}
		if report.PodMetrics.NetworkTx > 1000000 {
			recommendations = append(recommendations,
				"High network transmit traffic - ensure network policies and limits are appropriate")
		}
	}

	// Add Prometheus setup recommendation if metrics are unavailable
	if report.PodMetrics.Error != "" {
		recommendations = append(recommendations,
			"Prometheus metrics unavailable - set up Prometheus for enhanced monitoring")
	}

	// Combine with existing pod report recommendations
	recommendations = append(recommendations, report.PodReport.Recommendations...)
	report.Recommendations = recommendations
}

func (m *MetricsAnalyzer) calculateHealthScore(report *EnhancedPodReport) int {
	score := 100

	// Deduct points for pod issues
	if len(report.PodReport.Issues) > 0 {
		score -= len(report.PodReport.Issues) * 10
	}

	// Only deduct for metric issues if we have metrics
	if report.PodMetrics.Error == "" {
		if report.PodMetrics.CPUUsage > 0.9 {
			score -= 20
		} else if report.PodMetrics.CPUUsage > 0.8 {
			score -= 10
		}

		if report.PodMetrics.MemoryUsage > 2*1024*1024*1024 {
			score -= 15
		}
	} else {
		// Deduct points for missing metrics
		score -= 10
	}

	if score < 0 {
		score = 0
	}
	return score
}
