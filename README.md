# Distributed Load Balancer

A high-performance distributed load balancer written in Go that provides weighted round-robin distribution, health checking, and real-time metrics.

## Features

- ✨ Weighted round-robin load balancing
- 🏥 Automatic health checking of backend servers
- ⚡ Rate limiting with sliding window algorithm
- 📊 Real-time metrics and monitoring
- ⚖️ Configurable server weights
- 🔄 Dynamic server pool management

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Python 3.x (for test servers)
- PowerShell (for Windows)

### Installation

1. Clone the repository:

```powershell
git clone https://github.com/ykkalexx/loadbalancer
cd DistributedLoadBalancer
```

2. Install Deps

```powershell
cd src
go mod download
```

### Configuration

The load balancer is configured via config.json. Example configuration:

```json
{
  "port": 8080,
  "servers": [
    {
      "url": "http://localhost:5001",
      "weight": 2
    },
    {
      "url": "http://localhost:5002",
      "weight": 1
    },
    {
      "url": "http://localhost:5003",
      "weight": 1
    }
  ],
  "rate_limit": {
    "requests_per_second": 100,
    "burst_size": 20
  },
  "health_check": {
    "interval_seconds": 20,
    "timeout_seconds": 5,
    "max_failures": 3
  }
}
```

### Running the Load Balancer

Start the test servers:

```powershell
cd servers
.\start_servers.ps1
```

Start the load balancer:

```powershell
cd src
go run cmd/server/main.go
```

Test the load balancer:

```powershell
cd src/cmd/server
.\metrics.ps1
```

### Architecture

The load balancer consists of several key components:

- **Load Balancer Core**: Manages server selection and request distribution
- **Health Checker**: Monitors backend server health
- **Rate Limiter**: Prevents server overload
- **Metrics Collector**: Tracks performance metrics
