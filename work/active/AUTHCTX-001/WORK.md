# AUTHCTX-001 — Trusted authenticated tenant-context semantics

Contract version: `campus-lms-work/v1`

Status: `IMPLEMENTING`

Issue: `AUTHCTX-001` / GitHub issue `#1`

Issue URL: <https://github.com/farisakbar28/campus-lms/issues/1>

Issue specification digest:
`sha256:ff8cdd44535192565da4904b9b2b0eda60a9454bb6589777b690c5f4cba7137f`

Issue `updatedAt` observed: `2026-09-10T12:48:31Z` (informational only)

Plan author actor label: `Codex planner`

Plan author session label: `AUTHCTX-001-plan-r1-2026-09-10-codex-root`

Plan revision: `1`

Plan hash: `sha256:65f6e0bd70490a0a60e9ebf6a17560cc95d13e35331d961695d96aa39ec9dab0`

Risk: `HIGH`

Release required: `false`

Human plan-approval comment: `https://github.com/farisakbar28/campus-lms/issues/1#issuecomment-5618927606`

Implementation author actor label: `Codex implementer`

Implementation author session label: `AUTHCTX-001-implementation-r1-2026-09-10-codex-root`

Candidate Git commit SHA: `NONE`

This file is the approved execution binding for a documentation-only
implementation and local candidate commit. It authorizes no push, Pull Request
mutation, merge, external mutation, release, or deployment. A fresh independent
implementation review is required for the exact candidate commit.

<!-- PLAN-NORMATIVE-BEGIN -->

## Identity and input checkpoint

- Contract: `campus-lms-work/v1`.
- Work item: `AUTHCTX-001`.
- Canonical issue: GitHub issue `#1`,
  <https://github.com/farisakbar28/campus-lms/issues/1>.
- Issue specification digest algorithm: `campus-lms-issue-sha256-v1`.
- Issue specification digest:
  `sha256:ff8cdd44535192565da4904b9b2b0eda60a9454bb6589777b690c5f4cba7137f`.
- Base Git revision:
  `2a792695cd73903a30d9f94e58694b2e559d3980`.
- Plan revision: `1`.
- Plan hash algorithm: `campus-lms-plan-sha256-v1`.
- Risk: `HIGH`, because this decision defines the authentication-to-tenancy
  trust boundary; an incorrect choice could permit unauthorized or
  cross-tenant access when composed later.
- Release required: `false`; this item records an architecture contract only
  and creates no deployable runtime, schema, or configuration change.

Planning is bound to the current canonical Issue title/body, current source,
`docs/domain.md`, accepted ADR-0002, `docs/roadmap.md`, root `AGENTS.md`, and
`apps/api/AGENTS.md`. Current implementation facts are limited to the observed
repository base: access-token verification yields tenant-neutral user/session
identity; bearer middleware propagates only that identity; refresh-session
storage and lifecycle primitives exist; the running server does not compose
authentication; `Principal` currently contains string tenant/user IDs and is
created only by tests; and the roster repository performs tenant, active
membership, active lecturer-role, and active course-staff checks inside an
RLS-scoped transaction. This plan does not treat those partial primitives as a
complete authentication or authorization boundary.

## Goal

Record one accepted architecture decision that makes tenant selection
explicit without letting client input confer authority, defines server-side
identity/session/membership validation and fail-closed route semantics, and
gives dependent implementation issue `AUTH-002` an exact contract for creating
and consuming a trusted `Principal`.

## Decision to record

The implementation of this item will add ADR-0006 with the following complete
decision.

### Authentication identity and session

1. Short-lived access tokens remain tenant-neutral. Their authorization-bearing
   identity is limited to the global user ID and global authentication-session
   ID plus required protocol claims. Tenant IDs, memberships, tenant roles,
   course roles, and object permissions are not embedded as token claims.
2. Every protected request first passes the existing strict Bearer boundary.
   The server then resolves the verified `(user_id, session_id)` against global
   database state. Authentication succeeds only when the exact session belongs
   to that user, is unrevoked, and is unexpired, and the user has the explicit
   active status. Missing rows, unknown status values, mismatches, revocation,
   and expiry fail closed.
3. Refresh-token rotation makes the predecessor session ineligible for later
   protected requests; newly issued access tokens bind to the rotation child.
   An already admitted in-flight request uses its validated request snapshot,
   while the next request observes revocation or lifecycle changes.

### Tenant selection and validation

1. Tenant-scoped HTTP routes use an immutable tenant UUID in an explicit URL
   path segment, with the canonical shape `/tenants/{tenant_id}/...`. The path
   value is an untrusted selector and RLS partition-scoping input only. It never
   becomes authorization by itself.
2. The access token remains reusable across tenants in which the user is
   authorized. Each request independently selects and validates one tenant.
   No mutable "current tenant" is stored in the global auth session.
3. Request headers, query parameters, bodies, cookies, token claims, and
   client-supplied user/session values are not tenant-authorization inputs.
   An object ID alone is not used to discover tenant context through a global
   or RLS-bypassing lookup.
4. After Bearer verification, the server may use the syntactically valid path
   UUID to establish transaction-local `app.tenant_id` solely to constrain the
   membership-resolution query. It creates a trusted Principal only if that
   same request snapshot confirms an active user, the exact active auth
   session, a tenant with `status = 'active'` and no suspension timestamp, and
   the unique membership for `(tenant_id, user_id)` with `status = 'active'`.
   No tenant-scoped application data is returned or mutated before this
   validation succeeds.
5. Unknown or future user, tenant, and membership statuses are denied until a
   later approved decision explicitly permits them. Selection of another
   active membership requires another request using that tenant's URL.

### Principal and downstream authorization

1. Trusted middleware constructs the Principal from verifier and repository
   output only. The Principal contains typed, non-zero `UserID`, `SessionID`,
   `TenantID`, and `MembershipID` UUIDs. Handlers never construct or amend it
   from request fields.
2. Tenant roles and course-offering roles are deliberately not durable access
   token claims and are not treated as timeless Principal authority. At the
   point of each action, application/repository authorization rechecks the
   required unrevoked tenant-role grant, active course-staff assignment or
   active enrollment, object ownership, and action-specific lifecycle state.
3. Each tenant-scoped data operation performs those action/object checks in
   the same transaction that sets `SET LOCAL app.tenant_id` and reads or writes
   the protected data. It uses the trusted Principal IDs. RLS and composite
   tenant constraints remain defense in depth and never replace application
   authorization.
4. Holding an active membership or a broad tenant role is never sufficient for
   course-object access. The existing roster contract continues to require its
   authorized tenant role and offering assignment; this decision does not
   expand it to Teaching Assistants, students, tenant administrators, or other
   actors.
5. Public liveness/readiness endpoints remain outside tenant authentication.
   Global authentication endpoints remain tenant-neutral. Platform-operator
   and audited break-glass routes require a separate explicit Principal and
   authorization design and must not reuse a tenant membership as platform
   authority.

### Route and failure semantics

1. Missing, malformed, unverifiable, or expired Bearer credentials and invalid,
   revoked, expired, or mismatched user/session state return the existing
   generic `401` authentication response and Bearer challenge without exposing
   tokens or internal state.
2. A malformed tenant UUID is rejected before repository access with a stable
   safe `400` response. A syntactically valid but missing, suspended, inactive,
   or unauthorized tenant/membership selection uses one indistinguishable
   safe `404` response; it must not reveal whether the tenant or membership
   exists.
3. After tenant admission, `403` is used only for a tenant-level capability
   denial that discloses no additional object existence. A missing object,
   cross-tenant object, or object for which the Principal lacks access uses the
   same safe `404` response to avoid an object-enumeration oracle. Invalid
   object syntax remains `400`.
4. Dependency unavailability maps to `503`; unexpected internal failures map
   to `500`. Logs are structured and may contain bounded internal error and
   request-correlation information, but never credentials, request bodies,
   secrets, or student personal data.
5. Middleware order for tenant routes is fixed as Bearer verification,
   active user/session and selected-tenant membership resolution, Principal
   propagation, then route-specific role/object/lifecycle authorization and
   the data operation. Any missing stage fails closed.

## Alternatives and rationale

ADR-0006 will document these rejected alternatives:

1. A tenant claim in each access token is rejected because it makes tenant
   switching require token reissuance and leaves membership/role/lifecycle
   state stale unless the database is checked anyway.
2. A mutable tenant stored in the global auth session is rejected because it
   creates hidden cross-request state, multi-tab races, and surprising tenant
   switching.
3. A tenant header, query field, body field, or cookie is rejected as the
   canonical selector because tenant scope should be explicit in tenant
   resource URLs and must not become implicit proxy/cache/client state.
4. Deriving tenant solely from an object ID is rejected because object IDs are
   not authorization boundaries, many routes have no existing object, and a
   pre-RLS global lookup would weaken the accepted tenant architecture.
5. A separate access token per tenant is deferred because it adds issuance,
   refresh, switching, and revocation complexity without removing the need for
   current server-side authorization checks.

The chosen URL selector is visible, stable, compatible with tenant-local RLS,
and supports users with several memberships. Its cost is wider route changes
and repeated server-side validation. That cost is accepted because explicit
scope and current authorization state are more important than shorter URLs or
fewer lookups.

## Scope and affected areas

Implementation is documentation-only and limited to:

- add `docs/adr/0006-trusted-authenticated-tenant-context.md` with status
  `Accepted`, the decision above, alternatives, consequences, fail-closed
  semantics, and an explicit relationship to ADR-0002 and `AUTH-002`;
- add ADR-0006 to `docs/adr/README.md`; and
- narrowly update `docs/roadmap.md` so it identifies the accepted tenant-context
  decision while continuing to state that runtime Principal/authentication
  composition is not implemented.

No Go source, tests, configuration, dependency metadata, migration, database
state, API runtime behavior, deployment file, or external system is changed by
this work item. The later `AUTH-002` plan must independently bind its Issue and
translate this decision into code and focused composition/integration tests.

## Exclusions

- identity-provider selection, login/user-provisioning protocol, frontend, and
  user-interface tenant switching;
- implementation of middleware, repositories, handlers, routes, auth
  endpoints, or server startup wiring;
- new token/session primitives, schema changes, status migrations, RLS changes,
  or edits to applied migrations;
- expansion of roster permissions or definition of every future route's
  role/action matrix;
- platform-admin or break-glass implementation;
- caching authorization state, rate limiting, CSRF policy, audit-event schema,
  or observability design beyond the safe logging boundary;
- push, Pull Request mutation, merge, release, deployment, cloud mutation, or
  destructive data work; and
- any incremental paid infrastructure or tooling.

## Implementation plan

1. Re-fetch the canonical Issue and recompute its digest before implementation;
   stop and return to `DRAFT` if the normalized title/body differs.
2. Re-read the current domain contract, ADR-0002, roadmap, scoped engineering
   instructions, access-token/session/middleware/Principal code, tenant
   transaction boundary, roster authorization query, and relevant migrations.
   Treat any conflict affecting tenancy, security, authority, scope, or
   lifecycle semantics as material.
3. Add ADR-0006 using the repository ADR structure. State the selected URL-path
   mechanism, trust transitions, exact Principal contents, server-side
   lifecycle checks, transaction/RLS ordering, downstream authorization
   duties, safe HTTP semantics, rejected alternatives, consequences, and
   boundaries without claiming runtime implementation.
4. Update the ADR index and the minimum roadmap text needed to distinguish an
   accepted decision from still-missing composition. Do not rewrite the domain
   model or unrelated roadmap/ADR history.
5. Inspect the complete diff for consistency, accidental behavior claims,
   secret/personal-data exposure, and changes outside the documentation-only
   scope. Run the specified checks and record only their concise result in the
   mutable field below.
6. If implementation is later explicitly authorized, create a local candidate
   commit only if that separate instruction authorizes it, then stop for a fresh
   independent implementation review. This plan authorizes no commit by itself.

## Acceptance criteria

1. ADR-0006 is indexed and records one unambiguous trusted-selection mechanism:
   tenant-neutral `(user_id, session_id)` Bearer identity plus a tenant UUID URL
   selector that is untrusted until server-side validation succeeds.
2. The ADR requires exact active user/session ownership and lifecycle checks,
   active/non-suspended tenant state, and active membership before typed
   Principal creation; missing, mismatched, unknown, revoked, expired,
   suspended, or inactive state fails closed.
3. The ADR defines Principal propagation using non-zero typed user, session,
   tenant, and membership identities and forbids client fields from populating
   that Principal.
4. The ADR preserves the full authorization hierarchy: current tenant-role,
   course assignment/enrollment, object tenant, action permission, and
   lifecycle checks occur at the data operation, with transaction-local RLS and
   composite constraints retained as defense in depth.
5. The ADR defines route scoping and non-enumerating `400`/`401`/`403`/`404`,
   dependency `503`, and safe `500` behavior precisely enough for `AUTH-002` to
   test without inventing a trust rule.
6. Alternatives and trade-offs cover token claims, session-selected tenant,
   other client transports, object-derived tenant, and per-tenant tokens.
7. ADR-0006 explicitly preserves accepted ADR-0002 and identifies `AUTH-002` as
   the separate runtime-composition work; the roadmap makes no false claim that
   protected routes already receive a trusted Principal.
8. Only the new ADR, ADR index, narrow roadmap update, and active workflow pair
   change. No migration, application behavior, external state, dependency, or
   paid service is introduced.

## Verification plan

- Re-fetch GitHub issue #1 and recompute its canonical title/body digest before
  implementation and at later workflow gates.
- Run `make workflow-test`.
- Run `make workflow-check`.
- Run `git diff --check`.
- Inspect `git diff -- docs/adr/0006-trusted-authenticated-tenant-context.md
  docs/adr/README.md docs/roadmap.md work/active/AUTHCTX-001/WORK.md
  work/active/AUTHCTX-001/REVIEWS.md` and `git status --short` to confirm the
  documentation-only path boundary and absence of secrets or personal data.
- Verify the ADR index link resolves, its status/date/relationships match the
  ADR, every Issue acceptance topic is addressed, and the roadmap still says
  runtime authentication composition is incomplete.
- Go, database, and integration tests are not required for this documentation-
  only decision. If implementation changes any runtime or schema file, that is
  a material deviation requiring a revised plan, focused tests, fresh plan
  review, and new human approval.

## Documentation impact, risks, deviations, and release choice

The durable documentation impact is one accepted ADR, its index entry, and a
narrow roadmap reconciliation. `docs/domain.md` remains the domain source of
truth and ADR-0002 remains accepted without modification.

Primary risks are mistaking a client selector for authority, omitting session
revocation from access-token admission, allowing status values by default,
using RLS as the sole authorization layer, leaking tenant/object existence,
and falsely describing the decision as implemented. The specified trust
transition, deny-by-default lifecycle rules, same-transaction action checks,
safe response equivalence, and roadmap wording mitigate those risks. The
remaining accepted trade-offs are an explicit tenant prefix on protected URLs,
per-request database validation, and a future route migration in `AUTH-002`.

Material deviations: `NONE`. Any change to the chosen selector transport,
Principal fields, validation or error semantics, authorization/RLS layering,
affected files, acceptance, risk, dependencies, or release choice is material
and requires `DRAFT`, a new revision/hash, fresh independent plan review, and
new human approval.

Release required remains `false` because the item is documentation-only and
does not create a deployable artifact or authorize production work.

<!-- PLAN-NORMATIVE-END -->

- Material deviations: `NONE`
- Concise verification result: `PASS — Issue digest and approval revalidated; workflow-test (33 tests), workflow-check, diff checks, and scoped documentation inspection passed`
- Concise completion summary: `Added accepted ADR-0006, indexed it, and reconciled the roadmap; runtime and schema remain unchanged`
