ifeq ($(POSTGRES_SETUP_PROD),)
	POSTGRES_SETUP_PROD := user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=$(DB_SETUP_HOST) port=$(DB_PORT) sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(NAME)" sql

.PHONY: prod-migration-up
prod-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_PROD)" up

.PHONY: prod-migration-down
prod-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_PROD)" down

.PHONY: run
run:
	make prod-migration-up
	bin/api
