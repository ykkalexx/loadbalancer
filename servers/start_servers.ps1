$ports = @(5001, 5002, 5003)

# array to store background jobs
$jobs = @()

# Start each server in a new window
foreach ($port in $ports) {
    $windowTitle = "Flask Server on port $port"
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "python server.py $port" -WorkingDirectory $PSScriptRoot
    Write-Host "Started server on port $port"
}

Write-Host "`nAll test servers are running!"
Write-Host "To test, try these URLs:"
foreach ($port in $ports) {
    Write-Host "http://localhost:$port"
}
Write-Host "`nPress Ctrl+C in each window to stop the servers"

#make multiple requests:
for ($i = 1; $i -le 10; $i++) {
    Write-Host "Request $i"
    Invoke-WebRequest -Uri "http://localhost:8080" | Select-Object -ExpandProperty Content
    Start-Sleep -Milliseconds 500
}

# Make requests every second for 30 seconds
$start = Get-Date
while ((Get-Date) -lt ($start.AddSeconds(30))) {
    Invoke-WebRequest -Uri "http://localhost:8080" | Select-Object -ExpandProperty Content
    Start-Sleep -Seconds 1
}