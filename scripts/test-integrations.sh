#!/bin/bash

set -e

echo "🧪 Testing K8s Lens integrations..."

# Build the project
echo "🔨 Building K8s Lens..."
make build

# Check if monitoring is running
if ! curl -s http://localhost:9090/-/healthy > /dev/null; then
    echo "📊 Setting up monitoring..."
    ./scripts/setup-prometheus.sh
    sleep 10
fi

# Deploy a test application
echo "🚀 Deploying test nginx application..."
kubectl create deployment nginx-test --image=nginx:latest --dry-run=client -o yaml | kubectl apply -f -
kubectl scale deployment nginx-test --replicas=2
kubectl wait --for=condition=ready pod -l app=nginx-test --timeout=60s

# Get resource names
POD_NAME=$(kubectl get pods -l app=nginx-test -o jsonpath='{.items[0].metadata.name}')
NODE_NAME=$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')

echo ""
echo "📋 Test Resources:"
echo "   Pod:  $POD_NAME"
echo "   Node: $NODE_NAME"
echo ""

# Test cluster metrics
echo "🌐 Testing cluster metrics..."
echo "======================================"
./bin/k8s-lens integrations metrics cluster --prometheus-url http://localhost:9090

echo ""
echo "🖥️  Testing node metrics..."
echo "======================================"
./bin/k8s-lens integrations metrics node $NODE_NAME --prometheus-url http://localhost:9090

echo ""
echo "📦 Testing pod metrics..."
echo "======================================"
./bin/k8s-lens integrations metrics pod $POD_NAME -n default --prometheus-url http://localhost:9090

# Cleanup
echo ""
echo "🧹 Cleaning up test app..."
kubectl delete deployment nginx-test

echo ""
echo "🎉 All integration tests completed!"
echo ""
echo "💡 Monitoring is still running. To stop it, run: ./scripts/stop-monitoring.sh"