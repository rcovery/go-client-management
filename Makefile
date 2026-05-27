#!make
include .env

migration:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:$(DBPORT)/$(DBDATABASE)?sslmode=$(DBSSLMODE)"
	goose create -dir internal/infra/postgres/migrations $$NAME sql 

migrateup:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:$(DBPORT)/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose up -dir internal/infra/postgres/migrations

migratedown:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:$(DBPORT)/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose down -dir internal/infra/postgres/migrations

