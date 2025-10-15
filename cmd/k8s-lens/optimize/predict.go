package optimize

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/ai"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var predictCmd = &cobra.Command{
	Use:   "predict [deployment-name]",
	Short: "Predict potential failures and issues",
	Long:  `Use AI-powered analysis to predict potential failures and performance issues.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		deploymentName := args[0]

		utils.PrintInfo("Starting predictive analysis for deployment: %s in namespace: %s", deploymentName, namespace)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := ai.NewPredictiveAnalyzer(k8sClient)
		report, err := analyzer.PredictFailures(deploymentName, namespace)
		if err != nil {
			utils.PrintError("Error performing predictive analysis: %v", err)
			os.Exit(1)
		}

		fmt.Printf("K8s Lens Predictive Analysis: %s\n", deploymentName)
		fmt.Println("===")

		utils.PrintSection("Risk Assessment")
		fmt.Printf("Overall Risk: %s\n", report.OverallRisk)
		fmt.Printf("Confidence: %d%%\n", report.Confidence)

		if len(report.Predictions) > 0 {
			utils.PrintSection("Predictions")
			for i, prediction := range report.Predictions {
				color := "yellow"
				if prediction.Probability >= 70 {
					color = "red"
				} else if prediction.Probability >= 50 {
					color = "yellow"
				} else {
					color = "blue"
				}

				fmt.Printf("\nPrediction %d:\n", i+1)
				fmt.Printf("  Type: %s\n", prediction.Type)
				fmt.Printf("  Probability: %s\n", utils.Colorize(fmt.Sprintf("%d%%", prediction.Probability), color))
				fmt.Printf("  Timeframe: %s\n", prediction.Timeframe)
				fmt.Printf("  Description: %s\n", prediction.Description)
				fmt.Printf("  Evidence:\n")
				for _, evidence := range prediction.Evidence {
					fmt.Printf("    - %s\n", evidence)
				}
			}
		} else {
			utils.PrintSuccess("No significant risks predicted!")
		}

		utils.PrintSection("Recommendations")
		for i, rec := range report.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec)
		}

		// Risk interpretation
		utils.PrintSection("Risk Interpretation")
		switch report.OverallRisk {
		case "High":
			utils.PrintError("HIGH RISK: Immediate action recommended")
		case "Medium":
			utils.PrintWarning("MEDIUM RISK: Proactive measures recommended")
		case "Low":
			utils.PrintSuccess("LOW RISK: Continue monitoring")
		}
	},
}

func init() {
	predictCmd.Flags().StringP("namespace", "n", "default", "Namespace")
}
