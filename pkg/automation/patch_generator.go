package automation

import "fmt"

// PatchGenerator provides YAML patch generation for automated fixes
type PatchGenerator struct {
}

// NewPatchGenerator creates a new PatchGenerator
func NewPatchGenerator() *PatchGenerator {
	return &PatchGenerator{}
}

// GenerateResourceLimitsPatch generates patch for adding resource limits
func GenerateResourceLimitsPatch() string {
	return `
spec:
  template:
    spec:
      containers:
      - name: "*"
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
`
}

// GenerateProbePatch generates patch for adding health check probes
func GenerateProbePatch() string {
	return `
spec:
  template:
    spec:
      containers:
      - name: "*"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
`
}

// GenerateSecurityContextPatch generates patch for adding security context
func GenerateSecurityContextPatch() string {
	return `
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        seccompProfile:
          type: RuntimeDefault
      containers:
      - name: "*"
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
`
}

// GenerateResourceRightSizingPatch generates patch for resource optimization
func GenerateResourceRightSizingPatch(cpu, memory string) string {
	return fmt.Sprintf(`
spec:
  template:
    spec:
      containers:
      - name: "*"
        resources:
          requests:
            cpu: "%s"
            memory: "%s"
          limits:
            cpu: "%s"
            memory: "%s"
`, cpu, memory, cpu, memory)
}
