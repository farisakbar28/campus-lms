# Rule 40 — CI/CD and Testing

## Pipeline order (cheap and fast first)

1. lint (golangci-lint, gofmt, ruff, tsc)
2. unit tests (`-race`)
3. integration tests (testcontainers, real Postgres)
4. build (buildx, multi-arch)
5. security scan (Trivy on image; gosec + govulncheck on source)
6. push to GHCR tagged with git SHA

Target total: **under 8 minutes**. Slow pipelines get bypassed, and bypassed pipelines are worse than no pipelines.

## Requirements

- `concurrency` group so a new push cancels the stale run.
- Cache Go modules and build cache.
- Pin third-party actions by commit SHA, not by tag — tags are mutable and this is a supply-chain vector.
- Minimal `permissions:` per job; `packages: write` only on the publishing job.
- Never echo a secret. Never print full env.

## Testing philosophy

- Test behaviour, not implementation. A test that breaks on every refactor is a liability.
- Integration tests use real dependencies via testcontainers. Mocking Postgres hides constraint, transaction, and RLS bugs — the exact class of bug that matters here.
- Coverage is a diagnostic, not a target. Aim for meaningful coverage in `internal/domain` and `internal/repository`; ignore generated code.
- Write down what you deliberately do **not** test and why (`docs/notes/testing-strategy.md`). Interviewers ask this.

## Database migrations in CI

Run migrations against an ephemeral database in the pipeline. From Week 5, use Neon branches for per-PR preview databases — copy-on-write branching makes this nearly free and is a strong talking point.

## AI evaluation gate (from Week 9)

`ai-eval.yml` runs on PRs touching `apps/ai/**` or `prompts/**`:

- Run the golden dataset through the pipeline.
- Compute ragas metrics.
- **Fail the PR** if faithfulness drops more than 5% versus the stored baseline, or any critical case regresses.
- Post a comparison table as a PR comment.

This is the highest-signal artefact in the whole repository. Treat it accordingly.
