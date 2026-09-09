# Review records

Contract version: `campus-lms-work/v1`

Append concise records while the work item is active. A record must identify
the review, type, actor and session labels, relevant author labels, fresh
session attestation, exact plan/revision/hash and Issue digest, findings, and
verdict. Implementation reviews additionally bind the exact candidate Git
SHA. Reviewers are read-only and do not create a parallel evidence archive.

## Plan review record

```text
review_id=<ID>
type=PLAN
work_item_id=<subject work item>
actor_label=<reviewer actor>
session_label=<fresh reviewer session>
plan_author_actor_label=<plan author actor>
plan_author_session_label=<plan author session>
fresh_session_attestation=<concise attestation>
plan_revision=<integer>
plan_hash=sha256:<64 lowercase hexadecimal>
issue_digest=sha256:<64 lowercase hexadecimal>
finding_id=<ID>
severity=CRITICAL|HIGH|MEDIUM|LOW
status=OPEN|RESOLVED|ACCEPTED_RESIDUAL_RISK
summary=<concise summary>
residual_risk_comment_url=<exact canonical-Issue comment URL when accepted>
verdict=CHANGES_REQUIRED|APPROVED
```

## Implementation review record

```text
review_id=<ID>
type=IMPLEMENTATION
work_item_id=<subject work item>
actor_label=<reviewer actor>
session_label=<fresh reviewer session>
implementation_author_actor_label=<implementation author actor>
implementation_author_session_label=<implementation author session>
fresh_session_attestation=<concise attestation>
candidate_git_sha=<full 40-character SHA>
plan_revision=<integer>
plan_hash=sha256:<64 lowercase hexadecimal>
issue_digest=sha256:<64 lowercase hexadecimal>
finding_id=<ID>
severity=CRITICAL|HIGH|MEDIUM|LOW
status=OPEN|RESOLVED|ACCEPTED_RESIDUAL_RISK
summary=<concise summary>
residual_risk_comment_url=<exact canonical-Issue comment URL when accepted>
verdict=CHANGES_REQUIRED|APPROVED
```

Agents cannot author plan approval, residual-risk acceptance, production
approval, merge, or a review verdict for their own work.
