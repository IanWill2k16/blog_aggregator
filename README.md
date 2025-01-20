# Gator

Gator is a simple blog aggregator that doesn't actually work all that well.

## Description

An in-depth paragraph about your project and overview of use.

## Getting Started

### Prerequisites

* Go version 1.22.3

### Installing

* Postgres installation with Linux:
```
sudo apt update
sudo apt install postgresql postgresql-contrib

sudo passwd postgres

sudo -u postgres psql

CREATE DATABASE gator;

\c gator

ALTER USER postgres PASSWORD 'postgres';
```
* Goose installation:
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
* Create a ~/.gatorconfig.json file that contains:
```
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": "kahya"
}
```

### Executing program

* Register your user account
```
go run . register <user>
```
* Add a blog feed
```
go run . addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
```
* Start the aggreagation loop in a dedicated terminal
``` 
go run . agg 30s
```
* Retrieve entries
```
go run . browse 2
```
