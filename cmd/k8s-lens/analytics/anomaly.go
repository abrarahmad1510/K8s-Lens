package analytics

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/abrarahmad1510/k8s-lens/pkg/machinelearning"
	"github.com/spf13/cobra"
)

var anomalyCmd = &cobra.Command{
	Use:   "anomaly [namespace]",
	Short: "Detect anomalies in Kubernetes namespace",
	Long:  "Use machine learning to detect unusual patterns and potential issues in your Kubernetes namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]
		verbose, _ := cmd.Flags().GetBool("verbose")

		utils.PrintInfo("Starting anomaly detection for namespace: %s", namespace)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		detector := machinelearning.NewAnomalyDetector(k8sClient)
		report, err := detector.DetectNamespaceAnomalies(namespace)
		if err != nil {
			utils.PrintError("Error detecting anomalies: %v", err)
			os.Exit(1)
		}

		// Print report
		fmt.Printf("K8s Lens Anomaly Detection Report\n")
		fmt.Printf("=================================\n")
		fmt.Printf("Namespace: %s\n", report.Namespace)
		fmt.Printf("Total Pods: %d\n", report.TotalPods)
		fmt.Printf("Anomaly Score: %d/100\n", report.Score)
		fmt.Printf("Generated: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println()

		utils.PrintSection("Detected Anomalies")
		if len(report.Anomalies) == 0 {
			utils.PrintSuccess("No anomalies detected!")
		} else {
			for i, anomaly := range report.Anomalies {
				fmt.Printf("%d. [%s] %s - %s\n", i+1, anomaly.Severity, anomaly.Resource, anomaly.Message)
				fmt.Printf("   Type: %s, Confidence: %.0f%%\n", anomaly.Type, anomaly.Confidence*100)
				if verbose {
					fmt.Printf("   Detected: %s\n", anomaly.Timestamp.Format("15:04:05"))
				}
				fmt.Println()
			}
		}

		utils.PrintSection("Recommendations")
		for i, recommendation := range report.Recommendations {
			fmt.Printf("%d. %s\n", i+1, recommendation)
		}
	},
}

func init() {
	anomalyCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
