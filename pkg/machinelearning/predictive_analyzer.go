package machinelearning

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PredictiveAnalyzer predicts potential future issues
type PredictiveAnalyzer struct {
	client kubernetes.Interface
}

// NewPredictiveAnalyzer creates a new predictive analyzer
func NewPredictiveAnalyzer(client kubernetes.Interface) *PredictiveAnalyzer {
	return &PredictiveAnalyzer{
		client: client,
	}
}

// PredictionReport contains predictive insights
type PredictionReport struct {
	Namespace   string
	Predictions []Prediction
	Confidence  float64
	TimeHorizon time.Duration
	GeneratedAt time.Time
}

// Prediction represents a single predictive insight
type Prediction struct {
	Type           string
	Resource       string
	Message        string
	Probability    float64
	ExpectedTime   time.Time
	Impact         string // Low, Medium, High, Critical
	Recommendation string
}

// PredictDeploymentFailures analyzes deployment for potential future issues
func (p *PredictiveAnalyzer) PredictDeploymentFailures(deploymentName, namespace string) (*PredictionReport, error) {
	report := &PredictionReport{
		Namespace:   namespace,
		GeneratedAt: time.Now(),
		TimeHorizon: 24 * time.Hour,
	}

	// Get deployment
	deployment, err := p.client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	// Get related pods
	pods, err := p.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %v", err)
	}

	predictions := p.analyzeDeploymentPatterns(deployment, pods.Items)
	report.Predictions = predictions
	report.Confidence = p.calculateOverallConfidence(predictions)

	return report, nil
}

func (p *PredictiveAnalyzer) analyzeDeploymentPatterns(deployment *appsv1.Deployment, pods []corev1.Pod) []Prediction {
	var predictions []Prediction

	// Analyze resource trends
	resourcePredictions := p.predictResourceIssues(deployment, pods)
	predictions = append(predictions, resourcePredictions...)

	// Analyze scaling patterns
	scalingPredictions := p.predictScalingIssues(deployment, pods)
	predictions = append(predictions, scalingPredictions...)

	// Analyze availability risks
	availabilityPredictions := p.predictAvailabilityRisks(deployment, pods)
	predictions = append(predictions, availabilityPredictions...)

	return predictions
}

func (p *PredictiveAnalyzer) predictResourceIssues(deployment *appsv1.Deployment, pods []corev1.Pod) []Prediction {
	var predictions []Prediction

	totalCPURequest := int64(0)
	totalMemoryRequest := int64(0)
	containerCount := 0

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				cpu := container.Resources.Requests[corev1.ResourceCPU]
				memory := container.Resources.Requests[corev1.ResourceMemory]

				totalCPURequest += cpu.MilliValue()
				totalMemoryRequest += memory.Value() / (1024 * 1024) // MB
				containerCount++
			}
		}
	}

	if containerCount > 0 {
		avgCPU := float64(totalCPURequest) / float64(containerCount)
		avgMemory := float64(totalMemoryRequest) / float64(containerCount)

		// Predict resource exhaustion if average usage is high
		if avgCPU > 800 { // 800 millicores average
			predictions = append(predictions, Prediction{
				Type:           "ResourceExhaustion",
				Resource:       deployment.Name,
				Message:        "High CPU usage may lead to resource exhaustion",
				Probability:    0.7,
				ExpectedTime:   time.Now().Add(12 * time.Hour),
				Impact:         "High",
				Recommendation: "Consider horizontal pod autoscaling or resource optimization",
			})
		}

		if avgMemory > 2048 { // 2GB average
			predictions = append(predictions, Prediction{
				Type:           "MemoryPressure",
				Resource:       deployment.Name,
				Message:        "High memory usage may cause OOM kills",
				Probability:    0.6,
				ExpectedTime:   time.Now().Add(24 * time.Hour),
				Impact:         "High",
				Recommendation: "Optimize memory usage or increase limits",
			})
		}
	}

	return predictions
}

func (p *PredictiveAnalyzer) predictScalingIssues(deployment *appsv1.Deployment, pods []corev1.Pod) []Prediction {
	var predictions []Prediction

	// Check if deployment is near resource limits
	if deployment.Spec.Replicas != nil {
		currentReplicas := *deployment.Spec.Replicas

		// Simple prediction: if we're at max replicas and have resource issues, predict scaling failure
		if currentReplicas >= 10 { // Arbitrary threshold
			predictions = append(predictions, Prediction{
				Type:           "ScalingLimit",
				Resource:       deployment.Name,
				Message:        "Deployment may reach scaling limits soon",
				Probability:    0.5,
				ExpectedTime:   time.Now().Add(48 * time.Hour),
				Impact:         "Medium",
				Recommendation: "Consider cluster autoscaling or application optimization",
			})
		}
	}

	return predictions
}

func (p *PredictiveAnalyzer) predictAvailabilityRisks(deployment *appsv1.Deployment, pods []corev1.Pod) []Prediction {
	var predictions []Prediction

	// Analyze pod distribution and availability risks
	nodeDistribution := make(map[string]int)
	for _, pod := range pods {
		nodeDistribution[pod.Spec.NodeName]++
	}

	// Check if pods are concentrated on few nodes
	if len(pods) > 0 {
		maxPodsPerNode := 0
		for _, count := range nodeDistribution {
			if count > maxPodsPerNode {
				maxPodsPerNode = count
			}
		}

		if maxPodsPerNode > len(pods)/2 { // More than half pods on one node
			predictions = append(predictions, Prediction{
				Type:           "AvailabilityRisk",
				Resource:       deployment.Name,
				Message:        "Pods concentrated on few nodes - node failure risk",
				Probability:    0.4,
				ExpectedTime:   time.Now().Add(72 * time.Hour),
				Impact:         "High",
				Recommendation: "Use pod anti-affinity to spread pods across nodes",
			})
		}
	}

	return predictions
}

func (p *PredictiveAnalyzer) calculateOverallConfidence(predictions []Prediction) float64 {
	if len(predictions) == 0 {
		return 1.0 // High confidence when no issues predicted
	}

	totalProbability := 0.0
	for _, prediction := range predictions {
		totalProbability += prediction.Probability
	}

	return 1.0 - (totalProbability / float64(len(predictions)) / 2.0)
}
