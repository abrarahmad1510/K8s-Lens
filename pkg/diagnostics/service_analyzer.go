package diagnostics

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceAnalyzer provides analysis for Service resources
type ServiceAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewServiceAnalyzer creates a new ServiceAnalyzer
func NewServiceAnalyzer(client kubernetes.Interface, namespace string) *ServiceAnalyzer {
	return &ServiceAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// ServiceReport contains the analysis report for a Service
type ServiceReport struct {
	Name       string
	Namespace  string
	Type       corev1.ServiceType
	ClusterIP  string
	ExternalIP string
	Ports      []corev1.ServicePort
	Selector   map[string]string
	Endpoints  *corev1.Endpoints
	Events     []corev1.Event
	Analysis   ServiceAnalysis
}

// ServiceAnalysis contains diagnostic results
type ServiceAnalysis struct {
	Status          string
	Issues          []string
	Recommendations []string
}

// Analyze performs the analysis of a Service
func (s *ServiceAnalyzer) Analyze(serviceName string) (*ServiceReport, error) {
	// Get the service
	service, err := s.client.CoreV1().Services(s.namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %v", serviceName, err)
	}

	// Get endpoints
	endpoints, err := s.client.CoreV1().Endpoints(s.namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints for service %s: %v", serviceName, err)
	}

	// Get events
	events, err := s.client.CoreV1().Events(s.namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + serviceName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get events for service %s: %v", serviceName, err)
	}

	report := &ServiceReport{
		Name:       service.Name,
		Namespace:  service.Namespace,
		Type:       service.Spec.Type,
		ClusterIP:  service.Spec.ClusterIP,
		ExternalIP: s.getExternalIP(service),
		Ports:      service.Spec.Ports,
		Selector:   service.Spec.Selector,
		Endpoints:  endpoints,
		Events:     events.Items,
	}

	s.analyzeService(report)
	s.analyzeEndpoints(report)

	return report, nil
}

func (s *ServiceAnalyzer) getExternalIP(service *corev1.Service) string {
	if len(service.Status.LoadBalancer.Ingress) > 0 {
		if service.Status.LoadBalancer.Ingress[0].IP != "" {
			return service.Status.LoadBalancer.Ingress[0].IP
		}
		return service.Status.LoadBalancer.Ingress[0].Hostname
	}
	return ""
}

func (s *ServiceAnalyzer) analyzeService(report *ServiceReport) {
	// Check service type specific issues
	switch report.Type {
	case corev1.ServiceTypeLoadBalancer:
		if report.ExternalIP == "" {
			report.Analysis.Issues = append(report.Analysis.Issues,
				"LoadBalancer service has no external IP assigned")
		}
	case corev1.ServiceTypeClusterIP:
		if report.ClusterIP == "" {
			report.Analysis.Issues = append(report.Analysis.Issues,
				"ClusterIP service has no cluster IP assigned")
		}
	case corev1.ServiceTypeNodePort:
		if len(report.Ports) == 0 {
			report.Analysis.Issues = append(report.Analysis.Issues,
				"NodePort service has no ports configured")
		}
	}

	// Check selector
	if len(report.Selector) == 0 {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"Service has no selector configured")
	}

	// Check ports
	if len(report.Ports) == 0 {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"Service has no ports configured")
	}

	if len(report.Analysis.Issues) == 0 {
		report.Analysis.Status = "Healthy"
	} else {
		report.Analysis.Status = "Unhealthy"
	}
}

func (s *ServiceAnalyzer) analyzeEndpoints(report *ServiceReport) {
	if report.Endpoints == nil {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"No endpoints found for service")
		return
	}

	totalAddresses := 0
	for _, subset := range report.Endpoints.Subsets {
		totalAddresses += len(subset.Addresses)
	}

	if totalAddresses == 0 {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"Service has no active endpoints")
		report.Analysis.Recommendations = append(report.Analysis.Recommendations,
			"Check if pods matching the selector are running and ready")
	} else {
		report.Analysis.Recommendations = append(report.Analysis.Recommendations,
			fmt.Sprintf("Service has %d active endpoint(s)", totalAddresses))
	}
}
