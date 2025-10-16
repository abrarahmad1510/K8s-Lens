package enterprise

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/enterprise"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

func init() {
	// Add analyze and report subcommands to rbacCmd
	rbacCmd.AddCommand(&cobra.Command{
		Use:   "analyze [namespace]",
		Short: "Analyze RBAC configuration",
		Args:  cobra.RangeArgs(0, 1),
		Run:   analyzeRBAC,
	})

	rbacCmd.AddCommand(&cobra.Command{
		Use:   "report [namespace]",
		Short: "Generate RBAC compliance report",
		Args:  cobra.RangeArgs(0, 1),
		Run:   generateRBACReport,
	})
}

func analyzeRBAC(cmd *cobra.Command, args []string) {
	namespace := "default"
	if len(args) > 0 {
		namespace = args[0]
	}

	utils.PrintInfo("Starting RBAC analysis for namespace: %s", namespace)
	
	k8sClient, err := k8s.NewClient()
	if err != nil {
		utils.PrintError("Error creating Kubernetes client: %v", err)
		os.Exit(1)
	}

	analyzer := enterprise.NewRBACAnalyzer(k8sClient)
	report, err := analyzer.AnalyzeNamespaceRBAC(namespace)
	if err != nil {
		utils.PrintError("Error analyzing RBAC: %v", err)
		os.Exit(1)
	}

	printRBACReport(report)
}

func generateRBACReport(cmd *cobra.Command, args []string) {
	namespace := "default"
	if len(args) > 0 {
		namespace = args[0]
	}
	utils.PrintInfo("Generating RBAC compliance report for namespace: %s", namespace)
	fmt.Printf("RBAC compliance report for %s - Feature coming soon!\n", namespace)
}

func printRBACReport(report *enterprise.RBACReport) {
	fmt.Printf("K8s Lens RBAC Security Analysis Report\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("Namespace: %s\n", report.Namespace)
	fmt.Printf("Risk Level: %s\n", report.RiskLevel)
	
	fmt.Printf("\nRBAC Resources:\n")
	fmt.Printf("  Cluster Roles: %d\n", report.ClusterRoles)
	fmt.Printf("  Roles: %d\n", report.Roles)
	fmt.Printf("  Cluster Role Bindings: %d\n", report.ClusterRoleBindings)
	fmt.Printf("  Role Bindings: %d\n", report.RoleBindings)
	fmt.Printf("  Service Accounts: %d\n", report.ServiceAccounts)

	if len(report.SecurityIssues) > 0 {
		fmt.Printf("\nSecurity Issues Found:\n")
		for i, issue := range report.SecurityIssues {
			fmt.Printf("  %d. [%s] %s: %s\n", i+1, issue.Severity, issue.Type, issue.Description)
		}
	}

	if len(report.Recommendations) > 0 {
		fmt.Printf("\nRecommendations:\n")
		for i, rec := range report.Recommendations {
			fmt.Printf("  %d. %s\n", i+1, rec)
		}
	}
}
