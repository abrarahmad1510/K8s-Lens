package setup

import (
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/spf13/cobra"
)

// SetupCmd represents the setup command
var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup And Verify Kubernetes Connection",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintInfo("Setting Up K8s Lens")

		// Test Kubernetes Connection
		analyzer, err := diagnostics.NewResourceAnalyzer()
		if err != nil {
			utils.PrintError("Failed To Connect To Kubernetes: %s", err)
			os.Exit(1)
		}

		utils.PrintSuccess("Successfully Connected To Kubernetes Cluster")
		utils.PrintInfo("Testing Cluster Access")

		// Test Basic Operations
		if err := analyzer.TestConnection(); err != nil {
			utils.PrintError("Cluster Access Test Failed: %s", err)
			os.Exit(1)
		}

		// Get Cluster Info
		clusterInfo, err := analyzer.GetClusterInfo()
		if err != nil {
			utils.PrintWarning("Unable To Retrieve Cluster Version Information")
		} else {
			utils.PrintInfo("Cluster Version: %s", clusterInfo)
		}

		utils.PrintSuccess("K8s Lens Is Ready To Analyze Your Kubernetes Resources")
	},
}
