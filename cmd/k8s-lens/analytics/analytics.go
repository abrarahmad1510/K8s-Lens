package analytics

import (
	"github.com/spf13/cobra"
)

// AnalyticsCmd represents the analytics command
var AnalyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Advanced analytics and machine learning features",
	Long:  "Access advanced machine learning and predictive analytics for your Kubernetes cluster",
}

func init() {
	AnalyticsCmd.AddCommand(anomalyCmd)
	AnalyticsCmd.AddCommand(predictCmd)
	AnalyticsCmd.AddCommand(trendCmd)
}
