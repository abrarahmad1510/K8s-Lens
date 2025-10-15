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
