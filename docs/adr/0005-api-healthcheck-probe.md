# ADR-0005: API Binary Health-Check Probe

- Date: 2026-08-18
- Status: Accepted

## Context

The selected distroless API runtime intentionally does not include a shell or
HTTP client. Docker still needs an HTTP health check for the API liveness
endpoint. The repository contains an API binary mode that performs the probe
itself.

## Alternatives considered

1. Copy BusyBox into the distroless runtime and invoke its HTTP client. This
   provides a probe but adds a multi-call utility that the application does
   not otherwise need.
2. Use the API binary's health-check mode. This keeps the selected runtime
   minimal and makes the probe behavior part of the API code path.
3. Use a shell-based general-purpose runtime. This simplifies interactive
   debugging but adds a larger runtime and weakens the purpose of the
   distroless selection.
4. Omit Docker HEALTHCHECK. This avoids a probe but does not satisfy the
   liveness requirement.

## Decision

Use the API binary health-check mode:

/api -healthcheck

The mode performs a bounded GET request to the local /healthz endpoint and
exits successfully only for HTTP 200. The selected distroless runtime does not
include BusyBox for this purpose.

## Consequences and trade-offs

The final runtime contains only the API and its selected base image, with no
additional HTTP utility. The probe is a small code path that must remain
compatible with /healthz and must be covered by focused tests.

Distroless remains less convenient for interactive debugging because it does
not provide a shell or general diagnostic tools. Debugging therefore uses
external or ephemeral diagnostic tooling rather than adding utilities to the
selected runtime.

## Supersession and relationships

No later ADR supersedes this decision. Docker and API changes must preserve
the relationship between the health-check mode and /healthz.

## Verification assumptions

The API health-check mode, Docker HEALTHCHECK configuration, and focused probe
tests are IMPLEMENTATION_FACT supported by the repository. Later changes
must verify that the local endpoint, timeout, non-200 behavior, selected
runtime, and Docker health status remain aligned.
