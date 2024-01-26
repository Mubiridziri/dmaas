.PHONY: build
build: ## Build a version
	go build -v ./cmd/dmaas

.PHONY: clean
clean: ## Remove temporary files
	go clean

.PHONY: dev
dev: ## Go Run
	export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=app port=5432 sslmode=disable" && go run cmd/dmaas/main.go

.PHONY: swag
swag: ## Update swagger.json
	swag init -g ./cmd/dmaas/main.go

.PHONY: swag-fmt
swag-fmt: ## Formatter for GoDoc (Swagger)
	swag fmt -g ./cmd/dmaas/main.go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build