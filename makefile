# Makefile for sysfail project

# Binaries and directories
BINARY_NAME = sysfail
CMD_DIR     = ./cmd/sysfail
BIN_DIR     = /usr/local/bin

.PHONY: all build install fmt vet test clean tidy

all: fmt vet build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_DIR)

install: build
	@echo "Installing $(BINARY_NAME) to $(BIN_DIR)..."
	@mv $(BINARY_NAME) $(BIN_DIR)/

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)

tidy:
	@echo "Tidying go.mod..."
	@go mod tidy