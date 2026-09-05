# Repository Engineering Contract

## Product identity

`campus-lms` is a multi-tenant Learning Management System for higher education.
Institutional academic systems such as SIAKAD remain external authoritative
sources wherever the domain model assigns them that responsibility.

The product and architecture sources of truth are:

- `docs/domain.md` for LMS boundaries, entities, ownership, and authorization;
- `docs/domain-ai.md` for the staff-side AI direction and AI safety boundaries;
- accepted decisions in `docs/adr/` for architecture and operational constraints.

## Engineering workflow

The normal repository workflow is:

```text
short-lived branch → implementation → local verification → review
→ Pull Request → required CI → human squash merge into protected master
```

This workflow is the target operating model. Do not claim that branch
protection, required checks, or the complete remote workflow are active during
the migration.

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
- External or cloud mutations require explicit task authorization and human
  confirmation.
- Destructive database or data operations require explicit authorization.
- Do not edit an already-applied migration.
- The project constraint is zero incremental paid infrastructure or tooling.
  There is no automatic paid fallback.

## Hierarchical instructions

When working inside a narrower tree, also follow:

- `apps/api/AGENTS.md`
- `deploy/AGENTS.md`

The narrower file applies in addition to this repository-wide contract.
