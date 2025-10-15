package optimize

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/abrarahmad1510/k8s-lens/pkg/optimization"
	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource [namespace]",
	Short: "Optimize resource allocation and reduce costs",
	Long:  `Analyze and optimize Kubernetes resource allocation for cost savings.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]

		utils.PrintInfo("Starting resource optimization analysis for namespace: %s", namespace)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		optimizer := optimization.NewResourceOptimizer(k8sClient)
		report, err := optimizer.AnalyzeNamespace(namespace)
		if err != nil {
			utils.PrintError("Error analyzing resource optimization: %v", err)
			os.Exit(1)
		}

		fmt.Printf("K8s Lens Resource Optimization Report: %s\n", namespace)
		fmt.Println("===")

		utils.PrintSection("Namespace Overview")
		fmt.Printf("Total Pods: %d\n", report.TotalPods)
		fmt.Printf("Analyzed Pods: %d\n", report.AnalyzedPods)
		fmt.Printf("Total Optimizations: %d\n", report.Summary.TotalOptimizations)
		fmt.Printf("Estimated Monthly Savings: $%.2f\n", report.Summary.TotalMonthlySavings)
		fmt.Printf("Overall Confidence: %d%%\n", report.Summary.OverallConfidence)
		fmt.Printf("Risk Level: %s\n", report.Summary.RiskLevel)

		if len(report.Optimizations) > 0 {
			utils.PrintSection("Optimization Recommendations")
			for i, opt := range report.Optimizations {
				fmt.Printf("\nOptimization %d:\n", i+1)
				fmt.Printf("  Pod: %s | Container: %s\n", opt.PodName, opt.ContainerName)
				fmt.Printf("  Type: %s\n", opt.Type)
				fmt.Printf("  Current: CPU=%s, Memory=%s\n", opt.Current.CPU, opt.Current.Memory)
				fmt.Printf("  Recommended: CPU=%s, Memory=%s\n", opt.Recommended.CPU, opt.Recommended.Memory)
				fmt.Printf("  Monthly Savings: $%.2f (%.1f%%)\n", opt.Savings.MonthlySavings, opt.Savings.PercentSavings)
				fmt.Printf("  Confidence: %d%%\n", opt.Confidence)
				fmt.Printf("  Description: %s\n", opt.Description)
			}
		} else {
			utils.PrintSuccess("No optimization opportunities found. Resources are well configured!")
		}

		utils.PrintSection("Next Steps")
		if report.Summary.TotalMonthlySavings > 0 {
			utils.PrintInfo("Apply optimizations to save approximately $%.2f per month", report.Summary.TotalMonthlySavings)
			utils.PrintInfo("Use 'k8s-lens optimize fix' to generate automated patches")
		} else {
			utils.PrintSuccess("Your resource configuration is optimal!")
		}
	},
}
