package enterprise

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RBACAnalyzer analyzes Kubernetes RBAC configurations
type RBACAnalyzer struct {
	client kubernetes.Interface
}

// NewRBACAnalyzer creates a new RBAC analyzer
func NewRBACAnalyzer(client kubernetes.Interface) *RBACAnalyzer {
	return &RBACAnalyzer{
		client: client,
	}
}

// RBACReport contains RBAC analysis results
type RBACReport struct {
	Namespace           string
	ClusterRoles        int
	Roles               int
	ClusterRoleBindings int
	RoleBindings        int
	ServiceAccounts     int
	SecurityIssues      []SecurityIssue
	Recommendations     []string
	RiskLevel           string
}

// SecurityIssue represents a security concern in RBAC configuration
type SecurityIssue struct {
	Type           string
	Severity       string
	Resource       string
	Description    string
	Recommendation string
}

// AnalyzeNamespaceRBAC analyzes RBAC configuration in a namespace
func (r *RBACAnalyzer) AnalyzeNamespaceRBAC(namespace string) (*RBACReport, error) {
	report := &RBACReport{
		Namespace: namespace,
	}

	// Get cluster roles
	clusterRoles, err := r.client.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list cluster roles: %v", err)
	}
	report.ClusterRoles = len(clusterRoles.Items)

	// Get roles in namespace
	roles, err := r.client.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %v", err)
	}
	report.Roles = len(roles.Items)

	// Get cluster role bindings
	clusterRoleBindings, err := r.client.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list cluster role bindings: %v", err)
	}
	report.ClusterRoleBindings = len(clusterRoleBindings.Items)

	// Get role bindings in namespace
	roleBindings, err := r.client.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list role bindings: %v", err)
	}
	report.RoleBindings = len(roleBindings.Items)

	// Get service accounts
	serviceAccounts, err := r.client.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list service accounts: %v", err)
	}
	report.ServiceAccounts = len(serviceAccounts.Items)

	// Analyze security issues
	r.analyzeClusterRoles(report, clusterRoles.Items)
	r.analyzeRoles(report, roles.Items)
	r.analyzeClusterRoleBindings(report, clusterRoleBindings.Items)
	r.analyzeRoleBindings(report, roleBindings.Items)
	r.analyzeServiceAccounts(report, serviceAccounts.Items)

	// Determine overall risk level
	report.RiskLevel = r.calculateRiskLevel(report.SecurityIssues)
	report.Recommendations = r.generateRecommendations(report.SecurityIssues)

	return report, nil
}

func (r *RBACAnalyzer) analyzeClusterRoles(report *RBACReport, clusterRoles []rbacv1.ClusterRole) {
	for _, clusterRole := range clusterRoles {
		for _, rule := range clusterRole.Rules {
			// Check for wildcard permissions
			for _, verb := range rule.Verbs {
				if verb == "*" {
					report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
						Type:           "WildcardVerb",
						Severity:       "High",
						Resource:       clusterRole.Name,
						Description:    fmt.Sprintf("ClusterRole '%s' uses wildcard verb '*'", clusterRole.Name),
						Recommendation: "Replace wildcard verbs with specific actions",
					})
				}
			}

			for _, resource := range rule.Resources {
				if resource == "*" {
					report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
						Type:           "WildcardResource",
						Severity:       "High",
						Resource:       clusterRole.Name,
						Description:    fmt.Sprintf("ClusterRole '%s' uses wildcard resource '*'", clusterRole.Name),
						Recommendation: "Replace wildcard resources with specific resource types",
					})
				}
			}

			// Check for dangerous permissions
			for _, resource := range rule.Resources {
				if resource == "secrets" && containsAny(rule.Verbs, "create", "update", "patch", "delete") {
					report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
						Type:           "DangerousSecretPermission",
						Severity:       "High",
						Resource:       clusterRole.Name,
						Description:    fmt.Sprintf("ClusterRole '%s' has dangerous permissions on secrets", clusterRole.Name),
						Recommendation: "Review and restrict secret permissions",
					})
				}

				if resource == "pods/exec" && contains(rule.Verbs, "create") {
					report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
						Type:           "PodExecPermission",
						Severity:       "Medium",
						Resource:       clusterRole.Name,
						Description:    fmt.Sprintf("ClusterRole '%s' can execute commands in pods", clusterRole.Name),
						Recommendation: "Restrict pod exec permissions to trusted users",
					})
				}
			}
		}
	}
}

func (r *RBACAnalyzer) analyzeRoles(report *RBACReport, roles []rbacv1.Role) {
	for _, role := range roles {
		for _, rule := range role.Rules {
			// Check for wildcard permissions in namespace roles
			for _, verb := range rule.Verbs {
				if verb == "*" {
					report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
						Type:           "WildcardVerbInRole",
						Severity:       "Medium",
						Resource:       fmt.Sprintf("%s/%s", report.Namespace, role.Name),
						Description:    fmt.Sprintf("Role '%s' in namespace '%s' uses wildcard verb '*'", role.Name, report.Namespace),
						Recommendation: "Replace wildcard verbs with specific actions",
					})
				}
			}
		}
	}
}

func (r *RBACAnalyzer) analyzeClusterRoleBindings(report *RBACReport, bindings []rbacv1.ClusterRoleBinding) {
	for _, binding := range bindings {
		// Check for cluster-admin bindings
		if binding.RoleRef.Name == "cluster-admin" {
			report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
				Type:           "ClusterAdminBinding",
				Severity:       "Critical",
				Resource:       binding.Name,
				Description:    fmt.Sprintf("ClusterRoleBinding '%s' grants cluster-admin privileges", binding.Name),
				Recommendation: "Review cluster-admin bindings and use least privilege",
			})
		}

		// Check for bindings to default service accounts
		for _, subject := range binding.Subjects {
			if subject.Kind == "ServiceAccount" && strings.Contains(subject.Name, "default") {
				report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
					Type:           "DefaultServiceAccountBinding",
					Severity:       "High",
					Resource:       binding.Name,
					Description:    fmt.Sprintf("ClusterRoleBinding '%s' binds to default service account", binding.Name),
					Recommendation: "Avoid binding cluster roles to default service accounts",
				})
			}
		}
	}
}

func (r *RBACAnalyzer) analyzeRoleBindings(report *RBACReport, bindings []rbacv1.RoleBinding) {
	for _, binding := range bindings {
		// Check for admin role bindings
		if strings.Contains(binding.RoleRef.Name, "admin") {
			report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
				Type:           "AdminRoleBinding",
				Severity:       "Medium",
				Resource:       fmt.Sprintf("%s/%s", report.Namespace, binding.Name),
				Description:    fmt.Sprintf("RoleBinding '%s' in namespace '%s' uses admin role", binding.Name, report.Namespace),
				Recommendation: "Review admin role usage and apply least privilege",
			})
		}
	}
}

func (r *RBACAnalyzer) analyzeServiceAccounts(report *RBACReport, serviceAccounts []corev1.ServiceAccount) {
	// Check for service accounts without explicit secrets
	for _, sa := range serviceAccounts {
		if len(sa.Secrets) == 0 {
			report.SecurityIssues = append(report.SecurityIssues, SecurityIssue{
				Type:           "ServiceAccountWithoutSecrets",
				Severity:       "Low",
				Resource:       sa.Name,
				Description:    fmt.Sprintf("ServiceAccount '%s' has no explicitly defined secrets", sa.Name),
				Recommendation: "Consider defining explicit secrets for service accounts",
			})
		}
	}
}

func (r *RBACAnalyzer) calculateRiskLevel(issues []SecurityIssue) string {
	criticalCount := 0
	highCount := 0
	mediumCount := 0

	for _, issue := range issues {
		switch issue.Severity {
		case "Critical":
			criticalCount++
		case "High":
			highCount++
		case "Medium":
			mediumCount++
		}
	}

	if criticalCount > 0 {
		return "Critical"
	} else if highCount > 0 {
		return "High"
	} else if mediumCount > 0 {
		return "Medium"
	} else {
		return "Low"
	}
}

func (r *RBACAnalyzer) generateRecommendations(issues []SecurityIssue) []string {
	recommendations := []string{}

	hasWildcard := false
	hasClusterAdmin := false
	hasDangerousPermissions := false

	for _, issue := range issues {
		switch issue.Type {
		case "WildcardVerb", "WildcardResource", "WildcardVerbInRole":
			hasWildcard = true
		case "ClusterAdminBinding":
			hasClusterAdmin = true
		case "DangerousSecretPermission", "PodExecPermission":
			hasDangerousPermissions = true
		}
	}

	if hasWildcard {
		recommendations = append(recommendations,
			"Replace all wildcard permissions with specific verbs and resources")
	}

	if hasClusterAdmin {
		recommendations = append(recommendations,
			"Review and minimize cluster-admin bindings - use least privilege principles")
	}

	if hasDangerousPermissions {
		recommendations = append(recommendations,
			"Restrict dangerous permissions (secrets, pod exec) to trusted principals only")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"RBAC configuration appears secure - maintain current security practices")
	}

	return recommendations
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsAny(slice []string, items ...string) bool {
	for _, item := range items {
		if contains(slice, item) {
			return true
		}
	}
	return false
}
