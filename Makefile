# Inu - AI-Powered Text Anonymization Tool

.PHONY: help build build-all test lint clean

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build variables
BINARY_NAME := inu
BUILD_DIR := bin
CMD_DIR := cmd/inu
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Platform targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build binary for current platform
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build binaries for all platforms
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(platform))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(platform))))\
		$(eval OUTPUT := $(BUILD_DIR)/$(BINARY_NAME)-$(OS)-$(ARCH)$(if $(filter windows,$(OS)),.exe,))\
		echo "  Building $(OS)/$(ARCH)..." && \
		GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(OUTPUT) ./$(CMD_DIR) || exit 1;)
	@echo "Build complete for all platforms"
	@ls -lh $(BUILD_DIR)

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -cover ./...

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

.DEFAULT_GOAL := help
