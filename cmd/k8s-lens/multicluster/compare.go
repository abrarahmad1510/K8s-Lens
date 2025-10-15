package multicluster

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/multicluster"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare [resource-type]",
	Short: "Compare resources across clusters",
	Long:  `Compare Kubernetes resources across all available clusters.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]

		utils.PrintInfo("Comparing %s across all clusters", resourceType)

		manager := multicluster.NewClusterManager()
		err := manager.LoadContexts()
		if err != nil {
			utils.PrintError("Error loading cluster contexts: %v", err)
			os.Exit(1)
		}

		comparison, err := manager.CompareClusters(resourceType)
		if err != nil {
			utils.PrintError("Error comparing clusters: %v", err)
			os.Exit(1)
		}

		fmt.Println(comparison.GenerateReport())
	},
}
