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
	// Add scan and audit subcommands to securityCmd
	securityCmd.AddCommand(&cobra.Command{
		Use:   "scan [namespace]",
		Short: "Scan for security vulnerabilities",
		Args:  cobra.RangeArgs(0, 1),
		Run:   scanSecurity,
	})

	securityCmd.AddCommand(&cobra.Command{
		Use:   "audit [namespace]",
		Short: "Run comprehensive security audit",
		Args:  cobra.RangeArgs(0, 1),
		Run:   runSecurityAudit,
	})
}

func scanSecurity(cmd *cobra.Command, args []string) {
	namespace := "default"
	if len(args) > 0 {
		namespace = args[0]
	}

	utils.PrintInfo("Starting security scan for namespace: %s", namespace)
	
	k8sClient, err := k8s.NewClient()
	if err != nil {
		utils.PrintError("Error creating Kubernetes client: %v", err)
		os.Exit(1)
	}

	scanner := enterprise.NewSecurityScanner(k8sClient)
	report, err := scanner.ScanNamespace(namespace)
	if err != nil {
		utils.PrintError("Error scanning security: %v", err)
		os.Exit(1)
	}

	printSecurityReport(report)
}

func runSecurityAudit(cmd *cobra.Command, args []string) {
	namespace := "default"
	if len(args) > 0 {
		namespace = args[0]
	}
	utils.PrintInfo("Running comprehensive security audit for namespace: %s", namespace)
	fmt.Printf("Comprehensive security audit for %s - Feature coming soon!\n", namespace)
}

func printSecurityReport(report *enterprise.SecurityScanReport) {
	fmt.Printf("K8s Lens Security Scan Report\n")
	fmt.Printf("=============================\n")
	fmt.Printf("Namespace: %s\n", report.Namespace)
	fmt.Printf("Compliance Score: %d/100\n", report.ComplianceScore)
	fmt.Printf("Risk Level: %s\n", report.RiskLevel)
	
	fmt.Printf("\nResources Scanned:\n")
	fmt.Printf("  Pods: %d\n", report.TotalPods)
	fmt.Printf("  Services: %d\n", report.TotalServices)

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
