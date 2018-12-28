GO=go
GOBUILD=$(GO) build

build-test-helper:
	mkdir -p bin
	$(GOBUILD) -o bin/apm-test-helper ./cmd/apm-test-helper

build:
	mkdir -p bin
	$(GOBUILD) -o bin/apm ./cmd/apm

run:
	./bin/apm daemon --debug true
test: build-test-helper
	go test -v ./...
