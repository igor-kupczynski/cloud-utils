BINARY_NAME=cloud-utils

all: build

build:
	go build -o $(BINARY_NAME) main.go

run:
	go run main.go

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test ./...

lint:
	golint ./...

format:
	go fmt ./...

.PHONY: all build run clean test lint format
