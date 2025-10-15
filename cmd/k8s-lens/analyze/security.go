package analyze

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var securityCmd = &cobra.Command{
	Use:   "security [pod-name]",
	Short: "Analyze Pod Security",
	Long:  `Perform security analysis of Kubernetes pods and containers.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")

		utils.PrintInfo("Performing security analysis for pod: %s in namespace: %s", args[0], namespace)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := diagnostics.NewSecurityAnalyzer(k8sClient, namespace)
		report, err := analyzer.AnalyzePodSecurity(args[0])
		if err != nil {
			utils.PrintError("Error analyzing pod security: %v", err)
			os.Exit(1)
		}

		fmt.Printf("K8s Lens Security Analysis: %s\n", report.PodName)
		fmt.Println("---")

		utils.PrintSection("Security Assessment")
		fmt.Printf("Namespace: %s\n", report.Namespace)
		fmt.Printf("Security Status: %s\n", report.Analysis.Status)
		fmt.Printf("Risk Level: %s\n", report.Analysis.RiskLevel)
		fmt.Printf("Security Score: %d/100\n", report.Analysis.Score)

		if len(report.Issues) > 0 {
			utils.PrintSection("Security Issues")
			for _, issue := range report.Issues {
				color := "red"
				switch issue.Level {
				case "Critical":
					color = "red"
				case "High":
					color = "red"
				case "Medium":
					color = "yellow"
				case "Low":
					color = "blue"
				}
				fmt.Printf("- [%s] %s\n", utils.Colorize(issue.Level, color), issue.Title)
				fmt.Printf("  Description: %s\n", issue.Description)
				fmt.Printf("  Remediation: %s\n", issue.Remediation)
				fmt.Println()
			}
		}

		if len(report.Warnings) > 0 {
			utils.PrintSection("Security Warnings")
			for _, warning := range report.Warnings {
				color := "yellow"
				switch warning.Level {
				case "High":
					color = "red"
				case "Medium":
					color = "yellow"
				case "Low":
					color = "blue"
				}
				fmt.Printf("- [%s] %s\n", utils.Colorize(warning.Level, color), warning.Title)
				fmt.Printf("  Description: %s\n", warning.Description)
				fmt.Println()
			}
		}

		if len(report.Recommendations) > 0 {
			utils.PrintSection("Security Recommendations")
			for _, rec := range report.Recommendations {
				utils.PrintInfo("- %s", rec)
			}
		}

		// Security score interpretation
		utils.PrintSection("Score Interpretation")
		if report.Analysis.Score >= 80 {
			utils.PrintSuccess("Excellent security posture")
		} else if report.Analysis.Score >= 60 {
			utils.PrintWarning("Moderate security posture - improvements needed")
		} else {
			utils.PrintError("Poor security posture - immediate action required")
		}
	},
}

func init() {
	securityCmd.Flags().StringP("namespace", "n", "default", "Namespace")
}
