#!/bin/sh
set -e

if [ ! -f .env ]; then
  cp .env.example .env
fi

export $(grep -v '^#' .env | xargs)

go install github.com/pressly/goose/v3/cmd/goose@latest

GOOSE_DBSTRING="postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBDATABASE}?sslmode=${DBSSLMODE}" \
  goose up -dir internal/infra/postgres/migrations

exec go run .
