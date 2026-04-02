COMPOSE = docker compose
DB_URL = postgres://app:app@localhost:5432/app?sslmode=disable
MIGRATE = docker run --rm \
	-v $(PWD)/migrations:/migrations \
	--network host \
	migrate/migrate:v4.17.0

.PHONY: up down logs migrate-up migrate-down migrate-force

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

migrate-up:
	$(MIGRATE) -path=/migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path=/migrations -database "$(DB_URL)" down 1

migrate-force:
	$(MIGRATE) -path=/migrations -database "$(DB_URL)" force $(VERSION)
