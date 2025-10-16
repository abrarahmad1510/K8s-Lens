#!/bin/bash

echo "🔄 K8s Lens Regression Test Suite"
echo "================================"

PASSED=0
FAILED=0

test_regression() {
    local name="$1"
    local command="$2"
    
    echo -n "Regression: $name... "
    
    if eval "$command" > /dev/null 2>&1; then
        echo "✅ PASS"
        ((PASSED++))
    else
        echo "❌ FAIL"
        ((FAILED++))
    fi
}

echo "Running regression tests..."

# Test 1: Automation commands handle missing arguments gracefully
# Fixed: Check for cobra's error message pattern
test_regression "Automation missing args" "./bin/k8s-lens automation remediate pod 2>&1 | grep -q 'accepts 2 arg(s)'"

# Test 2: Help works for all subcommands
test_regression "Automation list actions" "./bin/k8s-lens automation remediate list-actions > /dev/null"

# Test 3: Build works after clean
test_regression "Clean and rebuild" "make clean && make build && ./bin/k8s-lens version > /dev/null"

echo ""
echo "Regression Summary:"
echo "✅ Passed: $PASSED"
echo "❌ Failed: $FAILED"
