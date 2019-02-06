COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --exact-match --always)
DATE=$(shell date +'%FT%TZ%z')

SERVER_DIR = cmd/server
GO  = GO111MODULE=on go

all: fmt lint build test

build: deps fmt
	CGO_ENABLED=0 ${GO} build -v -o bin/vegeta-server -a -tags=netgo \
		-ldflags '-s -w -extldflags "-static" -X vegeta-server/restapi.version=$(VERSION) -X vegeta-server/restapi.commit=$(COMMIT) -X vegeta-server/restapi.date=$(DATE)' ${SERVER_DIR}/main.go

clean:
	rm -f coverage.txt
	rm -rf bin

deps:
	${GO} mod vendor
	${GO} mod download

update-deps:
	${GO} mod verify
	${GO} mod tidy
	rm -rf vendor
	${GO} mod vendor

install:
	$(shell ./scripts/make-install.sh)

test:
	${GO} test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

fmt:
	${GO} fmt ./...

validate:
	golangci-lint run

lint:
	golint

ineffassign:
	ineffassign .

run: build
	$(shell bin/vegeta-server --scheme=http --host=localhost --port=8000)

.PHONY: all build clean deps update-deps install test fmt validate lint ineffassign run
