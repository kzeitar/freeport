.PHONY: build test fmt coverage

build:
	go build -o bin/freeport ./cmd/freeport

test:
	go test ./... -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...
