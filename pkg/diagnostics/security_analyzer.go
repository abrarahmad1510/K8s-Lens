package diagnostics

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SecurityAnalyzer provides security analysis for Pod resources
type SecurityAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewSecurityAnalyzer creates a new SecurityAnalyzer
func NewSecurityAnalyzer(client kubernetes.Interface, namespace string) *SecurityAnalyzer {
	return &SecurityAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// SecurityReport contains the security analysis report
type SecurityReport struct {
	PodName         string
	Namespace       string
	Analysis        SecurityAnalysis
	Issues          []SecurityIssue
	Warnings        []SecurityWarning
	Recommendations []string
}

// SecurityAnalysis contains security assessment results
type SecurityAnalysis struct {
	Status    string
	RiskLevel string
	Score     int
}

// SecurityIssue represents a security vulnerability
type SecurityIssue struct {
	Level       string
	Title       string
	Description string
	Remediation string
}

// SecurityWarning represents a security warning
type SecurityWarning struct {
	Level       string
	Title       string
	Description string
}

// AnalyzePodSecurity performs security analysis of a Pod
func (s *SecurityAnalyzer) AnalyzePodSecurity(podName string) (*SecurityReport, error) {
	pod, err := s.client.CoreV1().Pods(s.namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %v", podName, err)
	}

	report := &SecurityReport{
		PodName:   pod.Name,
		Namespace: pod.Namespace,
	}

	s.analyzeSecurityContext(report, pod)
	s.analyzeContainerSecurity(report, pod)
	s.calculateRiskScore(report)

	return report, nil
}

func (s *SecurityAnalyzer) analyzeSecurityContext(report *SecurityReport, pod *corev1.Pod) {
	// Analyze pod-level security context
	if pod.Spec.SecurityContext == nil {
		report.Issues = append(report.Issues, SecurityIssue{
			Level:       "High",
			Title:       "No Pod Security Context",
			Description: "Pod is running without any security context",
			Remediation: "Add securityContext with runAsNonRoot and seccompProfile",
		})
	} else {
		sc := pod.Spec.SecurityContext

		if sc.RunAsNonRoot == nil || !*sc.RunAsNonRoot {
			report.Issues = append(report.Issues, SecurityIssue{
				Level:       "High",
				Title:       "Running as Root",
				Description: "Pod may be running as root user",
				Remediation: "Set runAsNonRoot: true in securityContext",
			})
		}

		if sc.SeccompProfile == nil || sc.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
			report.Warnings = append(report.Warnings, SecurityWarning{
				Level:       "Medium",
				Title:       "No Seccomp Profile",
				Description: "Pod is not using runtime default seccomp profile",
			})
		}
	}
}

func (s *SecurityAnalyzer) analyzeContainerSecurity(report *SecurityReport, pod *corev1.Pod) {
	for i, container := range pod.Spec.Containers {
		// Check container security context
		if container.SecurityContext == nil {
			report.Issues = append(report.Issues, SecurityIssue{
				Level:       "High",
				Title:       fmt.Sprintf("Container %d: No Security Context", i),
				Description: "Container is running without security context",
				Remediation: "Add securityContext with readOnlyRootFilesystem and allowPrivilegeEscalation: false",
			})
			continue
		}

		sc := container.SecurityContext

		// Check privilege escalation
		if sc.AllowPrivilegeEscalation == nil || *sc.AllowPrivilegeEscalation {
			report.Issues = append(report.Issues, SecurityIssue{
				Level:       "High",
				Title:       fmt.Sprintf("Container %d: Privilege Escalation Allowed", i),
				Description: "Container can escalate privileges",
				Remediation: "Set allowPrivilegeEscalation: false",
			})
		}

		// Check read-only root filesystem
		if sc.ReadOnlyRootFilesystem == nil || !*sc.ReadOnlyRootFilesystem {
			report.Warnings = append(report.Warnings, SecurityWarning{
				Level:       "Medium",
				Title:       fmt.Sprintf("Container %d: Writable Root Filesystem", i),
				Description: "Container has writable root filesystem",
			})
		}

		// Check privileged mode
		if sc.Privileged != nil && *sc.Privileged {
			report.Issues = append(report.Issues, SecurityIssue{
				Level:       "Critical",
				Title:       fmt.Sprintf("Container %d: Privileged Mode", i),
				Description: "Container is running in privileged mode",
				Remediation: "Avoid running containers in privileged mode",
			})
		}

		// Check capabilities
		if sc.Capabilities != nil {
			for _, cap := range sc.Capabilities.Add {
				if isDangerousCapability(string(cap)) {
					report.Issues = append(report.Issues, SecurityIssue{
						Level:       "High",
						Title:       fmt.Sprintf("Container %d: Dangerous Capability %s", i, cap),
						Description: "Container has dangerous capability added",
						Remediation: "Remove unnecessary capabilities",
					})
				}
			}
		}
	}
}

func (s *SecurityAnalyzer) calculateRiskScore(report *SecurityReport) {
	score := 100

	// Deduct points for issues
	for _, issue := range report.Issues {
		switch issue.Level {
		case "Critical":
			score -= 30
		case "High":
			score -= 20
		case "Medium":
			score -= 10
		case "Low":
			score -= 5
		}
	}

	// Deduct points for warnings
	for _, warning := range report.Warnings {
		switch warning.Level {
		case "High":
			score -= 10
		case "Medium":
			score -= 5
		case "Low":
			score -= 2
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	report.Analysis.Score = score

	// Determine risk level
	if score >= 80 {
		report.Analysis.RiskLevel = "Low"
		report.Analysis.Status = "Secure"
	} else if score >= 60 {
		report.Analysis.RiskLevel = "Medium"
		report.Analysis.Status = "Needs Improvement"
	} else if score >= 40 {
		report.Analysis.RiskLevel = "High"
		report.Analysis.Status = "Vulnerable"
	} else {
		report.Analysis.RiskLevel = "Critical"
		report.Analysis.Status = "Highly Vulnerable"
	}

	// Generate recommendations based on score
	if score < 80 {
		report.Recommendations = append(report.Recommendations,
			"Implement security context with runAsNonRoot and readOnlyRootFilesystem")
		report.Recommendations = append(report.Recommendations,
			"Use seccomp profiles and disable privilege escalation")
	}
}

func isDangerousCapability(cap string) bool {
	dangerousCaps := []string{
		"CAP_SYS_ADMIN",
		"CAP_NET_RAW",
		"CAP_SYS_MODULE",
		"CAP_SYS_PTRACE",
	}
	for _, dangerous := range dangerousCaps {
		if cap == dangerous {
			return true
		}
	}
	return false
}
