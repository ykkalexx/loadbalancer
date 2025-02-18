3. Create the API documentation:

````markdown
# API Endpoints

## Health Check Endpoint

```http
GET /health
```

- Response:

```json
{
  "status": "healthy",
  "node_id": "node1",
  "is_primary": true,
  "timestamp": "2024-02-18T15:04:05Z"
}
```
````

```http
GET /metrics
```

- Response:

```json
{
  "total_requests": 1000,
  "successful_requests": 950,
  "failed_requests": 50,
  "server_count": 3,
  "healthy_servers": 3,
  "cluster_info": {
    "node_id": "node1",
    "is_primary": true,
    "healthy_nodes": ["http://localhost:8081", "http://localhost:8082"]
  }
}
```
