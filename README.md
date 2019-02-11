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

### Quick Start 

Start the server using the `vegeta-server` binary generated after the previous step.

```
Usage: main [<flags>]

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
      --ip="localhost"  Server IP Address.
      --port="8000"     Server Port.
  -v, --version         Version Info
      --debug           Enabled Debug
```

#### Example 
*Serve `HTTP` traffic at `localhost:8000/api/v1`*
```
./bin/vegeta-server --ip=localhost --port=8000 --debug
```

> **Bonus**
> 
> The `localhost:8000` example can be also started using the `run` target, which handles the build and starts the server.
> ```
> make run
> 
> INFO[0000] creating new dispatcher                       component=dispatcher
> INFO[0000] starting dispatcher                           component=dispatcher
> ```

### REST API Usage (`api/v1`)

#### Submit an attack - `POST api/v1/attack`

```
curl --header "Content-Type: application/json" --request POST --data '{"rate": 5,"duration": "3s","target":{"method": "GET","URL": "http://localhost:8000/api/v1/attack","scheme": "http"}}' http://localhost:8000/api/v1/attack
```
 
```
{
  "id": "494f98a2-7165-4d1b-8834-3226b49ab582",
  "status": "scheduled",
  "params": {
    "rate": 5,
    "duration": "3s",
    "target": {
      "method": "GET",
      "URL": "http://localhost:8000/api/v1/attack",
      "scheme": "http"
    }
  }
}
```
*The returned JSON body includes the **Attack ID** (`494f98a2-7165-4d1b-8834-3226b49ab582`) and the **Attack Status** (`scheduled`).*

#### View attack status by **Attack ID** - `GET api/v1/attack/<attackID>`

```
curl http://localhost:8000/api/v1/attack/494f98a2-7165-4d1b-8834-3226b49ab582
```

```
{
  "id": "494f98a2-7165-4d1b-8834-3226b49ab582",
  "status": "completed",
  "params": {
    "rate": 5,
    "duration": "3s",
    "target": {
      "method": "GET",
      "URL": "http://localhost:8000/api/v1/attack",
      "scheme": "http"
    }
  }
}
```

#### List all attacks `GET /api/v1/attack`

```
curl http://localhost:8000/api/v1/attack/
```

```
[
    {
        "id": "494f98a2-7165-4d1b-8834-3226b49ab582",
        "status": "completed",
        "params": {
            "rate": 5,
            "duration": "3s",
            "target": {
                "method": "GET",
                "URL": "http://localhost:8000/api/v1/attack",
                "scheme": "http"
            }
        }
    },
    {
        "id": "c6fbc450-434a-4082-86c0-2a00b09297cf",
        "status": "completed",
        "params": {
            "rate": 5,
            "duration": "1s",
            "target": {
                "method": "GET",
                "URL": "http://localhost:8000/api/v1/attack",
                "scheme": "http"
            }
        }
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
    "latencies": {
        "total": 44164990,
        "mean": 2944332,
        "max": 3394263,
        "50th": 2914967,
        "95th": 3391265,
        "99th": 3394263
    },
    "bytes_in": {
        "total": 0,
        "mean": 0
    },
    "bytes_out": {
        "total": 0,
        "mean": 0
    },
    "earliest": "2019-02-10T22:52:30.703235-05:00",
    "latest": "2019-02-10T22:52:33.50831-05:00",
    "end": "2019-02-10T22:52:33.511692272-05:00",
    "duration": 2805075000,
    "wait": 3382272,
    "requests": 15,
    "rate": 5.347450602925056,
    "success": 1,
    "status_codes": {
        "200": 15
    },
    "errors": []
}
```

#### List all attack reports - `GET api/v1/report`

```
curl http://localhost:8000/api/v1/report/
```

```
[
    {
        "latencies": {
            "total": 44164990,
            "mean": 2944332,
            "max": 3394263,
            "50th": 2914967,
            "95th": 3391265,
            "99th": 3394263
        },
        "bytes_in": {
            "total": 0,
            "mean": 0
        },
        "bytes_out": {
            "total": 0,
            "mean": 0
        },
        "earliest": "2019-02-10T22:52:30.703235-05:00",
        "latest": "2019-02-10T22:52:33.50831-05:00",
        "end": "2019-02-10T22:52:33.511692272-05:00",
        "duration": 2805075000,
        "wait": 3382272,
        "requests": 15,
        "rate": 5.347450602925056,
        "success": 1,
        "status_codes": {
            "200": 15
        },
        "errors": []
    },
    {
        "latencies": {
            "total": 14307169,
            "mean": 2861433,
            "max": 3409154,
            "50th": 3081794,
            "95th": 3409154,
            "99th": 3409154
        },
        "bytes_in": {
            "total": 0,
            "mean": 0
        },
        "bytes_out": {
            "total": 0,
            "mean": 0
        },
        "earliest": "2019-02-10T22:53:37.735724-05:00",
        "latest": "2019-02-10T22:53:38.537849-05:00",
        "end": "2019-02-10T22:53:38.540930794-05:00",
        "duration": 802125000,
        "wait": 3081794,
        "requests": 5,
        "rate": 6.233442418575659,
        "success": 1,
        "status_codes": {
            "200": 5
        },
        "errors": []
    }
]
```

### Running tests

Tests can be run using the `Makefile` target `test`

```make test```

## Contributing

Link to [CONTRIBUTING.md](https://github.com/nitishm/vegeta-server/blob/master/CONTRIBUTING.md)

## Road-map

Link to [road-map](https://github.com/nitishm/vegeta-server/projects/1)

## License

Link to [LICENSE](https://github.com/nitishm/vegeta-server/blob/master/LICENSE)

## Support

Contact Author at nitish.malhotra@gmail.com
