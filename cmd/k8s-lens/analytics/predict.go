package analytics

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/abrarahmad1510/k8s-lens/pkg/machinelearning"
	"github.com/spf13/cobra"
)

var predictCmd = &cobra.Command{
	Use:   "predict [deployment-name]",
	Short: "Predict potential future issues for deployment",
	Long:  "Use predictive analytics to forecast potential issues and capacity problems",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentName := args[0]
		namespace, _ := cmd.Flags().GetString("namespace")

		utils.PrintInfo("Starting predictive analysis for deployment: %s", deploymentName)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		predictor := machinelearning.NewPredictiveAnalyzer(k8sClient)
		report, err := predictor.PredictDeploymentFailures(deploymentName, namespace)
		if err != nil {
			utils.PrintError("Error generating predictions: %v", err)
			os.Exit(1)
		}

		// Print report
		fmt.Printf("K8s Lens Predictive Analysis Report\n")
		fmt.Printf("===================================\n")
		fmt.Printf("Deployment: %s\n", deploymentName)
		fmt.Printf("Namespace: %s\n", report.Namespace)
		fmt.Printf("Confidence: %.0f%%\n", report.Confidence*100)
		fmt.Printf("Time Horizon: %v\n", report.TimeHorizon)
		fmt.Printf("Generated: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))

		utils.PrintSection("Predictions")
		if len(report.Predictions) == 0 {
			utils.PrintSuccess("No potential issues predicted in the next %v!", report.TimeHorizon)
		} else {
			for i, prediction := range report.Predictions {
				fmt.Printf("%d. [%s] %s\n", i+1, prediction.Impact, prediction.Message)
				fmt.Printf("   Probability: %.0f%%, Expected: %s\n",
					prediction.Probability*100,
					prediction.ExpectedTime.Format("Jan 02, 15:04"))
				fmt.Printf("   Recommendation: %s\n", prediction.Recommendation)
				fmt.Println()
			}
		}
	},
}

func init() {
	predictCmd.Flags().StringP("namespace", "n", "default", "Namespace")
}
