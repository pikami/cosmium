GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

BINARY_NAME=cosmium
SERVER_LOCATION=./cmd/server

SHARED_LIB_LOCATION=./sharedlibrary
SHARED_LIB_OPT=-buildmode=c-shared

DIST_DIR=dist

all: test build-all

build-all: build-darwin-arm64 build-darwin-amd64 build-linux-amd64 build-linux-arm64 build-windows-amd64 build-windows-arm64

build-darwin-arm64:
	@echo "Building macOS ARM binary..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(SERVER_LOCATION)

build-darwin-amd64:
	@echo "Building macOS x64 binary..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(SERVER_LOCATION)

build-linux-amd64:
	@echo "Building Linux x64 binary..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(SERVER_LOCATION)

build-linux-arm64:
	@echo "Building Linux ARM binary..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(SERVER_LOCATION)

build-windows-amd64:
	@echo "Building Windows x64 binary..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(SERVER_LOCATION)

build-windows-arm64:
	@echo "Building Windows ARM binary..."
	@GOOS=windows GOARCH=arm64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe $(SERVER_LOCATION)

build-sharedlib-linux-amd64:
	@echo "Building shared library for Linux x64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(SHARED_LIB_OPT) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64.so $(SHARED_LIB_LOCATION)

generate-parser-nosql:
	pigeon -o ./parsers/nosql/nosql.go ./parsers/nosql/nosql.peg

test:
	@echo "Running unit tests..."
	@$(GOTEST) -v ./...

clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -rf $(DIST_DIR)

.PHONY: all test build-all build-macos build-linux clean generate-parser-nosql
