package ai

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PredictiveAnalyzer provides predictive failure analysis
type PredictiveAnalyzer struct {
	client kubernetes.Interface
}

// NewPredictiveAnalyzer creates a new PredictiveAnalyzer
func NewPredictiveAnalyzer(client kubernetes.Interface) *PredictiveAnalyzer {
	return &PredictiveAnalyzer{
		client: client,
	}
}

// PredictionReport contains predictive analysis results
type PredictionReport struct {
	PodName         string
	Namespace       string
	Predictions     []Prediction
	OverallRisk     string
	Confidence      int
	Recommendations []string
}

// Prediction represents a single failure prediction
type Prediction struct {
	Type        string
	Description string
	Probability int
	Timeframe   string
	Evidence    []string
}

// PredictFailures analyzes a deployment for potential failures
func (p *PredictiveAnalyzer) PredictFailures(deploymentName, namespace string) (*PredictionReport, error) {
	// Get the deployment
	deployment, err := p.client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s: %v", deploymentName, err)
	}

	// Get pods for the deployment
	pods, err := p.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods for deployment %s: %v", deploymentName, err)
	}

	// Get events for the namespace
	events, err := p.client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get events for namespace %s: %v", namespace, err)
	}

	report := &PredictionReport{
		PodName:   deploymentName,
		Namespace: namespace,
	}

	p.analyzeRestartPatterns(report, pods.Items)
	p.analyzeResourcePatterns(report, pods.Items, deployment)
	p.analyzeEventPatterns(report, events.Items)
	p.calculateOverallRisk(report)

	return report, nil
}

func (p *PredictiveAnalyzer) analyzeRestartPatterns(report *PredictionReport, pods []corev1.Pod) {
	totalRestarts := 0
	frequentRestarters := 0

	for _, pod := range pods {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			totalRestarts += int(containerStatus.RestartCount)
			if containerStatus.RestartCount > 5 {
				frequentRestarters++
			}
		}
	}

	if frequentRestarters > 0 {
		report.Predictions = append(report.Predictions, Prediction{
			Type:        "Container Crash",
			Description: "Containers are crashing frequently indicating potential stability issues",
			Probability: 70,
			Timeframe:   "Next 7 days",
			Evidence:    []string{fmt.Sprintf("%d containers have restarted more than 5 times", frequentRestarters)},
		})
	}

	if totalRestarts > 20 {
		report.Predictions = append(report.Predictions, Prediction{
			Type:        "Memory Leak",
			Description: "High number of restarts may indicate memory leaks or resource exhaustion",
			Probability: 60,
			Timeframe:   "Next 14 days",
			Evidence:    []string{fmt.Sprintf("Total of %d container restarts across all pods", totalRestarts)},
		})
	}
}

func (p *PredictiveAnalyzer) analyzeResourcePatterns(report *PredictionReport, pods []corev1.Pod, deployment interface{}) {
	// Analyze resource patterns that might lead to failures
	// This is a simplified implementation

	// Check for resource constraints
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests == nil || container.Resources.Limits == nil {
				report.Predictions = append(report.Predictions, Prediction{
					Type:        "Resource Exhaustion",
					Description: "Missing resource limits may lead to OOM kills or CPU throttling",
					Probability: 50,
					Timeframe:   "Next 30 days",
					Evidence:    []string{fmt.Sprintf("Container %s has no resource limits set", container.Name)},
				})
				break
			}
		}
	}
}

func (p *PredictiveAnalyzer) analyzeEventPatterns(report *PredictionReport, events []corev1.Event) {
	warningCount := 0
	recentWarnings := 0
	cutoffTime := time.Now().Add(-24 * time.Hour)

	for _, event := range events {
		if event.Type == "Warning" {
			warningCount++
			if event.LastTimestamp.Time.After(cutoffTime) {
				recentWarnings++
			}
		}
	}

	if recentWarnings > 10 {
		report.Predictions = append(report.Predictions, Prediction{
			Type:        "Cluster Issues",
			Description: "High number of recent warning events indicates cluster-level problems",
			Probability: 65,
			Timeframe:   "Next 3 days",
			Evidence:    []string{fmt.Sprintf("%d warning events in the last 24 hours", recentWarnings)},
		})
	}
}

func (p *PredictiveAnalyzer) calculateOverallRisk(report *PredictionReport) {
	if len(report.Predictions) == 0 {
		report.OverallRisk = "Low"
		report.Confidence = 90
		report.Recommendations = []string{"No significant risks detected. Continue monitoring."}
		return
	}

	// Calculate overall risk based on predictions
	totalProbability := 0
	highRiskCount := 0

	for _, prediction := range report.Predictions {
		totalProbability += prediction.Probability
		if prediction.Probability >= 70 {
			highRiskCount++
		}
	}

	averageProbability := totalProbability / len(report.Predictions)
	report.Confidence = averageProbability

	if highRiskCount > 0 || averageProbability >= 70 {
		report.OverallRisk = "High"
		report.Recommendations = []string{
			"Immediate action required. Review container configurations.",
			"Consider increasing resource limits or debugging application issues.",
			"Set up additional monitoring and alerts.",
		}
	} else if averageProbability >= 50 {
		report.OverallRisk = "Medium"
		report.Recommendations = []string{
			"Monitor closely and address resource configuration issues.",
			"Review application logs for patterns.",
			"Consider proactive maintenance.",
		}
	} else {
		report.OverallRisk = "Low"
		report.Recommendations = []string{
			"Continue normal monitoring operations.",
			"Review recommendations during next maintenance window.",
		}
	}
}
