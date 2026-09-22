# Production-readiness checklist

Use this checklist when a change can affect production data, availability,
security, configuration, deployment, or an externally reachable system. It is
an engineering review checklist, not a workflow engine. Human production
approval remains mandatory and is separate from implementation review.

## Binding and scope

- [ ] The exact Git revision, immutable image/artifact digest, or migration
      bundle is recorded.
- [ ] The target environment and approved scope are explicit.
- [ ] The proposed rollout matches the reviewed implementation and plan.
- [ ] Any material change after review or approval has returned to planning;
      the previous approval is treated as invalid.

## CI and verification

- [ ] Required hosted CI has passed for the exact revision and the observed
      check identity is recorded.
- [ ] Relevant local tests, formatting, vet/build, integration checks, and
      deployment validation have actual outputs.
- [ ] The complete diff and configuration changes have been reviewed.

## Configuration and secrets

- [ ] Required configuration is present, validated, and fail-closed.
- [ ] Secret sources, rotation, access, and redaction are known; no secret is
      copied into source, images, examples, logs, or tickets.
- [ ] Environment-specific values cannot silently select a weaker security
      mode.

## Migrations and data safety

- [ ] Every migration is forward, reviewed, versioned, and safe for the
      target data volume; applied migration files are unchanged.
- [ ] Locking, compatibility, backfill, constraint, tenant-isolation, and
      authorization effects are understood.
- [ ] Any destructive database/data action has explicit task authority and
      human confirmation.

## Backup and recovery

- [ ] Backup coverage, retention, encryption/access, and last verified restore
      evidence are known for the affected data.
- [ ] Restore targets are disposable and distinct from the source.
- [ ] Recovery steps are reviewed, safe, and do not rely on automatic
      destructive rollback or automatic migration reversal.
- [ ] Stop conditions and escalation owners are named.

## Dependencies and external state

- [ ] Runtime, image, module, provider, and infrastructure dependencies are
      pinned or otherwise verified for the approved scope.
- [ ] External provider, quota, region, cost, certificate, DNS, and account
      claims are freshly verified when material; historic claims are not
      treated as current state.
- [ ] Any external, cloud, or production mutation has explicit task authority
      and human confirmation. Normal task-scoped Issue, PR, and branch
      operations follow the selected and approved task authority in the
      [engineering workflow](workflow.md).

## Security

- [ ] Authentication, authorization, tenant isolation, RLS, logging, and
      privacy behavior have been checked for the changed path.
- [ ] Errors and logs do not expose credentials, tokens, request bodies,
      secrets, or student personal data.
- [ ] Threats from untrusted input, retrieved content, and operator access are
      addressed; security ambiguity fails closed.

## Rollout and recovery

- [ ] Rollout order, compatibility window, traffic/data impact, and resource
      limits are explicit.
- [ ] A safe recovery or stop procedure is available and approved; it does
      not imply automatic destructive rollback.
- [ ] The operator knows when to pause, isolate, or escalate rather than
      continuing on partial evidence.

## Health, readiness, and observability

- [ ] Liveness and dependency readiness checks cover the affected services.
- [ ] Logs, metrics, traces, alerts, dashboards, and retention are sufficient
      for the approved change, without sensitive data.
- [ ] Expected failure signals and ownership for response are known.

## Post-deploy verification

- [ ] Verify the exact deployed revision/artifact and target environment.
- [ ] Check health/readiness, migrations, authorization boundaries, and a
      small set of safe representative behaviors.
- [ ] Confirm observability signals and absence of unexpected errors.
- [ ] Record sanitized observed results and escalate anomalies; do not claim
      success from intention, local checks, or CI alone.

## Human approval

The human approver must approve the exact revision or immutable artifact,
target environment, and scope after the readiness review. The approver must
not rely on CI alone, and an approval is invalidated by a material change.
