package optimization

// CostCalculator provides cost estimation functionality
type CostCalculator struct {
	clusterCostPerCPUHour      float64
	clusterCostPerMemoryGBHour float64
}

// NewCostCalculator creates a new CostCalculator
func NewCostCalculator(cpuCost, memoryCost float64) *CostCalculator {
	return &CostCalculator{
		clusterCostPerCPUHour:      cpuCost,
		clusterCostPerMemoryGBHour: memoryCost,
	}
}

// CalculatePodCost estimates monthly cost for a pod
func (c *CostCalculator) CalculatePodCost(cpuRequest, memoryRequest string) (float64, error) {
	cpuCost, err := c.calculateCPUCost(cpuRequest)
	if err != nil {
		return 0, err
	}

	memoryCost, err := c.calculateMemoryCost(memoryRequest)
	if err != nil {
		return 0, err
	}

	// Monthly cost (730 hours in a month)
	monthlyCost := (cpuCost + memoryCost) * 730
	return monthlyCost, nil
}

func (c *CostCalculator) calculateCPUCost(cpuRequest string) (float64, error) {
	// Parse CPU request and calculate hourly cost
	// This is a simplified implementation
	if cpuRequest == "" {
		return 0, nil
	}

	// Default assumption: 1 CPU core
	return c.clusterCostPerCPUHour, nil
}

func (c *CostCalculator) calculateMemoryCost(memoryRequest string) (float64, error) {
	// Parse memory request and calculate hourly cost
	// This is a simplified implementation
	if memoryRequest == "" {
		return 0, nil
	}

	// Default assumption: 1 GB
	return c.clusterCostPerMemoryGBHour, nil
}

// CalculateNamespaceCost estimates total cost for a namespace
func (c *CostCalculator) CalculateNamespaceCost(resources []PodResources) (float64, error) {
	totalCost := 0.0

	for _, resource := range resources {
		cost, err := c.CalculatePodCost(resource.CPU, resource.Memory)
		if err != nil {
			return 0, err
		}
		totalCost += cost
	}

	return totalCost, nil
}

// PodResources represents pod resource requirements
type PodResources struct {
	PodName string
	CPU     string
	Memory  string
}
