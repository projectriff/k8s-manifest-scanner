.PHONY: build clean install test all

SCANNER_OUTPUT = ./bin/k8s-manifest-scanner
RESOLVER_OUTPUT = ./bin/k8s-tag-resolver
GO_SOURCES = $(shell find . -type f -name '*.go')
VERSION ?= $(shell cat VERSION)

GOBIN ?= $(shell go env GOPATH)/bin

all: build test verify-goimports

test:
	GO111MODULE=on go test ./... -race -coverprofile=coverage.txt -covermode=atomic

check-goimports:
	@which goimports > /dev/null || (echo goimports not found: issue \"GO111MODULE=off go get golang.org/x/tools/cmd/goimports\" && false)

goimports: check-goimports
	@goimports -w pkg cmd

verify-goimports: check-goimports
	@goimports -l pkg cmd | (! grep .) || (echo above files are not formatted correctly. please run \"make goimports\" && false)

clean:
	rm -rf bin
	rm coverage.txt

install: build
	cp bin/* $(GOBIN)

build: $(SCANNER_OUTPUT) $(RESOLVER_OUTPUT)

$(SCANNER_OUTPUT): $(GO_SOURCES) go.mod go.sum VERSION
	GO111MODULE=on go build -o $(SCANNER_OUTPUT) ./cmd/k8s-manifest-scanner

$(RESOLVER_OUTPUT): $(GO_SOURCES) go.mod go.sum VERSION
	GO111MODULE=on go build -o $(RESOLVER_OUTPUT) ./cmd/k8s-tag-resolver
