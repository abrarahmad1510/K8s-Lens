package multicluster

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/multicluster"
	"github.com/spf13/cobra"
)

var federatedCmd = &cobra.Command{
	Use:   "federated",
	Short: "Run federated analysis across all clusters",
	Long:  `Perform comprehensive analysis across all available Kubernetes clusters.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintInfo("Running federated analysis across all clusters")

		manager := multicluster.NewClusterManager()
		err := manager.LoadContexts()
		if err != nil {
			utils.PrintError("Error loading cluster contexts: %v", err)
			os.Exit(1)
		}

		report, err := manager.FederatedAnalysis()
		if err != nil {
			utils.PrintError("Error running federated analysis: %v", err)
			os.Exit(1)
		}

		fmt.Println(report.GenerateFederatedReport())
	},
}
