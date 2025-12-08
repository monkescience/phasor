# AI Agent Guide for Phasor

This document provides essential information for AI agents working with the Phasor project.

## Project Overview

Phasor is a Go-based microservices application consisting of:
- **Backend Service**: REST API that provides instance information (version, hostname, uptime, Go version, timestamp)
- **Frontend Service**: Web UI that displays backend instance information with colored tiles
- **Traefik**: Reverse proxy/load balancer for routing and multiple backend instances

The project uses OpenAPI-first design with code generation via oapi-codegen.

## Architecture

```
Frontend (port 8081) → Traefik (port 8080) → Backend (6 replicas)
```

- Frontend makes requests to backend through Traefik
- Backend exposes `/instance/info` endpoint
- Traefik load balances across multiple backend instances

## Technology Stack

- **Language**: Go 1.25.5
- **Build System**: Make + Go workspaces
- **API Spec**: OpenAPI 3.1.1
- **Code Generation**: oapi-codegen
- **Containerization**: Docker + Docker Compose
- **Orchestration**: Kubernetes/Helm
- **Reverse Proxy**: Traefik v3.2
- **Linting**: golangci-lint

## Project Structure

```
phasor/
├── backend/           # Backend service
│   ├── cmd/          # Main application entry point
│   ├── config/       # Configuration files
│   └── internal/     # Internal packages
│       ├── config/   # Configuration loading
│       └── instance/ # Generated API handlers
├── frontend/         # Frontend service
│   ├── cmd/          # Main application entry point
│   ├── config/       # Configuration files
│   └── internal/     # Internal packages
│       ├── config/   # Configuration loading
│       ├── frontend/ # HTTP handlers & templates
│       └── outgoing/ # Generated API client
├── openapi/          # OpenAPI specifications
│   ├── instance-api.yaml                    # API definition
│   ├── instance-api.oapi-codegen.server.yaml # Server generation config
│   └── instance-api.oapi-codegen.client.yaml # Client generation config
├── local/            # Local development setup
│   ├── docker-compose.yaml      # Docker Compose orchestration
│   ├── traefik.yaml             # Traefik static configuration
│   └── traefik-dynamic.yaml     # Traefik dynamic configuration
├── chart/            # Helm chart for deployment
├── Dockerfile        # Multi-stage build for both services
└── Makefile          # Build automation
```

## Building the Project

### Prerequisites
- Go 1.25.5 or later
- Docker (for containerized builds)
- golangci-lint (for linting)
- helm (for chart operations)

### Common Build Commands

```bash
# Display all available commands
make help

# Build both services
make build

# Build individual services
make build-backend
make build-frontend

# Generate code from OpenAPI specs
make generate

# Run tests
make test

# Run tests with coverage
make coverage

# Lint code
make lint

# Format code
make fmt

# Clean build artifacts
make clean
```

### Running Locally

```bash
# Option 1: Run with Go (requires two terminals)
# Terminal 1:
make run-backend

# Terminal 2:
make run-frontend

# Option 2: Run with Docker Compose (recommended)
make docker-up

# Stop Docker Compose services
make docker-down
```

When running with Docker Compose:
- Frontend: http://localhost:8081
- Backend (via Traefik): http://localhost:8080
- Traefik Dashboard: http://localhost:8090

### Docker Operations

```bash
# Build Docker images
make docker-build

# Build individual images
make docker-build-backend
make docker-build-frontend

# Start services with docker-compose
make docker-up

# Stop docker-compose services
make docker-down
```

### Helm Operations

```bash
# Lint Helm chart
make helm-lint

# Preview Helm templates
make helm-template
```

## Configuration

### Backend Configuration
- Location: `backend/config/config.yaml`
- Required environment variable: `VERSION`
- Configurable: log level, format, source inclusion

### Frontend Configuration
- Location: `frontend/config/config.yaml`
- Required setting: `backend_url` (must be set in config file)
- Configurable: backend URL, tile colors, log settings

## Code Generation

The project uses OpenAPI specifications to generate server and client code:

```bash
# Generate both server and client code
make generate
```

Generated files:
- `backend/internal/instance/server.gen.go` - Server interfaces and types
- `frontend/internal/outgoing/http/instance/client.gen.go` - Client code

**Important**: Always regenerate code after modifying OpenAPI specs in `openapi/instance-api.yaml`

## Development Workflow

1. Make changes to code or OpenAPI specs
2. If OpenAPI changed: `make generate`
3. Run tests: `make test`
4. Lint code: `make lint`
5. Build: `make build`
6. Test locally: `make docker-up`

## Dependencies

Dependencies are managed per-service using Go modules:

```bash
# Tidy dependencies for both services
make mod-tidy
```

The project uses a Go workspace (`go.work`) to coordinate the backend and frontend modules.

## Testing

```bash
# Run all tests with race detection
make test

# Generate coverage reports (HTML output in build/)
make coverage
```

## Key Files for AI Agents

When working with this project, pay attention to:

1. **Makefile** - All build commands and project automation
2. **openapi/instance-api.yaml** - API contract definition
3. **backend/internal/config/config.go** - Backend configuration structure
4. **frontend/internal/config/config.go** - Frontend configuration structure
5. **local/docker-compose.yaml** - Local development environment setup
6. **Dockerfile** - Multi-stage build configuration

## Common Tasks

### Adding a new API endpoint
1. Update `openapi/instance-api.yaml`
2. Run `make generate` to regenerate code
3. Implement handler in `backend/internal/instance/handler.go`
4. Update frontend client usage if needed
5. Run `make test` and `make lint`

### Modifying configuration
1. Update config structs in `backend/internal/config/config.go` or `frontend/internal/config/config.go`
2. Update corresponding `config/config.yaml` files
3. Update Helm chart values if needed

### Debugging deployment issues
1. Check `make helm-lint` output
2. Review `make helm-template` output
3. Verify configuration in `chart/values.yaml`
4. Test locally with `make docker-up`

## Notes for AI Agents

- The VERSION is determined from git tags via `git describe --tags --always --dirty`
- Backend must have VERSION environment variable set
- Frontend backend URL must be configured in `frontend/config/config.yaml`
- The project uses CGO_ENABLED=0 for static binaries
- All services have health check endpoints at `/health/live`
- Traefik runs 6 backend replicas by default in docker-compose
- Generated code should never be manually edited
