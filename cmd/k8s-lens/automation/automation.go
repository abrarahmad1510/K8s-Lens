package automation

import (
	"github.com/spf13/cobra"
)

// AutomationCmd represents the automation command
var AutomationCmd = &cobra.Command{
	Use:   "automation",
	Short: "Automated remediation and self-healing",
	Long:  "Automated issue remediation, self-healing, and predictive scaling for Kubernetes resources",
}

var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Automated remediation commands",
}

var scaleCmd = &cobra.Command{
	Use:   "scale", 
	Short: "Predictive scaling commands",
}

var healCmd = &cobra.Command{
	Use:   "heal",
	Short: "Self-healing commands",
}

func init() {
	AutomationCmd.AddCommand(remediateCmd)
	AutomationCmd.AddCommand(scaleCmd)
	AutomationCmd.AddCommand(healCmd)
}
