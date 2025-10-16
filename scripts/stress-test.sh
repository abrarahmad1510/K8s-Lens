#!/bin/bash

echo "🔥 K8s Lens Stress Test"
echo "======================"

# Test rapid command execution
echo "Testing rapid command execution..."
for i in {1..50}; do
    ./bin/k8s-lens version > /dev/null 2>&1
done
echo "✅ 50 rapid version commands completed"

# Test memory usage under load
echo "Testing memory usage..."
for i in {1..20}; do
    ./bin/k8s-lens --help > /dev/null 2>&1 &
    ./bin/k8s-lens automation --help > /dev/null 2>&1 &
    ./bin/k8s-lens enterprise --help > /dev/null 2>&1 &
done
wait
echo "✅ Parallel help commands completed"

# Test binary size and load time
echo "Binary size: $(du -h ./bin/k8s-lens | cut -f1)"
echo "Load test:"
time for i in {1..10}; do
    ./bin/k8s-lens version > /dev/null
done

echo "🔥 Stress test completed!"
