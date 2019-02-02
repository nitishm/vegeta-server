#!/bin/sh

# Install go binaries
GO111MODULE=on go get golang.org/x/lint/golint
GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint
GO111MODULE=on go get github.com/gordonklaus/ineffassign