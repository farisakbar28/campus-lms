---
name: production-release
description: Review and carry out an explicitly approved campus-lms production release after merge, then verify the deployed artifact.
---

# Production release

Use this skill only after the relevant PR is human-merged and a production
release is in scope. Read `deploy/AGENTS.md`, `SECURITY.md`, applicable
ADRs, and `docs/engineering/production-readiness.md`.

Have the independent `production_reviewer` inspect the exact revision or
immutable artifact, target environment, scope, rollout, recovery, current
external state, and readiness evidence in an effective read-only context.
Stop on a NOT_READY verdict or missing evidence. A reviewer cannot approve
deployment or accept residual risk for the human.

Obtain human production approval bound to that exact artifact, environment,
and scope before deployment. A material change invalidates the approval and
requires renewed readiness review and human decision. Destructive data
actions and external or cloud mutations require explicit task authority and
human confirmation; never infer them from merge or CI.

Perform only the authorized deployment actions. Verify the deployed artifact,
health and readiness, safe representative behavior, observability, and stop
conditions. Report sanitized observations and escalate anomalies. Never
claim production success from a plan or CI result alone.
