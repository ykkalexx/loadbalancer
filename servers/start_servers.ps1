$ports = @(5001, 5002, 5003)

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