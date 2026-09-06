# ENG-016 — Repository-native AI engineering workflow harness

Status: `IMPLEMENTATION_REVIEW`

Issue: `ENG-016` / GitHub issue `#16`

Issue URL: <https://github.com/farisakbar28/campus-lms/issues/16>

Issue specification digest:
`sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`

Issue `updatedAt` observed: `2026-09-05T18:16:42Z` (informational only)

Plan author actor label: `Codex planner`

Plan author session label: `ENG-016-plan-r4-2026-09-06`

Plan revision: `4`

Plan hash: `sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934`

Risk: `CRITICAL`

Release required: `false`

Human plan-approval comment: `https://github.com/farisakbar28/campus-lms/issues/16#issuecomment-5556469166`

Implementation author actor label: `Codex implementer`

Implementation author session label: `ENG-016-implement-r4-2026-09-06`

Candidate Git commit SHA: `c03eeeab7a9ad7508246214f64c42b7027c627c6`

This is the approved revision-4 execution contract after Independent PLAN
Review 5 and the exact human-authored canonical-Issue approval comment. It
authorizes implementation and an explicitly authorized local candidate commit
only. It authorizes no push, PR mutation, merge, external mutation, or
deployment.

## Repository and issue checkpoint

- Base Git revision:
  `92419ed6144e8d07b1392e0b931cdf32a044af77`.
- Branch: `chore/ai-engineering-harness`.
- GitHub issue #16 was re-fetched for revision 4. Its title and body produce the
  issue specification digest above under
  `campus-lms-issue-sha256-v1`. Comments and `updatedAt` are deliberately
  excluded.
- Independent PLAN Review 1 is preserved unchanged in `REVIEWS.md` with
  verdict `CHANGES_REQUIRED` and findings `ENG-016-PLAN-001` through
  `ENG-016-PLAN-008`.
- `CONTRIBUTING.md` explicitly rejects a per-task evidence-receipt system.
  The revised plan therefore uses the active files only for the minimum execution
  and review bindings; agents rerun relevant commands and inspect native
  Git/GitHub state instead of maintaining a parallel evidence ledger.
- The current tracked ADR index has no ADR-0006. This plan does not reserve an
  ADR number and does not propose an ADR for ENG-016.
- A successful hosted run for the base SHA was observed as workflow/job
  `CI / API`. Current branch-protection observations recorded by Review 1 are
  external state, not a repository guarantee, and ENG-016 will not change
  branch protection.

## Normative plan

Only the content between the following exact marker lines is normative and
covered by the plan hash. Status, approval reference, candidate identity,
review records, verification results, deviations, and completion summary are
mutable and excluded.

<!-- PLAN-NORMATIVE-BEGIN -->

### Identity and input checkpoint

- Schema: `campus-lms-work/v1`.
- Work-item ID: `ENG-016`; no duplicate `W-0001`-style identity is allowed.
- Canonical issue: GitHub issue `#16`,
  <https://github.com/farisakbar28/campus-lms/issues/16>.
- Issue digest algorithm: `campus-lms-issue-sha256-v1`.
- Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`.
- Base Git revision:
  `92419ed6144e8d07b1392e0b931cdf32a044af77`.
- Plan revision: `4`.
- Plan hash algorithm: `campus-lms-plan-sha256-v1`.
- Risk: `CRITICAL`, because this changes governance, review/approval rules,
  CI integration, and release-control semantics.
- Release required: `false`; ENG-016 creates no deployable LMS artifact and
  authorizes no release or production mutation.

The canonical issue checkpoint is calculated from its title and body:

1. Fetch the issue number, title, body, URL, and `updatedAt` from GitHub.
2. Decode title and body as Unicode and reject invalid Unicode.
3. For each field independently, convert CRLF and lone CR to LF, remove
   trailing ASCII spaces/tabs from each line, remove leading/trailing blank
   lines, and append exactly one LF when the normalized field is non-empty.
4. Encode both normalized fields as UTF-8.
5. Form the byte sequence
   `campus-lms-issue-sha256-v1`, NUL, the normalized-title byte length as an
   unsigned 64-bit big-endian integer, title bytes, the normalized-body byte
   length in the same form, and body bytes.
6. SHA-256 that sequence and store `sha256:<64 lowercase hexadecimal>`.

`updatedAt` is informational because ordinary comments and unrelated issue
metadata must not stale a plan. Before implementation and again during PR
validation, re-fetch and recompute the title/body digest. A mismatch blocks
progress. Compare the specifications; any material title/body change requires
`DRAFT`, a revised normative plan, a new plan revision/hash, independent plan
review, and new human approval. Under this conservative v1 contract, even a
normalized editorial title/body change requires refreshing the normative
checkpoint and approval; comments alone never change the digest.

### Goal

Install a lean, repository-native, tool-agnostic workflow contract supporting:

- temporary Product/Phase Planning followed by durable roadmap and GitHub Issue
  synchronization;
- issue-backed planning, independent plan review, human plan approval,
  implementation, independent implementation review, PR validation, and human
  merge;
- deterministic plan and issue-specification bindings;
- conditional GitHub-backed post-merge release handling; and
- short prompts that agents resolve from current repository and GitHub state,
  not conversational memory.

### Scope

In scope:

- one durable engineering workflow contract;
- minimal active `WORK.md`, `REVIEWS.md`, and phase templates;
- one small standard-library-only validator plus focused tests;
- concise routing and human procedure updates;
- integration of validator tests/checks into the existing hosted CI job;
- a concise Product/Phase Planning summary contract in the roadmap; and
- explicit procedural limitations where local files and current CI cannot
  prove human intent, reviewer freshness, or self-protect governance controls.

Out of scope:

- LMS/API behavior, authentication, authorization, tenant context, database
  schema/data, migrations, deployment topology, cloud state, or production;
- branch-protection, repository ruleset, CODEOWNERS, GitHub App, bot, or remote
  approval-policy mutation;
- a workflow engine, service, daemon, database, dashboard, per-task progress
  journal, command receipt store, or completed-work archive;
- duplicate work IDs or migration of completed issues into repository files;
- replacement/consolidation of domain, AI, roadmap, security, or ADR sources;
- automatic commit, push, PR mutation, merge, release, deployment, rollback,
  external/cloud mutation, or paid fallback; and
- creation or modification of `.agent/**`.

### Source-of-truth and conflict rules

1. Current source and directly observed command/runtime results determine
   implementation facts. A plan or document cannot manufacture evidence.
2. `docs/domain.md` and `docs/domain-ai.md` remain the canonical domain and
   staff-side AI contracts.
3. Accepted ADRs govern only their stated scope and are not silently rewritten.
   Proposed ADRs are not binding.
4. Root and applicable narrower `AGENTS.md` files jointly constrain work;
   narrower rules add to and cannot weaken repository-wide safety boundaries.
5. `docs/roadmap.md` remains durable product/phase direction. GitHub Issues
   remain the canonical operational backlog and retain existing planning IDs.
6. While active, an approved `WORK.md` is the issue execution contract but
   cannot override the canonical issue, domain sources, accepted ADRs, or
   engineering contracts. `REVIEWS.md` records review verdicts, not scope.
7. Explicit human decisions control genuinely open choices within existing
   safety/authority boundaries. Durable product or architecture changes must
   update their proper source through approved work.
8. A material conflict or ambiguity involving security, tenancy, authority,
   destructive action, evidence, release, production, or scope fails closed
   and returns the work to planning.

### Exact proposed file layout and documentation impact

New durable files:

```text
docs/engineering/ai-workflow.md
scripts/test_validate_ai_workflow.py
scripts/validate_ai_workflow.py
work/README.md
work/templates/PHASE.md
work/templates/REVIEWS.md
work/templates/WORK.md
```

Dynamic temporary files:

```text
work/phases/active/<phase-id>/PHASE.md
work/phases/active/<phase-id>/REVIEWS.md
work/active/<existing-planning-id>/WORK.md
work/active/<existing-planning-id>/REVIEWS.md
```

Existing files updated during implementation:

```text
.github/workflows/ci.yml
AGENTS.md
CONTRIBUTING.md
Makefile
docs/roadmap.md
```

Responsibilities:

- `docs/engineering/ai-workflow.md` is the single durable detailed contract.
- `work/README.md` explains temporary artifacts, cleanup, and durable GitHub
  handoff.
- Templates encode the minimal schemas in this plan.
- The validator and tests check only local machine-verifiable invariants.
- `AGENTS.md` maps short prompts to the contract and states authority and
  reviewer-role boundaries.
- `CONTRIBUTING.md` adds the concise human lifecycle while preserving and
  explicitly reaffirming its no-per-task-evidence-receipt rule.
- `docs/roadmap.md` adds the Product/Phase Planning synchronization contract
  and makes only a narrow, dated hosted-CI evidence correction.
- `Makefile` adds `workflow-check` and `workflow-test` entry points.
- `.github/workflows/ci.yml` runs those checks within the existing `API`
  job without changing triggers, permissions, existing application checks, or
  job identity.

Preserved without content changes:

```text
.env.example
.gitignore
README.md
SECURITY.md
apps/api/**
deploy/**
docs/domain.md
docs/domain-ai.md
docs/adr/**
```

No ADR is proposed. The harness is an engineering/process contract and does not
need a second durable architectural record duplicating that contract. If
implementation evidence reveals a genuinely separate architectural decision,
that is a material deviation: return to planning and choose any ADR identifier
only from the tracked ADR tree/index at that time.

### Product/Phase Planning contract

The temporary phase loop is:

```text
domain/AI contracts + accepted ADRs + roadmap + implementation + open issues
  -> temporary PHASE.md
  -> fresh independent phase review
  -> exact human phase approval
  -> roadmap current-truth update + authorized GitHub Issue creation/update
  -> PR validation and human merge of roadmap update
  -> temporary phase-artifact cleanup
```

A phase proposal contains its phase ID, author/session labels, revision/hash,
goal/outcomes, dependency order and rationale, exit criteria, major exclusions,
unresolved decisions, and proposed existing planning IDs with issue scopes and
risk. It uses the plan-hash normalization in this document.

The phase reviewer must differ from the phase author and use a fresh agent
session. The reviewer is read-only, independently checks repository and GitHub
evidence, and appends a minimal verdict. Human phase approval must be a
human-authored GitHub record binding phase ID, revision, hash, and `APPROVED`.
Agent-authored approval is prohibited.

Before the temporary phase directory is deleted, the accepted decision is
synchronized to:

- `docs/roadmap.md`: phase ID, approved revision/hash, concise goal/outcomes,
  dependency/order rationale, exit criteria, major exclusions, unresolved
  decisions, and exact planning-ID-to-GitHub-Issue mapping; and
- GitHub Issues: each executable backlog item with its existing planning ID,
  scope, acceptance criteria, risk, and dependency links.

Issue creation/update, roadmap PR mutation, and push require their explicit
human authorizations; roadmap merge is human-only. The phase artifact is
removed from the final roadmap PR tree after the durable issue mapping exists
and the roadmap summary is included. No `work/completed/` or permanent phase
archive is created.

### Minimal pre-merge state machine

Repository-local `WORK.md` and `REVIEWS.md` represent pre-merge delivery
only:

```text
DRAFT
  -> PLAN_REVIEW
  -> APPROVED
  -> IMPLEMENTING
  -> IMPLEMENTATION_REVIEW
  -> READY_FOR_PR
  -> PR_VALIDATION
  -> READY_TO_MERGE
```

Gates and backward transitions:

| State | Required gate | Backward transition |
|---|---|---|
| `DRAFT` | Open canonical issue, current issue digest, complete normative plan/hash | — |
| `PLAN_REVIEW` | Different plan-reviewer actor/session; fresh read-only review | `DRAFT` for findings or any normative/input change |
| `APPROVED` | Passing plan review plus human-authored issue approval comment matching ID/revision/hash | `DRAFT` if approval/input/plan becomes stale |
| `IMPLEMENTING` | Approval and issue digest rechecked; implementation remains in approved scope | `DRAFT` for material deviation |
| `IMPLEMENTATION_REVIEW` | Explicitly authorized candidate commit; different implementation reviewer in a fresh read-only session | `IMPLEMENTING` for in-plan fixes; `DRAFT` for material deviation |
| `READY_FOR_PR` | Passing implementation review for exact candidate SHA and finding gates satisfied | `IMPLEMENTATION_REVIEW` after relevant candidate change |
| `PR_VALIDATION` | Explicit push/PR authority; current issue digest; PR diff and hosted checks inspected | `IMPLEMENTING` or `DRAFT` for failures/deviation |
| `READY_TO_MERGE` | Exact reviewed candidate or cleanup-only descendant, applicable hosted CI observed, no blocking findings, durable pre-merge summary ready | `PR_VALIDATION` after any non-cleanup change |

The local state machine ends at `READY_TO_MERGE`. Human squash merge and
post-merge release are represented by native GitHub PR/Issue/Actions records,
not by a temporary repository state. Merge is always human-only.

Plan-review failure returns to `DRAFT`; implementation-review fixes that stay
inside approved scope return to `IMPLEMENTING`; material deviation always
returns to `DRAFT`, increments the plan revision, and requires independent
plan review and human approval again.

### Lean WORK.md and REVIEWS.md contracts

`WORK.md` contains only:

- issue identity and normalized input checkpoint;
- goal, scope, non-goals, acceptance criteria, risk, and affected areas;
- normative plan and verification plan;
- documentation impact;
- plan revision/hash;
- human approval comment URL/reference;
- implementation author actor/session labels and candidate Git SHA when known;
- material-deviation statement; and
- a concise completion summary before cleanup.

It does not contain milestone updates, a command-by-command log, an evidence
index, or parallel PR/release history.

`REVIEWS.md` is append-oriented only while active. Existing records are not
rewritten. New plan-review records contain only:

```text
review_id, type=PLAN, actor_label, session_label,
plan_author_actor_label, plan_author_session_label,
fresh_session_attestation, plan_revision, plan_hash, issue_digest,
findings(id,severity,status,summary), verdict
```

New implementation-review records contain only:

```text
review_id, type=IMPLEMENTATION, actor_label, session_label,
implementation_author_actor_label, implementation_author_session_label,
fresh_session_attestation, candidate_git_sha, plan_revision, plan_hash,
issue_digest, findings(id,severity,status,summary), verdict
```

Actor/session labels support relational audit:

- plan reviewer actor/session must differ from plan author actor/session;
- implementation reviewer actor/session must differ from implementation author
  actor/session; and
- each final review must come from a fresh agent session.

These labels and attestations do not cryptographically prove identity,
independence, freshness, or correct review. They are explicit process
requirements checked by participants and human review. The local validator can
detect equal/missing labels and malformed bindings only.

Reviewers read the current issue, repository contracts, complete diff/source,
and rerun or inspect relevant commands and GitHub checks. They do not fix the
plan/implementation and do not create a parallel evidence archive. Concise
verification results or native GitHub URLs may be cited only when needed to
support a finding/verdict.

### Deterministic plan hash and human approval provenance

The plan hash algorithm remains `campus-lms-plan-sha256-v1`:

1. Read `WORK.md` as UTF-8; reject invalid UTF-8.
2. Require exactly one line equal to `<!-- PLAN-NORMATIVE-BEGIN -->` and one
   later line equal to `<!-- PLAN-NORMATIVE-END -->`.
3. Extract only bytes between marker lines.
4. Convert CRLF and lone CR to LF.
5. Remove trailing ASCII spaces/tabs on every line.
6. Remove leading/trailing blank lines and end non-empty content with exactly
   one LF; empty content is invalid.
7. SHA-256 the resulting UTF-8 bytes and render
   `sha256:<64 lowercase hexadecimal>`.

The normative section includes the issue digest, revision, goal, scope,
non-goals, affected areas, plan, acceptance criteria, verification,
documentation impact, risks, deviations, and release choice. Mutable state,
approval URL, author/reviewer records, candidate SHA, verification results, and
completion summary are outside it. Any normative change increments
`plan_revision`, changes the hash, invalidates review/approval, and returns to
`DRAFT`.

Human plan approval is authoritative only as an explicit HUMAN-AUTHORED comment
on the canonical GitHub Issue. Its body must contain:

```text
approval_type=PLAN
work_item_id=<ID>
plan_revision=<integer>
plan_hash=sha256:<64 lowercase hexadecimal>
decision=APPROVED
```

The human authors the comment directly; an agent must not compose/post its own
approval, impersonate a human, or treat a chat statement/transcription as the
durable approval. `WORK.md` records only the exact approval comment URL.

At `APPROVED`, implementation start, and PR validation, fetch that URL and
compare the four values to current `WORK.md`. The validator may accept
freshly fetched GitHub comment JSON through an explicit temporary/stdin input
and deterministically validate URL shape and matching values. That input is not
stored as a receipt. Neither the validator nor CI can prove the commenter is a
human, that the account was uncompromised, or that the comment expresses
genuine intent; human authorship/intent remains a declared process trust
boundary.

### Candidate identity and stale implementation review

Candidate identity is a native, full Git commit SHA. A local commit may be
created only when the task explicitly authorizes it. Without commit authority,
implementation may be prepared but cannot advance to
`IMPLEMENTATION_REVIEW`; the agent must request authority rather than invent a
candidate identity.

Before review, record implementation-author actor/session labels and the exact
candidate SHA. The independent implementation reviewer checks out/reads that
candidate, confirms its parent/base and full diff, reruns applicable checks,
and binds the verdict to that SHA plus current plan revision/hash and issue
digest.

Any non-active-work source, test, configuration, durable documentation,
migration, CI, or tooling change after the reviewed candidate requires a new
candidate commit, relevant verification, and fresh independent implementation
review. Appending a review or updating mutable `WORK.md` fields does not alter
the reviewed commit.

For active-only cleanup, the PR head may be a descendant of the reviewed
candidate only when every path in `candidate_sha..head_sha` is deletion of
`work/active/<ID>/WORK.md` and `REVIEWS.md`. PR validation checks that exact
native Git diff. Any other path/content change makes the review stale. Hosted
CI must run on the final cleanup head SHA; this does not claim CI independently
proves review validity.

### Finding severity and residual-risk acceptance

- `CRITICAL` and `HIGH` findings always block.
- `MEDIUM` findings must be resolved or explicitly accepted as residual risk
  by a human.
- `LOW` findings are normally non-blocking but may be escalated when their
  aggregate effect is material.

Residual-risk acceptance is valid only as a HUMAN-AUTHORED canonical Issue or
PR comment containing:

```text
decision=ACCEPTED_RESIDUAL_RISK
work_item_id=<ID>
finding_id=<stable ID>
plan_revision=<integer>
plan_hash=<exact hash>
candidate_git_sha=<exact SHA or NONE for plan review>
rationale=<bounded rationale>
```

`REVIEWS.md` stores only the comment URL beside the finding disposition.
Agents cannot author/accept residual risk. Relevant plan, issue, or candidate
change invalidates the acceptance. A validator can compare supplied comment
values but cannot prove human intent/authorship.

Work risk is the highest applicable class:

- `NORMAL`: bounded ordinary behavior.
- `HIGH`: persistent data, database implementation, integration, concurrency,
  cache/background work, configuration, deployment implementation,
  observability, or performance-sensitive paths.
- `CRITICAL`: authentication/authorization/tenant isolation,
  secrets/cryptography, destructive migrations, production IAM/networking,
  backup/recovery controls, security controls, and release-control or
  governance-enforcement semantics.

### Material deviations

A deviation is material when it changes normative scope/acceptance/
verification, adds an unplanned affected area or dependency, changes behavior
or durable architecture, weakens a safety/review/CI/release control, raises
risk, changes the canonical issue specification, requires unapproved
external/destructive authority, or invalidates an assumption/recovery plan.

Formatting, typo correction, and implementation detail wholly within approved
scope are non-material only when behavior, acceptance, risk, and verification
remain unchanged. Security-sensitive uncertainty is treated as material.

On material deviation: stop; record one concise deviation statement; set
`DRAFT`; revise normative content; increment revision and hash; obtain new
independent plan review; and obtain a new human-authored approval comment.
Reviewers do not implement the deviation.

### Active-only cleanup and durable pre-merge handoff

Each active issue has exactly one `WORK.md` and, from first review onward, one
append-oriented `REVIEWS.md`. No completed/archive directory exists.

Before cleanup, the authorized PR/Issue native records must contain:

- work-item ID and canonical Issue/PR links;
- approved plan revision/hash and human approval-comment URL;
- issue specification digest;
- reviewed candidate SHA and final implementation-review verdict;
- residual-risk comment URLs, if any;
- final hosted CI run URL/head SHA when observed; and
- `release_required=true|false`.

This is a concise handoff, not a command/evidence ledger. Actual diffs, commits,
reviews, checks, and discussion remain in Git/GitHub.

At `READY_TO_MERGE`, delete the active directory in a cleanup-only commit
after explicit local-commit authority. Push/PR update still requires explicit
human authorization. PR validation confirms the reviewed-candidate-to-head diff
contains only those two deletions and observes hosted CI on the cleanup head.
The human squash-merges. `work/active/<ID>/` is not required or present after
merge.

### Separate GitHub-backed post-merge release procedure

Post-merge processing never recreates or relies on repository-local active or
completed work artifacts. The canonical Issue, PR, merge SHA, GitHub Actions,
and human-authored approval comments are the durable state.

Minimum common post-merge disposition comment on the canonical Issue:

```text
record_type=POST_MERGE_DISPOSITION
work_item_id=<ID>
pull_request_url=<URL>
merged_sha=<full SHA>
release_required=true|false
decision=NO_RELEASE|RELEASE_EVALUATION_REQUIRED
recorded_at=<RFC3339 UTC>
```

Posting/updating/closing GitHub records is a PR/Issue mutation and requires
explicit human authorization unless the human performs it directly.

#### Path: release_required = false

1. Human squash-merges the PR.
2. Re-fetch the PR and record its exact merged SHA.
3. Post the common disposition with `release_required=false` and
   `decision=NO_RELEASE`.
4. Close the canonical Issue only by a human or explicitly authorized agent.
   No release evaluation, production approval, deployment, or repository work
   archive follows.

#### Path: release_required = true

1. Human squash-merges; re-fetch the exact merged SHA and post
   `decision=RELEASE_EVALUATION_REQUIRED`.
2. In a fresh post-merge session, an authorized release evaluator inspects the
   canonical Issue/PR, merge SHA, hosted CI, deployable artifacts, current
   deployment contracts, and required external state.
3. Post a concise `RELEASE_EVALUATION` Issue comment containing work-item ID,
   merged SHA, decision `DEFERRED|NOT_RELEASABLE|PROPOSED`, rationale, target
   environment when proposed, and release-scope hash/block when proposed.
   Deferred/not-releasable work creates/follows an Issue as needed and performs
   no production mutation.
4. A proposed release scope is the exact text between unique
   `<!-- RELEASE-SCOPE-BEGIN -->` and `<!-- RELEASE-SCOPE-END -->` marker
   lines in that comment. It minimally specifies work-item ID, merged SHA,
   immutable artifact digests, target environment, ordered external mutations,
   configuration references (never secret values), migration/forward-recovery
   scope, approved rollback/remediation actions, post-deploy checks, stop
   conditions, and authority to create the required claim/result Issue comments.
   Normalize/hash it with the plan-hash line-ending/whitespace rules under
   algorithm `campus-lms-release-scope-sha256-v1`.
5. A human must directly author a separate production-approval comment on the
   canonical Issue containing:

   ```text
   approval_type=PRODUCTION
   work_item_id=<ID>
   merged_sha=<full SHA>
   target_environment=<exact environment>
   release_scope_hash=sha256:<64 lowercase hexadecimal>
   expires_at=<RFC3339 UTC timestamp>
   decision=APPROVED
   ```

   A fetched approval is eligible only as this exact unedited version tuple:
   native comment ID, native URL, native `created_at`, native `updated_at`, and
   a body digest. The approval body digest uses
   `campus-lms-comment-body-sha256-v1`: normalize the complete UTF-8 body with
   the plan-hash LF/trailing-whitespace/blank-boundary rules, SHA-256 the
   normalized bytes, and render `sha256:<64 lowercase hexadecimal>`. Eligibility
   requires native `created_at == updated_at`; GitHub comments are not claimed
   to be immutable. A correction, scope change, or expiry extension must be a
   new human-authored approval comment. The old comment is never edited to
   change authority.

   The authoritative approval time is native GitHub `created_at`, not a
   human-entered body field. The deployment agent must not author approval. The
   accepted human approval channel is a comment authored directly on the
   canonical Issue by a human account that current repository/environment
   governance recognizes as authorized for the target environment; agent, bot,
   service, impersonated, or transcribed comments are ineligible. Recognizing
   that account as human and authorized remains a process trust decision.
6. Before creating any deployment claim or performing any production mutation,
   re-fetch the canonical Issue, release-scope comment, and exact approval
   comment from live GitHub state. Record and compare the exact approval
   comment ID/URL, `created_at`, `updated_at`, and normalized body digest. The
   approval is valid only when:

   - the exact comment still exists on the canonical Issue, is unedited
     (`created_at == updated_at`), and is from the accepted human channel;
   - its work-item ID, merged SHA, target environment, release-scope hash, and
     `decision=APPROVED` match exactly;
   - `created_at` and `expires_at` are unambiguous valid RFC3339 instants;
   - `created_at <= current time < expires_at` and `created_at < expires_at`;
   - no prior or ambiguous deployment claim exists for this approval; and
   - the release scope, external state, and explicit authority still cover the
     exact intended operations.

   Current time is the deployment executor's UTC clock; if a trustworthy
   current time cannot be established, stop. Missing, deleted, edited,
   malformed, future-inconsistent, expired, mismatched, or ambiguous input
   fails closed before production mutation.
7. Treat each approval as intended for one deployment execution. Generate one
   unique opaque execution ID (UUIDv4 is the default), then perform a pre-claim
   live read of all canonical-Issue comments. If a prior or ambiguous claim or
   result references this approval, stop without posting a new claim or
   mutating production and require human resolution/new approval.
8. If no prior claim exists, the approval is still unexpired, and live GitHub
   remains reliable, create one durable `DEPLOYMENT_CLAIM` comment on the
   canonical Issue before the first production mutation:

   ```text
   record_type=DEPLOYMENT_CLAIM
   approval_comment_id=<exact native ID>
   approval_comment_url=<exact native URL>
   approval_created_at=<native created_at>
   approval_updated_at=<native updated_at>
   approval_body_digest=sha256:<64 lowercase hexadecimal>
   work_item_id=<ID>
   merged_sha=<full SHA>
   target_environment=<exact environment>
   release_scope_hash=sha256:<64 lowercase hexadecimal>
   execution_id=<unique execution ID>
   decision=CLAIMED
   ```

   The claim time is the claim comment's native GitHub `created_at`; no
   body-entered claim timestamp is authoritative. Immediately after GitHub
   acknowledges the comment, re-fetch all canonical-Issue comments with
   complete pagination. Proceed only if exactly one valid claim references the
   exact approval version and its execution ID equals the current executor's
   ID. A prior claim, multiple claims, malformed/ambiguous claim state, a
   foreign execution ID, missing just-created claim, incomplete pagination, or
   unreliable GitHub re-read requires stopping before production mutation and
   human resolution or a new approval as appropriate.
9. Immediately before the first production mutation, re-fetch the exact
   approval comment again and verify its ID/URL, `created_at == updated_at`,
   body digest, all bindings, `created_at <= current time < expires_at`, and
   the sole-owned claim. The claim itself must have been created while the
   approval was unexpired. If approval expiry is reached after claim creation
   but before this first mutation, perform no production mutation; create one
   `DEPLOYMENT_RESULT` with `result=ABORTED`,
   `reason=APPROVAL_EXPIRED_BEFORE_MUTATION`, and the exact approval/claim/
   execution bindings; do not reuse that approval; and require a new human
   approval and execution ID for another attempt.
10. The claim/re-fetch procedure reduces accidental duplicate execution but is
    not an atomic lock, compare-and-swap, queue, or cryptographic single-use
    guarantee. Two truly concurrent executors could both observe an apparently
    valid sole claim before GitHub state converges. Operational serialization,
    like reviewer freshness and human intent, remains an explicit process/trust
    boundary. ENG-016 adds no database, lock service, deployment coordinator,
    daemon, queue, or workflow engine to close that boundary.
11. Every claimed execution must create exactly one durable
    `DEPLOYMENT_RESULT` comment tied to its exact approval version and claim:

    ```text
    record_type=DEPLOYMENT_RESULT
    approval_comment_id=<exact native ID>
    approval_comment_url=<exact native URL>
    deployment_claim_url=<exact URL>
    work_item_id=<ID>
    execution_id=<same unique execution ID>
    merged_sha=<full SHA>
    target_environment=<exact environment>
    release_scope_hash=sha256:<64 lowercase hexadecimal>
    production_mutation_started=true|false
    result=SUCCEEDED|FAILED|ABORTED
    reason=<concise sanitized reason when not succeeded>
    verification_reference=<concise sanitized result or native reference>
    ```

    GitHub's native comment `created_at` is the result-record time. Duplicate,
    missing, or contradictory result records are ambiguous and block any later
    execution until human resolution. If GitHub is temporarily unavailable,
    production work stops as soon as safely possible and the result must be
    recorded when the approved path and GitHub are available; the approval may
    not meanwhile authorize another attempt.
12. The approval and claim are consumed once the claim authorizes the first
    production mutation, regardless of eventual result. A consumed approval can
    never authorize another independent deployment execution. Approval expiry
    is the admission boundary for forward production mutation. After mutation
    begins, expiry does not authorize additional forward work: every forward
    production mutation still requires the approval to be within its expiry
    window. If `expires_at` is reached after a partial mutation, stop all
    additional forward mutations and do not continue rollout merely because
    execution started earlier. Only rollback/remediation explicitly enumerated
    in the already-approved release scope may continue under the same execution
    ID when required to reach a safe state; those actions do not expand scope.
    Any new remediation, infrastructure mutation, different rollout strategy,
    continued forward deployment, or out-of-scope operation requires new
    explicit human approval.
    Record the final result as `FAILED` or `ABORTED` as appropriate.
13. If the exact approval comment becomes edited (`created_at != updated_at`),
    deleted, unavailable, or version-mismatched after the first mutation,
    stop forward production mutations, permit only already-approved
    rollback/remediation needed for a safe state, and record the condition in
    `DEPLOYMENT_RESULT`. A later edit cannot widen or extend the execution; any
    new forward or out-of-scope operation requires new human authorization.
14. A failed, partial, expired, edited, deleted, or pre-mutation-aborted
    execution cannot silently recycle its approval. A new attempt requires a
    new human-authored approval comment with a new exact ID/URL/version,
    current merged SHA, current target environment, a release-scope hash
    reassessed against then-current production state, and a new `expires_at`.
    Do not assume the previous scope hash remains valid after partial or failed
    deployment. Before any production mutation, an `ABORTED` claim may be
    superseded only through explicit human resolution/new authorization; the
    agent never deletes or edits claim/result history.
15. Close the Issue only after every claim has one result, the durable outcome
    and any necessary follow-up Issues exist, and ambiguity has been resolved.
    No repository-local post-merge state is created.

### Production procedure specification cases

The durable workflow contract must specify and its review must exercise this
decision table. These are release-agent/live-GitHub procedure cases, not local
repository-validator claims and not an ENG-016 production deployment:

| Case | Required outcome before or during production execution |
|---|---|
| Valid unedited approval | Continue to pre-claim check only when exact ID/URL, `created_at == updated_at`, body digest, bindings, native timestamp, expiry, scope, and authority validate. |
| Edited approval (`created_at != updated_at`) | Stop before claim/mutation; a correction or expiry extension must use a new human-authored comment. |
| Corrected approval using a new comment | Treat only the new exact comment version as eligible after all checks pass. |
| Expiry extension by editing old comment | Invalid; stop and require a new approval comment. |
| Valid new approval after correction | May proceed through pre-claim checks with its new ID/version and bindings. |
| Future-dated approval (`created_at > current time`) | Stop before claim/mutation; require a valid approval whose native creation time has arrived. |
| Approval expired before claim | Stop before claim/mutation; require new approval. |
| Approval expires after claim before first mutation | No mutation; record `ABORTED` with `APPROVAL_EXPIRED_BEFORE_MUTATION`; require new approval/execution ID. |
| Approval valid at first mutation | Consume approval/claim at first mutation, then enforce forward-expiry rules. |
| Approval expires after partial production mutation | Stop forward mutation; permit only pre-approved in-scope recovery and record `FAILED`/`ABORTED`. |
| Rollback/remediation after expiry inside approved scope | May continue under same execution ID only when needed for safe state and explicitly pre-approved. |
| Attempted forward mutation after expiry | Stop; new approval is required. |
| Approval edited/deleted after first mutation | Stop forward mutation; only pre-approved safe recovery may continue; record result. |
| Retry after `FAILED`/`ABORTED` execution | Do not reuse approval; require new approval and execution ID. |
| New approval after partial deployment | Reassess current production state and scope; use a new exact approval version/hash/expiry. |
| Mismatched merged SHA | Stop before claim/mutation. |
| Mismatched target environment | Stop before claim/mutation. |
| Mismatched release-scope hash | Stop before claim/mutation. |
| Prior, multiple, or ambiguous claims | Stop before mutation; require human resolution/new approval. |
| Own unique claim after re-fetch | Continue only when it is the sole valid claim for the exact approval version and all checks still pass. |
| GitHub re-fetch unavailable/incomplete | Stop before mutation; do not infer version or claim validity. |
| Operation outside approved scope | Stop; obtain a new scope and explicit human approval. |

### Deterministic validator, CI responsibility, and trust boundaries

The standard-library validator is intentionally a local artifact checker, not a
workflow engine.

It can deterministically verify:

- required durable templates/contract versions and active-file cardinality;
- directory ID and work-item ID equality;
- allowed state/risk/release values and field/hash/SHA/URL syntax;
- unique normative markers and exact plan/phase hash recomputation;
- issue-digest recomputation when explicitly supplied freshly fetched issue
  JSON;
- pre-merge plan-approval/residual-comment value matching when explicitly
  supplied freshly fetched comment JSON;
- presence of state-required local references/review fields;
- actor/session label inequality;
- finding severity/verdict gates represented in active files;
- rejection of `DONE`, post-merge states, `work/completed/`, and archive
  shapes in temporary artifacts; and
- when given base/candidate/head Git SHAs locally, whether the candidate-to-head
  diff is exactly the permitted active-file cleanup.

Tests cover valid/invalid templates, LF/CRLF/trailing whitespace, mutable versus
normative hash changes, issue title/body digest behavior and comment exclusion,
invalid IDs/enums/hashes, missing/stale bindings, equal actor/session labels,
severity gates, forbidden post-merge/DONE/archive state, approval/comment
matching inputs, and cleanup-only versus other Git diffs.

It cannot deterministically or cryptographically verify:

- that an actor label maps to a distinct real actor;
- that a session is fresh or a review is independent/thorough;
- human authorship, account integrity, intent, or informed approval;
- truth of a transcribed claim or external/production observation;
- semantic materiality, correct risk classification, or adequacy of tests;
- live production-approval validity, deployment-claim uniqueness, deployment
  serialization, claim consumption, or production-result truth;
- current branch protection/repository rules without separate external
  inspection; or
- that a changed validator/workflow/CI file has not weakened its own checks.

Those are declared human/process or remote-platform trust boundaries. Fresh
reviewers and humans must inspect source, diffs, native GitHub authorship/state,
and actual command output. Stronger CODEOWNERS, required-review, ruleset,
separate trusted workflow, or branch-protection enforcement may be proposed
only as a separate approved governance work item.

Production approval, claim, re-fetch, and result checks are deliberately not
implemented in the local repository validator. They are runtime release-agent
procedures against live GitHub state, and GitHub comments do not expose a sound
atomic serialization primitive for this contract.

CI will run validator tests and the validator against whatever active artifacts
are present, plus all existing repository checks, in the existing `API` job.
No base-aware PR-to-work-item enforcement is proposed because active artifacts
must be absent from the final tree and the current workflow can modify its own
validator invocation. The final cleanup tree still tests durable templates and
validator behavior, but `CI / API` alone is not tamper-proof proof of
governance approval. ENG-016 does not change branch protection and does not
claim this procedural gap is remotely enforced.

### Semantically protected control plane

A dedicated CRITICAL governance work item is required when a change alters the
meaning or enforcement of:

- source precedence, plan/issue/candidate binding, review independence,
  approval provenance, finding gates, material-deviation behavior, authority,
  active-only cleanup, or post-merge release/production approval;
- authentication, authorization, tenant isolation, secrets/cryptography,
  destructive migration, production IAM/networking, backup/recovery, or other
  security controls; or
- CI required-check semantics, security validation, release/deploy gates, or
  the validator's enforcement logic.

Conservative path defaults apply to the durable workflow contract,
`work/templates/**`, the workflow validator/tests, and
`.github/workflows/**`, because their primary purpose is governance/CI
enforcement. Changes to them require explicit governance classification unless
an independent review demonstrates the path change is purely non-semantic.

`Makefile`, `CONTRIBUTING.md`, general documentation, and general
`AGENTS.md` changes are not automatically CRITICAL by filename. They become
protected-control-plane changes only when the actual diff triggers the
semantics above. Local CI cannot make this classification tamper-proof; plan
author, independent reviewer, and human approver enforce it procedurally.

ENG-016 is the dedicated CRITICAL bootstrap item because its actual semantics,
not merely its filenames, establish these controls.

### Remote authority boundaries

- Local commit: only when explicitly authorized by the task.
- Push: requires explicit human authorization.
- PR creation, edits, comments, reviews, labels, or closure: require explicit
  human authorization unless directly performed by the human.
- GitHub Issue creation/update/comment/closure: requires explicit human
  authorization unless directly performed by the human; approval and
  residual-risk/production-approval comments must be human-authored.
- Merge: human-only.
- Production/external/cloud mutation: exact explicit human authorization,
  bounded scope, and all applicable repository/deployment controls.
- Destructive database/data work: separate explicit authorization.

Passing local tests, review, CI, merge, or a release evaluation grants none of
the later authorities.

### Ordered implementation plan

1. Re-fetch issue #16, recompute its digest, validate revision-4 plan hash, and
   verify a human-authored approval comment matches this revision/hash.
2. Add the single durable workflow contract and lean active templates/README;
   add no ADR.
3. Add focused standard-library validator tests and implement only the local
   deterministic checks enumerated above.
4. Add the two Make targets.
5. Update `AGENTS.md` and `CONTRIBUTING.md` concisely, preserving the
   no-receipt rule and authority boundaries.
6. Update `docs/roadmap.md` only for the Product/Phase synchronization
   contract and dated CI observation.
7. Add validator test/check steps to the existing CI job without changing its
   identity, triggers, permissions, or existing checks.
8. Inspect all changed files and full diff; run focused and existing checks;
   record only concise results in the completion summary.
9. With explicit local-commit authority, create the candidate commit and record
   implementation-author labels/SHA.
10. Obtain a fresh independent implementation review. With later explicit
    push/PR authority, validate the PR and hosted CI, create the cleanup-only
    commit, and stop at `READY_TO_MERGE` for human merge.

### Acceptance criteria

1. Existing domain, AI, security, roadmap, and one-file-per-ADR sources remain
   canonical; GitHub Issues remain the operational backlog and existing
   planning IDs remain identities.
2. Product/Phase Planning durably synchronizes approved phase goal/order/exit
   criteria/exclusions/decisions/hash and exact Issue mapping to the roadmap and
   GitHub before temporary cleanup.
3. The pre-merge state machine has only the eight specified states and handles
   review failure/material deviation without repository post-merge states.
4. `WORK.md` and `REVIEWS.md` carry only minimum plan, approval, candidate,
   finding, verdict, deviation, verification, and completion bindings; no
   progress/evidence receipt system is introduced.
5. Plans bind revision/hash to normative fields and normalized canonical Issue
   title/body; issue digest is rechecked before implementation and PR
   validation.
6. Plan approval, residual-risk acceptance, and production approval use exact
   human-authored durable GitHub comments; agents cannot author those decisions.
7. Reviewer independence is relationally defined with minimum actor/session
   labels and fresh-session requirements, with the non-cryptographic trust
   limitation explicit.
8. Implementation review binds to an authorized native candidate commit SHA and
   becomes stale after any non-active-file candidate change.
9. CRITICAL/HIGH block; MEDIUM is resolved or human-accepted by exact durable
   reference; LOW is normally non-blocking.
10. Active files are removed before human merge. Both
    `release_required=false` and `release_required=true` have exact
    GitHub-backed post-merge procedures and no `work/completed/`.
11. Production execution requires an exact, unedited, unexpired human-authored
    approval. The approval is revalidated before the first production mutation
    and forward mutations may not continue past its expiry. After production
    mutation has started, only explicitly pre-approved rollback/remediation may
    continue after expiry when necessary to reach a safe state. Any further
    forward or out-of-scope operation requires new human approval. GitHub
    comments do not provide an atomic distributed lock; prevention of truly
    concurrent executors remains an explicit process boundary.
12. The small validator checks only enumerated machine-verifiable invariants;
    CI runs it in the existing job while explicitly making no claim of reviewer
    freshness, human intent, self-protection, or tamper-proof governance.
13. Protected-control-plane treatment is semantic, with narrow conservative
    path defaults; ordinary Makefile/CONTRIBUTING/AGENTS edits are not
    automatically CRITICAL.
14. No LMS runtime, auth, authorization, tenancy, data/migration, deployment,
    branch-protection, external, or production behavior changes.

### Verification plan

During implementation, the relevant agent reruns and directly inspects:

```text
python3 -m unittest scripts/test_validate_ai_workflow.py
python3 scripts/validate_ai_workflow.py
make workflow-test
make workflow-check
make help
make test
make build
git diff --check
git status --short
```

It also:

- re-fetches issue #16 and approval comment inputs rather than relying on
  stored transcripts;
- confirms the full changed-path set and no application/deploy/ADR/secret/data
  scope;
- checks current CI YAML retained its job name, triggers, permissions, and all
  existing steps;
- reviews both release-required procedure shapes and every “Production
  procedure specification cases” row as contract scenarios, without adding
  production-comment logic to the local validator and without GitHub or
  production mutation;
- tests candidate staleness and cleanup-only Git diffs in temporary Git
  repositories; and
- records only a concise pass/fail/limitation summary, while reviewers rerun
  commands and inspect native diffs/checks themselves.

HOSTED_CI is claimed only after an explicitly authorized push/PR and direct
observation of the run URL, head SHA, job, and conclusion. No EXTERNAL_STATE or
PRODUCTION operation/evidence is required for ENG-016.

### Risks and recovery

- Procedural controls can be bypassed by a malicious/self-modifying governance
  change under current remote settings. This is disclosed, not presented as
  deterministically solved; stronger enforcement is separate work.
- Human/reviewer identity and intent cannot be proven by local labels/parsers.
  Native human-authored comments plus relational checks improve auditability
  without claiming cryptographic proof.
- GitHub deployment claims are durable and re-fetched but are not an atomic
  distributed lock. Operationally permitting only one executor per approval is
  required; a stronger coordinator is explicitly outside this bootstrap.
- Temporary cleanup could lose decisions. The required roadmap/Issue/PR
  handoffs precede deletion and are independently checked.
- CI could be disrupted by the validator. It uses Python standard library only,
  has focused tests, and leaves all existing checks in place.
- Before merge, recovery is limited to reverting ENG-016 harness paths on its
  short-lived branch; no application/data recovery is involved.
- After merge, changing the harness requires a new semantically classified
  governance issue and normal reviewed PR. Published history is not rewritten.

<!-- PLAN-NORMATIVE-END -->

## Mutable approval, candidate, deviation, and completion fields

- Human plan-approval comment URL:
  `https://github.com/farisakbar28/campus-lms/issues/16#issuecomment-5556469166`
- Implementation author actor label: `Codex implementer`
- Implementation author session label: `ENG-016-implement-r4-2026-09-06`
- Candidate Git commit SHA: `c03eeeab7a9ad7508246214f64c42b7027c627c6`
- Material deviations: `NONE`
- Concise verification result: `PASS: workflow tests, canonical/bootstrap/legacy/author probes, Go race tests, build, go vet, go mod verify, formatting, and diff checks; workflow-check expected-fails on OPEN ENG-016-IMPL-008 and ENG-016-IMPL-009 pending fresh review.`
- Concise completion summary: `CANDIDATE_CREATED_c03eeeab7a9ad7508246214f64c42b7027c627c6_STOPPED_FOR_INDEPENDENT_IMPLEMENTATION_REVIEW`

## Finding-resolution matrix for Independent PLAN Reviews 1–3

| Finding | Current disposition | Exact revised sections |
|---|---|---|
| `ENG-016-PLAN-001` | `RESOLVED` by revision 2 and unchanged. | “Minimal pre-merge state machine”; “Active-only cleanup and durable pre-merge handoff”; “Separate GitHub-backed post-merge release procedure” |
| `ENG-016-PLAN-002` | `RESOLVED` by revision 2 and unchanged. | “Deterministic validator, CI responsibility, and trust boundaries”; “Semantically protected control plane”; acceptance criteria 12–13 |
| `ENG-016-PLAN-003` | `RESOLVED` by revision 2 and unchanged. | “Lean WORK.md and REVIEWS.md contracts”; “Verification plan”; acceptance criterion 4 |
| `ENG-016-PLAN-004` | `RESOLVED` by revision 2 and unchanged. | “Lean WORK.md and REVIEWS.md contracts”; “Deterministic plan hash and human approval provenance” |
| `ENG-016-PLAN-005` | `RESOLVED` by revision 2 and unchanged. | “Identity and input checkpoint”; “Deterministic plan hash and human approval provenance” |
| `ENG-016-PLAN-006` | Addressed in revision 4, pending fresh independent review. | “Separate GitHub-backed post-merge release procedure”, release-required=true steps 4–15; “Production procedure specification cases”; “Deterministic validator, CI responsibility, and trust boundaries”; acceptance criterion 11 |
| `ENG-016-PLAN-007` | `RESOLVED` by revision 2 and unchanged. | “Product/Phase Planning contract” |
| `ENG-016-PLAN-008` | `RESOLVED` by revision 2 and unchanged. | “Exact proposed file layout and documentation impact”; “Ordered implementation plan” |

Reviews 2 and 3 are authoritative that findings 001–005 and 007–008 are
resolved. Finding 006 remains authoritative as open until a fresh independent
reviewer evaluates revision 4 and appends a verdict to `REVIEWS.md`.

### Concise resolution mapping for ENG-016-PLAN-006

- Approval version/time: requires exact native comment ID/URL, unedited
  `created_at == updated_at`, normalized body digest, and native `created_at` as
  the approval time satisfying `created_at <= current time < expires_at`.
- Exact validity: live canonical-Issue existence, accepted human channel, work
  ID, merged SHA, environment, release-scope hash, decision, timestamps, scope,
  and authority are all rechecked; malformed/missing/mismatched/ambiguous state
  stops before mutation.
- Intended one execution: added a pre-claim read, uniquely identified
  `DEPLOYMENT_CLAIM`, immediate complete re-fetch, and sole-own-claim gate.
- Accurate boundary: explicitly states GitHub comments are not an atomic lock
  or deterministic concurrency guarantee; serialization remains procedural and
  no coordinator/workflow engine is added.
- Outcome/retry: every claim requires one bound `DEPLOYMENT_RESULT`; first
  production mutation consumes the approval; in-scope rollback may continue,
  while retry/out-of-scope work requires new human approval.
- Validator/coverage: production comment checks remain live release-agent
  procedures outside the local validator, with every required scenario
  specified in the production decision table.

## Unresolved limitations

- Current remote branch/review settings and the existing CI job do not provide
  tamper-proof governance enforcement. ENG-016 documents this limitation and
  does not mutate branch protection.
- A repository-local validator cannot prove reviewer freshness, distinct real
  identities, human authorship/intent, or that it has not been maliciously
  weakened. Those remain explicit human/process and remote-platform trust
  boundaries.
- No finding is impossible to resolve as directed. The two limitations above
  are constraints the revised instructions explicitly require documenting, not
  unresolved blockers to independent review.

## Recommended risk classification

`CRITICAL`, based on the workflow’s governance, CI, approval, and release
control semantics. No application or production change is included.
