SHELL := /bin/bash

REAL_ARCHITECTURE := $(shell uname -m)
ARCHITECTURE := $(shell echo $(REAL_ARCHITECTURE) | sed "s/x86_64/amd64/g")
CREATED := $(shell date +%Y-%m-%dT%T%z)
GCP_PROJECT_ID ?= my-gcp-project-id
GIT_REPO := $(shell git config --get remote.origin.url)
OS := $(shell uname)
REPO_NAME := $(shell basename ${GIT_REPO} .git)
SHORT_SHA := $(shell git rev-parse --short HEAD)
TAG_NAME := $(shell git describe --exact-match --tags 2> /dev/null)
REVISION ?= $(if ${TAG_NAME},${TAG_NAME},${SHORT_SHA})

DOCKER_REPOSITORY := europe-docker.pkg.dev/${GCP_PROJECT_ID}/images
DOCKER_IMAGE := ${DOCKER_REPOSITORY}/${REPO_NAME}:${REVISION}

GOTESTSUM_VERSION := 1.12.0
GOTESTSUM_PATH := bin/gotestsum_v$(GOTESTSUM_VERSION)/gotestsum
GOTESTSUM_URL := https://github.com/gotestyourself/gotestsum/releases/download/v$(GOTESTSUM_VERSION)/gotestsum_$(GOTESTSUM_VERSION)_$(OS)_$(ARCHITECTURE).tar.gz

GOLANGCI_LINT_VERSION := v1.61.0
GOLANGCI_LINT := bin/golangci-lint_$(GOLANGCI_LINT_VERSION)/golangci-lint
GOLINTCI_URL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh

all: help

## clean: Clean up all build artifacts
clean:
	@echo "🚀 Cleaning up old artifacts"
	@rm -f bin/${REPO_NAME}

## build: Build the application artifacts. Linting can be skipped by setting env variable IGNORE_LINTING.
build: clean lint
	@echo "🚀 Building artifacts"
	@go build -ldflags="-s -w" -o bin/${REPO_NAME} ./cmd/server

## run: Run the application
run: build
	@echo "🚀 Running binary"
	@./bin/${REPO_NAME}

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
	@protoc --go_out=api/gen/. --go-grpc_out=api/gen/. api/proto/customer.proto

## docker: Build and publishes the docker image
docker: docker-build docker-publish docker-info

## docker-build: Build the docker file.
docker-build:
	@echo "🚀 Building docker image using Dockerfile"
	@docker build \
		-t "${DOCKER_IMAGE}" \
		-f "Dockerfile" \
		--build-arg CREATED=${CREATED} \
		--build-arg REVISION=${REVISION} \
		--build-arg TAG=${REVISION} .
	@echo
	@echo "Built '${DOCKER_IMAGE}"

## docker-info: Returns the name of the image that will be generated
docker-info:
	@echo ${DOCKER_IMAGE}

## docker-publish: Published the docker image
docker-publish:
	@echo "🚀 Pushing docker image"
	## @docker push ${DOCKER_IMAGE} skip until there is a real project

## test: Run Go tests
test: ${GOTESTSUM_PATH}
	@echo "🚀 Running tests"
	@set -o pipefail; ${GOTESTSUM_PATH} --format testname --no-color=false -- -race ./... | grep -v 'EMPTY'; exit $$?

${GOTESTSUM_PATH}:
	@echo "📦 Installing GoTestSum ${GOTESTSUM_VERSION}"
	@mkdir -p $(dir ${GOTESTSUM_PATH})
	@curl -sSL ${GOTESTSUM_URL} > bin/gotestsum.tar.gz
	@tar -xzf bin/gotestsum.tar.gz -C $(patsubst %/,%,$(dir ${GOTESTSUM_PATH}))
	@rm -f bin/gotestsum.tar.gz

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