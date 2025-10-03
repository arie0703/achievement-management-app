# Achievement Management Application Makefile

# Variables
APP_NAME := achievement-app
API_NAME := achievement-api
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go variables
GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)"

# Directories
BUILD_DIR := build
DIST_DIR := dist
CONFIG_DIR := config

# Default target
.PHONY: all
all: clean test build

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build binaries for current platform"
	@echo "  build-all      - Build binaries for all platforms"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Download dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  run-api        - Run API server"
	@echo "  run-cli        - Run CLI with help"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  install        - Install binaries to GOPATH/bin"
	@echo "  install-local  - Install binaries locally using install script"
	@echo "  package        - Create distribution packages"
	@echo "  version-info   - Show version information"
	@echo "  version-tag    - Create version tag (TYPE=patch|minor|major)"

# Build targets
.PHONY: build
build: build-api build-cli

.PHONY: build-api
build-api:
	@echo "Building API server..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(API_NAME) ./cmd/api

.PHONY: build-cli
build-cli:
	@echo "Building CLI application..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/cli

# Cross-compilation targets
.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)/linux-amd64
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/linux-amd64/$(API_NAME) ./cmd/api
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/linux-amd64/$(APP_NAME) ./cmd/cli

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)/darwin-amd64
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/darwin-amd64/$(API_NAME) ./cmd/api
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/darwin-amd64/$(APP_NAME) ./cmd/cli
	@mkdir -p $(BUILD_DIR)/darwin-arm64
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/darwin-arm64/$(API_NAME) ./cmd/api
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/darwin-arm64/$(APP_NAME) ./cmd/cli

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)/windows-amd64
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/windows-amd64/$(API_NAME).exe ./cmd/api
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/windows-amd64/$(APP_NAME).exe ./cmd/cli

# Test targets
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GO) test -v -tags=integration ./...

# Development targets
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run targets
.PHONY: run-api
run-api: build-api
	@echo "Starting API server..."
	./$(BUILD_DIR)/$(API_NAME)

.PHONY: run-cli
run-cli: build-cli
	@echo "Running CLI application..."
	./$(BUILD_DIR)/$(APP_NAME) --help

# Docker targets
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		-t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(APP_NAME):latest

.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image..."
	docker push $(APP_NAME):$(VERSION)
	docker push $(APP_NAME):latest

# Install targets
.PHONY: install
install: build
	@echo "Installing binaries..."
	$(GO) install $(LDFLAGS) ./cmd/api
	$(GO) install $(LDFLAGS) ./cmd/cli

# Package targets
.PHONY: package
package: build-all
	@echo "Creating distribution packages..."
	@mkdir -p $(DIST_DIR)
	
	# Linux package
	@tar -czf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64.tar.gz \
		-C $(BUILD_DIR)/linux-amd64 $(API_NAME) $(APP_NAME) \
		-C ../../$(CONFIG_DIR) . \
		-C .. README.md .env.example
	
	# macOS Intel package
	@tar -czf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64.tar.gz \
		-C $(BUILD_DIR)/darwin-amd64 $(API_NAME) $(APP_NAME) \
		-C ../../$(CONFIG_DIR) . \
		-C .. README.md .env.example
	
	# macOS Apple Silicon package
	@tar -czf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64.tar.gz \
		-C $(BUILD_DIR)/darwin-arm64 $(API_NAME) $(APP_NAME) \
		-C ../../$(CONFIG_DIR) . \
		-C .. README.md .env.example
	
	# Windows package
	@cd $(BUILD_DIR)/windows-amd64 && zip -r ../../$(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64.zip \
		$(API_NAME).exe $(APP_NAME).exe ../../$(CONFIG_DIR)/* ../../README.md ../../.env.example

# Clean targets
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean
	@echo "Cleaning all artifacts including dependencies..."
	$(GO) clean -modcache

# Development setup
.PHONY: setup
setup: deps
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Development environment setup complete!"

# Database setup (for local development)
.PHONY: setup-dynamodb
setup-dynamodb:
	@echo "Setting up local DynamoDB..."
	@if command -v docker >/dev/null 2>&1; then \
		docker run -d --name dynamodb-local -p 8000:8000 amazon/dynamodb-local; \
		echo "DynamoDB Local started on port 8000"; \
	else \
		echo "Docker not found. Please install Docker to run DynamoDB Local"; \
	fi

.PHONY: stop-dynamodb
stop-dynamodb:
	@echo "Stopping local DynamoDB..."
	@docker stop dynamodb-local || true
	@docker rm dynamodb-local || true

# Version info
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit Hash: $(COMMIT_HASH)"
	@echo "Go Version: $(shell go version)"
	@echo "Platform: $(GOOS)/$(GOARCH)"

# Version management
.PHONY: version-info
version-info:
	@./scripts/version.sh info

.PHONY: version-tag
version-tag:
	@./scripts/version.sh tag $(TYPE)

# Installation
.PHONY: install-local
install-local: build-all
	@./scripts/install.sh install