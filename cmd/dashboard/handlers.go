package main

import (
	"context"
	"net/http"

	"github.com/abrarahmad1510/k8s-lens/pkg/k8s"
	"github.com/abrarahmad1510/k8s-lens/pkg/multicluster"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
	})
}

func clusterInfoHandler(c *gin.Context) {
	client, err := k8s.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	version, err := client.GetServerVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get node count
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	nodeCount := 0
	if err == nil {
		nodeCount = len(nodes.Items)
	}

	// Get pod count
	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	podCount := 0
	if err == nil {
		podCount = len(pods.Items)
	}

	c.JSON(http.StatusOK, gin.H{
		"clusterVersion": version,
		"nodes":          nodeCount,
		"pods":           podCount,
		"status":         "healthy",
	})
}

func analysisHandler(c *gin.Context) {
	resourceType := c.Param("resourceType")
	resourceName := c.Param("resourceName")
	namespace := c.DefaultQuery("namespace", "default")

	// This would integrate with existing analysis capabilities
	c.JSON(http.StatusOK, gin.H{
		"resourceType": resourceType,
		"resourceName": resourceName,
		"namespace":    namespace,
		"analysis":     "Analysis results would be here",
		"status":       "completed",
	})
}

func optimizationHandler(c *gin.Context) {
	namespace := c.Param("namespace")

	// This would integrate with existing optimization capabilities
	c.JSON(http.StatusOK, gin.H{
		"namespace": namespace,
		"optimizations": []string{
			"Optimization 1: Add resource limits",
			"Optimization 2: Right-size CPU requests",
		},
		"estimatedSavings": 45.50,
	})
}

func multiclusterContextsHandler(c *gin.Context) {
	manager := multicluster.NewClusterManager()
	err := manager.LoadContexts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contexts := manager.ListContexts()
	currentContext, err := manager.GetCurrentContext()
	currentContextName := ""
	if err == nil {
		currentContextName = currentContext.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"contexts":       contexts,
		"currentContext": currentContextName,
		"total":          len(contexts),
	})
}

func multiclusterCompareHandler(c *gin.Context) {
	resourceType := c.Param("resourceType")

	manager := multicluster.NewClusterManager()
	err := manager.LoadContexts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comparison, err := manager.CompareClusters(resourceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": resourceType,
		"comparison":   comparison.GenerateReport(),
		"differences":  len(comparison.Differences),
	})
}

func multiclusterFederatedHandler(c *gin.Context) {
	manager := multicluster.NewClusterManager()
	err := manager.LoadContexts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	report, err := manager.FederatedAnalysis()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": report.GenerateFederatedReport(),
		"summary": gin.H{
			"totalClusters":   report.Summary.TotalClusters,
			"healthyClusters": report.Summary.HealthyClusters,
			"overallHealth":   report.Summary.OverallHealth,
		},
	})
}
