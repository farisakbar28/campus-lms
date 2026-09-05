# ADR-0002c: Azure Naming, Region, and Tagging Conventions

- Date: 2026-08-16
- Status: Accepted

## Context

Cloud resources need predictable names, grouping, and ownership metadata so
that operations, cost review, and cleanup remain understandable. The
convention is intentionally small:

- one production resource group;
- a project and environment marker on supported resources; and
- an owner marker on supported resources.

The original region and subscription observations below are historical
context. They do not assert current account, quota, availability, or cost
state.

### Historical context

On 2026-08-16, Southeast Asia (Singapore) was selected as the expected region
for Indonesian users, with a resource group named rg-campuslms-prod and the
tags project=campus-lms, env=prod, and owner=faris. This dated selection
explains the convention but is not a current Azure availability claim.

## Alternatives considered

### Region convention

The original alternatives were Southeast Asia for expected proximity to
Indonesia, East Asia as another Asia-Pacific option, and Australia East as a
further fallback. Region selection remains subject to subscription policy,
availability, quota, latency, and cost checks.

### Resource grouping

- One production group, rg-campuslms-prod — simple management and cleanup.
- Per-environment groups — clearer separation if additional environments are
  introduced, with more administrative overhead.

### Tagging

- Mandatory project, environment, and owner tags — readable ownership and cost
  reporting.
- No tags — less discipline, but poor accountability and cost attribution.

## Decision

Use predictable resource naming, group production resources under
rg-campuslms-prod, and apply project=campus-lms, env=prod, and owner=faris to
resources that support tags. Additional environments should use an explicit
environment suffix and separate group when they become an approved
requirement.

The naming and tagging decision does not select a current region, VM SKU,
provider entitlement, or deployment state.

## Consequences and trade-offs

The convention makes ownership, environment, cost review, and cleanup easier
without requiring a separate naming service. A single production group does
not isolate a future staging environment, and manual tagging can drift without
policy enforcement.

## Supersession and relationships

ADR-0002d supersedes only this ADR's original production-region and VM-sizing
assumptions. The resource-group and tagging conventions remain Accepted.
ADR-0002e changes bounded ingress assumptions but does not replace this
resource organization decision.

## Verification assumptions

Current Azure subscription policy, region availability, quota, VM SKU,
entitlements, resource state, and costs require
REVALIDATE_EXTERNAL_STATE. The dated Singapore selection is
HISTORICAL_CONTEXT, not proof of current availability.
