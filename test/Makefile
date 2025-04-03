# ! PUBLIC SCHEMA WILL BE DROPPED DURING INTEGRATION TESTS ! DO NOT CHANGE DATA SOURCE OR USE ONLY TEST DATABASE
DB_NAME = test
DB_USER = postgres
DB_PASSWORD = password
POSTGRES_PORT = 5432
POSTGRES_DSN = postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(POSTGRES_PORT)/$(DB_NAME)?sslmode=disable

COMPOSE_PATH = ./docker-compose.yml
COMPOSE_ENV = DB_NAME=$(DB_NAME) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) POSTGRES_PORT=$(POSTGRES_PORT)

integration-test: postgres-up
	go test . -dsn $(POSTGRES_DSN) -tags integration

postgres-up:
	$(COMPOSE_ENV) docker compose -f $(COMPOSE_PATH) up -d

postgres-down:
	$(COMPOSE_ENV) docker compose -f $(COMPOSE_PATH) down

.PHONY: integration-test postgres-up postgres-down