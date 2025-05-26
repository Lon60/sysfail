BINARY_NAME=sysfail
VERSION=1.0.0
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin
GO_FILES=$(shell find . -name "*.go" -type f)

.PHONY: all build install uninstall clean test fmt vet release

all: build

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sysfail
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"


install: build
	@echo "Installing $(BINARY_NAME) to $(DESTDIR)$(INSTALL_DIR)..."
	install -Dm755 $(BUILD_DIR)/$(BINARY_NAME) $(DESTDIR)$(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(DESTDIR)$(INSTALL_DIR)/$(BINARY_NAME)"

uninstall:
	@echo "Removing $(BINARY_NAME) from $(DESTDIR)$(INSTALL_DIR)..."
	rm -f $(DESTDIR)$(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Clean complete"

test:
	@echo "Running tests..."
	go test -v ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Vetting code..."
	go vet ./...

release: clean
	@echo "Creating release build..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sysfail
	@echo "Release build complete: $(BUILD_DIR)/$(BINARY_NAME)"

help:
	@echo "Available targets:"
	@echo "  build     - Build the binary only"
	@echo "  install   - Build and install the binary (requires sudo)"
	@echo "  uninstall - Remove installed binary (requires sudo)"
	@echo "  clean     - Remove build artifacts"
	@echo "  test      - Run tests"
	@echo "  fmt       - Format Go code"
	@echo "  vet       - Run go vet"
	@echo "  release   - Create optimized release build"
	@echo "  help      - Show this help"