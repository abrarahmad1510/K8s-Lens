package ai

// RecommendationEngine provides intelligent recommendations
type RecommendationEngine struct {
	knowledgeBase map[string]RecommendationRule
}

// NewRecommendationEngine creates a new RecommendationEngine
func NewRecommendationEngine() *RecommendationEngine {
	engine := &RecommendationEngine{
		knowledgeBase: make(map[string]RecommendationRule),
	}
	engine.initializeKnowledgeBase()
	return engine
}

// RecommendationRule defines a pattern and corresponding recommendation
type RecommendationRule struct {
	Pattern        string
	Condition      func(context map[string]interface{}) bool
	Recommendation string
	Priority       int
	Category       string
}

// GenerateRecommendations generates intelligent recommendations based on context
func (r *RecommendationEngine) GenerateRecommendations(context map[string]interface{}) []string {
	var recommendations []string

	for _, rule := range r.knowledgeBase {
		if rule.Condition(context) {
			recommendations = append(recommendations, rule.Recommendation)
		}
	}

	return recommendations
}

func (r *RecommendationEngine) initializeKnowledgeBase() {
	// Define recommendation rules based on common Kubernetes issues
	r.knowledgeBase["high_restarts"] = RecommendationRule{
		Pattern: "High container restart count",
		Condition: func(context map[string]interface{}) bool {
			restarts, ok := context["restart_count"].(int)
			return ok && restarts > 10
		},
		Recommendation: "Investigate application crashes. Check application logs and consider adding liveness probes.",
		Priority:       1,
		Category:       "Reliability",
	}

	r.knowledgeBase["missing_limits"] = RecommendationRule{
		Pattern: "Missing resource limits",
		Condition: func(context map[string]interface{}) bool {
			hasLimits, ok := context["has_limits"].(bool)
			return ok && !hasLimits
		},
		Recommendation: "Add resource limits to prevent resource exhaustion and ensure quality of service.",
		Priority:       2,
		Category:       "Performance",
	}

	r.knowledgeBase["image_pull_backoff"] = RecommendationRule{
		Pattern: "Image pull failures",
		Condition: func(context map[string]interface{}) bool {
			events, ok := context["events"].([]string)
			if !ok {
				return false
			}
			for _, event := range events {
				if event == "ImagePullBackOff" || event == "ErrImagePull" {
					return true
				}
			}
			return false
		},
		Recommendation: "Check image repository accessibility and image pull secrets configuration.",
		Priority:       1,
		Category:       "Configuration",
	}

	r.knowledgeBase["pending_pods"] = RecommendationRule{
		Pattern: "Pods stuck in pending state",
		Condition: func(context map[string]interface{}) bool {
			pendingPods, ok := context["pending_pods"].(int)
			return ok && pendingPods > 0
		},
		Recommendation: "Check node resources, affinity rules, and persistent volume claims.",
		Priority:       2,
		Category:       "Resource Management",
	}

	r.knowledgeBase["low_cpu_usage"] = RecommendationRule{
		Pattern: "Low CPU utilization",
		Condition: func(context map[string]interface{}) bool {
			cpuUsage, ok := context["cpu_usage_percent"].(float64)
			return ok && cpuUsage < 20.0
		},
		Recommendation: "Consider reducing CPU requests to optimize resource allocation and reduce costs.",
		Priority:       3,
		Category:       "Cost Optimization",
	}

	r.knowledgeBase["low_memory_usage"] = RecommendationRule{
		Pattern: "Low memory utilization",
		Condition: func(context map[string]interface{}) bool {
			memoryUsage, ok := context["memory_usage_percent"].(float64)
			return ok && memoryUsage < 30.0
		},
		Recommendation: "Consider reducing memory requests to optimize resource allocation.",
		Priority:       3,
		Category:       "Cost Optimization",
	}
}
