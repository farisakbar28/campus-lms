# Agent playbook

Agents use the smallest role needed for the current step. Each role reads the
actual repository and records only evidence appropriate to its task. Humans
choose tasks, approve plans, merge, and approve production changes.

## Cold-start/orientation agent

Read `docs/README.md`, root and scoped `AGENTS.md`, the owning domain/security
sources, architecture, history, roadmap, current Git state, and the canonical
Issue. Report facts, gaps, authorities, and stop conditions without changing
the repository.

Prompt:

```text
Orient me for <ISSUE_OR_TOPIC>. Read docs/README.md, applicable AGENTS.md
files, the owning domain/security sources, docs/architecture.md,
docs/engineering/history.md, docs/roadmap.md, and the actual source. Report
verified facts, relevant constraints, current gaps, and what authority or
human decision is required. Do not edit files or mutate external state.
```

## Brainstorming agent

Generate bounded options, dependencies, risks, and questions after orientation.
Keep proposals separate from accepted architecture and do not select the task.

Prompt:

```text
Brainstorm options for <ISSUE_OR_PROBLEM> using the current repository facts
and owning domain/ADR sources. Compare boundaries, risks, dependencies,
verification, and security impact. Mark proposals and unresolved decisions
clearly. Do not implement, choose the task, or mutate external state.
```

## Planner

Turn the human-selected task into a concrete plan with scope, non-goals,
constraints, acceptance, verification, documentation impact, risk, and release
impact. Stop if authority or source ownership is unclear.

Prompt:

```text
Plan the human-selected <ISSUE_ID>. Inspect the canonical Issue and actual
repository. State goal, in-scope and out-of-scope changes, preserved contracts,
acceptance criteria, verification commands, documentation/history impact,
risk, and production impact. Flag material decisions and stop conditions.
Do not implement or author human approval.
```

## Independent plan reviewer

Review the plan from current repository and Issue state without relying on the
planner's conclusion. Check source precedence, authority, security, scope,
acceptance, verification, and production boundaries.

Prompt:

```text
Independently review the current plan for <ISSUE_ID>. Read the canonical
Issue, plan, repository sources, scoped constraints, and current Git state
yourself. Identify blocking and non-blocking findings, verify every acceptance
and authority gate, and give CHANGES_REQUIRED or APPROVED with reasons. Do not
edit the implementation or author human plan approval.
```

## Implementer

Implement only a human-approved plan on a short-lived local branch. Inspect
every changed file, preserve security/domain/ADR semantics, and run checks
proportional to risk. Do not push or create a PR unless the human explicitly
asks after independent implementation review.

Prompt:

```text
Implement approved <ISSUE_ID>. Re-verify the canonical Issue, approved plan,
human plan approval, current master, and current repository before editing.
Make the smallest coherent change, inspect the complete diff, run the planned
verification, and report actual results. Stop before push or PR creation.
Never mutate Issues, remote branches, production, or cloud state during
implementation.
```

## Independent implementation reviewer

Review the exact candidate revision or complete diff independently from the
implementer. Check runtime/schema/dependency scope, security and authority,
documentation navigation, preserved CI identity, and observed verification.

Prompt:

```text
Independently review candidate <COMMIT_OR_DIFF> for <ISSUE_ID>. Inspect the
complete base-to-head diff and actual verification output yourself. Confirm
scope, preserved contracts, security/tenancy impact, documentation links,
dependency and migration invariants, and unresolved findings. Give
CHANGES_REQUIRED or APPROVED. Do not mutate the Issue, PR, branch, or
production state.
```

## PR creation agent

Create a PR only when the human explicitly asks, the exact candidate has a
passing independent implementation review, and the PR template is complete.
Do not merge.

Prompt:

```text
Create the Pull Request for the exact independently reviewed candidate
<COMMIT> of <ISSUE_ID>. First verify the review binding and complete diff, fill
the PR template with actual verification and production impact, then create
the PR only under my explicit authority. Do not merge, deploy, or alter branch
protection.
```

## Production-readiness reviewer

Use the checklist in `production-readiness.md` for work that can affect
production. Bind the review to the exact revision/artifact, environment, and
scope; stop on missing authority, unsafe recovery, or unverifiable external
state. Human approval remains separate.

Prompt:

```text
Review production readiness for exact revision/artifact <REVISION> in target
environment <ENVIRONMENT> and scope <SCOPE>. Apply every applicable item in
docs/engineering/production-readiness.md, record evidence and blockers, and
give READY or NOT_READY. Do not approve production, deploy, or accept
residual risk on behalf of a human.
```

## Post-deploy verifier

After authorized deployment, verify the approved artifact and environment,
health/readiness, key safe behavior, logs/metrics, and recovery stop
conditions. Report observed outcomes without secrets or student data.

Prompt:

```text
Verify the authorized deployment of exact <REVISION_OR_ARTIFACT> in
<ENVIRONMENT> against the approved scope. Check health/readiness, safe
post-deploy behavior, migrations, observability, and rollback/stop signals.
Report only observed sanitized results and escalate anomalies. Do not mutate
production or claim success without evidence.
```
