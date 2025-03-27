# gator
CLI RSS aggregator in Golang

# Requirements
- PostgreSQL version 14.15 or higher
- Goose version 3 or higher
- SQLC version 1.28.0 or higher

# Setup
## PostgreSQL stuff
## Goose migrations
Run the following from your `sql/schema` dir:
    goose postgres "connection_string" up

Where `"connection_string"` is the PostgreSQL connection string of the format `postgres://user:password@server:port/database`. With default settings this comes out to `postgres://postgres:postgres@localhost:5432/gator`.