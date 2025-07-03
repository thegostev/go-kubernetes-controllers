![Kubernetes Controllers on Go](https://github.com/user-attachments/assets/7cf40135-2f13-4204-85a3-5c1d8d20f44b)



[![Build Status](https://github.com/thegostev/go-kubernetes-controllers/actions/workflows/ci.yml/badge.svg)](https://github.com/thegostev/go-kubernetes-controllers/actions)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](https://github.com/thegostev/go-kubernetes-controllers/pulls)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

# Go Kubernetes Controllers

A toolkit for building Kubernetes controllers and operators in Go.  
Designed for teams who need production-grade automation and observability in their clusters.

---

## Features

- Easy CLI with Cobra.
- Leader election for HA deployments.
- Helm chart for Kubernetes deployment.
- Structured logging and Prometheus metrics.
- Event-driven controller-runtime architecture.

---

## Quick Start

```sh
git clone https://github.com/thegostev/go-kubernetes-controllers.git
cd go-kubernetes-controllers
go build -o controller .

# Run locally (no leader election)
./controller server --disable-leader-election

# Run in cluster (with leader election)
./controller server
```

---

## Usage

### List Deployments

```sh
./controller list --namespace default
./controller list --kubeconfig /path/to/kubeconfig
./controller list --timeout 60s
```

### Watch Deployment Events

```sh
./controller watch --namespace kube-system
./controller watch --workers 4 --resync 5m
./controller watch --in-cluster
```

### Start Controller Manager

```sh
./controller server --disable-leader-election --metrics-port 9000
```

### HTTP Server (legacy)

```sh
./controller server --port 8080
curl http://localhost:8080/  # Returns "Hello from FastHTTP!"
```

### Global Options

```sh
--log-level string   Set log level: trace, debug, info, warn, error (default "info")
```

---

## Configuration

| Flag                        | Description                          | Default   |
|-----------------------------|--------------------------------------|-----------|
| `--disable-leader-election` | Disable leader election (dev only)   | `false`   |
| `--metrics-port`            | Metrics endpoint port (if supported) | `8081`    |
| `--log-level`               | Log level (trace, debug, info, ...)  | `info`    |
| `--namespace`               | Namespace for list/watch commands    | `default` |
| `--kubeconfig`              | Path to kubeconfig file              | `~/.kube/config` |

---

## Architecture

Kubernetes API → controller-runtime Manager → Controllers → Reconcile Loop → Logging/Metrics

- Modular: `cmd/` for CLI, `pkg/` for controllers, informers, k8s clients, `internal/` for types/validation.
- Clean separation of concerns and testable interfaces.

---

## Development

```sh
# Run all tests and lints
make install test lint fmt vet build
```

Integration tests (envtest):
```sh
go test -tags=integration ./pkg/controller/
```

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on code style, PR process, and how to get started.

---

## License

[MIT](LICENSE)

---

## Links

- [Issues](https://github.com/thegostev/go-kubernetes-controllers/issues)
- [GitHub Actions CI](https://github.com/thegostev/go-kubernetes-controllers/actions)
- [Helm Chart](./charts/app/)

---

*Built for production Kubernetes workloads with reliability and observability in mind.*
