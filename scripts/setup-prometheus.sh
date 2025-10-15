#!/bin/bash

set -e

echo "🚀 Setting up Prometheus monitoring stack..."

# Check if we're using Minikube
if command -v minikube &> /dev/null; then
    echo "📦 Minikube detected - enabling metrics server..."
    minikube addons enable metrics-server
fi

# Create monitoring namespace
echo "📁 Creating monitoring namespace..."
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -

# Install Prometheus using Helm
echo "📊 Installing Prometheus stack..."
brew install helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Check if prometheus stack is already installed
if helm list -n monitoring | grep -q prometheus; then
    echo "✅ Prometheus stack already installed"
else
    echo "📥 Installing Prometheus stack..."
    helm install prometheus prometheus-community/kube-prometheus-stack \
        --namespace monitoring \
        --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false \
        --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
        --set grafana.enabled=true \
        --set alertmanager.enabled=false
fi

# Wait for Prometheus to be ready
echo "⏳ Waiting for Prometheus to be ready..."
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=prometheus -n monitoring --timeout=300s

# Wait for Grafana to be ready (optional)
echo "⏳ Waiting for Grafana to be ready..."
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=grafana -n monitoring --timeout=300s

# Set up port forwarding in background
echo "🔗 Setting up port forwarding..."
pkill -f "kubectl port-forward" || true

# Port forward Prometheus
kubectl port-forward -n monitoring service/prometheus-operated 9090:9090 &
PROMETHEUS_PID=$!

# Port forward Grafana (optional)
kubectl port-forward -n monitoring service/prometheus-grafana 8080:80 &
GRAFANA_PID=$!

# Save PIDs to file for later cleanup
echo $PROMETHEUS_PID > /tmp/prometheus_portforward.pid
echo $GRAFANA_PID > /tmp/grafana_portforward.pid

echo "✅ Prometheus setup complete!"
echo ""
echo "🌐 Access URLs:"
echo "   Prometheus: http://localhost:9090"
echo "   Grafana:    http://localhost:8080 (admin/prom-operator)"
echo ""
echo "🔧 K8s Lens commands:"
echo "   ./bin/k8s-lens integrations metrics pod <pod-name> -n <namespace>"
echo "   ./bin/k8s-lens integrations metrics node <node-name>"
echo "   ./bin/k8s-lens integrations metrics cluster"
echo ""
echo "⏹️  To stop port forwarding, run: pkill -f 'kubectl port-forward'"