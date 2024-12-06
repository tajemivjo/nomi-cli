# Default target builds the CLI
all: build

# Build the CLI
build:
	go build -o nomi-cli

# Clean build artifacts
clean:
	rm -f nomi-cli

# Install dependencies
deps:
	go mod download

.PHONY: all build clean deps
