// Dashboard functionality
document.addEventListener('DOMContentLoaded', function() {
    // Load cluster info on dashboard
    if (document.getElementById('cluster-version')) {
        loadClusterInfo();
    }

    // Load cluster contexts on multi-cluster page
    if (document.getElementById('cluster-contexts')) {
        loadClusterContexts();
        setupComparisonForm();
    }
});

function loadClusterInfo() {
    fetch('/api/cluster/info')
        .then(response => response.json())
        .then(data => {
            document.getElementById('cluster-version').textContent = data.clusterVersion || 'Unknown';
            document.getElementById('node-count').textContent = data.nodes || '0';
            document.getElementById('pod-count').textContent = data.pods || '0';
            
            const statusElement = document.getElementById('cluster-status');
            statusElement.textContent = data.status || 'Unknown';
            statusElement.className = 'status ' + (data.status === 'healthy' ? 'healthy' : 'degraded');
        })
        .catch(error => {
            console.error('Error loading cluster info:', error);
            document.getElementById('cluster-version').textContent = 'Error';
            document.getElementById('node-count').textContent = 'Error';
            document.getElementById('pod-count').textContent = 'Error';
            document.getElementById('cluster-status').textContent = 'Error';
        });
}

function loadClusterContexts() {
    fetch('/api/multicluster/contexts')
        .then(response => response.json())
        .then(data => {
            const container = document.getElementById('cluster-contexts');
            container.innerHTML = '';

            if (data.contexts && data.contexts.length > 0) {
                data.contexts.forEach(context => {
                    const isCurrent = context === data.currentContext;
                    const card = document.createElement('div');
                    card.className = `context-card ${isCurrent ? 'current' : ''}`;
                    card.innerHTML = `
                        <h3>${context}</h3>
                        <p>Status: <span class="status healthy">Healthy</span></p>
                        ${isCurrent ? '<p><strong>Current Context</strong></p>' : ''}
                    `;
                    container.appendChild(card);
                });
            } else {
                container.innerHTML = '<p>No cluster contexts found. Check your kubeconfig.</p>';
            }
        })
        .catch(error => {
            console.error('Error loading cluster contexts:', error);
            document.getElementById('cluster-contexts').innerHTML = 
                '<p class="error">Error loading cluster contexts: ' + error.message + '</p>';
        });
}

function setupComparisonForm() {
    const form = document.getElementById('comparison-form');
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            const resourceType = document.getElementById('resource-type').value;
            compareClusters(resourceType);
        });
    }
}

function compareClusters(resourceType) {
    const resultsDiv = document.getElementById('comparison-results');
    resultsDiv.innerHTML = '<div class="loading">Comparing clusters...</div>';

    fetch(`/api/multicluster/compare/${resourceType}`)
        .then(response => response.json())
        .then(data => {
            resultsDiv.innerHTML = `
                <h3>Comparison Results: ${resourceType}</h3>
                <pre>${data.comparison || 'No data'}</pre>
                <p><strong>Differences found:</strong> ${data.differences || 0}</p>
            `;
        })
        .catch(error => {
            console.error('Error comparing clusters:', error);
            resultsDiv.innerHTML = `<p class="error">Error comparing clusters: ${error.message}</p>`;
        });
}

function runFederatedAnalysis() {
    const resultsDiv = document.getElementById('federated-results');
    resultsDiv.innerHTML = '<div class="loading">Running federated analysis...</div>';

    fetch('/api/multicluster/federated')
        .then(response => response.json())
        .then(data => {
            resultsDiv.innerHTML = `
                <h3>Federated Analysis Results</h3>
                <pre>${data.report || 'No data'}</pre>
                <div class="summary">
                    <p><strong>Total Clusters:</strong> ${data.summary.totalClusters}</p>
                    <p><strong>Healthy Clusters:</strong> ${data.summary.healthyClusters}</p>
                    <p><strong>Overall Health:</strong> <span class="status ${data.summary.overallHealth.toLowerCase()}">${data.summary.overallHealth}</span></p>
                </div>
            `;
        })
        .catch(error => {
            console.error('Error running federated analysis:', error);
            resultsDiv.innerHTML = `<p class="error">Error running federated analysis: ${error.message}</p>`;
        });
}

function runQuickScan() {
    alert('Quick health scan would run here. This would integrate with the existing analysis capabilities.');
}

// Analysis page functionality
function analyzeResource() {
    const resourceType = document.getElementById('analysis-resource-type').value;
    const resourceName = document.getElementById('analysis-resource-name').value;
    const namespace = document.getElementById('analysis-namespace').value;

    const resultsDiv = document.getElementById('analysis-results');
    resultsDiv.innerHTML = '<div class="loading">Analyzing resource...</div>';

    fetch(`/api/analysis/${resourceType}/${resourceName}?namespace=${namespace}`)
        .then(response => response.json())
        .then(data => {
            resultsDiv.innerHTML = `
                <h3>Analysis Results: ${resourceType}/${resourceName}</h3>
                <pre>${JSON.stringify(data, null, 2)}</pre>
            `;
        })
        .catch(error => {
            console.error('Error analyzing resource:', error);
            resultsDiv.innerHTML = `<p class="error">Error analyzing resource: ${error.message}</p>`;
        });
}

function optimizeNamespace() {
    const namespace = document.getElementById('optimization-namespace').value;
    
    const resultsDiv = document.getElementById('optimization-results');
    resultsDiv.innerHTML = '<div class="loading">Analyzing optimizations...</div>';

    fetch(`/api/optimization/${namespace}`)
        .then(response => response.json())
        .then(data => {
            resultsDiv.innerHTML = `
                <h3>Optimization Results: ${namespace}</h3>
                <pre>${JSON.stringify(data, null, 2)}</pre>
                <p><strong>Estimated Savings:</strong> $${data.estimatedSavings || 0}</p>
            `;
        })
        .catch(error => {
            console.error('Error optimizing namespace:', error);
            resultsDiv.innerHTML = `<p class="error">Error optimizing namespace: ${error.message}</p>`;
        });
}