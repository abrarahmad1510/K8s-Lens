#!/bin/bash

echo "🛑 Stopping monitoring port forwarding..."

# Kill port forwarding processes
pkill -f "kubectl port-forward" || true

# Remove PID files
rm -f /tmp/prometheus_portforward.pid
rm -f /tmp/grafana_portforward.pid

echo "✅ Monitoring stopped"