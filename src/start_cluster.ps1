$nodes = @(
    @{
        port = 8080
        id = "node1"
        isPrimary = $true
    },
    @{
        port = 8081
        id = "node2"
        isPrimary = $false
    },
    @{
        port = 8082
        id = "node3"
        isPrimary = $false
    }
)

# Get absolute paths
$rootDir = Split-Path $PSScriptRoot -Parent
$configTemplate = Join-Path $rootDir "config.json"

Write-Host "Starting cluster nodes..."

foreach ($node in $nodes) {
    # Load and modify config
    $config = Get-Content $configTemplate | ConvertFrom-Json
    $config.port = $node.port
    $config.cluster.node_id = $node.id
    $config.cluster.is_primary = $node.isPrimary
    
    # Create node-specific config
    $nodeConfigPath = Join-Path $PSScriptRoot "config_$($node.id).json"
    $config | ConvertTo-Json -Depth 10 | Set-Content $nodeConfigPath
    Write-Host "Created config for node $($node.id) at $nodeConfigPath"

    # Start node
    Write-Host "Starting node $($node.id) on port $($node.port)..."
    $startInfo = @{
        FilePath = "powershell"
        ArgumentList = "-NoExit", "-Command", "go run cmd/server/main.go -config `"$nodeConfigPath`""
        WorkingDirectory = $PSScriptRoot
        WindowStyle = "Normal"
    }
    Start-Process @startInfo

    # Wait and verify node is running
    Write-Host "Waiting for node $($node.id) to start..."
    $retries = 5
    $nodeStarted = $false
    
    while ($retries -gt 0 -and -not $nodeStarted) {
        Start-Sleep -Seconds 2
        try {
            $health = Invoke-RestMethod -Uri "http://localhost:$($node.port)/health"
            Write-Host "Node $($node.id) started successfully" -ForegroundColor Green
            $nodeStarted = $true
        } catch {
            $retries--
            if ($retries -gt 0) {
                Write-Host "Waiting for node $($node.id) to start... ($retries retries left)"
            }
        }
    }
    
    if (-not $nodeStarted) {
        Write-Host "Failed to start node $($node.id)" -ForegroundColor Red
    }
}

Write-Host "`nCluster startup complete!" -ForegroundColor Green