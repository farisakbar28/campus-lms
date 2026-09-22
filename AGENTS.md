# Repository Engineering Contract

## Product identity

`campus-lms` is a multi-tenant Learning Management System for higher education.
Institutional academic systems such as SIAKAD remain external authoritative
sources wherever the domain model assigns them that responsibility.

The product and architecture sources of truth are:

- `docs/domain.md` for LMS boundaries, entities, ownership, and authorization;
- `docs/domain-ai.md` for the staff-side AI direction and AI safety boundaries;
- accepted decisions in `docs/adr/` for architecture and operational constraints.

## Cold start and engineering workflow

For every task, read [the documentation map](docs/README.md) first, then read
the applicable scoped `AGENTS.md` files before entering a narrower tree.
Inspect the actual source, Git history, GitHub state, and CI evidence before
making claims. Use [current architecture](docs/architecture.md) for the
implementation snapshot, [engineering history](docs/engineering/history.md)
for verified milestones and changed assumptions, [the roadmap](docs/roadmap.md)
for future direction, and GitHub Issues for the operational backlog.

The repository-native sequence is documented in
[the workflow](docs/engineering/workflow.md): orient, brainstorm, have a
human choose the task, plan, obtain a fresh independent plan review, obtain
direct human plan approval, implement, obtain a fresh independent
implementation review, create a Pull Request, observe CI and human review,
and have a human merge. Human task selection, plan approval, merge, and
production approval are never delegated to an agent.

After a human approves the plan, an agent may create a short-lived local
branch, edit files, and create local commits when the task authorizes them.
Push or Pull Request creation is allowed only after the independent
implementation review passes and the human explicitly asks for the PR. Never
automatically merge or deploy production.

For production-impacting work, follow the separate readiness checklist and
require human approval for the exact revision or immutable artifact, target
environment, and scope. A material change invalidates that approval.
Destructive database or data operations require explicit task authority and
human confirmation. External or cloud operations require the same two gates.
Agents must not author human approvals, accept residual risk, merge, or perform
production mutations. Issue and Pull Request mutations, pushes, and remote
branch mutations require explicit task authority and human confirmation.

## Repository-wide rules

- Make the smallest coherent change and keep the repository runnable.
- A behavior change requires a focused test.
- Read the actual source before making claims about repository state or
  capabilities.
- Never invent test, benchmark, runtime, security, or cost results.
- Prefer existing or standard-library capability when it is reasonable.
- Justify every new dependency by its concrete problem, maintenance status,
  license, and removal or operational cost.
- Use English for repository engineering surfaces.
- Use structured logging (`log/slog` in Go); never log secrets, tokens,
  request bodies, or student personal data.
- Security-sensitive configuration must fail closed when it is missing or
  invalid.
- Applied database migrations are immutable; add a new migration for a schema
  change.
- Destructive database or data operations require explicit task authority and
  human confirmation.
- Never weaken a security or acceptance check merely to obtain a passing
  result.

## Safety boundaries

- Never read or write `.env`; use `.env.example` for non-secret configuration
  shape.
- Never expose credentials, private keys, tokens, student personal data, or
  secret values in code, logs, issues, commits, or images.
- Do not rewrite published history, force-push, or delete a protected/default
  branch unless a separately authorized destructive procedure explicitly
  requires it.
- External or cloud mutations require explicit task authority and human
  confirmation.
- Do not edit an already-applied migration.
- The project constraint is zero incremental paid infrastructure or tooling.
  There is no automatic paid fallback.

## Hierarchical instructions

When working inside a narrower tree, also follow:

- `apps/api/AGENTS.md`
- `deploy/AGENTS.md`

The narrower file applies in addition to this repository-wide contract.
