# ADR-0002d: Azure Production Region and VM SKU

- Date: 2026-08-29
- Status: Accepted

## Context

ADR-0002c recorded Southeast Asia as the original Azure region convention.
This ADR recorded a later region and VM-sizing decision for the Azure
deployment design.

### Historical context

At the time of this decision, the subscription policy and discovery results
were recorded as allowing East Asia and Korea Central, with
Standard_B2ats_v2 (2 vCPU, 1 GiB RAM) and Standard_B2als_v2 (2 vCPU, 4 GiB
RAM) considered candidates for an x86-64 image path. The recorded observations
also included an Azure for Students B2ats v2 allowance and a possible Neon
Singapore database placement.

These values are HISTORICAL_CONTEXT as of 2026-08-29. They do not establish
current quota, availability, entitlement, VM state, provider state, cost,
latency, or whole-host capacity.

Container-level memory observations were not sufficient to prove that a
complete host, including the operating system, Docker, Caddy, API, migrations,
swap, and deployment overlap, would be safe at 1 GiB.

## Alternatives considered

### A. East Asia with Standard_B2ats_v2

The 1 GiB candidate minimized the recorded compute footprint but required a
real whole-host memory and stability gate. It left limited headroom.

### B. East Asia with Standard_B2als_v2

The 4 GiB candidate provided more headroom for host and deployment work but
could consume more of the recorded allowance. It was not an automatic or
zero-cost fallback.

### C. Korea Central equivalent

The equivalent candidate provided a regional alternative if East Asia became
unusable, but required fresh latency, availability, quota, and cost checks.

## Decision

The recorded architecture selected East Asia and an initial
Standard_B2ats_v2 trial, subject to whole-host memory and stability
validation before relying on that shape. Standard_B2als_v2 was identified as a
possible same-region reconsideration only after measured failure and an
explicit decision.

This ADR does not authorize an automatic paid resize or any paid fallback.
A regional change or VM-size change requires fresh factual verification and an
explicit decision under the zero-incremental-cost constraint.

## Consequences and trade-offs

The smaller trial shape made capacity decisions measurement-led and retained a
known larger-memory alternative. The trade-offs were limited headroom,
possible swap pressure, provider and regional dependency, and the need to
measure whole-host behavior rather than infer it from container-only results.

The API and database were recorded in different cloud regions/providers in the
historical design, so actual network behavior would need measurement before
relying on geographic assumptions.

## Supersession and relationships

ADR-0002d supersedes ADR-0002c only for its original production-region and
VM-sizing assumptions. ADR-0002c's resource-group and tagging conventions
remain Accepted.

ADR-0002b's cost constraint narrows this ADR: a B2als_v2 resize is not
automatic and must stop for explicit review if the initial trial fails.
ADR-0002e changes bounded ingress assumptions, not this region/SKU rationale.

## Verification assumptions

Before any present use, revalidate Azure policy, region and SKU availability,
quota, entitlement, VM state, provider placement, cost, latency, memory,
stability, and deployment behavior:

REVALIDATE_EXTERNAL_STATE

The recorded values and dates remain HISTORICAL_CONTEXT only.
