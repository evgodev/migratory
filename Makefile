BIN := "./bin/migratory"
DOCKER_IMG = "migratory:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/migratory/main.go

run: build
	$(BIN)

lint:
	golangci-lint run ./... -v

test:
	go test -race ./internal/... -cover

.PHONY: build run lint test