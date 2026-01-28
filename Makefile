all: generate-code api-docs lint build test

SWAG = github.com/swaggo/swag/cmd/swag@v1.16.4
MOCKERY = github.com/vektra/mockery/v3@v3.5.1
GOLANGCI-LINT = github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: help
help: ## show this help.
	@grep --no-filename -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## lint code
	@echo "\033[36m  Format code  \033[0m"
	go run $(GOLANGCI-LINT) run

.PHONY: install-code-generation-tools
install-code-generation-tools: ## Install code generation tools
	@ echo "\033[36m  Install code generation tools  \033[0m"
	go install $(SWAG)
	go install $(MOCKERY)

.PHONY: generate-code
generate-code: ## generate mocks
	@ echo "\033[36m  Generate mocks  \033[0m"
	go run $(MOCKERY) --log-level warn

.PHONY: api-docs
api-docs: ## generate swagger
	@ echo "\033[36m  Generate swagger docs  \033[0m"
	go run $(SWAG) init -d pkg/web --exclude pkg/web/healthcontroller -o api/notificationservice --parseDependency --generalInfo api.go --outputTypes yaml,go --instanceName notificationservice
	go run $(SWAG) init -d pkg/web/healthcontroller -o api/health --parseDependency --generalInfo api.go --outputTypes yaml,go --instanceName health

.PHONY: build
build: ## build app
	@echo "\033[36m  Build app  \033[0m"
	go build -o ./bin/ ./cmd/notification-service/

.PHONY: test
test: start-postgres-test-service ## run all tests
	@echo "\033[36m  Run tests  \033[0m"
	go test ./... -coverprofile=cov-unit-tests.txt

.PHONY: start-services
start-services: ## start service and dependencies with docker
	docker compose -f docker-compose.yml -f docker-compose.service.yml up --build --abort-on-container-exit

.PHONY: cleanup-services
cleanup-services: ## delete service, dependencies and all persistent data
	docker compose -f docker-compose.yml -f docker-compose.service.yml down -v

.PHONY: start-postgres-test-service
start-postgres-test-service: ## start test postgresql
	docker compose -f ./pkg/pgtesting/compose.yml -p postgres-test up -d --wait

.PHONY: stop-postgres-test-service
stop-postgres-test-service: ## stop test postgresql
	docker compose -f ./pkg/pgtesting/compose.yml -p postgres-test down
