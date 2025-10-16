package automation

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/automation"
	"github.com/abrarahmad1510/k8s-lens/pkg/automation/remediators"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

func init() {
	podRemediateCmd := &cobra.Command{
		Use:   "pod [pod-name] [issue-type]",
		Short: "Remediate pod issues automatically",
		Args:  cobra.ExactArgs(2),
		Run:   remediatePod,
	}
	podRemediateCmd.Flags().StringP("namespace", "n", "default", "Namespace of the pod")
	
	remediateCmd.AddCommand(podRemediateCmd)

	remediateCmd.AddCommand(&cobra.Command{
		Use:   "list-actions",
		Short: "List available remediation actions",
		Run:   listRemediationActions,
	})
}

func remediatePod(cmd *cobra.Command, args []string) {
	podName := args[0]
	issueType := args[1]
	namespace, _ := cmd.Flags().GetString("namespace")

	utils.PrintInfo("Attempting automated remediation for pod %s (issue: %s) in namespace %s", podName, issueType, namespace)
	
	k8sClient, err := k8s.NewClient()
	if err != nil {
		utils.PrintError("Error creating Kubernetes client: %v", err)
		os.Exit(1)
	}

	// Create automation engine and register remediators
	engine := automation.NewAutomationEngine(k8sClient)
	engine.RegisterRemediator(remediators.NewPodRestartRemediator(k8sClient))

	result, err := engine.AutoRemediate(cmd.Context(), issueType, podName, namespace)
	if err != nil {
		utils.PrintError("Remediation failed: %v", err)
		os.Exit(1)
	}

	if result.Success {
		utils.PrintSuccess("Remediation successful!")
		fmt.Printf("Action: %s\n", result.Action)
		fmt.Printf("Resource: %s\n", result.Resource)
		fmt.Printf("Message: %s\n", result.Message)
		fmt.Printf("Duration: %v\n", result.Duration)
	} else {
		utils.PrintWarning("Remediation attempted but didn't succeed")
		fmt.Printf("Message: %s\n", result.Message)
	}
}

func listRemediationActions(cmd *cobra.Command, args []string) {
	fmt.Printf("Available Remediation Actions\n")
	fmt.Printf("=============================\n")
	
	k8sClient, err := k8s.NewClient()
	if err != nil {
		utils.PrintError("Error creating Kubernetes client: %v", err)
		return
	}

	engine := automation.NewAutomationEngine(k8sClient)
	engine.RegisterRemediator(remediators.NewPodRestartRemediator(k8sClient))

	fmt.Printf("\nPod Remediation Actions:\n")
	podRemediator := remediators.NewPodRestartRemediator(k8sClient)
	for _, action := range podRemediator.GetRemediationActions() {
		fmt.Printf("  • %s: %s (Risk: %s)\n", action.Type, action.Description, action.Risk)
		fmt.Printf("    Command: %s\n", action.Command)
	}
}
