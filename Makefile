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

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests and show coverage in terminal
test-coverage-text:
	go test -v -cover ./...

.PHONY: all build clean deps test test-coverage test-coverage-text
