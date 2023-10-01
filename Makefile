GO ?= go
DOCKER ?= docker
LINTER := golangci-lint
GOFMT := gofmt

DIST_DIR := $(CURDIR)/dist
INTERNAL_DIR := $(CURDIR)/internal
CMD_DIR := $(CURDIR)/cmd

TEST_MODULES := $(shell $(GO) list $(INTERNAL_DIR)/...)

GOFLAGS :=
# Set to 1 to use static linking for all builds (including tests).
STATIC :=

ifeq ($(STATIC),1)
LDFLAGS += -s -w -extldflags "-static"
endif

## help: Print this message
.PHONY: help
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'

## all: Download dependencies, generate mocks, fmt, run unit tests, build binary.
.PHONY: all
all: vendor fmt test build

## build: Create the binary 
.PHONY: build
build:
	@mkdir -p $(DIST_DIR)
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -mod=vendor -o $(DIST_DIR) $(CMD_DIR)/...

## vendor: Download the vendored dependencies 
.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

## lint: Lint the project 
.PHONY: lint
lint:
	$(LINTER) run

## test: Run the unit tests for the project 
.PHONY: test
test:
	@$(GO) test $(TEST_MODULES) -coverprofile=$(CURDIR)/coverage.out coverpkg=$(INTERNAL_DIR)
	@$(GO) tool cover -html=$(CURDIR)/coverage.out -o $(CURDIR)/test-coverage.html
	@$(GO) tool cover -func=$(CURDIR)/coverage.out \
		| awk '$$1 == "total:" {printf("Total coverage: %.2f%% of statements\n", $$3)}'

## fmt: Format all code for the project
.PHONY: fmt
fmt: 
	$(GOFMT) -s -w $(CURDIR)
