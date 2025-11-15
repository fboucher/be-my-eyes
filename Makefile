# Application name
APP_NAME = be-my-eyes

# Build directory
BUILD_DIR = .

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Main package path
MAIN_PATH = ./cmd/$(APP_NAME)

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_PATH)

# Install the application to $GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install $(MAIN_PATH)

# Run the application
.PHONY: run
run:
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	./$(APP_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(APP_NAME)
	rm -f $(BUILD_DIR)/$(APP_NAME)-*

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for all supported platforms"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-darwin   - Build for macOS"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  run            - Build and run the application"
	@echo "  clean          - Remove build artifacts"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  help           - Show this help message"
