package automation

import (
	"fmt"
)

// FixEngine provides automated fix generation
type FixEngine struct {
}

// NewFixEngine creates a new FixEngine
func NewFixEngine() *FixEngine {
	return &FixEngine{}
}

// FixPlan contains the automated fix plan
type FixPlan struct {
	ResourceType string
	ResourceName string
	Namespace    string
	Fixes        []Fix
	Risks        []string
	Confidence   int
	Preview      string
}

// Fix represents a single automated fix
type Fix struct {
	Type        string
	Description string
	Action      string
	YAMLPatch   string
	RiskLevel   string
	BackupPlan  string
}

// GenerateFix generates automated fixes for identified issues
func (f *FixEngine) GenerateFix(resourceType, resourceName, namespace string, issues []string) (*FixPlan, error) {
	plan := &FixPlan{
		ResourceType: resourceType,
		ResourceName: resourceName,
		Namespace:    namespace,
		Confidence:   75,
	}

	// Generate fixes based on issue types
	for _, issue := range issues {
		fix := f.generateFixForIssue(issue, resourceType, resourceName, namespace)
		if fix != nil {
			plan.Fixes = append(plan.Fixes, *fix)
		}
	}

	// Generate preview
	plan.Preview = f.generatePreview(plan)

	// Assess risks
	plan.Risks = f.assessRisks(plan)

	return plan, nil
}

func (f *FixEngine) generateFixForIssue(issue, resourceType, resourceName, namespace string) *Fix {
	switch issue {
	case "Missing resource limits":
		return &Fix{
			Type:        "Resource Configuration",
			Description: "Add default resource limits to prevent resource exhaustion",
			Action:      "Patch resource configuration",
			YAMLPatch:   GenerateResourceLimitsPatch(), // Fixed: Use exported function
			RiskLevel:   "Low",
			BackupPlan:  "Rollback to previous resource configuration",
		}
	case "High restart count":
		return &Fix{
			Type:        "Probe Configuration",
			Description: "Add liveness and readiness probes to improve container health checks",
			Action:      "Add health check probes",
			YAMLPatch:   GenerateProbePatch(), // Fixed: Use exported function
			RiskLevel:   "Medium",
			BackupPlan:  "Remove probes and restart containers",
		}
	case "Security context missing":
		return &Fix{
			Type:        "Security Hardening",
			Description: "Add security context to run as non-root user",
			Action:      "Update security context",
			YAMLPatch:   GenerateSecurityContextPatch(), // Fixed: Use exported function
			RiskLevel:   "Low",
			BackupPlan:  "Revert security context changes",
		}
	}

	return nil
}

func (f *FixEngine) generatePreview(plan *FixPlan) string {
	preview := fmt.Sprintf("Fix Plan for %s/%s in namespace %s\n", plan.ResourceType, plan.ResourceName, plan.Namespace)
	preview += "=================================================================\n\n"

	for i, fix := range plan.Fixes {
		preview += fmt.Sprintf("Fix %d: %s\n", i+1, fix.Type)
		preview += fmt.Sprintf("Description: %s\n", fix.Description)
		preview += fmt.Sprintf("Action: %s\n", fix.Action)
		preview += fmt.Sprintf("Risk Level: %s\n", fix.RiskLevel)
		preview += fmt.Sprintf("YAML Patch:\n%s\n", fix.YAMLPatch)
		preview += fmt.Sprintf("Backup Plan: %s\n", fix.BackupPlan)
		preview += "---\n"
	}

	preview += fmt.Sprintf("Overall Confidence: %d%%\n", plan.Confidence)
	preview += "Risks to Consider:\n"
	for _, risk := range plan.Risks {
		preview += fmt.Sprintf("- %s\n", risk)
	}

	return preview
}

func (f *FixEngine) assessRisks(plan *FixPlan) []string {
	var risks []string

	for _, fix := range plan.Fixes {
		switch fix.RiskLevel {
		case "High":
			risks = append(risks, fmt.Sprintf("High risk fix: %s - %s", fix.Type, fix.Description))
		case "Medium":
			risks = append(risks, fmt.Sprintf("Medium risk fix: %s - may cause temporary downtime", fix.Type))
		}
	}

	if len(plan.Fixes) > 3 {
		risks = append(risks, "Multiple changes being applied at once - consider applying fixes incrementally")
	}

	return risks
}
