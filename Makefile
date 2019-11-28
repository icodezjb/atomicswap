# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

SOURCEPKG=github.com/icodezjb/atomicswap
ASWAP_BINARY=aswap
ASWAP_ADMIN_BINARY=aswap-admin

ASWAP_GITCOMMIT=$(shell git rev-parse --short HEAD 2> /dev/null || true)
ASWAP_BUILDTIME=$(shell date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')

# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
EXIST_LINT := $(shell ls $(GOPATH)/bin/golangci-lint 2> /dev/null)

.PHONY:all
all: test build

.PHONY:build
build:
	$(GOBUILD) -ldflags "\
               -X $(SOURCEPKG)/cmd.Commit=$(ASWAP_GITCOMMIT) \
               -X $(SOURCEPKG)/cmd.BuildTime=$(ASWAP_BUILDTIME) \
               " -v -o build/bin/$(ASWAP_BINARY) $(SOURCEPKG)/cmd/$(ASWAP_BINARY)

	$(GOBUILD) -ldflags "\
               -X $(SOURCEPKG)/cmd.Commit=$(ASWAP_GITCOMMIT) \
               -X $(SOURCEPKG)/cmd.BuildTime=$(ASWAP_BUILDTIME) \
               " -v -o build/bin/$(ASWAP_ADMIN_BINARY) $(SOURCEPKG)/cmd/$(ASWAP_ADMIN_BINARY)

.PHONY:test
test:
	$(GOTEST) -v ./contract/...

.PHONY:clean
clean:
	$(GOCLEAN)
	rm -f build/bin/$(ASWAP_BINARY)
	rm -f build/bin/$(ASWAP_ADMIN_BINARY)

.PHONY:fmt
fmt:
	$(GOCMD) list -f {{.Dir}} ./... | xargs gofmt -w -s -d

.PHONY:lint
lint:
# fmt project
	@ $(GOCMD) list -f {{.Dir}} ./... | xargs gofmt -w -s -d

ifndef EXIST_LINT
	@ echo "waiting several minutes for installing golangci-lint"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.18.0
	@ echo ""
endif

# more info about `GOGC` env: https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint
	@GOGC=50 $(GOPATH)/bin/golangci-lint run \
	        --enable=goimports \
	        --enable=varcheck \
	        --enable=vet \
	        --enable=gofmt \
	        --enable=misspell \
            --max-same-issues=100