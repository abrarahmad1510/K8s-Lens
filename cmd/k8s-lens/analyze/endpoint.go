package analyze

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var endpointCmd = &cobra.Command{
	Use:   "endpoint [service-name]",
	Short: "Analyze Service Endpoints",
	Long:  `Analyze and validate service endpoints connectivity.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")

		utils.PrintInfo("Analyzing endpoints for service: %s in namespace: %s", args[0], namespace)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := diagnostics.NewEndpointAnalyzer(k8sClient, namespace)
		report, err := analyzer.ValidateEndpoints(args[0])
		if err != nil {
			utils.PrintError("Error analyzing endpoints: %v", err)
			os.Exit(1)
		}

		fmt.Printf("K8s Lens Endpoint Analysis: %s\n", report.ServiceName)
		fmt.Println("---")

		utils.PrintSection("Endpoint Status")
		fmt.Printf("Namespace: %s\n", report.Namespace)
		fmt.Printf("Ready Pods: %d/%d\n", report.Analysis.ReadyPods, report.Analysis.TotalPods)
		fmt.Printf("Status: %s\n", report.Analysis.Status)

		utils.PrintSection("Pod Readiness")
		if report.Analysis.TotalPods > 0 {
			for i, pod := range report.Pods {
				ready := "Not Ready"
				if isPodReady(&pod) {
					ready = "Ready"
				}
				fmt.Printf("- Pod %d: %s (%s)\n", i+1, pod.Name, ready)
			}
		} else {
			utils.PrintWarning("No pods found matching service selector")
		}

		if len(report.Analysis.Issues) > 0 {
			utils.PrintSection("Connectivity Issues")
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
	},
}

func init() {
	endpointCmd.Flags().StringP("namespace", "n", "default", "Namespace")
}

func isPodReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
