SWAGGERCLI := $(shell command -v swagger-cli 2>/dev/null)

.PHONY: help ensure-swagger-cli bundle generate build run clean test install-tools docker-build docker-run docker-compose-up docker-compose-down docker-compose-logs docker-compose-build docker-compose-restart docker-clean

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install-tools: ## Install required tools
	@echo "Installing oapi-codegen..."
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Tools installed successfully"

ensure-swagger-cli:
ifndef SWAGGERCLI
	@echo "Installing @apidevtools/swagger-cli globally..."
	@npm install -g @apidevtools/swagger-cli
else
	@echo "swagger-cli is already installed: $(SWAGGERCLI)"
endif

bundle: ensure-swagger-cli
	@echo "Bundling OpenAPI spec..."
	@swagger-cli bundle ./api/openapi.yaml -o ./api/openapi.bundled.json

generate: bundle ## Generate API code from OpenAPI spec
	@echo "Generating API code from split specs..."
	@go generate ./...

run: generate ## Run the application
	@if [ -n "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Running application with config: $(filter-out $@,$(MAKECMDGOALS))"; \
		go run ./cmd/server/main.go -config $(filter-out $@,$(MAKECMDGOALS)); \
	else \
		echo "Running application..."; \
		go run ./cmd/server/main.go; \
	fi

%:
	@:

clean: ## Clean build artifacts and generated files
	@echo "Cleaning..."
	@rm -rf build/
	@find pkg/api -name "*.gen.go" -type f -delete
	@rm -f coverage.out coverage.html
	@rm -f api/swagger.json
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Tests complete. Coverage report: coverage.html"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Formatting complete"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Linting complete"

dev: generate ## Run in development mode with hot reload (requires air)
	@echo "Starting development server..."
	@air

PACKAGE_NAME = github.com/quanNaNbk03/monitor-demo

VERSION ?= $(shell git describe --abbrev=0 --tags)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
SHA ?= $(shell git describe --match=none --always --abbrev=8)
BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%z)
LICENSE_PUBLIC_KEY ?=

ORGANIZATION = OCN

LDFLAGS = "-s -w -X ${PACKAGE_NAME}/buildinfo.Version=${VERSION} \
         -X ${PACKAGE_NAME}/buildinfo.GitBranch=${BRANCH}.${SHA} \
         -X ${PACKAGE_NAME}/buildinfo.BuildDate=${BUILD_TIME} \
         -X ${PACKAGE_NAME}/buildinfo.Organization=${ORGANIZATION} \
         -X $(PACKAGE_NAME)/buildinfo.LicensePublicKey=$(LICENSE_PUBLIC_KEY)"

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

BUILD_DIR ?= ./cmd/server/
OUTPUT_FILE ?= ./odns-api-admin/opt/ocn/odns-api-admin/bin/odns-api-admin

GO_BUILD = go build -buildvcs=false -a -installsuffix cgo

.PHONY: build
build:
	@if [ -z "$(BUILD_DIR)" ] || [ -z "$(OUTPUT_FILE)" ]; then \
  		echo "Error: please set BUILD_DIR and OUTPUT_FILE(VD: BUILD_DIR=./cmd/myapp/ OUTPUT_FILE=./myapp)"; \
  		exit 1; \
  	fi
	@mkdir -p $(dir $(OUTPUT_FILE))
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) -ldflags=$(LDFLAGS) -o $(OUTPUT_FILE) $(BUILD_DIR)