# Go - RSS Feed Aggregator Tutorial

Source:

[FreeCodeCamp - Go Programming – Golang Course with Bonus Projects](https://www.youtube.com/watch?v=un6ZyFkqFKo&t=22176s)


## Setup

```
go mod tidy
go mod vendor
go build
./tutorial-go-rss-server

# PG Server
docker-compose up
```

## Migrations

Run these from `sql/schema` dir.

To migrate:

```
./bin/migrate.sh
```

Or, manually:

```
goose postgres $DB_URL up
# e.g.
goose postgres postgres://superuser:tacotuesdays@localhost:5432/golang_rss up
```

To rollback:

```
goose postgres $DB_URL down
```

## Generate DB Queries

Run these from project root.

```
sqlc generate
```

## Test

```
# GET request
curl -i localhost:8000/v1/healthz

# POST request
curl -d '{"key1":"value1", "key2":"value2"}' -H "Content-Type: application/json" -X POST -i http://localhost:8000/v1/something
```

## Deps / Tools

- [sqlc](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html)([install](https://docs.sqlc.dev/en/latest/overview/install.html))
- [goose](https://github.com/pressly/goose)
- [PGAdmin](https://www.pgadmin.org/)
