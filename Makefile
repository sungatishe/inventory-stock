COMPOSE_PATH = app
COMPOSE_FILES = -f $(COMPOSE_PATH)/docker-compose.yaml -f $(COMPOSE_PATH)/docker-compose.db.yaml -f $(COMPOSE_PATH)/docker-compose.tools.yaml

# Цели
up:
	docker-compose $(COMPOSE_FILES) up
up-build:
	docker-compose $(COMPOSE_FILES) up --build
down:
	docker-compose $(COMPOSE_FILES) down

# Удаление неиспользуемых ресурсов
prune:
	docker system prune -f --volumes