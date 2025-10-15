package version

import (
	"fmt"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	versionNum = "0.3.0"
	green      = color.New(color.FgGreen).SprintFunc()
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Version Information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("K8s Lens Version: %s\n", green(versionNum))
		fmt.Printf("Go Version: %s\n", green(utils.GetGoVersion()))
		fmt.Printf("Platform: %s\n", green(utils.GetPlatform()))
		fmt.Printf("Kubernetes Client: %s\n", green("Enabled"))
		fmt.Printf("Built With: %s\n", green("For The Kubernetes Community"))
	},
}
