.PHONY: dev build test migrate seed docker-up docker-down lint fmt help

# ── Local dev ──────────────────────────────────────────────────────────────────

dev: ## Run server + worker in dev mode (requires local postgres + redis)
	@echo "Starting server..."
	@air -c .air.toml & \
	 go run ./cmd/worker & \
	 wait

server: ## Run server only
	go run ./cmd/server

worker: ## Run worker only
	go run ./cmd/worker

# ── Build ──────────────────────────────────────────────────────────────────────

build: ## Build Go binaries
	CGO_ENABLED=0 go build -o bin/server ./cmd/server
	CGO_ENABLED=0 go build -o bin/worker ./cmd/worker

build-tool: ## Build React tool (dev bundle)
	cd tool && npm install && npm run build

build-embed: ## Build embeddable IIFE bundle
	cd tool && npm install && npm run build:embed

build-admin: ## Build admin dashboard
	cd admin && npm install && npm run build

build-all: build build-tool build-embed build-admin ## Build everything

# ── Database ───────────────────────────────────────────────────────────────────

migrate-up: ## Run DB migrations up
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" up

migrate-down: ## Roll back last migration
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" down 1

seed: ## Seed database with sample data
	bash scripts/seed.sh

provision: ## Provision a new tenant (SLUG= NAME= EMAIL= required)
	bash scripts/provision-tenant.sh "$(SLUG)" "$(NAME)" "$(EMAIL)"

# ── Docker ─────────────────────────────────────────────────────────────────────

docker-up: ## Start all services with Docker Compose
	docker compose up --build -d

docker-down: ## Stop all Docker services
	docker compose down

docker-logs: ## Tail server logs
	docker compose logs -f server worker

# ── Quality ────────────────────────────────────────────────────────────────────

test: ## Run Go tests
	go test ./... -v -timeout 30s

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format Go code
	gofmt -w .
	goimports -w .

tidy: ## Tidy Go modules
	go mod tidy

# ── Help ───────────────────────────────────────────────────────────────────────

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
