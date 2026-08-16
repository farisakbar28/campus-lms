# Rule 10 — Go API (`apps/api`)

## Layout

```
cmd/api/            entrypoint only, wiring, no business logic
internal/config/    env parsing, validation, fail-fast
internal/http/      handlers, routing, request/response DTOs
internal/domain/    entities and business rules, no framework imports
internal/repository/ data access (pgx), SQL lives here
internal/middleware/ auth, tenant context, logging, recovery, rate limit
migrations/         versioned SQL
```

`internal/domain` must not import HTTP or database packages. If it does, the layering is broken.

## Non-negotiables

- **Server timeouts are mandatory:** `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`. A server without them is vulnerable to slowloris — and this is a common interview question.
- **Graceful shutdown:** catch SIGTERM/SIGINT via `signal.NotifyContext`, stop accepting connections, drain in-flight requests up to `APP_SHUTDOWN_TIMEOUT`, close DB/Redis, exit 0. Docker and Kubernetes send SIGTERM before SIGKILL; without this, every deploy cuts live requests.
- **`context.Context` is the first parameter** of anything that does I/O, and it is honoured (no ignoring cancellation).
- **`/healthz` vs `/readyz`:** liveness returns 200 while the process lives; readiness returns 200 only when dependencies are reachable. They behave differently under Kubernetes — liveness failure kills the pod, readiness failure only removes it from the load balancer.
- **No global mutable state.** Dependencies are injected explicitly.
- **No ORM magic.** Explicit SQL via `pgx`, or `sqlc`-generated code. The owner must be able to read every query that runs.

## Tenancy

Every request that touches tenant data must:

1. Resolve `tenant_id` from the authenticated principal (never from a client-supplied header or body field).
2. Open a transaction and `SET LOCAL app.tenant_id = ...` so Postgres RLS applies.
3. Never build a query that relies solely on application-level filtering for isolation.

Defence in depth: RLS is the backstop, application filtering is the fast path. Both, always.

## Testing

- Table-driven tests, `t.Parallel()` where safe.
- Handlers tested with `httptest`.
- Repository tests use **real Postgres via testcontainers** — mocks hide exactly the bugs that matter (constraints, RLS, transaction semantics).
- Race detector on in CI: `go test -race`.

## Performance habits

- Avoid N+1: batch or join. Add a test that counts queries for list endpoints.
- Use keyset pagination for large lists, not `OFFSET`.
- Reuse the `pgxpool`; do not open connections per request.
- Profile before optimising (`pprof`), and record the before/after numbers as evidence.
