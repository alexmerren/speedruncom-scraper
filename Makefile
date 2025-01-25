GO ?= go
LINTER := golangci-lint
GOFMT := gofmt

DIST_DIR := $(CURDIR)/dist
CMD_DIR := $(CURDIR)/cmd
VENDOR_DIR := $(CURDIR)/vendor
DATA_DIR := $(CURDIR)/data

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

## all: Download dependencies, format code, build binary.
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

## clean: Remove vendored code and built executables
.PHONY: clean
clean:
	rm -r $(DIST_DIR)
	rm -r $(VENDOR_DIR)

## run: Run all executables in required order
.PHONY: run
run:
	$(DIST_DIR)/games-list && \
	$(DIST_DIR)/games-data && \
	$(DIST_DIR)/leaderboards-data && \
	$(DIST_DIR)/users-list && \
	$(DIST_DIR)/users-data && \
	$(DIST_DIR)/runs-data

## monitor: Monitor the progress of `make run`
.PHONY: monitor
monitor:
	@sh -c "find $(DATA_DIR)/v1 $(DATA_DIR)/v2 -type f -exec wc -l {} \;"
