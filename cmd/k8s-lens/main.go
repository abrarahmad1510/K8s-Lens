package main

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/diagnostics"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version = "0.2.0"
	blue    = color.New(color.FgBlue).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
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
  k8s-lens analyze namespace production
  k8s-lens version
`),
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(createAnalyzeCommand())
	rootCmd.AddCommand(createVersionCommand())
	rootCmd.AddCommand(createCompletionCommand())
	rootCmd.AddCommand(createSetupCommand())

	if err := rootCmd.Execute(); err != nil {
		utils.PrintError("Command Execution Failed: %s", err)
		os.Exit(1)
	}
}

func createAnalyzeCommand() *cobra.Command {
	var namespace string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "analyze [resource-type] [resource-name]",
		Short: blue("Analyze A Kubernetes Resource For Issues"),
		Long: yellow(`
Analyze Kubernetes Resources And Get Intelligent Diagnostics

Supported Resource Types:
• pod, pods, po
• deployment, deployments, deploy  
• service, services, svc
• namespace, namespaces, ns
• node, nodes, no

Examples:
  k8s-lens analyze pod my-app-pod-12345
  k8s-lens analyze deployment web-api -n production
  k8s-lens analyze node worker-1 --verbose
`),
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			resourceType := args[0]
			resourceName := args[1]

			if verbose {
				utils.PrintInfo("Verbose Mode Enabled")
				utils.PrintInfo("Resource Type: %s", resourceType)
				utils.PrintInfo("Resource Name: %s", resourceName)
				utils.PrintInfo("Namespace: %s", namespace)
			}

			utils.PrintSuccess("Analyzing %s/%s In Namespace %s", resourceType, resourceName, namespace)

			// Real Kubernetes Analysis
			result, err := diagnostics.AnalyzeResource(resourceType, resourceName, namespace)
			if err != nil {
				utils.PrintError("Analysis Failed: %s", err)
				os.Exit(1)
			}

			// Display Results
			fmt.Println(result.Report)

			if verbose && len(result.Recommendations) > 0 {
				utils.PrintSection("Intelligent Recommendations")
				for _, rec := range result.Recommendations {
					utils.PrintInfo("%s", rec)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes Namespace")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable Verbose Output")
	return cmd
}

func createSetupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: blue("Setup And Verify Kubernetes Connection"),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintInfo("Setting Up K8s Lens")

			// Test Kubernetes Connection
			analyzer, err := diagnostics.NewResourceAnalyzer()
			if err != nil {
				utils.PrintError("Failed To Connect To Kubernetes: %s", err)
				os.Exit(1)
			}

			utils.PrintSuccess("Successfully Connected To Kubernetes Cluster")
			utils.PrintInfo("Testing Cluster Access")

			// Test Basic Operations
			if err := analyzer.TestConnection(); err != nil {
				utils.PrintError("Cluster Access Test Failed: %s", err)
				os.Exit(1)
			}

			// Get Cluster Info
			clusterInfo, err := analyzer.GetClusterInfo()
			if err != nil {
				utils.PrintWarning("Unable To Retrieve Cluster Version Information")
			} else {
				utils.PrintInfo("Cluster Version: %s", clusterInfo)
			}

			utils.PrintSuccess("K8s Lens Is Ready To Analyze Your Kubernetes Resources")
		},
	}
}

func createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: blue("Print Version Information"),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("K8s Lens Version: %s\n", green(version))
			fmt.Printf("Go Version: %s\n", green(utils.GetGoVersion()))
			fmt.Printf("Platform: %s\n", green(utils.GetPlatform()))
			fmt.Printf("Kubernetes Client: %s\n", green("Enabled"))
			fmt.Printf("Built With: %s\n", green("For The Kubernetes Community"))
		},
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
