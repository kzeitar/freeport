.PHONY: build test fmt coverage release clean

BINARY_NAME=freeport
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

build:
	go build -o bin/$(BINARY_NAME) ./cmd/freeport

test:
	go test ./... -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

# Cross-platform build for release
release: clean
	@echo "Building release binaries..."
	@mkdir -p bin

	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/freeport
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/freeport

	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/freeport
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/freeport

	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/freeport

	@echo "Release binaries built in bin/"
	@ls -lh bin/$(BINARY_NAME)-*

clean:
	rm -rf bin/

