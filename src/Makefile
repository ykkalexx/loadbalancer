# Binary name
BINARY_NAME=loadbalancer

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Git related variables
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean test coverage lint run docker

all: clean build

## Build:
build:
    @echo "Building..."
    go build -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" -o $(GOBIN)/$(BINARY_NAME) ./cmd/loadbalancer

## Run:
run:
    go run ./cmd/loadbalancer

## Clean:
clean:
    @echo "Cleaning..."
    @rm -rf $(GOBIN)
    @go clean -cache

## Test:
test:
    go test -v ./...

## Coverage:
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

## Lint:
lint:
    @echo "Running linter..."
    golangci-lint run

## Docker:
docker:
    docker build -t $(BINARY_NAME) .

## Install dependencies:
deps:
    go mod download
    go mod tidy

## Help:
help:
    @echo "Make commands:"
    @echo "build    - Build the application"
    @echo "run      - Run the application"
    @echo "clean    - Clean build files"
    @echo "test     - Run tests"
    @echo "coverage - Generate test coverage"
    @echo "lint     - Run linter"
    @echo "docker   - Build Docker image"
    @echo "deps     - Install dependencies"