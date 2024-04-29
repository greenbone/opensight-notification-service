all: api-docs build test

.PHONY: install-code-generation-tools api-docs generate-code build test start-services

SWAG = github.com/swaggo/swag/cmd/swag@v1.16.2

install-code-generation-tools:
	go install $(SWAG)

api-docs:
	go run $(SWAG) init -d pkg/web --exclude pkg/web/healthcontroller -o api/notificationservice --parseDependency --generalInfo api.go --instanceName notificationservice
	go run $(SWAG) init -d pkg/web/healthcontroller -o api/health --parseDependency --generalInfo api.go --instanceName health

build:
	go build -o ./bin/ ./cmd/notification-service/

test: # run unit tests
	go test ./... -cover

start-services: ## start service and dependencies with docker
	docker compose -f docker-compose.yml -f docker-compose.service.yml up --build --abort-on-container-exit

cleanup-services: # delete service, dependencies and all persistent data
	docker compose -f docker-compose.yml -f docker-compose.service.yml down -v
