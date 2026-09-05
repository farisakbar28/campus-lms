# Go API Engineering Instructions

These instructions apply to the `apps/api` Go module:
`github.com/farisakbar28/campus-lms/apps/api`.

## Structure and dependencies

- `cmd/api` contains process entrypoint and dependency wiring only.
- `internal/config` loads and validates environment-backed configuration.
- `internal/http` contains standard-library `net/http` handlers, routing, and
  transport DTOs.
- `internal/domain` contains transport- and database-independent entities and
  rules.
- `internal/repository` owns PostgreSQL access through `pgx` and explicit SQL.
- `internal/middleware` owns trusted request identity and cross-cutting HTTP
  boundaries.
- `migrations` contains versioned SQL.

Keep domain code independent of HTTP and database packages. Inject I/O
dependencies explicitly; do not introduce global mutable state or ORM magic.

## Go and HTTP rules

- Pass `context.Context` as the first parameter for I/O and honor cancellation.
- Wrap errors with operation context using `%w`.
- Return stable, safe error codes/messages at HTTP boundaries and log internal
  details without exposing them to callers.
- Use `log/slog` with structured fields. Never log secrets, tokens, full
  request bodies, or student personal data.
- Keep `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` on
  the HTTP server.
- Handle SIGINT/SIGTERM with bounded graceful shutdown and close injected
  resources.
- `/healthz` reports process liveness; `/readyz` reports dependency readiness.

## Authentication and tenancy

Tenant identity must come from a trusted authenticated principal. Never accept
`tenant_id` from a request header, query parameter, or request body as an
authorization input.

For tenant-scoped database work:

1. authorize the requested object and tenant in application code;
2. use a transaction and `SET LOCAL app.tenant_id` for PostgreSQL RLS;
3. retain RLS as a defense-in-depth backstop rather than relying only on
   application filters.

The repository currently contains access-token verification, bearer parsing,
and refresh-session primitives. Production authentication composition is not
complete: the running server does not yet provide the full Principal/auth
endpoint wiring. Do not paper over this gap with tests or documentation; track
it as product work.

## Database changes and testing

- Migrations are forward/versioned and already-applied migration files are
  immutable. Add a new migration for a schema change.
- Use explicit SQL and column lists.
- Repository/integration tests use real PostgreSQL through Testcontainers where
  applicable; mocks must not replace RLS or constraint coverage.
- Run the race detector with `go test -race`.
- Test handlers with `httptest`.
- Add query-count or equivalent coverage for list paths to prevent N+1
  regressions.

Secrets and student personal data must never appear in source, logs, test
output, HTTP errors, or committed configuration.
