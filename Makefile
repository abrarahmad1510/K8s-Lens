# K8s Lens Build System
BINARY_NAME=k8s-lens
VERSION=0.1.0

.PHONY: build install test clean release help

build:
	@echo "BUILD: Compiling K8s Lens version $(VERSION)"
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) cmd/k8s-lens/main.go

install:
	@echo "INSTALL: Installing K8s Lens to GOPATH"
	go install github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens

test:
	@echo "TEST: Running basic compilation test"
	@./bin/k8s-lens version > /dev/null && echo "TEST PASS: Binary works" || echo "TEST FAIL: Binary broken"

clean:
	@echo "CLEAN: Removing build artifacts"
	rm -rf bin/

release: build
	@echo "RELEASE: Building multi-platform binaries"
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 cmd/k8s-lens/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 cmd/k8s-lens/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 cmd/k8s-lens/main.go

# Monitoring targets
setup-monitoring:
	@echo "Setting up monitoring stack..."
	@chmod +x scripts/setup-prometheus.sh
	@./scripts/setup-prometheus.sh

stop-monitoring:
	@echo "Stopping monitoring..."
	@chmod +x scripts/stop-monitoring.sh
	@./scripts/stop-monitoring.sh

test-integrations:
	@echo "Testing integrations..."
	@chmod +x scripts/test-integrations.sh
	@./scripts/test-integrations.sh

# Quick test commands
test-cluster-metrics:
	@./bin/k8s-lens integrations metrics cluster --prometheus-url http://localhost:9090

test-node-metrics:
	@./bin/k8s-lens integrations metrics node $(shell kubectl get nodes -o jsonpath='{.items[0].metadata.name}') --prometheus-url http://localhost:9090

test-pod-metrics:
	@./bin/k8s-lens integrations metrics pod $(shell kubectl get pods -o jsonpath='{.items[0].metadata.name}') -n default --prometheus-url http://localhost:9090

# Week 5-6: Enterprise Features
test-week5-6:
	@echo "Testing Week 5-6: Enterprise Features"
	@chmod +x scripts/test-week5-6.sh
	@./scripts/test-week5-6.sh

# Individual enterprise tests
test-rbac:
	@echo "Testing RBAC analysis..."
	@./bin/k8s-lens enterprise rbac default

test-security:
	@echo "Testing security scanning..."
	@./bin/k8s-lens enterprise security default

# Create test namespace
create-test-ns:
	@kubectl create namespace k8s-lens-test --dry-run=client -o yaml | kubectl apply -f -

# Clean test namespace
clean-test-ns:
	@kubectl delete namespace k8s-lens-test --ignore-not-found=true
	
help:
	@echo "K8s Lens Build System"
	@echo ""
	@echo "Available commands:"
	@echo "  make build    - Build the binary"
	@echo "  make install  - Install to GOPATH"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make release  - Build release binaries"
	@echo "  make help     - Show this help"
