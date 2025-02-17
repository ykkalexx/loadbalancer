Invoke-RestMethod -Uri "http://localhost:8080/metrics" | ConvertTo-Json

# making some test requests
for ($i = 1; $i -le 10; $i++) {
    Invoke-WebRequest -Uri "http://localhost:8080"
    Start-Sleep -Milliseconds 500
}

# check metrics again to see the changes
Invoke-RestMethod -Uri "http://localhost:8080/metrics" | ConvertTo-Json