# Rule 30 — Docker, Compose, Deployment

## Images

- Multi-stage always. Build stage compiles; runtime stage carries only the binary.
- Go runtime: `distroless/static` or `scratch`. Target: **API image < 25 MB**.
- Non-root user, always. Verify with `docker inspect`.
- `HEALTHCHECK` present.
- Multi-arch builds (`buildx`, amd64 + arm64) — production is x86 today, but portability is cheap now and expensive later.
- Layer ordering: dependency download before source copy, so rebuilds hit cache.
- Never put a secret in a build ARG or ENV — it persists in the image layers and is trivially extractable.

## Compose

- Use **profiles**: `core` (api, web, postgres, redis), `storage` (minio), `obs` (otel, prometheus, grafana, loki, tempo).
- Never start `obs` together with `core` without warning — the laptop may have only 8 GB.
- `healthcheck` on every service; `depends_on` with `condition: service_healthy`.
- **Memory limits on every service.** This is not optional on this hardware.
- Named volumes for data; bind mounts only for dev hot-reload.
- Production uses an explicit `docker-compose.prod.yml`; the dev override must never reach the server.

## Production constraints (Azure B1s: 1 vCPU, 1 GB RAM)

- Postgres does **not** run on the VM — it is on Neon.
- Frontend does **not** run on the VM — it is on Cloudflare Pages.
- A 2 GB swap file exists; expect swap usage, monitor for thrashing.
- Before adding any service to production, state its memory footprint and what you will remove to make room.
- Prefer scheduled batch work over resident daemons.

## Deployment

- Images referenced by **digest**, never `:latest`. Reproducibility beats convenience.
- Deploy = pull → migrate → up → health check → verify. Automatic rollback if `/readyz` fails 3× within 60 s.
- Zero-downtime: two replicas behind Caddy, or blue-green. Prove it with a `curl` loop during deploy showing 0 failed requests.
- Record deploy duration; it becomes the "lead time" metric in the report.

## Cost discipline (real money, limited credit)

- `az vm deallocate` when idle — `stop` from inside the OS still bills compute.
- Standard SSD, not Premium.
- Scale-up is temporary and must be scaled back down the same week; note the cost in the report.
- Never propose AKS, Application Gateway, or any always-on managed service. They will consume the $100 credit quickly.
