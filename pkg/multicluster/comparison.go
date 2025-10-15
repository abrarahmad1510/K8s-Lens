package multicluster

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// ClusterComparison contains comparison data across clusters
type ClusterComparison struct {
	ResourceType string
	ClusterData  map[string]ClusterResources
	Differences  []ClusterDifference
	Summary      ComparisonSummary
}

// ClusterResources represents resources in a cluster
type ClusterResources struct {
	Pods  []corev1.Pod
	Nodes []corev1.Node
	Count int
}

// ClusterDifference represents a difference between clusters
type ClusterDifference struct {
	ResourceType string
	ClusterA     string
	ClusterB     string
	Field        string
	ValueA       interface{}
	ValueB       interface{}
	Severity     string
}

// ComparisonSummary provides a summary of the comparison
type ComparisonSummary struct {
	TotalClusters  int
	TotalResources int
	MaxResources   int
	MinResources   int
	Differences    int
}

// ClusterReport contains analysis of a single cluster
type ClusterReport struct {
	Name         string
	TotalNodes   int
	HealthyNodes int
	TotalPods    int
	HealthStatus string
}

// FederatedReport contains analysis across all clusters
type FederatedReport struct {
	ClusterReports map[string]ClusterReport
	Summary        FederatedSummary
}

// FederatedSummary provides overall summary
type FederatedSummary struct {
	TotalClusters   int
	HealthyClusters int
	TotalNodes      int
	TotalPods       int
	OverallHealth   string
}

func (c *ClusterComparison) analyzeDifferences() {
	clusters := make([]string, 0, len(c.ClusterData))
	for cluster := range c.ClusterData {
		clusters = append(clusters, cluster)
	}

	// Compare each pair of clusters
	for i := 0; i < len(clusters); i++ {
		for j := i + 1; j < len(clusters); j++ {
			c.compareClusterPair(clusters[i], clusters[j])
		}
	}

	// Calculate summary
	c.Summary.TotalClusters = len(c.ClusterData)
	c.Summary.Differences = len(c.Differences)

	maxResources := 0
	minResources := -1
	totalResources := 0

	for _, resources := range c.ClusterData {
		totalResources += resources.Count
		if resources.Count > maxResources {
			maxResources = resources.Count
		}
		if minResources == -1 || resources.Count < minResources {
			minResources = resources.Count
		}
	}

	c.Summary.MaxResources = maxResources
	c.Summary.MinResources = minResources
	c.Summary.TotalResources = totalResources
}

func (c *ClusterComparison) compareClusterPair(clusterA, clusterB string) {
	resourcesA := c.ClusterData[clusterA]
	resourcesB := c.ClusterData[clusterB]

	// Compare resource counts
	if resourcesA.Count != resourcesB.Count {
		severity := "low"
		diff := abs(resourcesA.Count - resourcesB.Count)
		if diff > 10 {
			severity = "high"
		} else if diff > 5 {
			severity = "medium"
		}

		c.Differences = append(c.Differences, ClusterDifference{
			ResourceType: c.ResourceType,
			ClusterA:     clusterA,
			ClusterB:     clusterB,
			Field:        "count",
			ValueA:       resourcesA.Count,
			ValueB:       resourcesB.Count,
			Severity:     severity,
		})
	}
}

func (f *FederatedReport) generateSummary() {
	f.Summary.TotalClusters = len(f.ClusterReports)
	f.Summary.TotalNodes = 0
	f.Summary.TotalPods = 0
	f.Summary.HealthyClusters = 0

	for _, report := range f.ClusterReports {
		f.Summary.TotalNodes += report.TotalNodes
		f.Summary.TotalPods += report.TotalPods
		if report.HealthStatus == "Healthy" {
			f.Summary.HealthyClusters++
		}
	}

	if f.Summary.HealthyClusters == f.Summary.TotalClusters {
		f.Summary.OverallHealth = "Healthy"
	} else if f.Summary.HealthyClusters > f.Summary.TotalClusters/2 {
		f.Summary.OverallHealth = "Degraded"
	} else {
		f.Summary.OverallHealth = "Critical"
	}
}

// GenerateReport generates a human-readable comparison report
func (c *ClusterComparison) GenerateReport() string {
	report := fmt.Sprintf("Multi-Cluster Comparison Report: %s\n", c.ResourceType)
	report += "============================================\n\n"

	report += "Cluster Resource Counts:\n"
	for clusterName, resources := range c.ClusterData {
		report += fmt.Sprintf("  %s: %d %s\n", clusterName, resources.Count, c.ResourceType)
	}

	if len(c.Differences) > 0 {
		report += "\nDifferences Found:\n"
		for _, diff := range c.Differences {
			report += fmt.Sprintf("  - %s: %s (%v) vs %s (%v) [%s]\n",
				diff.Field, diff.ClusterA, diff.ValueA, diff.ClusterB, diff.ValueB, diff.Severity)
		}
	} else {
		report += "\nNo significant differences found across clusters.\n"
	}

	report += fmt.Sprintf("\nSummary:\n")
	report += fmt.Sprintf("  Total Clusters: %d\n", c.Summary.TotalClusters)
	report += fmt.Sprintf("  Total Resources: %d\n", c.Summary.TotalResources)
	report += fmt.Sprintf("  Differences Found: %d\n", c.Summary.Differences)

	return report
}

// GenerateFederatedReport generates a federated analysis report
func (f *FederatedReport) GenerateFederatedReport() string {
	report := "Federated Cluster Analysis Report\n"
	report += "==================================\n\n"

	for clusterName, clusterReport := range f.ClusterReports {
		report += fmt.Sprintf("Cluster: %s\n", clusterName)
		report += fmt.Sprintf("  Nodes: %d/%d healthy\n", clusterReport.HealthyNodes, clusterReport.TotalNodes)
		report += fmt.Sprintf("  Pods: %d\n", clusterReport.TotalPods)
		report += fmt.Sprintf("  Status: %s\n", clusterReport.HealthStatus)
		report += "  ---\n"
	}

	report += fmt.Sprintf("\nOverall Summary:\n")
	report += fmt.Sprintf("  Total Clusters: %d\n", f.Summary.TotalClusters)
	report += fmt.Sprintf("  Healthy Clusters: %d\n", f.Summary.HealthyClusters)
	report += fmt.Sprintf("  Total Nodes: %d\n", f.Summary.TotalNodes)
	report += fmt.Sprintf("  Total Pods: %d\n", f.Summary.TotalPods)
	report += fmt.Sprintf("  Overall Health: %s\n", f.Summary.OverallHealth)

	return report
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
