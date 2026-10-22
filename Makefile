COMPOSE := docker compose -f docker-compose.infra.yaml

.PHONY: run infra-up infra-down infra-reset

# Запустить приложение (сначала поднимите инфраструктуру: make infra-up)
run:
	go run ./cmd/api

# Поднять Postgres и Redis и дождаться, пока они будут готовы
infra-up:
	$(COMPOSE) up -d --wait

# Остановить. Данные базы сохраняются
infra-down:
	$(COMPOSE) down

# Удалить всё вместе с данными и поднять заново: миграции применятся с нуля
infra-reset:
	$(COMPOSE) down -v
	$(COMPOSE) up -d --wait
