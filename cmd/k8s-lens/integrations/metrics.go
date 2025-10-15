package integrations

import (
	"fmt"
	"os"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
	"github.com/abrarahmad1510/k8s-lens/pkg/integrations"
	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics [resource-type] [resource-name]",
	Short: "Analyze resources with Prometheus metrics",
	Long:  "Enhanced analysis using Prometheus metrics for pods, nodes, and clusters",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		resourceType := args[0]
		var resourceName string
		if len(args) > 1 {
			resourceName = args[1]
		}

		namespace, _ := cmd.Flags().GetString("namespace")
		prometheusURL, _ := cmd.Flags().GetString("prometheus-url")

		utils.PrintInfo("Starting metrics analysis for %s: %s", resourceType, resourceName)

		k8sClient, err := k8s.NewClient()
		if err != nil {
			utils.PrintError("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		analyzer := integrations.NewMetricsAnalyzer(k8sClient, prometheusURL)

		switch resourceType {
		case "pod", "pods":
			if resourceName == "" {
				utils.PrintError("Pod name is required for pod metrics analysis")
				os.Exit(1)
			}
			report, err := analyzer.AnalyzePodWithMetrics(resourceName, namespace)
			if err != nil {
				utils.PrintError("Error analyzing pod with metrics: %v", err)
				os.Exit(1)
			}
			printPodMetricsReport(report)

		case "node", "nodes":
			if resourceName == "" {
				utils.PrintError("Node name is required for node metrics analysis")
				os.Exit(1)
			}
			metrics, err := analyzer.AnalyzeNodeWithMetrics(resourceName)
			if err != nil {
				utils.PrintError("Error analyzing node with metrics: %v", err)
				os.Exit(1)
			}
			printNodeMetricsReport(metrics)

		case "cluster":
			metrics, err := analyzer.AnalyzeClusterWithMetrics()
			if err != nil {
				utils.PrintError("Error analyzing cluster with metrics: %v", err)
				os.Exit(1)
			}
			printClusterMetricsReport(metrics)

		default:
			utils.PrintError("Unsupported resource type: %s. Supported types: pod, node, cluster", resourceType)
			os.Exit(1)
		}
	},
}

func printPodMetricsReport(report *integrations.EnhancedPodReport) {
	fmt.Printf("K8s Lens Enhanced Pod Analysis Report\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("Pod: %s\n", report.PodReport.Name)
	fmt.Printf("Namespace: %s\n", report.PodReport.Namespace)
	fmt.Printf("Health Score: %d/100\n", report.HealthScore)

	utils.PrintSection("Resource Metrics")

	if report.PodMetrics.Error != "" {
		utils.PrintWarning("Metrics Collection Issues: %s", report.PodMetrics.Error)
		fmt.Println("\nTo fix this:")
		fmt.Println("1. Ensure Prometheus is running in your cluster")
		fmt.Println("2. Use the correct Prometheus URL with --prometheus-url flag")
		fmt.Println("3. Set up port forwarding: kubectl port-forward -n monitoring service/prometheus-operated 9090:9090")
	} else {
		fmt.Printf("CPU Usage: %.3f cores\n", report.PodMetrics.CPUUsage)
		fmt.Printf("Memory Usage: %.2f MB\n", report.PodMetrics.MemoryUsage/(1024*1024))
		fmt.Printf("Network Receive: %.2f KB/s\n", report.PodMetrics.NetworkRx/1024)
		fmt.Printf("Network Transmit: %.2f KB/s\n", report.PodMetrics.NetworkTx/1024)
	}

	utils.PrintSection("Pod Status")
	fmt.Printf("Phase: %s\n", report.PodReport.Phase)
	fmt.Printf("Status: %s\n", report.PodReport.Status)
	fmt.Printf("Restart Count: %d\n", report.PodReport.RestartCount)

	if len(report.PodReport.Issues) > 0 {
		utils.PrintSection("Issues")
		for _, issue := range report.PodReport.Issues {
			utils.PrintWarning("- %s", issue)
		}
	}

	utils.PrintSection("Recommendations")
	if len(report.Recommendations) > 0 {
		for i, rec := range report.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec)
		}
	} else {
		fmt.Println("No recommendations - pod is healthy!")
	}
}

func printNodeMetricsReport(metrics *integrations.NodeMetrics) {
	fmt.Printf("K8s Lens Node Metrics Report\n")
	fmt.Printf("============================\n")
	fmt.Printf("Node: %s\n", metrics.NodeName)
	fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))

	utils.PrintSection("Resource Utilization")
	if metrics.Error != "" {
		utils.PrintWarning("Metrics Collection Issues: %s", metrics.Error)
		fmt.Println("\nCommon fixes:")
		fmt.Println("1. Ensure node-exporter is running in your cluster")
		fmt.Println("2. Check if Prometheus has proper node monitoring")
		fmt.Println("3. Verify node name matches Prometheus labels")
	} else {
		fmt.Printf("CPU Usage: %.1f%%\n", metrics.CPUUsage)
		fmt.Printf("Memory Usage: %.1f%%\n", metrics.MemoryUsage)
		fmt.Printf("Disk Usage: %.1f%%\n", metrics.DiskUsage)
		fmt.Printf("Pod Count: %d\n", metrics.PodCount)

		utils.PrintSection("Health Assessment")
		if metrics.CPUUsage > 80 {
			utils.PrintWarning("High CPU usage - consider scaling or optimization")
		} else {
			utils.PrintSuccess("CPU usage: %.1f%%", metrics.CPUUsage)
		}

		if metrics.MemoryUsage > 85 {
			utils.PrintWarning("High memory usage - monitor for memory pressure")
		} else {
			utils.PrintSuccess("Memory usage: %.1f%%", metrics.MemoryUsage)
		}

		if metrics.DiskUsage > 90 {
			utils.PrintWarning("High disk usage - consider cleanup or expansion")
		} else {
			utils.PrintSuccess("Disk usage: %.1f%%", metrics.DiskUsage)
		}

		if metrics.PodCount > 100 {
			utils.PrintWarning("High pod density - consider node scaling")
		} else {
			utils.PrintSuccess("Pod count: %d", metrics.PodCount)
		}
	}
}

func printClusterMetricsReport(metrics *integrations.ClusterMetrics) {
	fmt.Printf("K8s Lens Cluster Metrics Report\n")
	fmt.Printf("===============================\n")
	fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))

	utils.PrintSection("Cluster Overview")
	if metrics.Error != "" {
		utils.PrintWarning("Metrics Collection Issues: %s", metrics.Error)
		fmt.Println("\nTo fix this:")
		fmt.Println("1. Ensure Prometheus is running in your cluster")
		fmt.Println("2. Use the correct Prometheus URL with --prometheus-url flag")
		fmt.Println("3. Set up port forwarding: kubectl port-forward -n monitoring service/prometheus-operated 9090:9090")
	} else {
		fmt.Printf("Total Nodes: %d\n", metrics.TotalNodes)
		fmt.Printf("Total Pods: %d\n", metrics.TotalPods)

		utils.PrintSection("Resource Capacity & Usage")
		if metrics.CPUCapacity > 0 {
			fmt.Printf("CPU Capacity: %.1f cores\n", metrics.CPUCapacity)
			fmt.Printf("CPU Usage: %.1f cores (%.1f%%)\n", metrics.CPUUsage, (metrics.CPUUsage/metrics.CPUCapacity)*100)
		} else {
			fmt.Printf("CPU Capacity: Not available\n")
			fmt.Printf("CPU Usage: %.1f cores\n", metrics.CPUUsage)
		}

		if metrics.MemoryCapacity > 0 {
			fmt.Printf("Memory Capacity: %.1f GB\n", metrics.MemoryCapacity)
			fmt.Printf("Memory Usage: %.1f GB (%.1f%%)\n", metrics.MemoryUsage, (metrics.MemoryUsage/metrics.MemoryCapacity)*100)
		} else {
			fmt.Printf("Memory Capacity: Not available\n")
			fmt.Printf("Memory Usage: %.1f GB\n", metrics.MemoryUsage)
		}

		utils.PrintSection("Cluster Health")
		if metrics.CPUCapacity > 0 {
			cpuUtilization := (metrics.CPUUsage / metrics.CPUCapacity) * 100
			if cpuUtilization > 80 {
				utils.PrintWarning("High cluster CPU utilization: %.1f%%", cpuUtilization)
			} else if cpuUtilization > 60 {
				utils.PrintInfo("Moderate cluster CPU utilization: %.1f%%", cpuUtilization)
			} else {
				utils.PrintSuccess("Good cluster CPU utilization: %.1f%%", cpuUtilization)
			}
		}

		if metrics.MemoryCapacity > 0 {
			memoryUtilization := (metrics.MemoryUsage / metrics.MemoryCapacity) * 100
			if memoryUtilization > 80 {
				utils.PrintWarning("High cluster memory utilization: %.1f%%", memoryUtilization)
			} else if memoryUtilization > 60 {
				utils.PrintInfo("Moderate cluster memory utilization: %.1f%%", memoryUtilization)
			} else {
				utils.PrintSuccess("Good cluster memory utilization: %.1f%%", memoryUtilization)
			}
		}
	}
}

func init() {
	metricsCmd.Flags().StringP("namespace", "n", "default", "Namespace (for pods)")
	metricsCmd.Flags().StringP("prometheus-url", "p", "http://localhost:9090", "Prometheus URL")
}
