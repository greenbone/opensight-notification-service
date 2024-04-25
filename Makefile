all: api-docs build test

SWAG = github.com/swaggo/swag/cmd/swag@v1.16.2
GOLANGCI-LINT = github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	go run $(GOLANGCI-LINT) run

.PHONY: install-code-generation-tools
install-code-generation-tools:
	go install $(SWAG)

.PHONY: api-docs
api-docs:
	go run $(SWAG) init -d pkg/web --exclude pkg/web/healthcontroller -o api/notificationservice --parseDependency --generalInfo api.go --instanceName notificationservice
	go run $(SWAG) init -d pkg/web/healthcontroller -o api/health --parseDependency --generalInfo api.go --instanceName health

.PHONY: build
build:
	go build -o ./bin/ ./cmd/notification-service/

.PHONY: test
test: # run unit tests
	go test ./... -cover

.PHONY: start-services
start-services: ## start service and dependencies with docker
	docker compose -f docker-compose.yml -f docker-compose.service.yml up --build --abort-on-container-exit

.PHONY: cleanup-services
cleanup-services: # delete service, dependencies and all persistent data
	docker compose -f docker-compose.yml -f docker-compose.service.yml down -v
