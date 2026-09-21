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

## Start here

- [Documentation map](docs/README.md) — what we build, how it is designed,
  security boundaries, engineering workflow, current architecture, and future
  direction.
- [Current architecture](docs/architecture.md) — implementation facts,
  accepted decisions, known gaps, and future work.
- [Domain contract](docs/domain.md) and [AI contract](docs/domain-ai.md) —
  ownership, tenancy, authorization, and future staff-side AI boundaries.
- [Architecture decisions](docs/adr/README.md) — accepted and proposed ADRs.
- [Engineering history](docs/engineering/history.md) — concise, verified
  milestones and changed assumptions.
- [Roadmap](docs/roadmap.md) — future product and engineering direction.
- [Security requirements](SECURITY.md) — enduring security rules and known
  limitations.
- [Engineering workflow](docs/engineering/workflow.md) and
  [agent playbook](docs/engineering/agent-playbook.md) — how humans and
  agents plan, review, implement, and hand off work.

The architecture and history documents distinguish repository facts from
future plans and external state. Hosted CI status must be read from the
current GitHub Actions run, and no local command proves cloud state or
production readiness.

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
