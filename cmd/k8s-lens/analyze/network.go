package analyze

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network [policy-name]",
	Short: "Analyze Kubernetes Network Policies",
	Long:  `Analyze Kubernetes Network Policies and provide security assessment.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := diagnostics.NewNetworkAnalyzer(k8sClient, namespace)

		if len(args) == 1 {
			// Analyze specific network policy
			utils.PrintInfo("Analyzing network policy: %s in namespace: %s", args[0], namespace)
			report, err := analyzer.AnalyzeNetworkPolicy(args[0])
			if err != nil {
				utils.PrintError("Error analyzing network policy: %v", err)
				os.Exit(1)
			}

			fmt.Printf("K8s Lens Network Policy Analysis: %s\n", report.Name)
			fmt.Println("---")

			utils.PrintSection("Policy Configuration")
			fmt.Printf("Namespace: %s\n", report.Namespace)
			fmt.Printf("Policy Type: %s\n", report.PolicyType)

			utils.PrintSection("Pod Selector")
			if len(report.PodSelector) > 0 {
				for key, value := range report.PodSelector {
					fmt.Printf("- %s: %s\n", key, value)
				}
			} else {
				fmt.Println("  Applies to all pods in namespace")
			}

			utils.PrintSection("Security Assessment")
			fmt.Printf("Status: %s\n", report.Analysis.Status)

			if len(report.Analysis.Issues) > 0 {
				utils.PrintSection("Security Issues")
				for _, issue := range report.Analysis.Issues {
					utils.PrintWarning("- %s", issue)
				}
			}

			if len(report.Analysis.Recommendations) > 0 {
				utils.PrintSection("Recommendations")
				for _, rec := range report.Analysis.Recommendations {
					utils.PrintInfo("- %s", rec)
				}
			}

		} else {
			// Analyze all network policies in namespace
			utils.PrintInfo("Analyzing all network policies in namespace: %s", namespace)
			report, err := analyzer.AnalyzeNamespaceNetworkPolicies()
			if err != nil {
				utils.PrintError("Error analyzing network policies: %v", err)
				os.Exit(1)
			}

			fmt.Printf("K8s Lens Network Policy Analysis - Namespace: %s\n", report.Namespace)
			fmt.Println("---")

			utils.PrintSection("Namespace Overview")
			fmt.Printf("Total Policies: %d\n", report.TotalPolicies)
			fmt.Printf("Coverage Status: %s\n", report.CoverageStatus)

			if report.TotalPolicies > 0 {
				utils.PrintSection("Policy Details")
				for _, policyReport := range report.PolicyReports {
					fmt.Printf("\nPolicy: %s\n", policyReport.Name)
					fmt.Printf("  Type: %s\n", policyReport.PolicyType)
					fmt.Printf("  Status: %s\n", policyReport.Analysis.Status)
					if len(policyReport.Analysis.Issues) > 0 {
						fmt.Printf("  Issues: %d\n", len(policyReport.Analysis.Issues))
					}
				}
			}

			if len(report.Recommendations) > 0 {
				utils.PrintSection("Recommendations")
				for _, rec := range report.Recommendations {
					utils.PrintInfo("- %s", rec)
				}
			}
		}
	},
}

func init() {
	networkCmd.Flags().StringP("namespace", "n", "default", "Namespace")
}
