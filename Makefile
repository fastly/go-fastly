SHELL := /bin/bash -o pipefail
GO ?= go

# List of tests to run
FILES ?= ./...

# List all our actual files, excluding vendor
GOPKGS ?= $(shell $(GO) list $(FILES) | grep -v /vendor/)

# Tags specific for building
GOTAGS ?=

# Number of procs to use
GOMAXPROCS ?= 4

NAME := $(notdir $(shell pwd))

# Test Resource IDs
FASTLY_TEST_DELIVERY_SERVICE_ID ?=
DEFAULT_FASTLY_TEST_DELIVERY_SERVICE_ID = kKJb5bOFI47uHeBVluGfX1
FASTLY_TEST_COMPUTE_SERVICE_ID ?=
DEFAULT_FASTLY_TEST_COMPUTE_SERVICE_ID = XsjdElScZGjmfCcTwsYRC1
FASTLY_TEST_NGWAF_WORKSPACE_ID ?=
DEFAULT_FASTLY_TEST_NGWAF_WORKSPACE_ID = Am2qjXkgamuYp3u54rQkLD
FASTLY_API_KEY ?=

# Enables support for tools such as https://github.com/rakyll/gotest
TEST_COMMAND ?= $(GO) test

# Tooling versions
GOLANGCI_LINT_VERSION = v2.4.0
BIN_DIR := $(CURDIR)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

all: mod-download tidy fmt lint test semgrep ## Runs all of the required cleaning and verification targets.
.PHONY: all

mod-download: ## Downloads the Go module.
	@echo "==> Downloading Go module"
	@$(GO) mod download
.PHONY: mod-download

tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@$(GO) mod tidy
.PHONY: tidy

install-linter: ## Installs golangci-lint v2.4.0 into ./bin (cross-platform)
	@echo "==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)"
	@mkdir -p $(BIN_DIR)
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		OS=$$(uname | tr '[:upper:]' '[:lower:]'); \
		ARCH=$$(uname -m); \
		case "$$ARCH" in \
			x86_64) ARCH="amd64" ;; \
			arm64|aarch64) ARCH="arm64" ;; \
			*) echo "Unsupported architecture: $$ARCH" && exit 1 ;; \
		esac; \
		URL="https://github.com/golangci/golangci-lint/releases/download/$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-$$OS-$$ARCH.tar.gz"; \
		echo "Downloading: $$URL"; \
		curl -sSfL "$$URL" | tar -xz -C $(BIN_DIR) --strip-components=1; \
	fi
.PHONY: install-linter

check-linter-version: ## Verifies installed golangci-lint version matches expected
	@echo "==> Checking golangci-lint version"
	@EXPECTED="$(GOLANGCI_LINT_VERSION)"; \
	INSTALLED=$$($(GOLANGCI_LINT) version --format short | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+'); \
	if [ "$$INSTALLED" != "$$EXPECTED" ]; then \
		echo "Expected golangci-lint $$EXPECTED but found $$INSTALLED"; \
		exit 1; \
	fi
.PHONY: check-linter-version

fmt: install-linter ## Properly formats Go files and orders dependencies.
	@echo "==> Running golangci-lint fmt"
	@$(GOLANGCI_LINT) fmt
.PHONY: fmt

lint: install-linter check-linter-version ## Runs golangci lint
	@echo "==> Running golangci-lint"
	@$(GOLANGCI_LINT) run
.PHONY: lint

semgrep: ## Run semgrep checker.
	@if [[ "$$(uname)" == 'Darwin' ]]; then brew install semgrep; fi
	@if command -v semgrep &> /dev/null; then semgrep ci --config auto --exclude-rule generic.secrets.security.detected-private-key.detected-private-key $(SEMGREP_ARGS); fi
.PHONY: semgrep

test: ## Runs the test suite with VCR mocks enabled.
	@echo "==> Testing ${NAME}"
	@$(TEST_COMMAND) -timeout=30s -parallel=20 -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test

test-race: ## Runs the test suite with the -race flag to identify race conditions, if they exist.
	@echo "==> Testing ${NAME} (race)"
	@$(TEST_COMMAND) -timeout=60s -race -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test-race

test-full: ## Runs the tests with VCR disabled (i.e., makes external calls).
	@echo "==> Testing ${NAME} with VCR disabled"
	@VCR_DISABLE=1 \
		bash -c \
		'${GO} test -timeout=60s -parallel=20 ${GOPKGS} ${TESTARGS}'
.PHONY: test-full

fix-delivery-fixtures: ## Updates test fixtures with a specified default Delivery service ID.
	@echo "==> Updating fixtures"
	@$(shell pwd)/scripts/fixFixtures.sh ${FASTLY_TEST_DELIVERY_SERVICE_ID} ${DEFAULT_FASTLY_TEST_DELIVERY_SERVICE_ID}
.PHONY: fix-delivery-fixtures

fix-compute-fixtures: ## Updates test fixtures with a specified default Compute service ID.
	@echo "==> Updating fixtures"
	@$(shell pwd)/scripts/fixFixtures.sh ${FASTLY_TEST_COMPUTE_SERVICE_ID} ${DEFAULT_FASTLY_TEST_COMPUTE_SERVICE_ID}
.PHONY: fix-compute-fixtures

fix-ngwaf-fixtures: ## Updates test fixtures with a specified default Next-Gen WAF workspace ID.
	@echo "==> Updating fixtures"
	@$(shell pwd)/scripts/fixFixtures.sh ${FASTLY_TEST_NGWAF_WORKSPACE_ID} ${DEFAULT_FASTLY_TEST_NGWAF_WORKSPACE_ID}
.PHONY: fix-ngwaf-fixtures

nilaway: ## Run nilaway
	@nilaway ./...
.PHONY: nilaway

clean-bin: ## Removes locally installed binaries
	@echo "==> Cleaning ./bin directory"
	@rm -rf $(BIN_DIR)
.PHONY: clean-bin

help: ## Prints this help menu.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help
