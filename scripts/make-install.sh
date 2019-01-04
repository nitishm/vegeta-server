#!/bin/sh

# Install swagger
download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url')
if [ ! -d "bin" ]; then
  mkdir bin
fi
curl -o bin/swagger -L'#' "$download_url"
chmod +x bin/swagger

# Install go binaries
GO111MODULE=on go get golang.org/x/lint/golint
GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint
GO111MODULE=on go get github.com/gordonklaus/ineffassign