# Temporary workflow artifacts

Contract version: `campus-lms-work/v1`

The `work/` tree contains temporary planning, review, and implementation
bindings. It is not a completed-work archive or a command/evidence receipt
system.

## Active layout

An active issue has exactly one pair:

```text
work/active/<planning-id>/WORK.md
work/active/<planning-id>/REVIEWS.md
```

Temporary Product/Phase Planning artifacts use:

```text
work/phases/active/<phase-id>/PHASE.md
work/phases/active/<phase-id>/REVIEWS.md
```

Use the matching files under `work/templates/` as the minimum schema. Existing
planning IDs, such as `AUTHCTX-001`, remain identities. Do not create duplicate
IDs, `work/completed/`, an archive, a progress journal, or a per-task evidence
receipt.

## Lifecycle

Active issue work follows:

```text
DRAFT -> PLAN_REVIEW -> APPROVED -> IMPLEMENTING
      -> IMPLEMENTATION_REVIEW -> READY_FOR_PR
      -> PR_VALIDATION -> READY_TO_MERGE
```

The current canonical Issue, repository contracts, plan hash, review records,
native Git commits, GitHub checks, and human-authored approval comments are the
evidence sources. Agents must re-fetch relevant GitHub state instead of
trusting chat or stored transcripts.

Before active files are removed, the approved decision and candidate/review
bindings must be handed to the native Issue/PR/Actions records and any required
roadmap update. Active cleanup is pre-merge and cleanup-only; the human
maintainer performs the merge. Post-merge disposition and any conditional
release procedure remain in GitHub rather than recreating local artifacts.

See [the durable workflow contract](../docs/engineering/ai-workflow.md) for
source precedence, approval schemas, severity gates, release boundaries, and
validator limitations.
