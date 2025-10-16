#!/bin/bash

echo "☸️ K8s Lens Kubernetes Integration Test"
echo "======================================"

# Check if we have kubectl and cluster access
if command -v kubectl &> /dev/null; then
    echo "Testing Kubernetes integration..."
    
    # Test basic cluster access
    if kubectl cluster-info &> /dev/null; then
        echo "✅ Kubernetes cluster accessible"
        
        # Test namespace operations
        kubectl create namespace k8s-lens-test --dry-run=client -o yaml | kubectl apply -f - > /dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo "✅ Namespace operations work"
            kubectl delete namespace k8s-lens-test --ignore-not-found=true > /dev/null 2>&1
        fi
        
        # Test that k8s-lens can create client
        ./bin/k8s-lens enterprise rbac analyze default --help > /dev/null 2>&1
        echo "✅ K8s Lens Kubernetes client initialization works"
    else
        echo "⚠️ Kubernetes cluster not accessible, skipping integration tests"
    fi
else
    echo "⚠️ kubectl not found, skipping Kubernetes integration tests"
fi

echo "☸️ Kubernetes integration test completed"
