package automation

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
)

// AutomationEngine provides self-healing and automated remediation
type AutomationEngine struct {
	client    kubernetes.Interface
	remediators []Remediator
	scalers    []Scaler
	healers    []Healer
}

// NewAutomationEngine creates a new automation engine
func NewAutomationEngine(client kubernetes.Interface) *AutomationEngine {
	return &AutomationEngine{
		client:    client,
		remediators: []Remediator{},
		scalers:    []Scaler{},
		healers:    []Healer{},
	}
}

// RemediationResult represents the outcome of an automated fix
type RemediationResult struct {
	Success    bool
	Action     string
	Resource   string
	Message    string
	Duration   time.Duration
}

// RegisterRemediator adds a new remediation capability
func (a *AutomationEngine) RegisterRemediator(remediator Remediator) {
	a.remediators = append(a.remediators, remediator)
}

// RegisterScaler adds a new scaling capability
func (a *AutomationEngine) RegisterScaler(scaler Scaler) {
	a.scalers = append(a.scalers, scaler)
}

// RegisterHealer adds a new healing capability
func (a *AutomationEngine) RegisterHealer(healer Healer) {
	a.healers = append(a.healers, healer)
}

// AutoRemediate attempts to automatically fix detected issues
func (a *AutomationEngine) AutoRemediate(ctx context.Context, issueType, resource, namespace string) (*RemediationResult, error) {
	for _, remediator := range a.remediators {
		if remediator.CanFix(issueType) {
			return remediator.Remediate(ctx, resource, namespace)
		}
	}
	
	return &RemediationResult{
		Success:  false,
		Action:   "none",
		Resource: resource,
		Message:  fmt.Sprintf("No automediation available for issue type: %s", issueType),
	}, nil
}

// PredictiveScale analyzes metrics and suggests scaling actions
func (a *AutomationEngine) PredictiveScale(ctx context.Context, deployment, namespace string) (*ScaleRecommendation, error) {
	for _, scaler := range a.scalers {
		if scaler.CanScale(deployment) {
			return scaler.PredictScale(ctx, deployment, namespace)
		}
	}
	
	return nil, fmt.Errorf("no predictive scaling available for %s", deployment)
}

// SelfHeal attempts to automatically heal failing resources
func (a *AutomationEngine) SelfHeal(ctx context.Context, resource, namespace string) (*RemediationResult, error) {
	for _, healer := range a.healers {
		if healer.CanHeal(resource) {
			return healer.Heal(ctx, resource, namespace)
		}
	}
	
	return &RemediationResult{
		Success:  false,
		Action:   "none",
		Resource: resource,
		Message:  fmt.Sprintf("No self-healing available for resource: %s", resource),
	}, nil
}
