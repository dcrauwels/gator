# gator
CLI RSS aggregator in Golang

# Requirements
- PostgreSQL version 14.15 or higher
- Goose version 3 or higher
- SQLC version 1.28.0 or higher
- Go version?

# Setup
## PostgreSQL stuff
## Goose migrations
Run the following from your `sql/schema` dir:
    goose postgres "connection_string" up

Where `"connection_string"` is the PostgreSQL connection string of the format `postgres://user:password@server:port/database`. With default settings this comes out to `postgres://postgres:postgres@localhost:5432/gator`.

## Install
Run the following from the dir containing `main.go`:
    go install

# Usage
    gator <command> <argument>

## Commands
    register: takes 1 argument (username). Adds a user with the provided username. Automatically sets current user to this user.
    login: takes 1 argument (username). Sets the current user to a previously registered user.
    reset: takes 0 arguments. Resets the users database (populated through register).
    users: takes 0 arguments. Lists users.
    addfeed: takes 1 argument (url). Adds an RSS feed to the database and follows it as the current user.
    feeds: takes 0 arguments. Lists all added RSS feeds.
    follow: takes 1 argument (url). Follows a previously added RSS feed as the current user.
    unfollow: takes 1 argument (url). Unfollows a previously followed RSS feed as the current user.
    following: takes 0 arguments. Lists all feeds followed by the current user.
    agg: takes 1 argument (time interval e.g. 10s, 1m, 1h). Aggregates posts from all added RSS feeds.
    browse: takes 0 or 1 argument (number of posts displayed, default 2). Shows posts from feeds followed by current user provided they have been aggregated through agg previously.
    resetposts: takes 0 arguments. Resets the posts database (populated through agg).