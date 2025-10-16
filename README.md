# K8s Lens - AI-Powered Kubernetes Troubleshooting & Automation

## Table of Contents
1. [About](#about)
2. [Features](#features)
3. [Quick Start](#quick-start)
4. [Installation](#installation)
5. [Usage](#usage)
6. [Enterprise Features](#enterprise-features)
7. [Automation & Self-Healing](#automation--self-healing)
8. [Contributing](#contributing)
9. [Technical Stack](#technical-stack)
10. [Development](#development)
11. [License](#license)

## About

K8s Lens is an advanced Kubernetes CLI tool that provides intelligent diagnostics, automated remediation, and comprehensive monitoring capabilities. Built with AI-powered analysis and enterprise-grade security features, it helps DevOps teams and SREs efficiently manage and troubleshoot Kubernetes clusters at scale.

The project combines machine learning insights with practical automation to reduce manual toil, prevent outages, and optimize resource utilization across multi-cluster environments.

## Features

### 🔍 Intelligent Analysis
- **AI-Powered Diagnostics**: Machine learning-driven analysis of Kubernetes resources
- **Comprehensive Health Reports**: Actionable insights with SRE best practices
- **Deep Resource Inspection**: Pods, deployments, services, and statefulsets
- **Multi-Cluster Support**: Unified view across multiple Kubernetes clusters

### 🛡️ Enterprise Security
- **RBAC Security Analysis**: Risk assessment and permission auditing
- **Security Scanning**: Compliance scoring and vulnerability detection
- **Policy Enforcement**: Automated security policy validation
- **Audit Logging**: Comprehensive operation tracking

### 🤖 Automation & Self-Healing
- **Automated Remediation**: Auto-fix common Kubernetes issues
- **Predictive Scaling**: ML-based resource optimization
- **Self-Healing Mechanisms**: Automatic recovery from failures
- **Intelligent Rollbacks**: Safe deployment management

### 📊 Advanced Monitoring
- **Prometheus Integration**: Real-time metrics collection and analysis
- **Performance Analytics**: Resource utilization and bottleneck detection
- **Custom Dashboards**: Tailored monitoring views
- **Alert Integration**: Smart alerting and notification system

## Quick Start

### Prerequisites
- Kubernetes cluster (v1.20+)
- kubectl configured with cluster access
- Go 1.19+ (for development)

### Basic Installation
```bash
# Download and install
go install github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens@latest

# Or build from source
git clone https://github.com/abrarahmad1510/k8s-lens
cd k8s-lens
make build
