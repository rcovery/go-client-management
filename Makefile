#!make
include .env

setup:
	docker compose up -d
	cp .env.example .env
	go mod tidy
	go build -o app
	make migrateup

build:
	go build -o app

start:
	./app

dev:
	go run .

migration:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:$(DBPORT)/$(DBDATABASE)?sslmode=$(DBSSLMODE)"
	goose create -dir internal/infra/postgres/migrations $$NAME sql 

migrateup:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:$(DBPORT)/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose up -dir internal/infra/postgres/migrations

migratedown:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:$(DBPORT)/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose down -dir internal/infra/postgres/migrations

