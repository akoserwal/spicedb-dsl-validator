# Requires golangci-lint to be installed @ $(go env GOPATH)/bin/golangci-lint
# https://golangci-lint.run/usage/install/

GO := go
GOFMT := gofmt

VERSION=$(shell git describe --tags --always --long --dirty)
EXECUTABLE=spicedb-dsl-validator
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64

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
	$(GO) build -o -ldflags="-s -w -X main.version=$(VERSION)"
.PHONY: binary

binary/darwin:
	env GOOS=darwin GOARCH=amd64 $(GO) build -o $(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)"
.PHONY: binary

binary/linux:
	env GOOS=linux GOARCH=amd64 $(GO) build -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"
.PHONY: binary/linux

binary/windows:
	env GOOS=windows GOARCH=amd64 $(GO) build -o $(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)"
.PHONY: binary/windows

test:
	$(GO) clean -testcache && go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)
.PHONY: test

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)
