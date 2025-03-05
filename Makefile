VERSION?=latest
APPLICATION_NAME = cool

# Self-documented Makefile https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run Go tests
	@go test ./... -v -coverprofile coverage.out -covermode count

integration-test: ## Run Go Integration tests
	@go test ./internal/integrationtest/... -v -coverprofile coverage.out -covermode count --tags=integration

lint: ## Run linter
	@golangci-lint run

docker: ## Build docker
	docker build --platform linux/$(shell uname -m) --ssh default -t $(APPLICATION_NAME):${VERSION} . -f Dockerfile

generate-sql: ## Generate SQL
	sqlc generate