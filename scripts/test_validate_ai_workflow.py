"""Focused tests for the campus-lms-work/v1 repository workflow validator."""

from __future__ import annotations

import json
import subprocess
import tempfile
import unittest
from pathlib import Path

from scripts.validate_ai_workflow import (
    CONTRACT_VERSION,
    PLAN_BEGIN,
    PLAN_END,
    issue_digest,
    normalize_document,
    plan_hash_from_text,
    validate_approval_comment,
    validate_cleanup_diff,
    validate_repository,
)


def bind_fixture_work(
    root: Path,
    *,
    candidate: str = "b" * 40,
    status: str = "IMPLEMENTING",
) -> dict[str, str | int | None]:
    work = root / "work" / "active" / "ENG-016" / "WORK.md"
    text = work.read_text(encoding="utf-8")
    text = text.replace("Status: `IMPLEMENTING`", f"Status: `{status}`")
    text = text.replace(
        "Candidate Git commit SHA: `NONE`", f"Candidate Git commit SHA: `{candidate}`"
    )
    work.write_text(text, encoding="utf-8")
    from scripts.validate_ai_workflow import parse_work

    return parse_work(text)


def plan_review_record(
    root: Path,
    *,
    verdict: str = "APPROVED",
    omitted: str | None = None,
    **overrides: str,
) -> str:
    from scripts.validate_ai_workflow import parse_work

    work = parse_work(
        (root / "work" / "active" / "ENG-016" / "WORK.md").read_text(
            encoding="utf-8"
        )
    )
    fields = {
        "review_id": "ENG-016-PLAN-REVIEW-099",
        "type": "PLAN",
        "work_item_id": str(work["work_id"]),
        "plan_revision": str(work["plan_revision"]),
        "plan_hash": str(work["plan_hash"]),
        "issue_digest": str(work["issue_digest"]),
        "actor_label": "Fixture plan reviewer",
        "session_label": "Fixture plan reviewer session",
        "plan_author_actor_label": str(work["plan_author_actor"]),
        "plan_author_session_label": str(work["plan_author_session"]),
        "fresh_session_attestation": "Fresh independent plan review session.",
        "verdict": verdict,
    }
    fields.update(overrides)
    if omitted is not None:
        fields.pop(omitted, None)
    return "\n".join(f"{key}={value}" for key, value in fields.items())


def replace_with_plan_review(
    root: Path,
    *,
    verdict: str = "APPROVED",
    omitted: str | None = None,
    **overrides: str,
) -> None:
    reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
    reviews.write_text(
        "# ENG-016 reviews\n\n## Review 1\n\n"
        + plan_review_record(
            root, verdict=verdict, omitted=omitted, **overrides
        )
        + "\n",
        encoding="utf-8",
    )


def legacy_plan_review_record(
    root: Path, *, omitted: str | None = None, **overrides: str
) -> str:
    from scripts.validate_ai_workflow import parse_work

    work = parse_work(
        (root / "work" / "active" / "ENG-016" / "WORK.md").read_text(
            encoding="utf-8"
        )
    )
    fields = {
        "Review ID": "ENG-016-PLAN-REVIEW-099",
        "Review type": "PLAN",
        "Subject work item": str(work["work_id"]),
        "Subject plan revision": str(work["plan_revision"]),
        "Exact reviewed plan checkpoint SHA": str(work["plan_hash"]),
        "Independently recomputed normative SHA": str(work["plan_hash"]),
        "Exact Issue specification digest": str(work["issue_digest"]),
        "Independently recomputed Issue specification digest": str(
            work["issue_digest"]
        ),
        "Actor": "Fixture plan reviewer",
        "Session label": "Fixture plan reviewer session",
        "Plan author actor label": str(work["plan_author_actor"]),
        "Plan author session label": str(work["plan_author_session"]),
        "Fresh-session attestation": "Fresh independent plan review session.",
        "Verdict": "APPROVED",
    }
    fields.update(overrides)
    if omitted is not None:
        fields.pop(omitted, None)
    return "\n".join(f"- {key}: `{value}`" for key, value in fields.items())


def replace_with_legacy_plan_review(
    root: Path, *, omitted: str | None = None, **overrides: str
) -> None:
    reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
    reviews.write_text(
        "# ENG-016 reviews\n\n## Review 1\n\n"
        + legacy_plan_review_record(root, omitted=omitted, **overrides)
        + "\n",
        encoding="utf-8",
    )


def append_exact_implementation_review(
    root: Path,
    *,
    candidate: str,
    verdict: str = "APPROVED",
    **overrides: str,
) -> None:
    from scripts.validate_ai_workflow import parse_work

    work = parse_work(
        (root / "work" / "active" / "ENG-016" / "WORK.md").read_text(
            encoding="utf-8"
        )
    )
    fields = {
        "review_id": "ENG-016-IMPLEMENTATION-REVIEW-002",
        "type": "IMPLEMENTATION",
        "work_item_id": str(work["work_id"]),
        "plan_revision": str(work["plan_revision"]),
        "plan_hash": str(work["plan_hash"]),
        "issue_digest": str(work["issue_digest"]),
        "actor_label": "Fixture reviewer",
        "session_label": "Fixture reviewer session",
        "implementation_author_actor_label": str(work["implementation_actor"]),
        "implementation_author_session_label": str(work["implementation_session"]),
        "fresh_session_attestation": "Fresh independent review session.",
        "candidate_git_sha": candidate,
        "verdict": verdict,
    }
    fields.update(overrides)
    record = "\n".join(f"{key}={value}" for key, value in fields.items())
    reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
    reviews.write_text(
        reviews.read_text(encoding="utf-8").rstrip() + "\n\n## Review 7\n\n" + record + "\n",
        encoding="utf-8",
    )


def append_finding(
    root: Path,
    *,
    finding_id: str = "ENG-016-IMPL-TEST",
    severity: str = "HIGH",
    status: str = "OPEN",
    summary: str = "unsafe",
    residual_url: str | None = None,
) -> None:
    lines = [
        f"finding_id={finding_id}",
        f"severity={severity}",
        f"status={status}",
        f"summary={summary}",
    ]
    if residual_url is not None:
        lines.append(f"residual_risk_comment_url={residual_url}")
    reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
    reviews.write_text(
        reviews.read_text(encoding="utf-8").rstrip() + "\n\n" + "\n".join(lines) + "\n",
        encoding="utf-8",
    )


class WorkflowValidatorTests(unittest.TestCase):
    def test_plan_hash_normalizes_line_endings_and_mutable_text(self) -> None:
        document = (
            "header\r\n"
            f"{PLAN_BEGIN}\r\n"
            "\t normative line  \r\n"
            "\r\n"
            f"{PLAN_END}\r\n"
            "Status: IMPLEMENTING\r\n"
        )
        equivalent = (
            "different mutable header\n"
            f"{PLAN_BEGIN}\n"
            "\t normative line\n"
            f"{PLAN_END}\n"
            "Status: READY_FOR_PR\n"
        )
        self.assertEqual(plan_hash_from_text(document), plan_hash_from_text(equivalent))

    def test_plan_hash_rejects_duplicate_markers(self) -> None:
        with self.assertRaises(ValueError):
            plan_hash_from_text(
                f"{PLAN_BEGIN}\nbody\n{PLAN_BEGIN}\n{PLAN_END}\n"
            )

    def test_issue_digest_ignores_comments_and_normalizes_fields(self) -> None:
        issue = {"title": "Title  \r\n", "body": "Body\r\n", "comments": []}
        with_comments = dict(issue, comments=[{"body": "not part of the digest"}])
        equivalent = {"title": "Title\n", "body": "Body\n\n"}
        self.assertEqual(issue_digest(issue), issue_digest(with_comments))
        self.assertEqual(issue_digest(issue), issue_digest(equivalent))
        self.assertNotEqual(issue_digest(issue), issue_digest({"title": "Other", "body": "Body"}))

    def test_repository_fixture_passes_and_contract_is_machine_checkable(self) -> None:
        with temporary_repository() as root:
            self.assertEqual(validate_repository(root), [])

    def test_canonical_plan_review_requires_exact_current_bindings(self) -> None:
        required = (
            "review_id",
            "type",
            "work_item_id",
            "plan_revision",
            "plan_hash",
            "issue_digest",
            "actor_label",
            "session_label",
            "plan_author_actor_label",
            "plan_author_session_label",
            "fresh_session_attestation",
            "verdict",
        )
        for field in required:
            with self.subTest(missing=field), temporary_repository() as root:
                replace_with_plan_review(root, omitted=field)
                errors = validate_repository(root)
                self.assertTrue(
                    any(f"missing {field}" in error for error in errors), errors
                )

        mismatches = {
            "work_item_id": "ENG-017",
            "plan_revision": "3",
            "plan_hash": "sha256:" + "0" * 64,
            "issue_digest": "sha256:" + "0" * 64,
            "plan_author_actor_label": "Foreign planner",
            "plan_author_session_label": "Foreign planner session",
        }
        for field, value in mismatches.items():
            with self.subTest(mismatch=field), temporary_repository() as root:
                replace_with_plan_review(root, **{field: value})
                errors = validate_repository(root)
                self.assertTrue(
                    any(
                        f"{field} does not match current WORK.md" in error
                        for error in errors
                    ),
                    errors,
                )

    def test_canonical_plan_review_requires_independence_freshness_and_verdict(self) -> None:
        with temporary_repository() as root:
            from scripts.validate_ai_workflow import parse_work

            work = parse_work(
                (root / "work" / "active" / "ENG-016" / "WORK.md").read_text(
                    encoding="utf-8"
                )
            )
            replace_with_plan_review(
                root,
                actor_label=str(work["plan_author_actor"]),
                session_label=str(work["plan_author_session"]),
            )
            errors = validate_repository(root)
            self.assertTrue(
                any("reviewer actor equals author actor" in error for error in errors)
            )
            self.assertTrue(
                any("reviewer session equals author session" in error for error in errors)
            )

        with temporary_repository() as root:
            replace_with_plan_review(root, verdict="CHANGES_REQUIRED")
            errors = validate_repository(root)
            self.assertTrue(
                any("verdict is not APPROVED" in error for error in errors)
            )

        with temporary_repository() as root:
            replace_with_plan_review(root)
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                reviews.read_text(encoding="utf-8")
                + "plan_hash=sha256:"
                + "0" * 64
                + "\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("repeats plan_hash" in error for error in errors))

    def test_minimal_prose_plan_review_is_not_an_approval(self) -> None:
        with temporary_repository() as root:
            from scripts.validate_ai_workflow import parse_work

            work = parse_work(
                (root / "work" / "active" / "ENG-016" / "WORK.md").read_text(
                    encoding="utf-8"
                )
            )
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                "# ENG-016 reviews\n\n## Review 1\n\n"
                "- Review ID: `ENG-016-PLAN-REVIEW-099`\n"
                "- Review type: `PLAN`\n"
                f"- Exact reviewed plan checkpoint SHA: `{work['plan_hash']}`\n"
                "- Verdict: `APPROVED`\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("missing work_item_id" in error for error in errors))
            self.assertTrue(any("missing plan_revision" in error for error in errors))

    def test_current_bootstrap_legacy_plan_review_requires_explicit_bindings(self) -> None:
        with temporary_repository() as root:
            replace_with_legacy_plan_review(root)
            self.assertEqual(validate_repository(root), [])

        with temporary_repository() as root:
            record = legacy_plan_review_record(
                root, omitted="Subject plan revision"
            )
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                "# ENG-016 reviews\n\n## Review 1\n\n" + record + "\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("missing plan_revision" in error for error in errors))

        with temporary_repository() as root:
            record = legacy_plan_review_record(
                root, omitted="Independently recomputed normative SHA"
            )
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                "# ENG-016 reviews\n\n## Review 1\n\n" + record + "\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(
                any("missing plan_hash_recomputed" in error for error in errors)
            )

    def test_invalid_state_and_equal_review_identity_are_rejected(self) -> None:
        with temporary_repository() as root:
            work = root / "work" / "active" / "ENG-016" / "WORK.md"
            text = work.read_text(encoding="utf-8").replace(
                "Status: `IMPLEMENTING`", "Status: `DONE`"
            )
            work.write_text(text, encoding="utf-8")
            reviews = work.parent / "REVIEWS.md"
            reviews.write_text(
                reviews.read_text(encoding="utf-8")
                + "\n## Review 1\n"
                + "- Review type: `IMPLEMENTATION`\n"
                + "- Actor: `Codex implementer`\n"
                + "- Session label: `ENG-016-implement-r4`\n"
                + "- Implementation author actor label: `Codex implementer`\n"
                + "- Implementation author session label: `ENG-016-author`\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("post-merge/DONE" in error for error in errors))
            self.assertTrue(any("reviewer actor equals" in error for error in errors))

    def test_finding_severity_gate_is_fail_closed(self) -> None:
        with temporary_repository() as root:
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                reviews.read_text(encoding="utf-8")
                + "\nfinding_id=F-1\nseverity=HIGH\nstatus=OPEN\nsummary=unsafe\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("HIGH finding F-1" in error for error in errors))

    def test_invalid_enum_hash_archive_and_stale_issue_are_rejected(self) -> None:
        with temporary_repository() as root:
            work = root / "work" / "active" / "ENG-016" / "WORK.md"
            work.write_text(
                work.read_text(encoding="utf-8").replace(
                    "Risk: `CRITICAL`", "Risk: `UNSAFE`"
                ),
                encoding="utf-8",
            )
            (root / "work" / "completed").mkdir()
            errors = validate_repository(
                root, issue={"title": "stale", "body": "Fixture body"}
            )
            self.assertTrue(any("invalid risk" in error for error in errors))
            self.assertTrue(any("archive shape" in error for error in errors))
            self.assertTrue(any("digest" in error for error in errors))

    def test_implementation_review_state_requires_candidate_binding(self) -> None:
        with temporary_repository() as root:
            work = root / "work" / "active" / "ENG-016" / "WORK.md"
            work.write_text(
                work.read_text(encoding="utf-8").replace(
                    "Status: `IMPLEMENTING`", "Status: `IMPLEMENTATION_REVIEW`"
                ),
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("requires a candidate Git SHA" in error for error in errors))

    def test_ready_requires_exact_implementation_review_bindings(self) -> None:
        required = (
            "review_id",
            "type",
            "work_item_id",
            "plan_revision",
            "plan_hash",
            "issue_digest",
            "actor_label",
            "session_label",
            "implementation_author_actor_label",
            "implementation_author_session_label",
            "fresh_session_attestation",
            "candidate_git_sha",
            "verdict",
        )
        for missing in required:
            with self.subTest(missing=missing), temporary_repository() as root:
                candidate = "b" * 40
                bind_fixture_work(root, candidate=candidate, status="READY_FOR_PR")
                append_exact_implementation_review(root, candidate=candidate)
                reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
                reviews.write_text(
                    reviews.read_text(encoding="utf-8").replace(
                        f"{missing}=", f"omitted_{missing}="
                    ),
                    encoding="utf-8",
                )
                errors = validate_repository(root)
                self.assertTrue(
                    any(
                        "requires an APPROVED implementation review" in error
                        or f"missing {missing}" in error
                        for error in errors
                    )
                )

        with temporary_repository() as root:
            candidate = "b" * 40
            bind_fixture_work(root, candidate=candidate, status="READY_FOR_PR")
            append_exact_implementation_review(root, candidate=candidate)
            self.assertEqual(validate_repository(root), [])

    def test_implementation_review_mismatch_duplicate_and_prose_are_rejected(self) -> None:
        mismatches = {
            "work_item_id": "ENG-017",
            "plan_revision": "5",
            "plan_hash": "sha256:" + "0" * 64,
            "issue_digest": "sha256:" + "0" * 64,
            "candidate_git_sha": "c" * 40,
            "implementation_author_actor_label": "Foreign implementer",
        }
        for field, value in mismatches.items():
            with self.subTest(field=field), temporary_repository() as root:
                candidate = "b" * 40
                bind_fixture_work(root, candidate=candidate, status="READY_FOR_PR")
                append_exact_implementation_review(
                    root, candidate=candidate, **{field: value}
                )
                errors = validate_repository(root)
                self.assertTrue(
                    any(
                        "requires an APPROVED implementation review" in error
                        for error in errors
                    )
                )

        with temporary_repository() as root:
            candidate = "b" * 40
            bind_fixture_work(root, candidate=candidate, status="READY_FOR_PR")
            append_exact_implementation_review(root, candidate=candidate)
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                reviews.read_text(encoding="utf-8") + f"candidate_git_sha={candidate}\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("repeats candidate_git_sha" in error for error in errors))

        with temporary_repository() as root:
            candidate = "b" * 40
            bind_fixture_work(root, candidate=candidate, status="READY_FOR_PR")
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                reviews.read_text(encoding="utf-8")
                + f"\n## Review 7\n\n- Review type: `IMPLEMENTATION`\n"
                + f"- Candidate Git SHA: `{candidate}`\n- Verdict: `APPROVED`\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("missing work_item_id" in error for error in errors))

    def test_approval_comment_requires_exact_binding(self) -> None:
        with temporary_repository() as root:
            work = root / "work" / "active" / "ENG-016" / "WORK.md"
            from scripts.validate_ai_workflow import parse_work

            parsed = parse_work(work.read_text(encoding="utf-8"))
            comment = {
                "html_url": parsed["approval_url"],
                "body": "\n".join(
                    (
                        "approval_type=PLAN",
                        "work_item_id=ENG-016",
                        "plan_revision=4",
                        f"plan_hash={parsed['plan_hash']}",
                        "decision=APPROVED",
                    )
                ),
                "user": {"type": "User"},
                "created_at": "2026-09-06T02:53:15Z",
                "updated_at": "2026-09-06T02:53:15Z",
            }
            self.assertEqual(validate_approval_comment(comment, parsed), [])
            comment["body"] = comment["body"].replace("APPROVED", "REJECTED")
            self.assertTrue(validate_approval_comment(comment, parsed))

    def test_approval_comment_must_bind_to_canonical_issue_repository(self) -> None:
        with temporary_repository() as root:
            from scripts.validate_ai_workflow import parse_work

            work = root / "work" / "active" / "ENG-016" / "WORK.md"
            parsed = parse_work(work.read_text(encoding="utf-8"))
            comment = {
                "html_url": parsed["approval_url"],
                "body": "\n".join(
                    (
                        "approval_type=PLAN",
                        "work_item_id=ENG-016",
                        "plan_revision=4",
                        f"plan_hash={parsed['plan_hash']}",
                        "decision=APPROVED",
                    )
                ),
                "user": {"type": "User"},
            }
            self.assertEqual(validate_approval_comment(comment, parsed), [])
            parsed["approval_url"] = (
                "https://github.com/foreign/repository/issues/16#issuecomment-1"
            )
            comment["html_url"] = parsed["approval_url"]
            self.assertTrue(validate_approval_comment(comment, parsed))

            work.write_text(
                work.read_text(encoding="utf-8").replace(
                    "https://github.com/example/campus-lms/issues/16#issuecomment-1",
                    "https://github.com/foreign/repository/issues/16#issuecomment-1",
                ),
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("not on the canonical Issue" in error for error in errors))

    def test_live_issue_json_distinguishes_api_and_web_urls(self) -> None:
        with temporary_repository() as root:
            issue = {
                "title": "Fixture issue",
                "body": "Fixture body",
                "number": 16,
                "url": "https://api.github.com/repos/example/campus-lms/issues/16",
                "html_url": "https://github.com/example/campus-lms/issues/16",
                "repository_url": "https://api.github.com/repos/example/campus-lms",
            }
            self.assertEqual(validate_repository(root, issue=issue), [])

            issue["url"] = issue["html_url"]
            errors = validate_repository(root, issue=issue)
            self.assertTrue(any("API url" in error for error in errors))

    def test_malformed_live_issue_json_fails_closed(self) -> None:
        with temporary_repository() as root:
            issue = {
                "title": "Fixture issue",
                "body": "Fixture body",
                "number": "16",
                "url": "https://api.github.com/repos/foreign/repository/issues/16",
            }
            errors = validate_repository(root, issue=issue)
            self.assertTrue(any("number is missing or malformed" in error for error in errors))
            self.assertTrue(any("repository identity" in error for error in errors))
            self.assertTrue(any("html_url" in error for error in errors))

    def test_residual_risk_requires_current_candidate_human_comment_and_finding(self) -> None:
        candidate = "b" * 40
        comment_url = "https://github.com/example/campus-lms/issues/16#issuecomment-42"
        comment = {
            "id": 42,
            "url": "https://api.github.com/repos/example/campus-lms/issues/comments/42",
            "html_url": comment_url,
            "body": "\n".join(
                (
                    "decision=ACCEPTED_RESIDUAL_RISK",
                    "work_item_id=ENG-016",
                    "finding_id=ENG-016-IMPL-TEST",
                    "plan_revision=4",
                    "plan_hash=PLACEHOLDER",
                    f"candidate_git_sha={candidate}",
                    "rationale=Human decision for bounded residual risk.",
                )
            ),
            "user": {"type": "User"},
        }
        with temporary_repository() as root:
            work = bind_fixture_work(root, candidate=candidate)
            comment["body"] = comment["body"].replace(
                "plan_hash=PLACEHOLDER", f"plan_hash={work['plan_hash']}"
            )
            append_finding(
                root,
                severity="MEDIUM",
                status="ACCEPTED_RESIDUAL_RISK",
                residual_url=comment_url,
            )
            self.assertEqual(validate_repository(root, residual_comments=[comment]), [])

            wrong_candidate = dict(comment)
            wrong_candidate["body"] = wrong_candidate["body"].replace(candidate, "c" * 40)
            errors = validate_repository(root, residual_comments=[wrong_candidate])
            self.assertTrue(any("candidate_git_sha" in error for error in errors))

            wrong_finding = dict(comment)
            wrong_finding["body"] = wrong_finding["body"].replace(
                "finding_id=ENG-016-IMPL-TEST", "finding_id=ENG-016-IMPL-OTHER"
            )
            errors = validate_repository(root, residual_comments=[wrong_finding])
            self.assertTrue(any("applicable finding" in error for error in errors))

            bot_comment = dict(comment)
            bot_comment["user"] = {"type": "Bot"}
            errors = validate_repository(root, residual_comments=[bot_comment])
            self.assertTrue(any("not a GitHub User" in error for error in errors))

        with temporary_repository() as root:
            bind_fixture_work(root, candidate=candidate)
            append_finding(
                root,
                severity="MEDIUM",
                status="ACCEPTED_RESIDUAL_RISK",
                residual_url=comment_url,
            )
            errors = validate_repository(root)
            self.assertTrue(any("requires exactly one supplied" in error for error in errors))

    def test_finding_status_enum_and_high_critical_residual_bypass_fail_closed(self) -> None:
        for severity in ("HIGH", "CRITICAL"):
            with self.subTest(severity=severity), temporary_repository() as root:
                append_finding(root, severity=severity, status="CLOSED")
                errors = validate_repository(root)
                self.assertTrue(any("invalid finding status" in error for error in errors))
                self.assertTrue(any(f"{severity} finding" in error for error in errors))

            with self.subTest(severity=f"{severity}-accepted"), temporary_repository() as root:
                append_finding(
                    root,
                    severity=severity,
                    status="ACCEPTED_RESIDUAL_RISK",
                    residual_url="https://github.com/example/campus-lms/issues/16#issuecomment-42",
                )
                errors = validate_repository(root)
                self.assertTrue(any("cannot accept residual risk" in error for error in errors))

        with temporary_repository() as root:
            append_finding(root, status="OPEN")
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                reviews.read_text(encoding="utf-8") + "status=OPEN\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(any("repeats status" in error for error in errors))

    def test_cleanup_diff_allows_only_active_pair_deletions(self) -> None:
        with temporary_git_repository() as root:
            base = git_commit(root, "base", {"README.md": "base\n"})
            candidate = git_commit(
                root,
                "candidate",
                {
                    "work/active/ENG-016/WORK.md": "work\n",
                    "work/active/ENG-016/REVIEWS.md": "reviews\n",
                },
            )
            head = git_commit(
                root,
                "cleanup",
                {
                    "work/active/ENG-016/WORK.md": None,
                    "work/active/ENG-016/REVIEWS.md": None,
                },
            )
            self.assertEqual(
                validate_cleanup_diff(root, base, candidate, head, work_id="ENG-016"), []
            )
            git_commit(root, "unpermitted", {"README.md": "changed\n"})
            bad_head = subprocess.check_output(
                ["git", "rev-parse", "HEAD"], cwd=root, text=True
            ).strip()
            self.assertTrue(validate_cleanup_diff(root, base, candidate, bad_head))

    def test_cleanup_rejects_multiple_or_wrong_active_work_items(self) -> None:
        with temporary_git_repository() as root:
            base = git_commit(root, "base", {"README.md": "base\n"})
            candidate = git_commit(
                root,
                "candidate",
                {
                    "work/active/ENG-016/WORK.md": "work 16\n",
                    "work/active/ENG-016/REVIEWS.md": "reviews 16\n",
                    "work/active/ENG-017/WORK.md": "work 17\n",
                    "work/active/ENG-017/REVIEWS.md": "reviews 17\n",
                },
            )
            head = git_commit(
                root,
                "cleanup both",
                {
                    "work/active/ENG-016/WORK.md": None,
                    "work/active/ENG-016/REVIEWS.md": None,
                    "work/active/ENG-017/WORK.md": None,
                    "work/active/ENG-017/REVIEWS.md": None,
                },
            )
            errors = validate_cleanup_diff(root, base, candidate, head, work_id="ENG-016")
            self.assertTrue(any("exactly one reviewed work item" in error for error in errors))

        with temporary_git_repository() as root:
            base = git_commit(root, "base", {"README.md": "base\n"})
            candidate = git_commit(
                root,
                "candidate",
                {
                    "work/active/ENG-016/WORK.md": "work 16\n",
                    "work/active/ENG-016/REVIEWS.md": "reviews 16\n",
                    "work/active/ENG-017/WORK.md": "work 17\n",
                    "work/active/ENG-017/REVIEWS.md": "reviews 17\n",
                },
            )
            head = git_commit(
                root,
                "wrong cleanup",
                {
                    "work/active/ENG-017/WORK.md": None,
                    "work/active/ENG-017/REVIEWS.md": None,
                },
            )
            errors = validate_cleanup_diff(root, base, candidate, head, work_id="ENG-016")
            self.assertTrue(any("does not delete the reviewed work item" in error for error in errors))

    def test_nested_phase_archive_is_rejected(self) -> None:
        with temporary_repository() as root:
            (root / "work" / "phases" / "archive").mkdir(parents=True)
            errors = validate_repository(root)
            self.assertTrue(any("work/phases/archive" in error for error in errors))


def temporary_repository():
    return _TemporaryRepository()


class _TemporaryRepository:
    def __enter__(self) -> Path:
        self._directory = tempfile.TemporaryDirectory()
        self.root = Path(self._directory.name)
        create_fixture(self.root)
        return self.root

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self._directory.cleanup()


def create_fixture(root: Path) -> None:
    for relative in (
        "docs/engineering/ai-workflow.md",
        "scripts/validate_ai_workflow.py",
        "scripts/test_validate_ai_workflow.py",
        "work/README.md",
        "work/templates/PHASE.md",
        "work/templates/REVIEWS.md",
        "work/templates/WORK.md",
    ):
        path = root / relative
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(CONTRACT_VERSION + "\n", encoding="utf-8")

    issue = {"title": "Fixture issue", "body": "Fixture body"}
    work_id = "ENG-016"
    base_sha = "a" * 40
    approval_url = "https://github.com/example/campus-lms/issues/16#issuecomment-1"
    normative = "\n".join(
        (
            "### Fixture normative plan",
            "- This fixture tests the repository contract.",
        )
    )
    header = f"""# {work_id} — Fixture

Status: `IMPLEMENTING`

Issue: `{work_id}` / GitHub issue `#16`

Issue URL: <https://github.com/example/campus-lms/issues/16>

Issue specification digest:
`{issue_digest(issue)}`

Plan author actor label: `Fixture planner`

Plan author session label: `Fixture planner session`

Plan revision: `4`

Plan hash: `PLACEHOLDER`

Risk: `CRITICAL`

Release required: `false`

Human plan-approval comment: `{approval_url}`

Implementation author actor label: `Fixture implementer`

Implementation author session label: `Fixture session`

Base Git revision:
`{base_sha}`

Candidate Git commit SHA: `NONE`

<!-- PLAN-NORMATIVE-BEGIN -->

{normative}

<!-- PLAN-NORMATIVE-END -->

- Material deviations: `NONE`
"""
    plan_hash = plan_hash_from_text(header)
    work_text = header.replace("Plan hash: `PLACEHOLDER`", f"Plan hash: `{plan_hash}`")
    active = root / "work" / "active" / work_id
    active.mkdir(parents=True)
    (active / "WORK.md").write_text(work_text, encoding="utf-8")
    (active / "REVIEWS.md").write_text(
        "# ENG-016 reviews\n\n## Review 1\n\n"
        + plan_review_record(root)
        + "\n",
        encoding="utf-8",
    )


def temporary_git_repository():
    return _TemporaryGitRepository()


class _TemporaryGitRepository:
    def __enter__(self) -> Path:
        self._directory = tempfile.TemporaryDirectory()
        self.root = Path(self._directory.name)
        subprocess.run(["git", "init", "-q"], cwd=self.root, check=True)
        subprocess.run(
            ["git", "config", "user.email", "test@example.invalid"],
            cwd=self.root,
            check=True,
        )
        subprocess.run(
            ["git", "config", "user.name", "Workflow Test"],
            cwd=self.root,
            check=True,
        )
        return self.root

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self._directory.cleanup()


def git_commit(root: Path, message: str, files: dict[str, str | None]) -> str:
    for relative, content in files.items():
        path = root / relative
        if content is None:
            path.unlink()
        else:
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(content, encoding="utf-8")
    subprocess.run(["git", "add", "-A"], cwd=root, check=True)
    subprocess.run(["git", "commit", "-qm", message], cwd=root, check=True)
    return subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=root, text=True).strip()


if __name__ == "__main__":
    unittest.main()
