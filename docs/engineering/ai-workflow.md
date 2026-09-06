# Repository-native AI engineering workflow

Contract version: `campus-lms-work/v1`

This document is the durable process contract for AI-assisted engineering in
`campus-lms`. It governs repository-local planning and review artifacts while
preserving human authority over intent, material changes, merge, release, and
production. It does not change LMS runtime behavior or make a claim that the
current GitHub settings enforce this process remotely.

## Source precedence and boundaries

Current source and directly observed command or runtime results determine
implementation facts. The following sources retain their existing authority:

- `docs/domain.md` defines LMS entities, ownership, tenancy, and authorization.
- `docs/domain-ai.md` defines the staff-side AI boundary and safety rules.
- Accepted files under `docs/adr/` govern only their stated decisions.
- `docs/roadmap.md` is durable product and phase direction.
- GitHub Issues are the canonical operational backlog and existing planning IDs
  remain work-item identities.
- Root and narrower `AGENTS.md` files jointly constrain every change.
- An active, approved `WORK.md` is the execution contract for one Issue, but it
  cannot override the Issue, domain sources, accepted ADRs, or engineering
  contracts. `REVIEWS.md` records review verdicts and does not expand scope.

A conflict involving security, tenancy, authority, destructive work, evidence,
release, production, or scope fails closed and returns the item to planning.
ENG-016 introduces no new domain ownership, tenant behavior, authorization
behavior, AI runtime, database schema, migration, deployment behavior, or ADR.

## Product and Phase Planning

Temporary phase planning connects current contracts and implementation reality
to durable backlog work:

```text
contracts + roadmap + implementation + open issues
  -> temporary PHASE.md
  -> fresh independent phase review
  -> exact human phase approval
  -> roadmap and GitHub Issue synchronization
  -> human-reviewed roadmap PR and merge
  -> temporary phase cleanup
```

`work/phases/active/<phase-id>/PHASE.md` contains the phase ID, author and
session labels, revision and deterministic hash, goal and outcomes, dependency
order and rationale, exit criteria, major exclusions, unresolved decisions,
and proposed existing planning IDs with Issue scopes and risk. Its reviewer
must use a fresh session and differ from the author. Human approval is a
directly authored GitHub record binding the phase ID, revision, hash, and
`APPROVED`; an agent cannot author it.

Before phase artifacts are removed, the approved decision is synchronized to
`docs/roadmap.md` with the revision/hash, outcomes, ordering rationale, exit
criteria, exclusions, unresolved decisions, and exact planning-ID-to-Issue
mapping. Authorized GitHub Issue creation or updates and the roadmap PR are
external mutations; a human must authorize or perform them. There is no
permanent phase archive.

## Issue-backed work items

An active Issue uses exactly:

```text
work/active/<planning-id>/WORK.md
work/active/<planning-id>/REVIEWS.md
```

The templates in `work/templates/` are schemas, not task history. Active files
hold only the minimum identity, plan, approval, candidate, review, deviation,
verification, and completion bindings. They are not a command receipt store,
progress journal, evidence ledger, dashboard, or completed-work archive. No
`work/completed/` directory is allowed.

The normal human prompt surface is intentionally short:

```text
Plan AUTHCTX-001
Independently review the current plan for AUTHCTX-001
Implement approved AUTHCTX-001
Independently review the implementation for AUTHCTX-001
```

Each prompt resolves the current repository and GitHub state. Conversational
memory is not evidence.

## Pre-merge lifecycle

Repository-local state ends before merge:

```text
DRAFT -> PLAN_REVIEW -> APPROVED -> IMPLEMENTING
      -> IMPLEMENTATION_REVIEW -> READY_FOR_PR
      -> PR_VALIDATION -> READY_TO_MERGE
```

The gates are:

| State | Required gate | Failure transition |
| --- | --- | --- |
| `DRAFT` | Open canonical Issue, current Issue digest, complete normative plan/hash | — |
| `PLAN_REVIEW` | Different plan-review actor/session and fresh read-only review | `DRAFT` |
| `APPROVED` | Passing plan review and exact human-authored Issue approval | `DRAFT` |
| `IMPLEMENTING` | Approval and Issue digest rechecked; work stays in scope | `DRAFT` for material deviation |
| `IMPLEMENTATION_REVIEW` | Authorized candidate commit and different fresh reviewer | `IMPLEMENTING` for in-plan fixes; `DRAFT` for deviation |
| `READY_FOR_PR` | Approved implementation review bound to the exact candidate and finding gates | `IMPLEMENTATION_REVIEW` |
| `PR_VALIDATION` | Explicit push/PR authority, current Issue, PR diff, and hosted checks inspected | `IMPLEMENTING` or `DRAFT` |
| `READY_TO_MERGE` | Exact reviewed candidate or cleanup-only descendant, applicable CI observed, no blocking finding | `PR_VALIDATION` |

Human squash merge is human-only and is represented by native GitHub PR and
Actions records. The repository never creates `DONE`, merged, released, or
post-merge local states.

## Deterministic bindings

### Plan hash

`WORK.md` is UTF-8. It must contain exactly one later pair of lines:

```text
<!-- PLAN-NORMATIVE-BEGIN -->
<!-- PLAN-NORMATIVE-END -->
```

Only bytes between those lines are hashed under
`campus-lms-plan-sha256-v1`. Normalize CRLF and lone CR to LF, remove trailing
ASCII spaces and tabs from every line, remove leading and trailing blank lines,
append exactly one LF to non-empty content, then SHA-256 the UTF-8 bytes and
render `sha256:<64 lowercase hexadecimal>`. The normative section includes the
Issue digest, revision, goal, scope, exclusions, affected areas, plan,
acceptance, verification, documentation impact, risks, deviations, and
release choice. Mutable status, approval URL, actor labels, candidate SHA,
verification results, and completion summary remain outside it. A normative
change increments the revision and requires fresh plan review and human
approval.

### Canonical Issue digest

For the canonical Issue, normalize title and body independently with the same
line-ending, trailing-whitespace, blank-boundary, and final-LF rules. Encode
them as UTF-8 and hash:

```text
campus-lms-issue-sha256-v1 || NUL || uint64be(title length) || title bytes
  || uint64be(body length) || body bytes
```

Comments and ordinary `updatedAt` changes are excluded. A title/body mismatch
at implementation or PR validation is a planning failure, even when the
change is editorial.

### Human plan approval

The exact canonical-Issue comment body is:

```text
approval_type=PLAN
work_item_id=<ID>
plan_revision=<integer>
plan_hash=sha256:<64 lowercase hexadecimal>
decision=APPROVED
```

`WORK.md` stores only its exact URL. The comment must be directly authored by
an accepted human account; an agent may not compose or post the approval. A
fresh fetch at approval, implementation start, and PR validation compares the
URL and all four bindings. Local parsing cannot prove identity, account
integrity, informed intent, or freshness, so those remain explicit human and
remote-process trust boundaries.

## Reviews, findings, and deviations

`REVIEWS.md` is append-oriented while active. A plan record contains its ID,
type, actor/session labels, plan-author labels, fresh-session attestation,
revision, plan hash, Issue digest, findings, and verdict. An implementation
record contains the equivalent implementation-author labels and the exact
candidate Git SHA. A reviewer actor and session must differ from the relevant
author labels, and each final review must be from a fresh session. Labels and
attestations improve auditability but are not cryptographic proof.

`CRITICAL` and `HIGH` findings block. `MEDIUM` findings require resolution or
an exact human-authored residual-risk comment; `LOW` findings are normally
non-blocking. Agents cannot accept residual risk. A residual-risk comment is
bound to the work item, plan revision/hash, finding ID, candidate SHA, and a
bounded rationale.

A deviation is material when it changes scope, acceptance, verification,
dependencies, behavior, architecture, safety, authority, risk, Issue
specification, or external/destructive authority. Security-sensitive
uncertainty is material. Stop, record one concise deviation, return to `DRAFT`,
revise the normative plan and revision/hash, and obtain fresh independent plan
review and human approval. Formatting and implementation detail are
non-material only when those properties remain unchanged.

## Candidate, PR, and cleanup boundaries

An implementation reviewer checks the exact authorized native candidate commit,
its parent/base and full diff, relevant checks, current Issue digest, and
current plan binding. Any non-active-file change after review makes the review
stale. An active-only descendant is allowed only when the candidate-to-head
diff deletes exactly `work/active/<ID>/WORK.md` and `REVIEWS.md`.

Before cleanup, native Issue/PR/Actions records must carry the plan approval,
Issue digest, reviewed candidate SHA and verdict, residual-risk references,
hosted CI head/run when observed, and `release_required`. Cleanup is a
separate authorized local commit before human merge. Push, PR mutation, and
merge are not authorized by local validation.

## Conditional post-merge release procedure

Post-merge state lives in native GitHub records, never in a recreated local
work artifact. A common human-authorized Issue disposition is:

```text
record_type=POST_MERGE_DISPOSITION
work_item_id=<ID>
pull_request_url=<URL>
merged_sha=<full SHA>
release_required=true|false
decision=NO_RELEASE|RELEASE_EVALUATION_REQUIRED
recorded_at=<RFC3339 UTC>
```

When `release_required=false`, the human merges, the exact merged SHA is
re-fetched, the disposition is recorded as `NO_RELEASE`, and the Issue may be
closed by a human or explicitly authorized agent. No deployment or production
approval follows.

When `release_required=true`, a fresh authorized evaluator inspects the merge,
CI, artifacts, deployment contracts, and required external state, then records
`DEFERRED`, `NOT_RELEASABLE`, or `PROPOSED`. A proposed release scope is the
exact normalized text between unique release-scope markers and binds the work
item, merged SHA, immutable artifacts, environment, ordered mutations,
configuration references without secrets, migration/forward-recovery scope,
approved safe recovery, post-deploy checks, stop conditions, and authority for
native claim/result comments.

A separate directly human-authored production approval binds the merged SHA,
environment, release-scope hash, RFC3339 `expires_at`, and `APPROVED`. Its
native comment ID/URL, `created_at`, `updated_at`, and normalized body digest
are an exact version. It is eligible only when unedited, human-authored,
unexpired, and `created_at <= now < expires_at`; the native timestamp is
authoritative. Editing an old comment never extends authority; a correction or
expiry extension requires a new comment.

Before any claim or production mutation, the release agent re-fetches the
Issue and exact scope/approval comments with complete pagination. It fails
closed on deletion, edit, malformed or future-dated input, mismatch, expiry,
unavailable/incomplete GitHub state, ambiguous claims, or missing authority.
It creates one unique `DEPLOYMENT_CLAIM`, re-fetches to confirm exactly one
valid claim owned by that executor, revalidates immediately before the first
mutation, and creates exactly one bound `DEPLOYMENT_RESULT`. A claim consumes
the approval for that execution. A pre-mutation expiry records
`APPROVAL_EXPIRED_BEFORE_MUTATION` and performs no mutation.

The claim binds the exact approval comment ID/URL, native creation/update
timestamps, normalized body digest, work item, merged SHA, target environment,
release-scope hash, unique execution ID, and `CLAIMED`. The result binds that
same approval and execution, the claim URL, production-mutation-started flag,
`SUCCEEDED|FAILED|ABORTED`, a sanitized reason, and a verification reference.
Every claim has exactly one result. Missing, duplicate, contradictory, or
unavailable records stop the execution and require human resolution.

After mutation starts, expiry, edit, or deletion stops forward work. Only
rollback/remediation explicitly enumerated in the already approved scope may
continue when required to reach a safe state; a new forward or out-of-scope
operation requires a new approval and execution ID. Failed, partial, expired,
edited, deleted, or aborted executions never recycle their approval. GitHub
comments are not an atomic lock or compare-and-swap primitive, so truly
concurrent executors remain a process boundary; this harness adds no
coordinator, queue, daemon, lock service, database, or workflow engine.

The release agent applies this decision table before and during execution:

| Case | Required outcome |
| --- | --- |
| Valid unedited approval | Continue only after exact version, binding, digest, expiry, scope, and authority checks. |
| Edited approval or edited expiry extension | Stop; use a new human-authored approval. |
| Corrected approval in a new comment | Validate only the new exact version. |
| Future-dated or expired approval | Stop before claim/mutation; require a new approval. |
| Approval expires after claim before first mutation | No mutation; record `ABORTED` with `APPROVAL_EXPIRED_BEFORE_MUTATION`. |
| Approval valid at first mutation | Consume it for that execution and enforce forward-expiry rules. |
| Expiry after partial mutation | Stop forward work; only pre-approved safe recovery may continue; record `FAILED`/`ABORTED`. |
| Edited/deleted approval after mutation | Stop forward work; only pre-approved safe recovery may continue; record the condition. |
| Retry after failure/abort or partial deployment | Never reuse the approval; reassess state and obtain a new approval, scope hash, and execution ID. |
| Mismatched SHA, environment, or scope hash | Stop before claim/mutation. |
| Prior, multiple, foreign, or ambiguous claim | Stop before mutation and require human resolution. |
| Sole current executor claim after complete re-fetch | Continue only while every approval and authority check still passes. |
| GitHub re-fetch unavailable or incomplete | Stop before mutation; do not infer validity. |
| Operation outside approved scope or attempted forward work after expiry | Stop and obtain new explicit approval. |

## Validator and CI

`scripts/validate_ai_workflow.py` uses only the Python standard library. With
no external input it checks durable-file presence and version, active artifact
cardinality, ID/state/risk/release/hash/SHA/URL syntax, marker uniqueness and
hashes, required state bindings, reviewer actor/session inequality, structured
finding gates, and rejection of archives or post-merge states. Optional fresh
JSON inputs support Issue digest recomputation, exact plan-approval matching,
and residual-risk binding checks without storing a receipt. Optional base,
candidate, and head SHAs check that a diff is cleanup-only.

The validator cannot prove human authorship or intent, reviewer freshness or
thoroughness, semantic materiality, current branch protection, self-protection
against a malicious governance change, live production state, claim
uniqueness, deployment serialization, or production-result truth. CI is
therefore a useful local artifact check, not tamper-proof governance.

The existing `CI / API` job retains its name, triggers, permissions, and
application checks, and additionally runs `make workflow-test` and
`make workflow-check`. No new dependency or paid service is introduced.

## Human operating procedure

1. Plan from the canonical Issue and repository sources; create a temporary
   phase or active `WORK.md` with a deterministic revision/hash.
2. Obtain a fresh independent plan review, then the exact human approval
   comment. Re-fetch both before implementation.
3. Implement only approved scope, run relevant checks, and create a local
   candidate commit only when explicitly authorized.
4. Stop for a fresh independent implementation review. Resolve findings or
   return to planning; do not give the implementation its own verdict.
5. With explicit push/PR authority, validate the exact candidate and hosted CI,
   perform only authorized active cleanup, and stop at `READY_TO_MERGE`.
6. The human squash-merges and follows the applicable GitHub-backed post-merge
   release path. Remove temporary artifacts only after durable handoff.

Push, PR/Issue mutation, merge, cloud/external operations, production, and
destructive data work each require their own explicit authority. Passing a
local check grants none of those authorities.
