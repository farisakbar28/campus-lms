# ADR-0001: API and Persistence Stack Direction

- Date: 2026-08-17
- Status: Proposed

## Context

campus-lms is a multi-tenant higher-education LMS. The stack must support
tenant isolation, versioned academic records, explicit SQL migrations,
auditable workflows, bounded resource use, and a maintainable deployment
surface. The project also has a zero-incremental-cost constraint.

This ADR is a proposal, not a retrospective acceptance of every technology
mentioned in an earlier plan. Current repository implementation facts and
future choices are intentionally separated.

## Current repository implementation facts

The repository currently contains:

- a Go API under apps/api;
- standard-library HTTP and Go package boundaries;
- PostgreSQL migrations and explicit SQL access;
- Docker and Docker Compose definitions; and
- a modular package structure separating configuration, HTTP, domain,
  middleware, database, and repository concerns.

These facts describe what is present in the repository. They do not prove a
particular production cloud deployment, provider entitlement, capacity
threshold, or future AI architecture.

## Alternatives considered

### API language

- Go — a small compiled runtime, standard-library HTTP support, explicit
  concurrency and cancellation, and a good fit for the current API code.
- Node.js or TypeScript — broad web tooling, with a different runtime and
  dependency-management profile.
- Python — strong general-purpose data and AI tooling, with a different
  runtime and dependency profile.

### Persistence

- PostgreSQL — relational constraints, transactions, RLS, explicit migrations,
  and SQL that matches the domain's integrity requirements.
- A managed PostgreSQL service — potentially less operational work, but with
  provider lifecycle and account assumptions that require separate decisions.

### Application shape

- Modular monolith — one API deployment with explicit internal module
  boundaries.
- Multiple services — independent deployment boundaries at the cost of
  additional contracts, operations, and failure modes.

## Decision

Propose retaining the repository's Go API, PostgreSQL persistence model, and
modular package boundaries as the current implementation direction, subject to
formal acceptance of this ADR. The proposal does not select a frontend,
cloud/database provider, AI provider, AI model, or separate AI-service
architecture. Those choices remain future decisions and must be documented
separately if they acquire an approved scope.

## Consequences and trade-offs

The current direction keeps the API resource-conscious and makes migrations,
constraints, transactions, RLS, and repository behavior explicit. A modular
monolith limits network and deployment overhead while preserving boundaries
that can be tested in one repository.

The trade-offs are that module boundaries require discipline, PostgreSQL
operations remain an explicit responsibility, and a later service split would
require deliberate contracts and ownership changes. The proposal does not
claim a measured capacity limit or a production provider entitlement.

## Supersession and relationships

No later ADR supersedes this proposal. The accepted multi-tenancy and
health-check decisions constrain relevant implementation details. Future AI,
provider, and infrastructure choices require their own decision records.

## Verification assumptions

The Go API, migrations, Docker definitions, and package boundaries are
IMPLEMENTATION_FACT supported by the repository. Any current cloud,
provider, quota, cost, or capacity statement requires
REVALIDATE_EXTERNAL_STATE. Unimplemented technologies must not be presented
as committed architecture.
