.PHONY: all build-cli build-server build-all test lint clean

all: build-cli build-server

build-cli:
	@echo "Building CLI..."
	@go build -o freecode ./cmd/freecode

build-server:
	@echo "Building server..."
	@go build -o freecode-server ./cmd/freecode-server

build-all: all

test:
	@echo "Running tests..."
	@go test ./...

lint:
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not found, skipping lint"

clean:
	@rm -f freecode freecode-server

fmt:
	@go fmt ./...

mod:
	@go mod tidy