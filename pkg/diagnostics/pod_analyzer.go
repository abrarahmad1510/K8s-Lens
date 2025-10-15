package diagnostics

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodAnalyzer provides analysis for Pod resources
type PodAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewPodAnalyzer creates a new PodAnalyzer
func NewPodAnalyzer(client kubernetes.Interface, namespace string) *PodAnalyzer {
	return &PodAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// PodReport contains the analysis report for a Pod
type PodReport struct {
	Name                string
	Namespace           string
	UID                 string
	Phase               string
	Node                string
	PodIP               string
	ServiceAccount      string
	Created             time.Time
	Status              string
	Containers          []ContainerStatus
	Events              []corev1.Event
	Issues              []string
	Recommendations     []string
	ResourceLimitsSet   bool
	ResourceRequestsSet bool
	RestartCount        int32
}

// ContainerStatus represents the status of a container
type ContainerStatus struct {
	Name    string
	Image   string
	Status  string
	Ready   bool
	Reason  string
	Message string
}

// Analyze performs the analysis of a Pod
func (p *PodAnalyzer) Analyze(podName string) (*PodReport, error) {
	// Get the pod
	pod, err := p.client.CoreV1().Pods(p.namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %v", podName, err)
	}

	// Get events for the pod
	events, err := p.client.CoreV1().Events(p.namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + podName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get events for pod %s: %v", podName, err)
	}

	report := &PodReport{
		Name:           pod.Name,
		Namespace:      pod.Namespace,
		UID:            string(pod.UID),
		Phase:          string(pod.Status.Phase),
		Node:           pod.Spec.NodeName,
		PodIP:          pod.Status.PodIP,
		ServiceAccount: pod.Spec.ServiceAccountName,
		Created:        pod.CreationTimestamp.Time,
		Events:         events.Items,
	}

	// Analyze container statuses
	p.analyzeContainers(report, pod)

	// Analyze resource configuration
	p.analyzeResources(report, pod)

	// Generate recommendations
	p.generateRecommendations(report)

	return report, nil
}

func (p *PodAnalyzer) analyzeContainers(report *PodReport, pod *corev1.Pod) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		container := ContainerStatus{
			Name:  containerStatus.Name,
			Image: containerStatus.Image,
			Ready: containerStatus.Ready,
		}

		// Determine container status
		if containerStatus.State.Running != nil {
			container.Status = "Running"
		} else if containerStatus.State.Waiting != nil {
			container.Status = fmt.Sprintf("Waiting - %s: %s",
				containerStatus.State.Waiting.Reason,
				containerStatus.State.Waiting.Message)
			container.Reason = containerStatus.State.Waiting.Reason
			container.Message = containerStatus.State.Waiting.Message
		} else if containerStatus.State.Terminated != nil {
			container.Status = fmt.Sprintf("Terminated - %s: %s",
				containerStatus.State.Terminated.Reason,
				containerStatus.State.Terminated.Message)
			container.Reason = containerStatus.State.Terminated.Reason
			container.Message = containerStatus.State.Terminated.Message
		}

		report.Containers = append(report.Containers, container)
		report.RestartCount += containerStatus.RestartCount

		// Check for issues
		if containerStatus.State.Waiting != nil {
			switch containerStatus.State.Waiting.Reason {
			case "ImagePullBackOff", "ErrImagePull":
				report.Issues = append(report.Issues,
					fmt.Sprintf("Container %s cannot pull image: %s",
						containerStatus.Name, containerStatus.State.Waiting.Message))
			case "CrashLoopBackOff":
				report.Issues = append(report.Issues,
					fmt.Sprintf("Container %s is crashing: %s",
						containerStatus.Name, containerStatus.State.Waiting.Message))
			}
		}
	}

	// Determine overall pod status
	if pod.Status.Phase == corev1.PodRunning {
		allReady := true
		for _, container := range report.Containers {
			if !container.Ready {
				allReady = false
				break
			}
		}
		if allReady {
			report.Status = "Running"
		} else {
			report.Status = "Running but not all containers ready"
		}
	} else {
		report.Status = string(pod.Status.Phase)
	}
}

func (p *PodAnalyzer) analyzeResources(report *PodReport, pod *corev1.Pod) {
	report.ResourceLimitsSet = true
	report.ResourceRequestsSet = true

	for _, container := range pod.Spec.Containers {
		if container.Resources.Limits == nil || len(container.Resources.Limits) == 0 {
			report.ResourceLimitsSet = false
		}
		if container.Resources.Requests == nil || len(container.Resources.Requests) == 0 {
			report.ResourceRequestsSet = false
		}
	}
}

func (p *PodAnalyzer) generateRecommendations(report *PodReport) {
	if !report.ResourceLimitsSet {
		report.Recommendations = append(report.Recommendations,
			"Add resource limits to prevent OOM kills and ensure quality of service")
	}

	if !report.ResourceRequestsSet {
		report.Recommendations = append(report.Recommendations,
			"Add resource requests to help the scheduler make better placement decisions")
	}

	if report.RestartCount > 5 {
		report.Recommendations = append(report.Recommendations,
			fmt.Sprintf("Investigate why container has restarted %d times", report.RestartCount))
	}

	// Check for common issues in events
	for _, event := range report.Events {
		if event.Type == "Warning" {
			switch event.Reason {
			case "FailedScheduling":
				report.Recommendations = append(report.Recommendations,
					"Check node resources and affinity rules")
			case "FailedMount":
				report.Recommendations = append(report.Recommendations,
					"Verify volume configurations and storage class availability")
			}
		}
	}
}
