GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

BINARY_NAME=cosmium
SERVER_LOCATION=./cmd/server

SHARED_LIB_LOCATION=./sharedlibrary
SHARED_LIB_OPT=-buildmode=c-shared
XGO_TARGETS=linux/amd64,linux/arm64,windows/amd64,windows/arm64,darwin/amd64,darwin/arm64
GOVERSION=1.22.0

DIST_DIR=dist

SHARED_LIB_TEST_CC=gcc
SHARED_LIB_TEST_CFLAGS=-Wall -ldl
SHARED_LIB_TEST_TARGET=$(DIST_DIR)/sharedlibrary_test
SHARED_LIB_TEST_DIR=./sharedlibrary/tests
SHARED_LIB_TEST_SOURCES=$(wildcard $(SHARED_LIB_TEST_DIR)/*.c)

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

build-sharedlib-tests: build-sharedlib-linux-amd64
	@echo "Building shared library tests..."
	@$(SHARED_LIB_TEST_CC) $(SHARED_LIB_TEST_CFLAGS) -o $(SHARED_LIB_TEST_TARGET) $(SHARED_LIB_TEST_SOURCES)

run-sharedlib-tests: build-sharedlib-tests
	@echo "Running shared library tests..."
	@$(SHARED_LIB_TEST_TARGET) $(DIST_DIR)/$(BINARY_NAME)-linux-amd64.so

xgo-compile-sharedlib:
	@echo "Building shared libraries using xgo..."
	@mkdir -p $(DIST_DIR)
	@xgo -targets=$(XGO_TARGETS) -go $(GOVERSION) -buildmode=c-shared -dest=$(DIST_DIR) -out=$(BINARY_NAME) -pkg=$(SHARED_LIB_LOCATION) .

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
