package integrations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/abrarahmad1510/k8s-lens/internal/utils"
)

// PrometheusClient represents a client to interact with Prometheus
type PrometheusClient struct {
	baseURL string
	client  *http.Client
}

// NewPrometheusClient creates a new Prometheus client
func NewPrometheusClient(baseURL string) *PrometheusClient {
	return &PrometheusClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// TestConnection tests if Prometheus is accessible
func (p *PrometheusClient) TestConnection() error {
	u, err := url.Parse(p.baseURL + "/api/v1/query")
	if err != nil {
		return fmt.Errorf("invalid Prometheus URL: %v", err)
	}

	q := u.Query()
	q.Set("query", "up")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot connect to Prometheus at %s: %v", p.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Prometheus returned status %d", resp.StatusCode)
	}

	return nil
}

// QueryResult represents the result of a Prometheus query
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// PodMetrics contains metrics for a pod
type PodMetrics struct {
	PodName     string
	Namespace   string
	CPUUsage    float64
	MemoryUsage float64
	NetworkRx   float64
	NetworkTx   float64
	Timestamp   time.Time
	Error       string
}

// NodeMetrics contains metrics for a node
type NodeMetrics struct {
	NodeName    string
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
	PodCount    int
	Timestamp   time.Time
	Error       string
}

// ClusterMetrics contains cluster-level metrics
type ClusterMetrics struct {
	TotalNodes     int
	TotalPods      int
	CPUCapacity    float64
	CPUUsage       float64
	MemoryCapacity float64
	MemoryUsage    float64
	Timestamp      time.Time
	Error          string
}

// GetPodMetrics retrieves metrics for a specific pod
func (p *PrometheusClient) GetPodMetrics(podName, namespace string) (*PodMetrics, error) {
	utils.PrintInfo("Fetching metrics for pod %s in namespace %s", podName, namespace)

	metrics := &PodMetrics{
		PodName:   podName,
		Namespace: namespace,
		Timestamp: time.Now(),
	}

	// Test connection first
	if err := p.TestConnection(); err != nil {
		metrics.Error = fmt.Sprintf("Prometheus connection failed: %v", err)
		return metrics, fmt.Errorf("Prometheus connection failed: %v", err)
	}

	// Query for CPU usage
	cpuQuery := fmt.Sprintf(`rate(container_cpu_usage_seconds_total{pod="%s", namespace="%s"}[5m])`, podName, namespace)
	cpuValue, err := p.queryPrometheus(cpuQuery)
	if err != nil {
		utils.PrintWarning("Failed to query CPU metrics: %v", err)
		metrics.Error = fmt.Sprintf("CPU metrics unavailable: %v", err)
	} else if len(cpuValue) > 0 {
		metrics.CPUUsage = cpuValue[0]
	}

	// Query for memory usage
	memoryQuery := fmt.Sprintf(`container_memory_usage_bytes{pod="%s", namespace="%s"}`, podName, namespace)
	memoryValue, err := p.queryPrometheus(memoryQuery)
	if err != nil {
		utils.PrintWarning("Failed to query memory metrics: %v", err)
		if metrics.Error != "" {
			metrics.Error += "; " + fmt.Sprintf("Memory metrics unavailable: %v", err)
		} else {
			metrics.Error = fmt.Sprintf("Memory metrics unavailable: %v", err)
		}
	} else if len(memoryValue) > 0 {
		metrics.MemoryUsage = memoryValue[0]
	}

	// Only proceed with network metrics if we have basic connectivity
	if metrics.Error == "" {
		// Query for network receive
		networkRxQuery := fmt.Sprintf(`rate(container_network_receive_bytes_total{pod="%s", namespace="%s"}[5m])`, podName, namespace)
		networkRxValue, err := p.queryPrometheus(networkRxQuery)
		if err != nil {
			utils.PrintWarning("Failed to query network RX metrics: %v", err)
		} else if len(networkRxValue) > 0 {
			metrics.NetworkRx = networkRxValue[0]
		}

		// Query for network transmit
		networkTxQuery := fmt.Sprintf(`rate(container_network_transmit_bytes_total{pod="%s", namespace="%s"}[5m])`, podName, namespace)
		networkTxValue, err := p.queryPrometheus(networkTxQuery)
		if err != nil {
			utils.PrintWarning("Failed to query network TX metrics: %v", err)
		} else if len(networkTxValue) > 0 {
			metrics.NetworkTx = networkTxValue[0]
		}
	}

	return metrics, nil
}

// GetNodeMetrics retrieves metrics for a specific node
func (p *PrometheusClient) GetNodeMetrics(nodeName string) (*NodeMetrics, error) {
	utils.PrintInfo("Fetching metrics for node %s", nodeName)

	metrics := &NodeMetrics{
		NodeName:  nodeName,
		Timestamp: time.Now(),
	}

	// Test connection first
	if err := p.TestConnection(); err != nil {
		metrics.Error = fmt.Sprintf("Prometheus connection failed: %v", err)
		return metrics, fmt.Errorf("Prometheus connection failed: %v", err)
	}

	// Query for node CPU usage - fixed query for kube-prometheus-stack
	cpuQuery := `100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
	cpuValue, err := p.queryPrometheus(cpuQuery)
	if err != nil {
		utils.PrintWarning("Failed to query node CPU metrics: %v", err)
		metrics.Error = fmt.Sprintf("CPU metrics unavailable: %v", err)
	} else if len(cpuValue) > 0 {
		metrics.CPUUsage = cpuValue[0]
	}

	// Query for node memory usage - fixed query for kube-prometheus-stack
	memoryQuery := `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`
	memoryValue, err := p.queryPrometheus(memoryQuery)
	if err != nil {
		utils.PrintWarning("Failed to query node memory metrics: %v", err)
		if metrics.Error != "" {
			metrics.Error += "; " + fmt.Sprintf("Memory metrics unavailable: %v", err)
		} else {
			metrics.Error = fmt.Sprintf("Memory metrics unavailable: %v", err)
		}
	} else if len(memoryValue) > 0 {
		metrics.MemoryUsage = memoryValue[0]
	}

	// Query for node disk usage - fixed query for kube-prometheus-stack
	diskQuery := `(1 - (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"})) * 100`
	diskValue, err := p.queryPrometheus(diskQuery)
	if err != nil {
		utils.PrintWarning("Failed to query node disk metrics: %v", err)
	} else if len(diskValue) > 0 {
		metrics.DiskUsage = diskValue[0]
	}

	// Query for pod count on node - fixed query for kube-prometheus-stack
	podCountQuery := `count(kube_pod_info) by (node)`
	podCountValue, err := p.queryPrometheus(podCountQuery)
	if err != nil {
		utils.PrintWarning("Failed to query pod count metrics: %v", err)
	} else if len(podCountValue) > 0 {
		// Find the value for our specific node
		podCountForNode, err := p.queryPrometheus(fmt.Sprintf(`count(kube_pod_info{node="%s"})`, nodeName))
		if err == nil && len(podCountForNode) > 0 {
			metrics.PodCount = int(podCountForNode[0])
		}
	}

	return metrics, nil
}

// GetClusterMetrics retrieves cluster-level metrics
func (p *PrometheusClient) GetClusterMetrics() (*ClusterMetrics, error) {
	utils.PrintInfo("Fetching cluster-level metrics")

	metrics := &ClusterMetrics{
		Timestamp: time.Now(),
	}

	// Test connection first
	if err := p.TestConnection(); err != nil {
		metrics.Error = fmt.Sprintf("Prometheus connection failed: %v", err)
		return metrics, fmt.Errorf("Prometheus connection failed: %v", err)
	}

	// Query for total nodes - fixed query
	nodeCountQuery := `count(kube_node_info)`
	nodeCountValue, err := p.queryPrometheus(nodeCountQuery)
	if err != nil {
		utils.PrintWarning("Failed to query node count: %v", err)
		metrics.Error = fmt.Sprintf("Node count unavailable: %v", err)
	} else if len(nodeCountValue) > 0 {
		metrics.TotalNodes = int(nodeCountValue[0])
	}

	// Query for total pods - fixed query
	podCountQuery := `count(kube_pod_info)`
	podCountValue, err := p.queryPrometheus(podCountQuery)
	if err != nil {
		utils.PrintWarning("Failed to query pod count: %v", err)
		if metrics.Error != "" {
			metrics.Error += "; " + fmt.Sprintf("Pod count unavailable: %v", err)
		} else {
			metrics.Error = fmt.Sprintf("Pod count unavailable: %v", err)
		}
	} else if len(podCountValue) > 0 {
		metrics.TotalPods = int(podCountValue[0])
	}

	// Query for cluster CPU capacity - FIXED with multiple fallbacks
	var cpuCapacity float64
	cpuQueries := []string{
		`sum(kube_node_status_capacity_cpu_cores)`,
		`sum(machine_cpu_cores)`,
		`sum(kube_node_status_allocatable_cpu_cores)`,
	}

	for _, query := range cpuQueries {
		capacityValue, err := p.queryPrometheus(query)
		if err == nil && len(capacityValue) > 0 && capacityValue[0] > 0 {
			cpuCapacity = capacityValue[0]
			break
		}
	}
	metrics.CPUCapacity = cpuCapacity

	// Query for cluster CPU usage
	cpuUsageQuery := `sum(rate(container_cpu_usage_seconds_total[5m]))`
	cpuUsageValue, err := p.queryPrometheus(cpuUsageQuery)
	if err != nil {
		utils.PrintWarning("Failed to query CPU usage: %v", err)
	} else if len(cpuUsageValue) > 0 {
		metrics.CPUUsage = cpuUsageValue[0]
	}

	// Query for cluster memory capacity - FIXED with multiple fallbacks
	var memoryCapacity float64
	memoryQueries := []string{
		`sum(kube_node_status_capacity_memory_bytes) / (1024 * 1024 * 1024)`,
		`sum(machine_memory_bytes) / (1024 * 1024 * 1024)`,
		`sum(kube_node_status_allocatable_memory_bytes) / (1024 * 1024 * 1024)`,
	}

	for _, query := range memoryQueries {
		capacityValue, err := p.queryPrometheus(query)
		if err == nil && len(capacityValue) > 0 && capacityValue[0] > 0 {
			memoryCapacity = capacityValue[0]
			break
		}
	}
	metrics.MemoryCapacity = memoryCapacity

	// Query for cluster memory usage
	memoryUsageQuery := `sum(container_memory_working_set_bytes) / (1024 * 1024 * 1024)`
	memoryUsageValue, err := p.queryPrometheus(memoryUsageQuery)
	if err != nil {
		utils.PrintWarning("Failed to query memory usage: %v", err)
	} else if len(memoryUsageValue) > 0 {
		metrics.MemoryUsage = memoryUsageValue[0]
	}

	return metrics, nil
}

// queryPrometheus executes a Prometheus query and returns the values
func (p *PrometheusClient) queryPrometheus(query string) ([]float64, error) {
	u, err := url.Parse(p.baseURL + "/api/v1/query")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Prometheus returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result QueryResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("Prometheus query failed: %s", string(body))
	}

	var values []float64
	for _, res := range result.Data.Result {
		if len(res.Value) >= 2 {
			if str, ok := res.Value[1].(string); ok {
				if f, err := strconv.ParseFloat(str, 64); err == nil {
					values = append(values, f)
				}
			}
		}
	}

	return values, nil
}
