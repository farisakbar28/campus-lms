# Living engineering history

This document records concise, durable milestones and changed assumptions. It
is not a command log, daily diary, task archive, or copy of Pull Requests.
Implementation facts belong in [architecture.md](../architecture.md), future
direction belongs in [the roadmap](../roadmap.md), and current hosted CI
status belongs in GitHub Actions.

## Current state

At the ENG-017 migration point on 2026-09-21, `master` was at
`3cc3a32bbf339816a98ac0c8dc865a2f64d065bb`, the merge of PR #20. The
repository contained the Go/PostgreSQL API foundation, tenant/RLS and
authentication/session primitives, deployment tooling, and the accepted
tenant-context contract described in [architecture.md](../architecture.md).
The repository was preparing the transition from the high-ceremony
`campus-lms-work/v1` harness to the simpler repository-native v2 documents;
those v2 documents become authoritative when their reviewed Pull Request is
human-merged.

The dated state above is historical evidence. It does not assert current
hosted CI, branch-protection, provider, cloud, runtime, or production state.

## Completed milestones

- **Engineering baseline — 2026-09-05.** Commit
  `92419ed6144e8d07b1392e0b931cdf32a044af77` established the tracked
  repository baseline used by the subsequent engineering work.
- **Hosted CI evidence checkpoint — 2026-09-06.** A successful hosted
  `CI / API` run was observed for the exact revision
  `92419ed6144e8d07b1392e0b931cdf32a044af77`. This is historical evidence for
  that revision and job only; current hosted status, required-check settings,
  and branch-protection state must be fetched from GitHub before relying on
  them.
- **AUTHCTX-001 durable architecture milestone — 2026-09-10 to 2026-09-20.**
  Commit `847760a92c3519c4fc0ea1c0b1ba7a1b56b345a2` recorded the trusted
  authenticated tenant-context contract, and merge commit
  `27eeecef73776ac26e47149ad9a3ab3414f14bb4` brought ADR-0006 into `master`.
  Commit `3cc3a32bbf339816a98ac0c8dc865a2f64d065bb` then removed the temporary
  AUTHCTX-001 work records in PR #20. The durable milestone is the accepted
  architecture contract, not a claim that runtime authentication composition
  is complete.
- **ENG-017 workflow-v2 transition — 2026-09-21.** GitHub Issue #21 defines
  the simpler repository-native workflow and its human approval boundary. The
  migration retires the local v1 harness while preserving independent review,
  human authority, fail-closed security, and a separate production approval
  path.

## Current gaps

Runtime authentication composition remains a separate gap from AUTHCTX-001:
the repository has token verification, Bearer parsing, session primitives,
and accepted admission semantics, but the running server does not yet provide
the complete trusted Principal/auth endpoint wiring. Other gaps are listed in
[architecture.md](../architecture.md) and sequenced in
[the roadmap](../roadmap.md), including product workflows, current-schema
recovery, production observability, and permanent ingress.

## Changed or superseded assumptions

- The accepted bounded Quick Tunnel decision supersedes the scoped custom
  domain, Origin-CA, and permanent-ingress assumptions recorded in earlier
  ingress planning; unaffected private-origin, recovery-boundary, and
  zero-incremental-cost constraints remain.
- The roadmap no longer owns a competing implementation snapshot. Current
  implementation facts belong in `docs/architecture.md`; dated CI evidence is
  historical and current hosted CI is read from GitHub Actions.
- AUTHCTX-001 established the trusted tenant-context contract, but it did not
  implement runtime authentication composition. That distinction remains
  intentional.

## How to maintain this document

Add only concise milestones that can be verified from native Git, GitHub,
tracked CI, or an explicitly cited repository source. Include a date and exact
revision when a revision is material. Describe outcomes and durable changed
assumptions, not a sequence of commands or every Pull Request detail. Never
include secrets, tokens, credentials, or student personal data. When a current
implementation fact changes, update `architecture.md`; when future direction
changes, update `roadmap.md`; when a security or architectural contract
changes, update its owning source or ADR through the normal reviewed workflow.
