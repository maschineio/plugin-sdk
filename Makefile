.PHONY: all build test clean install manifest-gen

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary names
MANIFEST_GEN_BINARY=manifest-gen

# Build directories
BIN_DIR=bin

all: test build

build: manifest-gen

manifest-gen:
	@echo "Building manifest generator..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(MANIFEST_GEN_BINARY) ./cmd/manifest-gen

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	$(GOCMD) clean

install: manifest-gen
	@echo "Installing manifest generator..."
	$(GOCMD) install ./cmd/manifest-gen

fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		sdk/plugin.proto

# Run all checks
check: fmt vet test

# Help
help:
	@echo "Available targets:"
	@echo "  make              - Run tests and build all binaries"
	@echo "  make build        - Build all binaries"
	@echo "  make manifest-gen - Build manifest generator tool"
	@echo "  make test         - Run all tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make install      - Install tools"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Run go vet"
	@echo "  make deps         - Update dependencies"
	@echo "  make proto        - Generate protobuf files"
	@echo "  make check        - Run all checks (fmt, vet, test)"
	@echo "  make help         - Show this help message"