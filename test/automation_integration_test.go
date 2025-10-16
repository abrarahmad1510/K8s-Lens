package test

import (
	"testing"

	"github.com/abrarahmad1510/k8s-lens/pkg/automation"
	"github.com/abrarahmad1510/k8s-lens/pkg/automation/remediators"
	"github.com/stretchr/testify/assert"
)

// TestAutomationEngineCreation tests that the automation engine can be created
func TestAutomationEngineCreation(t *testing.T) {
	// This is a basic test to ensure the package structure is correct
	// In a real test, we would use a mocked Kubernetes client
	engine := automation.NewAutomationEngine(nil)
	assert.NotNil(t, engine, "Automation engine should be created successfully")
}

// TestPodRestartRemediator tests the pod restart remediator
func TestPodRestartRemediator(t *testing.T) {
	remediator := remediators.NewPodRestartRemediator(nil)
	
	// Test supported issues
	assert.True(t, remediator.CanFix("CrashLoopBackOff"), "Should support CrashLoopBackOff")
	assert.True(t, remediator.CanFix("ImagePullBackOff"), "Should support ImagePullBackOff")
	assert.False(t, remediator.CanFix("UnknownIssue"), "Should not support unknown issues")
	
	// Test actions
	actions := remediator.GetRemediationActions()
	assert.Greater(t, len(actions), 0, "Should have remediation actions")
	assert.Equal(t, "PodRestart", actions[0].Type, "First action should be PodRestart")
}

// TestRemediationInterfaces tests that interfaces are properly implemented
func TestRemediationInterfaces(t *testing.T) {
	var _ automation.Remediator = (*remediators.PodRestartRemediator)(nil)
}
