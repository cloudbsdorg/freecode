# macOS Makefile for Freecode

BINARY_NAME=freecode
INSTALL_DIR=/usr/local/bin

.PHONY: all build test clean install uninstall fmt tidy package

all: build

build:
	@echo "Building for macOS..."
	go build -o $(BINARY_NAME) ./cmd/freecode

test:
	@echo "Running tests on macOS..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf dist

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	install -m 755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

uninstall:
	@echo "Uninstalling from $(INSTALL_DIR)..."
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

fmt:
	go fmt ./...

tidy:
	go mod tidy

package: build
	@echo "Packaging for macOS..."
	mkdir -p dist/macos
	tar -czf dist/macos/freecode-macos.tar.gz $(BINARY_NAME)
	@echo "Package created at dist/macos/freecode-macos.tar.gz"
	@if [ -f packaging/macos/freecode.rb ]; then \
		echo "Homebrew formula found at packaging/macos/freecode.rb"; \
	fi
