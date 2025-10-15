package integration

import (
	"testing"

	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeploymentAnalysis(t *testing.T) {
	// Create fake client
	client := fake.NewSimpleClientset()

	analyzer := diagnostics.NewDeploymentAnalyzer(client, "default")

	// Test with non-existent deployment
	_, err := analyzer.Analyze("test-deployment")
	if err == nil {
		t.Error("Expected error for non-existent deployment, got nil")
	}
}
