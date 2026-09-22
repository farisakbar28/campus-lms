---
name: implementation-review
description: Obtain and act on an independent, read-only review of the exact campus-lms implementation candidate before shipping.
---

# Independent implementation review

Use the project `implementation_reviewer` agent in a fresh review context.
Require an effective read-only sandbox and no mutating tool calls. If an
independent reviewer or effective read-only context is unavailable, stop and
report the missing review instead of self-approving.

Give the reviewer the canonical Issue, approved plan and human decision,
base and candidate revisions or complete diff, actual verification output,
and relevant sources. The reviewer must inspect primary evidence rather than
accepting the implementer's summary.

Resolve blocking findings, then ask for a fresh review of the changed
candidate. Bind any passing verdict to an exact revision or complete diff.
After roadmap finalization, obtain reviewer confirmation of the final
candidate if that edit changed the reviewed diff. Do not turn reviewer
approval into human merge or production approval.
