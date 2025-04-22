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
DEFAULT_FASTLY_TEST_NGWAF_WORKSPACE_ID = alk6DTsYKHKucJCOIavaJM
FASTLY_API_KEY ?=
#
# Enables support for tools such as https://github.com/rakyll/gotest
TEST_COMMAND ?= $(GO) test

all: mod-download tidy fmt fiximports test vet staticcheck semgrep ## Runs all of the required cleaning and verification targets.
.PHONY: all

mod-download: ## Downloads the Go module.
	@echo "==> Downloading Go module"
	@$(GO) mod download
.PHONY: mod-download

tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@$(GO) mod tidy
.PHONY: tidy

fmt: ## Properly formats Go files and orders dependencies.
	@echo "==> Running gofmt"
	@gofmt -s -w fastly internal
.PHONY: fmt

fiximports: ## Properly formats and orders imports.
	@echo "==> Fixing imports"
	@$(GO) tool goimports -w fastly internal
.PHONY: fiximports

test: ## Runs the test suite with VCR mocks enabled.
	@echo "==> Testing ${NAME}"
	@$(TEST_COMMAND) -timeout=30s -parallel=20 -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test

vet: ## Identifies common errors.
	@echo "==> Running go vet"
	@$(GO) vet ./...
.PHONY: vet

staticcheck: ## Runs the staticcheck linter.
	@echo "==> Running staticcheck"
	@$(GO) tool staticcheck -version
	@$(GO) tool staticcheck ./...
.PHONY: staticcheck

semgrep: ## Run semgrep checker.
	@if [[ "$$(uname)" == 'Darwin' ]]; then brew install semgrep; fi
	@if command -v semgrep &> /dev/null; then semgrep ci --config auto --exclude-rule generic.secrets.security.detected-private-key.detected-private-key $(SEMGREP_ARGS); fi
.PHONY: semgrep

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

check-imports: ## A check which lists improperly-formatted imports, if they exist.
	@$(shell pwd)/scripts/check-imports.sh
.PHONY: check-imports

check-fmt: ## A check which lists improperly-formatted files, if they exist.
	@$(shell pwd)/scripts/check-gofmt.sh
.PHONY: check-fmt

check-mod: ## A check which lists extraneous dependencies, if they exist.
	@$(shell pwd)/scripts/check-mod.sh
.PHONY: check-mod

nilaway: ## Run nilaway
	@nilaway ./...

.PHONY: help
help: ## Prints this help menu.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
