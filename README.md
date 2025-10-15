# K8s Lens 🔍

> AI-Powered Kubernetes Troubleshooting Assistant - Professional CLI Foundation

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Phase](https://img.shields.io/badge/Phase-1%20Complete-success)](https://github.com/abrarahmad1510/k8s-lens)

## 🎯 Overview

K8s Lens is an intelligent command-line tool designed to help developers and operators troubleshoot Kubernetes issues efficiently. This Phase 1 release establishes the professional CLI foundation for future Kubernetes diagnostic capabilities.

## ✨ Phase 1 Features

- **Professional CLI Framework** - Built with Cobra for enterprise-grade command structure
- **Capital Case Messaging** - Clean, professional output without emojis
- **Multi-Resource Support** - Analyze pods, deployments, services, nodes, and namespaces
- **Shell Completions** - Full support for bash, zsh, fish, and PowerShell
- **Cross-Platform Builds** - Linux, macOS (Intel/Apple Silicon), Windows
- **Verbose Mode** - Detailed analysis output for debugging
- **Namespace Awareness** - Multi-tenant cluster support

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- Git

### Installation: 

```bash
# Clone the repository
git clone https://github.com/abrarahmad1510/k8s-lens
cd k8s-lens

# Build from source
make build

# Install to your PATH (optional)
sudo cp bin/k8s-lens /usr/local/bin/
```
### Basic Usage: 

```bash
# Check version
k8s-lens version

# Analyze a pod
k8s-lens analyze pod my-app-pod

# Analyze with verbose output
k8s-lens analyze deployment web-service --verbose

# Analyze in specific namespace
k8s-lens analyze service frontend -n production

# Generate shell completions
k8s-lens completion zsh > ~/.zsh/completion/_k8s-lens
```
### Example Output: 

```bash
$ k8s-lens analyze pod test-pod --verbose

INFO: Verbose mode enabled
INFO: Resource type: pod
INFO: Resource name: test-pod
INFO: Namespace: default
ANALYZING: pod/test-pod in namespace 'default'
STATUS: K8s Lens analysis engine initialized
NEXT: Kubernetes cluster integration pending

--- SIMULATION RESULTS ---
PASS: Pod spec validation completed
WARNING: Container resource limits not set
FAIL: Liveness probe configuration issue detected
RECOMMENDATION: Check application health endpoint configuration
```
### Project Architecture: 

```bash
k8s-lens/
├── cmd/k8s-lens/
│   └── main.go              # CLI entry point
├── internal/utils/
│   └── helpers.go           # Utility functions
├── pkg/                     # Future packages
├── Makefile                 # Build system
├── go.mod                   # Dependencies
└── README.md               # Documentation
```






