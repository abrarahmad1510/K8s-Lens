package diagnostics

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EndpointAnalyzer provides analysis for Service endpoints
type EndpointAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewEndpointAnalyzer creates a new EndpointAnalyzer
func NewEndpointAnalyzer(client kubernetes.Interface, namespace string) *EndpointAnalyzer {
	return &EndpointAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// EndpointReport contains the analysis report
type EndpointReport struct {
	ServiceName string
	Namespace   string
	Endpoints   *corev1.Endpoints
	Pods        []corev1.Pod
	Analysis    EndpointAnalysis
}

// EndpointAnalysis contains diagnostic results
type EndpointAnalysis struct {
	Status          string
	Issues          []string
	Recommendations []string
	ReadyPods       int
	TotalPods       int
}

// ValidateEndpoints analyzes endpoints for a service
func (e *EndpointAnalyzer) ValidateEndpoints(serviceName string) (*EndpointReport, error) {
	// Get endpoints
	endpoints, err := e.client.CoreV1().Endpoints(e.namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints for service %s: %v", serviceName, err)
	}

	// Get the service to find selector
	service, err := e.client.CoreV1().Services(e.namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %v", serviceName, err)
	}

	// Get pods matching the service selector
	var pods []corev1.Pod
	if len(service.Spec.Selector) > 0 {
		labelSelector := metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: service.Spec.Selector,
		})

		podList, err := e.client.CoreV1().Pods(e.namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get pods for service %s: %v", serviceName, err)
		}
		pods = podList.Items
	}

	report := &EndpointReport{
		ServiceName: serviceName,
		Namespace:   e.namespace,
		Endpoints:   endpoints,
		Pods:        pods,
	}

	e.analyzeEndpoints(report)
	e.analyzePodReadiness(report)

	return report, nil
}

func (e *EndpointAnalyzer) analyzeEndpoints(report *EndpointReport) {
	if report.Endpoints == nil {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"No endpoints object found for service")
		return
	}

	totalAddresses := 0
	for _, subset := range report.Endpoints.Subsets {
		totalAddresses += len(subset.Addresses)
	}

	if totalAddresses == 0 {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"Service has no active endpoints")
	} else {
		report.Analysis.Recommendations = append(report.Analysis.Recommendations,
			fmt.Sprintf("Service has %d active endpoint(s)", totalAddresses))
	}
}

func (e *EndpointAnalyzer) analyzePodReadiness(report *EndpointReport) {
	readyPods := 0
	totalPods := len(report.Pods)

	for _, pod := range report.Pods {
		if isPodReady(&pod) {
			readyPods++
		}
	}

	report.Analysis.ReadyPods = readyPods
	report.Analysis.TotalPods = totalPods

	if totalPods == 0 {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"No pods found matching service selector")
	} else if readyPods == 0 {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"No pods are ready to serve traffic")
	} else if readyPods < totalPods {
		report.Analysis.Issues = append(report.Analysis.Issues,
			fmt.Sprintf("Only %d of %d pods are ready", readyPods, totalPods))
	}

	if len(report.Analysis.Issues) == 0 {
		report.Analysis.Status = "Healthy"
	} else {
		report.Analysis.Status = "Unhealthy"
	}
}

func isPodReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
