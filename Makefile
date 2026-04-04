.PHONY: all build test test-race fmt vet tidy

all: fmt vet build test

build:
	go build ./...

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy
