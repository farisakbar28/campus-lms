# campus-lms Engineering and Product Roadmap

Status: active product and engineering planning direction.

This document describes product direction and engineering dependencies. It is
Operational work is tracked in GitHub Issues rather than in this roadmap. A
capability is not presented as implemented merely because it is described here.

## Purpose and planning rules

The roadmap is dependency-driven and implementation-state-aware. Each item is
classified as an existing foundation, partial work, a future capability, a
stale implementation that needs reconciliation, or a decision that remains
open. Sequencing follows authorization, data ownership, safety, and operational
dependencies rather than calendar milestones.

Current engineering work follows the repository contracts in `AGENTS.md`,
`CONTRIBUTING.md`, the scoped engineering instructions, `docs/domain.md`,
`docs/domain-ai.md`, and accepted ADRs. Applied migrations remain immutable;
security-sensitive configuration fails closed; and no automatic paid fallback
may be introduced.

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

Backup and restore validation remains pinned to migration 5, the pre-0006
schema, and 11 application tables. It is not current-schema validation.

Historic external infrastructure, provider, region, quota, storage, and
capacity facts require fresh revalidation before reliance. The project retains
the zero-incremental-paid-infrastructure constraint and must fail closed,
defer, or disable a capability when verified zero-cost capacity is exhausted.

## Current implementation snapshot

| Classification | Current repository state |
|---|---|
| Implemented foundation | Go API foundation; fail-closed configuration; `/healthz` and dependency-backed `/readyz`; migrations through `0006`; tenant/RLS and composite-constraint foundations; roster repository and handler; access-token and refresh-session primitives; Docker, Compose, health probes, and deployment tooling foundations. |
| Partial / requires composition | Authentication: primitives exist, but the running server does not yet produce the trusted `Principal` required by protected routes. Backup/restore: local tooling exists, but its validator is stale and does not cover the current schema. |
| Not implemented | Frontend application; AI subsystem; AI evaluation; complete course-content workflow; assessment and gradebook workflow; attendance workflow; production observability; permanent production ingress; and a current hosted CI baseline. |
| Stale / requires reconciliation | Legacy deployment wiring using loopback 8443, Origin-CA files, and hostname-based Caddy TLS conflicts with the later bounded Quick Tunnel validation ADR. |
| Decision required | AI role allowlist; institutional AI data-processing boundary; permanent production ingress; and fresh external infrastructure/provider validation. |

These classifications are bounded by the current source and documentation.
Existing authentication and session primitives do not imply complete
authenticated API behavior, production readiness, or a complete authorization
boundary.

## Dependency-oriented capability roadmap

### 1. Engineering baseline

Preserve the runnable Go/PostgreSQL repository, explicit migrations, focused
tests, safe deployment tooling, and clear governance. Treat hosted CI as
validated only after a successful GitHub Actions run has been observed for the
current repository history.

### 2. Identity, authentication, trusted tenant context, and authorization

First decide how an authenticated user and session select or represent tenant
context. The decision must require server-side validation against identity,
active membership, role and lifecycle rules, route/object authorization, and
fail-closed behavior. Then compose the existing token and session primitives
into the running server and protected routes. Tenant context must never gain
authority merely from an unvalidated client selector.

### 3. SIAKAD integration contract

Define the provider-neutral boundary for inbound academic reference data and
future outbound results. The contract must preserve SIAKAD ownership of
programs, terms, courses, offerings, staff assignments, and enrollment; make
synchronization idempotent; retain history; audit conflicts; and represent
errors without turning ordinary LMS operations into an alternate academic
master-data path.

### 4. Course delivery and content

Build tenant-safe course-offering content workflows for modules, lessons,
materials, publication state, and file metadata. Application authorization and
database constraints must protect offering and tenant boundaries. Binary
storage remains a separate decision when it is actually introduced; no storage
provider is selected here.

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

Reconcile backup, restore, and safety validation with migration `0006` and the
current application tables. Add proportionate operational validation,
structured signals, recovery procedures, and measured reliability checks for
capabilities that actually exist. Do not claim production observability or
recovery readiness before the relevant checks run.

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
