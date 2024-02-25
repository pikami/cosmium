GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

BINARY_NAME=cosmium

DIST_DIR=dist

all: test build-all

build-all: build-darwin-arm64 build-darwin-amd64 build-linux-amd64 build-windows-amd64

build-darwin-arm64:
	@echo "Building macOS ARM binary..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-darwin-amd64:
	@echo "Building macOS x64 binary..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .

build-linux-amd64:
	@echo "Building Linux x64 binary..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .

build-windows-amd64:
	@echo "Building Windows x64 binary..."
    @GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .

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
