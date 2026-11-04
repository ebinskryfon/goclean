.PHONY: docker-build run dev build migrate migrate-down test tidy clean help

help: ## Show this help menu
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run API server directly
	go run ./cmd/api

dev: ## Run with hot-reloading (air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Install via: go install github.com/air-verse/air@latest"; \
		echo "Falling back to go run..."; \
		go run ./cmd/api; \
	fi

build: ## Build API binary
	@mkdir -p bin
	go build -o bin/api ./cmd/api

docker-build: ## Build the application Docker image
	docker build -t {{.ImageName}}:latest .

migrate: ## Run pending database migrations
	go run ./cmd/migrate up

migrate-down: ## Rollback last database migration
	go run ./cmd/migrate down

test: ## Run unit and integration tests
	go test -v -race ./...

tidy: ## Download and tidy dependencies
	go mod tidy

clean: ## Clean built binaries and temporary files
	rm -rf bin tmp
