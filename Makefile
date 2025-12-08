# Makefile for reference-service-go monorepo

# Variables
APP_NAME := reference-service
BUILD_DIR := ./build
CHART_PATH := ./chart
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOGENERATE := $(GOCMD) generate

# Build flags
CGO_ENABLED := 0

# Service-specific variables
BACKEND_DIR := ./backend
FRONTEND_DIR := ./frontend
BACKEND_BINARY := backend-service
FRONTEND_BINARY := frontend-service

# Phony targets
.PHONY: all build build-backend build-frontend run run-backend run-frontend generate test coverage fmt lint clean docker-build docker-build-backend docker-build-frontend docker-up docker-down helm-lint helm-template mod-tidy help

# Default target
all: help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  build              - Build both backend and frontend binaries"
	@echo "  build-backend      - Build the backend binary"
	@echo "  build-frontend     - Build the frontend binary"
	@echo "  run                - Run both services locally (use separate terminals)"
	@echo "  run-backend        - Run the backend service locally"
	@echo "  run-frontend       - Run the frontend service locally"
	@echo "  generate           - Run go generate to generate code from OpenAPI specs"
	@echo "  test               - Run tests with race detection for both services"
	@echo "  coverage           - Run tests with coverage report"
	@echo "  fmt                - Format Go code using golangci-lint"
	@echo "  lint               - Run golangci-lint on both services"
	@echo "  clean              - Remove build artifacts"
	@echo "  docker-build       - Build both Docker images"
	@echo "  docker-build-backend  - Build the backend Docker image"
	@echo "  docker-build-frontend - Build the frontend Docker image"
	@echo "  docker-up          - Start both services with docker-compose"
	@echo "  docker-down        - Stop and remove docker-compose services"
	@echo "  helm-lint          - Lint the Helm chart"
	@echo "  helm-template      - Render Helm chart templates"
	@echo "  mod-tidy           - Tidy Go module dependencies for both services"

## build: Build both backend and frontend binaries
build: build-backend build-frontend

## build-backend: Build the backend binary
build-backend:
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -o ../$(BUILD_DIR)/$(BACKEND_BINARY) ./cmd/main.go

## build-frontend: Build the frontend binary
build-frontend:
	@mkdir -p $(BUILD_DIR)
	cd $(FRONTEND_DIR) && CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -o ../$(BUILD_DIR)/$(FRONTEND_BINARY) ./cmd/main.go

## run: Instructions to run both services
run:
	@echo "To run both services, open two terminals:"
	@echo "Terminal 1: make run-backend"
	@echo "Terminal 2: make run-frontend"

## run-backend: Run the backend service locally
run-backend:
	cd $(BACKEND_DIR) && VERSION=$(VERSION) $(GORUN) ./cmd/main.go -config config/config.yaml

## run-frontend: Run the frontend service locally
run-frontend:
	cd $(FRONTEND_DIR) && VERSION=$(VERSION) $(GORUN) ./cmd/main.go -config config/config.yaml

## generate: Run go generate to generate code from OpenAPI specs
generate:
	cd $(BACKEND_DIR) && $(GOGENERATE) ./...

## test: Run tests with race detection for both services
test:
	cd $(BACKEND_DIR) && $(GOTEST) -v -race ./...
	cd $(FRONTEND_DIR) && $(GOTEST) -v -race ./...

## coverage: Run tests with coverage report
coverage:
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && $(GOTEST) -v -race -coverprofile=../$(BUILD_DIR)/coverage-backend.out -covermode=atomic ./...
	cd $(FRONTEND_DIR) && $(GOTEST) -v -race -coverprofile=../$(BUILD_DIR)/coverage-frontend.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(BUILD_DIR)/coverage-backend.out -o $(BUILD_DIR)/coverage-backend.html
	$(GOCMD) tool cover -html=$(BUILD_DIR)/coverage-frontend.out -o $(BUILD_DIR)/coverage-frontend.html

## fmt: Format Go code
fmt:
	cd $(BACKEND_DIR) && golangci-lint fmt
	cd $(FRONTEND_DIR) && golangci-lint fmt

## lint: Run golangci-lint on both services
lint:
	cd $(BACKEND_DIR) && golangci-lint run --timeout=5m
	cd $(FRONTEND_DIR) && golangci-lint run --timeout=5m

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
	cd $(BACKEND_DIR) && $(GOCLEAN)
	cd $(FRONTEND_DIR) && $(GOCLEAN)

## docker-build: Build both Docker images
docker-build: docker-build-backend docker-build-frontend

## docker-build-backend: Build the backend Docker image
docker-build-backend:
	docker build --build-arg VERSION=$(VERSION) --target backend-runtime -t $(APP_NAME)-backend:$(VERSION) -t $(APP_NAME)-backend:latest .

## docker-build-frontend: Build the frontend Docker image
docker-build-frontend:
	docker build --build-arg VERSION=$(VERSION) --target frontend-runtime -t $(APP_NAME)-frontend:$(VERSION) -t $(APP_NAME)-frontend:latest .

## docker-up: Start both services with docker-compose
docker-up:
	docker-compose -f local/docker-compose.yaml up --build

## docker-down: Stop and remove docker-compose services
docker-down:
	docker-compose -f local/docker-compose.yaml down

## helm-lint: Lint the Helm chart
helm-lint:
	helm lint $(CHART_PATH)

## helm-template: Render Helm chart templates
helm-template:
	helm template $(APP_NAME) $(CHART_PATH)

## mod-tidy: Tidy Go module dependencies for both services
mod-tidy:
	cd $(BACKEND_DIR) && $(GOMOD) tidy
	cd $(FRONTEND_DIR) && $(GOMOD) tidy
