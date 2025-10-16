package enterprise

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SecurityScanner provides comprehensive security scanning
type SecurityScanner struct {
	client kubernetes.Interface
}

// NewSecurityScanner creates a new security scanner
func NewSecurityScanner(client kubernetes.Interface) *SecurityScanner {
	return &SecurityScanner{
		client: client,
	}
}

// SecurityScanReport contains security scan results
type SecurityScanReport struct {
	Namespace       string
	TotalPods       int
	TotalServices   int
	SecurityIssues  []SecurityIssue
	ComplianceScore int
	RiskLevel       string
	Recommendations []string
}

// ScanNamespace performs a comprehensive security scan of a namespace
func (s *SecurityScanner) ScanNamespace(namespace string) (*SecurityScanReport, error) {
	report := &SecurityScanReport{
		Namespace: namespace,
	}

	// Get all pods in the namespace
	pods, err := s.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}
	report.TotalPods = len(pods.Items)

	// Get all services in the namespace
	services, err := s.client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}
	report.TotalServices = len(services.Items)

	// Perform security scans
	s.scanPodSecurity(report, pods.Items)
	s.scanServiceSecurity(report, services.Items)
	s.scanNetworkPolicies(report, namespace)

	// Calculate compliance score and risk level
	report.ComplianceScore = s.calculateComplianceScore(report.SecurityIssues)
	report.RiskLevel = s.calculateRiskLevelFromScore(report.ComplianceScore)
	report.Recommendations = s.generateSecurityRecommendations(report.SecurityIssues)

	return report, nil
}

func (s *SecurityScanner) scanPodSecurity(report *SecurityScanReport, pods []corev1.Pod) {
	for _, pod := range pods {
		// Check pod security context
		if pod.Spec.SecurityContext == nil {
			report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
				Type:           "MissingPodSecurityContext",
				Severity:       "Medium",
				Resource:       pod.Name,
				Description:    "Pod does not have a security context defined",
				Recommendation: "Define pod-level security context with reasonable defaults",
			})
		} else {
			// Check specific security context settings
			if pod.Spec.SecurityContext.RunAsNonRoot == nil || !*pod.Spec.SecurityContext.RunAsNonRoot {
				report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
					Type:           "RunAsRootAllowed",
					Severity:       "Medium",
					Resource:       pod.Name,
					Description:    "Pod can run as root user",
					Recommendation: "Set runAsNonRoot to true in security context",
				})
			}

			if pod.Spec.SecurityContext.SeccompProfile == nil {
				report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
					Type:           "MissingSeccompProfile",
					Severity:       "Low",
					Resource:       pod.Name,
					Description:    "Pod does not have seccomp profile defined",
					Recommendation: "Define seccomp profile for enhanced security",
				})
			}
		}

		// Check container security
		for _, container := range pod.Spec.Containers {
			if container.SecurityContext == nil {
				report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
					Type:           "MissingContainerSecurityContext",
					Severity:       "Medium",
					Resource:       fmt.Sprintf("%s/%s", pod.Name, container.Name),
					Description:    "Container does not have security context defined",
					Recommendation: "Define container security context with least privilege",
				})
				continue
			}

			// Check for privileged mode
			if container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
				report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
					Type:           "PrivilegedContainer",
					Severity:       "High",
					Resource:       fmt.Sprintf("%s/%s", pod.Name, container.Name),
					Description:    "Container is running in privileged mode",
					Recommendation: "Avoid running containers in privileged mode",
				})
			}

			// Check for root user
			if container.SecurityContext.RunAsUser != nil && *container.SecurityContext.RunAsUser == 0 {
				report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
					Type:           "RunAsRoot",
					Severity:       "Medium",
					Resource:       fmt.Sprintf("%s/%s", pod.Name, container.Name),
					Description:    "Container is running as root user",
					Recommendation: "Run containers as non-root user",
				})
			}

			// Check for read-only root filesystem
			if container.SecurityContext.ReadOnlyRootFilesystem == nil || !*container.SecurityContext.ReadOnlyRootFilesystem {
				report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
					Type:           "WritableRootFilesystem",
					Severity:       "Low",
					Resource:       fmt.Sprintf("%s/%s", pod.Name, container.Name),
					Description:    "Container has writable root filesystem",
					Recommendation: "Set readOnlyRootFilesystem to true",
				})
			}

			// Check for dangerous capabilities
			if container.SecurityContext.Capabilities != nil {
				for _, cap := range container.SecurityContext.Capabilities.Add {
					if isDangerousCapability(string(cap)) {
						report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
							Type:           "DangerousCapability",
							Severity:       "High",
							Resource:       fmt.Sprintf("%s/%s", pod.Name, container.Name),
							Description:    fmt.Sprintf("Container has dangerous capability: %s", cap),
							Recommendation: "Remove dangerous capabilities",
						})
					}
				}
			}
		}
	}
}

func (s *SecurityScanner) scanServiceSecurity(report *SecurityScanReport, services []corev1.Service) {
	for _, service := range services {
		// Check for services with external IPs
		if len(service.Spec.ExternalIPs) > 0 {
			report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
				Type:           "ServiceWithExternalIP",
				Severity:       "Medium",
				Resource:       service.Name,
				Description:    "Service has external IPs configured",
				Recommendation: "Review external IP usage for security implications",
			})
		}

		// Check for LoadBalancer services
		if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
			report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
				Type:           "LoadBalancerService",
				Severity:       "Low",
				Resource:       service.Name,
				Description:    "Service uses LoadBalancer type",
				Recommendation: "Consider using Ingress instead of LoadBalancer for external access",
			})
		}
	}
}

func (s *SecurityScanner) scanNetworkPolicies(report *SecurityScanReport, namespace string) {
	// Check if namespace has network policies
	networkPolicies, err := s.client.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// Network policies might not be available in all clusters
		return
	}

	if len(networkPolicies.Items) == 0 {
		report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
			Type:           "NoNetworkPolicies",
			Severity:       "Medium",
			Resource:       namespace,
			Description:    "Namespace has no network policies defined",
			Recommendation: "Implement network policies for network segmentation",
		})
	}
}

func (s *SecurityScanner) calculateComplianceScore(issues []SecurityIssue) int {
	baseScore := 100

	for _, issue := range issues {
		switch issue.Severity {
		case "Critical":
			baseScore -= 10
		case "High":
			baseScore -= 7
		case "Medium":
			baseScore -= 4
		case "Low":
			baseScore -= 1
		}
	}

	if baseScore < 0 {
		return 0
	}
	return baseScore
}

func (s *SecurityScanner) calculateRiskLevelFromScore(score int) string {
	if score >= 90 {
		return "Low"
	} else if score >= 70 {
		return "Medium"
	} else if score >= 50 {
		return "High"
	} else {
		return "Critical"
	}
}

func (s *SecurityScanner) generateSecurityRecommendations(issues []SecurityIssue) []string {
	recommendations := []string{}

	hasPrivilegedContainers := false
	hasRootContainers := false
	hasMissingSecurityContexts := false

	for _, issue := range issues {
		switch issue.Type {
		case "PrivilegedContainer":
			hasPrivilegedContainers = true
		case "RunAsRoot", "RunAsRootAllowed":
			hasRootContainers = true
		case "MissingPodSecurityContext", "MissingContainerSecurityContext":
			hasMissingSecurityContexts = true
		}
	}

	if hasPrivilegedContainers {
		recommendations = append(recommendations,
			"Eliminate all privileged containers - they pose significant security risks")
	}

	if hasRootContainers {
		recommendations = append(recommendations,
			"Run containers as non-root users to minimize attack surface")
	}

	if hasMissingSecurityContexts {
		recommendations = append(recommendations,
			"Define security contexts for all pods and containers")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"Security posture is good - maintain current security practices")
	}

	return recommendations
}

func isDangerousCapability(cap string) bool {
	dangerousCaps := []string{
		"SYS_ADMIN", "SYS_PTRACE", "SYS_MODULE", "SYS_RAWIO",
		"NET_ADMIN", "NET_RAW", "IPC_LOCK", "DAC_READ_SEARCH",
	}

	for _, dangerous := range dangerousCaps {
		if cap == dangerous {
			return true
		}
	}
	return false
}
