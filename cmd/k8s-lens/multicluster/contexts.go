package multicluster

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/multicluster"
	"github.com/spf13/cobra"
)

var contextsCmd = &cobra.Command{
	Use:   "contexts",
	Short: "List available Kubernetes contexts",
	Long:  `List all available Kubernetes contexts from kubeconfig.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintInfo("Loading available Kubernetes contexts")

		manager := multicluster.NewClusterManager()
		err := manager.LoadContexts()
		if err != nil {
			utils.PrintError("Error loading cluster contexts: %v", err)
			os.Exit(1)
		}

		contexts := manager.ListContexts()
		currentContext, err := manager.GetCurrentContext()

		fmt.Println("Available Kubernetes Contexts:")
		fmt.Println("===============================")

		for i, context := range contexts {
			prefix := "  "
			if currentContext != nil && context == currentContext.Name {
				prefix = "→ "
			}
			fmt.Printf("%s%d. %s\n", prefix, i+1, context)
		}

		if currentContext != nil {
			fmt.Printf("\nCurrent Context: %s\n", currentContext.Name)
		}

		fmt.Printf("\nTotal Contexts: %d\n", len(contexts))
	},
}
