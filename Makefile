DOCKER_COMPOSE = /usr/libexec/docker/cli-plugins/docker-compose

.PHONY: help build up down restart logs test proto seed

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0ms %s\n", $$1, $$2}'

build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

run:
	$(DOCKER_COMPOSE) up --build

down:
	$(DOCKER_COMPOSE) down -v

restart: down run

logs:
	$(DOCKER_COMPOSE) logs -f

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/auth.proto

seed:
	docker exec -it $$(docker ps -qf "name=postgres") psql -U postgres -d auth_db -c "INSERT INTO users (token, user_id, rate_limit, tier) VALUES ('gold-token', 'user_pro_1', 100, 'gold') ON CONFLICT DO NOTHING;"

status:
	$(DOCKER_COMPOSE) ps
