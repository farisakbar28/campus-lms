# ADR-0005: API binary health check versus BusyBox probe

- **Date:** 2026-08-18
- **Status:** Proposed
- **Roadmap week:** 2

## Context

The Week 2 Dockerfile requires both a distroless runtime and an HTTP
`HEALTHCHECK`. Distroless intentionally contains no shell or HTTP client.
The initial implementation copied BusyBox into the final runtime and invoked
its `wget` applet. That made the check work, but also introduced a multi-call
binary with utilities that distroless intentionally omits.

The implementation now provides an API binary mode, `/api -healthcheck`, that
performs a bounded HTTP GET to `http://127.0.0.1:$APP_PORT/healthz` and exits
zero only for HTTP 200. The BusyBox implementation remains available as the
`busybox-runtime` Docker target solely for a measured comparison.

## Options considered

1. **Copy BusyBox and run `wget` in the final distroless image**
   - Provides an HTTP client without changing application code.
   - `docker images` reports a final size of 16.6MB.
   - The copied `/busybox` file is 1,214,736 bytes. Docker history reports its
     COPY layer as 1.22MB.
   - The final image reports 4,054,678 bytes through `docker image inspect`.
   - It adds a multi-call binary, including utilities not needed by the API.

2. **Use an API binary health-check mode**
   - Keeps the final runtime to the API binary and the distroless base.
   - `docker images` reports a final size of 14.6MB.
   - The final image reports 3,285,268 bytes through `docker image inspect`.
   - This is 769,410 bytes smaller than the BusyBox runtime by that Docker
     image-size measurement.
   - It adds a small HTTP client code path that requires unit tests and must
     stay compatible with the health endpoint.

3. **Use `golang:alpine` as the production runtime**
   - Makes interactive debugging easier because a shell and utilities exist.
   - A previous local comparison measured 377MB, which is outside the API
     image-size target.

4. **Omit Docker HEALTHCHECK**
   - Avoids adding any probe utility or code path.
   - Does not meet the Week 2 task brief.

## Evidence pending human decision

- [BusyBox versus API-probe exact image sizes](../progress/evidence/week-02/busybox-vs-api-probe-size.txt)
- [BusyBox and API-probe layer histories](../progress/evidence/week-02/busybox-layer-history.txt)
- [API-probe runtime health check](../progress/evidence/week-02/api-probe-healthcheck.txt)
- [Earlier Alpine versus distroless comparison](../progress/evidence/week-02/image-size.txt)

This ADR intentionally has no Decision or Consequences section. The human
owner chooses and records those sections after reviewing the evidence.
