COMPOSE=docker compose

.PHONY: up down logs ps api worker

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

api:
	cd backend && go run ./cmd/api

worker:
	cd backend && go run ./cmd/worker
