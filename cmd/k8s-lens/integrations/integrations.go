package integrations

import (
	"github.com/spf13/cobra"
)

// IntegrationsCmd represents the integrations command
var IntegrationsCmd = &cobra.Command{
	Use:   "integrations",
	Short: "Third-party integrations and enhanced metrics",
	Long:  "Access integrations with Prometheus and other monitoring systems",
}

func init() {
	IntegrationsCmd.AddCommand(metricsCmd)
}
