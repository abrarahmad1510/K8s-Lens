package analyze

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Analyze a Kubernetes Service",
	Long:  `Analyze a Kubernetes Service and provide diagnostic information.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		verbose, _ := cmd.Flags().GetBool("verbose")

		utils.PrintInfo("Starting service analysis for: %s in namespace: %s", args[0], namespace)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := diagnostics.NewServiceAnalyzer(k8sClient, namespace)
		report, err := analyzer.Analyze(args[0])
		if err != nil {
			utils.PrintError("Error analyzing service: %v", err)
			os.Exit(1)
		}

		// Print the report
		fmt.Printf("K8s Lens Analysis Report For Service: %s\n", report.Name)
		fmt.Println("---")

		utils.PrintSection("Service Configuration")
		fmt.Printf("Namespace: %s\n", report.Namespace)
		fmt.Printf("Type: %s\n", report.Type)
		fmt.Printf("Cluster IP: %s\n", report.ClusterIP)
		if report.ExternalIP != "" {
			fmt.Printf("External IP: %s\n", report.ExternalIP)
		}

		utils.PrintSection("Port Configuration")
		if len(report.Ports) > 0 {
			for _, port := range report.Ports {
				fmt.Printf("- Port: %d/%s -> TargetPort: %v\n", port.Port, port.Protocol, port.TargetPort)
			}
		} else {
			utils.PrintWarning("No ports configured")
		}

		utils.PrintSection("Selector")
		if len(report.Selector) > 0 {
			for key, value := range report.Selector {
				fmt.Printf("- %s: %s\n", key, value)
			}
		} else {
			utils.PrintWarning("No selector configured")
		}

		utils.PrintSection("Endpoints Analysis")
		if report.Endpoints != nil {
			totalAddresses := 0
			for _, subset := range report.Endpoints.Subsets {
				totalAddresses += len(subset.Addresses)
			}
			if totalAddresses > 0 {
				utils.PrintSuccess("Active endpoints: %d", totalAddresses)
				for _, subset := range report.Endpoints.Subsets {
					for _, address := range subset.Addresses {
						fmt.Printf("- %s\n", address.IP)
					}
				}
			} else {
				utils.PrintWarning("No active endpoints found")
			}
		} else {
			utils.PrintWarning("No endpoints found")
		}

		utils.PrintSection("Status")
		fmt.Printf("Overall Status: %s\n", report.Analysis.Status)

		if len(report.Analysis.Issues) > 0 {
			utils.PrintSection("Issues")
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

		if verbose {
			utils.PrintSection("Verbose Information")
			fmt.Println("Recent Events:")
			if len(report.Events) > 0 {
				for _, event := range report.Events {
					fmt.Printf("- [%s] %s: %s\n",
						event.LastTimestamp.Format("15:04:05"),
						event.Reason,
						event.Message)
				}
			} else {
				fmt.Println("  No recent events")
			}
		}
	},
}

func init() {
	serviceCmd.Flags().StringP("namespace", "n", "default", "Namespace")
	serviceCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
