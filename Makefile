APP_NAME := gosmi-json-exporter
CMD_DIR := ./cmd/exporter
OUT_DIR := ./bin
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*")

.PHONY: all build run test clean fmt vet lint

all: build

build:
	@echo ">> Building..."
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(APP_NAME) $(CMD_DIR)

run:
	@$(OUT_DIR)/$(APP_NAME) --help

test:
	@echo ">> Running tests..."
	go test -v ./internal/...

clean:
	@echo ">> Cleaning up..."
	rm -rf $(OUT_DIR)

fmt:
	@echo ">> Formatting..."
	go fmt ./...

vet:
	@echo ">> Vetting..."
	go vet ./...

lint:
	@echo ">> Linting (requires golangci-lint)..."
	golangci-lint run ./...
