package diagnostics

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EventsAnalyzer provides analysis for Kubernetes events
type EventsAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewEventsAnalyzer creates a new EventsAnalyzer
func NewEventsAnalyzer(client kubernetes.Interface, namespace string) *EventsAnalyzer {
	return &EventsAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// EventAnalysis contains the analysis of events
type EventAnalysis struct {
	TotalEvents   int
	WarningEvents []corev1.Event
	NormalEvents  []corev1.Event
	RecentEvents  []corev1.Event
	Issues        []string
}

// AnalyzeEvents analyzes events for a specific resource
func (e *EventsAnalyzer) AnalyzeEvents(resourceName string, resourceType string) (*EventAnalysis, error) {
	// Get events for the resource
	events, err := e.client.CoreV1().Events(e.namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=%s", resourceName, resourceType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get events for %s %s: %v", resourceType, resourceName, err)
	}

	analysis := &EventAnalysis{
		TotalEvents: len(events.Items),
	}

	// Categorize events
	now := time.Now()
	for _, event := range events.Items {
		// Check if event is recent (last 24 hours)
		if event.LastTimestamp.Time.After(now.Add(-24 * time.Hour)) {
			analysis.RecentEvents = append(analysis.RecentEvents, event)
		}

		// Categorize by type
		if event.Type == corev1.EventTypeWarning {
			analysis.WarningEvents = append(analysis.WarningEvents, event)
		} else {
			analysis.NormalEvents = append(analysis.NormalEvents, event)
		}
	}

	// Generate issues from warning events
	for _, event := range analysis.WarningEvents {
		if event.LastTimestamp.Time.After(now.Add(-1 * time.Hour)) {
			analysis.Issues = append(analysis.Issues,
				fmt.Sprintf("Recent warning: %s - %s", event.Reason, event.Message))
		}
	}

	return analysis, nil
}

// AnalyzeNamespaceEvents analyzes all events in a namespace
func (e *EventsAnalyzer) AnalyzeNamespaceEvents() (*EventAnalysis, error) {
	events, err := e.client.CoreV1().Events(e.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get events for namespace %s: %v", e.namespace, err)
	}

	analysis := &EventAnalysis{
		TotalEvents: len(events.Items),
	}

	// Categorize events
	now := time.Now()
	for _, event := range events.Items {
		// Check if event is recent (last 24 hours)
		if event.LastTimestamp.Time.After(now.Add(-24 * time.Hour)) {
			analysis.RecentEvents = append(analysis.RecentEvents, event)
		}

		// Categorize by type
		if event.Type == corev1.EventTypeWarning {
			analysis.WarningEvents = append(analysis.WarningEvents, event)
		} else {
			analysis.NormalEvents = append(analysis.NormalEvents, event)
		}
	}

	// Generate issues from recent warning events
	for _, event := range analysis.WarningEvents {
		if event.LastTimestamp.Time.After(now.Add(-1 * time.Hour)) {
			analysis.Issues = append(analysis.Issues,
				fmt.Sprintf("[%s] %s: %s - %s",
					event.InvolvedObject.Kind,
					event.InvolvedObject.Name,
					event.Reason,
					event.Message))
		}
	}

	return analysis, nil
}
