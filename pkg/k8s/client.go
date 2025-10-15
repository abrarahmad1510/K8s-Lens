package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client Manages Kubernetes API Connections and embeds kubernetes.Interface
type Client struct {
	kubernetes.Interface
	Config *rest.Config
}

// NewClient Creates A New Kubernetes Client
func NewClient() (*Client, error) {
	var kubeconfig string

	// Try To Find Kubeconfig In Home Directory
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return nil, fmt.Errorf("Unable To Find Home Directory For Kubeconfig")
	}

	// Build Config From Kubeconfig File
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// Fall Back To In-Cluster Config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("Failed To Get Kubernetes Config: %v", err)
		}
	}

	// Create Clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Failed To Create Kubernetes Client: %v", err)
	}

	return &Client{
		Interface: clientset,
		Config:    config,
	}, nil
}

// TestConnection Verifies Kubernetes API Connectivity
func (c *Client) TestConnection() error {
	_, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed To Connect To Kubernetes API: %v", err)
	}
	return nil
}

// GetServerVersion Returns Kubernetes Server Version
func (c *Client) GetServerVersion() (string, error) {
	version, err := c.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("Failed To Get Server Version: %v", err)
	}
	return version.String(), nil
}
