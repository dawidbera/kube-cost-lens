# Makefile for kube-cost-lens
# Provides targets for building, testing, linting, and running the application.

APP_NAME := kube-cost-lens
BIN_DIR := bin
BINARY := $(BIN_DIR)/$(APP_NAME)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "manual")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X github.com/dawidbera/kube-cost-lens/cmd.Version=$(VERSION) \
	-X github.com/dawidbera/kube-cost-lens/cmd.GitCommit=$(COMMIT) \
	-X github.com/dawidbera/kube-cost-lens/cmd.BuildDate=$(BUILD_DATE)

.PHONY: all build test test-coverage lint run run-demo clean docker-build help

all: build

## build: Compile the binary for the local platform
build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) main.go

## test: Run unit tests with race detection
test:
	go test -v -race -cover ./...

## test-coverage: Run unit tests and generate coverage report
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## lint: Run go vet and golangci-lint if installed
lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, ran go vet successfully."; \
	fi

## run: Build and run against local kubeconfig
run: build
	$(BINARY)

## run-demo: Build and run in synthetic offline demo mode
run-demo: build
	$(BINARY) --demo

## clean: Remove build artifacts and temporary files
clean:
	rm -rf $(BIN_DIR) coverage.out patch-*.yaml

## docker-build: Build minimal Docker container
docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

## help: Display available make targets
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':'
