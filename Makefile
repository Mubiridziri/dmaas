.PHONY: build
build: ## Build a version
	go build -v ./cmd/dmaas

.PHONY: clean
clean: ## Remove temporary files
	go clean

.PHONY: dev-start
dev-start:
	export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=app port=5432 sslmode=disable" && go run cmd/dmaas/main.go

.PHONY: swag # Update swagger.json
swag:
	swag init -g ./cmd/dmaas/main.go

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build