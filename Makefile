SHELL := /bin/bash

OS := $(shell uname)
REAL_ARCHITECTURE := $(shell uname -m)
ARCHITECTURE := $(shell echo $(REAL_ARCHITECTURE) | sed "s/x86_64/amd64/g")
BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
CREATED := $(shell date +%Y-%m-%dT%T%z)
GIT_REPO := $(shell git config --get remote.origin.url)
GIT_TOKEN ?= $(shell cat git-token.txt)
REPO_NAME := $(shell basename ${GIT_REPO} .git)
SHORT_SHA := $(shell git rev-parse --short HEAD)
TAG_NAME := $(shell git describe --exact-match --tags 2> /dev/null)
REVISION ?= $(if ${TAG_NAME},${TAG_NAME},${SHORT_SHA})
VERSION_PATH := github.com/ingka-group-digital/${REPO_NAME}/internal/version

ACTIONLINT_VERSION := 1.6.27
ACTIONLINT_PATH := bin/actionlint_v$(ACTIONLINT_VERSION)/actionlint
ACTIONLINT_URL :=  https://github.com/rhysd/actionlint/releases/download/v$(ACTIONLINT_VERSION)/actionlint_$(ACTIONLINT_VERSION)_$(OS)_$(ARCHITECTURE).tar.gz

DOCKER_REPOSITORY := europe-docker.pkg.dev/${GCP_PROJECT_ID}/images
DOCKER_IMAGE := ${DOCKER_REPOSITORY}/${REPO_NAME}:${REVISION}

GOLANGCI_LINT_VERSION := v1.57.2
GOLANGCI_LINT := bin/golangci-lint_$(GOLANGCI_LINT_VERSION)/golangci-lint
GOLINTCI_URL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh

GOTESTSUM_VERSION := 1.8.0
GOTESTSUM_PATH := bin/gotestsum_v$(GOTESTSUM_VERSION)/gotestsum
GOTESTSUM_URL := https://github.com/gotestyourself/gotestsum/releases/download/v$(GOTESTSUM_VERSION)/gotestsum_$(GOTESTSUM_VERSION)_$(OS)_$(ARCHITECTURE).tar.gz

SPANCHECK_VERSION := v1.1.3
SPANCHECK_PATH := bin/spectre-go-linter-spancheck-$(SPANCHECK_VERSION)/spectre-go-linter-spancheck
SPANCHECK_URL := https://github.com/ingka-group-digital/spectre-go-linter-spancheck

SHELLCHECK_VERSION := 0.9.0
SHELLCHECK_PATH := bin/shellcheck-v$(SHELLCHECK_VERSION)/shellcheck
SHELLCHECK_URL := https://github.com/koalaman/shellcheck/releases/download/v$(SHELLCHECK_VERSION)/shellcheck-v$(SHELLCHECK_VERSION).$(OS).x86_64.tar.xz

TFLINT_VERSION := 0.47.0
TFLINT_PATH := bin/tflint_v$(TFLINT_VERSION)/tflint
TFLINT_URL :=  https://github.com/terraform-linters/tflint/releases/download/v$(TFLINT_VERSION)/tflint_$(OS)_$(ARCHITECTURE).zip

SPECTRAL_VERSION := 6.11.0
SPECTRAL_PATH := bin/spectral-v$(SPECTRAL_VERSION)/spectral-cli
SPECTRAL_URL := https://github.com/stoplightio/spectral/releases/download/v$(SPECTRAL_VERSION)/spectral-$(shell echo $(OS) | sed "s/Darwin/macos/g")-$(shell echo $(REAL_ARCHITECTURE) | sed "s/x86_64/x64/g")

DOC_SERVER_PORT := 8282

all: help

## docker: Build and publishes the docker image
docker: docker-build docker-publish docker-info

## docker-build: Build the docker file and inject the current git revision as the version. The DOCKERFILE_SUFFIX is used to specify a dockerfile with a suffix.
docker-build:
	@echo "🚀 Building docker image using Dockerfile$(if $(DOCKERFILE_SUFFIX),".$(DOCKERFILE_SUFFIX)",)"
	@docker build \
		-t "${DOCKER_IMAGE}$(if $(DOCKERFILE_SUFFIX),-$(DOCKERFILE_SUFFIX),)" \
		-f "Dockerfile$(if $(DOCKERFILE_SUFFIX),.$(DOCKERFILE_SUFFIX),)" \
		--build-arg CREATED=${CREATED} \
		--build-arg GIT_TOKEN=${GIT_TOKEN} \
		--build-arg REVISION=${REVISION} \
		--build-arg TAG=${REVISION} .
	@echo
	@echo "Built '${DOCKER_IMAGE}$(if $(DOCKERFILE_SUFFIX),-$(DOCKERFILE_SUFFIX),)'"

## docker-info: Returns the name of the image that will be generated
docker-info:
	@echo ${DOCKER_IMAGE}

## docker-publish: Published the docker image
docker-publish:
	@echo "🚀 Pushing docker image"
	@docker push ${DOCKER_IMAGE}

## clean: Clean up all build artifacts
clean:
	@echo "🚀 Cleaning up old artifacts"
	@rm -f bin/${REPO_NAME}

## build: Build the application artifacts. Linting can be skipped by setting env variable IGNORE_LINTING.
build: clean $(if $(IGNORE_LINTING), , lint)
	@echo "🚀 Building artifacts"
	@go build -ldflags="-s -w -X '${VERSION_PATH}.Version=${REVISION}' -X '${VERSION_PATH}.Commit=${SHORT_SHA}'" -o bin/${REPO_NAME} ./cmd/server

## run: Run the application
run: build
	@echo "🚀 Running binary"
	@./bin/${REPO_NAME}

## lint: Lint the source code
lint: ${GOLANGCI_LINT}
	@echo "🚀 Linting code"
	@$(GOLANGCI_LINT) run
	@$(MAKE) lint-spancheck
	@$(MAKE) lint-openapi-specs

## lint-conventional-commits: checks the commits to ensure they follow the conventional commits format
lint-conventional-commits:
	@echo "🚀 Checking commits"
	@source "scripts/conventional-commits/branch.sh"

## lint-spancheck: Check the source code for invalid span trace names.
lint-spancheck: ${SPANCHECK_PATH}
	@echo "🚀 Checking trace spans"
	@${SPANCHECK_PATH} ./... || { echo "❌ Invalid span(s). Run 'make lint-spancheck-fix' to fix them. "; exit 1; }

## lint-spancheck-fix: Fix all invalid span trace names.
lint-spancheck-fix: ${SPANCHECK_PATH}
	@echo "🚀 Fixing trace spans"
	@${SPANCHECK_PATH} -fix ./...

## lint-github-actions: Lint the GitHub actions code
lint-github-actions: ${ACTIONLINT_PATH} ${SHELLCHECK_PATH}
	@echo "🚀 Linting GitHub actions code"
	@$(ACTIONLINT_PATH) -shellcheck=${SHELLCHECK_PATH}

## lint-github-actions-info: Returns information about the current GitHub actions linter being used
lint-github-actions-info:
	@echo ${ACTIONLINT_PATH}

## lint-shellcheck-info: Returns information about the current Shellcheck linter being used
lint-shellcheck-info:
	@echo ${SHELLCHECK_PATH}

## lint-terraform: Lint Terraform code using tflint
lint-terraform: ${TFLINT_PATH}
	@echo "🚀 Linting Terraform code"
	@$(TFLINT_PATH) --init > /dev/null 2>&1
	@$(TFLINT_PATH) --recursive

## lint-terraform-info: Returns information about the current Terraform linter being used
lint-terraform-info:
	@echo ${TFLINT_PATH}

## lint-openapi-specs: Lint OpenAPI docs
lint-openapi-specs: ${SPECTRAL_PATH}
	@echo "🚀 Linting OpenAPI docs"
	@$(SPECTRAL_PATH) lint docs/openapi.yml

## test: Run Go tests
test: ${GOTESTSUM_PATH}
	@echo "🚀 Running tests"
	@set -o pipefail; ${GOTESTSUM_PATH} --format testname --no-color=false -- -race ./... | grep -v 'EMPTY'; exit $$?

## test-benchmark: Run Go benchmark tests
test-benchmark:
	@echo "🚀 Running benchmark tests"
	@go test -bench=. -benchmem ./...

## install-hooks: Install Git hooks
install-hooks:
	@echo "🚀 Installing Git hooks"
	@cp hooks/commit-msg .git/hooks/commit-msg
	@cp hooks/pre-push .git/hooks/pre-push

## uninstall-hooks: Uninstall Git hooks
uninstall-hooks:
	@echo "🚀 Uninstalling Git hooks"
	@rm -f .git/hooks/commit-msg
	@rm -f .git/hooks/pre-push

## update-hooks: Updates the installed hooks to the latest version
update-hooks: uninstall-hooks install-hooks

## coverage: Create a test coverage report in HTML format
coverage:
	@echo "🚀 Creating coverage report in HTML format"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

## doc: Start a local pkgserver that displays local documentation
doc: pkgsite
	@echo "📜 Starting local documentation server"
	@echo "Go to http://localhost:${DOC_SERVER_PORT}/github.com/ingka-group-digital/${REPO_NAME} to view documentation"
	@pkgsite -http localhost:${DOC_SERVER_PORT}

pkgsite:
ifeq (, $(shell which pkgsite))
	 @echo "📦 Installing pkgsite"
	 @go install golang.org/x/pkgsite/cmd/pkgsite@latest
endif

${ACTIONLINT_PATH}:
	@echo "📦 Installing actionlint ${ACTIONLINT_VERSION}"
	@mkdir -p $(dir ${ACTIONLINT_PATH})
	@curl -sSL ${ACTIONLINT_URL} > bin/actionlint.tar.gz
	@tar -xzf bin/actionlint.tar.gz -C $(patsubst %/,%,$(dir ${ACTIONLINT_PATH}))
	@rm -f bin/actionlint.tar.gz

${GOLANGCI_LINT}:
	@echo "📦 Installing golangci-lint ${GOLANGCI_LINT_VERSION}"
	@mkdir -p $(dir ${GOLANGCI_LINT})
	@curl -sfL ${GOLINTCI_URL} | sh -s -- -b ./$(patsubst %/,%,$(dir ${GOLANGCI_LINT})) ${GOLANGCI_LINT_VERSION} > /dev/null 2>&1

${GOTESTSUM_PATH}:
	@echo "📦 Installing GoTestSum ${GOTESTSUM_VERSION}"
	@mkdir -p $(dir ${GOTESTSUM_PATH})
	@curl -sSL ${GOTESTSUM_URL} > bin/gotestsum.tar.gz
	@tar -xzf bin/gotestsum.tar.gz -C $(patsubst %/,%,$(dir ${GOTESTSUM_PATH}))
	@rm -f bin/gotestsum.tar.gz

${SPANCHECK_PATH}:
	@echo "📦 Installing linter spancheck into ${SPANCHECK_PATH}"
	@mkdir -p $(dir ${SPANCHECK_PATH})
	@gh release download ${SPANCHECK_VERSION} --clobber --repo ${SPANCHECK_URL} --output bin/spancheck.tar.gz --pattern '*spancheck_${OS}_${REAL_ARCHITECTURE}.tar.gz'
	@tar -xzf bin/spancheck.tar.gz -C $(patsubst %/,%,$(dir ${SPANCHECK_PATH}))
	@rm -f bin/spancheck.tar.gz

${SHELLCHECK_PATH}:
	@echo "📦 Installing Shellcheck ${SHELLCHECK_VERSION}"
	@mkdir -p $(dir ${SHELLCHECK_PATH})
	@curl -sSL ${SHELLCHECK_URL} > bin/shellcheck.tar.xz
	@tar -xJf bin/shellcheck.tar.xz -C bin/
	@rm -f bin/shellcheck.tar.xz

${TFLINT_PATH}:
	@echo "📦 Installing tflint ${TFLINT_VERSION}"
	@mkdir -p $(dir ${TFLINT_PATH})
	@curl -sSL $(TFLINT_URL) > bin/tflint.zip
	@unzip -d $(dir ${TFLINT_PATH}) bin/tflint.zip
	@rm -f bin/tflint.zip

${SPECTRAL_PATH}:
	@echo "📦 Installing spectral ${SPECTRAL_VERSION} into '${SPECTRAL_PATH}'"
	@mkdir -p $(dir ${SPECTRAL_PATH})
	@echo "🌐 Downloading ${SPECTRAL_URL}"
	@curl -sSL $(SPECTRAL_URL) --output $(SPECTRAL_PATH)
	@chmod +x $(SPECTRAL_PATH)

help: Makefile
	@echo
	@echo "📗 Choose a command run in "${REPO_NAME}":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
