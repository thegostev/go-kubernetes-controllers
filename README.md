# Go Kubernetes Controllers

A toolkit for building custom Kubernetes controllers and operators in Go. Designed for teams that need more control over their cluster automation than what's available off the shelf.

## What this is

This project gives you a solid foundation for writing your own Kubernetes controllers. Instead of wrestling with boilerplate code every time you need custom cluster behavior, you get working examples of the patterns that can be leveraged in production.

The codebase covers everything from basic CLI tools to multi-cluster management, with real implementations you can modify for your specific needs. No theoretical examples or toy demos.

## Getting started

```bash
# Get the code
git clone https://github.com/yourusername/go-kubernetes-controllers
cd go-kubernetes-controllers

# Build it
make build

# Run the server
./bin/controller server --log-level info

# Try some commands
./bin/controller list deployments
```

## What's inside

```
cmd/                    # Main application code
├── root.go            # CLI setup and global flags
├── server.go          # HTTP API server
└── controllers/       # Custom controller logic

pkg/                   # Reusable libraries
├── clients/           # Kubernetes API clients
├── informers/         # Resource watchers
└── reconcilers/       # Controller business logic

charts/                # Helm deployment
deploy/                # Raw Kubernetes manifests
```

## Key components

**CLI interface** - Built with Cobra. Handles the usual suspects: listing resources, managing deployments, debugging cluster state. Nothing fancy, just reliable tools for day-to-day operations.

**HTTP API** - FastHTTP-based server that exposes controller functionality over REST. Includes health checks, metrics endpoints, and proper error handling. Swagger docs included.

**Controller framework** - The meat of the project. Handles resource watching, event processing, and state reconciliation. Supports custom resources and multi-cluster scenarios.

**Observability** - Structured logging with zerolog, Prometheus metrics, and OpenTelemetry tracing. Everything you need to understand what your controllers are doing in production.

## Configuration

Controllers read from YAML config files or environment variables:

```yaml
server:
  port: 8080
  
controllers:
  deployment_sync:
    enabled: true
    workers: 5
    
logging:
  level: info
  format: json
```

You can also configure everything through CLI flags if that's more your style.

## Deployment options

**Local development:**
```bash
make run-dev
```

**Docker:**
```bash
make docker-build
docker run -p 8080:8080 your-controller:latest
```

**Kubernetes:**
```bash
helm install my-controllers charts/app/
```

The Helm chart includes RBAC, service accounts, and monitoring setup. It's designed to work with standard Kubernetes distributions without requiring custom operators or CRDs (unless you want them).

## Built with

Go 1.24+, controller-runtime, client-go for the Kubernetes parts. FastHTTP for the web server because it's fast and doesn't get in your way. Standard observability stack with Prometheus and OpenTelemetry.

Security-wise, everything runs in distroless containers with minimal privileges. JWT authentication is available if you need it, along with RBAC integration.

## Performance notes

The controller can handle several thousand events per second on modest hardware. Informer caching keeps memory usage reasonable even with large clusters. Leader election works properly for HA deployments.

Response times for the HTTP API are typically under 5ms for simple operations. More complex multi-cluster queries might take longer depending on network latency.

## Why another controller framework

Most existing solutions either oversimplify things to the point of being useless, or they're so complex that you spend more time fighting the framework than solving your actual problems.

This project hits the middle ground. It gives you working code that handles the boring parts (RBAC, metrics, leader election) while staying out of your way when you need to implement custom logic.

---

*Made for teams that run real Kubernetes workloads and need reliable automation.*
