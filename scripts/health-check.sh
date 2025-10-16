#!/bin/bash

echo "🏥 K8s Lens Health Check"
echo "========================"

echo ""
echo "1. Build Status:"
make build && echo "✅ Build successful" || echo "❌ Build failed"

echo ""
echo "2. Basic Commands:"
./bin/k8s-lens version > /dev/null && echo "✅ Version command works" || echo "❌ Version command broken"
./bin/k8s-lens --help > /dev/null && echo "✅ Help command works" || echo "❌ Help command broken"

echo ""
echo "3. Feature Modules:"
./bin/k8s-lens analyze --help > /dev/null && echo "✅ Analysis module OK" || echo "❌ Analysis module broken"
./bin/k8s-lens enterprise --help > /dev/null && echo "✅ Enterprise module OK" || echo "❌ Enterprise module broken"
./bin/k8s-lens automation --help > /dev/null && echo "✅ Automation module OK" || echo "❌ Automation module broken"

echo ""
echo "4. Key Features:"
./bin/k8s-lens automation remediate list-actions > /dev/null && echo "✅ Remediation actions OK" || echo "❌ Remediation actions broken"

echo ""
echo "🏥 Health Check Complete!"
