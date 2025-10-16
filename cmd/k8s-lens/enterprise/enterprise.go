package enterprise

import "github.com/spf13/cobra"

// EnterpriseCmd represents the enterprise command
var EnterpriseCmd = &cobra.Command{
	Use:   "enterprise",
	Short: "Enterprise security and RBAC analysis",
	Long:  "Advanced security scanning, RBAC analysis, and compliance reporting for enterprise Kubernetes clusters",
}

var rbacCmd = &cobra.Command{
	Use:   "rbac",
	Short: "RBAC analysis commands",
}

var securityCmd = &cobra.Command{
	Use:   "security", 
	Short: "Security scanning commands",
}

func init() {
	EnterpriseCmd.AddCommand(rbacCmd)
	EnterpriseCmd.AddCommand(securityCmd)
}
