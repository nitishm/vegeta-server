# REST API Usage (`api/v1`)

## Submit an attack - `POST api/v1/attack`

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

### With Request Body

The request body is passed along as **[base64](https://en.wikipedia.org/wiki/Base64)** encoded string, generated from the JSON request body.

Example

- **Original JSON Request Body**
```json
{
	"rate": 1,
	"duration": "5s",
	"target": {
		"method": "POST",
		"URL": "http://localhost:80/api/v1/attack",
		"scheme": "http"
	}
}
```

- **Convert to base64**
```
$ echo '{
        "rate": 1,
        "duration": "5s",
        "target": {
                "method": "POST",
                "URL": "http://localhost:80/api/v1/attack",
                "scheme": "http"
        }
}' | base64
ewoJInJhdGUiOiAxLAoJImR1cmF0aW9uIjogIjVzIiwKCSJ0YXJnZXQiOiB7CgkJIm1ldGhvZCI6ICJQT1NUIiwKCQkiVVJMIjogImh0dHA6Ly9sb2NhbGhvc3Q6ODAvYXBpL3YxL2F0dGFjayIsCgkJInNjaGVtZSI6ICJodHRwIgoJfQp9Cg==
```

- **Submit Attack**

```
curl --header "Content-Type: application/json" --request POST --data '{"rate": 5, "duration": "10s", "target": {"method": "POST", "URL": "http://localhost:80/api/v1/attack", "scheme": "http"}, "body": "ewoJInJhdGUiOiAxLAoJImR1cmF0aW9uIjogIjVzIiwKCSJ0YXJnZXQiOiB7CgkJIm1ldGhvZCI6ICJQT1NUIiwKCQkiVVJMIjogImh0dHA6Ly9sb2NhbGhvc3Q6ODAvYXBpL3YxL2F0dGFjayIsCgkJInNjaGVtZSI6ICJodHRwIgoJfQp9Cg=="}' http://0.0.0.0:80/api/v1/attack
```
```json
{
  "id": "443101cb-ded8-4e39-aa6b-c745516d1ca7",
  "status": "scheduled",
  "params": {
    "rate": 5,
    "duration": "10s",
    "body": "ewoJInJhdGUiOiAxLAoJImR1cmF0aW9uIjogIjVzIiwKCSJ0YXJnZXQiOiB7CgkJIm1ldGhvZCI6ICJQT1NUIiwKCQkiVVJMIjogImh0dHA6Ly9sb2NhbGhvc3Q6ODAvYXBpL3YxL2F0dGFjayIsCgkJInNjaGVtZSI6ICJodHRwIgoJfQp9Cg==",
    "target": {
      "method": "POST",
      "URL": "http://localhost:80/api/v1/attack",
      "scheme": "http"
    }
  },
  "created_at": "Sun, 03 Mar 2019 20:55:12 EST",
  "updated_at": "Sun, 03 Mar 2019 20:55:12 EST"
}
```

## Cancel an attack by **Attack ID** - `POST api/v1/attack/<attackID>/cancel`

> SUCCESS - Returns Status Code 200 OK

```
curl --header "Content-Type: application/json" --request POST --data '{"cancel": true}' http://0.0.0.0:80/api/v1/attack/5ebdfe2a-5c98-4cd9-a9ce-a1af89f20d53/cancel
```

## View attack status by **Attack ID** - `GET api/v1/attack/<attackID>`

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

## List all attacks `GET /api/v1/attack[?{parameters}]`

Availables parameters :
* status : `scheduled | running | canceled | completed | failed`
* created_before : `YYYY-mm-dd`
* created_after : `YYYY-mm-dd`

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

## View attack report by **Attack ID** - `GET /api/v1/report/<attackID>[?format=json/text/binary/histogram]`

> The report endpoint only returns results for **Completed** attacks

### JSON Format

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

### Text Format

```
curl http://0.0.0.0:80/api/v1/report/9aea25c6-3dcf-4f14-808f-5e499d1d0074?format=text
```

```text
Id 9aea25c6-3dcf-4f14-808f-5e499d1d0074
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

### Histogram Format `Default`

```
curl http://0.0.0.0/api/v1/report/b39cf62a-0141-4919-a9e0-38a007e59d8f?format=histogram
```

```text
ID b39cf62a-0141-4919-a9e0-38a007e59d8f
Bucket           #   %        Histogram
[0s,     500ms]  0   0.00%    
[500ms,  1s]     0   0.00%    
[1s,     1.5s]   0   0.00%    
[1.5s,   2s]     0   0.00%    
[2s,     2.5s]   0   0.00%    
[2.5s,   3s]     0   0.00%    
[3s,     +Inf]   15  100.00%  ###########################################################################
```

### Histogram Format

```
curl http://0.0.0.0/api/v1/report/b39cf62a-0141-4919-a9e0-38a007e59d8f?format=histogram&bucket=0,2s,4s,6s,8s
```

```text
ID b39cf62a-0141-4919-a9e0-38a007e59d8f
Bucket         #   %       Histogram
[0s,    2s]    0   0.00%   
[2s,    4s]    3   20.00%  ###############
[4s,    6s]    10  66.67%  ##################################################
[6s,    8s]    2   13.33%  ##########
[8s,    +Inf]  0   0.00%   
```

## List all attack reports - `GET api/v1/report`

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
