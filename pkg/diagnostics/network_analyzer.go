package diagnostics

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NetworkAnalyzer provides analysis for NetworkPolicy resources
type NetworkAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewNetworkAnalyzer creates a new NetworkAnalyzer
func NewNetworkAnalyzer(client kubernetes.Interface, namespace string) *NetworkAnalyzer {
	return &NetworkAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// NetworkPolicyReport contains the analysis report
type NetworkPolicyReport struct {
	Name        string
	Namespace   string
	PolicyType  string
	PodSelector map[string]string
	Ingress     []networkingv1.NetworkPolicyIngressRule
	Egress      []networkingv1.NetworkPolicyEgressRule
	Analysis    NetworkPolicyAnalysis
}

// NetworkPolicyAnalysis contains diagnostic results
type NetworkPolicyAnalysis struct {
	Status          string
	Issues          []string
	Recommendations []string
	Coverage        string
}

// AnalyzeNetworkPolicy analyzes a specific NetworkPolicy
func (n *NetworkAnalyzer) AnalyzeNetworkPolicy(policyName string) (*NetworkPolicyReport, error) {
	policy, err := n.client.NetworkingV1().NetworkPolicies(n.namespace).Get(context.TODO(), policyName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get network policy %s: %v", policyName, err)
	}

	report := &NetworkPolicyReport{
		Name:        policy.Name,
		Namespace:   policy.Namespace,
		PodSelector: policy.Spec.PodSelector.MatchLabels,
		Ingress:     policy.Spec.Ingress,
		Egress:      policy.Spec.Egress,
	}

	n.analyzePolicy(report, policy)

	return report, nil
}

// AnalyzeNamespaceNetworkPolicies analyzes all NetworkPolicies in a namespace
func (n *NetworkAnalyzer) AnalyzeNamespaceNetworkPolicies() (*NamespaceNetworkReport, error) {
	policies, err := n.client.NetworkingV1().NetworkPolicies(n.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list network policies: %v", err)
	}

	report := &NamespaceNetworkReport{
		Namespace:     n.namespace,
		TotalPolicies: len(policies.Items),
		PolicyReports: []NetworkPolicyReport{},
	}

	for _, policy := range policies.Items {
		policyReport := NetworkPolicyReport{
			Name:        policy.Name,
			Namespace:   policy.Namespace,
			PodSelector: policy.Spec.PodSelector.MatchLabels,
			Ingress:     policy.Spec.Ingress,
			Egress:      policy.Spec.Egress,
		}
		n.analyzePolicy(&policyReport, &policy)
		report.PolicyReports = append(report.PolicyReports, policyReport)
	}

	n.analyzeNamespaceCoverage(report)

	return report, nil
}

func (n *NetworkAnalyzer) analyzePolicy(report *NetworkPolicyReport, policy *networkingv1.NetworkPolicy) {
	// Determine policy type
	if len(policy.Spec.Ingress) > 0 && len(policy.Spec.Egress) > 0 {
		report.PolicyType = "Ingress and Egress"
	} else if len(policy.Spec.Ingress) > 0 {
		report.PolicyType = "Ingress only"
	} else if len(policy.Spec.Egress) > 0 {
		report.PolicyType = "Egress only"
	} else {
		report.PolicyType = "No rules (default deny)"
	}

	// Check for overly permissive rules
	for _, ingress := range policy.Spec.Ingress {
		if len(ingress.From) == 0 {
			report.Analysis.Issues = append(report.Analysis.Issues,
				"Ingress rule allows traffic from all sources")
		}
	}

	for _, egress := range policy.Spec.Egress {
		if len(egress.To) == 0 {
			report.Analysis.Issues = append(report.Analysis.Issues,
				"Egress rule allows traffic to all destinations")
		}
	}

	// Check pod selector
	if len(policy.Spec.PodSelector.MatchLabels) == 0 {
		report.Analysis.Issues = append(report.Analysis.Issues,
			"Policy applies to all pods in namespace (no pod selector)")
	}

	if len(report.Analysis.Issues) == 0 {
		report.Analysis.Status = "Secure"
	} else {
		report.Analysis.Status = "Needs Review"
	}
}

func (n *NetworkAnalyzer) analyzeNamespaceCoverage(report *NamespaceNetworkReport) {
	if report.TotalPolicies == 0 {
		report.CoverageStatus = "No network policies"
		report.Recommendations = append(report.Recommendations,
			"Consider implementing network policies for namespace isolation")
	} else {
		report.CoverageStatus = fmt.Sprintf("%d policies active", report.TotalPolicies)
	}

	// Count policies with issues
	issuesCount := 0
	for _, policyReport := range report.PolicyReports {
		if len(policyReport.Analysis.Issues) > 0 {
			issuesCount++
		}
	}

	if issuesCount > 0 {
		report.Recommendations = append(report.Recommendations,
			fmt.Sprintf("%d policies have configuration issues that need review", issuesCount))
	}
}

// NamespaceNetworkReport contains analysis of all network policies in a namespace
type NamespaceNetworkReport struct {
	Namespace       string
	TotalPolicies   int
	PolicyReports   []NetworkPolicyReport
	CoverageStatus  string
	Recommendations []string
}
