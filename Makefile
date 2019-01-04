
COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --exact-match --always)
DATE=$(shell date +'%FT%TZ%z')

GO  = GO111MODULE=on go
all: fmt lint build test

build: deps fmt
	CGO_ENABLED=0 ${GO} build -v -a -tags=netgo \
		-ldflags '-s -w -extldflags "-static" -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' -o bin/vegeta-server

clean:
	rm -f coverage.txt
	rm -rf bin

deps:
	${GO} mod vendor
	${GO} mod download

install:
	$(shell ./scripts/make-install.sh)

swagger:
	bin/swagger generate server --spec=spec/swagger.yaml --name=vegeta --exclude-main

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

.PHONY: all build clean deps install swagger test fmt validate lint ineffassign run
