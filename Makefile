.DEFAULT_GOAL := help
SHELL := /bin/bash

COMPOSE := docker compose -f deploy/compose/docker-compose.yml
COMPOSE_DEV := $(COMPOSE) -f deploy/compose/docker-compose.override.yml

.PHONY: help
help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ---------------------------------------------------------------- dev loop

.PHONY: up
up: ## Start dev stack (core profile)
	$(COMPOSE) --profile core up -d

.PHONY: up-obs
up-obs: ## Start observability stack (RAM heavy - not alongside core on 8GB)
	$(COMPOSE) --profile obs up -d

.PHONY: down
down: ## Stop all services
	$(COMPOSE) --profile core --profile obs --profile storage down

.PHONY: logs
logs: ## Follow logs from all services
	$(COMPOSE) logs -f --tail=100

.PHONY: ps
ps: ## Container status
	$(COMPOSE) ps

.PHONY: restart
restart: down up ## Restart the stack

# ---------------------------------------------------------------- go

.PHONY: run
run: ## Run API directly on host (no Docker)
	cd apps/api && go run ./cmd/api

.PHONY: build
build: ## Build API binary
	cd apps/api && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api

.PHONY: test
test: ## Unit tests with race detector
	cd apps/api && go test -race -count=1 ./...

.PHONY: test-cover
test-cover: ## Tests with coverage report
	cd apps/api && go test -race -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1

.PHONY: lint
lint: ## Lint Go (requires golangci-lint)
	cd apps/api && golangci-lint run ./...

.PHONY: fmt
fmt: ## Format Go code
	cd apps/api && gofmt -w . && go mod tidy

# ---------------------------------------------------------------- database

.PHONY: db-shell
db-shell: ## Open psql in the container
	$(COMPOSE) exec postgres psql -U campus -d campus_lms

.PHONY: migrate-up
migrate-up: ## Run migrations (Week 3)
	@echo "TODO Week 3: install goose/golang-migrate, then fill this target"
	@exit 1

.PHONY: migrate-down
migrate-down: ## Roll back one migration (Week 3)
	@echo "TODO Week 3"
	@exit 1

.PHONY: seed
seed: ## Seed dummy data (Week 3)
	@echo "TODO Week 3: 3 tenants, 50 lecturers, 2000 students, 200 courses"
	@exit 1

.PHONY: db-backup
db-backup: ## Back up the database (Week 4)
	bash deploy/scripts/backup.sh

.PHONY: db-restore
db-restore: ## Restore from backup - MUST be tested, not just written (Week 4)
	bash deploy/scripts/restore.sh

# ---------------------------------------------------------------- docker

.PHONY: docker-build
docker-build: ## Build API image locally
	docker build -t campus-lms-api:dev -f apps/api/Dockerfile apps/api

.PHONY: docker-size
docker-size: ## Check image size (Week 2 target: < 25MB)
	@docker images campus-lms-api:dev --format "{{.Repository}}:{{.Tag}} = {{.Size}}"

.PHONY: buildx
buildx: ## Multi-arch build amd64+arm64 (Week 2)
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t ghcr.io/CHANGE_ME/campus-lms-api:dev -f apps/api/Dockerfile apps/api

.PHONY: prune
prune: ## Clean Docker junk - run every Friday (256GB disk!)
	docker system prune -af --volumes

# ---------------------------------------------------------------- utilities

.PHONY: todo
todo: ## List all outstanding TODOs
	@grep -rn "TODO" --include="*.go" --include="*.py" --include="*.yml" \
		--include="*.yaml" --include="*.md" --include="*.sh" --include="go.mod" \
		--include="Dockerfile" --include="Makefile" --include="Caddyfile" \
		--exclude-dir=.git . 2>/dev/null \
		| grep -v "^./README.md" | grep -v Binary || echo "  No TODOs left"

.PHONY: journal
journal: ## Create today's journal entry
	@f=docs/journal/$$(date +%Y-%m-%d).md; \
	if [ -f $$f ]; then echo "already exists: $$f"; else \
	printf '# %s\n\n## 3 target hari ini\n1. \n2. \n3. \n\n## Yang berhasil\n\n## Yang macet\n\n## Satu hal yang saya pelajari\n\n## Angka hari ini (kalau ada)\n\n' "$$(date +%Y-%m-%d)" > $$f; \
	echo "created: $$f"; fi

.PHONY: health
health: ## Check API health endpoints
	@curl -sf http://localhost:8080/healthz && echo " ✓ healthz OK" || echo " ✗ healthz GAGAL"
	@curl -sf http://localhost:8080/readyz  && echo " ✓ readyz OK"  || echo " ✗ readyz GAGAL"

# ---------------------------------------------------------------- learning system

.PHONY: week-init
week-init: ## Scaffold a new week: report, quiz, evidence dir (make week-init W=01)
	@test -n "$(W)" || (echo "usage: make week-init W=01" && exit 1)
	@mkdir -p docs/progress/evidence/week-$(W)
	@test -f docs/progress/week-$(W).md || cp agent/templates/weekly-report.md docs/progress/week-$(W).md
	@test -f docs/progress/quiz/week-$(W).md || cp agent/templates/quiz.md docs/progress/quiz/week-$(W).md
	@echo "ready: docs/progress/week-$(W).md, quiz/week-$(W).md, evidence/week-$(W)/"

.PHONY: evidence
evidence: ## Capture a command as evidence (make evidence W=01 SLUG=image-size CMD="docker images")
	@test -n "$(W)" -a -n "$(SLUG)" -a -n "$(CMD)" || 		(echo 'usage: make evidence W=01 SLUG=image-size CMD="docker images"' && exit 1)
	@mkdir -p docs/progress/evidence/week-$(W)
	@f=docs/progress/evidence/week-$(W)/$(SLUG).txt; 	{ 	  echo "=== EVIDENCE ==="; 	  echo "CLAIM:    $${CLAIM:-<fill in>}"; 	  echo "COMMAND:  $(CMD)"; 	  echo "CWD:      $$(pwd)"; 	  echo "RUN AT:   $$(date -Iseconds)"; 	  echo "COMMIT:   $$(git rev-parse --short HEAD 2>/dev/null || echo 'no-git')"; 	} > $$f; 	set +e; out=$$(eval "$(CMD)" 2>&1); code=$$?; set -e; 	echo "EXIT:     $$code" >> $$f; 	echo "=== RAW OUTPUT ===" >> $$f; 	echo "$$out" >> $$f; 	echo "=== END ===" >> $$f; 	echo "saved: $$f (exit $$code)"

.PHONY: gate
gate: ## Show the week-completion gate checklist
	@echo ""
	@echo "  GATE — minggu berikutnya tidak dimulai sebelum semua ini beres:"
	@echo "    1. Laporan mingguan ditandatangani"
	@echo "    2. Skor quiz >= 70%"
	@echo "    3. 3 file bukti sudah di-spot-check acak"
	@echo "    4. Explain-back 3 menit terekam"
	@echo ""
	@echo "  Checklist lengkap: agent/checklists/human-verification.md"
	@echo ""

.PHONY: agent-context
agent-context: ## Print the files an agent must read before starting work
	@echo "Required reading before any task:"
	@echo "  AGENTS.md"
	@echo "  agent/policy.md"
	@echo "  agent/evidence-protocol.md"
	@echo "  agent/rules/00-global.md + the rule file for your area"
	@echo "  the current week section of the roadmap"
