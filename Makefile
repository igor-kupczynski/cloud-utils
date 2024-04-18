BINARY_NAME=cloud-utils

.PHONY: all
all: install

.PHONY: build
build:
	go build -o $(BINARY_NAME) main.go

.PHONY: run
run:
	go run main.go

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: format
format:
	go fmt ./...
