# Metadata about this makefile and position
MKFILE_PATH := $(lastword $(MAKEFILE_LIST))
CURRENT_DIR := $(patsubst %/,%,$(dir $(realpath $(MKFILE_PATH))))

# Ensure GOPATH
GOPATH ?= $(HOME)/go

# List all our actual files, excluding vendor
GOPKGS ?= $(shell go list $(FILES) | grep -v /vendor/)
GOFILES ?= $(shell find . -name '*.go' | grep -v /vendor/)

# Tags specific for building
GOTAGS ?=

# Number of procs to use
GOMAXPROCS ?= 4

PROJECT := $(CURRENT_DIR:$(GOPATH)/src/%=%)
OWNER := $(notdir $(patsubst %/,%,$(dir $(PROJECT))))
NAME := $(notdir $(PROJECT))
EXTERNAL_TOOLS = \
	github.com/golang/dep/cmd/dep

# Current system information
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOCACHE ?= $(shell go env GOCACHE)

# List of tests to run
FILES ?= ./...

# Test Service ID
FASTLY_TEST_SERVICE_ID ?=

# bootstrap installs the necessary go tools for development or build.
bootstrap:
	@echo "==> Bootstrapping ${PROJECT}"
	@for t in ${EXTERNAL_TOOLS}; do \
		echo "--> Installing $$t" ; \
		go get -u "$$t"; \
	done
.PHONY: bootstrap

clean:
	@echo "==> Cleaning ${NAME}"
	@rm -rf pkg
.PHONY: clean

# build builds the binary into pkg/
build:
	@echo "==> Building ${NAME} for ${GOOS}/${GOARCH}"
	@env \
		-i \
		PATH="${PATH}" \
		CGO_ENABLED="0" \
		GOOS="${GOOS}" \
		GOARCH="${GOARCH}" \
		GOCACHE="${GOCACHE}" \
		GOPATH="${GOPATH}" \
		go build -a -o "pkg/${GOOS}_${GOARCH}/${NAME}" ${GOPKGS}
.PHONY: build

# deps updates all dependencies for this project.
deps:
	@echo "==> Updating deps for ${PROJECT}"
	@dep ensure -update
	@dep prune
.PHONY: deps

# dev builds and installs the
dev:
	@echo "==> Installing ${NAME} for ${GOOS}/${GOARCH}"
	@env \
		-i \
		PATH="${PATH}" \
		CGO_ENABLED="0" \
		GOOS="${GOOS}" \
		GOARCH="${GOARCH}" \
		GOCACHE="${GOCACHE}" \
		GOPATH="${GOPATH}" \
		go install ${GOPKGS}
.PHONY: dev

# linux builds the linux binary
linux:
	@env \
		GOOS="linux" \
		GOARCH="amd64" \
		$(MAKE) -f "${MKFILE_PATH}" build
.PHONY: linux

# test runs the test suite.
test:
	@echo "==> Testing ${NAME}"
	@go test -timeout=30s -parallel=20 -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test

# test-race runs the test suite.
test-race:
	@echo "==> Testing ${NAME} (race)"
	@go test -timeout=60s -race -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test-race

# test without VCR
test-full:
	@echo "==> Testing ${NAME} with VCR disabled"
	@env \
		VCR_DISABLE=1 \
		go test -timeout=60s -parallel=20 ${GOPKGS} ${TESTARGS}
.PHONY: test-full

# update fixtures default service ID
fix-fixtures:
	@echo "==> Updating fixtures"
	@$(CURRENT_DIR)/scripts/fixFixtures.sh ${FASTLY_TEST_SERVICE_ID}
.PHONY: fix-fixtures

# lists improperly-formatted files if they exist
check-fmt:
	@gofmt -s -l ${GOFILES}
	@printf "\nPlease, run 'make fmt'.\n"
.PHONY: check-fmt

# properly formats go files
fmt:
	@gofmt -s -w ${GOFILES}

# runs go vet
vet:
	@go vet ${GOFILES}
.PHONY: vet

# runs the staticcheck linter
staticcheck:
	@staticcheck ${GOFILES}
.PHONY: staticcheck

# generates the full project changelog
changelog:
	@$(CURRENT_DIR)/scripts/changelog.sh
.PHONY: changelog

# generates the changelog for a specific release
release-changelog:
	@$(CURRENT_DIR)/scripts/release-changelog.sh
.PHONY: release-changelog
