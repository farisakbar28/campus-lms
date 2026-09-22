---
name: task-plan
description: Plan a human-selected campus-lms Issue for direct human approval, revision, or rejection before implementation.
---

# Task planning

Use this skill after the human selects a task. Confirm or create its canonical
GitHub Issue within that selection, then read the current repository and owning
sources. The Issue records task intent; it is not plan approval.

Produce a plan with the goal, exact scope and non-scope, affected files,
preserved domain/security/ADR contracts, acceptance criteria, verification,
documentation and roadmap impact, risks, dependencies, and release or
production impact. Separate verified facts from decisions or assumptions.
Identify any destructive, external, cloud, or production operation that needs
its own authority.

Present the current plan directly to the human for approval, revision, or
rejection. There is no independent plan review. Do not implement or create a
task branch before direct human approval. An agent must not write or imply
human approval on the human's behalf. If a future cold session cannot verify
the approved plan and direct human decision, it must stop and request that
evidence before implementation.
