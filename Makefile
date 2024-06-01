GO ?= go
DOCKER ?= docker
LINTER := golangci-lint
GOFMT := gofmt

DIST_DIR := $(CURDIR)/dist
INTERNAL_DIR := $(CURDIR)/internal
CMD_DIR := $(CURDIR)/cmd
PKG_DIR := $(CURDIR)/pkg

TEST_MODULES := $(shell $(GO) list $(INTERNAL_DIR)/... $(PKG_DIR)/...)

GOFLAGS :=

## help: Print this message
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'

## all: Download dependencies and build executables
all: vendor build

## build: Create the binary 
build:
	@mkdir -p $(DIST_DIR)
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -mod=vendor -o $(DIST_DIR) $(CMD_DIR)/...

## vendor: Download the vendored dependencies 
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

## deps: Remove, upgrade, then vendor dependencies
deps:
	$(GO) get -u $(CURDIR)/...
	$(GO) mod tidy
	$(GO) mod vendor

## lint: Lint the project 
lint:
	$(LINTER) run

## test: Run the unit tests for the project 
test:
	@$(GO) test $(TEST_MODULES) -coverprofile=$(CURDIR)/coverage.out coverpkg=$(INTERNAL_DIR)
	@$(GO) tool cover -html=$(CURDIR)/coverage.out -o $(CURDIR)/test-coverage.html
	@$(GO) tool cover -func=$(CURDIR)/coverage.out \
		| awk '$$1 == "total:" {printf("Total coverage: %.2f%% of statements\n", $$3)}'

## fmt: Format all code for the project
fmt: 
	$(GOFMT) -s -w $(CURDIR)

.PHONY: help all build vendor lint test fmt deps