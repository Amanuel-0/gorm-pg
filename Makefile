

ENV ?= dev

.PHONY: help create-migration migrate-up migrate-down migrate-force

DATABASE_DSN="mysql://amanuel:password@tcp(gormpg-db:3306)/gormpg"
MIGRATIONS_DIR=./backend/internals/db/migrations


up:
	@set -a; [ -f .env.$(ENV) ] && . .env.$(ENV); set +a; \
	COMPOSE_PROJECT_ENV=$(ENV) docker-compose \
		-f docker-compose.yml \
		-f docker-compose.$(ENV).yml up --build

up-detached:
	@set -a; [ -f .env.$(ENV) ] && . .env.$(ENV); set +a; \
	COMPOSE_PROJECT_ENV=$(ENV) docker-compose \
		-f docker-compose.yml \
		-f docker-compose.$(ENV).yml up --build -d

stop:
	@set -a; [ -f .env.$(ENV) ] && . .env.$(ENV); set +a; \
	COMPOSE_PROJECT_ENV=$(ENV) docker-compose -f docker-compose.yml -f docker-compose.$(ENV).yml stop

down:
	@set -a; [ -f .env.$(ENV) ] && . .env.$(ENV); set +a; \
	COMPOSE_PROJECT_ENV=$(ENV) docker-compose -f docker-compose.yml -f docker-compose.$(ENV).yml down -v

logs:
	@set -a; [ -f .env.$(ENV) ] && . .env.$(ENV); set +a; \
	COMPOSE_PROJECT_ENV=$(ENV) docker-compose -f docker-compose.yml -f docker-compose.$(ENV).yml logs -f

rebuild:
	@set -a; [ -f .env.$(ENV) ] && . .env.$(ENV); set +a; \
	COMPOSE_PROJECT_ENV=$(ENV) docker-compose -f docker-compose.yml -f docker-compose.$(ENV).yml build --no-cache

help: ## Show help
	@echo "\n\033[1mAvailable commands:\033[0m\n"
	@@awk 'BEGIN {FS = ":.*##";} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


create-migration: ## Create an empty migration
	@read -p "Enter the sequence name: " SEQ; \
    docker run --rm -v ${MIGRATIONS_DIR}:/migrations migrate/migrate \
        create -ext sql -dir /migrations -seq $${SEQ}

migrate-up: ## Migration up
	@docker run --rm -v ${MIGRATIONS_DIR}:/migrations --network gormpgcare_gormpg_db_net migrate/migrate \
        -path=/migrations -database ${DATABASE_DSN} up

migrate-down: ## Migration down
	@read -p "Number of migrations you want to rollback (default: 1): " NUM; NUM=$${NUM:-1}; \
	docker run --rm -it -v ${MIGRATIONS_DIR}:/migrations --network gormpgcare_gormpg_db_net migrate/migrate \
        -path=/migrations -database ${DATABASE_DSN} down $${NUM}

migrate-force: ## Migration force version
	@read -p "Enter the version to force: " VERSION; \
	docker run --rm -it -v ${MIGRATIONS_DIR}:/migrations --network gormpgcare_gormpg_db_net migrate/migrate \
        -path=/migrations -database ${DATABASE_DSN} force $${VERSION}