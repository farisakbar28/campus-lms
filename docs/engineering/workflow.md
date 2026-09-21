# Engineering workflow

This is the simple repository-native workflow for human and AI-assisted work.
GitHub Issues, Git branches and commits, Pull Requests, reviews, and CI are
the evidence sources. This document does not define a local state machine,
hash protocol, receipt ledger, or temporary task-artifact protocol.

## Normal sequence

```text
ORIENT
→ BRAINSTORM
→ HUMAN CHOOSES TASK
→ PLAN
→ INDEPENDENT PLAN REVIEW
→ HUMAN APPROVES
→ IMPLEMENT
→ INDEPENDENT IMPLEMENTATION REVIEW
→ CREATE PR
→ CI + HUMAN REVIEW
→ HUMAN MERGE
→ HISTORY SYNC
→ NEXT TASK
```

For production-impacting work, add:

```text
PRODUCTION READINESS REVIEW
→ HUMAN APPROVAL
→ DEPLOY
→ VERIFY
```

### Orient

Read [the documentation map](../README.md), applicable scoped
`AGENTS.md` files, the owning domain and security sources, the current
architecture, and relevant history. Inspect actual code, migrations,
deployment files, Git state, the canonical Issue, and applicable CI evidence.
Do not infer capabilities from plans or stale documentation.

### Brainstorm and human task choice

Explore options, boundaries, dependencies, and risks. The human chooses the
Issue and scope after that orientation. Agents may suggest work but may not
turn a suggestion into authorized implementation, production work, or an
external mutation.

### Plan and independent plan review

A plan states the goal, scope, non-goals, constraints, acceptance criteria,
verification, documentation impact, risk, and release/production impact. A
reviewer independent from the planner checks the current repository and
canonical Issue, tests authority and security boundaries, and records findings
in the normal review channel. Blocking findings must be resolved; security
ambiguity fails closed.

### Human approval and implementation

Implementation starts only after directly human-authored plan approval. Once
approved, the implementer may create a short-lived local branch, edit files,
and create local commits when the task authorizes them. A change to scope,
behavior, architecture, security, authority, acceptance, or production impact
returns to planning, fresh independent review, and human approval.

### Independent implementation review

The implementation reviewer is independent from the implementer and reviews
the exact candidate revision or complete candidate diff, applicable test and
CI results, documentation links, preserved contracts, and security impact.
Blocking findings stop progression. Any relevant change after review requires
a fresh review of the new candidate.

### Pull Request, CI, and merge

Creating or mutating a Pull Request or remote branch requires explicit human
authority. The human must explicitly ask for PR creation after the independent
implementation review passes. Observe the hosted CI result for the exact PR
revision, review the complete diff, and leave merge to the human. CI success
does not authorize merge, release, deployment, destructive work, or external
mutation.

### History sync

When a milestone is known, update the living history in the task Pull Request
or through a separately reviewed short-lived branch and Pull Request. CI and
human merge still apply. Do not write directly to protected `master`, create a
task archive, or turn history into a command transcript.

## Authority stops

- Agents cannot author human task choice, plan approval, residual-risk
  acceptance, merge approval, or production approval.
- Destructive database or data operations require explicit task authority and
  human confirmation. Never use CI success as that authority.
- External or cloud mutations require explicit task authority and human
  confirmation. This includes Issue, Pull Request, remote-branch, provider,
  and production mutations.
- Production deployment requires the separate readiness checklist and human
  approval bound to the exact revision or immutable artifact, target
  environment, and approved scope. A material change invalidates that
  approval.
- Never push, merge, or deploy automatically. Stop when authority, source
  ownership, review independence, or required evidence is unclear.
