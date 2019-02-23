### REST API Usage (`api/v1`)

#### Submit an attack - `POST api/v1/attack`

```
curl --header "Content-Type: application/json" --request POST --data '{"rate": 5,"duration": "3s","target":{"method": "GET","URL": "http://0.0.0.0:80/api/v1/attack","scheme": "http"}}' http://0.0.0.0:80/api/v1/attack
```
 
```json
{
  "id": "494f98a2-7165-4d1b-8834-3226b49ab582",
  "status": "scheduled",
  "params": {
    "rate": 5,
    "duration": "3s",
    "target": {
      "method": "GET",
      "URL": "http://0.0.0.0:80/api/v1/attack",
      "scheme": "http"
    }
  },
  "created_at": "Mon, 18 Feb 2019 19:48:19 EST",
  "updated_at": "Mon, 18 Feb 2019 19:48:33 EST"
}
```
*The returned JSON body includes the **Attack ID** (`494f98a2-7165-4d1b-8834-3226b49ab582`) and the **Attack Status** (`scheduled`).*

#### View attack status by **Attack ID** - `GET api/v1/attack/<attackID>`

```
curl http://0.0.0.0:80/api/v1/attack/494f98a2-7165-4d1b-8834-3226b49ab582
```

```json
{
  "id": "494f98a2-7165-4d1b-8834-3226b49ab582",
  "status": "completed",
  "params": {
    "rate": 5,
    "duration": "3s",
    "target": {
      "method": "GET",
      "URL": "http://0.0.0.0:80/api/v1/attack",
      "scheme": "http"
    }
  },
  "created_at": "Mon, 18 Feb 2019 19:48:19 EST",
  "updated_at": "Mon, 18 Feb 2019 19:48:33 EST"

}
```

#### List all attacks `GET /api/v1/attack`

```
curl http://0.0.0.0:80/api/v1/attack/
```

```json
[
    {
        "id": "494f98a2-7165-4d1b-8834-3226b49ab582",
        "status": "completed",
        "params": {
            "rate": 5,
            "duration": "3s",
            "target": {
                "method": "GET",
                "URL": "http://0.0.0.0:80/api/v1/attack",
                "scheme": "http"
            }
        },
        "created_at": "Mon, 18 Feb 2019 19:48:19 EST",
        "updated_at": "Mon, 18 Feb 2019 19:48:33 EST"
    },
    {
        "id": "c6fbc450-434a-4082-86c0-2a00b09297cf",
        "status": "completed",
        "params": {
            "rate": 5,
            "duration": "1s",
            "target": {
                "method": "GET",
                "URL": "http://0.0.0.0:80/api/v1/attack",
                "scheme": "http"
            }
        },
        "created_at": "Mon, 18 Feb 2019 19:48:19 EST",
        "updated_at": "Mon, 18 Feb 2019 19:48:33 EST"
    }
]
```

#### View attack report by **Attack ID** - `GET /api/v1/report/<attackID>[?format=json/text/binary]`

> The report endpoint only returns results for **Completed** attacks

- *JSON Format*
```
curl http://0.0.0.0:80/api/v1/report/d9788d4c-1bd7-48e9-92e4-f8d53603a483?format=json
```

```json
{
    "id": "d9788d4c-1bd7-48e9-92e4-f8d53603a483",
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

- *Text Format*
```
curl http://0.0.0.0:80/api/v1/report/9aea25c6-3dcf-4f14-808f-5e499d1d0074?format=text
```

```text
Requests      [total, rate]            200, 100.47
Duration      [total, attack, wait]    1.993288918s, 1.990719s, 2.569918ms
Latencies     [mean, 50, 95, 99, max]  2.136603ms, 1.642011ms, 4.151042ms, 9.884504ms, 15.338328ms
Bytes In      [total, mean]            0, 0.00
Bytes Out     [total, mean]            0, 0.00
Success       [ratio]                  0.00%
Status Codes  [code:count]             404:200  
Error Set:
404 Not Found
```
#### List all attack reports - `GET api/v1/report`

```
curl http://0.0.0.0:80/api/v1/report/
```

```json
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
