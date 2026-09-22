# campus-lms documentation map

This is the starting point for a new engineer, reviewer, or agent. It routes
questions to the authority that owns the answer and keeps implementation facts,
architecture decisions, product direction, and operational evidence distinct.

## What are we building?

`campus-lms` is a multi-tenant Learning Management System for higher
education. The LMS owns learning content and activity workflows. Institutional
academic systems such as SIAKAD remain authoritative for the academic master
and reference data assigned to them by the domain model.

Staff-side AI is a future, optional, non-authoritative capability. Students are
not AI users, and AI must never decide grades, attendance, enrollment,
pass/fail, discipline, or another authoritative academic outcome.

Read the [domain contract](domain.md) for entities, ownership, tenancy, and
authorization, and the [AI contract](domain-ai.md) for the staff-side AI
boundary.

## How is it designed?

The current implementation is summarized in
[architecture.md](architecture.md). It distinguishes code facts, accepted
decisions, proposed decisions, known gaps, and future work. Detailed domain
semantics remain in the domain documents; durable architectural choices remain
in the [accepted ADRs](adr/README.md).

The current direction is a Go API with explicit SQL and PostgreSQL tenant
isolation, surrounded by deliberately bounded deployment and recovery tooling.
That summary does not claim a production deployment, provider entitlement, or
complete authentication composition.

## What are the security constraints?

Read [SECURITY.md](../SECURITY.md) for enduring security requirements. The
most important boundaries are:

- tenant identity comes from a trusted authenticated principal, not a client
  header, query parameter, or body field;
- application authorization, transaction-local PostgreSQL RLS, and composite
  tenant constraints work together;
- missing or invalid security configuration fails closed;
- secrets, tokens, request bodies, and student personal data never appear in
  source, logs, errors, tests, issues, commits, or images;
- applied migrations are immutable and data-destructive work requires explicit
  task authority plus human confirmation; and
- external or cloud mutations require explicit task authority plus human
  confirmation.

The scoped rules in [apps/api/AGENTS.md](../apps/api/AGENTS.md) and
[deploy/AGENTS.md](../deploy/AGENTS.md) add API and deployment constraints.

## How do engineers and agents work?

Use the simple repository-native sequence in
[engineering/workflow.md](engineering/workflow.md). The
[agent playbook](engineering/agent-playbook.md) provides copy-paste prompts
for common roles. Humans choose tasks, approve plans, review intent, merge,
and approve production changes. Independent plan and implementation reviews
are required; a material change returns the work to planning and fresh review.

The [production-readiness checklist](engineering/production-readiness.md) is a
checklist for production-impacting work, not a workflow engine.

## Where are we now?

[architecture.md](architecture.md) owns the current repository implementation
snapshot. [engineering/history.md](engineering/history.md) records concise,
verified milestones and changed assumptions. Current hosted CI status belongs
to GitHub Actions for the relevant revision, not to a static status claim.

The repository has a runnable API and PostgreSQL/RLS foundations, token and
session primitives, local deployment tooling, and accepted tenant-context
architecture. Runtime authentication composition, product capabilities beyond
the current foundation, production observability, current-schema recovery, and
permanent ingress remain separate gaps or future work.

## Where are we going?

[roadmap.md](roadmap.md) contains future product and engineering direction,
dependency order, classifications, and unresolved decisions. GitHub Issues
contain the canonical operational backlog and task intent. Git commits, Pull
Requests, reviews, and CI provide implementation evidence.

## Recommended reading order

1. [Domain contract](domain.md) and [security requirements](../SECURITY.md).
2. [Current architecture](architecture.md) and [accepted ADRs](adr/README.md).
3. [Engineering workflow](engineering/workflow.md) and
   [agent playbook](engineering/agent-playbook.md).
4. [Engineering history](engineering/history.md) and
   [future roadmap](roadmap.md).
5. [Contributing guide](../CONTRIBUTING.md) and applicable scoped
   `AGENTS.md` files.
