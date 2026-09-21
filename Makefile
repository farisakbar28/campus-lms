.DEFAULT_GOAL := help
SHELL := /bin/bash

COMPOSE := docker compose -f deploy/compose/docker-compose.yml
COMPOSE_DEV := $(COMPOSE) -f deploy/compose/docker-compose.override.yml

.PHONY: help
help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ---------------------------------------------------------------- development

.PHONY: up
up: ## Start the development core stack
	$(COMPOSE_DEV) --profile core up -d

.PHONY: down
down: ## Stop all development services
	$(COMPOSE) --profile core --profile obs --profile storage down

.PHONY: logs
logs: ## Follow service logs
	$(COMPOSE) logs -f --tail=100

.PHONY: ps
ps: ## Show container status
	$(COMPOSE) ps

.PHONY: restart
restart: down up ## Restart the development stack

# ---------------------------------------------------------------- go

.PHONY: run
run: ## Run the API directly on the host
	cd apps/api && go run ./cmd/api

.PHONY: build
build: ## Build the API binary
	cd apps/api && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api

.PHONY: test
test: ## Run Go tests with the race detector
	cd apps/api && go test -race -count=1 -p 1 ./...

.PHONY: test-cover
test-cover: ## Run tests and print coverage summary
	cd apps/api && go test -race -p 1 -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1

.PHONY: lint
lint: ## Lint Go (requires golangci-lint)
	cd apps/api && golangci-lint run ./...

.PHONY: fmt
fmt: ## Format Go code and tidy module metadata
	cd apps/api && gofmt -w . && go mod tidy

# ---------------------------------------------------------------- database

.PHONY: db-shell
db-shell: ## Open psql in the database container
	$(COMPOSE) exec postgres psql -U campus -d campus_lms

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@$(COMPOSE_DEV) exec -T api sh -ec '\
		MIGRATE_PATH="/go/bin/migrate"; \
		WANT_VER="github.com/golang-migrate/migrate/v4 v4.18.3"; \
		if [ -x "$$MIGRATE_PATH" ]; then \
			CUR_VER=$$(go version -m "$$MIGRATE_PATH" 2>/dev/null | grep "$$WANT_VER" || true); \
			if [ -z "$$CUR_VER" ]; then \
				echo "migrate binary found but version mismatch or error. Reinstalling..."; \
				rm -f "$$MIGRATE_PATH"; \
			fi; \
		fi; \
		if [ ! -x "$$MIGRATE_PATH" ]; then \
			echo "Installing golang-migrate..."; \
			go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3; \
		fi; \
		test -n "$$MIGRATE_DATABASE_URL" || (echo "MIGRATE_DATABASE_URL is not set" && exit 1); \
		exec /go/bin/migrate -path /src/migrations -database "$$MIGRATE_DATABASE_URL" up \
	'

.PHONY: migrate-down
migrate-down: ## Roll back exactly one migration
	@$(COMPOSE_DEV) exec -T api sh -ec '\
		MIGRATE_PATH="/go/bin/migrate"; \
		WANT_VER="github.com/golang-migrate/migrate/v4 v4.18.3"; \
		if [ -x "$$MIGRATE_PATH" ]; then \
			CUR_VER=$$(go version -m "$$MIGRATE_PATH" 2>/dev/null | grep "$$WANT_VER" || true); \
			if [ -z "$$CUR_VER" ]; then \
				echo "migrate binary found but version mismatch or error. Reinstalling..."; \
				rm -f "$$MIGRATE_PATH"; \
			fi; \
		fi; \
		if [ ! -x "$$MIGRATE_PATH" ]; then \
			echo "Installing golang-migrate..."; \
			go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3; \
		fi; \
		test -n "$$MIGRATE_DATABASE_URL" || (echo "MIGRATE_DATABASE_URL is not set" && exit 1); \
		exec /go/bin/migrate -path /src/migrations -database "$$MIGRATE_DATABASE_URL" down 1 \
	'

.PHONY: migrate-version
migrate-version: ## Show the migration version and dirty state
	@$(COMPOSE_DEV) exec -T api sh -ec '\
		MIGRATE_PATH="/go/bin/migrate"; \
		WANT_VER="github.com/golang-migrate/migrate/v4 v4.18.3"; \
		if [ -x "$$MIGRATE_PATH" ]; then \
			CUR_VER=$$(go version -m "$$MIGRATE_PATH" 2>/dev/null | grep "$$WANT_VER" || true); \
			if [ -z "$$CUR_VER" ]; then \
				echo "migrate binary found but version mismatch or error. Reinstalling..."; \
				rm -f "$$MIGRATE_PATH"; \
			fi; \
		fi; \
		if [ ! -x "$$MIGRATE_PATH" ]; then \
			echo "Installing golang-migrate..."; \
			go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3; \
		fi; \
		test -n "$$MIGRATE_DATABASE_URL" || (echo "MIGRATE_DATABASE_URL is not set" && exit 1); \
		exec /go/bin/migrate -path /src/migrations -database "$$MIGRATE_DATABASE_URL" version \
	'

.PHONY: seed
seed: ## Seed the local database with test data
	@echo "Running seed script..."
	@cat apps/api/testdata/seed.sql | $(COMPOSE_DEV) exec -T postgres sh -ec '\
		psql -v ON_ERROR_STOP=1 -U "$$POSTGRES_USER" -d "$$POSTGRES_DB" \
	'
	@echo "Seeding completed."

.PHONY: db-backup
db-backup: ## Create a local custom-format database backup
	bash deploy/scripts/backup.sh

.PHONY: db-restore
db-restore: ## Restore a local backup into a disposable database
	bash deploy/scripts/restore.sh

# ---------------------------------------------------------------- docker

.PHONY: docker-build
docker-build: ## Build the API image locally
	docker build -t campus-lms-api:dev -f apps/api/Dockerfile apps/api

.PHONY: docker-size
docker-size: ## Show the local API image size
	@docker images campus-lms-api:dev --format "{{.Repository}}:{{.Tag}} = {{.Size}}"

.PHONY: buildx
buildx: ## Build the API image for amd64 and arm64
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t ghcr.io/farisakbar28/campus-lms-api:dev -f apps/api/Dockerfile apps/api

.PHONY: prune
prune: ## Remove unused Docker artifacts while keeping volumes
	docker system prune -af
	@echo ""
	@echo "Volumes were not touched. Database data is safe."
	@echo "Use 'make prune-hard' only when you intend to destroy all data."

.PHONY: prune-hard
prune-hard: ## DANGER: remove unused Docker artifacts and volumes
	@echo "This deletes all Docker volumes, including campus-lms_postgres-data."
	@read -p "Type DESTROY to confirm: " c && [ "$$c" = "DESTROY" ] || (echo "Aborted."; exit 1)
	docker system prune -af --volumes

# ---------------------------------------------------------------- utilities

.PHONY: todo
todo: ## List outstanding TODO markers
	@grep -rn "TODO" --include="*.go" --include="*.py" --include="*.yml" \
		--include="*.yaml" --include="*.md" --include="*.sh" --include="go.mod" \
		--include="Dockerfile" --include="Makefile" --include="Caddyfile" \
		--exclude-dir=.git . 2>/dev/null \
		| grep -v "^./README.md" | grep -v Binary || echo "  No TODOs left"

.PHONY: health
health: ## Check API health endpoints
	@curl -sf http://localhost:8080/healthz && echo " healthz OK" || echo " healthz FAILED"
	@curl -sf http://localhost:8080/readyz  && echo " readyz OK"  || echo " readyz FAILED"
