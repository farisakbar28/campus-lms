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
implementation snapshot, [the roadmap](docs/roadmap.md) for future direction,
completed progress, and durable material changes, and GitHub Issues for the
operational backlog. Tests and CI provide automated evidence.

Follow [the engineering workflow](docs/engineering/workflow.md) and the
repository skills in `.agents/skills/`. Orient and brainstorm with the human;
the human chooses the task. The agent may then create or update its canonical
Issue and propose a plan. A human directly approves, revises, or rejects the
plan. There is no independent plan review. Implementation starts only after
human plan approval on a short-lived branch, followed by tests and an
implementation review independent from the implementer.

Human task selection authorizes task-scoped Issue creation or updates. After
plan approval, agents may make normal task-scoped Issue and PR creations or
updates, comments, branches, commits, and pushes without separate approval
for each action. Finalize roadmap changes and obtain a passing independent
review of the final candidate before pushing or opening a PR. Observe CI for
the PR head; only a human merges. A closing reference in the PR lets the
Issue close on merge.

If requirements, architecture, security, scope, acceptance, or production
impact materially changes, stop affected work and ask the human to decide;
revise the plan and Issue as needed before continuing. Relevant changes after
implementation review require a fresh review. Agents must not author human
task choice, plan approval, material-deviation approval, residual-risk
acceptance, merge, or production approval.

After merge, production-impacting work follows the separate
[readiness checklist](docs/engineering/production-readiness.md) and needs
human approval for the exact revision or immutable artifact, target
environment, and scope before deployment. A material change invalidates that
approval. Destructive database or data operations and external or cloud
mutations require explicit task authority and human confirmation. Never
automatically merge or deploy production.

When changing this workflow, the authority rules on the default branch
govern that migration until its Pull Request is human-merged.

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
