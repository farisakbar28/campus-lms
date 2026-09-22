---
name: task-ship
description: Ship an approved and independently reviewed campus-lms task through roadmap finalization, PR, CI, and human merge.
---

# Task shipping

Confirm the human-selected Issue and approved plan. Finalize applicable
roadmap progress and durable material changes before push, without claiming
an unmerged change is already effective on `master`. Confirm that a passing
independent implementation review covers the final complete candidate;
renew review after relevant edits.

Within the selected and approved task, push the branch and create or update
the PR without per-action confirmation. Complete the PR template with actual
scope, risk, verification, review, roadmap, and release impact. Link the
canonical Issue and include `Closes #<issue-number>`.

Observe hosted CI for the exact PR head. If CI or human review leads to a
candidate change, run relevant checks and renew independent review before
presenting the new candidate. The human reviews and merges; GitHub closes
the Issue through the PR reference. Never merge or infer production approval
from CI success.

During a workflow migration, the authority rules on the default branch
remain effective until the migration PR is human-merged.
