# Deployment Engineering Instructions

These instructions apply to `deploy/**`.

## Cost and secrets

- The project targets zero incremental paid infrastructure and tooling.
- Never put secret material in images, Compose files, environment examples,
  logs, or deployment state.
- Never add an automatic paid fallback.
- Treat every free-tier, allowance, quota, price, and provider claim as
  unverified until checked against current external state. Use the marker
  `REVALIDATE_EXTERNAL_STATE` for claims that need fresh verification.

## Images and runtime validation

- Use immutable image references for deployment; do not deploy mutable tags.
- Preserve non-root containers where the existing image design supports them.
- Validate image health and application readiness before declaring deployment
  success.
- Keep explicit memory/resource limits where they are part of the current
  architecture.
- Keep development Compose and production Compose separate. The development
  override must not be used as the production deployment definition.

## Database operations

- Preserve safe backup and restore behavior and verify restore targets are
  disposable and distinct from the source database.
- Never perform an automatic destructive rollback.
- Never automatically roll back database migrations; use an explicitly
  reviewed forward migration or a separately authorized recovery procedure.
- External or cloud mutations require explicit authorization and human
  confirmation.

## Current ingress boundary

The repository's deployment and ingress implementation contains topology that
is stale relative to the accepted bounded Quick Tunnel architecture. Until
that reconciliation is separately authorized:

- do not change ingress topology or Caddy ports;
- do not start `cloudflared`;
- do not represent the current deployment as production-ready;
- preserve current behavior while documenting the stale state honestly.

The permanent-ingress decision and any Quick Tunnel runtime reconciliation are
future authorized infrastructure work.
