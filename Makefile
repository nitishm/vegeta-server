COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --exact-match --always)
DATE=$(shell date +'%FT%TZ%z')

SERVER_DIR = cmd/server
GO  = GO111MODULE=on go


	

all: fmt lint build test

build: deps fmt
	CGO_ENABLED=0 ${GO} build -v -o bin/vegeta-server -a -tags=netgo \
		-ldflags '-s -w -extldflags "-static" -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)' ${SERVER_DIR}/main.go

clean:
	rm -f coverage.txt
	rm -rf bin
	rm -rf vendor

deps:
	${GO} mod vendor
	${GO} mod download

update-deps:
	${GO} mod verify
	${GO} mod tidy

install:
	$(shell ./scripts/make-install.sh)

test:
	${GO} test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	${GO} test -covermode=count -coverprofile=profile.cov ./...

fmt:
	${GO} fmt ./...

validate:
	golangci-lint run

lint:
	golint ./...

ineffassign:
	ineffassign .

run: build
	$(shell bin/vegeta-server --scheme=http --host=localhost --port=8000)


container:
	docker build -t vegeta-server:latest .

container_run: container
	@docker run -d -p 8000:80 --name vegeta vegeta-server:latest --rm

container_clean:
	@docker rm -f vegeta || true
	@docker image rm vegeta-server:latest || true

.PHONY: all build clean deps update-deps install test fmt validate lint ineffassign run
