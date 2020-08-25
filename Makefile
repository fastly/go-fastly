SHELL := /bin/bash -o pipefail

# List of tests to run
FILES ?= ./...

# List all our actual files, excluding vendor
GOPKGS ?= $(shell go list $(FILES) | grep -v /vendor/)
GOFILES ?= $(shell find . -name '*.go' | grep -v /vendor/)

# Tags specific for building
GOTAGS ?=

# Number of procs to use
GOMAXPROCS ?= 4

NAME := $(notdir $(shell pwd))

# Test Service ID
FASTLY_TEST_SERVICE_ID ?=
FASTLY_API_KEY ?=

all: mod-download dev-dependencies tidy fmt fiximports test vet staticcheck ## Runs all of the required cleaning and verification targets.
.PHONY: all

tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@go mod tidy
.PHONY: tidy

mod-download: ## Downloads the Go module.
	@echo "==> Downloading Go module"
	@go mod download
.PHONY: mod-download

dev-dependencies: ## Downloads the necessesary dev dependencies.
	@echo "==> Downloading development dependencies"
	@go install honnef.co/go/tools/cmd/staticcheck
	@go install golang.org/x/tools/cmd/goimports
.PHONY: dev-dependencies

test: ## Runs the test suite with VCR mocks enabled.
	@echo "==> Testing ${NAME}"
	@go test -timeout=30s -parallel=20 -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test

test-race: ## Runs the test suite with the -race flag to identify race conditions, if they exist.
	@echo "==> Testing ${NAME} (race)"
	@go test -timeout=60s -race -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test-race

test-full: ## Runs the tests with VCR disabled (i.e., makes external calls).
	@echo "==> Testing ${NAME} with VCR disabled"
	@VCR_DISABLE=1 \
		bash -c \
		'go test -timeout=60s -parallel=20 ${GOPKGS} ${TESTARGS}'
.PHONY: test-full

fix-fixtures: ## Updates test fixtures with a specified default service ID.
	@echo "==> Updating fixtures"
	@$(shell pwd)/scripts/fixFixtures.sh ${FASTLY_TEST_SERVICE_ID}
.PHONY: fix-fixtures

check-imports: ## A check which lists improperly-formatted imports, if they exist.
	@$(shell pwd)/scripts/check-imports.sh
.PHONY: check-imports

check-fmt: ## A check which lists improperly-formatted files, if they exist.
	@$(shell pwd)/scripts/check-gofmt.sh
.PHONY: check-fmt

check-mod: ## A check which lists extraneous dependencies, if they exist.
	@$(shell pwd)/scripts/check-mod.sh
.PHONY: check-mod

fiximports: ## Properly formats and orders imports.
	@echo "==> Fixing imports"
	@goimports -w {fastly,tools}
.PHONY: fiximports

fmt: ## Properly formats Go files and orders dependencies.
	@echo "==> Running gofmt"
	@gofmt -s -w ${GOFILES}
.PHONY: fmt

vet: ## Identifies common errors.
	@echo "==> Running go vet"
	@go vet ./...
.PHONY: vet

staticcheck: ## Runs the staticcheck linter.
	@echo "==> Running staticcheck"
	@staticcheck ./...
.PHONY: staticcheck

changelog: ## Generates the full project changelog.
	@$(shell pwd)/scripts/changelog.sh
.PHONY: changelog

release-changelog: ## Generates the changelog for a specific release.
	@$(shell pwd)/scripts/release-changelog.sh
.PHONY: release-changelog

.PHONY: help
help: ## Prints this help menu.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
