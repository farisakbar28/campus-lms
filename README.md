# campus-lms

`campus-lms` is a multi-tenant Learning Management System for higher education.
It owns learning content and activity workflows while institutional academic
systems such as SIAKAD remain authoritative for the records assigned to them
by the domain model. Tenant isolation is a core correctness and security
property, enforced through application authorization and PostgreSQL RLS.

Staff-side AI is a future differentiator for authorized, non-student
institutional users. Possible uses include permitted information retrieval,
attendance lookup, submission or response summarization, and staff
productivity assistance. AI will provide assistance, search, summaries, or
drafts only; it must never decide grades, attendance, pass/fail status,
disciplinary outcomes, or other authoritative academic results.

The exact staff role allow-list is unresolved:
`DECISION_REQUIRED: AI_ROLE_ALLOWLIST`.

No student-facing AI chat, RAG assistant, or study-planner functionality is
implemented or promised. The AI service itself is not implemented yet.

## Current status

| Area | Status | What is actually present |
|---|---|---|
| Go API foundation | Implemented | Standard-library `net/http` server, package layering, dependency wiring, and structured JSON logging. |
| Configuration | Implemented | Fail-closed environment parsing and validation for the current API. |
| Health and lifecycle | Implemented | `/healthz`, dependency-backed `/readyz`, server timeouts, and graceful shutdown. |
| PostgreSQL schema | Implemented | Versioned tenant, identity, academic, membership, enrollment, audit, and auth-session migrations. |
| Tenant/RLS foundation | Implemented | Tenant-scoped constraints, RLS policies, transaction-local context, and repository tests. |
| Authentication/session primitives | Partially implemented | Access-token verification, bearer parsing, and refresh-session lifecycle primitives exist. |
| Production authentication composition | Incomplete | The running server does not yet provide complete Principal/auth endpoint wiring. |
| Containers and local deployment tooling | Materially present | Multi-stage API/migrator image, development Compose, production Compose, health probes, and deployment scripts exist. |
| Production deployment/ingress | Incomplete and stale | Existing Caddy/Compose implementation is stale relative to the accepted bounded Quick Tunnel validation architecture and requires separately authorized implementation reconciliation; it is not production-ready. The bounded-validation decision is accepted, while permanent production ingress remains unresolved and is not implemented. |
| Backup and restore | Partial | Local backup and disposable restore tooling exists; production off-machine operations are not complete. |
| Frontend | Not implemented | No frontend application is present. |
| AI service | Not implemented | No AI application or model integration is present. |
| CI | Defined in repository | `.github/workflows/ci.yml` defines API validation for pull requests and pushes to `master`; hosted execution status must be determined from current GitHub Actions results rather than inferred from documentation. |
| CD | Not implemented | No working remote release workflow is present. |
| AI evaluation | Not implemented | No evaluation pipeline or golden dataset runner is present. |
| Full production observability | Not implemented | Structured API logs exist; the planned metrics, traces, and dashboards do not. |
| Permanent production ingress | Not implemented | The bounded validation topology is not a permanent ingress solution. |

## Architecture and domain references

- [Domain model](docs/domain.md)
- [AI domain boundaries](docs/domain-ai.md)
- [Architecture decisions](docs/adr/)

These documents describe ownership, tenant isolation, authorization boundaries,
and decisions that must remain aligned with the implementation. External
provider, cloud, runtime, and cost claims require current verification before
they are presented as facts.

## Development

Prerequisites are Go 1.23 or later and Docker with Compose v2. The integration
tests use real PostgreSQL through Testcontainers when their environment is
available.

```bash
make help
make build
make test
make lint       # only when golangci-lint is already installed
make docker-build
```

For local Compose operation, provide a private, uncommitted `.env` whose shape
matches `.env.example`, then use:

```bash
make up
make health
make logs
make down
```

The API currently uses the standard-library HTTP mux; no `chi` dependency is
claimed here. Do not treat local commands as proof of cloud deployment,
external provider state, or permanent production readiness.

## Cost policy

The project targets zero incremental infrastructure and tooling spend. The
owner's existing ChatGPT Plus/Codex subscription is outside that incremental
project-cost accounting. No external vendor free tier is claimed as currently
verified, and the application must not introduce an automatic paid fallback.

## License

MIT. See [LICENSE](LICENSE).
