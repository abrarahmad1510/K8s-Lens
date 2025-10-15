package diagnostics

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AnalyzePod Performs Comprehensive Pod Analysis
func (a *ResourceAnalyzer) AnalyzePod(podName, namespace string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Healthy:         true,
		Confidence:      0.8,
		Recommendations: []string{},
		Warnings:        []string{},
		Errors:          []string{},
	}

	// Get Pod Information From Kubernetes API
	pod, err := a.client.Clientset.CoreV1().Pods(namespace).Get(a.ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed To Get Pod %s: %v", podName, err)
	}

	var report strings.Builder

	// Report Header
	report.WriteString(fmt.Sprintf("K8s Lens Analysis Report For Pod: %s\n", podName))
	report.WriteString(strings.Repeat("=", 60) + "\n\n")

	// Pod Status Analysis
	podStatus := a.analyzePodStatus(pod, result)
	report.WriteString(podStatus)

	// Container Analysis
	containerStatus := a.analyzeContainerStatuses(pod, result)
	report.WriteString(containerStatus)

	// Resource Analysis
	resourceStatus := a.analyzePodResources(pod, result)
	report.WriteString(resourceStatus)

	// Event Analysis
	eventStatus := a.analyzePodEvents(pod, namespace, result)
	report.WriteString(eventStatus)

	// Generate Summary
	report.WriteString(a.generateSummary(result))

	result.Report = report.String()
	return result, nil
}

// AnalyzePodStatus Analyzes Pod Phase And Conditions
func (a *ResourceAnalyzer) analyzePodStatus(pod *corev1.Pod, result *AnalysisResult) string {
	var status strings.Builder

	status.WriteString("Pod Status Analysis:\n")
	status.WriteString(fmt.Sprintf("  Phase: %s\n", pod.Status.Phase))
	status.WriteString(fmt.Sprintf("  Node: %s\n", pod.Spec.NodeName))
	status.WriteString(fmt.Sprintf("  Created: %s\n", pod.CreationTimestamp.Format(time.RFC1123)))

	switch pod.Status.Phase {
	case corev1.PodPending:
		result.Healthy = false
		result.Warnings = append(result.Warnings, "Pod Is Stuck In Pending State")
		status.WriteString("  Warning: Pod Is Pending - Check Resource Availability\n")
	case corev1.PodFailed:
		result.Healthy = false
		result.Errors = append(result.Errors, "Pod Has Failed")
		status.WriteString("  Critical: Pod Has Failed\n")
	case corev1.PodRunning:
		status.WriteString("  Status: Pod Is Running Normally\n")
	case corev1.PodSucceeded:
		status.WriteString("  Status: Pod Completed Successfully\n")
	case corev1.PodUnknown:
		result.Healthy = false
		result.Warnings = append(result.Warnings, "Pod Status Is Unknown")
		status.WriteString("  Warning: Pod Status Is Unknown\n")
	}

	status.WriteString("\n")
	return status.String()
}

// AnalyzeContainerStatuses Checks Container Health And Readiness
func (a *ResourceAnalyzer) analyzeContainerStatuses(pod *corev1.Pod, result *AnalysisResult) string {
	var status strings.Builder

	status.WriteString("Container Status Analysis:\n")

	for _, container := range pod.Spec.Containers {
		status.WriteString(fmt.Sprintf("  Container: %s\n", container.Name))
		status.WriteString(fmt.Sprintf("    Image: %s\n", container.Image))

		// Find Container Status
		var containerStatus *corev1.ContainerStatus
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Name == container.Name {
				containerStatus = &cs
				break
			}
		}

		if containerStatus == nil {
			status.WriteString("    Status: Not Available\n")
			continue
		}

		// Check Container State
		if containerStatus.State.Waiting != nil {
			waiting := containerStatus.State.Waiting
			status.WriteString(fmt.Sprintf("    Status: Waiting - %s: %s\n", waiting.Reason, waiting.Message))

			// Provide Intelligent Recommendations Based On Waiting Reason
			switch waiting.Reason {
			case "ImagePullBackOff", "ErrImagePull":
				result.Recommendations = append(result.Recommendations,
					fmt.Sprintf("Check Image Availability: %s", container.Image))
				result.Recommendations = append(result.Recommendations,
					"Verify Image Pull Secrets Are Configured")
			case "CrashLoopBackOff":
				result.Recommendations = append(result.Recommendations,
					fmt.Sprintf("Check Container Logs: kubectl logs %s -c %s", pod.Name, container.Name))
				result.Recommendations = append(result.Recommendations,
					"Verify Application Configuration And Dependencies")
			case "ContainerCreating":
				result.Recommendations = append(result.Recommendations,
					"Container Is Still Being Created - Check Node Resources")
			}
		}

		if containerStatus.State.Running != nil {
			runningTime := time.Since(containerStatus.State.Running.StartedAt.Time).Round(time.Second)
			status.WriteString(fmt.Sprintf("    Status: Running For %s\n", runningTime))
		}

		if containerStatus.State.Terminated != nil {
			terminated := containerStatus.State.Terminated
			status.WriteString(fmt.Sprintf("    Status: Terminated - %s (Exit Code: %d)\n",
				terminated.Reason, terminated.ExitCode))
		}

		// Check Readiness
		if !containerStatus.Ready {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Container %s Is Not Ready", container.Name))
			status.WriteString("    Warning: Container Is Not Ready\n")
		} else {
			status.WriteString("    Status: Container Is Ready\n")
		}

		status.WriteString("\n")
	}

	return status.String()
}
