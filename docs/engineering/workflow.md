# Engineering workflow

This is the durable workflow for human and Codex-assisted engineering.
GitHub Issues and Pull Requests hold operational work; Git commits, reviews,
tests, and CI provide implementation evidence. The repository skills in
`.agents/skills/` provide reusable stage instructions. No local task state
machine, receipt ledger, or independent plan-review gate is required.
During a workflow migration, the authority rules on the default branch
remain effective until the migration PR is human-merged.

## Normal sequence

```text
COLD CODEX
→ ORIENT
→ BRAINSTORM WITH HUMAN
→ HUMAN CHOOSES TASK
→ AGENT CREATES/UPDATES ISSUE
→ AGENT MAKES PLAN
→ HUMAN APPROVES / REVISES / REJECTS PLAN
→ AGENT CREATES SHORT-LIVED BRANCH
→ IMPLEMENT + TEST
→ MATERIAL DEVIATION? ASK HUMAN
→ INDEPENDENT IMPLEMENTATION REVIEW
→ FIX + REVIEW LOOP IF NEEDED
→ FINALIZE ROADMAP
→ PUSH + CREATE/UPDATE PR
→ CI
→ HUMAN REVIEWS PR
→ HUMAN MERGES
→ ISSUE AUTO-CLOSES
→ NEXT TASK
```

### Orient, brainstorm, and select

Cold Codex reads [the documentation map](../README.md) first and follows
applicable root and scoped `AGENTS.md` instructions. Inspect the owning
domain, security, and ADR sources, current architecture and roadmap, actual
source, Git and GitHub state, and relevant CI evidence. Distinguish observed
facts from proposals and external claims.

Brainstorm options, dependencies, risks, and unresolved decisions with the
human. Only the human chooses the task. After that choice, an agent may create
or update the canonical GitHub Issue for that task. The Issue states intent,
scope, non-goals, acceptance, and relevant authority; it is not an approval
of an implementation plan.

### Plan and human decision

The agent proposes a plan grounded in the Issue and current repository. It
covers goal, exact scope and non-scope, preserved contracts, file impact,
acceptance, verification, documentation and roadmap impact, risk, and
release/production impact. The human directly approves, revises, or rejects
the plan. There is no independent plan review. Implementation cannot begin
until the human approves the current plan. Keep that decision visible in the
normal task record or current conversation; an agent must not author approval
on the human's behalf.

### Implement, test, and handle deviations

From fresh `master`, create a short-lived task branch after plan approval.
Implement the smallest coherent change, add focused tests for behavior
changes, run checks appropriate to risk, and inspect the complete diff.

If requirements, architecture, security, scope, acceptance criteria, or
production impact materially changes, stop affected work and present the
change and its consequences to the human. Continue only after a human
decision and a revised plan and Issue where needed. The prior approval does
not cover the changed scope. Routine fixes within the approved scope need no
new human decision.

### Independent implementation review and roadmap

A reviewer independent from the implementer reads the canonical Issue,
approved plan, complete candidate diff or exact revision, owning sources,
and actual test evidence. The reviewer reports concrete findings and an
approval or changes-required verdict bound to the candidate. Blocking
findings require fixes and another review. Relevant changes after a review
make its verdict stale.

Before pushing, finalize the roadmap's future direction, completed progress,
and durable material-change record as applicable. A roadmap edit after a
passing review changes the candidate; obtain reviewer confirmation of the
final complete diff before shipping it. Do not record an unmerged change as
already effective on `master`.

### Pull Request, CI, and human merge

Within the human-selected task and approved plan, agents may update the
Issue, comment, push, and create or update its PR without per-action
confirmation. The PR links the Issue, includes `Closes #<issue-number>`,
describes actual verification and final independent review, and follows
[the PR template](../../.github/pull_request_template.md).

Observe hosted CI for the exact PR head. Fix failures, rerun relevant checks,
and renew implementation review when the candidate changes. The human
reviews and merges the PR. The closing reference causes GitHub to close the
Issue when the PR is merged into the default branch. CI success and agent
review never authorize merge, release, deployment, or destructive work.

## Production sequence

```text
MERGED
→ PRODUCTION READINESS REVIEW
→ HUMAN PRODUCTION APPROVAL
→ DEPLOY
→ POST-DEPLOY VERIFY
```

Use [the production-readiness checklist](production-readiness.md) when a
merged change can affect production. The independent production reviewer
binds findings to the exact revision or immutable artifact, target
environment, and scope. A human then approves that same binding before
deployment. A material change invalidates the approval. Deployment and
post-deploy verification follow only the authorized scope; report observed
results without secrets or student personal data.

Destructive database or data actions and external or cloud mutations require
explicit task authority and human confirmation. The normal task-scoped
GitHub authority above does not authorize those operations.

## Authority and evidence

- Humans alone choose tasks, decide plans and material deviations, accept
  residual risk, merge PRs, and approve production deployment.
- Agents may carry out normal operations within the selected and approved
  task. They stop when source ownership, authority, reviewer independence, or
  required evidence is unclear.
- `docs/domain.md`, `docs/domain-ai.md`, `SECURITY.md`, and accepted ADRs
  own their respective contracts. `docs/architecture.md` owns the current
  implementation snapshot. `docs/roadmap.md` owns future direction,
  completed progress, and durable material changes.
- GitHub Issues and PRs hold operational work; Git and reviews show changes
  and judgments; tests and CI show automated evidence. None is a substitute
  for another source's authority.
