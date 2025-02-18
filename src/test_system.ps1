Write-Host "Starting Distributed Load Balancer Test Suite" -ForegroundColor Cyan

# Define paths
$serversScript = "C:\Projects\DistributedLoadBalancer\servers\start_servers.ps1"
$clusterScript = "C:\Projects\DistributedLoadBalancer\src\start_cluster.ps1"

# 1. Start backend servers first
Write-Host "`nStep 1: Starting backend servers..." -ForegroundColor Green
if (Test-Path $serversScript) {
    & $serversScript
    Start-Sleep -Seconds 3
} else {
    Write-Host "Error: Cannot find start_servers.ps1 at $serversScript" -ForegroundColor Red
    exit 1
}

# 2. Start the load balancer cluster
Write-Host "`nStep 2: Starting load balancer cluster..." -ForegroundColor Green
if (Test-Path $clusterScript) {
    & $clusterScript
    Start-Sleep -Seconds 3
} else {
    Write-Host "Error: Cannot find start_cluster.ps1 at $clusterScript" -ForegroundColor Red
    exit 1
}

# 3. Basic health check for all nodes
Write-Host "`nStep 3: Testing node health..." -ForegroundColor Green
$ports = @(8080, 8081, 8082)
foreach ($port in $ports) {
    Write-Host "`nChecking node on port $port"
    try {
        $health = Invoke-RestMethod -Uri "http://localhost:$port/health"
        Write-Host "Health check: OK - Node $($health.node_id)" -ForegroundColor Green
    } catch {
        Write-Host "Health check: Failed for port $port" -ForegroundColor Red
    }
}

# 4. Test load distribution
Write-Host "`nStep 4: Testing load distribution..." -ForegroundColor Green
Write-Host "Making 10 requests to each node..."

foreach ($port in $ports) {
    Write-Host "`nTesting load balancer on port $port"
    for ($i = 1; $i -le 10; $i++) {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:$port"
            Write-Host "Request $i : $($response.StatusCode) - $($response.Content)"
        } catch {
            Write-Host "Request $i : Failed" -ForegroundColor Red
        }
        Start-Sleep -Milliseconds 200
    }
}

# 5. Check metrics on all nodes
Write-Host "`nStep 5: Checking metrics on all nodes..." -ForegroundColor Green
foreach ($port in $ports) {
    Write-Host "`nMetrics for node on port $port"
    try {
        $metrics = Invoke-RestMethod -Uri "http://localhost:$port/metrics" | ConvertTo-Json
        Write-Host $metrics
    } catch {
        Write-Host "Failed to get metrics from port $port" -ForegroundColor Red
    }
}

Write-Host "`nTest suite completed!" -ForegroundColor Cyan