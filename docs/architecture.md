# Current architecture overview

This is a concise snapshot of what the repository currently implements. It is
not a replacement for [the domain contract](domain.md) or the detailed
[architecture decisions](adr/README.md). Labels below are deliberate:
implementation facts come from source, migrations, tests, and tracked
deployment files; accepted decisions come from accepted ADRs; proposed
decisions are not commitments; and gaps are not implied capabilities.

## Implementation facts

### API and package boundaries

- `apps/api` is a Go module with a standard-library `net/http` transport.
- `cmd/api` loads configuration, opens the database pool, wires the roster
  repository, starts the server, handles signals, and supports the binary
  health-check probe.
- `internal/config` validates required environment-backed configuration and
  requires TLS-capable database URLs in production.
- `internal/http` owns transport handlers and exposes `/healthz`, dependency-
  backed `/readyz`, and the current course-offering participant route.
- `internal/domain` contains transport- and database-independent domain types;
  `internal/repository` owns explicit PostgreSQL queries and transactions;
  `internal/middleware` owns bearer parsing and trusted-context helpers.
- Structured `log/slog` logging, server timeouts, bounded startup, and bounded
  graceful shutdown are part of the current API foundation.

### Data and tenant isolation

- PostgreSQL is accessed through `pgx` and an injected connection pool.
- Versioned migrations currently reach `0006`, covering tenant, identity,
  academic, membership, enrollment, audit, and authentication-session data.
- Tenant-scoped data uses explicit `tenant_id`, composite consistency
  constraints, PostgreSQL RLS, and transaction-local context. Repository and
  database tests exercise authorization, RLS, constraints, and query shape.
- The application authorization layer remains responsible for active
  membership, roles, course authority, object access, and lifecycle checks;
  RLS is defense in depth rather than a substitute for authorization.

### Authentication boundary

Access-token verification, strict Bearer parsing, refresh-session lifecycle,
and a `Principal` context helper exist as primitives. The running server does
not yet compose the full user/session/membership admission boundary, typed
trusted Principal propagation, or complete authentication endpoints. This is
an implementation gap, not a reason to weaken the accepted security contract.

### Deployment and operations

The repository contains multi-stage API and migrator images, development and
production Compose definitions, health probes, local backup/restore tooling,
and deployment scripts. The production Compose/Caddy path still contains
legacy loopback `8443`, Origin-CA, and hostname-based TLS wiring. It remains
stale relative to the accepted bounded Quick Tunnel validation decision and is
not a production-readiness claim.

## Accepted architectural decisions

- [ADR-0002](adr/0002-multi-tenancy.md) accepts shared PostgreSQL schema,
  explicit tenant identity, RLS, composite constraints, and application
  authorization as defense in depth.
- [ADR-0002e](adr/0002e-zero-domain-quick-tunnel-validation.md) accepts only
  a bounded no-domain Quick Tunnel validation path and does not establish
  permanent production ingress.
- [ADR-0005](adr/0005-api-healthcheck-probe.md) accepts `/api -healthcheck`
  probing `/healthz` for the minimal API image.
- [ADR-0006](adr/0006-trusted-authenticated-tenant-context.md) accepts
  tenant-neutral access tokens, explicit untrusted tenant URL selection, and
  the fail-closed transition to a trusted Principal.

Other accepted and superseded decisions are indexed in
[docs/adr/README.md](adr/README.md). The proposed status of
[ADR-0001](adr/0001-pilihan-stack.md) remains proposed.

## Proposed or unresolved decisions

The repository does not currently select a frontend, cloud/database provider,
AI provider or model, permanent ingress design, AI role allow-list, or
institutional AI processing boundary. These choices require their own
decision and verification; this document must not turn a proposal or historic
observation into an implementation claim.

## Known gaps and future work

- Compose the accepted authentication and tenant-context contract into the
  running server and protected routes.
- Build the remaining LMS content, assessment, gradebook, attendance, and
  frontend capabilities while preserving domain ownership and authorization.
- Reconcile backup/restore validation with the current schema and add measured
  operational observability.
- Reconcile deployment ingress with the accepted bounded validation decision
  and separately decide permanent production ingress.
- Define the approved staff-side AI role/data boundary before implementing a
  narrow, optional, non-authoritative capability.

For dependency order and future classifications, see the
[roadmap](roadmap.md). For dated, verifiable milestones, see
[engineering history](engineering/history.md).
