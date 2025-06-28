![Kubernetes Controllers on Go](https://github.com/user-attachments/assets/7cf40135-2f13-4204-85a3-5c1d8d20f44b)

# Go Kubernetes Controllers

A toolkit for building custom Kubernetes controllers and operators in Go. Designed for teams that need more control over their cluster automation than what's available off the shelf.

## What this is

This project gives you a solid foundation for writing your own Kubernetes controllers. Instead of wrestling with boilerplate code every time you need custom cluster behavior, you get working examples of the patterns that can be leveraged in production.

The codebase covers everything from basic CLI tools to multi-cluster management, with real implementations you can modify for your specific needs. No theoretical examples or toy demos.

## Getting started

```bash
# Build
go build -o controller .

# List deployments
./controller list

# Watch deployment events
./controller watch

# Start HTTP server
./controller server
```

## Commands

### List Deployments
```bash
./controller list --namespace default
./controller list --kubeconfig /path/to/kubeconfig
./controller list --timeout 60s
```

### Watch Events
```bash
./controller watch --namespace kube-system
./controller watch --workers 4 --resync 5m
./controller watch --in-cluster
```

### HTTP Server
```bash
./controller server --port 8080
curl http://localhost:8080/  # Returns "Hello from FastHTTP!"
```

### Global Options
```bash
--log-level string   Set log level: trace, debug, info, warn, error (default "info")
```

## Key Features

- **Event-Driven**: Kubernetes informer with multi-worker event processing
- **Error Handling**: Custom error types with proper validation
- **Structured Logging**: Zerolog integration with configurable levels
- **CLI Interface**: Cobra-based commands for list, watch, and server
- **Testing**: Unit and integration tests with envtest

## Development

```bash
# Run tests
go test ./...

# Integration tests
go test -tags=integration ./pkg/informer/

# Development
go run main.go watch --namespace default --log-level debug
```

## Architecture

```
Kubernetes API → Informer → Event Queue → Workers → Processing
```

## Dependencies

- **Go 1.24+** - Language runtime
- **k8s.io/client-go v0.33.0** - Kubernetes client
- **github.com/spf13/cobra** - CLI framework
- **github.com/rs/zerolog** - Structured logging
- **github.com/valyala/fasthttp** - HTTP server

---

*Built for production Kubernetes workloads with reliability and observability in mind.*
