---
name: task-implement
description: Implement a human-approved campus-lms task on a short-lived branch with focused verification and material-deviation handling.
---

# Task implementation

Confirm the canonical Issue, current approved plan, direct human approval,
and fresh `master`. Read applicable root and scoped instructions and owning
sources. Create a short-lived task branch from the verified base.

Make the smallest coherent change within the approved scope. Add a focused
test for a behavior change. Run the planned and risk-appropriate checks,
inspect every changed file and the complete base-to-head diff, and report
actual results and limits without secrets or student personal data.

If requirements, architecture, security, scope, acceptance, or production
impact materially changes, stop affected work. Explain the proposed change
and consequences to the human; continue only after their decision and any
needed Issue and plan revision. Routine fixes within the approved scope do
not need a new human decision.

Hand the candidate and observed verification to an implementation reviewer
independent from the implementer. Do not ship an unreviewed candidate. A
relevant change after review requires renewed review.
