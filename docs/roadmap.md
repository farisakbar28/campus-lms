# campus-lms Engineering and Product Roadmap

Status: active product and engineering direction.

This document owns future direction, concise completed progress, and durable
material changes. GitHub Issues are the operational backlog;
[the current architecture](architecture.md) is the implementation snapshot;
Git and Pull Requests show delivered changes; tests and GitHub Actions hold
verification evidence. A capability is not implemented merely because it
appears below.

## Planning rules

Roadmap sequencing follows authorization, data ownership, safety, and
operational dependencies rather than calendar milestones. Future work must
preserve the contracts in `docs/domain.md`, `docs/domain-ai.md`, `SECURITY.md`,
accepted ADRs, and scoped engineering instructions. Applied migrations remain
immutable, security-sensitive configuration fails closed, and no automatic
paid fallback may be introduced.

## Completed progress

These milestones record durable outcomes, not current hosted CI, cloud, or
production status. Detailed task and review history remains in GitHub.

- **Engineering baseline (2026-09-05):** commit
  `92419ed6144e8d07b1392e0b931cdf32a044af77` established the tracked
  repository baseline. A successful hosted `CI / API` run was observed for
  that exact revision on 2026-09-05 at 12:00 UTC. This is historical evidence
  for that revision and job only; check GitHub Actions for current CI status.
- **Trusted tenant-context contract (2026-09-20):**
  [ADR-0006](adr/0006-trusted-authenticated-tenant-context.md) entered
  `master` through [PR #19](https://github.com/farisakbar28/campus-lms/pull/19)
  at `27eeecef73776ac26e47149ad9a3ab3414f14bb4`. It defines admission
  semantics for the later runtime composition.
- **Trusted tenant-context runtime composition (2026-09-23):**
  [AUTH-002 / PR #25](https://github.com/farisakbar28/campus-lms/pull/25)
  merged at `d904784ca32292a441100e648fcc3bb0badcae9d`. The running server now
  composes access-token verification, active session and membership admission,
  and typed Principal propagation for the protected roster route. Login,
  initial session creation, HTTP refresh flows, and identity-provider
  integration remain future work.
- **Repository-native workflow v2 (2026-09-22):**
  [ENG-017 / PR #22](https://github.com/farisakbar28/campus-lms/pull/22)
  merged at `4d70aea2b5a4cd67a332b713b569af4bc8b0fc26`, replacing the
  earlier v1 workflow harness. ENG-018 / [Issue #23](https://github.com/farisakbar28/campus-lms/issues/23)
  was delivered by [PR #24](https://github.com/farisakbar28/campus-lms/pull/24)
  and merged at `4330df83ed96e28485371cd0c9d979ab1fc81e34`.
- **Current-schema local recovery validation (2026-09-24, OPS-003):** the
  change set updates the local backup, disposable restore, normalized state,
  and source-safety checks to migration `0006`, including the global
  `auth_sessions` table and its current constraints. This remains local
  operational validation, not production recovery readiness or off-machine
  backup evidence.
- **Tenant-safe course content slice (2026-09-24, CONTENT-008):** the
  implementation adds migration `0007`, modules, lessons, materials, file
  metadata, protected API routes, staff/student offering authorization,
  publication audit records, availability filtering, and PostgreSQL/RLS
  integration coverage. Binary storage, signed delivery, malware scanning,
  and frontend workflows remain separate work.

## Durable material changes

- **2026-08-31 — bounded ingress:** accepted
  [ADR-0002e](adr/0002e-zero-domain-quick-tunnel-validation.md) superseded
  ADR-0002b only for its ingress, domain, Origin-CA, and normal
  administration assumptions. Its private-origin, recovery, and zero-cost
  boundaries remain. The accepted path is temporary validation, not
  permanent production ingress or evidence of a running tunnel.
- **2026-09-20 — tenant-context authority:** ADR-0006 established trusted
  admission semantics. The runtime Principal and authentication admission
  composition were subsequently delivered by AUTH-002 / PR #25; login and
  initial session entrypoints remain separate future work.
- **2026-09-22 — implementation snapshot ownership:** current code facts
  belong in [architecture.md](architecture.md), while this roadmap owns
  direction and durable progress. Current hosted CI status must be read
  from GitHub Actions for the relevant revision.

## Product boundary

`campus-lms` is a multi-tenant university Learning Management System. It is
not SIAKAD.

SIAKAD remains authoritative for academic master and reference data, including:

- programs;
- academic terms;
- the course master;
- course offerings;
- staff assignments; and
- enrollment.

The LMS remains authoritative for LMS activity and content records, including:

- course content;
- assignments and submissions;
- quizzes;
- attendance;
- progress;
- gradebook; and
- feedback.

`ATTENDANCE_SOURCE_OF_TRUTH=LMS`

Global person and user identity is distinct from tenant membership and tenant
roles. Tenant roles are distinct from course-offering roles. A course is
distinct from a course offering. Tenant ownership and isolation remain
explicit through application authorization, PostgreSQL RLS, and composite
tenant-consistency constraints.

Academic history is retained when membership or enrollment state changes.
Submission versions are immutable. A Teaching Assistant may draft a score,
while publication and locking must remain governed by the appropriate
authorized LMS or institutional workflow.

## AI boundary

AI is future, optional, staff-only, and non-authoritative. Students are not AI
users. Subject to an approved role and capability policy, authorized
institutional staff may receive assistance such as:

- attendance or participation lookup;
- summaries of authorized student answers;
- course or activity summaries; and
- authorized information navigation and analysis.

AI must never be authoritative for grades, attendance decisions, pass/fail,
disciplinary outcomes, enrollment decisions, or another outcome owned by an
authoritative LMS or institutional workflow.

Authorization, tenant checks, object access, and lifecycle checks must happen
before retrieval, prompt construction, model calls, or tool calls. Retrieved
documents and submissions are untrusted prompt content and cannot override
system or authorization policy. Student-visible material originating from AI
must pass through ordinary LMS review and publication workflows. Ephemeral
staff assistance may be shown to its authorized requester without a second
human approval workflow, but remains non-authoritative.

The exact role allowlist remains unresolved:

`DECISION_REQUIRED: AI_ROLE_ALLOWLIST`

The institutional data-processing boundary also remains unresolved. No AI
tables, vector indexes, providers, models, orchestration frameworks, tracing
platforms, or student-facing AI workflow is committed by this roadmap.

## Infrastructure and operations boundary

The accepted Quick Tunnel path is bounded temporary validation only. It is not
permanent production ingress. The current deployment implementation still
contains legacy 8443 and Origin-CA assumptions and requires separately
authorized reconciliation. Permanent production ingress remains a separate
future architecture decision.

The local backup and restore validation now targets migration `0007` and the
current application tables, including global `auth_sessions` and the first
course-content tables. It verifies a distinct disposable restore and
source-state stability, but does not establish off-machine retention,
production recovery readiness, or provider state.

Historic external infrastructure, provider, region, quota, storage, and
capacity facts require fresh revalidation before reliance. The project retains
the zero-incremental-paid-infrastructure constraint and must fail closed,
defer, or disable a capability when verified zero-cost capacity is exhausted.

## Dependency-oriented capability roadmap

### 1. Engineering baseline

Preserve the runnable Go/PostgreSQL repository, explicit migrations, focused
tests, safe deployment tooling, and clear governance. Treat hosted CI as
validated only after a successful GitHub Actions run has been observed for the
current repository history.

### 2. Identity, authentication, trusted tenant context, and authorization

ADR-0006 accepts tenant-neutral user/session access tokens and an explicit
`/tenants/{tenant_id}/...` URL selector that remains untrusted until the server
validates the active user, exact active session, active non-suspended tenant,
and active membership. AUTH-002 / PR #25 now creates and propagates the typed
trusted Principal, composes the existing token and session primitives into the
protected roster route, and retains current role, object, and lifecycle
authorization at each data operation. The remaining identity work is
authentication entrypoints for initial session creation and HTTP login or
refresh flows.

### 3. SIAKAD integration contract

Define the provider-neutral boundary for inbound academic reference data and
future outbound results. The contract must preserve SIAKAD ownership of
programs, terms, courses, offerings, staff assignments, and enrollment; make
synchronization idempotent; retain history; audit conflicts; and represent
errors without turning ordinary LMS operations into an alternate academic
master-data path.

### 4. Course delivery and content

CONTENT-008 establishes the first tenant-safe API vertical slice for modules,
lessons, materials, publication state, availability windows, and file
metadata. Application authorization, composite database constraints, and RLS
protect offering and tenant boundaries; publication changes are audited.
Binary storage, signed file access, malware scanning, richer TA permission
configuration, activities, and frontend workflows remain follow-up work. No
storage provider is selected here.

### 5. Assessment, immutable submissions, and gradebook

Build assignments, submission records with immutable versions, grade items,
grades, feedback, and the existing authorization hierarchy. Preserve the
Teaching Assistant draft-score limitation and authorized publication/locking
workflow. This capability may evolve in parallel with course content because
the current repository has no invariant requiring the course-content capability
to precede it;
both still require trusted authentication and a verified hosted CI baseline. AI is
never the grade authority.

### 6. LMS-authoritative attendance and progress

Build attendance sessions/records and distinct activity-completion and
progress workflows. Preserve `ATTENDANCE_SOURCE_OF_TRUTH=LMS`, tenant and
offering authorization, correction history, and retention. Core attendance
does not depend on completing a SIAKAD adapter. Only a later outbound
attendance synchronization path may depend on the SIAKAD integration contract.

### 7. Operational reliability and current-schema recovery

The local backup, restore, normalized state, and source-safety validation cover
migration `0006` and the current application tables, including `auth_sessions`.
Remaining work is proportionate operational validation, structured signals,
recovery procedures, and measured reliability checks for capabilities that
actually exist. Do not claim production observability or recovery readiness
before the relevant checks run.

### 8. External and integration adapters

Implement only approved adapters after their ownership, data minimization,
failure, retry, idempotency, and audit contracts are defined. Provider-specific
representations must remain outside the canonical LMS model. External facts
must be revalidated before implementation relies on them.

### 9. Delivery interface and frontend foundation

After the API authorization contract is stable, establish the smallest
maintainable interface for implemented LMS capabilities. The interface must
consume authorized contracts and represent safe failure states. No frontend
framework, visual scope, or complete application is selected or claimed by
this roadmap.

### 10. Authorized staff-side AI

Only after the role allowlist, data-processing boundary, trusted authorization,
and relevant LMS data paths are approved should one narrow staff-side
assistance capability be implemented. It must remain optional, tenant-safe,
non-authoritative, source-traceable where meaningful, and removable without
changing core LMS semantics. AI artifacts intended for ordinary LMS content
must use the normal review and publication workflow.

### 11. Permanent production infrastructure

Reconcile the current deployment implementation with the accepted bounded
validation decision, then separately decide the requirements for permanent
ingress, recovery, access control, rate limiting, cost, and operations. Do not
select a provider, region, SKU, storage system, database service, or capacity
target from historic planning text.

## Deferred decisions

- `DECISION_REQUIRED: AI_ROLE_ALLOWLIST`
- Institutional AI data-processing boundary: local versus permitted external
  processing, including minimization, retention, and tenant policy.
- Permanent production ingress after bounded validation.
- External infrastructure and provider revalidation before reliance.

Unresolved decisions are not implementation commitments or migration blockers
unless a later work item explicitly depends on them.
