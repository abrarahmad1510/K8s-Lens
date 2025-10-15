package ai

import (
	"strings"
)

// RecommendationEngine Provides Intelligent Troubleshooting Suggestions
type RecommendationEngine struct {
	knowledgeBase map[string][]string
}

// NewRecommendationEngine Creates A New Recommendation Engine
func NewRecommendationEngine() *RecommendationEngine {
	return &RecommendationEngine{
		knowledgeBase: map[string][]string{
			"ImagePullBackOff": {
				"Check If The Image Name And Tag Are Correct",
				"Verify Docker Registry Accessibility",
				"Check Image Pull Secrets: kubectl get secrets",
				"Try Pulling The Image Manually: docker pull <image>",
			},
			"CrashLoopBackOff": {
				"Check Application Logs: kubectl logs <pod> --previous",
				"Verify Environment Variables And Config Maps",
				"Check Resource Limits And Requests",
				"Test The Application Locally With The Same Configuration",
			},
			"Pending": {
				"Check Node Resources: kubectl describe nodes",
				"Verify Persistent Volume Claims",
				"Check Resource Quotas: kubectl describe quota",
				"Look At Scheduler Events: kubectl get events",
			},
			"OOMKilled": {
				"Increase Memory Limits In The Container Spec",
				"Check For Memory Leaks In The Application",
				"Monitor Application Memory Usage",
				"Consider Adding Resource Requests And Limits",
			},
		},
	}
}

// GetRecommendations Returns Troubleshooting Suggestions For Specific Issues
func (r *RecommendationEngine) GetRecommendations(issueType string) []string {
	if recs, exists := r.knowledgeBase[issueType]; exists {
		return recs
	}
	return []string{"Check Kubernetes Documentation And Application Logs For More Details"}
}

// AnalyzePatterns Analyzes Log Patterns For Intelligent Recommendations
func (r *RecommendationEngine) AnalyzePatterns(logs string) []string {
	logsLower := strings.ToLower(logs)
	var recommendations []string

	if strings.Contains(logsLower, "connection refused") {
		recommendations = append(recommendations,
			"Application Cannot Connect To Dependent Service - Check Service Discovery And Networking")
	}

	if strings.Contains(logsLower, "out of memory") || strings.Contains(logsLower, "oom") {
		recommendations = append(recommendations,
			"Container Is Hitting Memory Limits - Consider Increasing Memory Requests And Limits")
	}

	if strings.Contains(logsLower, "permission denied") {
		recommendations = append(recommendations,
			"Check Security Context And File Permissions In The Container")
	}

	if strings.Contains(logsLower, "no such file or directory") {
		recommendations = append(recommendations,
			"Check Container File System And Volume Mounts")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"Review Application Logs For Specific Error Patterns")
	}

	return recommendations
}
