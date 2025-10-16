#!/bin/bash

set -e

echo "🔧 WEEK 7-8: Testing Automation & Self-healing Features"
echo "========================================================"

# Build the project
echo ""
echo "✔ Building K8s Lens..."
make build
echo "Build completed"

# Test Automation Commands
echo ""
echo "■ TEST 1: Automation Help"
echo "=========================="
./bin/k8s-lens automation --help

echo ""
echo "■ TEST 2: Remediation Help"
echo "==========================="
./bin/k8s-lens automation remediate --help

echo ""
echo "■ TEST 3: List Remediation Actions"
echo "==================================="
./bin/k8s-lens automation remediate list-actions

echo ""
echo "■ TEST 4: Validate Command Structure"
echo "===================================="
./bin/k8s-lens automation remediate pod --help

echo ""
echo "WEEK 7-8 TESTING COMPLETE!"
echo "=========================="
echo "✅ Automation engine built successfully"
echo "✅ CLI commands available"
echo "✅ Remediation actions listed"
echo "✅ Ready for real-world testing with Kubernetes"
