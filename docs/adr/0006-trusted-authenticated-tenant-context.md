# ADR-0006: Trusted Authenticated Tenant Context

- Date: 2026-09-10
- Status: Accepted

## Context

The API has tenant-neutral access-token verification, bearer parsing, and
global refresh-session lifecycle primitives. It does not yet compose those
pieces into the running server or create the trusted `Principal` required by
protected routes. The existing `Principal` is string-based and populated only
by tests.

ADR-0002 requires tenant-scoped database work to combine application
authorization, transaction-local PostgreSQL RLS context, and composite tenant
constraints. A client therefore needs a way to select a tenant without that
selector becoming authority. Users may have active memberships in several
tenants, while roles, course assignments, enrollments, object ownership, and
lifecycle state may change independently of an access token.

This ADR defines the authentication-to-tenancy trust transition and the
contract that later runtime composition must implement. It does not implement
middleware, repositories, routes, or authentication endpoints.

## Decision

### Authentication identity and session

Short-lived access tokens remain tenant-neutral. Their authorization-bearing
identity is limited to a global user UUID and global authentication-session
UUID, plus required protocol claims. Tenant IDs, memberships, tenant roles,
course roles, and object permissions are not token claims.

Every protected request first passes the strict Bearer boundary. The server
then resolves the verified `(user_id, session_id)` against global database
state. Authentication succeeds only when all of these conditions hold in the
request's validation snapshot:

- the user exists and has exactly `status = 'active'`;
- the exact session exists and belongs to that user;
- the session has not been revoked; and
- the session has not expired.

Missing rows, an ownership mismatch, revocation, expiry, and unknown or future
user statuses fail closed.

Refresh-token rotation makes the predecessor session ineligible for later
protected requests. A newly issued access token binds to the rotation child.
An already admitted in-flight request uses its validated request snapshot; the
next request must observe revocation or other lifecycle changes.

### Tenant selection and admission

Tenant-scoped HTTP routes use an immutable tenant UUID in an explicit path
segment with the canonical shape:

```text
/tenants/{tenant_id}/...
```

The path UUID is an untrusted selector and RLS partition-scoping input only.
It never confers authority. The tenant-neutral access token can be reused for
tenants in which the user is authorized, but every request independently
selects and validates exactly one tenant. The global authentication session
does not store mutable current-tenant state.

Headers, query parameters, request bodies, cookies, token claims, and
client-supplied user or session values are not tenant-authorization inputs. An
object ID alone must not be used to discover tenant context through a global
or RLS-bypassing lookup.

After Bearer verification, the server may use a syntactically valid path UUID
to set transaction-local `app.tenant_id` solely to constrain the
membership-resolution query. That selector becomes trusted tenant context only
after the same request snapshot confirms all of these conditions:

- the active user and exact active auth session checks above pass;
- the selected tenant exists, has exactly `status = 'active'`, and has no
  suspension timestamp; and
- the unique membership for `(tenant_id, user_id)` exists and has exactly
  `status = 'active'`.

No tenant-scoped application data may be returned or mutated before admission
succeeds. Missing data, suspension, inactive state, and unknown or future
tenant or membership statuses fail closed. Selecting another membership
requires another request with that tenant's URL.

### Principal and downstream authorization

Only trusted middleware constructs and propagates a `Principal`, using
verifier and repository output rather than request fields. It contains typed,
non-zero UUID values for `UserID`, `SessionID`, `TenantID`, and `MembershipID`.
Handlers must not construct or amend it from client input.

The `Principal` proves only the validated user, session, tenant, and membership
snapshot. Tenant roles and course-offering roles are neither durable access
token claims nor timeless Principal authority. At each action, application and
repository authorization must recheck the applicable:

- unrevoked tenant-role grant;
- active course-staff assignment or active enrollment;
- object ownership and tenant relationship;
- action permission; and
- lifecycle state.

Each tenant-scoped data operation performs those action and object checks in
the same transaction that sets transaction-local `app.tenant_id` and reads or
writes the protected data. It uses the trusted Principal IDs. RLS and
composite tenant constraints remain defense in depth and do not replace
application authorization.

An active membership or broad tenant role is never sufficient for
course-object access. In particular, the existing roster contract continues
to require an unrevoked `lecturer` tenant-role grant and an active
`instructor` or `lead_instructor` assignment for the offering. This ADR does
not extend roster access to Teaching Assistants, students, tenant
administrators, or other actors.

Public liveness and readiness endpoints remain outside tenant authentication.
Global authentication endpoints remain tenant-neutral. Platform-operator and
audited break-glass routes require a separate explicit Principal and
authorization design; tenant membership must not be reused as platform
authority.

### Route and failure semantics

Tenant routes apply middleware in this fixed order:

1. Bearer verification;
2. active user/session and selected-tenant membership resolution;
3. trusted Principal propagation; and
4. route-specific role, object, and lifecycle authorization together with the
   data operation.

Any missing or failed stage fails closed. HTTP boundaries use these safe
semantics:

- missing, malformed, unverifiable, or expired Bearer credentials, and
  invalid, revoked, expired, or mismatched user/session state, return the
  existing generic `401` response and Bearer challenge;
- a malformed tenant UUID is rejected before repository access with a stable,
  safe `400` response;
- a syntactically valid but missing, suspended, inactive, or unauthorized
  tenant or membership selection returns one indistinguishable safe `404`
  response;
- after tenant admission, `403` is used only for a tenant-level capability
  denial that reveals no additional object existence;
- a missing object, cross-tenant object, or object the Principal cannot access
  returns the same safe `404`, while invalid object syntax remains `400`;
- dependency unavailability maps to `503`; and
- unexpected internal failures map to `500`.

Responses must not reveal tokens or internal state. Logs remain structured and
may contain bounded internal errors and request-correlation information, but
never credentials, request bodies, secrets, or student personal data.

## Alternatives considered

1. Put a tenant claim in each access token. Rejected because switching tenants
   would require token reissuance, while current membership, role, and
   lifecycle state would still require database validation.
2. Store a mutable current tenant in the global auth session. Rejected because
   it creates hidden cross-request state, multi-tab races, and surprising
   tenant switching.
3. Use a header, query field, request body field, or cookie as the canonical
   selector. Rejected because tenant scope should be explicit in tenant
   resource URLs and must not become implicit proxy, cache, or client state.
4. Derive tenant context only from an object ID. Rejected because object IDs
   are not authorization boundaries, many routes do not start from an existing
   object, and a pre-RLS global lookup would weaken the accepted tenant
   architecture.
5. Issue a separate access token per tenant. Deferred because it adds
   issuance, refresh, switching, and revocation complexity without removing
   the need for current server-side authorization checks.

## Consequences and trade-offs

The selected path makes tenant scope visible and stable, supports users with
several memberships, and aligns membership lookup and later data operations
with tenant-local RLS. It also establishes one testable, fail-closed contract
for `AUTH-002` without treating any client value as authority.

The accepted costs are a tenant prefix on protected URLs, per-request database
validation, and future route and Principal migration work. Authorization state
cannot be inferred solely from a previously issued token or a previously
validated request. Implementations must preserve non-enumerating failure
equivalence and keep authorization checks next to protected data operations.

## Supersession and relationships

This decision preserves ADR-0002. ADR-0002 continues to govern shared-schema
tenant isolation, transaction-local RLS, composite tenant constraints, and
application authorization. This ADR specifies how an HTTP request earns the
trusted Principal whose tenant identity can enter that boundary.

`AUTH-002` is the separate runtime-composition work item. It must independently
bind its Issue and translate this decision into middleware, repositories,
route composition, typed Principal changes, and focused integration tests.
This ADR does not claim that the running server currently performs those
steps.

## Verification assumptions

Tenant-neutral access-token verification, Bearer identity propagation, global
refresh-session lifecycle primitives, tenant transactions, and the current
roster query are IMPLEMENTATION_FACT supported by repository code and tests.
The full user/session/membership admission boundary, typed four-ID Principal,
tenant-prefixed protected routes, and production authentication composition
are not yet implemented and must not be inferred from this accepted decision.
