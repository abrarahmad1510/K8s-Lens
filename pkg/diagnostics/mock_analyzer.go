package diagnostics

import (
	"fmt"
	"time"
)

// MockAnalyzer provides mock analysis for testing
type MockAnalyzer struct{}

// NewMockAnalyzer creates a new MockAnalyzer
func NewMockAnalyzer() *MockAnalyzer {
	return &MockAnalyzer{}
}

// Analyze runs a mock analysis
func (m *MockAnalyzer) Analyze(podName string) string {
	return fmt.Sprintf(`
K8s Lens Mock Analysis Report For Pod: %s
---
Pod Status Analysis:
Phase: Running
Node: minikube-mock
Created: %s
Status: Pod Is Running Normally

Container Status Analysis:
Container: nginx
Image: nginx:latest
Status: Running For 2h35m
Status: Container Is Ready

Resource Analysis:
Warning: No Resource Limits Configured
Status: Resource Requests Configured

Recent Events Analysis:
Status: No Recent Events Found

Summary And Recommendations:
Overall Health: Needs Attention
Warnings:
• No Resource Limits Set
Recommended Actions:
• Add Resource Limits To Prevent OOM Kills
• Monitor Container Restart Patterns
• Consider Adding Liveness And Readiness Probes
`, podName, time.Now().Format("Mon, 15 Oct 2024 14:32:15 UTC"))
}
