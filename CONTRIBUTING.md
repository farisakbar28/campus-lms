# Contributing

## Development workflow

The development workflow is:

1. Start from an updated `master`.
2. Create a short-lived, descriptive branch.
3. Implement one coherent change.
4. Run targeted local checks.
5. Request review and address review findings.
6. Open a Pull Request.
7. Ensure the applicable CI checks pass.
8. Have the human maintainer squash-merge the change.

Changes intended for master should be developed on short-lived branches and
submitted through pull requests.

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
cost state without verification appropriate to that claim. There is no
per-task evidence-receipt system; review the actual command output and diff.

Before review, inspect every changed file, ensure no secret or personal data is
present, and confirm that behavior changes have focused tests.
