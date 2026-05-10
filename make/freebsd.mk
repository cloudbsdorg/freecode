# FreeBSD Makefile for Freecode
# Use with bmake or gmake

BINARY_NAME=freecode
INSTALL_DIR=/usr/local/bin

.PHONY: all build test clean install uninstall fmt tidy package

all: build

build:
	@echo "Building for FreeBSD..."
	go build -o $(BINARY_NAME) ./cmd/freecode

test:
	@echo "Running tests on FreeBSD..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf dist

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	install -c -m 755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

uninstall:
	@echo "Uninstalling from $(INSTALL_DIR)..."
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

fmt:
	go fmt ./...

tidy:
	go mod tidy

package: build
	@echo "Packaging for FreeBSD..."
	mkdir -p dist/freebsd
	tar -czf dist/freebsd/freecode-freebsd.tar.gz $(BINARY_NAME)
	@echo "Package created at dist/freebsd/freecode-freebsd.tar.gz"
	@if [ -d packaging/freebsd ]; then \
		echo "FreeBSD port files found at packaging/freebsd/"; \
	fi
