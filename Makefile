.PHONY: build clean test all

OUTPUT = ./bin/k8s-manifest-scanner
GO_SOURCES = $(shell find . -type f -name '*.go')
VERSION ?= $(shell cat VERSION)

GOBIN ?= $(shell go env GOPATH)/bin

all: build test verify-goimports

test:
	GO111MODULE=on go test ./... -race -coverprofile=coverage.txt -covermode=atomic

check-goimports:
	@which goimports > /dev/null || (echo goimports not found: issue \"GO111MODULE=off go get golang.org/x/tools/cmd/goimports\" && false)

goimports: check-goimports
	@goimports -w cmd

verify-goimports: check-goimports
	@goimports -l pkg main.go | (! grep .) || (echo above files are not formatted correctly. please run \"make goimports\" && false)

install: build
	cp $(OUTPUT) $(GOBIN)

build: $(GO_SOURCES) VERSION
	GO111MODULE=on go build -o $(OUTPUT) ./cmd/k8s-manifest-scanner
