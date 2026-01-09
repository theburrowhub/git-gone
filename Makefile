.PHONY: help build install test lint clean run dev docker-test release goreleaser-check goreleaser-build goreleaser-release-dry

# Variables
BINARY_NAME=git-gone
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-X git-gone/cmd.Version=$(VERSION) -X git-gone/cmd.CommitHash=$(COMMIT_HASH) -X git-gone/cmd.BuildTime=$(BUILD_TIME)
INSTALL_DIR?=$(HOME)/.local/bin

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)git-gone Makefile$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@go build -ldflags "$(LDFLAGS) -s -w" -o $(BINARY_NAME) .
	@echo "$(GREEN)✓ Built $(BINARY_NAME)$(NC)"

build-all: ## Build for all platforms
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS) -s -w" -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS) -s -w" -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS) -s -w" -o dist/$(BINARY_NAME)-macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS) -s -w" -o dist/$(BINARY_NAME)-macos-arm64 .
	@echo "$(GREEN)✓ Built binaries in dist/$(NC)"

install: build ## Build and install to ~/.local/bin (or INSTALL_DIR)
	@echo "$(BLUE)Installing $(BINARY_NAME) to $(INSTALL_DIR)...$(NC)"
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✓ Installed to $(INSTALL_DIR)/$(BINARY_NAME)$(NC)"

uninstall: ## Remove installed binary
	@echo "$(BLUE)Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)...$(NC)"
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run --timeout=5m
	@echo "$(GREEN)✓ Linting passed$(NC)"

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet passed$(NC)"

tidy: ## Tidy go modules
	@echo "$(BLUE)Tidying go modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Cleaned$(NC)"

run: build ## Build and run locally
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	@./$(BINARY_NAME)

dev: ## Run in development mode (no build cache)
	@echo "$(BLUE)Running in development mode...$(NC)"
	@go run -ldflags "$(LDFLAGS)" .

docker-test: ## Run installation tests in Docker
	@echo "$(BLUE)Running Docker installation tests...$(NC)"
	@./test-install.sh

version: ## Show version information
	@echo "Version:     $(VERSION)"
	@echo "Commit:      $(COMMIT_HASH)"
	@echo "Build Time:  $(BUILD_TIME)"

deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

update-deps: ## Update dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "$(GREEN)✓ All checks passed$(NC)"

release: clean check build-all ## Prepare release (clean, check, build all platforms)
	@echo "$(GREEN)✓ Release ready in dist/$(NC)"

ci: lint test ## Run CI checks
	@echo "$(GREEN)✓ CI checks passed$(NC)"

# ========== GORELEASER TARGETS ==========

goreleaser-check: ## Check goreleaser configuration
	@echo "$(BLUE)Checking goreleaser configuration...$(NC)"
	@goreleaser check
	@echo "$(GREEN)✓ GoReleaser config valid$(NC)"

goreleaser-build: ## Build snapshot with goreleaser (no publish)
	@echo "$(BLUE)Building with goreleaser (snapshot)...$(NC)"
	@goreleaser build --snapshot --clean
	@echo "$(GREEN)✓ Snapshot build complete$(NC)"

goreleaser-release-dry: ## Simulate full release with goreleaser (no publish)
	@echo "$(BLUE)Simulating release with goreleaser...$(NC)"
	@goreleaser release --snapshot --clean
	@echo "$(GREEN)✓ Dry-run release complete$(NC)"

.DEFAULT_GOAL := help
