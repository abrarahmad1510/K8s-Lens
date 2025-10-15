package main

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/analytics"
	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/analyze"
	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/integrations"
	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/multicluster"
	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/optimize"
	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/setup"
	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/test"
	"github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens/version"
	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	versionNum = "1.0.0"
	blue       = color.New(color.FgBlue).SprintFunc()
	green      = color.New(color.FgGreen).SprintFunc()
	red        = color.New(color.FgRed).SprintFunc()
	yellow     = color.New(color.FgYellow).SprintFunc()
)

func main() {
	// Print ASCII Art Banner For Non-Completion Commands
	if len(os.Args) > 1 && os.Args[1] != "completion" {
		fig := figure.NewFigure("K8s Lens", "slant", true)
		fig.Print()
		fmt.Println()
	}

	var rootCmd = &cobra.Command{
		Use:   "k8s-lens",
		Short: blue("AI Powered Kubernetes Troubleshooting Assistant"),
		Long: yellow(`
K8s Lens Provides Intelligent Diagnostics And Recommendations For Kubernetes Issues

Features:
• AI Powered Analysis Of Kubernetes Resources
• Comprehensive Health Reports With Actionable Insights  
• Deep Dive Into Pod, Deployment, And Service Issues
• Intelligent Recommendations Based On SRE Best Practices
• Multi Cluster Support And Real Time Monitoring

Examples:
  k8s-lens analyze pod my-app-pod
  k8s-lens analyze deployment my-web-service
  k8s-lens analyze statefulset database
  k8s-lens setup
  k8s-lens version
`),
		Version:       versionNum,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add commands from the new command structure
	rootCmd.AddCommand(analyze.AnalyzeCmd)
	rootCmd.AddCommand(setup.SetupCmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(test.TestCmd)
	rootCmd.AddCommand(createCompletionCommand())
	rootCmd.AddCommand(optimize.OptimizeCmd)
	rootCmd.AddCommand(multicluster.MulticlusterCmd)
	rootCmd.AddCommand(analytics.AnalyticsCmd)
	rootCmd.AddCommand(integrations.IntegrationsCmd)

	if err := rootCmd.Execute(); err != nil {
		utils.PrintError("Command Execution Failed: %s", err)
		os.Exit(1)
	}
}

func createCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: blue("Generate Completion Script"),
		Long: yellow(`
Generate Completion Script For K8s Lens

To Load Completions:

Bash:
  $ source <(k8s-lens completion bash)
  # To Load Completions For Each Session, Execute Once:
  # Linux:
  $ k8s-lens completion bash > /etc/bash_completion.d/k8s-lens
  # macOS:
  $ k8s-lens completion bash > /usr/local/etc/bash_completion.d/k8s-lens

Zsh:
  # If Shell Completion Is Not Already Enabled In Your Environment,
  # You Will Need To Enable It. You Can Execute The Following Once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To Load Completions For Each Session, Execute Once:
  $ k8s-lens completion zsh > "${fpath[1]}/_k8s-lens"

  # You Will Need To Start A New Shell For This Setup To Take Effect.

Fish:
  $ k8s-lens completion fish | source

  # To Load Completions For Each Session, Execute Once:
  $ k8s-lens completion fish > ~/.config/fish/completions/k8s-lens.fish

PowerShell:
  PS> k8s-lens completion powershell | Out-String | Invoke-Expression

  # To Load Completions For Every New Session, Run:
  PS> k8s-lens completion powershell > k8s-lens.ps1
  # And Source This File From Your PowerShell Profile.
`),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
}
