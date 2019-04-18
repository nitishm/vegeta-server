COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --exact-match --always)
DATE=$(shell date +'%FT%TZ%z')
CONTAINER_NAME ?= vegeta

SERVER_DIR = cmd/server

all: fmt lint build test

build: deps fmt
	CGO_ENABLED=0 go build -v -o bin/vegeta-server -a -tags=netgo \
		-ldflags '-s -w -extldflags "-static" -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)' ${SERVER_DIR}/main.go

clean:
	rm -f coverage.txt
	rm -rf bin
	rm -rf vendor

deps:
	go mod vendor
	go mod download

update-deps:
	go mod verify
	go mod tidy

install:
	$(shell ./scripts/make-install.sh)

test:
	go clean -testcache ./...
	go test -v -race -covermode=atomic ./...

	go clean -testcache ./...
	go test -v -covermode=count -coverprofile=profile.cov ./...

fmt:
	go fmt ./...

validate:
	golangci-lint run

lint:
	golint cmd/... internal/... models/... pkg/...
	go vet ${SERVER_DIR}/main.go

ineffassign:
	ineffassign .

run: build
	$(shell bin/vegeta-server --scheme=http --host=localhost --port=8000)

container:
	docker build -t vegeta-server:latest .

container_stop:
	@docker rm -f '$(CONTAINER_NAME)' || true

container_run: container
	@docker run --rm -d -p 8000:80 --name '$(CONTAINER_NAME)' vegeta-server:latest

container_clean: container_stop
	@docker image rm vegeta-server:latest || true

.PHONY: all build clean deps update-deps install test fmt validate lint ineffassign run
