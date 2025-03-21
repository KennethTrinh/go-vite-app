.PHONY: docker-ps docker-services docker-dev-up docker-dev-down docker-prod-up docker-prod-down docker-rollout


docker-ps:
	@docker ps -a --format "table {{.Names}}\t{{.CreatedAt}}\t{{.Status}}\t{{.Image}}\t"

docker-services:
	@docker compose ps --services

docker-dev-up:
	@echo "Starting development environment..."
	docker compose -f ./docker/docker-compose.dev.yml up --build -d

docker-dev-down:
	@echo "Stopping development environment..."
	docker compose -f ./docker/docker-compose.dev.yml down

docker-prod-up:
	# @echo "Starting production environment..."
	# @echo "Recording deployment info..."
	# @echo "DEPLOY" > .lastdeploy
	# @git log -1 --format="%H%n%an%n%ad%n%s" >> .lastdeploy
	# @cat .lastdeploy
	docker compose --env-file .env -f ./docker/docker-compose.yml up --build -d

docker-prod-down:
	@echo "Stopping production environment..."
	docker compose --env-file .env -f ./docker/docker-compose.yml down

docker-rollout:
ifndef service
	$(error Please specify a service: make docker-rollout service=<service-name>)
endif
	@echo "Rolling out ${service}..."
	@echo "Recording deployment info..."
	@echo "ROLLOUT: ${service}" > .lastdeploy
	@git log -1 --format="%H%n%an%n%ad%n%s" >> .lastdeploy
	@cat .lastdeploy
	docker compose build ${service}
	docker rollout ${service}