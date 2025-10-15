package diagnostics

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AnalyzePodEvents Examines Recent Pod Events
func (a *ResourceAnalyzer) analyzePodEvents(pod *corev1.Pod, namespace string, result *AnalysisResult) string {
	var eventsStr strings.Builder

	eventsStr.WriteString("Recent Events Analysis:\n")

	events, err := a.client.Clientset.CoreV1().Events(namespace).List(a.ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", pod.Name),
	})

	if err != nil || len(events.Items) == 0 {
		eventsStr.WriteString("  Status: No Recent Events Found\n\n")
		return eventsStr.String()
	}

	// Show Last 5 Events
	eventCount := 0
	for i := len(events.Items) - 1; i >= 0 && eventCount < 5; i-- {
		event := events.Items[i]
		if time.Since(event.LastTimestamp.Time).Hours() > 24 {
			continue // Skip Events Older Than 24 Hours
		}

		eventType := "Info"
		if event.Type == "Warning" {
			eventType = "Warning"
			result.Warnings = append(result.Warnings, fmt.Sprintf("Event: %s", event.Message))
		}

		eventsStr.WriteString(fmt.Sprintf("  [%s] %s: %s - %s\n",
			event.LastTimestamp.Format("15:04:05"), eventType, event.Reason, event.Message))
		eventCount++
	}

	eventsStr.WriteString("\n")
	return eventsStr.String()
}

// AnalyzePodResources Checks Resource Configuration
func (a *ResourceAnalyzer) analyzePodResources(pod *corev1.Pod, result *AnalysisResult) string {
	var resources strings.Builder

	resources.WriteString("Resource Analysis:\n")

	hasLimits := false
	hasRequests := false

	for _, container := range pod.Spec.Containers {
		if container.Resources.Limits != nil {
			hasLimits = true
		}
		if container.Resources.Requests != nil {
			hasRequests = true
		}
	}

	if !hasLimits {
		result.Warnings = append(result.Warnings, "No Resource Limits Set")
		resources.WriteString("  Warning: No Resource Limits Configured\n")
	} else {
		resources.WriteString("  Status: Resource Limits Configured\n")
	}

	if !hasRequests {
		result.Warnings = append(result.Warnings, "No Resource Requests Set")
		resources.WriteString("  Warning: No Resource Requests Configured\n")
	} else {
		resources.WriteString("  Status: Resource Requests Configured\n")
	}

	resources.WriteString("\n")
	return resources.String()
}

// GenerateSummary Creates Final Analysis Summary
func (a *ResourceAnalyzer) generateSummary(result *AnalysisResult) string {
	var summary strings.Builder

	summary.WriteString("Summary And Recommendations:\n")
	summary.WriteString(fmt.Sprintf("  Overall Health: "))

	if result.Healthy {
		summary.WriteString("Healthy\n")
	} else {
		summary.WriteString("Needs Attention\n")
	}

	if len(result.Errors) > 0 {
		summary.WriteString("  Critical Issues:\n")
		for _, err := range result.Errors {
			summary.WriteString(fmt.Sprintf("    • %s\n", err))
		}
	}

	if len(result.Warnings) > 0 {
		summary.WriteString("  Warnings:\n")
		for _, warning := range result.Warnings {
			summary.WriteString(fmt.Sprintf("    • %s\n", warning))
		}
	}

	if len(result.Recommendations) > 0 {
		summary.WriteString("  Recommended Actions:\n")
		for _, rec := range result.Recommendations {
			summary.WriteString(fmt.Sprintf("    • %s\n", rec))
		}
	} else if result.Healthy {
		summary.WriteString("  Status: No Actions Needed - Everything Looks Good\n")
	}

	return summary.String()
}
