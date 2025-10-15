package analyze

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment [name]",
	Short: "Analyze a Kubernetes Deployment",
	Long:  `Analyze a Kubernetes Deployment and provide diagnostic information.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		verbose, _ := cmd.Flags().GetBool("verbose")

		client, err := k8s.NewClient()
		if err != nil {
			fmt.Printf("Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		analyzer := diagnostics.NewDeploymentAnalyzer(client, namespace)
		report, err := analyzer.Analyze(args[0])
		if err != nil {
			fmt.Printf("Error analyzing deployment: %v\n", err)
			os.Exit(1)
		}

		// Print the report
		fmt.Printf("K8s Lens Analysis Report For Deployment: %s\n", report.Name)
		fmt.Println("---")
		fmt.Printf("Namespace: %s\n", report.Namespace)
		fmt.Printf("Desired Replicas: %d\n", report.DesiredReplicas)
		fmt.Printf("Current Replicas: %d\n", report.CurrentReplicas)
		fmt.Printf("Ready Replicas: %d\n", report.ReadyReplicas)
		fmt.Printf("Available Replicas: %d\n", report.AvailableReplicas)
		fmt.Printf("Updated Replicas: %d\n", report.UpdatedReplicas)
		fmt.Printf("Status: %s\n", report.Analysis.Status)
		fmt.Printf("Rollout Status: %s\n", report.Analysis.RolloutStatus)

		if len(report.Analysis.Issues) > 0 {
			fmt.Println("Issues:")
			for _, issue := range report.Analysis.Issues {
				fmt.Printf("  - %s\n", issue)
			}
		}

		if len(report.Analysis.Recommendations) > 0 {
			fmt.Println("Recommendations:")
			for _, rec := range report.Analysis.Recommendations {
				fmt.Printf("  - %s\n", rec)
			}
		}

		if verbose {
			fmt.Println("Conditions:")
			for _, condition := range report.Conditions {
				fmt.Printf("  - %s: %s (%s)\n", condition.Type, condition.Status, condition.Message)
			}
			fmt.Println("Recent Events:")
			for _, event := range report.Events {
				fmt.Printf("  - [%s] %s: %s\n", event.LastTimestamp.Format("15:04:05"), event.Reason, event.Message)
			}
		}
	},
}

func init() {
	// Add flags
	deploymentCmd.Flags().StringP("namespace", "n", "default", "Namespace")
	deploymentCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
