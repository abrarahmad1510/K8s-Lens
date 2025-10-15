package analyze

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var podCmd = &cobra.Command{
	Use:   "pod [name]",
	Short: "Analyze a Kubernetes Pod",
	Long:  `Analyze a Kubernetes Pod and provide diagnostic information.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		verbose, _ := cmd.Flags().GetBool("verbose")

		utils.PrintInfo("Starting pod analysis for: %s in namespace: %s", args[0], namespace)

		client, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := diagnostics.NewPodAnalyzer(client, namespace)
		report, err := analyzer.Analyze(args[0])
		if err != nil {
			utils.PrintError("Error analyzing pod: %v", err)
			os.Exit(1)
		}

		// Print the report
		fmt.Printf("K8s Lens Analysis Report For Pod: %s\n", report.Name)
		fmt.Println("---")

		utils.PrintSection("Pod Status Analysis")
		fmt.Printf("Phase: %s\n", report.Phase)
		fmt.Printf("Node: %s\n", report.Node)
		fmt.Printf("Created: %s\n", report.Created.Format("Mon, 02 Jan 2006 15:04:05 UTC"))

		if report.Status == "Running" {
			utils.PrintSuccess("Status: Pod Is Running Normally")
		} else {
			utils.PrintWarning("Status: %s", report.Status)
		}

		utils.PrintSection("Container Status Analysis")
		for _, container := range report.Containers {
			fmt.Printf("Container: %s\n", container.Name)
			fmt.Printf("Image: %s\n", container.Image)
			fmt.Printf("Status: %s\n", container.Status)

			if container.Ready {
				utils.PrintSuccess("Status: Container Is Ready")
			} else {
				utils.PrintWarning("Status: Container Is Not Ready")
			}
			fmt.Println()
		}

		utils.PrintSection("Resource Analysis")
		if report.ResourceLimitsSet {
			utils.PrintSuccess("Status: Resource Limits Configured")
		} else {
			utils.PrintWarning("Warning: No Resource Limits Configured")
		}

		if report.ResourceRequestsSet {
			utils.PrintSuccess("Status: Resource Requests Configured")
		} else {
			utils.PrintWarning("Warning: No Resource Requests Configured")
		}

		utils.PrintSection("Recent Events Analysis")
		if len(report.Events) > 0 {
			for _, event := range report.Events {
				fmt.Printf("[%s] %s: %s\n",
					event.LastTimestamp.Format("15:04:05"),
					event.Reason,
					event.Message)
			}
		} else {
			utils.PrintSuccess("Status: No Recent Events Found")
		}

		utils.PrintSection("Summary And Recommendations")
		if len(report.Issues) == 0 {
			utils.PrintSuccess("Overall Health: Healthy")
		} else {
			utils.PrintWarning("Overall Health: Needs Attention")
			fmt.Println("Warnings:")
			for _, issue := range report.Issues {
				fmt.Printf("• %s\n", issue)
			}
		}

		if len(report.Recommendations) > 0 {
			fmt.Println("Recommended Actions:")
			for _, rec := range report.Recommendations {
				fmt.Printf("• %s\n", rec)
			}
		}

		if verbose {
			utils.PrintSection("Verbose Debug Information")
			fmt.Printf("Pod UID: %s\n", report.UID)
			fmt.Printf("Pod IP: %s\n", report.PodIP)
			fmt.Printf("Service Account: %s\n", report.ServiceAccount)
			fmt.Printf("Restart Count: %d\n", report.RestartCount)
		}
	},
}

func init() {
	podCmd.Flags().StringP("namespace", "n", "default", "Namespace")
	podCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
