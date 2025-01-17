#! /bin/bash

goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator" down
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator" up