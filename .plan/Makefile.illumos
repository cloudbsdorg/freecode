# IllumOS Makefile for Freecode

BINARY_NAME=freecode
SERVER_BINARY_NAME=freecode-server
INSTALL_DIR=/usr/local/bin

.PHONY: all build test clean install uninstall fmt tidy package

all: build

build:
	@echo "Building for IllumOS..."
	go build -o $(BINARY_NAME) ./cmd/freecode
	go build -o $(SERVER_BINARY_NAME) ./cmd/freecode-server

test:
	@echo "Running tests on IllumOS..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) $(SERVER_BINARY_NAME)
	rm -rf dist

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	# IllumOS install command might vary, standard install usually works
	install -m 755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	install -m 755 $(SERVER_BINARY_NAME) $(INSTALL_DIR)/$(SERVER_BINARY_NAME)

uninstall:
	@echo "Uninstalling from $(INSTALL_DIR)..."
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	rm -f $(INSTALL_DIR)/$(SERVER_BINARY_NAME)

fmt:
	go fmt ./...

tidy:
	go mod tidy

package: build
	@echo "Packaging for IllumOS..."
	mkdir -p dist/illumos
	tar -czf dist/illumos/freecode-illumos.tar.gz $(BINARY_NAME) $(SERVER_BINARY_NAME)
	@echo "Package created at dist/illumos/freecode-illumos.tar.gz"
	@if [ -f packaging/illuminos/install.sh ]; then \
		echo "IllumOS installation script found at packaging/illuminos/install.sh"; \
	fi
