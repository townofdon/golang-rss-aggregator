# Go - RSS Feed Tutorial

Source:

[FreeCodeCamp - Go Programming – Golang Course with Bonus Projects](https://www.youtube.com/watch?v=un6ZyFkqFKo&t=22176s)


## Setup

```
go mod tidy
go mod vendor
go build
./tutorial-go-rss-server
```

## Test

```
# GET request
curl -i localhost:8000/v1/healthz

# POST request
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST -i http://localhost:8000/v1/something
```
