# ENG-016 reviews

## Review 1

- Review type: `PLAN`
- Recorded at: `2026-09-05T18:44:32Z`
- Actor: `Codex plan reviewer`
- Actor role: `independent plan reviewer`
- Fresh-session attestation: I derived this review from the repository and
  current read-only GitHub evidence. I did not use previous agent conversation,
  conclusions, or delegated review work.
- Subject work item: `ENG-016`
- Subject plan revision: `1`
- Exact reviewed plan checkpoint SHA:
  `sha256:e623cf47c59767b9d968911f240ddd0fa0afff9e19eeb2dbc25af062bdff231d`
- Independently recomputed normative SHA:
  `sha256:e623cf47c59767b9d968911f240ddd0fa0afff9e19eeb2dbc25af062bdff231d`
- Full `WORK.md` file SHA-256:
  `87eedf4eb50db198e56aa487254cbc64426051d5f28ca71e02837e0ee91373bf`
- Reviewed Git base/HEAD: `92419ed6144e8d07b1392e0b931cdf32a044af77`
- Verdict: `CHANGES_REQUIRED`

### Executive assessment

The plan correctly classifies ENG-016 as CRITICAL, preserves the LMS and AI
domain sources, keeps existing GitHub planning IDs, verifies its normative
hash, distinguishes local/hosted/external/production observations, and states
the requested human authority boundaries for issue/PR writes, push, merge,
external mutation, and production. Its proposed application and deployment
scope is compatible with the current repository and does not alter LMS,
authentication, tenancy, migration, or deployment behavior.

It is not approvable at revision 1. The cleanup mechanism makes the post-merge
release states unrepresentable, and the proposed CI check can be bypassed by
omitting or deleting the active artifact that it validates. The plan also
reintroduces a mandatory per-task progress/evidence ledger despite the current
repository's explicit rejection of a per-task evidence-receipt system. Those
are blocking design defects, not implementation details.

### Findings

#### ENG-016-PLAN-001 — HIGH — Pre-merge cleanup destroys the post-merge release state machine

Status: `OPEN` / blocking.

Evidence:

- `WORK.md:259-265` requires transitions from `READY_TO_MERGE` through
  `MERGED`, conditional `RELEASE_EVALUATION`, production approval/deploy,
  post-deploy verification, and finally `DONE`.
- `WORK.md:390-415` places release evaluation, production approval, and
  post-deploy verification records in the active `REVIEWS.md`.
- `WORK.md:555-573` instead requires deletion of the entire active directory
  after `READY_TO_MERGE`, before the human squash merge.

Impact:

For every work item with `EVALUATE_AFTER_MERGE`, the only defined state and
review-record files no longer exist when the workflow reaches the states that
need them. GitHub is named as the durable record, but no GitHub-side state
schema, transition owner, approval record, or validator binding is defined.
The issue's conditional-release and production-approval properties therefore
cannot be executed as specified.

Required resolution:

Choose one coherent lifecycle. Define an explicit durable GitHub-backed schema
and transition procedure for all post-merge states and bind production approval
to it, or redesign cleanup/state storage so release-applicable work remains
representable without merging a completed active archive. Add focused
acceptance tests for both `NONE` and `EVALUATE_AFTER_MERGE` paths.

#### ENG-016-PLAN-002 — HIGH — The required CI/control-plane gate is structurally bypassable

Status: `OPEN` / blocking.

Evidence:

- `WORK.md:582-591` validates the shape and gates of active artifacts that are
  present, but does not require an active item for a relevant PR diff.
- `WORK.md:566-573` deliberately deletes the active item before the final
  merge candidate.
- `WORK.md:604-608` adds only the test/check invocations to the existing job;
  `WORK.md:636-638` says protected-path reporting occurs only “when comparison
  context is available,” without specifying or requiring that context.
- `.github/workflows/ci.yml:23-26` performs the current checkout, and the plan
  does not define a base-ref/diff input for `workflow-check`.
- Current GitHub branch protection, inspected read-only at review time,
  requires the `API` context and enforces protection for administrators, but
  requires zero approving reviews, does not dismiss stale reviews, and does not
  require approval of the last push.

Impact:

A PR can change a protected control-plane path while omitting an active work
item, or can delete the item during cleanup, and the repository-only validator
has no required subject on which to enforce plan/review/approval bindings. A PR
that changes the validator or workflow can also validate its own weakened
version. The final required `API` result would therefore not prove the claimed
control-plane rule, independent review, or human-approval gate.

Required resolution:

Specify and test a base-aware CI invariant that maps every relevant diff,
especially protected-control-plane changes, to an exact approved work item and
current review even on the cleanup commit. Define a trust boundary for changes
to the validator/workflow themselves. If remote review/branch-protection
changes remain outside ENG-016 authority, record the unenforced controls as an
explicit limitation or separately approved prerequisite instead of describing
them as protected or deterministically gated.

#### ENG-016-PLAN-003 — HIGH — The plan recreates the evidence/progress bureaucracy the repository rejects

Status: `OPEN` / blocking.

Evidence:

- `CONTRIBUTING.md:64-68` explicitly says there is no per-task evidence-receipt
  system and directs reviewers to actual command output and the diff.
- `WORK.md:329-340` makes implementation notes, a deviation log, an evidence
  index, findings summary, approval references, PR/merge references, and
  cleanup status mandatory mutable fields for every work item.
- `WORK.md:390-415` mandates a 15-field append-oriented record spanning eleven
  record types.
- `WORK.md:533-553` mandates a structured receipt with identity, class, time,
  observer, subject, command/source, result, location, and limitations for
  every evidence entry.
- `WORK.md:246-265` adds fourteen workflow states; `WORK.md:630-634` also makes
  every change to broad files such as the whole `Makefile` and
  `CONTRIBUTING.md` a dedicated CRITICAL governance item.
- `WORK.md:168-170` says the `CONTRIBUTING.md` update will not replace existing
  guidance, but the planned evidence model directly reverses its current
  no-receipt rule.

Impact:

Keeping receipts inside two files instead of separate per-step files does not
make them cease to be per-task progress/evidence receipts. The mandatory
ledger, custom candidate format, state transitions, and overbroad CRITICAL path
classification are disproportionate for this one-commit, small repository and
conflict with Issue #16's requirement for a lean harness rather than a workflow
engine. They also create duplicate status sources across `WORK.md`,
`REVIEWS.md`, GitHub Issues, PRs, GitHub reviews, and CI.

Required resolution:

Reconcile the proposal explicitly with `CONTRIBUTING.md` and reduce the active
schema to the minimum bindings the issue actually requires: plan revision/hash,
approval provenance, candidate revision, findings/verdict, residual-risk
decision, and hosted check reference. Prefer direct review of command output,
the repository diff, native Git/GitHub revisions, PR reviews, and CI statuses
over a second progress/evidence ledger. Narrow protected paths to semantic
control-plane changes or justify each whole-file classification.

#### ENG-016-PLAN-004 — MEDIUM — Reviewer independence and human approval provenance are asserted but not bound

Status: `OPEN`.

Evidence:

- `WORK.md:269-273` prohibits an implementer from authoring passing reviews but
  does not explicitly prohibit the plan author from issuing the plan verdict,
  and the schema contains no plan-author or implementation-author identity to
  compare with the reviewer.
- `WORK.md:384-386` permits the human approval to be transcribed into a local
  active file, while `WORK.md:390-423` relies on actor strings and a
  fresh-session attestation that the validator explicitly cannot authenticate.
- Current branch protection requires no approving GitHub review.

Impact:

Revision/hash matching prevents content drift but does not establish that the
review was independent or that the named human actually issued the approval.
An agent-authored assertion can satisfy the proposed file shape.

Required resolution:

Define independence relationally: the plan reviewer must differ from the plan
author, and the implementation reviewer must differ from the implementation
author, with stable recorded session/actor identifiers. Require human approval
to cite an authoritative human-authored source (for example, the exact GitHub
comment/review URL or another explicitly accepted channel) and state which
parts remain human judgment rather than deterministic enforcement.

#### ENG-016-PLAN-005 — MEDIUM — Approval does not become stale when the canonical GitHub Issue changes

Status: `OPEN`.

Evidence:

- `WORK.md:252-255` requires an open canonical issue, and
  `WORK.md:342-345` validates only its number/URL-to-ID mapping.
- `WORK.md:347-386` binds approval to the normative `WORK.md` content, not to an
  observed issue version, update timestamp, or body digest.
- `WORK.md:456-469` comprehensively invalidates implementation review for
  repository changes but defines no equivalent invalidation when the canonical
  issue body changes.

Impact:

Issue scope or acceptance criteria can change after plan approval while the
plan hash remains valid. The conflict prose tells a diligent agent to stop once
it notices the change, but the plan omits a reproducible stale-input check for
the source it calls the canonical operational backlog.

Required resolution:

Record a reproducible issue checkpoint (at minimum issue number, observed
`updatedAt`, and a normalized body hash or immutable GitHub reference) among
the normative inputs, re-check it at implementation and PR gates, and require
re-planning when a material canonical-issue change is observed.

#### ENG-016-PLAN-006 — MEDIUM — Production approval has no defined validity/expiry record

Status: `OPEN`.

Evidence:

- `WORK.md:263` requires an “unexpired” production approval.
- `WORK.md:411-415` permits a `PRODUCTION_APPROVAL` record but defines no
  production-specific required fields.
- `WORK.md:668-673` defines what the release-scope hash covers, but not
  `approved_at`, `expires_at`, approver identity/provenance, one-time use,
  environment binding validation, or invalidation after failed/partial use.

Impact:

The workflow cannot deterministically decide whether a production approval is
current, expired, reusable, or authentic. The general fail-closed conflict rule
reduces immediate risk but leaves a required Issue #16 boundary incomplete.

Required resolution:

Define the production-approval schema and its validation semantics, including
authoritative human provenance, exact environment/scope hash, approval and
expiry times, invalidation events, and reuse/failure behavior. Preserve the
separate explicit authorization required for external/cloud execution.

#### ENG-016-PLAN-007 — MEDIUM — Product/Phase Planning cleanup loses the approved phase decision

Status: `OPEN`.

Evidence:

- `WORK.md:218-229` puts dependency order, proposed scopes, acceptance
  criteria, exclusions, unresolved decisions, and rollout assumptions in the
  temporary `PHASE.md`.
- `WORK.md:237-243` requires only the issue mapping and review summary to be
  durably present in GitHub before deleting the phase directory.

Impact:

Issue bodies can preserve per-issue scope, but the phase-level sequencing,
cross-issue rationale, exclusions, and rollout assumptions can disappear. The
roadmap remains intentionally broad, so neither it nor the issue mapping is a
defined durable replacement for the approved phase decision.

Required resolution:

Define the minimal phase summary that must be transferred to a durable GitHub
location before cleanup, including revision/hash, outcomes, ordering/rationale,
exclusions, unresolved decisions, and the exact issue mapping. Avoid retaining
a repository phase archive.

#### ENG-016-PLAN-008 — LOW — ADR numbering relies on non-repository IDE state

Status: `OPEN`.

Evidence:

- `WORK.md:847-850` chooses ADR-0007 because an IDE tab names an absent
  ADR-0006.
- `docs/adr/README.md:47-49` says a missing number is not a placeholder or a
  requirement to create one.
- Neither the tracked tree nor the worktree contains
  `docs/adr/0006-trusted-authenticated-tenant-context.md` at the reviewed
  checkpoint.

Impact:

The planned durable filename is influenced by non-repository state even though
the harness is intended to derive decisions from repository/GitHub evidence.
The gap is harmless technically but the rationale violates the proposed source
discipline.

Required resolution:

Choose the ADR identifier from repository state at implementation time, or
record a repository-grounded reason for intentionally selecting 0007. Treat a
new conflicting ADR path as a material deviation.

### Evidence actually inspected

Repository and Git evidence:

- `HEAD`, local `master`, local `origin/master`, and current remote
  `refs/heads/master` all resolved to
  `92419ed6144e8d07b1392e0b931cdf32a044af77` at review time.
- Branch: `chore/ai-engineering-harness`.
- Tracked `git diff master --`: empty.
- Before this review file was created, the sole worktree difference was the
  untracked 865-line `work/active/ENG-016/WORK.md`; its complete new-file diff
  against `/dev/null` was inspected.
- The repository history reachable from ordinary branches contains one commit,
  `92419ed6144e8d07b1392e0b931cdf32a044af77` (`chore: establish engineering
  baseline`).
- The normative plan hash was independently recomputed according to the plan's
  LF/trailing-whitespace/blank-boundary algorithm and matched the declared
  checkpoint.
- Repository structure was inspected with `rg --files`, `git ls-tree`, and
  scoped file enumeration. No tracked workflow harness, `work/` artifact,
  `docs/engineering/`, root `scripts/`, CODEOWNERS file, or ADR-0006 exists at
  the base revision.

GitHub evidence, read-only at `2026-09-05T18:44:32Z`:

- Issue #16 is open, is identified as `ENG-016`, has CRITICAL risk, no comments,
  and contains the requirements/non-goals used for this review.
- The complete issue list (#1 through #16) was inspected for planning IDs,
  overlap, and sequencing. `AUTHCTX-001` is open as issue #1, so the proposed
  pilot preserves an existing canonical ID.
- No PR exists for head branch `chore/ai-engineering-harness`.
- GitHub Actions run `33964729866` is a successful `CI` push run for exact SHA
  `92419ed6144e8d07b1392e0b931cdf32a044af77`; its `API` job and all existing
  steps succeeded.
- Current `master` protection requires status context `API`, has
  `enforce_admins=true`, `required_conversation_resolution=true`,
  `allow_force_pushes=false`, and `allow_deletions=false`; it has
  `required_approving_review_count=0`, `dismiss_stale_reviews=false`,
  `require_last_push_approval=false`, and non-strict status checks. These are
  time-bounded external observations, not repository guarantees.

Files read:

- Root and scoped contracts: `AGENTS.md`, `apps/api/AGENTS.md`,
  `deploy/AGENTS.md`, `CONTRIBUTING.md`, `SECURITY.md`, and `README.md`.
- Product/architecture sources: all of `docs/domain.md`, `docs/domain-ai.md`,
  `docs/roadmap.md`, and every file currently under `docs/adr/`, including the
  index and template.
- Automation: `.github/workflows/ci.yml` and the complete `Makefile`.
- Relevant implementation/verification scripts: both scripts under
  `apps/api/scripts/`, `apps/api/testdata/auth-sessions-verify.sh`, and all
  scripts under `deploy/scripts/`, with particular attention to
  `deploy.sh`, its local safety test, production-Compose verification, and the
  current backup/restore safety boundary.
- Repository layout, migrations, API package/test layout, Compose/Caddy files,
  and testdata paths were enumerated. No application or governance file was
  changed by this review.

### Residual concerns

- The reachable repository has only the single squashed baseline commit, so a
  claim about exactly which older learning/progress/evidence artifacts were
  previously removed cannot be independently reconstructed from ordinary Git
  history. The current `CONTRIBUTING.md` no-receipt rule is direct evidence of
  the intended present baseline and is sufficient for finding 003.
- Reviewer freshness and human intent can never be fully proven by a local
  parser. After the schema is tightened, these should remain explicitly named
  human/process assumptions rather than automated guarantees.
- GitHub protection and hosted-check configuration are external state and may
  change. They must be re-read before relying on them; this review did not
  mutate them.
- The proposed roadmap correction is evidence-bounded and the existing domain,
  AI, security, API, migration, and deployment documents are otherwise
  appropriately preserved. Any implementation expansion beyond the listed
  governance/documentation/CI paths would be a material deviation.

### Verdict

`CHANGES_REQUIRED`

Open HIGH findings `ENG-016-PLAN-001`, `ENG-016-PLAN-002`, and
`ENG-016-PLAN-003` block approval. Plan revision 1 must not advance to human
approval or implementation. The planner must revise the normative plan,
increment the revision, recompute the plan hash, and obtain a fresh independent
PLAN review.

## Review 2

- Review ID: `ENG-016-PLAN-REVIEW-002`
- Review type: `PLAN`
- Recorded at: `2026-09-05T19:21:16Z`
- Actor: `Codex plan reviewer r2`
- Actor role: `independent plan reviewer`
- Session label: `ENG-016-plan-review-r2-2026-09-06`
- Plan author actor label: `Codex planner`
- Plan author session label: `ENG-016-plan-r2-2026-09-06`
- Fresh-session attestation: I derived this review independently from the
  canonical Issue, current repository and GitHub state, the complete revision-2
  plan, and the preserved Review 1 record. I did not adopt the planner's
  finding-resolution matrix or previous agent conclusions as evidence.
- Subject work item: `ENG-016`
- Subject plan revision: `2`
- Exact reviewed plan checkpoint SHA:
  `sha256:4e46736d291879d53ea33619d2a639243bb7ef086815e1a7f9f21c40e3305b4d`
- Independently recomputed normative SHA:
  `sha256:4e46736d291879d53ea33619d2a639243bb7ef086815e1a7f9f21c40e3305b4d`
- Exact Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Independently recomputed Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Full reviewed `WORK.md` SHA-256:
  `d6feee08ca53030eaf4152298ba8da886f8569dd71ca2143fecef24a1359ea9b`
- Review 1 pre-append file SHA-256:
  `dfea908ba74fdaa841666918da19eada9013a047ffd700eebeba200841492ec7`
- Reviewed Git base/HEAD: `92419ed6144e8d07b1392e0b931cdf32a044af77`
- Verdict: `CHANGES_REQUIRED`

### Integrity and input verification

- `WORK.md` declares plan revision `2`, and the normative section also binds
  revision `2`.
- Independent byte-level normalization of the unique normative marker section
  produced 35,620 bytes and exactly matched the declared plan hash.
- GitHub Issue #16 was re-fetched read-only. It is open, its title/body and URL
  match ENG-016, its observed `updatedAt` remains
  `2026-09-05T18:16:42Z`, and its normalized 58-byte title plus 3,925-byte body
  exactly reproduce the declared Issue digest.
- Before this append, `REVIEWS.md` contained exactly one review header and the
  eight Review 1 finding IDs, in order. Its complete SHA-256 is recorded above.
  Its filesystem modification/change time preceded the revision-2 `WORK.md`
  modification/change time, and the revision-1 `WORK.md` full-file SHA recorded
  by Review 1 independently matches the recoverable Git object
  `94f5db7ec65d91fb02139e665748c621f22296a8`. No alteration of Review 1 was
  detected. Because Review 1 did not record a self-hash and was not tracked at
  that checkpoint, the repository cannot supply cryptographic proof against an
  alteration made before the observed filesystem checkpoint; the pre-append
  hash above establishes the exact Review 1 bytes reviewed here.

### Independent resolution audit of Review 1 findings

| Finding | Review 2 status | Independent basis |
|---|---|---|
| `ENG-016-PLAN-001` | `RESOLVED` | Local state now ends at `READY_TO_MERGE`; active files are deleted before merge, while exact GitHub Issue/PR schemas and separate procedures represent both post-merge release paths. |
| `ENG-016-PLAN-002` | `RESOLVED` | Revision 2 no longer claims base-aware or tamper-proof enforcement. It limits CI to present-artifact/template checks and explicitly states that absent active artifacts and self-modifying CI remain procedural gaps. |
| `ENG-016-PLAN-003` | `RESOLVED` | The active schema removes milestone logs, command receipts, evidence indexes, and parallel release history; reviews retain only identity/binding/findings/verdict data, with native Git/GitHub/CI as evidence. |
| `ENG-016-PLAN-004` | `RESOLVED` | Reviewer independence is relationally bound to author actor/session labels, human plan approval must be a directly human-authored canonical-Issue comment URL, and non-cryptographic identity/intent limits are explicit. |
| `ENG-016-PLAN-005` | `RESOLVED` | The normalized Issue title/body digest is normative, indirectly bound by the plan hash, and must be re-fetched at implementation start and PR validation; any digest change invalidates approval. |
| `ENG-016-PLAN-006` | `OPEN` | Revision 2 adds the requested fields and failure/retry prose, but its validity and single-use checks remain underspecified as detailed below. |
| `ENG-016-PLAN-007` | `RESOLVED` | Before phase cleanup, the roadmap must retain revision/hash, outcomes, ordering rationale, exit criteria, exclusions, unresolved decisions, and exact Issue mapping; Issues retain executable scopes and dependencies. |
| `ENG-016-PLAN-008` | `RESOLVED` | ENG-016 no longer proposes or reserves an ADR number; any later architectural decision is a material deviation selected from then-current tracked state. |

### Findings

#### ENG-016-PLAN-006 — MEDIUM — Production approval validity and single-use semantics remain underspecified

Status: `OPEN` / blocking until resolved or explicitly accepted by a human
under the exact residual-risk procedure.

Evidence:

- `WORK.md:578-590` defines `approved_at` and `expires_at` values inside a
  human-authored comment, but does not bind either value to the native GitHub
  comment timestamp or require `approved_at <= current_time`; the stated check
  at `WORK.md:592-596` requires only a valid interval and current time before
  `expires_at`.
- `WORK.md:592-601` requires the deployment agent to verify an undefined
  “unused status,” then says consumption begins with the first production
  mutation and is recorded on the Issue. It defines neither a pre-mutation
  claim/consumption record nor an ordering or serialization rule for two fresh
  deployment sessions that read the same still-unused approval concurrently.
- The deterministic validator inventory at `WORK.md:615-631` contains no
  production-approval, release-scope, native-comment-time, or consumption
  validation. `WORK.md:639-648` correctly disclaims proof of production facts,
  but does not identify single-use enforcement as process-only and non-atomic.
- `WORK.md:762-764` nevertheless makes the unqualified acceptance claim that
  production approval is single-use and fails closed after use or failure.
- Current `deploy/scripts/deploy.sh` validates deployment inputs and performs
  mutation directly; ENG-016 preserves `deploy/**`, so no executable gate or
  lock supplies the missing semantics.

Impact:

A future executor could accept an approval before its stated `approved_at`, and
two sessions could both pass the same pre-mutation “unused” observation before
either records consumption. The separate exact external-operation authority
and required human approval reduce immediate exposure, so this is MEDIUM rather
than HIGH, but the plan's deterministic/process boundary and its claimed
single-use production guarantee are not yet accurate. Review 2 cannot accept
that residual risk on behalf of the human.

Required resolution:

Define the approval time predicate completely, including the relationship to
the native comment timestamp and `approved_at <= current_time < expires_at`.
For one-time use, either define a durable pre-mutation claim and deterministic
selection/serialization/replay rule that prevents concurrent reuse, or state
plainly that GitHub-comment consumption is a non-atomic human/process control
and remove the deterministic/fail-closed single-use claim. Bind any claim and
outcome to the approval URL, release-scope hash, merged SHA, environment, and a
unique execution identity, and add focused fixtures for future-dated approval,
expired approval, prior consumption, concurrent claims, and retry after
failure.

### Assessment of the emphasized risks

- Active-only cleanup is now coherent with post-merge handling: repository
  state stops before merge and GitHub records carry later disposition. This is
  procedural, not CI-enforced, and the plan says so.
- CI and validator claims are appropriately narrowed. Current master
  protection still requires the `API` status but requires zero approving
  reviews, does not dismiss stale reviews, and does not require last-push
  approval; revision 2 does not present those settings as guarantees or propose
  to mutate them.
- The proposed active files no longer form a per-task command/progress/evidence
  receipt system. The remaining pre-merge bindings are proportionate to Issue
  #16's explicit approval/review/staleness requirements.
- The detailed workflow has one declared durable contract. The roadmap keeps
  durable product/phase outcomes, Issues keep operational work, templates encode
  schemas, and routing/human documents have narrower responsibilities; no
  conflicting domain, AI, security, or ADR source is created.
- The eight-state local lifecycle and stateless standard-library validator are
  substantial but bounded: no service, database, daemon, archive, automatic
  GitHub mutation, or release executor is proposed. Apart from the production
  validity gap above, revision 2 does not cross the Issue's workflow-engine
  non-goal.

### Evidence actually inspected

- Complete `WORK.md`, preserved Review 1, root/scoped engineering contracts,
  `CONTRIBUTING.md`, `SECURITY.md`, `README.md`, `docs/domain.md`,
  `docs/domain-ai.md`, `docs/roadmap.md`, every tracked ADR and the ADR index,
  current repository tree/status/history, `.github/workflows/ci.yml`, the full
  `Makefile`, and relevant API/deployment verification and mutation scripts.
- Current HEAD, local `master`, and `origin/master` all resolve to
  `92419ed6144e8d07b1392e0b931cdf32a044af77`; `work/` is the sole untracked
  tree and no ENG-016 implementation exists.
- `make help`, `git diff --check`, plan-hash recomputation, Issue-digest
  recomputation, and Review 1 preservation checks completed successfully. No
  application, governance implementation, external, deployment, or production
  mutation was performed.
- Read-only GitHub state observed for this review: Issue #16 is open with no
  comments; no PR exists for `chore/ai-engineering-harness`; Actions run
  `33964729866` is a successful `CI / API` run for the exact base SHA; current
  master protection requires `API`, enforces protection for administrators and
  conversation resolution, disallows force-push/deletion, and retains the
  review-setting limitations stated above.

### Verdict

`CHANGES_REQUIRED`

Revision 2 resolves seven of the eight Review 1 findings. The remaining open
MEDIUM finding `ENG-016-PLAN-006` requires planner resolution or an exact
human-authored residual-risk acceptance; this reviewer does not accept it on
the human's behalf. The plan must not advance to approval or implementation in
its current form.

## Review 3

- Review ID: `ENG-016-PLAN-REVIEW-003`
- Review type: `PLAN`
- Recorded at: `2026-09-06T01:36:39Z`
- Actor: `Codex plan reviewer r3`
- Actor role: `independent plan reviewer`
- Session label: `ENG-016-plan-review-r3-2026-09-06`
- Plan author actor label: `Codex planner`
- Plan author session label: `ENG-016-plan-r3-2026-09-06`
- Fresh-session attestation: I derived this review independently from the
  canonical Issue, current repository and GitHub state, the complete revision-3
  plan, the complete preserved Reviews 1 and 2, and the applicable governance,
  release, deployment, and external-mutation contracts. I treated earlier
  conclusions and the planner's resolution matrix as claims to verify, not as
  evidence.
- Subject work item: `ENG-016`
- Subject plan revision: `3`
- Exact reviewed plan checkpoint SHA:
  `sha256:49e2e6fbfe3cc60dd0310147b29cafc9af2113a8dd37e7f44315803ca2b890da`
- Independently recomputed normative SHA:
  `sha256:49e2e6fbfe3cc60dd0310147b29cafc9af2113a8dd37e7f44315803ca2b890da`
- Exact Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Independently recomputed Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Full reviewed `WORK.md` SHA-256:
  `85997bb7c01fec976f3b25c87ab2bc3250ec7dbf200282080f8b1624890b554a`
- Reviews 1–2 pre-append file SHA-256:
  `659cdff29b735b9b7a61a42eebc80743306b7cb8dcbbba5de8027fa05e3bde83`
- Reviews 1–2 pre-append Git blob:
  `a48fcdc254ff20adf42a77425e78bae93b81c4e7`
- Reviewed Git base/HEAD: `92419ed6144e8d07b1392e0b931cdf32a044af77`
- Verdict: `CHANGES_REQUIRED`

### Integrity and input verification

- Both the mutable header and normative checkpoint declare plan revision `3`.
- Independent byte-level normalization of the single normative marker section
  produced 42,941 bytes and exactly reproduced the declared revision-3 plan
  hash.
- GitHub Issue #16 was re-fetched read-only. It is open, has no comments, and
  its current title, body, URL, and observed `updatedAt` of
  `2026-09-05T18:16:42Z` match the ENG-016 checkpoint. The normalized 58-byte
  title and 3,925-byte body exactly reproduced the declared Issue digest.
- Before this append, `REVIEWS.md` contained exactly Review 1 and Review 2 in
  order. The first 376 lines exactly match archived Git blob
  `1f1e0f87c62dd6244d5e64490d2cd097d9496b84`, whose SHA-256 is the Review 1
  pre-append hash recorded by Review 2. The complete Reviews 1–2 file exactly
  matched the pre-existing Git blob and full-file SHA recorded above, and its
  filesystem modification/change time preceded revision 3 of `WORK.md`. No
  alteration of either preserved review was detected. Review 2 did not record
  its own prior self-hash, so the combined pre-append blob and filesystem
  checkpoint are the strongest repository-local preservation evidence.

### Executive assessment

Revision 3 substantially improves `ENG-016-PLAN-006`. It binds approval,
claim, and result records to the canonical Issue, exact merged SHA, exact target
environment, release-scope hash, approval URL, and execution identity. It
handles visible prior, multiple, malformed, foreign, missing, and
incompletely-paginated claims conservatively; consumes attempted executions;
requires new approval for retry; and separates in-scope recovery from new or
out-of-scope operations.

It also accurately states that the comment claim/re-fetch sequence is not an
atomic lock and cannot deterministically prevent truly concurrent execution.
Production checks remain live process procedures rather than local-validator
claims, and no service, daemon, queue, deployment coordinator, workflow engine,
or deployment implementation is introduced. The plan remains compatible with
the repository's current non-production-ready deployment boundary and requires
fresh external-state inspection and separate mutation authority.

The finding is nevertheless not fully resolved. Native GitHub comments are
editable and expose distinct creation and update timestamps, but revision 3
treats only `created_at` as the effective approval time. It also leaves expiry
behavior internally inconsistent once production mutation has started. These
are validity and failure-recovery semantics within the original MEDIUM finding,
not accepted residual risk.

### Finding

#### ENG-016-PLAN-006 — MEDIUM — Edited-approval timing and post-start expiry semantics remain incomplete

Status: `OPEN` / blocking until resolved or explicitly accepted by a human
under the exact residual-risk procedure.

Evidence:

- `WORK.md:580-618` declares the approval comment's native GitHub `created_at`
  authoritative and checks only `created_at <= current time < expires_at`.
  GitHub's official REST contract exposes an Update Issue Comment operation
  and returns both `created_at` and `updated_at`. An authorized human can create
  a placeholder or different comment at time T0 and edit its body into the
  approval at T1, or later extend `expires_at`, while `created_at` remains T0.
  Revision 3 would misstate T0 as the approval time and does not bind the claim
  to the exact native comment version. See
  <https://docs.github.com/en/rest/issues/comments?apiVersion=2022-11-28#update-an-issue-comment>.
- `WORK.md:601-618` says a valid, unexpired approval is required before
  performing **any** production mutation. `WORK.md:678-684` and the decision
  row at `WORK.md:714` separately allow already-enumerated rollback/remediation
  to continue under the same execution after mutation begins. If approval
  expires after a partial mutation, the first rule says to stop before the
  recovery mutation while the latter rule says recovery may continue. The
  plan does not define whether expiry is an admission deadline for claim/first
  mutation or a cutoff applied to every forward and recovery operation.
- The decision table covers future-inconsistent and expired creation-time
  inputs, expiry before claim/mutation, and retry after mutation, but it has no
  case for approval-body edits, expiry between claim and first mutation, or
  expiry after a partial mutation requiring in-scope recovery.

Impact:

The live procedure can assign the wrong effective approval time after a native
comment edit and can reach contradictory instructions during a failed or
partial deployment that crosses the expiry instant. Exact SHA, environment,
scope, and execution bindings limit the blast radius, and the current
repository is not production-ready, so this remains MEDIUM rather than HIGH.
It still blocks approval because no human-authored residual-risk acceptance
exists.

Required resolution:

Define one native-comment version rule. For example, require an eligible
production-approval comment to be unedited (`created_at == updated_at`) and
require a new comment for every correction or extension, or use native
`updated_at` as the effective approval time and bind the claim to that exact
comment version. Add contract cases for an approval granted or extended by
edit.

Then define expiry precedence explicitly: either make expiry an admission gate
through the first production mutation while permitting only pre-approved
rollback/remediation to finish the already-started execution, or apply expiry
before every mutation and specify a safe recovery path that does not conflict
with the rollback rule. Cover expiry after claim but before first mutation and
expiry after partial mutation in the procedure cases.

### Regression check of Review 2 resolved findings

| Finding | Review 3 status | Independent regression basis |
|---|---|---|
| `ENG-016-PLAN-001` | `RESOLVED` | The local lifecycle still ends at `READY_TO_MERGE`, cleanup remains pre-merge, and both post-merge release choices retain explicit GitHub-backed procedures. |
| `ENG-016-PLAN-002` | `RESOLVED` | CI/validator guarantees remain narrowly stated; absent active artifacts, self-modification, and lack of tamper-proof remote enforcement remain disclosed rather than claimed solved. |
| `ENG-016-PLAN-003` | `RESOLVED` | Revision 3 adds release protocol fields only to native GitHub comments; it does not restore a repository command/progress/evidence ledger or completed-work archive. |
| `ENG-016-PLAN-004` | `RESOLVED` | Relational reviewer labels and direct human-authored approval provenance remain unchanged, with identity/intent limits explicit. |
| `ENG-016-PLAN-005` | `RESOLVED` | The normalized Issue title/body digest remains normative and is still re-fetched before implementation and PR validation. |
| `ENG-016-PLAN-007` | `RESOLVED` | The durable roadmap phase summary and exact Issue mapping still precede temporary phase cleanup. |
| `ENG-016-PLAN-008` | `RESOLVED` | No ADR or ADR number is proposed; any later architectural decision remains a material deviation selected from then-current tracked state. |

No revision-3 change reopened a finding that Review 2 marked `RESOLVED`.

### Assessment of the emphasized risks

- Native timestamp and approval expiry semantics retain the MEDIUM gap above.
- Exact SHA/environment/release-scope binding is coherent across the release
  evaluation, approval, claim, and result records.
- Visible prior and ambiguous claims, incomplete pagination, lost claim
  acknowledgement, failed/aborted outcomes, and retries fail closed. GitHub
  comment edit/deletion resistance is a process/platform trust boundary, not a
  deterministic guarantee; the executor is explicitly prohibited from
  rewriting claim/result history.
- The non-atomic concurrent-executor limitation is stated plainly and is not
  disguised as deterministic single-use enforcement.
- Deterministic local checks and live process-only release checks remain
  clearly separated.
- No workflow engine or deployment system is introduced, and ENG-016 itself
  retains `release_required=false` with no external or production mutation.

### Evidence actually inspected

- Complete revision-3 `WORK.md`; complete preserved Reviews 1 and 2; the exact
  revision-2-to-revision-3 diff recovered from Git blob
  `e3ce2dad6e53b156a3291bf576073fda7bdd6ba6`; root and scoped `AGENTS.md`;
  `CONTRIBUTING.md`; `SECURITY.md`; `README.md`; `docs/roadmap.md`; the ADR
  index and deployment-relevant accepted ADRs; `.github/workflows/ci.yml`;
  `Makefile`; `deploy/AGENTS.md`; production Compose; and deployment,
  backup/restore, and safety scripts.
- Current HEAD, local `master`, and `origin/master` all resolve to
  `92419ed6144e8d07b1392e0b931cdf32a044af77`. The branch is
  `chore/ai-engineering-harness`; `work/` is the sole untracked tree; no ENG-016
  implementation or PR exists.
- Read-only GitHub state showed Issue #16 open with no comments, the successful
  `CI / API` run `33964729866` for the exact base SHA, and current master
  protection still requiring `API` while requiring zero approvals and neither
  stale-review dismissal nor last-push approval. These are observations, not
  guarantees, and no remote state was mutated.
- `git diff --check`, plan-hash recomputation, Issue-digest recomputation,
  preserved-review comparisons, repository-state inspection, and the
  revision-2-to-3 regression diff completed successfully. No implementation,
  external, deployment, or production operation was performed.

### Verdict

`CHANGES_REQUIRED`

Revision 3 does not reopen any finding Review 2 resolved and it resolves most
of the original production-approval gap, including exact bindings, conservative
claim/retry handling, and honest non-atomic concurrency disclosure. The
remaining MEDIUM `ENG-016-PLAN-006` timestamp-edit and expiry/recovery ambiguity
requires planner resolution or exact human-authored residual-risk acceptance.
This reviewer does not accept that risk on the human's behalf. The plan must
not advance to approval or implementation in its current form.

## Review 4

- Review ID: `ENG-016-PLAN-REVIEW-004`
- Review type: `PLAN`
- Recorded at: `2026-09-06T09:50:21+08:00`
- Actor: `Codex plan reviewer r4`
- Actor role: `independent plan reviewer`
- Session label: `ENG-016-plan-review-r4-2026-09-06`
- Plan author actor label: `Codex planner`
- Plan author session label: `ENG-016-plan-r4-2026-09-06`
- Fresh-session attestation: I derived this review from the canonical Issue,
  complete revision-4 `WORK.md`, complete preserved reviews, applicable
  repository/release/deployment contracts, and current read-only repository
  and GitHub state. I did not treat the planner's resolution matrix or prior
  reviewer conclusions as evidence.
- Subject work item: `ENG-016`
- Subject plan revision: `4`
- Exact reviewed plan checkpoint SHA:
  `sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934`
- Independently recomputed normative SHA:
  `sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934`
- Exact Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Independently recomputed Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Full reviewed `WORK.md` SHA-256:
  `84c627c9c6fd35910370161934f2b729da9cb55643afc31de1ce07471e96637d`
- Reviews 1-3 pre-append SHA-256:
  `e04297c7f3b1b8479daee9aac5a64ea65da19c2f830c142de0c8583895a27a2f`
- Reviewed Git base/HEAD: `92419ed6144e8d07b1392e0b931cdf32a044af77`
- Verdict: `APPROVED`

### Integrity and input verification

- Both the mutable header and normative checkpoint declare revision `4`. The
  unique marker section normalized to 47,508 UTF-8 bytes under the specified
  LF, trailing-whitespace, and blank-boundary rules and reproduced the exact
  declared plan hash above.
- I re-fetched GitHub Issue #16 read-only. It is open at the canonical URL,
  has no comments, and has observed `updatedAt`
  `2026-09-05T18:16:42Z`. Its normalized 58-byte title and 3,925-byte body
  reproduce the exact Issue specification digest above. Comments and
  `updatedAt` were not used in that digest.
- Reviews 1 and 2 remain byte-for-byte preserved: the current Review 1 prefix
  hashes to `dfea908ba74fdaa841666918da19eada9013a047ffd700eebeba200841492ec7`
  and equals Git blob `1f1e0f87c62dd6244d5e64490d2cd097d9496b84`; the current
  Reviews 1-2 prefix hashes to
  `659cdff29b735b9b7a61a42eebc80743306b7cb8dcbbba5de8027fa05e3bde83` and
  equals Git blob `a48fcdc254ff20adf42a77425e78bae93b81c4e7`.
- Review 3 remains the single contiguous Review 3 record after that preserved
  prefix, ending immediately before this append; its current standalone SHA-256
  is `85882eb2feb0fc8fffd0036af8acad6988d86a8e7f83d6acb27c8a12efc2614a`.
  Review 3 did not leave a prior self-hash or Git blob for a stronger historic
  byte comparison. The current `REVIEWS.md` modification time precedes the
  revision-4 `WORK.md` modification time, and no alteration was detected.
  As with any untracked active artifact, this cannot cryptographically prove a
  modification made before this review's observable checkpoint.

### Independent assessment of ENG-016-PLAN-006

The revision-4 release procedure resolves the remaining production-approval
finding. It is a live GitHub/release-agent procedure, expressly not a local
validator claim or an ENG-016 production authorization.

- An eligible approval is one exact native comment version: ID, URL,
  `created_at`, `updated_at`, and normalized complete-body digest. It must be
  unedited (`created_at == updated_at`); a correction or expiry extension is a
  new human-authored approval comment, not an edit. Native `created_at` is the
  authoritative approval time and must satisfy
  `created_at <= current time < expires_at` and `created_at < expires_at`.
- Before claim and again immediately before the first production mutation, the
  executor re-fetches the canonical Issue, release scope, exact approval, and
  all required comment state. Exact work-item, merged-SHA, environment,
  release-scope-hash, decision, provenance, version, time, scope, and external
  authority bindings must match. Missing, edited, deleted, malformed,
  future-dated, expired, mismatched, incomplete, or ambiguous state stops.
- A pre-claim read, one UUID execution ID, a bound `DEPLOYMENT_CLAIM`, and an
  immediate complete re-fetch require the executor's claim to be the sole
  valid claim for that exact approval version. Claims include the native
  approval timestamps and body digest. Prior, foreign, multiple, malformed,
  missing, result-only, or ambiguously paginated claim/result state blocks
  mutation and requires human resolution or a new approval as applicable.
- Expiry before claim stops. Expiry after claim but before the first mutation
  produces the bound `ABORTED` result with
  `APPROVAL_EXPIRED_BEFORE_MUTATION`, performs no production mutation, and
  requires a new approval and execution ID. At the first mutation the claim
  and approval are consumed. Thereafter every forward mutation still requires
  an in-window approval; expiry after partial work stops all forward mutation.
  Only rollback/remediation already enumerated in the approved release scope
  may continue after expiry, solely as needed to reach a safe state, and must
  finish with `FAILED` or `ABORTED`.
- Edit, deletion, unavailability, or version mismatch after mutation starts
  likewise stops forward work and allows only that pre-approved safe recovery.
  A failed, partial, expired, edited, deleted, or pre-mutation-aborted attempt
  cannot reuse its approval. A subsequent attempt requires a new exact
  approval version, execution ID, and a release scope reassessed against the
  current production state.
- The procedure accurately discloses that Issue comments are not an atomic
  distributed lock: two genuinely concurrent executors can still observe
  apparently valid state. It makes serialization a human/process boundary and
  adds no lock service, queue, coordinator, daemon, database, or workflow
  engine. This is an explicit limitation, not an unaccepted MEDIUM residual
  risk or a claim of deterministic single-use enforcement.

The production decision table exercises each of the above cases, including
edited approvals, new correction comments, edited expiry extensions,
future-dated inputs, expiry before claim, expiry between claim and first
mutation, expiry after partial mutation, in-scope recovery, attempted forward
work after expiry, post-start edit/deletion, retry/new approval after partial
state, exact SHA/environment/scope mismatches, claim ambiguity, and unavailable
GitHub re-fetches. It is compatible with the repository's separate explicit
external-mutation authority, immutable-image requirement, no automatic
destructive rollback, forward-only migration recovery, and current
non-production-ready ingress boundary.

### Regression review of previously resolved findings

| Finding | Status | Independent basis |
|---|---|---|
| `ENG-016-PLAN-001` | `RESOLVED` | The local lifecycle still stops at `READY_TO_MERGE`; active files are removed before merge while native Issue/PR/Actions records and both explicit GitHub-backed release paths carry post-merge disposition. |
| `ENG-016-PLAN-002` | `RESOLVED` | Validator/CI scope remains intentionally limited to local machine-verifiable artifacts. Missing active artifacts, self-modifying checks, and remote enforcement limits remain disclosed rather than represented as tamper-proof gates. |
| `ENG-016-PLAN-003` | `RESOLVED` | Active files retain only plan/review bindings and concise summaries. The release protocol uses native GitHub comments, not a local command log, receipt ledger, progress journal, or completed-work archive. |
| `ENG-016-PLAN-004` | `RESOLVED` | Reviewer/author actor and session labels remain relational requirements; human plan approval is an exact canonical-Issue comment URL, with identity, intent, and freshness honestly left as process trust boundaries. |
| `ENG-016-PLAN-005` | `RESOLVED` | The normalized canonical Issue title/body digest remains normative, is recomputable, and is re-fetched at implementation and PR validation; a mismatch returns the work to planning. |
| `ENG-016-PLAN-007` | `RESOLVED` | The durable roadmap transfer requires phase revision/hash, outcomes, ordering rationale, exit criteria, exclusions, unresolved decisions, and exact Issue mapping before temporary phase cleanup. |
| `ENG-016-PLAN-008` | `RESOLVED` | No ADR or ADR number is proposed or reserved. A subsequently needed architecture decision is a material deviation and uses then-current tracked ADR state. |

No regression reopens a resolved finding. The plan is detailed where the
post-merge safety boundary needs exact native GitHub semantics, but remains
lean in the relevant architectural sense: it proposes one documentation
contract, templates, a standard-library local validator, focused tests, and
existing-CI integration only. It does not introduce a workflow engine,
service, daemon, database, dashboard, archive, coordinator, paid dependency,
or ENG-016 production/deployment mutation.

### Verdict

`APPROVED`

Revision 4 resolves `ENG-016-PLAN-006` without an unaccepted MEDIUM residual
risk and retains the earlier resolved controls. This is a passing independent
plan-review verdict only. The work remains in `PLAN_REVIEW` until a separate
human-authored canonical-Issue plan-approval comment exactly binds revision 4
and its plan hash; this review grants no implementation, GitHub mutation,
deployment, or production authority.

## Review 5

- Review ID: `ENG-016-PLAN-REVIEW-005`
- Review type: `PLAN`
- Recorded at: `2026-09-06T02:42:22Z`
- Actor: `Codex plan reviewer r5`
- Actor role: `independent plan reviewer`
- Session label: `ENG-016-plan-review-r5-2026-09-06`
- Plan author actor label: `Codex planner`
- Plan author session label: `ENG-016-plan-r4-2026-09-06`
- Fresh-session attestation: I derived this review independently from the
  canonical Issue, complete revision-4 `WORK.md`, complete Reviews 1–3,
  historical Review 4, applicable repository and scoped contracts, product and
  architecture sources, CI, release, deployment, external-mutation surfaces,
  and current read-only repository/GitHub state. I did not use Review 4’s
  `APPROVED` verdict as evidence and did not perform implementation or remote
  mutation.
- Subject work item: `ENG-016`
- Subject plan revision: `4`
- Exact reviewed plan checkpoint SHA:
  `sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934`
- Independently recomputed normative SHA:
  `sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934`
- Exact Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Independently recomputed Issue specification digest:
  `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Full reviewed `WORK.md` SHA-256:
  `84c627c9c6fd35910370161934f2b729da9cb55643afc31de1ce07471e96637d`
- Review 1 preserved Git blob:
  `1f1e0f87c62dd6244d5e64490d2cd097d9496b84`; current prefix content matches
  it after removing one separator LF only.
- Reviews 1–2 preserved Git blob:
  `a48fcdc254ff20adf42a77425e78bae93b81c4e7`; current prefix content matches
  it after removing one separator LF only.
- Review 3 independently recomputed standalone SHA-256:
  `85882eb2feb0fc8fffd0036af8acad6988d86a8e7f83d6acb27c8a12efc2614a`
- Reviewed Git base/HEAD: `92419ed6144e8d07b1392e0b931cdf32a044af77`
- Verdict: `APPROVED`

### Integrity and input verification

- Both the mutable header and the normative checkpoint declare plan revision
  `4`. The exact normative marker section contains 47,508 normalized UTF-8
  bytes and independently reproduces the pinned plan hash.
- The complete current `WORK.md` SHA matches the previously reviewed full-file
  SHA. No material plan drift was found.
- GitHub Issue #16 was fetched read-only from its canonical URL. It is open;
  the normalized title is 58 bytes, the normalized body is 3,925 bytes, and
  the independently recomputed title/body digest matches the pinned Issue
  digest. Its current `updatedAt` changed because a comment was added, which
  does not affect the intentionally comment-excluded digest.
- The live Issue has one unedited comment by `farisakbar28` that semantically
  identifies a revision-4 plan approval, but its `PLAN_APPROVAL`/`Work item:`
  prose is not the exact `approval_type=PLAN` key-value schema required by the
  plan. I did not treat it as a valid plan-approval binding or copy it into
  `WORK.md`; the current `Human plan-approval comment: NONE` remains the safe
  local state. This is a current external gate observation, not evidence for
  this plan verdict.
- Reviews 1 and 2 match their preserved Git blobs exactly apart from one
  additional separator LF at the next-review boundary. Review 3’s contiguous
  record matches its previously demonstrated standalone SHA after excluding
  the separator LF before Review 4. No substantive drift was found in Reviews
  1–3.
- The current worktree contains only the two untracked ENG-016 active files;
  current `HEAD`, local `master`, `origin/master`, and the remote `master` tip
  are the pinned base SHA. No ENG-016 implementation files, durable harness
  files, candidate commit, PR, or remote feature branch exist.

### Independent plan assessment

The source-of-truth model is coherent. Current source and observed results
govern implementation facts; `docs/domain.md` and `docs/domain-ai.md` remain
canonical; accepted ADRs retain scoped authority; `docs/roadmap.md` remains
durable product/phase direction; GitHub Issues remain the operational backlog;
and active `WORK.md`/`REVIEWS.md` are temporary execution/review bindings that
cannot override those sources or the engineering contracts. The plan changes
no SIAKAD/LMS ownership, tenant, authorization, AI, data, migration, or
application behavior.

Product/Phase Planning has the required temporary proposal, independent review,
exact human approval, roadmap synchronization, exact planning-ID/Issue
mapping, explicit remote-write authority, human roadmap merge, and active-only
cleanup. The durable transfer includes outcomes, ordering rationale, exit
criteria, exclusions, unresolved decisions, revision/hash, and Issue mapping;
it does not create a completed phase archive or workflow engine.

The eight-state pre-merge lifecycle is internally coherent and stops at
`READY_TO_MERGE`. It handles plan-review failure, in-plan implementation fixes,
material deviations, exact candidate review, PR validation, cleanup-only
descendants, human-only squash merge, and native GitHub-backed post-merge
disposition without inventing repository-local post-merge states.

Revision/hash binding and Issue staleness are complete: the normative section
contains revision 4 and the exact plan hash, the canonical title/body digest is
normatively included, and the digest is recomputed before implementation and
PR validation. Comments and ordinary `updatedAt` changes are correctly
excluded. Human plan approval is required as an exact, directly human-authored
canonical-Issue comment; reviewer/author actor and session labels are
relationally distinct and fresh-session requirements are explicit. Their
identity, intent, and freshness limitations are honestly left as process
trust boundaries.

Material-deviation rules fail closed for scope, acceptance, verification,
dependencies, behavior, architecture, safety, authority, risk, Issue
specification, external/destructive authority, and recovery changes. The
revision/hash/independent-review/new-human-approval reset is explicit. Finding
severity is also coherent: CRITICAL/HIGH block; MEDIUM requires resolution or
exact human residual-risk acceptance; LOW is normally non-blocking; and this
review accepts no residual risk for the human.

The validator remains a small standard-library local artifact checker. Its
enumerated scope covers hashes, bindings, cardinality, enums, state-required
fields, actor/session inequality, finding gates, issue/comment inputs supplied
freshly, forbidden archives/post-merge states, and cleanup-only Git diffs. CI
adds validator tests/checks to the existing `API` job while preserving its
identity, triggers, permissions, and existing application checks. The plan
explicitly disclaims proof of reviewer freshness, human intent, live
production state, deployment serialization, self-protection, and tamper-proof
remote governance. The live master protection state likewise has `API` as its
required status but zero required approving reviews, no stale-review dismissal,
and no last-push approval; the plan documents and does not mutate those
external limitations.

Protected-control-plane treatment is semantic, with conservative defaults for
workflow/validator/CI paths and no automatic CRITICAL classification merely
from changing ordinary `Makefile`, `CONTRIBUTING.md`, or `AGENTS.md` files.
Remote authority is explicit for commits, pushes, PR/Issue mutations, merge,
external/cloud operations, production, and destructive data work.

Both conditional release paths are specified. The `false` path records native
post-merge disposition and performs no release/deployment. The `true` path
defines exact release scope, immutable artifacts, human production approval,
native comment ID/URL/timestamps/body digest, unedited-version and expiry
rules, live re-fetches, pre-claim and sole-claim checks, exact execution IDs,
result records, retry/partial/failure handling, post-start expiry behavior,
safe in-scope recovery, and new approval requirements. It explicitly states
that GitHub comments are not an atomic lock and leaves concurrent-executor
serialization as a declared process boundary, without adding a coordinator,
queue, daemon, database, or workflow engine.

The procedure is compatible with the current deployment contracts: it requires
fresh external-state checks and explicit authority, preserves immutable image
references and readiness validation, permits only explicitly approved safe
recovery, does not imply database-migration rollback, and leaves the stale
legacy Caddy/Origin-CA ingress topology unchanged. No deployment, cloud,
production, or external mutation is part of ENG-016.

### Coverage of Issue requirements and non-goals

The Issue requirements are all covered by the plan: canonical domain/AI/
security/roadmap/ADR sources and the GitHub Issue backlog are preserved;
existing planning IDs remain identities; Product/Phase Planning is introduced;
independent PLAN and IMPLEMENTATION reviews use fresh sessions; human plan
approval binds exact revision/hash; material deviations reset planning;
severity gates and residual-risk provenance are defined; active artifacts are
temporary and reviews append while active; implementation verdicts bind exact
candidate SHAs and stale on relevant changes; release evaluation is
conditional; production approval is explicit; merge is human-only; push/PR and
local-commit authority is explicit; external/cloud mutations fail closed and
remain separately scoped; validator/CI coverage and limitations are stated; no
workflow engine is added; existing CI, security, tenancy, branch-protection,
and governance behavior is not weakened; and the normal prompt surface remains
short and repository/GitHub-state driven.

The non-goals are also preserved: no LMS behavior, authentication or
tenant-context implementation, database migration, production deployment,
replacement of canonical documents/ADRs/roadmap, or branch-protection
mutation; no duplicate IDs, completed-work archive, progress/evidence ledger,
automatic commit/push/PR/merge/release/deployment/rollback, external/cloud
mutation, paid fallback, workflow service, or `.agent/**` change is proposed.

### Regression review of ENG-016 findings 001–008

| Finding | Status | Independent regression basis |
|---|---|---|
| `ENG-016-PLAN-001` | `RESOLVED` | The local lifecycle ends before merge; active cleanup is pre-merge; native Issue/PR/Actions records and separate `release_required=false`/`true` procedures carry post-merge disposition. |
| `ENG-016-PLAN-002` | `RESOLVED` | Validator/CI claims are limited to present machine-checkable artifacts; absent active artifacts, self-modifying checks, and current remote non-enforcement are explicitly disclosed. |
| `ENG-016-PLAN-003` | `RESOLVED` | Active files contain minimum bindings rather than command receipts, progress logs, evidence indexes, or completed-work archives; native Git/GitHub/CI remain evidence sources. |
| `ENG-016-PLAN-004` | `RESOLVED` | Reviewer/author actor and session labels are relationally distinct, fresh sessions are required, and human approval must be a directly authored canonical-Issue comment; non-cryptographic limits are explicit. |
| `ENG-016-PLAN-005` | `RESOLVED` | The normalized title/body digest is normative and re-fetched at implementation and PR gates; any specification mismatch returns work to planning. |
| `ENG-016-PLAN-006` | `RESOLVED` | Revision 4 requires an unedited native approval version, native creation-time/expiry predicate, pre-claim and re-fetch semantics, explicit pre-mutation expiry abort, post-start expiry/recovery rules, exact claim/result bindings, retry invalidation, and an honest non-atomic concurrency limitation. |
| `ENG-016-PLAN-007` | `RESOLVED` | Phase cleanup requires durable roadmap fields and exact Issue mapping for outcomes, order/rationale, exit criteria, exclusions, unresolved decisions, and revision/hash. |
| `ENG-016-PLAN-008` | `RESOLVED` | No ADR number is proposed or reserved; any later architectural decision must use then-current tracked ADR state as a material deviation. |

No CRITICAL, HIGH, MEDIUM, or LOW finding remains open. This fresh result is
independent of Review 4; Review 4 is historical context only and is not used as
the prerequisite for implementation.

### Verdict

`APPROVED`

Revision 4 is approved as a fresh independent PLAN review. This verdict does
not constitute human plan approval, implementation authority, commit/push/PR
authority, merge authority, release authority, deployment authority, or
production authority. Before implementation, the human must supply the exact
canonical-Issue approval schema required by the plan and the implementation
agent must re-fetch the Issue, approval, plan hash, and current repository
state.

## Review 6

review_id=ENG-016-IMPLEMENTATION-REVIEW-001
review_type=IMPLEMENTATION
type=IMPLEMENTATION
work_item_id=ENG-016
plan_revision=4
plan_hash=sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934
issue_digest=sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4
implementation_candidate_sha=959307b75fffa6986c2587af763f1ce2dfbcd3a7
candidate_git_sha=959307b75fffa6986c2587af763f1ce2dfbcd3a7
candidate_parent_sha=92419ed6144e8d07b1392e0b931cdf32a044af77
reviewer_actor_label=Codex implementation reviewer r1
actor_label=Codex implementation reviewer r1
fresh_reviewer_session_label=ENG-016-implementation-review-r1-2026-09-06
session_label=ENG-016-implementation-review-r1-2026-09-06
implementation_author_actor_label=Codex implementer
implementation_author_session_label=ENG-016-implement-r4-2026-09-06
fresh_session_attestation=Fresh independent implementation-review session; I reconstructed the requirements from the repository and live read-only GitHub state, inspected the exact candidate and complete diff, and did not rely on planner, implementer, prior plan-review conclusions, or delegated work.
material_deviation=NONE

findings:
- finding_id=ENG-016-IMPL-001; severity=HIGH; status=OPEN; summary=The READY_FOR_PR implementation-review gate accepts a block based only on prose review type, the candidate SHA appearing anywhere, and an APPROVED verdict. It does not require or bind the work item, current plan revision/hash, Issue digest, reviewer/author labels, or fresh-session attestation; an independently mutated fixture with no reviewer identity or bindings returned no validator errors.
- finding_id=ENG-016-IMPL-002; severity=HIGH; status=OPEN; summary=Approval URLs are validated only against a generic GitHub URL pattern and the mutable WORK.md value, not the canonical Issue owner/repository/number. The supplied fixture's approval URL is github.com/example/... and the exact-schema matcher accepts it, so a foreign-repository comment can satisfy the local approval check.
- finding_id=ENG-016-IMPL-003; severity=HIGH; status=OPEN; summary=Residual-risk validation checks candidate_git_sha syntax but not equality with the current candidate or an exact comment URL. An independently supplied residual-risk comment for a different candidate SHA returned no errors, allowing a risk acceptance to bind to the wrong implementation.
- finding_id=ENG-016-IMPL-004; severity=HIGH; status=OPEN; summary=The finding gate treats unsupported status=CLOSED as resolving HIGH/CRITICAL findings, although the approved schema permits only OPEN, RESOLVED, or ACCEPTED_RESIDUAL_RISK and HIGH/CRITICAL findings must block. An independently mutated HIGH/CLOSED fixture returned no validator errors.
- finding_id=ENG-016-IMPL-005; severity=MEDIUM; status=OPEN; summary=Fresh standard GitHub Issue JSON is rejected because validation compares its API url field with the web html_url form while ignoring html_url. Running the validator with raw gh api Issue #16 JSON failed on URL mapping despite the correct number and digest, so the documented live-input path cannot pass.
- finding_id=ENG-016-IMPL-006; severity=HIGH; status=OPEN; summary=Cleanup validation accepts deletion of multiple active work-item pairs, not exactly the reviewed work item, and archive rejection covers only names directly below work/, so work/phases/archive passes. Independent temporary-repository checks returned no errors for both cases.
- finding_id=ENG-016-IMPL-007; severity=MEDIUM; status=OPEN; summary=The 10 bundled tests do not cover the required negative paths for current review bindings, missing reviewer fields, canonical approval URLs, residual candidate binding, finding-status enums, phase/archive handling, malformed inputs, or multi-item cleanup; the reproduced false passes above are consequently untested.

verdict=CHANGES_REQUIRED

- Review ID: `ENG-016-IMPLEMENTATION-REVIEW-001`
- Review type: `IMPLEMENTATION`
- Actor: `Codex implementation reviewer r1`
- Session label: `ENG-016-implementation-review-r1-2026-09-06`
- Implementation author actor label: `Codex implementer`
- Implementation author session label: `ENG-016-implement-r4-2026-09-06`
- Fresh-session attestation: This record is from a fresh independent implementation-review session and binds to the exact candidate above.
- Candidate Git SHA: `959307b75fffa6986c2587af763f1ce2dfbcd3a7`
- Subject work item: `ENG-016`
- Plan revision: `4`
- Plan hash: `sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934`
- Issue digest: `sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4`
- Verdict: `CHANGES_REQUIRED`

## Review 7

review_id=ENG-016-IMPLEMENTATION-REVIEW-002
review_type=IMPLEMENTATION
type=IMPLEMENTATION
work_item_id=ENG-016
plan_revision=4
plan_hash=sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934
issue_digest=sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4
implementation_candidate_sha=4491883fde2f108802954602189bcd02878c3923
candidate_git_sha=4491883fde2f108802954602189bcd02878c3923
candidate_parent_sha=925b483dc043a43c24bc22845cc5a74705d6eb51
reviewer_actor_label=Codex implementation reviewer r2
actor_label=Codex implementation reviewer r2
fresh_reviewer_session_label=ENG-016-implementation-review-r2-2026-09-06
session_label=ENG-016-implementation-review-r2-2026-09-06
implementation_author_actor_label=Codex implementer
implementation_author_session_label=ENG-016-implement-r4-2026-09-06
fresh_session_attestation=Fresh independent implementation-review session; I re-fetched the canonical Issue and approval read-only, recomputed the plan and Issue bindings, inspected the complete base-to-candidate diff, reproduced the prior findings independently, and did not rely on prior reviewer conclusions as evidence.
material_deviation=NONE
recorded_at=2026-09-06T04:36:37Z

findings:
finding_id=ENG-016-IMPL-001
severity=HIGH
status=RESOLVED
summary=READY_FOR_PR now requires one complete implementation review with current work-item, plan, Issue, candidate, author, reviewer, session, attestation, and APPROVED bindings.

finding_id=ENG-016-IMPL-002
severity=HIGH
status=RESOLVED
summary=Approval references and supplied Issue JSON now bind web/API URLs and repository/Issue identity to the canonical Work mapping.

finding_id=ENG-016-IMPL-003
severity=HIGH
status=RESOLVED
summary=Residual-risk comments now bind the exact current candidate, finding, work item, plan, and canonical comment URL with matching native comment identity.

finding_id=ENG-016-IMPL-004
severity=HIGH
status=RESOLVED
summary=Finding status is allow-listed and unsupported states cannot resolve HIGH or CRITICAL findings.

finding_id=ENG-016-IMPL-005
severity=MEDIUM
status=RESOLVED
summary=Standard GitHub Issue JSON now distinguishes API url, html_url, and repository_url while retaining identity and digest checks.

finding_id=ENG-016-IMPL-006
severity=HIGH
status=RESOLVED
summary=Cleanup validation now requires exactly one reviewed active work item and rejects nested completed/archive structures under work/.

finding_id=ENG-016-IMPL-007
severity=MEDIUM
status=RESOLVED
summary=Focused tests now cover the required positive and negative paths for review bindings, URLs, residual risk, finding states, cleanup, archives, and malformed input.

finding_id=ENG-016-IMPL-008
severity=HIGH
status=OPEN
summary=The plan-review gate does not validate the prescribed key/value plan-review record or bind an approved legacy prose review to the current plan revision, plan hash, Issue digest, reviewer/author identities, or fresh session; a minimal unbound APPROVED plan review passes while a valid template-shaped record is rejected.

verdict=CHANGES_REQUIRED

## Review 8

review_id=ENG-016-IMPLEMENTATION-REVIEW-003
type=IMPLEMENTATION
work_item_id=ENG-016
plan_revision=4
plan_hash=sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934
issue_digest=sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4
actor_label=Codex implementation reviewer r3
session_label=ENG-016-implementation-review-r3-2026-09-06
implementation_author_actor_label=Codex implementer
implementation_author_session_label=ENG-016-implement-r4-2026-09-06
fresh_session_attestation=Fresh independent implementation-review session; I reconstructed the approved contract and current source, fetched the canonical Issue and approval read-only, inspected the exact candidate and ancestry, reran applicable verification, and did not rely on prior reviewer conclusions as evidence.
candidate_git_sha=e32e035461ac3b83c1887e88f1d932e21bf1d730

finding_id=ENG-016-IMPL-001
severity=HIGH
status=RESOLVED
summary=The READY_FOR_PR gate still requires exactly one complete implementation review bound to the current work item, plan, Issue, candidate, implementation author, reviewer, session, attestation, and APPROVED verdict.

finding_id=ENG-016-IMPL-002
severity=HIGH
status=RESOLVED
summary=Approval and supplied Issue URL validation still binds web/API forms and repository/Issue identity to the canonical WORK.md mapping.

finding_id=ENG-016-IMPL-003
severity=HIGH
status=RESOLVED
summary=Residual-risk validation still binds the exact current candidate, finding, work item, plan, and canonical comment identity.

finding_id=ENG-016-IMPL-004
severity=HIGH
status=RESOLVED
summary=Finding status remains allow-listed and unsupported status values cannot resolve HIGH or CRITICAL findings.

finding_id=ENG-016-IMPL-005
severity=MEDIUM
status=RESOLVED
summary=Live GitHub Issue JSON continues to distinguish API url, html_url, and repository_url while retaining identity and digest checks.

finding_id=ENG-016-IMPL-006
severity=HIGH
status=RESOLVED
summary=Cleanup validation still requires exactly one reviewed active work-item pair and rejects nested completed/archive structures.

finding_id=ENG-016-IMPL-007
severity=MEDIUM
status=RESOLVED
summary=Focused tests and independent negative probes cover the required review-binding, URL, residual-risk, finding-state, cleanup, archive, and malformed-input paths.

finding_id=ENG-016-IMPL-008
severity=HIGH
status=OPEN
summary=The legacy PLAN-review adapter is not restricted to bootstrap Review 5: an independently appended later Review block containing the historical labelled fields is selected and accepted, allowing a future review to bypass the strict canonical key/value schema.

finding_id=ENG-016-IMPL-009
severity=HIGH
status=OPEN
summary=The validator does not require plan-author actor/session fields in WORK.md. Removing both fields lets a canonical review bind literal None values through str(None), and validate_repository returns no error despite missing author identity bindings.

verdict=CHANGES_REQUIRED
