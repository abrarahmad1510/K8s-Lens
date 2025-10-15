package test

import (
	"fmt"

	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/spf13/cobra"
)

// TestCmd represents the test command
var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run mock analysis for testing",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running K8s Lens Test Mode")
		fmt.Println("=== Mock Analysis Demonstration ===")

		// Run mock analysis
		mockAnalyzer := diagnostics.NewMockAnalyzer()
		report := mockAnalyzer.Analyze("test-pod")

		fmt.Println(report)
	},
}
