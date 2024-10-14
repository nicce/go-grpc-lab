SHELL := /bin/bash

ARCHITECTURE := $(shell echo $(REAL_ARCHITECTURE) | sed "s/x86_64/amd64/g")
CREATED := $(shell date +%Y-%m-%dT%T%z)
OS := $(shell uname)
REAL_ARCHITECTURE := $(shell uname -m)
REVISION ?= $(if ${TAG_NAME},${TAG_NAME},${SHORT_SHA})

GOLANGCI_LINT_VERSION := v1.61.0
GOLANGCI_LINT := bin/golangci-lint_$(GOLANGCI_LINT_VERSION)/golangci-lint
GOLINTCI_URL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh

all: help

## lint: Lint the source code
lint: ${GOLANGCI_LINT}
	@echo "🚀 Linting code"
	@$(GOLANGCI_LINT) run

lint-fix: ${GOLANGCI_LINT}
	@echo "🚀 Linting code and fixing issues"
	@$(GOLANGCI_LINT) run --fix

## test-benchmark: Run Go benchmark tests
test-benchmark:
	@echo "🚀 Running benchmark tests"
	@go test -bench=. -benchmem ./...

## generate-proto: Generate server and client code from proto files
generate-proto: generate-proto-server-code

## generate-proto-server-code: Generate server code from proto files
generate-proto-server-code:
	@echo "🚀 Generating server code from proto files"
	@protoc --go_out=api/. --go-grpc_out=api/. proto/customer.proto

${GOLANGCI_LINT}:
	@echo "📦 Installing golangci-lint ${GOLANGCI_LINT_VERSION}"
	@mkdir -p $(dir ${GOLANGCI_LINT})
	@curl -sfL ${GOLINTCI_URL} | sh -s -- -b ./$(patsubst %/,%,$(dir ${GOLANGCI_LINT})) ${GOLANGCI_LINT_VERSION} > /dev/null 2>&1

help: Makefile
	@echo
	@echo "📗 Choose a command run in "${REPO_NAME}":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo