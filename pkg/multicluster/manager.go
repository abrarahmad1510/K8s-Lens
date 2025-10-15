package multicluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ClusterManager manages multiple Kubernetes clusters
type ClusterManager struct {
	contexts       map[string]*ClusterContext
	currentContext string
}

// ClusterContext represents a Kubernetes cluster context
type ClusterContext struct {
	Name   string
	Client kubernetes.Interface
	Config clientcmd.ClientConfig
}

// NewClusterManager creates a new ClusterManager
func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		contexts: make(map[string]*ClusterContext),
	}
}

// LoadContexts loads all available Kubernetes contexts
func (c *ClusterManager) LoadContexts() error {
	kubeconfig := getKubeconfigPath()
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	for contextName := range config.Contexts {
		client, clientConfig, err := c.createClientForContext(contextName)
		if err != nil {
			fmt.Printf("Warning: Failed to create client for context %s: %v\n", contextName, err)
			continue
		}

		c.contexts[contextName] = &ClusterContext{
			Name:   contextName,
			Client: client,
			Config: clientConfig,
		}
	}

	// Set current context
	c.currentContext = config.CurrentContext

	return nil
}

// SwitchContext switches the current context
func (c *ClusterManager) SwitchContext(contextName string) error {
	if _, exists := c.contexts[contextName]; !exists {
		return fmt.Errorf("context %s not found", contextName)
	}
	c.currentContext = contextName
	return nil
}

// GetCurrentContext returns the current cluster context
func (c *ClusterManager) GetCurrentContext() (*ClusterContext, error) {
	return c.GetContext(c.currentContext)
}

// GetContext returns a cluster context by name
func (c *ClusterManager) GetContext(contextName string) (*ClusterContext, error) {
	context, exists := c.contexts[contextName]
	if !exists {
		return nil, fmt.Errorf("context %s not found", contextName)
	}
	return context, nil
}

// ListContexts returns all available contexts
func (c *ClusterManager) ListContexts() []string {
	var contextNames []string
	for name := range c.contexts {
		contextNames = append(contextNames, name)
	}
	return contextNames
}

// CompareClusters compares resources across clusters
func (c *ClusterManager) CompareClusters(resourceType string) (*ClusterComparison, error) {
	comparison := &ClusterComparison{
		ResourceType: resourceType,
		ClusterData:  make(map[string]ClusterResources),
	}

	for contextName, context := range c.contexts {
		resources, err := c.getResourcesForType(context.Client, resourceType)
		if err != nil {
			return nil, fmt.Errorf("failed to get resources for %s in context %s: %v", resourceType, contextName, err)
		}

		comparison.ClusterData[contextName] = resources
	}

	comparison.analyzeDifferences()
	return comparison, nil
}

// FederatedAnalysis performs analysis across all clusters
func (c *ClusterManager) FederatedAnalysis() (*FederatedReport, error) {
	report := &FederatedReport{
		ClusterReports: make(map[string]ClusterReport),
	}

	for contextName, context := range c.contexts {
		clusterReport, err := c.analyzeCluster(context)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze cluster %s: %v", contextName, err)
		}
		report.ClusterReports[contextName] = *clusterReport
	}

	report.generateSummary()
	return report, nil
}

func (c *ClusterManager) createClientForContext(contextName string) (kubernetes.Interface, clientcmd.ClientConfig, error) {
	kubeconfig := getKubeconfigPath()
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	overrides := &clientcmd.ConfigOverrides{
		CurrentContext: contextName,
	}

	clientConfig := clientcmd.NewNonInteractiveClientConfig(*config, contextName, overrides, nil)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}

	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}

	return client, clientConfig, nil
}

func (c *ClusterManager) getResourcesForType(client kubernetes.Interface, resourceType string) (ClusterResources, error) {
	resources := ClusterResources{}

	switch resourceType {
	case "pods":
		pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return resources, err
		}
		resources.Pods = pods.Items
		resources.Count = len(pods.Items)

	case "nodes":
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return resources, err
		}
		resources.Nodes = nodes.Items
		resources.Count = len(nodes.Items)

	case "deployments":
		deployments, err := client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return resources, err
		}
		resources.Count = len(deployments.Items)

	default:
		return resources, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return resources, nil
}

func (c *ClusterManager) analyzeCluster(clusterContext *ClusterContext) (*ClusterReport, error) {
	// Get basic cluster info
	nodes, err := clusterContext.Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := clusterContext.Client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Calculate cluster health
	healthyNodes := 0
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				healthyNodes++
				break
			}
		}
	}

	report := &ClusterReport{
		Name:         clusterContext.Name,
		TotalNodes:   len(nodes.Items),
		HealthyNodes: healthyNodes,
		TotalPods:    len(pods.Items),
		HealthStatus: "Healthy",
	}

	if healthyNodes < len(nodes.Items) {
		report.HealthStatus = "Degraded"
	}

	return report, nil
}

func getKubeconfigPath() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}
	return filepath.Join(homedir.HomeDir(), ".kube", "config")
}
