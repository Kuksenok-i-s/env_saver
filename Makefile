BINARY_NAME="env_saver"
MAIN_PATH=cmd/env_saver/main.go

.PHONY: all build clean build-docker migrate

all: build

build:
	go build -ldflags="-X main.commitID=`git rev-parse HEAD`" -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	go clean
	rm -f $(BINARY_NAME)

build-docker:
	docker build -t $(BINARY_NAME) .

test:
	go test ./...

run:
	go run $(MAIN_PATH)

deps:
	go mod download

lint:
	golangci-lint run

.DEFAULT_GOAL := build
