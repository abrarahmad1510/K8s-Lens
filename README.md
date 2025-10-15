cat > README.md << 'EOF'
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

### Installation

```bash
# Clone the repository
git clone https://github.com/abrarahmad1510/k8s-lens
cd k8s-lens

# Build from source
make build

# Install to your PATH (optional)
sudo cp bin/k8s-lens /usr/local/bin/
