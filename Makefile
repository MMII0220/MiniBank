include .env
export

DB_URL = postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
MIGRATION_DIR = migration

MIGRATE = migrate -path $(MIGRATION_DIR) -database $(DB_URL)

NAME ?= create_transactions_table

# Steps to rollback (default: 1); override with N=3
N ?= 1

.PHONY: migrate-up migrate-down migrate-reset version migrate-create ensure-migrations-dir migrate-force migrate-status

migrate:
	migrate create -ext sql -dir $(MIGRATION_DIR) -digits 3 -seq $(NAME)

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down $(N)

# Rollback all migrations to version 0 (use with caution)
migrate-reset:
	$(MIGRATE) down -all

version:
	$(MIGRATE) version

# Alias for version
migrate-status: version

# Force-set migration version (use to clear dirty state)
# Usage: make migrate-force VERSION=1
migrate-force:
	$(if $(strip $(VERSION)),,$(error VERSION is required. Usage: make migrate-force VERSION=1))
	$(MIGRATE) force $(VERSION)
