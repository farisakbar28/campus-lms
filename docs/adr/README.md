# Architecture Decision Records

Architecture Decision Records (ADRs) capture durable architectural choices,
the context that led to them, the alternatives considered, and their
consequences. They are concise decision records, not implementation tickets or
general project notes.

## When to write an ADR

Write an ADR when a choice has meaningful consequences for system boundaries,
data ownership, security, operations, cost, or long-term maintenance. A normal
code change does not require an ADR. A changed architecture should normally
receive a new ADR or an explicit supersession update so that the decision
history remains understandable.

An accepted decision is not silently rewritten when circumstances change.
Create a new ADR or document an explicit relationship. A later ADR may
supersede only part of an earlier decision; unaffected portions remain active
and must be identified clearly.

## Status

Use one of:

- Proposed — under consideration and not yet accepted;
- Accepted — the current decision for the scope stated in the ADR; or
- Superseded by ADR-XXXX — replaced by a later decision, with any retained
  scope described explicitly.

Dates record when the decision was recorded. External account, provider,
quota, cost, and runtime observations must be labelled as historical context
or as assumptions requiring revalidation rather than presented as permanent
repository facts.

## Tracked ADRs

| ADR | Title | Status |
|---|---|---|
| [0001](0001-pilihan-stack.md) | API and persistence stack direction | Proposed |
| [0002](0002-multi-tenancy.md) | Multi-tenancy strategy | Accepted |
| [0002b](0002b-arsitektur-hemat-biaya.md) | Cost-constrained private-origin architecture | Superseded by 0002e for scoped ingress assumptions |
| [0002c](0002c-azure-conventions.md) | Azure naming, region, and tagging conventions | Accepted |
| [0002d](0002d-azure-production-region-sku.md) | Azure production region and VM SKU | Accepted |
| [0002e](0002e-zero-domain-quick-tunnel-validation.md) | Zero-domain Quick Tunnel bounded validation | Accepted |
| [0005](0005-api-healthcheck-probe.md) | API binary health-check probe | Accepted |

Only ADR files tracked in this directory are listed here. A future decision
may add a new ADR; a missing number is not a placeholder or a requirement to
create one.
