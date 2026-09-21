# Contributing

## Normal workflow

Start from an updated `master`, create a short-lived descriptive branch, and
make one coherent change. Follow
[the repository workflow](docs/engineering/workflow.md): orient and
brainstorm, have a human choose the task, plan it, obtain an independent plan
review and human approval, implement, obtain an independent implementation
review, then create a Pull Request only when the human explicitly asks for it.
CI and human review precede a human squash merge. History updates use the
same branch, Pull Request, CI, and human-merge boundary.

AI-assisted work follows the same workflow and evidence model as human work.
Agents inspect actual source and native GitHub/Git/CI records, never invent
results, and never create a parallel task archive or command-receipt system.
Human task choice, plan approval, merge, and production approval remain
mandatory. A material change returns the work to planning and fresh review.

## Branch names

Use a simple engineering category:

- `feat/<topic>`
- `fix/<topic>`
- `chore/<topic>`
- `docs/<topic>`

Do not use personal-development branch prefixes, schedule-based names, or
other non-engineering workflow names.

## Commits

Use Conventional Commits such as `feat:`, `fix:`, `chore:`, `docs:`,
`refactor:`, `test:`, `build:`, `ci:`, `perf:`, and `style:`. Keep one logical
concern per final change.

## Verification

Run the checks that match the changed scope:

- `make test` runs Go tests with the race detector; integration tests require
  the available PostgreSQL/Testcontainers environment.
- `make build` compiles the API binary.
- `make lint` is applicable only when `golangci-lint` is already available;
  this repository does not install tools automatically.
- `make help` confirms the available command surface.
- Use the relevant Docker or configuration verification script when its
  prerequisites are available.

Configure required checks only from check identities observed in successful
GitHub Actions runs; do not guess hosted check names from local workflow
configuration.

## Dependencies and architecture

Every dependency must solve a concrete problem that existing or standard
library capability cannot reasonably solve. Check its maintenance and
compatible license, and consider removal and operational cost before adding
it.

Meaningful architecture decisions belong in an ADR. Keep implementation and
documentation aligned with actual source and runtime behavior.

## Claims and review

Do not claim functionality, performance, runtime, security, cloud state, or
cost state without verification appropriate to that claim. Review the actual
command output, hosted CI result, and complete diff. Local checks do not prove
cloud state, branch protection, or production readiness.

Before review, inspect every changed file, ensure no secret or personal data is
present, and confirm that behavior changes have focused tests.
