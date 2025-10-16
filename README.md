# Open Source Contribution - Containerisation (K8s Lens)

## Table of Contents
- [About](#about)
- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [Enterprise Features](#enterprise-features)
- [Automation & Self-Healing](#automation--self-healing)
- [Contributing](#contributing)
- [Technical Stack](#technical-stack)
- [Development](#development)
- [License](#license)
    
## About

K8s Lens is an advanced Kubernetes CLI tool that provides intelligent diagnostics, automated remediation, and comprehensive monitoring capabilities. 
Built with AI-powered analysis and enterprise-grade security features, it helps DevOps teams and SREs efficiently manage and troubleshoot Kubernetes clusters at scale.

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

## Basic Installation
```bash
# Download and install
go install github.com/abrarahmad1510/k8s-lens/cmd/k8s-lens@latest

# Or build from source
git clone https://github.com/abrarahmad1510/k8s-lens
cd k8s-lens
make build
```
### Verify Installation 
```bash
k8s-lens version
k8s-lens --help
```
### Binary Download 
```bash
# Linux
curl -L https://github.com/abrarahmad1510/k8s-lens/releases/latest/download/k8s-lens-linux-amd64 -o k8s-lens
chmod +x k8s-lens
sudo mv k8s-lens /usr/local/bin/

# macOS
curl -L https://github.com/abrarahmad1510/k8s-lens/releases/latest/download/k8s-lens-darwin-amd64 -o k8s-lens
chmod +x k8s-lens
sudo mv k8s-lens /usr/local/bin/
```
### Docker
```bash
docker run -v ~/.kube:/root/.kube abrarahmad1510/k8s-lens:latest analyze cluster
```

### Helm Chart
```bash
helm repo add k8s-lens https://abrarahmad1510.github.io/k8s-lens
helm install k8s-lens k8s-lens/k8s-lens
```
## Usage
### Basic Analysis
```bash
# Analyze a specific pod
k8s-lens analyze pod my-app-pod -n production

# Analyze deployment health
k8s-lens analyze deployment my-web-service

# Comprehensive cluster analysis
k8s-lens analyze cluster
```

### Enterprise Security 
```bash
# RBAC risk analysis
k8s-lens enterprise rbac analyze default

# Security scanning
k8s-lens enterprise security scan production
```
### Automation & Self-Healing
```bash
# Automatically remediate pod issues
k8s-lens automation remediate pod my-pod CrashLoopBackOff -n default

# List available remediation actions
k8s-lens automation remediate list-actions

# Predictive scaling
k8s-lens automation scale predictive my-deployment
```

### Advanced Analysis
```bash
# Prometheus metrics integration
k8s-lens integrations metrics cluster --prometheus-url http://localhost:9090

# Multi-cluster operations
k8s-lens multicluster status
```
## Enterprise Features 
### 🔐 RBAC Security Analysis
- **Risk Assessment**: Identify dangerous permissions and service accounts
- **Compliance Scoring**: Measure against security benchmarks
- **Recommendation Engine**: Automated security improvements
- **Audit Reports**:  Comprehensive permission documentation

### 🚨 Security Scanning
- **Vulnerability Detection**: CVEs and security misconfigurations
- **Network Policy Validation**: Ensure proper isolation
- **Secrets Management**: Audit secret usage and exposure
- **Compliance Checks**:  HIPAA, SOC2, PCI-DSS standards

## Automation & Self-Heealing 

### Core Automation Engine
```bash
// Extensible plugin architecture
type Remediator interface {
    CanFix(issueType string) bool
    Remediate(ctx context.Context, resource, namespace string) (*RemediationResult, error)
}
```

### Supported Remediations
- **Pod Restarts**: CrashLoopBackOff, ImagePullBackOff, ErrImagePull
- **Resource Optimization**: Auto-scaling, resource limit adjustments
- **Network Healing**: Service endpoint regeneration
- **Storage Recovery**:  Persistent volume claim management

### Predictive Scaling
- Machine learning-based workload forecasting
- Horizontal and vertical pod autoscaling
- Cost-optimized resource allocation
- Real-time metric analysis

## Contributing 
We welcome contributions from the open source community! K8s Lens is built for container orchestration in 
Kubernetes cloud environments and thrives on community input.

## Getting Started With Developement 
```bash
# Fork and clone the repository
git clone https://github.com/abrarahmad1510/k8s-lens
cd k8s-lens

# Set up development environment
make build
make test-all

# Run comprehensive tests
make test-phase4-complete
```
## Areas for Contributions: 
- **New Remediations**: Add automated fixes for common issues
- **Analysis Plugins**: Extend diagnostic capabilities
- **Integration Adapters**: Support for additional monitoring tools
- **Performance Optimization**:  Enhance scaling and resource usage

## Developement Workflow 
1. Fork the repository
2. Create a feature branch (git checkout -b feature/amazing-feature)
3. Commit your changes (git commit -m 'Add amazing feature')
4. Push to the branch (git push origin feature/amazing-feature)
5. Open a Pull Request

## Testing Standards
```bash
# Run the complete test suite
make test-all

# Specific test categories
make test-unit
make test-integration
make test-e2e
make test-security
```
## Technical Stack
### ⚙️ Backend & Core
- **Go**: High-performance CLI development
- **Cobra**: Modern CLI framework
- **Kubernetes Client-go**: Official Kubernetes client library
- **Prometheus**: Metrics collection and analysis

### 🤖 Machine Learning 
- **TensorFlow**: Predictive analytics and anomaly detection
- **Custom ML Models**: Workload forecasting and pattern recognition

### ✅ Monitoring & Integration
- **Prometheus Integration**: Real-time metrics
- **Grafana Dashboards**: Visualization and alerting
- **Multiple Cloud Providers**: AWS, GCP, Azure, and on-prem

### 🏅 Testing & Quality
- **Testify**: Assertion and mocking framework
- **Ginkgo**: BDD testing framework
- **Industrial Test Suite**: Comprehensive validation

## Developement 
### Project Structure 
```text
k8s-lens/
├── cmd/k8s-lens/          # CLI command definitions
├── pkg/                   # Core packages
│   ├── automation/        # Self-healing engine
│   ├── enterprise/        # Security features
│   ├── integrations/      # Third-party integrations
│   └── k8s/              # Kubernetes utilities
├── scripts/              # Build and test scripts
└── tests/               # Comprehensive test suites
```
### Build System
```bash
# Build the binary
make build

# Run tests
make test-all

# Create release binaries
make release

# Development build with hot-reload
make dev
```
### Testing Infrastructure
```bash
# Comprehensive test suite
./scripts/test-suite.sh

# Performance and stress testing
make test-stress

# Security validation
make test-security

# Regression testing
make test-regression
```
## License 
This project is licensed under the MIT Open Use License - see the LICENSE file for details.

## Acknowledgments
- Kubernetes Community for the amazing ecosystem
- Prometheus for robust metrics collection
- The Go community for excellent tooling
- All contributors who help improve K8s Lens



















