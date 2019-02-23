[![GoDoc](https://godoc.org/github.com/shazow/ssh-chat?status.svg)](https://godoc.org/github.com/nitishm/vegeta-server) 
[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/nitishm/vegeta-server) 
[![Build Status](https://travis-ci.org/shazow/ssh-chat.svg?branch=master)](https://travis-ci.org/nitishm/vegeta-server) 
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/nitishm/vegeta-server/blob/master/LICENSE) 
[![Coverage Status](https://coveralls.io/repos/github/nitishm/vegeta-server/badge.svg?branch=master)](https://coveralls.io/github/nitishm/vegeta-server?branch=master)
# Vegeta Server - A RESTful load-testing service

> NOTE: This is currently a work-in-progress. This means that features will be added/modified/removed at will. Do not use this in production.

A RESTful API server for [vegeta](https://github.com/tsenart/vegeta), a load testing tool written in [Go](https://github.com/golang/go).

Vegeta is a versatile HTTP load testing tool built out of a need to drill HTTP services with a constant request rate. The vegeta library is written in Go, which makes it ideal to implement server in Go.  

The REST API model enables users to, asynchronously, submit multiple `attack`s targeting the same (or varied) endpoints, lookup pending or historical `report`s and update/cancel/re-submit `attack`s, using a simple RESTful API.

## Getting Started

### Installing

```
make all
```

> NOTE: `make all` resolves all the dependencies, formats the code using `gofmt`, validates and lints using `golangci-lint` and `golint`, builds the `vegeta-server` binary and drops it in the `/bin` directory and finally runs tests using `go test`.

### Quick Start 

Start the server using the `vegeta-server` binary generated after the previous step.

```
Usage: main [<flags>]

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
      --ip="0.0.0.0"  Server IP Address.
      --port="80"     Server Port.
  -v, --version         Version Info
      --debug           Enabled Debug
```

#### Example 
*Serve `HTTP` traffic at `0.0.0.0:80/api/v1`*
```
./bin/vegeta-server --ip=0.0.0.0 --port=80 --debug
```

> Try it out using `make run`
> ```
> make run
> 
> INFO[0000] creating new dispatcher                       component=dispatcher
> INFO[0000] starting dispatcher                           component=dispatcher
> ```

### Using Docker
*Build the docker image using local Dockerfile*
```
docker build .
```
*Run the docker container*
```
docker run -d -p 8000:80 --name vegeta {container id}
```
*You can also build and run a docker container using make*
```
make container_run
```
> NOTE: `make container` and `make container_clean` can be used to build the Dockerfile and delete the container and image.

### Running tests

```
make test
```

## Contributing

Link to [CONTRIBUTING.md](https://github.com/nitishm/vegeta-server/blob/master/CONTRIBUTING.md)

### Project Structure
- `/`: Extraneous setup and configuration files. No go code should exist at this level.
- `/cmd/server`: Comprises of `package main` serving as an entry point to the code.

- `/models`: Includes the model definitions used by the DB and the API endpoints.
    - `/db.go`: Provides the storage interface, which is implemented by the configured database.

- `/internal`: Internal only packages used by the server to run attacks and serve reports.
    - `/dispatcher`: Defines and implements the dispatcher interface, with the primary responsibility to carry out concurrent attacks.
    - `/reporter`: Defines and implements the reporter interface, with the primary responsibility to generate reports from previously completed attacks, in supported formats (JSON/Text/Binary).
    - `/endpoints`: Responsible for defining and registering the REST API endpoint handlers.

- `/pkg/vegeta`: [Vegeta library](https://github.com/tsenart/vegeta/tree/master/lib)  specific, wrapper methods and definitions. (*Keep these isolated from the internals of the server, to support more load-testing tools/libraries in the future.*)

- `/scripts`: Helper installation scripts.

## Road-map

Link to [road-map](https://github.com/nitishm/vegeta-server/projects/1)

## License

Link to [LICENSE](https://github.com/nitishm/vegeta-server/blob/master/LICENSE)

## Support

Contact Author at nitish.malhotra@gmail.com
