GO ?= go
LINTER := golangci-lint
GOFMT := gofmt

DIST_DIR := $(CURDIR)/dist
CMD_DIR := $(CURDIR)/cmd

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
all: vendor fmt build

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

## fmt: Format all code for the project
.PHONY: fmt
fmt: 
	$(GOFMT) -s -w $(CURDIR)
