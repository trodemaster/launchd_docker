# Default target architecture is the host system
GOARCH ?= $(shell go env GOARCH)
GOOS ?= $(shell go env GOOS)

# Binary name
BINARY_NAME=launchd_docker

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all clean build build-amd64 build-arm64

all: build

build:
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build $(LDFLAGS) -o $(BINARY_NAME) cmd/launchd_docker/main.go

build-amd64:
	GOARCH=amd64 GOOS=darwin go build $(LDFLAGS) -o $(BINARY_NAME)_amd64 cmd/launchd_docker/main.go

build-arm64:
	GOARCH=arm64 GOOS=darwin go build $(LDFLAGS) -o $(BINARY_NAME)_arm64 cmd/launchd_docker/main.go

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)_amd64 $(BINARY_NAME)_arm64

# Help target
help:
	@echo "Available targets:"
	@echo "  all        - Build for current architecture (default)"
	@echo "  build      - Build for current architecture"
	@echo "  build-amd64 - Build for macOS x86_64"
	@echo "  build-arm64 - Build for macOS arm64"
	@echo "  clean      - Remove built binaries"
	@echo "  help       - Show this help message" 