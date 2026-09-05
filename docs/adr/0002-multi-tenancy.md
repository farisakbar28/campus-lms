# ADR-0002: Multi-Tenancy Strategy

- Date: 2026-08-24
- Status: Accepted

## Context

campus-lms is a multi-tenant SaaS Learning Management System for higher
education institutions. It is not SIAKAD: SIAKAD remains authoritative for
institutional academic facts, while the LMS owns learning processes after
those facts are available.

A tenant represents one institution. Users are global platform identities.
Memberships connect users to tenants, membership_roles grant tenant-scoped
roles, and course_staff grants course-offering authority.

The accepted implementation uses one shared PostgreSQL database and shared
schema, explicit tenant_id values on tenant-scoped rows, PostgreSQL Row-Level
Security (RLS), tenant-aware composite foreign keys, and transaction-local
tenant context. The repository's migrations and tests are
IMPLEMENTATION_FACT evidence for these mechanisms; they are not proof of a
particular cloud deployment or production capacity.

These layers solve different problems:

- RLS controls row visibility for the database role;
- composite foreign keys prevent structurally invalid cross-tenant
  relationships; and
- application authorization checks active membership, tenant role, course
  staff authority, object access, and lifecycle.

RLS is not physical database isolation, composite foreign keys do not replace
RLS, and RLS does not replace application authorization.

## Alternatives considered

### A. Shared database, shared schema, tenant_id, and RLS

Tenants share physical database and schema objects. One migration stream
reduces schema-version drift and simplifies pooling, monitoring, onboarding,
backup, and platform operations. The trade-off is a shared blast radius:
unsafe policies, privileged access, migrations, resource contention, or
database incidents can affect several tenants. Per-tenant restoration requires
application-aware selection and validation of related rows.

### B. Shared database with one schema per tenant

Separate schemas add a namespace and permission boundary but still share the
database server. Every schema needs migration and validation, creating
search-path, onboarding, backup, and schema-version-drift concerns. It adds
orchestration without a verified requirement in the current scope.

### C. Database per tenant

Separate databases provide a stronger infrastructure and administrative
boundary and can localize some incidents. They also multiply provisioning,
credentials, migrations, pools, monitoring, backups, recovery, and reporting.
The application must reliably route every request to the correct database.

## Decision

Use Option A: shared database, shared schema, explicit tenant_id, and
PostgreSQL RLS.

Tenant isolation is defense in depth:

- every tenant-scoped row carries tenant_id;
- RLS restricts reads and writes for the application role;
- composite foreign keys reject cross-tenant references;
- trusted tenant context is set only within a transaction; and
- application authorization enforces membership, role, course, object, and
  lifecycle rules.

This decision does not claim physical tenant isolation, a measured tenant
limit, a measured throughput or latency threshold, or a current provider
entitlement.

## Consequences and trade-offs

Positive consequences:

- one schema and migration stream;
- a smaller connection, monitoring, and onboarding surface;
- data-provisioned tenant onboarding;
- platform operations over one deliberate data plane; and
- local testability of RLS, roles, transactions, and constraints.

Accepted trade-offs:

- a database or migration incident can affect multiple tenants;
- RLS or privileged-access mistakes can have a broad blast radius;
- noisy-neighbor thresholds are not measured;
- every repository, job, cache, and new table must propagate tenant scope;
- per-tenant restore is more complex than restoring a dedicated database; and
- residency or regulatory requirements may later require a different boundary.

RLS does not make an authorized tenant member privileged for every action.
Course-scoped and role-specific authorization remains application/domain
logic.

## Supersession and relationships

No later ADR supersedes the shared-schema tenant-isolation decision. ADRs
0002b, 0002c, 0002d, and 0002e concern infrastructure and ingress scopes and
do not weaken this tenant boundary.

## Verification assumptions

The shared-schema, explicit-tenant, RLS, composite-FK, and transaction-local
mechanisms are IMPLEMENTATION_FACT supported by repository code,
migrations, and tests. Current database provider, account, quota, cost,
capacity, and production-runtime claims require
REVALIDATE_EXTERNAL_STATE. Any future change to this architecture requires a
new ADR or an explicit supersession relationship.
