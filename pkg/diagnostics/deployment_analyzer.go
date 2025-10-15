package diagnostics

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeploymentAnalyzer provides analysis for Deployment resources
type DeploymentAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewDeploymentAnalyzer creates a new DeploymentAnalyzer
func NewDeploymentAnalyzer(client kubernetes.Interface, namespace string) *DeploymentAnalyzer {
	return &DeploymentAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// DeploymentReport contains the analysis report for a Deployment
type DeploymentReport struct {
	Name              string
	Namespace         string
	DesiredReplicas   int32
	CurrentReplicas   int32
	ReadyReplicas     int32
	AvailableReplicas int32
	UpdatedReplicas   int32
	Conditions        []appsv1.DeploymentCondition
	PodTemplate       corev1.PodTemplateSpec
	ReplicaSets       []appsv1.ReplicaSet
	Events            []corev1.Event
	Analysis          DeploymentAnalysis
}

// DeploymentAnalysis contains diagnostic results
type DeploymentAnalysis struct {
	Status          string
	Issues          []string
	Recommendations []string
	RolloutStatus   string
}

// Analyze performs the analysis of a Deployment
func (d *DeploymentAnalyzer) Analyze(deploymentName string) (*DeploymentReport, error) {
	// Get deployment
	deployment, err := d.client.AppsV1().Deployments(d.namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s: %v", deploymentName, err)
	}

	// Get related ReplicaSets
	rsList, err := d.client.AppsV1().ReplicaSets(d.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list replicasets for deployment %s: %v", deploymentName, err)
	}

	// Get events
	events, err := d.client.CoreV1().Events(d.namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + deploymentName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get events for deployment %s: %v", deploymentName, err)
	}

	report := &DeploymentReport{
		Name:              deployment.Name,
		Namespace:         deployment.Namespace,
		DesiredReplicas:   *deployment.Spec.Replicas,
		CurrentReplicas:   deployment.Status.Replicas,
		ReadyReplicas:     deployment.Status.ReadyReplicas,
		AvailableReplicas: deployment.Status.AvailableReplicas,
		UpdatedReplicas:   deployment.Status.UpdatedReplicas,
		Conditions:        deployment.Status.Conditions,
		PodTemplate:       deployment.Spec.Template,
		ReplicaSets:       rsList.Items,
		Events:            events.Items,
	}

	d.analyzeConditions(report)
	d.analyzeReplicaSets(report)
	d.analyzeRolloutStatus(report)

	return report, nil
}

func (d *DeploymentAnalyzer) analyzeConditions(report *DeploymentReport) {
	for _, condition := range report.Conditions {
		switch condition.Type {
		case appsv1.DeploymentAvailable:
			if condition.Status == corev1.ConditionFalse {
				report.Analysis.Issues = append(report.Analysis.Issues,
					fmt.Sprintf("Deployment not available: %s", condition.Message))
			}
		case appsv1.DeploymentProgressing:
			if condition.Status == corev1.ConditionFalse {
				report.Analysis.Issues = append(report.Analysis.Issues,
					fmt.Sprintf("Deployment not progressing: %s", condition.Message))
			}
		}
	}

	if report.ReadyReplicas != report.DesiredReplicas {
		report.Analysis.Issues = append(report.Analysis.Issues,
			fmt.Sprintf("Ready replicas (%d) does not match desired replicas (%d)",
				report.ReadyReplicas, report.DesiredReplicas))
	}

	if len(report.Analysis.Issues) == 0 {
		report.Analysis.Status = "Healthy"
	} else {
		report.Analysis.Status = "Unhealthy"
	}
}

func (d *DeploymentAnalyzer) analyzeReplicaSets(report *DeploymentReport) {
	for _, rs := range report.ReplicaSets {
		if *rs.Spec.Replicas > 0 && rs.Status.Replicas > 0 {
			if rs.CreationTimestamp.Time.Before(time.Now().Add(-24 * time.Hour)) {
				report.Analysis.Issues = append(report.Analysis.Issues,
					fmt.Sprintf("Old ReplicaSet %s still has %d replicas", rs.Name, rs.Status.Replicas))
			}
		}
	}
}

func (d *DeploymentAnalyzer) analyzeRolloutStatus(report *DeploymentReport) {
	if report.UpdatedReplicas == report.DesiredReplicas &&
		report.ReadyReplicas == report.DesiredReplicas {
		report.Analysis.RolloutStatus = "Complete"
	} else if report.UpdatedReplicas < report.DesiredReplicas {
		report.Analysis.RolloutStatus = "Progressing"
	} else {
		report.Analysis.RolloutStatus = "Degraded"
	}
}
