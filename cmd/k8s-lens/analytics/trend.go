package analytics

import (
	"fmt"
	"os"
	"time"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/analytics"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var trendCmd = &cobra.Command{
	Use:   "trend [namespace]",
	Short: "Analyze resource and performance trends",
	Long:  "Analyze historical trends and patterns in resource usage and performance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]
		periodStr, _ := cmd.Flags().GetString("period")

		// Parse period
		period, err := time.ParseDuration(periodStr)
		if err != nil {
			utils.PrintError("Invalid period format: %v", err)
			os.Exit(1)
		}

		utils.PrintInfo("Starting trend analysis for namespace: %s (period: %v)", namespace, period)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := analytics.NewTrendAnalyzer(k8sClient)
		report, err := analyzer.AnalyzeNamespaceTrends(namespace, period)
		if err != nil {
			utils.PrintError("Error analyzing trends: %v", err)
			os.Exit(1)
		}

		// Print report
		fmt.Printf("K8s Lens Trend Analysis Report\n")
		fmt.Printf("==============================\n")
		fmt.Printf("Namespace: %s\n", report.Namespace)
		fmt.Printf("Analysis Period: %v\n", report.AnalysisPeriod)
		fmt.Printf("Generated: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))

		utils.PrintSection("Resource Trends")
		if len(report.ResourceTrends) == 0 {
			fmt.Println("No resource trend data available")
		} else {
			for _, trend := range report.ResourceTrends {
				trendIndicator := "→"
				if trend.Trend == "Increasing" {
					trendIndicator = "↑"
				} else if trend.Trend == "Decreasing" {
					trendIndicator = "↓"
				}

				fmt.Printf("%s %s: %.1f (from %.1f) %s %.1f%%\n",
					trendIndicator, trend.Metric, trend.CurrentValue,
					trend.PreviousValue, trend.Trend, trend.ChangePercent)
			}
		}

		utils.PrintSection("Performance Trends")
		if len(report.PerformanceTrends) == 0 {
			fmt.Println("No performance trends detected")
		} else {
			for i, trend := range report.PerformanceTrends {
				fmt.Printf("%d. [%s] %s\n", i+1, trend.Impact, trend.Pattern)
				fmt.Printf("   Component: %s, Metric: %s\n", trend.Component, trend.Metric)
				fmt.Printf("   Confidence: %.0f%%\n", trend.Confidence*100)
			}
		}

		utils.PrintSection("Recommendations")
		for i, recommendation := range report.Recommendations {
			fmt.Printf("%d. %s\n", i+1, recommendation)
		}
	},
}

func init() {
	trendCmd.Flags().StringP("period", "p", "24h", "Analysis period (e.g., 24h, 7d, 30d)")
}
