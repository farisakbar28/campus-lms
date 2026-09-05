# ADR-0002e: Zero-Domain Quick Tunnel Bounded Validation

- Date: 2026-08-31
- Status: Accepted

## Context

The project needed a temporary public-HTTPS validation path without requiring
a custom domain, custom DNS mutation, or an authorized paid network service.
This ADR defines a bounded validation route only. It is not permanent
production ingress and it does not claim that the route has run.

Any account, provider, allowance, hostname, connector, or runtime observation
is external state and must be checked again before present use:

REVALIDATE_EXTERNAL_STATE

## Alternatives considered

### A. Named Tunnel with a custom domain

This would provide a stable hostname and a durable access design, but requires
a domain, DNS control, credentials, and a different administration
arrangement.

### B. Quick Tunnel bounded validation

This uses the random hostname generated at runtime under
*.trycloudflare.com. It needs no custom domain for this bounded path and
supports a temporary public-HTTPS validation window.

### C. Public Azure ingress or paid network services

A public IP, NAT gateway, load balancer, or bastion would change the approved
zero-incremental-cost and private-origin posture.

## Decision

Accept Option B for bounded validation:

Internet
→ HTTPS at Cloudflare edge
→ random *.trycloudflare.com Quick Tunnel
→ host-level cloudflared
→ loopback HTTP Caddy at 127.0.0.1:8081
→ Docker service DNS api:8080

cloudflared remains host-level infrastructure outside application Compose. The
API host mapping 127.0.0.1:8080 is a readiness and deployment probe only. The
bounded path requires no custom domain, no Origin CA, and no public Azure
ingress. It is temporary and non-production.

The hostname is generated at runtime, may change after connector restart, and
must never be hard-coded in the repository.

## Consequences and trade-offs

The path permits temporary public HTTPS validation without a custom domain or
direct public origin. The trade-offs are an ephemeral hostname, no stable
production identity, no inherent access-control layer on this chosen path,
dependence on the connector and edge, and the need for a new external
validation after a restart.

This decision does not satisfy a permanent configurable WAF or application
rate-limit requirement. A provider concurrency ceiling is not equivalent to an
audited rate-limit control. It also does not provide permanent production
ingress.

No statement here proves that Quick Tunnel has run, that a hostname was
discovered, or that the current Caddy/Compose deployment already matches this
topology.

## Supersession and relationships

This ADR supersedes ADR-0002b only for ingress, domain, Origin CA, and normal
administration assumptions. ADR-0002b's unaffected private-origin,
zero-incremental-cost, recovery, and other boundaries remain recorded there.
ADR-0002e does not supersede the multi-tenant domain decision or the Azure
region/SKU rationale.

## Verification assumptions

Before use, independently verify the external provider state, connector
availability, generated hostname, host-level process, Caddy loopback route,
Docker service route, public HTTPS response, and operational limitations:

REVALIDATE_EXTERNAL_STATE

The accepted decision remains bounded validation, not a production-readiness
claim.
