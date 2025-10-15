package optimize

import (
	"github.com/spf13/cobra"
)

// OptimizeCmd represents the optimize command
var OptimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize Kubernetes resources and performance",
	Long:  `Optimize Kubernetes resources, predict failures, and generate automated fixes.`,
}

func init() {
	// Add subcommands
	OptimizeCmd.AddCommand(resourceCmd)
	OptimizeCmd.AddCommand(predictCmd)
	OptimizeCmd.AddCommand(fixCmd)
}
