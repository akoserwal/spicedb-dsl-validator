# Requires golangci-lint to be installed @ $(go env GOPATH)/bin/golangci-lint
# https://golangci-lint.run/usage/install/

GO := go
GOFMT := gofmt
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell $(GO) env GOBIN))
GOBIN=$(shell $(GO) env GOPATH)/bin
else
GOBIN=$(shell $(GO) env GOBIN)
endif

DOCKER ?= podman
DOCKER_CONFIG="${PWD}/.podman"

PKG         := ./cmd/...
TAGS        :=
TESTS       := .
TESTFLAGS   := -race -v -failfast
LDFLAGS     := -w -s
GOFLAGS     :=
CGO_ENABLED ?= 0

lint:
	golangci-lint run ./...
.PHONY: lint

binary:
	$(GO) build
.PHONY: binary

binary/linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o spicedb-dsl-validator-linux-amd64
.PHONY: binary/linux

test:
	$(GO) clean -testcache && go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)
.PHONY: test

