package multicluster

import (
	"github.com/spf13/cobra"
)

// MulticlusterCmd represents the multicluster command
var MulticlusterCmd = &cobra.Command{
	Use:   "multicluster",
	Short: "Multi-cluster management and analysis",
	Long:  `Manage and analyze multiple Kubernetes clusters.`,
}

func init() {
	MulticlusterCmd.AddCommand(contextsCmd)
	MulticlusterCmd.AddCommand(compareCmd)
	MulticlusterCmd.AddCommand(federatedCmd)
}
