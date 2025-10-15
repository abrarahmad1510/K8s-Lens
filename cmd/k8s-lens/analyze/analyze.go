package analyze

import (
	"github.com/spf13/cobra"
)

// AnalyzeCmd represents the analyze command
var AnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze Kubernetes resources",
	Long:  `Analyze various Kubernetes resources and provide diagnostic information.`,
}

func init() {
	// Add subcommands
	AnalyzeCmd.AddCommand(podCmd)
	AnalyzeCmd.AddCommand(deploymentCmd)
	AnalyzeCmd.AddCommand(statefulsetCmd)
}
