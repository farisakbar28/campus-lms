# Audit Repo — campus-lms

**Waktu audit:** 2026-08-18 (Asia/Makassar)  
**Batas audit:** baca-saja. Tidak ada build, test, Docker build/Compose, migrasi, atau aksi Git tulis yang dijalankan. Satu-satunya berkas yang dibuat adalah laporan ini.

**Konvensi:** setiap fakta di bawah merujuk perintah audit yang benar-benar dijalankan. `TIDAK TERVERIFIKASI` berarti audit ini tidak memperoleh bukti yang diperlukan; itu bukan dugaan gagal/lulus.

## A. Lingkungan

Perintah: `go version; docker version; docker compose version; docker buildx version; make --version; git --version; uname -a; cat /etc/os-release; free -h; systemctl is-system-running; systemctl --version; du -sh --exclude=.git .`

| Item | Hasil |
|---|---|
| Path repo | `/home/hype/projects/campus-lms` |
| Go | `go1.23.6 linux/amd64` |
| Docker client | Docker Engine Community `29.7.2`, API `1.55`, `linux/amd64` |
| Docker server | **TIDAK TERVERIFIKASI** — `docker version` tidak dapat menyambung ke `/var/run/docker.sock`: `permission denied`. |
| Docker Compose | `v5.5.0` |
| Buildx | `v0.36.1` |
| Make | GNU Make `4.4.1` |
| Git | `2.53.0` |
| OS/kernel | Ubuntu `26.04 LTS` (Resolute Raccoon); `Linux HYPE 6.6.87.2-microsoft-standard-WSL2 ... x86_64` |
| Memori (`free -h`) | RAM 4.8 GiB total, 1.8 GiB terpakai, 2.5 GiB bebas, 3.0 GiB tersedia; swap 8.0 GiB, 1.8 MiB terpakai. |
| systemd aktif | **TIDAK TERVERIFIKASI** — binary `systemd 259` ada, tetapi `systemctl is-system-running` gagal mengakses system bus: `Operation not permitted`. |
| Ukuran repo tanpa `.git` | `920K` pada saat diukur (sebelum laporan ini dibuat). |

## B. Inventaris file

Perintah: `find . -path './.git' -prune -o -print | sort` dan `rg -l -i --hidden --glob '!.git/**' 'TODO|FIXME|PLACEHOLDER|CHANGE_ME' | sort`.

Pohon berikut adalah snapshot sebelum laporan ini ditambahkan; isi `.env` **tidak dibaca**.

```text
.
./.agents
./.codex
./.editorconfig
./.env
./.env.example
./.github/workflows/ai-eval.yml
./.github/workflows/cd.yml
./.github/workflows/ci.yml
./.gitignore
./AGENTS.md
./Makefile
./README.md
./WORKFLOW.md
./agent/{README.md,checklists,evidence-protocol.md,policy.md,prompts,rules,templates}
./agent/checklists/{human-verification.md,pre-commit.md}
./agent/prompts/{README.md,debug.md,dev.md,heavy.md,plan.md,quick.md,review.md,spark.md,teach.md}
./agent/rules/{00-global.md,10-go-api.md,20-database.md,30-docker-deploy.md,40-ci-cd.md,50-python-ai.md,60-security.md,70-docs.md}
./agent/templates/{quiz.md,session-log.md,weekly-report.md}
./apps/{ai,api,web}
./apps/ai/README.md
./apps/api/{.dockerignore,Dockerfile,go.mod,cmd,internal,migrations,scripts}
./apps/api/cmd/api/main.go
./apps/api/internal/{config,domain,healthcheck,http,middleware,repository}
./apps/api/internal/config/{.gitkeep,config.go,config_test.go}
./apps/api/internal/domain/.gitkeep
./apps/api/internal/healthcheck/{probe.go,probe_test.go}
./apps/api/internal/http/{.gitkeep,server.go,server_test.go}
./apps/api/internal/{middleware,repository}/.gitkeep
./apps/api/migrations/.gitkeep
./apps/api/scripts/{verify-graceful-shutdown.sh,verify-healthz.sh}
./apps/web/README.md
./deploy/{caddy,compose,scripts}
./deploy/caddy/Caddyfile
./deploy/compose/{docker-compose.override.yml,docker-compose.yml}
./deploy/scripts/{backup.sh,deploy.sh,restore.sh}
./docs/{adr,domain-ai.md,domain.md,notes,progress,roadmap.md,runbook,setup}
./docs/adr/{0001-pilihan-stack.md,0002c-azure-conventions.md,0005-api-healthcheck-probe.md,README.md,template.md}
./docs/notes/{README.md,docker-internals.md}
./docs/progress/{README.md,evidence,explain,quiz,sessions,week-00.md,week-01.md,week-02.md}
./docs/progress/evidence/{EXAMPLE-format.txt,README.md,week-00,week-01,week-02}
./docs/progress/explain/{.gitkeep,week-00.txt,week-01.txt}
./docs/progress/quiz/{week-00.md,week-01.md,week-02.md}
./docs/progress/sessions/.gitkeep
./docs/runbook/incident.md
./docs/setup/azure-day-0.md
```

Ekspansi path file yang berada di direktori yang diringkas oleh notasi `{...}` di atas (sehingga inventaris tetap lengkap):

```text
./.github
./.github/workflows
./agent/checklists
./agent/prompts
./agent/rules
./agent/templates
./apps
./apps/ai
./apps/api/cmd
./apps/api/cmd/api
./apps/api/internal
./apps/api/internal/config
./apps/api/internal/domain
./apps/api/internal/healthcheck
./apps/api/internal/http
./apps/api/internal/middleware
./apps/api/internal/repository
./apps/api/migrations
./apps/api/scripts
./apps/web
./deploy
./deploy/caddy
./deploy/compose
./deploy/scripts
./docs
./docs/adr
./docs/notes
./docs/progress
./docs/progress/evidence
./docs/progress/evidence/week-00
./docs/progress/evidence/week-00/.gitkeep
./docs/progress/evidence/week-00/azure-account-show.txt
./docs/progress/evidence/week-00/env-gitignore.txt
./docs/progress/evidence/week-00/make-todo.txt
./docs/progress/evidence/week-00/ssh-key-permissions.txt
./docs/progress/evidence/week-00/wsl-resources.txt
./docs/progress/evidence/week-00/wslconfig.txt
./docs/progress/evidence/week-01
./docs/progress/evidence/week-01/go-test-report-rerun.txt
./docs/progress/evidence/week-01/go-test.txt
./docs/progress/evidence/week-01/graceful-shutdown.txt
./docs/progress/evidence/week-01/healthz-curl.txt
./docs/progress/evidence/week-01/no-println.txt
./docs/progress/evidence/week-01/report-source-git-log.txt
./docs/progress/evidence/week-01/ssh-keyonly.txt
./docs/progress/evidence/week-02
./docs/progress/evidence/week-02/api-probe-healthcheck.txt
./docs/progress/evidence/week-02/api-probe-runtime-build.txt
./docs/progress/evidence/week-02/build-context-after.txt
./docs/progress/evidence/week-02/build-context-before.txt
./docs/progress/evidence/week-02/buildx-multiarch.txt
./docs/progress/evidence/week-02/busybox-layer-history.txt
./docs/progress/evidence/week-02/busybox-runtime-build.txt
./docs/progress/evidence/week-02/busybox-vs-api-probe-size.txt
./docs/progress/evidence/week-02/docker-build-current.txt
./docs/progress/evidence/week-02/docker-inspect-user.txt
./docs/progress/evidence/week-02/docker-size-api-probe.txt
./docs/progress/evidence/week-02/go-test-healthcheck.txt
./docs/progress/evidence/week-02/healthcheck-distroless.txt
./docs/progress/evidence/week-02/image-size.txt
./docs/progress/explain
./docs/progress/quiz
./docs/progress/sessions
./docs/runbook
./docs/setup
```

| Kategori hasil pencarian placeholder/TODO | File yang cocok |
|---|---|
| Implementasi yang masih placeholder/TODO | `.github/workflows/{ai-eval.yml,cd.yml,ci.yml}`; `Makefile`; `apps/api/go.mod`; `apps/web/README.md`; `deploy/caddy/Caddyfile`; `deploy/compose/{docker-compose.yml,docker-compose.override.yml}`; `deploy/scripts/{backup.sh,deploy.sh,restore.sh}`; `docs/adr/README.md`; `docs/progress/week-02.md` |
| Referensi/rencana/arsip yang mengandung kata itu (bukan otomatis pekerjaan aktif) | `AGENTS.md`, `README.md`, `WORKFLOW.md`, `agent/policy.md`, `agent/prompts/README.md`, `docs/progress/evidence/week-00/make-todo.txt`, `docs/progress/explain/week-00.txt`, `docs/progress/quiz/week-00.md`, `docs/progress/week-00.md`, `docs/roadmap.md`, `docs/setup/azure-day-0.md` |

## C. Status Git

Perintah: `git remote -v; git branch --show-current; git rev-list --count HEAD; git status --short; git log --format='%H - %s'; git for-each-ref ...; git rev-list --left-right --count master...origin/master; git reflog show --all ...`.

| Item | Hasil |
|---|---|
| Remote | `origin https://github.com/farisakbar28/campus-lms.git` (fetch dan push) |
| Branch aktif | `master` |
| Jumlah commit `HEAD` | 19 |
| Tracking ref lokal | `origin/master` menunjuk `a8c99bf`; `master` menunjuk `ae32c04` |
| Selisih terhadap `origin/master` yang tersimpan lokal | `5 0` (master 5 commit di depan, 0 di belakang). Remote **tidak di-fetch** karena itu mengubah ref lokal. |
| Perubahan belum commit | `M Makefile`; untracked: `docs/notes/docker-internals.md`, `docs/progress/evidence/week-02/`, `docs/progress/quiz/week-02.md`, `docs/progress/week-02.md` |
| Pernah push | Ya, **secara historis**: reflog memuat beberapa entri `update by push`. Status remote saat ini **TIDAK TERVERIFIKASI** tanpa fetch/network. |

Semua commit:

```text
ae32c04627e18a63164f712725d193e975e9804c - docs(compose): defer web from core profile to Week 4
db71c5c96ed92f2b4b68a85dcb28d403890ca645 - docs(adr): draft healthcheck probe options
3bb7565beacd38dffca7e94666e7a1d73ed38040 - fix(api): replace busybox health probe
8ac8e6bf4a5c8073048c11b8dd886da8c595b70b - build(api): restrict Docker build context
ad1bc11aed2885bc106fb8a29943121af6753d57 - feat(api): add production Dockerfile
a8c99bf83bce677dcff71eb1a7bcd95c55d0cd92 - docs: add go.mod for dependency guide.
1ff2747cc92785473e39c9ff860e07a6f415f615 - docs: mark Week 1 complete
f37226985ecc6ee774985dd3103406d0127134bd - docs: verification and selected stack
d90795a2d80d89a60584f7faec7074da33c85e98 - docs(progress): complete Week 1 report, quiz, and sign-off
73f24f7731afdb43785b558536c8abe76b973d96 - docs(evidence): capture SSH key-only auth verification
e7db51a76bee9cee0d94f90584c6faafa9564d16 - docs(domain): add implementation tiers, amendments, and AI layer
e4556de6b7e9717aeeb57e7dd9b6a7797d547e3a - docs(agent): forbid editing task briefs and metric-gaming
3195c7b9ca12f6b916c152733c6afe07459bfa06 - docs(evidence): refine Week 1 stdout check
eda3ba3c572bc12004b4bf6056356fe06cd24c8d - fix(api): restore task brief text and use precise stdout check
fe7133125d7c31f918dde9b8dacd6d0f2f5718b0 - docs(evidence): capture Week 1 API verification
46939229967c5569369e3d4de764dc27981e4ef7 - feat(api): add Week 1 HTTP server
197c42d9b694542082dd9e5ea4943765a2762f46 - docs: bring roadmap into repo, clarify agent context
50c84b6d19f8e5f1e87c4b238f8e27624f636466 - only just setup & week 0 complete
a0d2890b4b4c5a9b453a01216fbdac35ff34f247 - initial commit
```

Dengan pola yang disebut AGENTS (`feat|fix|docs|chore|refactor|test`, scope opsional): **16 patuh, 3 tidak patuh**. Yang tidak patuh: `8ac8e6b build(api): ...`, `50c84b6 only just setup & week 0 complete`, `a0d2890 initial commit`.

## D. Kode aplikasi

Perintah: `find apps/api -type f -name '*.go' ... wc -l`; `find ... awk '/^package/'`; `sed apps/api/go.mod`; `test -f apps/api/go.sum`; `rg` untuk route/fungsi/env; `sed apps/api/internal/config/config.go`; `sed .env.example`.

| File Go | Baris | Package |
|---|---:|---|
| `cmd/api/main.go` | 150 | `main` |
| `internal/config/config.go` | 93 | `config` |
| `internal/config/config_test.go` | 48 | `config` |
| `internal/healthcheck/probe.go` | 32 | `healthcheck` |
| `internal/healthcheck/probe_test.go` | 43 | `healthcheck` |
| `internal/http/server.go` | 59 | `http` |
| `internal/http/server_test.go` | 45 | `http` |
| **Total** | **470** | `main`, `config`, `healthcheck`, `http` |

| Item | Hasil |
|---|---|
| Isi `go.mod` efektif | Module `github.com/farisakbar28/campus-lms/apps/api`; Go `1.23`; tidak ada directive `require`. |
| Dependency eksternal | 0 (komentar roadmap dependency tidak dihitung sebagai dependency Go). |
| `go.sum` | Tidak ada. |
| Test files | `internal/config/config_test.go`; `internal/healthcheck/probe_test.go`; `internal/http/server_test.go` |
| Endpoint HTTP terdaftar | `GET /healthz` → `healthz(logger)`; `GET /readyz` → `readyz(logger)`, keduanya dalam `internal/http/server.go:24-25`. |
| Fungsi produksi | `main`, `newLogger`, `bootstrapLogger`, `buildVersion`, `Load`, `Address`, `required`, `parsePort`, `Probe`, `NewServer`, `healthz`, `readyz`, `writeStatus`. |

Variabel yang **benar-benar dibaca** oleh `internal/config` (`os.LookupEnv` melalui `required`): `APP_ENV`, `APP_PORT`, `APP_LOG_LEVEL`, `APP_SHUTDOWN_TIMEOUT`.

| Perbandingan `internal/config` vs `.env.example` | Nama |
|---|---|
| Ada di kedua sisi | `APP_ENV`, `APP_PORT`, `APP_LOG_LEVEL`, `APP_SHUTDOWN_TIMEOUT` |
| Dibaca kode tetapi tidak ada di `.env.example` | Tidak ada |
| Ada di `.env.example` tetapi belum dibaca `internal/config` | `DATABASE_URL`, `DB_MAX_CONNS`, `DB_MIN_CONNS`, `REDIS_URL`, `JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_BUCKET`, `AI_SERVICE_URL`, `GEMINI_API_KEY`, `GROQ_API_KEY`, `CEREBRAS_API_KEY`, `OPENROUTER_API_KEY`, `LLM_MONTHLY_TOKEN_BUDGET`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME`, `LANGFUSE_HOST`, `LANGFUSE_PUBLIC_KEY`, `LANGFUSE_SECRET_KEY` |

## E. Docker

Perintah: `sed apps/api/Dockerfile`; `sed apps/api/.dockerignore`; `sed deploy/compose/docker-compose.yml`; `docker images`; `docker ps -a`.

`apps/api/Dockerfile` (kutipan utuh):

```dockerfile
# -----------------------------------------------------------------------------
# TASK BRIEF — Week 2
# Agent: read agent/rules/30-docker-deploy.md before implementing.
# -----------------------------------------------------------------------------
#
# TARGETS
#   - multi-stage build (builder + minimal runtime)
#   - final image < 25 MB   (verify: make docker-size)
#   - runs as a NON-ROOT user (verify with docker inspect)
#   - HEALTHCHECK present
#   - builds for both linux/amd64 and linux/arm64
#   - layer caching: `go mod download` in its own layer before copying source
#
# HARD RULES
#   - never place a secret in ARG or ENV (it persists in image layers)
#   - CGO_ENABLED=0 for a static binary
#
# QUESTIONS THE HUMAN MUST BE ABLE TO ANSWER AFTERWARDS
#   - Why does CGO_ENABLED=0 matter for scratch/distroless images?
#   - distroless vs alpine: what is the debugging trade-off?
#   - Why does COPY order dominate rebuild speed?
#   - Why is a secret in a build ARG trivially extractable?
#
# EVIDENCE REQUIRED
#   docs/progress/evidence/week-02/image-size.txt
#   docs/progress/evidence/week-02/docker-inspect-user.txt
#   docs/progress/evidence/week-02/buildx-multiarch.txt

FROM --platform=$BUILDPLATFORM golang:1.23.6-alpine3.21 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Copy the module definition before application source so dependency downloads
# remain cached when only Go source files change.
COPY go.mod go.sum* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# distroless intentionally has neither a shell nor an HTTP client. BusyBox is
# copied solely to run the Docker health probe without adding a full OS layer.
FROM busybox:1.37.0-musl AS healthcheck

# This target is retained only for the Week 2 size comparison. It deliberately
# keeps the Go toolchain in the runtime image; production uses the final stage.
FROM golang:1.23.6-alpine3.21 AS alpine-runtime

COPY --from=builder /out/api /api

USER 65532:65532

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["wget", "--spider", "--quiet", "http://127.0.0.1:8080/healthz"]

ENTRYPOINT ["/api"]

FROM gcr.io/distroless/static-debian12:nonroot AS busybox-runtime

COPY --from=builder /out/api /api
COPY --from=healthcheck /bin/busybox /busybox

USER nonroot:nonroot

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/busybox", "wget", "--spider", "--quiet", "http://127.0.0.1:8080/healthz"]

ENTRYPOINT ["/api"]

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/api /api

USER nonroot:nonroot

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/api", "-healthcheck"]

ENTRYPOINT ["/api"]
```

`apps/api/.dockerignore`:

```text
# Send only module metadata and production Go source to the build context.
*
!go.mod
!go.sum
!cmd/
!cmd/**
!internal/
!internal/**

# Test inputs and local build artefacts never participate in `go build`.
**/*_test.go
**/testdata/
bin/
coverage.out
```

| Item | Hasil |
|---|---|
| Compose utama | Masih TODO: hanya TASK BRIEF dan `# TODO Week 2`, **tidak ada** key/service Compose. Tidak ada isi service untuk dikutip. |
| Override Compose | Masih TODO: hanya TASK BRIEF dan `# TODO Week 2`; tidak ada service/hot-reload mount. |
| Image lokal relevan | **TIDAK TERVERIFIKASI** — `docker images` gagal: `permission denied while trying to connect to the docker API at unix:///var/run/docker.sock`. |
| Semua container | **TIDAK TERVERIFIKASI** — `docker ps -a` gagal dengan galat identik. |

## F. Makefile

Perintah: `sed -n '1,360p' Makefile`; `rg -n ... TODO/CHANGE_ME`; pemeriksaan keberadaan `find deploy/compose`, `find apps/api/{migrations,internal/repository,internal/middleware}`.

| Target | Perintah/resep |
|---|---|
| `help` | `grep ... $(MAKEFILE_LIST) | awk ...` |
| `up` | `$(COMPOSE_DEV) --profile core up -d` |
| `up-obs` | `$(COMPOSE) --profile obs up -d` |
| `down` | `$(COMPOSE) --profile core --profile obs --profile storage down` |
| `logs` | `$(COMPOSE) logs -f --tail=100` |
| `ps` | `$(COMPOSE) ps` |
| `restart` | prerequisite `down up` |
| `run` | `cd apps/api && go run ./cmd/api` |
| `build` | `cd apps/api && CGO_ENABLED=0 go build ... -o bin/api ./cmd/api` |
| `test` | `cd apps/api && go test -race -count=1 ./...` |
| `test-cover` | `cd apps/api && go test -race -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1` |
| `lint` | `cd apps/api && golangci-lint run ./...` |
| `fmt` | `cd apps/api && gofmt -w . && go mod tidy` |
| `db-shell` | `$(COMPOSE) exec postgres psql -U campus -d campus_lms` |
| `migrate-up` | mencetak `TODO Week 3...`; `exit 1` |
| `migrate-down` | mencetak `TODO Week 3`; `exit 1` |
| `seed` | mencetak `TODO Week 3...`; `exit 1` |
| `db-backup` | `bash deploy/scripts/backup.sh` |
| `db-restore` | `bash deploy/scripts/restore.sh` |
| `docker-build` | `docker build -t campus-lms-api:dev -f apps/api/Dockerfile apps/api` |
| `docker-size` | `docker images campus-lms-api:dev --format ...` |
| `buildx` | `docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/<username-github-mu>/campus-lms-api:dev ...` |
| `prune` | `docker system prune -af --volumes` |
| `todo` | `grep -rn "TODO" ... | grep ...` |
| `journal` | membuat `docs/journal/YYYY-MM-DD.md` jika belum ada |
| `health` | dua `curl -sf` ke `/healthz` dan `/readyz` |
| `week-init` | membuat direktori evidence dan menyalin template report/quiz |
| `evidence` | membuat file evidence dan menjalankan `$(CMD)` yang diberikan |
| `gate` | mencetak checklist gate |
| `agent-context` | mencetak daftar konteks wajib agent |

| Pemeriksaan khusus | Hasil |
|---|---|
| `up` memakai apa? | **`COMPOSE_DEV`**, bukan `COMPOSE`; perubahan ini belum commit (`git diff -- Makefile`). |
| Target yang merujuk service belum ada | `db-shell` menyebut service `postgres`, tetapi kedua file Compose tidak mempunyai service. `up`, `up-obs`, `down`, `logs`, dan `ps` juga tidak memiliki definisi service yang dapat dijalankan. |
| Target yang belum diimplementasi | `migrate-up`, `migrate-down`, `seed`; skrip yang dipanggil `db-backup`/`db-restore` ada tetapi masing-masing masih TODO Week 4. |
| Referensi file/direktori yang belum ada | `journal` akan membuat `docs/journal/...`; pada audit, direktori `docs/journal` belum ada. Ini perilaku penciptaan target, bukan file yang wajib sudah ada. |
| `CHANGE_ME` | Tidak ditemukan pada source konfigurasi aktif. Pencarian menemukannya hanya di dokumen/arsip: `docs/setup/azure-day-0.md`, `docs/roadmap.md`, `docs/progress/week-00.md`, dan evidence historis. |

## G. Dokumentasi

Perintah: `find docs -type f ... wc -l`; `rg '^# ' docs/domain.md`; `rg '^# (37|38)' docs/domain.md`; pembacaan ADR; `test -f`.

Daftar file `docs/` dan jumlah baris (snapshot sebelum laporan ini):

| File | Baris |
|---|---:|
| `adr/0001-pilihan-stack.md` | 184 |
| `adr/0002c-azure-conventions.md` | 137 |
| `adr/0005-api-healthcheck-probe.md` | 56 |
| `adr/README.md` | 30 |
| `adr/template.md` | 36 |
| `domain-ai.md` | 337 |
| `domain.md` | 3,187 |
| `notes/README.md` | 16 |
| `notes/docker-internals.md` | 1,390 |
| `progress/README.md` | 45 |
| `progress/evidence/EXAMPLE-format.txt` | 16 |
| `progress/evidence/README.md` | 25 |
| `progress/evidence/EXAMPLE-format.txt` | 16 |
| `progress/evidence/README.md` | 25 |
| `progress/evidence/week-00/.gitkeep` | 0 |
| `progress/evidence/week-00/azure-account-show.txt` | 9 |
| `progress/evidence/week-00/env-gitignore.txt` | 7 |
| `progress/evidence/week-00/make-todo.txt` | 37 |
| `progress/evidence/week-00/ssh-key-permissions.txt` | 8 |
| `progress/evidence/week-00/wsl-resources.txt` | 13 |
| `progress/evidence/week-00/wslconfig.txt` | 9 |
| `progress/evidence/week-01/go-test-report-rerun.txt` | 13 |
| `progress/evidence/week-01/go-test.txt` | 15 |
| `progress/evidence/week-01/graceful-shutdown.txt` | 15 |
| `progress/evidence/week-01/healthz-curl.txt` | 33 |
| `progress/evidence/week-01/no-println.txt` | 10 |
| `progress/evidence/week-01/report-source-git-log.txt` | 16 |
| `progress/evidence/week-01/ssh-keyonly.txt` | 17 |
| `progress/evidence/week-02/api-probe-healthcheck.txt` | 13 |
| `progress/evidence/week-02/api-probe-runtime-build.txt` | 67 |
| `progress/evidence/week-02/build-context-after.txt` | 71 |
| `progress/evidence/week-02/build-context-before.txt` | 68 |
| `progress/evidence/week-02/buildx-multiarch.txt` | 113 |
| `progress/evidence/week-02/busybox-layer-history.txt` | 46 |
| `progress/evidence/week-02/busybox-runtime-build.txt` | 79 |
| `progress/evidence/week-02/busybox-vs-api-probe-size.txt` | 12 |
| `progress/evidence/week-02/docker-build-current.txt` | 76 |
| `progress/evidence/week-02/docker-inspect-user.txt` | 11 |
| `progress/evidence/week-02/docker-size-api-probe.txt` | 12 |
| `progress/evidence/week-02/go-test-healthcheck.txt` | 13 |
| `progress/evidence/week-02/healthcheck-distroless.txt` | 12 |
| `progress/evidence/week-02/image-size.txt` | 14 |
| `progress/explain/.gitkeep` | 0 |
| `progress/explain/week-00.txt` | 17 |
| `progress/explain/week-01.txt` | 29 |
| `progress/quiz/week-00.md` | 106 |
| `progress/quiz/week-01.md` | 128 |
| `progress/quiz/week-02.md` | 68 |
| `progress/sessions/.gitkeep` | 0 |
| `progress/week-00.md` | 262 |
| `progress/week-01.md` | 189 |
| `progress/week-02.md` | 193 |
| `roadmap.md` | 1,188 |
| `runbook/incident.md` | 32 |
| `setup/azure-day-0.md` | 75 |

| ADR | Status dalam file | Decision/Keputusan |
|---|---|---|
| `0001-pilihan-stack.md` | `Proposed` | Bagian `## Keputusan` dan `## Konsekuensi` ada dan berisi teks. |
| `0002c-azure-conventions.md` | `Accepted` | `## 3. Decision` dan `## 4. Consequences` ada dan berisi teks. |
| `0005-api-healthcheck-probe.md` | `Proposed` | Tidak ada bagian Decision/Consequences; file menyatakan sengaja menunggu keputusan manusia. |

Heading level 1 di `docs/domain.md`: `Domain Model — campus-lms`; `1. Aktor dan perannya`; `2. Boundary multi-tenant dan identity`; `3. Entitas inti dan relasinya`; `4. Academic reference dari SIAKAD`; `5. RPS dan capaian pembelajaran`; `6. Struktur pembelajaran`; `7. Materials dan file`; `8. Learning activities`; `9. Assignment dan submission`; `10. Group learning`; `11. Quiz dan exam`; `12. Rubrics`; `13. Gradebook`; `14. Grade change dan audit`; `15. Attendance`; `16. Activity completion dan progress`; `17. Discussion forum`; `18. Announcements`; `19. Copy/import course antarsemester`; `20. Integration boundary`; `21. SPADA dan PDDikti`; `22. SCORM, xAPI, dan LTI`; `23. Audit logs`; `24. Aturan bisnis yang tidak boleh dilanggar`; `25. RLS classification`; `26. Authorization hierarchy`; `27. Lifecycle penting`; `28. ERD sederhana`; `29. Invariant database yang harus dipaksakan`; `30. Prinsip waktu`; `31. Prinsip keamanan`; `32. Prinsip reporting dan learning analytics`; `33. Future extensibility`; `34. Ringkasan ownership`; `35. Keputusan domain yang telah dikunci`; `36. Dasar desain yang diverifikasi`; `37. Implementation Tiers`; `38. Amendemen terhadap keputusan yang telah dikunci`.

| Pemeriksaan | Hasil |
|---|---|
| Section 37 dan 38 domain | Ada, pada baris 2957 dan 3035. |
| `docs/domain-ai.md` | Ada (337 baris). |
| `docs/domain-amendments.md` | Tidak ada. |
| `docs/notes/docker-internals.md` | Ada (1,390 baris), tetapi untracked menurut status Git. |

## H. Progress & bukti

Perintah: `find docs/progress ...`; `rg` heading/section/signature; pembacaan report/quiz; `find docs/progress/evidence ... awk ...`.

| Laporan | Section “Belum Terverifikasi” | Tanda tangan |
|---|---|---|
| `week-00.md` | Ada dan terisi (7 baris isu dalam tabel). | `FfFfFfFf`, 16 Agustus 2026. |
| `week-01.md` | Ada dan terisi (6 isu serta 2 asumsi). | `fFfFfF`, 17 Agustus 2026. |
| `week-02.md` | Ada dan terisi (5 isu serta 1 asumsi). | Kosong (`______________`). |

| Quiz | Soal | Jawaban manusia | Skor |
|---|---:|---|---|
| `week-00.md` | 7 | Terisi | 86/100 = 86% |
| `week-01.md` | 10 | Terisi | 70/100 = 70% |
| `week-02.md` | 3 placeholder (dokumen meminta dilanjutkan hingga 8–12) | Belum diisi | Placeholder, tidak ada skor nyata |

`docs/progress/explain/` berisi `.gitkeep`, `week-00.txt`, dan `week-01.txt`; tidak ada explain Week 02.

Semua file di `docs/progress/evidence/**` — `CLAIM | COMMIT | EXIT` (`<tidak ada>` berarti field memang tidak ditemukan):

| File | CLAIM | COMMIT | EXIT |
|---|---|---|---|
| `EXAMPLE-format.txt` | Makefile has a working help target | `no-git` | `0` |
| `README.md` | `<klaim satu kalimat>` | `<git sha>` | `<exit code>` |
| `week-00/.gitkeep` | `<tidak ada>` | `<tidak ada>` | `<tidak ada>` |
| `week-00/azure-account-show.txt` | Azure CLI is authenticated to the Azure for Students subscription and it is enabled/default. | `a0d2890` | `<tidak ada>` |
| `week-00/env-gitignore.txt` | The local .env file is ignored by Git. | `a0d2890` | `<tidak ada>` |
| `week-00/make-todo.txt` | The repository Makefile todo target runs successfully in the Week 0 development environment. | `a0d2890` | `<tidak ada>` |
| `week-00/ssh-key-permissions.txt` | The Azure SSH private key is protected with owner-only permissions. | `a0d2890` | `<tidak ada>` |
| `week-00/wsl-resources.txt` | WSL2 runtime reflects the configured memory, swap, and CPU limits. | `a0d2890` | `<tidak ada>` |
| `week-00/wslconfig.txt` | WSL2 is configured with memory=5GB, processors=8, and swap=8GB. | `a0d2890` | `<tidak ada>` |
| `week-01/go-test-report-rerun.txt` | Unit tests pass with the Go race detector during weekly-report preparation | `73f24f7` | `0` |
| `week-01/go-test.txt` | Unit tests pass with the Go race detector | `4693922` | `0` |
| `week-01/graceful-shutdown.txt` | SIGTERM triggers graceful shutdown and exits with code 0 | `4693922` | `0` |
| `week-01/healthz-curl.txt` | GET /healthz returns HTTP 200 and the API emits JSON logs | `4693922` | `0` |
| `week-01/no-println.txt` | No unstructured stdout printing (fmt.Print*/log.Print*) in apps/api | `eda3ba3` | `0` |
| `week-01/report-source-git-log.txt` | Seven commits exist in the Week 1 report period | `73f24f7` | `0` |
| `week-01/ssh-keyonly.txt` | SSH localhost: key auth succeeds, password auth is rejected | `e7db51a` | `0` |
| `week-02/api-probe-healthcheck.txt` | API binary healthcheck is healthy without BusyBox and the API process runs as UID 65532 | `3bb7565` | `0` |
| `week-02/api-probe-runtime-build.txt` | Final API probe runtime builds successfully without BusyBox | `3bb7565` | `0` |
| `week-02/build-context-after.txt` | Docker build context size after adding apps/api/.dockerignore, measured with a fresh Buildx builder | `8ac8e6b` | `0` |
| `week-02/build-context-before.txt` | Docker build context size before adding apps/api/.dockerignore, measured with a fresh Buildx builder | `ad1bc11` | `0` |
| `week-02/buildx-multiarch.txt` | API image builds for linux/amd64 and linux/arm64 without a registry push | `ad1bc11` | `0` |
| `week-02/busybox-layer-history.txt` | BusyBox runtime adds a 1.22 MB COPY layer according to Docker layer history | `3bb7565` | `0` |
| `week-02/busybox-runtime-build.txt` | BusyBox comparison runtime builds successfully from the API healthcheck commit | `3bb7565` | `0` |
| `week-02/busybox-vs-api-probe-size.txt` | Exact image sizes compare the BusyBox runtime and API probe runtime from the same healthcheck commit | `3bb7565` | `0` |
| `week-02/docker-build-current.txt` | Current API Dockerfile builds successfully with the restricted build context | `8ac8e6b` | `0` |
| `week-02/docker-inspect-user.txt` | API container is configured and runs as a non-root UID | `ad1bc11` | `0` |
| `week-02/docker-size-api-probe.txt` | make docker-size confirms the current API probe image is under the 25 MB target | `3bb7565` | `0` |
| `week-02/go-test-healthcheck.txt` | API healthcheck probe unit tests pass with the Go race detector | `3bb7565` | `0` |
| `week-02/healthcheck-distroless.txt` | Distroless API container is healthy and runs as a non-root UID through the copied BusyBox probe | `8ac8e6b` | `0` |
| `week-02/image-size.txt` | Current API probe image is below the 25 MB target and remains comparable with BusyBox and Alpine targets | `3bb7565` | `0` |

## I. Status DoD

Sumber: `sed -n '500,578p' docs/roadmap.md`, kode saat ini, dan file evidence yang dibaca. Kolom terakhir adalah **USULAN audit**, bukan checkbox dan bukan keputusan manusia.

| Minggu | Item DoD | Ada bukti? | File bukti | Status menurut audit (USULAN) |
|---|---|---|---|---|
| 1 | Explain-back DNS → TCP → TLS → HTTP, rekaman 3 menit | Ada file explain Week 01; isi/kecukupan 3 menit tidak diukur saat audit. | `docs/progress/explain/week-01.txt` | Sebagian / durasi TIDAK TERVERIFIKASI |
| 1 | `/healthz` 200 + JSON log trace-ready | Ada evidence exit 0 pada commit Week 1. | `evidence/week-01/healthz-curl.txt` | Kandidat terpenuhi secara historis |
| 1 | SIGTERM graceful + exit 0 | Ada evidence exit 0 pada commit Week 1. | `evidence/week-01/graceful-shutdown.txt` | Kandidat terpenuhi secara historis |
| 1 | Minimal 15 conventional commits | 19 commit total, tetapi 3 tidak sesuai pola yang dinyatakan AGENTS. | Git log audit | Belum menurut pola AGENTS (16 patuh) |
| 1 | SSH localhost key-only | Ada evidence exit 0. | `evidence/week-01/ssh-keyonly.txt` | Kandidat terpenuhi secara historis |
| 2 | `make up` membentuk stack sehat; web → api via service name | Tidak ada service Compose; `make up` tidak dijalankan sesuai batas audit. | Tidak ada | Belum / kondisi file menghalangi |
| 2 | Image API <25MB dan UID non-root | Ada evidence historis exit 0; Docker daemon saat audit tidak dapat diakses. | `evidence/week-02/docker-size-api-probe.txt`; `api-probe-healthcheck.txt` | Kandidat terpenuhi historis; kondisi kini TIDAK TERVERIFIKASI |
| 2 | Build amd64/arm64 dan push GHCR | Build lokal didokumentasikan; evidence menyatakan tanpa registry push. | `evidence/week-02/buildx-multiarch.txt` | Sebagian — push belum dibuktikan |
| 2 | Limit 256MB dan demonstrasi OOM kill | Tidak ada Compose/limit atau evidence OOM. | Tidak ada | Belum |
| 2 | ADR distroless + trade-off | ADR-0005 ada, tetapi Proposed dan tanpa Decision/Consequences. | `docs/adr/0005-api-healthcheck-probe.md` | Sebagian — keputusan manusia belum ada |

## J. Temuan

| Keparahan | Temuan | Bukti/perintah |
|---|---|---|
| blocker | Compose utama dan override hanya berisi TASK BRIEF/TODO, tanpa `services`. Karena itu `make up`, networking antar-service, `db-shell`, dan target DoD Compose tidak punya stack untuk dijalankan. Ini langsung menghambat Minggu 2 hari Compose dan fondasi Postgres Minggu 3. | `sed deploy/compose/{docker-compose.yml,docker-compose.override.yml}`; `find deploy/compose` |
| blocker | Sesi audit tidak memiliki izin ke Docker daemon. Image/container saat ini serta Compose tidak dapat diverifikasi atau dijalankan oleh sesi ini. | `docker version`; `docker images`; `docker ps -a` semuanya menghasilkan `permission denied ... /var/run/docker.sock`. |
| blocker | Belum ada migration nyata (`apps/api/migrations/.gitkeep` saja), repository DB (`.gitkeep`), atau Compose Postgres. RLS Minggu 3 belum memiliki tempat implementasi maupun runtime lokal. | `find apps/api/migrations`; `find apps/api/internal/repository`; Compose audit. |
| should-fix | Roadmap Minggu 2 mensyaratkan `apps/web/Dockerfile`, core `api, web, postgres, redis`, dan target `make build-arm`; repo hanya punya `apps/web/README.md`, Compose task brief menunda web ke Minggu 4, dan Makefile menyediakan `buildx` (bukan `build-arm`). Ini konflik roadmap vs task brief/status commit. | `sed 500,578p docs/roadmap.md`; `find apps/web`; `sed Makefile`; `git log`. |
| should-fix | Target `buildx` memuat placeholder literal `ghcr.io/<username-github-mu>/...`. Ini bukan nilai registry final; shell/registry behavior belum diuji pada audit. | `Makefile` baris target `buildx`. |
| should-fix | `docs/adr/README.md` menandai 0001 dan 0002c sebagai TODO, sedangkan file 0001 berstatus Proposed dengan keputusan terisi dan 0002c berstatus Accepted dengan Decision/Consequences terisi. Indeks ADR tidak konsisten dengan file ADR. | `sed docs/adr/*.md`; `docs/adr/README.md`. |
| should-fix | Report Week 02, quiz Week 02, evidence Week 02, dan `docker-internals.md` masih untracked; laporan Week 02 belum ditandatangani dan quiz masih template. | `git status --short`; `sed docs/progress/week-02.md`; `sed docs/progress/quiz/week-02.md`. |
| should-fix | Enam evidence Week 00 memiliki CLAIM/COMMIT tetapi tidak memiliki field `EXIT`, tidak sesuai format evidence yang ditunjukkan template/Makefile. | pembacaan semua evidence dengan `awk` untuk CLAIM/COMMIT/EXIT. |
| catatan | `go.mod` masih memuat TODO mengganti username, namun module dan remote sudah sama-sama memakai `farisakbar28`. TODO tampak stale; kepemilikan username tidak dapat dipastikan audit ini. | `sed apps/api/go.mod`; `git remote -v`. |
| catatan | CI/CD, Caddy, backup/restore masih placeholder. Jadwal TODO-nya Week 4/5/9 sehingga bukan blocker langsung Compose/DB, tetapi belum dapat dianggap deliverable selesai. | `rg -n TODO...` |
| catatan | Ada 3 subject commit di luar pola prefix yang tercantum di AGENTS; total memenuhi angka 15 bila semua commit dihitung, tetapi tidak jika “format conventional” memakai pola proyek tersebut. | `git log` + `awk` audit. |

## K. Pertanyaan terbuka

| Pertanyaan yang membutuhkan keputusan manusia | Alasan tidak dapat dipastikan sendiri |
|---|---|
| Apakah scope Week 2 yang benar harus mengikuti roadmap (core mencakup web) atau TASK BRIEF/commit terbaru (web ditunda Minggu 4)? | Kedua sumber menyatakan scope berbeda; audit tidak berwenang memilih spesifikasi. |
| Siapa yang akan mengisi Decision dan Consequences ADR-0005, dan apakah status ADR-0001 perlu diselaraskan dengan isi Decision yang sudah ada? | Kebijakan repo menyatakan keputusan ADR milik manusia. |
| Apakah akses Docker daemon harus tersedia bagi user/sesi kerja ini sebelum implementasi Compose dilanjutkan? | Gejala permission tercatat, tetapi penyebab dan kebijakan akses host tidak dapat disimpulkan dari repo. |
| Apakah perubahan uncommitted/untracked saat audit ini memang work-in-progress yang ingin dipertahankan sebagai satu perubahan Minggu 2? | Audit hanya membaca status; tidak dapat menentukan niat pemilik. |
| Apakah `farisakbar28` adalah username pemilik yang dimaksud TODO `go.mod`? | Remote memakai nilai itu, tetapi identitas pemilik merupakan keputusan/fakta manusia. |

## Lampiran — output mentah perintah utama

### Lingkungan

```text
go version go1.23.6 linux/amd64
Client: Docker Engine - Community
 Version:           29.7.2
 API version:       1.55
 Go version:        go1.26.5
 Git commit:        a7dcaa6
 Built:             Wed Aug  5 18:28:40 2026
 OS/Arch:           linux/amd64
 Context:           default
permission denied while trying to connect to the docker API at unix:///var/run/docker.sock
Docker Compose version v5.5.0
github.com/docker/buildx v0.36.1 1d8dde89b8aba914e05e45366770736fea1fd690
GNU Make 4.4.1
git version 2.53.0
Linux HYPE 6.6.87.2-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Thu Jun 5 18:30:46 UTC 2025 x86_64 GNU/Linux
PRETTY_NAME="Ubuntu 26.04 LTS"
               total        used        free      shared  buff/cache   available
Mem:           4.8Gi       1.8Gi       2.5Gi        10Mi       747Mi       3.0Gi
Swap:          8.0Gi       1.8Mi       8.0Gi
Failed to connect to system scope bus via local transport: Operation not permitted (consider using --machine=<user>@.host --user to connect to bus of other user)
systemd 259 (259.5-0ubuntu3.4)
920K	.
```

### Git status dan ref

```text
origin  https://github.com/farisakbar28/campus-lms.git (fetch)
origin  https://github.com/farisakbar28/campus-lms.git (push)
master
19
 M Makefile
?? docs/notes/docker-internals.md
?? docs/progress/evidence/week-02/
?? docs/progress/quiz/week-02.md
?? docs/progress/week-02.md
master ae32c04627e18a63164f712725d193e975e9804c origin/master
origin/master a8c99bf83bce677dcff71eb1a7bcd95c55d0cd92
5	0
```

### Akses daemon Docker

```text
$ docker images
permission denied while trying to connect to the docker API at unix:///var/run/docker.sock

$ docker ps -a
permission denied while trying to connect to the docker API at unix:///var/run/docker.sock
```

### Status Compose

```text
# TODO Week 2
```

Baris di atas adalah satu-satunya konten non-TASK-BRIEF pada `deploy/compose/docker-compose.yml`; tidak ada `services:` atau definisi service.
