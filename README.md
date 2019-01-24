[![GoDoc](https://godoc.org/github.com/shazow/ssh-chat?status.svg)](https://godoc.org/github.com/nitishm/vegeta-server/pkg/vegeta) 
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

### Usage

Start the server using the `vegeta-server` binary generated after the previous step.

The `vegeta-server` supports `flags` to pass configuration options to the server such as , **scheme** [`http` / `https`], **host**, **port** and `TLS` configurations, among others.

```
Usage:
  main [OPTIONS]

This is a RESTful API for the vegeta load-testing utility. Vegeta is a versatile HTTP load testing tool built out of a need to drill HTTP services with a constant
request rate.


Application Options:
      --scheme=            the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=   grace period for which to wait before killing idle connections (default: 10s)
      --graceful-timeout=  grace period for which to wait before shutting down the server (default: 15s)
      --max-header-size=   controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It
                           does not limit the size of the request body. (default: 1MiB)
      --socket-path=       the unix socket to listen on (default: /var/run/vegeta.sock)
      --host=              the IP to listen on (default: localhost) [$HOST]
      --port=              the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=      limit the number of outstanding requests
      --keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default:
                           3m)
      --read-timeout=      maximum duration before timing out read of the request (default: 30s)
      --write-timeout=     maximum duration before timing out write of the response (default: 60s)
      --tls-host=          the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=          the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=   the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=           the private key to use for secure conections [$TLS_PRIVATE_KEY]
      --tls-ca=            the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=  limit the number of outstanding requests
      --tls-keep-alive=    sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=  maximum duration before timing out read of the request
      --tls-write-timeout= maximum duration before timing out write of the response

Version:
      --version            Show vegeta-server version details

Help Options:
  -h, --help               Show this help message
```

#### Example 
*Serve `HTTP` traffic at `localhost:8000/api/v1`*
```
./bin/vegeta-server --scheme=http --host=localhost --port=8000
```

> **Bonus**
> 
> The `localhost:8000` example can be also started using the `run` target, which handles the build and starts the server.
> ```
> make run
> 
> INFO[0000] Serving vegeta at http://127.0.0.1:8000
> ```

### REST API Usage (`api/v1`)

#### Submit an attack - `POST api/v1/attack`

```
curl --header "Content-Type: application/json" --request POST --data '{"rate": 5,"duration": "3s","target":{"method": "GET","URL": "http://localhost:8000/api/v1/attack","scheme": "http"}}' http://localhost:8000/api/v1/attack
```
 
```
{
    "id":"d9788d4c-1bd7-48e9-92e4-f8d53603a483",
    "status":"scheduled"
}
```
*The returned JSON body includes the **Attack ID** (`d9788d4c-1bd7-48e9-92e4-f8d53603a483`) and the **Attack Status** (`scheduled`).*

#### View attack status by **Attack ID** - `GET api/v1/attack/<attackID>`

```
curl http://localhost:8000/api/v1/attack/d9788d4c-1bd7-48e9-92e4-f8d53603a483
```

```
{
    "id": "d9788d4c-1bd7-48e9-92e4-f8d53603a483",
    "status": "completed"
}
```

#### List all attacks `GET /api/v1/attack`

```
curl http://localhost:8000/api/v1/attack/
```

```
[
    {
        "id": "d9788d4c-1bd7-48e9-92e4-f8d53603a483",
        "status": "completed"
    },
    {
        "id": "8300c02f-4836-4458-b0a9-1493d8a32409",
        "status": "completed"
    }
]
```

#### View attack report by **Attack ID** - `GET /api/v1/report/<attackID>`

> The report endpoint only returns results for **Completed** attacks


```
curl http://localhost:8000/api/v1/report/d9788d4c-1bd7-48e9-92e4-f8d53603a483
```

```
{
    "id": "d9788d4c-1bd7-48e9-92e4-f8d53603a483",
    "report": {
        "bytes_in": {
            "mean": 67,
            "total": 1005
        },
        "bytes_out": {},
        "duration": 2802988000,
        "earliest": "2019-01-21T13:40:50.441-05:00",
        "end": "2019-01-21T13:40:53.244-05:00",
        "errors": [],
        "latencies": {
            "50th": 297520,
            "95th": 991371,
            "99th": 1211546,
            "max": 1211546,
            "mean": 350855,
            "total": 5262839
        },
        "latest": "2019-01-21T13:40:53.244-05:00",
        "rate": 5.3514321145862915,
        "requests": 15,
        "success": 1,
        "wait": 315104
    }
}
```

#### List all attack reports - `GET api/v1/report`

```
curl http://localhost:8000/api/v1/report/
```

```
[
    {
        "id": "d9788d4c-1bd7-48e9-92e4-f8d53603a483",
        "report": {
            "bytes_in": {
                "mean": 67,
                "total": 1005
            },
            "bytes_out": {},
            "duration": 2802988000,
            "earliest": "2019-01-21T13:40:50.441-05:00",
            "end": "2019-01-21T13:40:53.244-05:00",
            "errors": [],
            "latencies": {
                "50th": 297520,
                "95th": 991371,
                "99th": 1211546,
                "max": 1211546,
                "mean": 350855,
                "total": 5262839
            },
            "latest": "2019-01-21T13:40:53.244-05:00",
            "rate": 5.3514321145862915,
            "requests": 15,
            "success": 1,
            "wait": 315104
        }
    },
    {
        "id": "8300c02f-4836-4458-b0a9-1493d8a32409",
        "report": {
            "bytes_in": {
                "mean": 433,
                "total": 216500
            },
            "bytes_out": {},
            "duration": 9983230000,
            "earliest": "2019-01-21T13:47:19.597-05:00",
            "end": "2019-01-21T13:47:29.581-05:00",
            "errors": [],
            "latencies": {
                "50th": 328599,
                "95th": 546664,
                "99th": 1148068,
                "max": 3619612,
                "mean": 368702,
                "total": 184351113
            },
            "latest": "2019-01-21T13:47:29.580-05:00",
            "rate": 50.08399085265991,
            "requests": 500,
            "success": 1,
            "wait": 425382
        }
    }
]
```

### Running tests

Tests can be run using the `Makefile` target `test`

```make test```

## Benchmark

*TODO*

## Guides & Tutorials

*TODO*

## Contributing

Link to [CONTRIBUTING.md](https://github.com/nitishm/vegeta-server/blob/master/CONTRIBUTING.md)

---

### Swagger - API Specification

The API's [swagger](https://swagger.io/) specification is formatted using the [OpenAPI 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md) specification. The spec can be found at [`vegeta-server/spec/swagger.yml`](https://github.com/nitishm/vegeta-server/tree/master/spec/swagger.yaml)

The server code is generated using [go-swagger](https://github.com/go-swagger/go-swagger), a tool, written in Go, that implements the OpenAPI 2.0 specification.

---

### Generate the server code

To generate the **server** code using the `go-swagger` CLI tool, use `make swagger`. 
> NOTE: Install the `go-swagger` tool binary perform using `make install`)

```
make swagger

bin/swagger generate server --spec=spec/swagger.yaml --name=vegeta --exclude-main
2019/01/05 17:57:04 validating spec /Users/nitishm/vegeta-server/spec/swagger.yaml
(...truncated for brevity...)
2019/01/05 17:57:05 Generation completed!

For this generation to compile you need to have some packages in your GOPATH:

	* github.com/go-openapi/runtime
	* github.com/jessevdk/go-flags

You can get these now with: go get -u -f ./...
```

---

### Code Structure 
**Generated Packages (_DO NOT EDIT_)**
1. [`/restapi`](https://github.com/nitishm/vegeta-server/tree/master/restapi) (except for [`configure_vegeta.go`](https://github.com/nitishm/vegeta-server/blob/master/restapi/configure_vegeta.go)) : Generated server code and specification object.
2.  [`/restapi/operations`](https://github.com/nitishm/vegeta-server/tree/master/restapi/operations) : Generated API params/responses/handlers.
3. [`/models`](https://github.com/nitishm/vegeta-server/tree/master/models) : Generated models for the `swagger` *`definitions`* component.

**Editable files & packages**
1. [`/internal`](https://github.com/nitishm/vegeta-server/tree/master/internal) :
  Internal packages that cannot be exported as part of the package.
2. [`/pkg`](https://github.com/nitishm/vegeta-server/tree/master/pkg) :
   Exportable packages used across the project
3. [`configure_vegeta.go`](https://github.com/nitishm/vegeta-server/blob/master/restapi/configure_vegeta.go) :
  API handlers that utilize the `internal` and `pkg` packages.
  
## Road-map

Link to [backend](https://github.com/nitishm/vegeta-server/projects/1)
Link to [frontend](https://github.com/nitishm/vegeta-server/projects/2)
Link to [documentation](https://github.com/nitishm/vegeta-server/projects/3)

## License

Link to [LICENSE](https://github.com/nitishm/vegeta-server/blob/master/LICENSE)

## Support

Contact Author at nitish.malhotra@gmail.com
