package remediators

import (
	"context"
	"fmt"
	"time"

	"github.com/abrarahmad1510/k8s-lens/pkg/automation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodRestartRemediator automatically restarts pods with issues
type PodRestartRemediator struct {
	client kubernetes.Interface
}

// NewPodRestartRemediator creates a new pod restart remediator
func NewPodRestartRemediator(client kubernetes.Interface) *PodRestartRemediator {
	return &PodRestartRemediator{
		client: client,
	}
}

// CanFix checks if this remediator can fix the given issue type
func (p *PodRestartRemediator) CanFix(issueType string) bool {
	supportedIssues := []string{
		"CrashLoopBackOff",
		"ImagePullBackOff",
		"ErrImagePull", 
		"RunContainerError",
		"PodStuckTerminating",
		"HighRestartCount",
	}
	
	for _, issue := range supportedIssues {
		if issue == issueType {
			return true
		}
	}
	return false
}

// Remediate attempts to fix the pod issue by restarting it
func (p *PodRestartRemediator) Remediate(ctx context.Context, resource, namespace string) (*automation.RemediationResult, error) {
	startTime := time.Now()
	
	// Delete the pod to trigger restart (Deployment will recreate it)
	err := p.client.CoreV1().Pods(namespace).Delete(ctx, resource, metav1.DeleteOptions{})
	if err != nil {
		return &automation.RemediationResult{
			Success:  false,
			Action:   "restart",
			Resource: resource,
			Message:  fmt.Sprintf("Failed to delete pod: %v", err),
			Duration: time.Since(startTime),
		}, err
	}

	return &automation.RemediationResult{
		Success:  true,
		Action:   "restart",
		Resource: resource,
		Message:  fmt.Sprintf("Successfully restarted pod %s in namespace %s", resource, namespace),
		Duration: time.Since(startTime),
	}, nil
}

// GetSupportedIssues returns the types of issues this remediator can fix
func (p *PodRestartRemediator) GetSupportedIssues() []string {
	return []string{
		"CrashLoopBackOff",
		"ImagePullBackOff", 
		"ErrImagePull",
		"RunContainerError",
		"PodStuckTerminating",
		"HighRestartCount",
	}
}

// GetRemediationActions returns available remediation actions
func (p *PodRestartRemediator) GetRemediationActions() []automation.RemediationAction {
	return []automation.RemediationAction{
		{
			Type:        "PodRestart",
			Description: "Restart the pod to resolve container issues",
			Command:     "kubectl delete pod <pod-name> -n <namespace>",
			Risk:        "low",
		},
	}
}
