---
name: project-orient
description: Orient a new campus-lms engineering task or topic from current repository and GitHub evidence before brainstorming or planning.
---

# Project orientation

Use this skill at cold start or when the relevant repository state is unknown.

1. Read `docs/README.md` first, then root and applicable scoped `AGENTS.md`
   files before entering a narrower tree.
2. Read the owning domain, AI, security, and accepted ADR sources for the
   topic. Use `docs/architecture.md` for current implementation facts and
   `docs/roadmap.md` for direction, completed progress, and material changes.
3. Inspect actual code, migrations, deployment files, Git state, the canonical
   Issue if one exists, relevant PRs, and CI evidence before making claims.
   Recheck external or cloud facts before relying on them.
4. Report verified facts, source ownership, gaps, dependencies, risks, and
   decisions that need a human. Mark proposals and uncertain claims clearly.
5. Brainstorm bounded options with the human when requested. The human alone
   chooses the task.

Orientation is read-only. Do not create an Issue, edit files, choose a task,
or mutate external state during this stage. The human's task selection
authorizes task-scoped Issue creation or updates under the engineering
workflow.
