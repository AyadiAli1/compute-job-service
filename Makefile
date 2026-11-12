# Einheitlicher Einstiegspunkt für lokale Entwicklung und CI.
BINARY  := bin/server
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.DEFAULT_GOAL := help
.PHONY: help fmt lint test build run clean


proto: ## Formatiert, lintet und generiert Code aus den .proto-Dateien
	buf format -w
	buf lint
	buf generate 

help: ## Zeigt diese Hilfe
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-8s\033[0m %s\n", $$1, $$2}'

fmt: ## Formatiert den Go-Code
	golangci-lint fmt

lint: ## Führt alle Linter aus
	golangci-lint run ./...

test: ## Go-Tests mit Race-Detector und Coverage
	go test -race -cover ./...

build: ## Baut das Server-Binary nach bin/
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/server

run: build ## Baut und startet den Server
	./$(BINARY)

clean: ## Entfernt Build-Artefakte
	rm -rf bin/