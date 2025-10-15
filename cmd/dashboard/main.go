package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", healthHandler)
		api.GET("/cluster/info", clusterInfoHandler)
		api.GET("/analysis/:resourceType/:resourceName", analysisHandler)
		api.GET("/optimization/:namespace", optimizationHandler)
		api.GET("/multicluster/contexts", multiclusterContextsHandler)
		api.GET("/multicluster/compare/:resourceType", multiclusterCompareHandler)
		api.GET("/multicluster/federated", multiclusterFederatedHandler)
	}

	// Web routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "K8s Lens Dashboard",
		})
	})

	router.GET("/cluster", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cluster.html", gin.H{
			"title": "Cluster Overview - K8s Lens",
		})
	})

	router.GET("/analysis", func(c *gin.Context) {
		c.HTML(http.StatusOK, "analysis.html", gin.H{
			"title": "Resource Analysis - K8s Lens",
		})
	})

	router.GET("/multicluster", func(c *gin.Context) {
		c.HTML(http.StatusOK, "multicluster.html", gin.H{
			"title": "Multi-Cluster - K8s Lens",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("K8s Lens Dashboard starting on port %s", port)
	log.Fatal(router.Run(":" + port))
}
