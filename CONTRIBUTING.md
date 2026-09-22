# Contributing

## Normal workflow

Follow [the repository workflow](docs/engineering/workflow.md): orient and
brainstorm with the human, have the human choose a task, create or update its
canonical Issue, propose a plan, and obtain direct human plan approval. Then
start a short-lived branch from updated `master`, implement and test, obtain
an independent implementation review, finalize the roadmap, and push and
create or update the PR. Observe CI for the PR head; the human reviews and
merges. The PR closes its Issue on merge. There is no independent plan review.

AI-assisted work follows the same workflow and evidence model as human work.
Agents inspect actual source and native GitHub/Git/CI records, never invent
results, and never create a parallel task archive or command-receipt system.
The human chooses tasks, decides plans and material deviations, merges PRs,
and approves production. Normal Issue, comment, branch, push, and PR
operations within a selected and approved task need no separate approval for
each action. Material changes return to the human and a revised plan;
relevant candidate changes require renewed implementation review.

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
