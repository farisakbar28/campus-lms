# ADR-0002b: Cost-Constrained Private-Origin Architecture

- Date: 2026-08-30
- Status: Superseded by ADR-0002e (scoped: ingress, domain, Origin CA, and
  normal administration only)

## Context

This ADR records the cost and private-origin architecture considered before
the bounded no-domain validation route in ADR-0002e. It remains the historical
record for decisions not replaced by that ADR.

The enduring constraints were a zero-incremental-cost policy, a private
application origin, explicit separation between host infrastructure and
application Compose, and a recovery boundary outside the application process.
Those constraints remain relevant unless a later ADR changes them.

### Historical context

At the time this decision was recorded, the proposed resource and service
context included an Azure VM trial, a Cloudflare Tunnel, Caddy, a PostgreSQL
provider, and an off-machine logical-backup candidate. Specific region, SKU,
disk, domain, provider, allowance, and account observations were
HISTORICAL_CONTEXT as of 2026-08-30. They are not current runtime or entitlement
claims.

## Alternatives considered

### A. Public Azure origin

A public IP and direct Caddy ingress would provide conventional routing and
administration but would expose the origin and conflict with the
zero-incremental-cost/private-origin constraints.

### B. Private origin with a managed Tunnel

Cloudflare edge would reach a host-level connector, a loopback Caddy origin,
and the Docker API service. This avoided public origin ingress but depended on
domain, connector, credential, and administration arrangements.

### C. Private origin with deterministic paid egress

Explicit paid egress could make outbound behavior more predictable, but was
not eligible under the zero-incremental-cost constraint.

### D. Recovery choices

Native database recovery was simpler but might not satisfy a desired
off-machine retention objective. A separate logical-backup destination and
restore drill added operational work but could preserve an eligible
zero-incremental-cost recovery candidate.

## Decision

The original decision selected a private origin, host-managed connector
infrastructure, no public Azure application ingress, a zero-incremental-cost
policy, and a separately designed recovery boundary. Host infrastructure was
kept outside the application Compose lifecycle, and the application origin
was intended to remain reachable only through the approved edge path.

ADR-0002e supersedes the original named-Tunnel, custom-domain, Origin CA, and
normal-administration assumptions. It does not erase the private-origin,
zero-incremental-cost, recovery-boundary, or other unaffected decisions
recorded here.

There is no automatic paid resize, paid model, paid egress, or paid domain
renewal fallback implied by this ADR.

## Consequences and trade-offs

The private-origin posture avoids direct Internet exposure of the application
origin and keeps host infrastructure separate from application lifecycle.
The cost constraint limits available capacity and provider choices, and the
recovery boundary requires a tested backup and restore procedure.

The trade-offs include dependence on an edge and connector path, limited
availability while the host is unavailable, operational complexity around
secret bootstrap, and the need to verify external service limits before use.
Those are design consequences, not evidence that any runtime path has been
activated.

## Supersession and relationships

ADR-0002e is the scoped successor for:

- ingress topology;
- domain requirements;
- Origin CA assumptions; and
- normal administration.

ADR-0002e does not provide permanent production ingress. ADR-0002d remains
the historical record of its own region and VM-sizing decision, subject to
the cost and fallback constraints documented here.

## Verification assumptions

Any present Azure, Cloudflare, database-provider, domain, quota, free-tier,
cost, account, connector, backup, or runtime assertion requires
REVALIDATE_EXTERNAL_STATE. Historical values may be retained only as
HISTORICAL_CONTEXT with their decision date. The repository does not claim
that the Quick Tunnel, backup destination, or recovery drill has run.
