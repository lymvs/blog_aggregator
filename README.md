# Gator
This is an RSS feed aggregator in Go

## Requirements
- Go 1.25.1 (or later)
- Postgres

## Installation
```
go install github.com/lymvs/blog_aggregator
psql -c "CREATE DATABASE gator;"
```

## Config
To use gator, you must first create a .gatorconfig.json within your home directory with the following minimal config (url needs tweaking to work with your configuration of postgres)

```
{
    "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable"
}
```