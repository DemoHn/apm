GO=go
GOBUILD=$(GO) build

build-test-helper:
	mkdir -p bin
	$(GOBUILD) -o bin/apm-test-helper ./cmd/apm-test-helper

build:
	mkdir -p bin
	$(GOBUILD) -o bin/apm ./cmd/apm/main.go

test: build-test-helper
	go test ./...