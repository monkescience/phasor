APP_NAME := phasor
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO_MODULES := backend frontend test
SERVICES := backend frontend

.PHONY: build test lint fmt clean docker-build docker-up docker-down helm-lint mod-tidy generate help

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build all binaries
	@for svc in $(SERVICES); do \
		mkdir -p build && cd $$svc && CGO_ENABLED=0 go build -o ../build/$$svc-service ./cmd/main.go && cd ..; \
	done

test: ## Run tests
	@cd test && go test -race ./...

lint: ## Run linter
	golangci-lint run --timeout=5m $(foreach mod,$(GO_MODULES),./$(mod)/...)

fmt: ## Format code
	golangci-lint fmt $(foreach mod,$(GO_MODULES),./$(mod)/...)

clean: ## Clean build artifacts
	rm -rf build

docker-build: ## Build Docker images
	@for svc in $(SERVICES); do \
		docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME)-$$svc:$(VERSION) $$svc; \
	done

docker-up: ## Start with docker-compose
	docker-compose -f local/docker-compose.yaml up --build

docker-down: ## Stop docker-compose
	docker-compose -f local/docker-compose.yaml down

helm-lint: ## Lint Helm chart
	helm lint chart

mod-tidy: ## Tidy Go modules
	@for mod in $(GO_MODULES); do cd $$mod && go mod tidy && cd ..; done

generate: ## Generate OpenAPI code
	cd backend && go generate ./...
