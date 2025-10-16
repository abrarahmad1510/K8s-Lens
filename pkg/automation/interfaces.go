package automation

import (
	"context"
	"time"
)

// RemediationAction represents a specific remediation action
type RemediationAction struct {
	Type        string
	Description string
	Command     string
	Risk        string // low, medium, high
}

// Remediator defines the interface for automated remediation
type Remediator interface {
	CanFix(issueType string) bool
	Remediate(ctx context.Context, resource, namespace string) (*RemediationResult, error)
	GetSupportedIssues() []string
	GetRemediationActions() []RemediationAction
}

// Scaler defines the interface for predictive scaling
type Scaler interface {
	CanScale(resource string) bool
	PredictScale(ctx context.Context, deployment, namespace string) (*ScaleRecommendation, error)
	GetScalingStrategies() []string
}

// Healer defines the interface for self-healing
type Healer interface {
	CanHeal(resource string) bool
	Heal(ctx context.Context, resource, namespace string) (*RemediationResult, error)
	GetHealingCapabilities() []string
}

// ScaleRecommendation represents a scaling recommendation
type ScaleRecommendation struct {
	Resource      string
	Namespace     string
	CurrentReplicas int32
	RecommendedReplicas int32
	Confidence    float64 // 0.0 to 1.0
	Reason        string
	Metrics       map[string]float64
	Timestamp     time.Time
}
