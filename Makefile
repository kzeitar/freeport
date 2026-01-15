.PHONY: build test fmt

build:
	go build -o bin/freeport ./cmd/freeport

test:
	go test ./...

fmt:
	go fmt ./...
