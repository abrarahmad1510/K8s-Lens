#!/bin/bash

echo "🧪 K8s Lens - Test Suite"
echo "========================"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

test_command() {
    local name="$1"
    local command="$2"
    
    echo -n "Testing: $name... "
    
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ FAIL${NC}"
        ((FAILED++))
        return 1
    fi
}

echo ""
echo "📦 Build System Tests"
echo "====================="

test_command "Project builds" "make build"
test_command "Binary is executable" "test -x ./bin/k8s-lens"
test_command "Version command works" "./bin/k8s-lens version"

echo ""
echo "🚀 Core CLI Tests"
echo "================="

test_command "Root help" "./bin/k8s-lens --help"
test_command "Analyze commands" "./bin/k8s-lens analyze --help"
test_command "Enterprise commands" "./bin/k8s-lens enterprise --help"
test_command "Automation commands" "./bin/k8s-lens automation --help"

echo ""
echo "🔧 Feature Tests"
echo "================"

test_command "RBAC analysis help" "./bin/k8s-lens enterprise rbac --help"
test_command "Security scan help" "./bin/k8s-lens enterprise security --help"
test_command "Remediation actions" "./bin/k8s-lens automation remediate list-actions"
test_command "Integrations help" "./bin/k8s-lens integrations --help"

echo ""
echo "⚡ Makefile Targets"
echo "==================="

test_command "Test Week 3-4" "make test-week3-4"
test_command "Test Week 5-6" "make test-week5-6"
test_command "Test Week 7-8" "make test-week7-8"
test_command "Phase 4 complete" "make test-phase4-complete"

echo ""
echo "📊 Test Summary"
echo "==============="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed! K8s Lens is ready! 🚀${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}❌ $FAILED test(s) failed${NC}"
    exit 1
fi
