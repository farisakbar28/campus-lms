"""Focused tests for the campus-lms-work/v1 repository workflow validator."""

from __future__ import annotations

import re
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from scripts.validate_ai_workflow import (
    BOOTSTRAP_PLAN_REVIEW_BLOCK_SHA256,
    BOOTSTRAP_PLAN_REVIEW_INDEX,
    BOOTSTRAP_PLAN_REVIEW_WORK,
    CONTRACT_VERSION,
    PLAN_BEGIN,
    PLAN_END,
    _matches_bootstrap_plan_review_identity,
    _plan_review_gate_errors,
    issue_digest,
    normalize_document,
    plan_hash_from_text,
    validate_approval_comment,
    validate_cleanup_diff,
    validate_repository,
)

HISTORICAL_REVIEW_5_SHA256 = (
    "828ef0618dd74baf2edc6262e273ebae0e124ea5a9dec7c46c24c221502f9e58"
)
HISTORICAL_REVIEW_5_WORK = {
    "work_id": "ENG-016",
    "issue_number": 16,
    "issue_url": "https://github.com/farisakbar28/campus-lms/issues/16",
    "plan_revision": "4",
    "plan_hash": "sha256:d6c0464e030f68eb2a9f1089c233e091725a306e9af5d3a7caa0bffa9c7cf934",
    "issue_digest": "sha256:c1cf78cad608524cec118bc469fc223f2cf3d556f37fef0d17c91fc00850ccf4",
    "plan_author_actor": "Codex planner",
    "plan_author_session": "ENG-016-plan-r4-2026-09-06",
}


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


def append_review_block(root: Path, body: str, *, review_number: str = "9") -> None:
    reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
    reviews.write_text(
        reviews.read_text(encoding="utf-8").rstrip()
        + f"\n\n## Review {review_number}\n\n"
        + body
        + "\n",
        encoding="utf-8",
    )


def replace_work_identity_field(root: Path, label: str, value: str) -> None:
    work = root / "work" / "active" / "ENG-016" / "WORK.md"
    lines = work.read_text(encoding="utf-8").splitlines()
    prefix = f"{label}:"
    for index, line in enumerate(lines):
        if line.startswith(prefix):
            lines[index] = f"{label}: `{value}`"
            break
    else:
        raise AssertionError(f"missing fixture field {label}")
    work.write_text("\n".join(lines) + "\n", encoding="utf-8")


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


def populated_review_template(
    root: Path,
    review_type: str,
    *,
    candidate: str = "b" * 40,
    omitted: str | None = None,
) -> str:
    from scripts.validate_ai_workflow import parse_work

    template = (
        Path(__file__).resolve().parents[1] / "work" / "templates" / "REVIEWS.md"
    ).read_text(encoding="utf-8")
    section = "Plan" if review_type == "PLAN" else "Implementation"
    match = re.search(
        rf"## {section} review record\n\n```text\n(?P<record>.*?)\n```",
        template,
        re.DOTALL,
    )
    if match is None:
        raise AssertionError(f"missing {section} review template")

    work = parse_work(
        (root / "work" / "active" / "ENG-016" / "WORK.md").read_text(
            encoding="utf-8"
        )
    )
    values = {
        "review_id": f"ENG-016-{review_type}-REVIEW-TEMPLATE",
        "type": review_type,
        "work_item_id": str(work["work_id"]),
        "actor_label": "Template reviewer",
        "session_label": "Template reviewer session",
        "plan_author_actor_label": str(work["plan_author_actor"]),
        "plan_author_session_label": str(work["plan_author_session"]),
        "implementation_author_actor_label": str(work["implementation_actor"]),
        "implementation_author_session_label": str(work["implementation_session"]),
        "fresh_session_attestation": "Fresh independent template review session.",
        "candidate_git_sha": candidate,
        "plan_revision": str(work["plan_revision"]),
        "plan_hash": str(work["plan_hash"]),
        "issue_digest": str(work["issue_digest"]),
        "finding_id": "ENG-016-TEMPLATE-001",
        "severity": "LOW",
        "status": "RESOLVED",
        "summary": "Canonical summary may mention an APPROVED PLAN as narrative.",
        "residual_risk_comment_url": (
            "https://github.com/example/campus-lms/issues/16#issuecomment-99"
        ),
        "verdict": "APPROVED",
    }
    lines = []
    for line in match.group("record").splitlines():
        key, separator, _ = line.partition("=")
        if not separator:
            raise AssertionError(f"noncanonical template line: {line!r}")
        if key == omitted:
            continue
        if key not in values:
            raise AssertionError(f"unpopulated template field: {key}")
        lines.append(f"{key}={values[key]}")
    return "\n".join(lines)


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
                    any(
                        f"missing {field}" in error
                        or (field == "type" and "missing canonical review type" in error)
                        for error in errors
                    ),
                    errors,
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
            self.assertTrue(errors)
            self.assertTrue(any("canonical" in error for error in errors))

    def test_bootstrap_legacy_plan_review_is_exact_and_non_extensible(self) -> None:
        self.assertEqual(BOOTSTRAP_PLAN_REVIEW_BLOCK_SHA256, HISTORICAL_REVIEW_5_SHA256)
        self.assertEqual(BOOTSTRAP_PLAN_REVIEW_WORK, HISTORICAL_REVIEW_5_WORK)
        bootstrap_work = dict(HISTORICAL_REVIEW_5_WORK)
        self.assertTrue(
            _matches_bootstrap_plan_review_identity(
                work=bootstrap_work,
                review_index=BOOTSTRAP_PLAN_REVIEW_INDEX,
                heading="## Review 5",
                block_sha256=HISTORICAL_REVIEW_5_SHA256,
            )
        )
        for change in (
            {"review_index": BOOTSTRAP_PLAN_REVIEW_INDEX + 1},
            {"heading": "## Review 6"},
            {"block_sha256": "0" * 64},
        ):
            arguments = {
                "work": bootstrap_work,
                "review_index": BOOTSTRAP_PLAN_REVIEW_INDEX,
                "heading": "## Review 5",
                "block_sha256": HISTORICAL_REVIEW_5_SHA256,
            }
            arguments.update(change)
            self.assertFalse(_matches_bootstrap_plan_review_identity(**arguments))

        foreign_work = dict(bootstrap_work, work_id="AUTHCTX-001")
        self.assertFalse(
            _matches_bootstrap_plan_review_identity(
                work=foreign_work,
                review_index=BOOTSTRAP_PLAN_REVIEW_INDEX,
                heading="## Review 5",
                block_sha256=HISTORICAL_REVIEW_5_SHA256,
            )
        )

        historical_prefix = "# reviews\n\n" + "\n\n".join(
            f"## Review {number}\n\nhistorical" for number in range(1, 6)
        )
        with patch(
            "scripts.validate_ai_workflow._has_exact_bootstrap_plan_review",
            return_value=True,
        ):
            self.assertEqual(
                _plan_review_gate_errors(historical_prefix, bootstrap_work), []
            )

        with temporary_repository() as root:
            replace_with_plan_review(root)
            self.assertEqual(validate_repository(root), [])

    def test_canonical_plan_and_implementation_reviews_pass_after_bootstrap(self) -> None:
        historical_prefix = "# reviews\n\n" + "\n\n".join(
            f"## Review {number}\n\nhistorical" for number in range(1, 6)
        )
        common = {
            "work_item_id": "ENG-016",
            "plan_revision": "4",
            "plan_hash": HISTORICAL_REVIEW_5_WORK["plan_hash"],
            "issue_digest": HISTORICAL_REVIEW_5_WORK["issue_digest"],
            "actor_label": "Fresh reviewer",
            "session_label": "Fresh reviewer session",
            "fresh_session_attestation": "Fresh independent review session.",
            "verdict": "APPROVED",
        }
        plan = {
            "review_id": "ENG-016-PLAN-REVIEW-006",
            "type": "PLAN",
            **common,
            "plan_author_actor_label": HISTORICAL_REVIEW_5_WORK["plan_author_actor"],
            "plan_author_session_label": HISTORICAL_REVIEW_5_WORK["plan_author_session"],
        }
        implementation = {
            "review_id": "ENG-016-IMPLEMENTATION-REVIEW-001",
            "type": "IMPLEMENTATION",
            **common,
            "implementation_author_actor_label": "Codex implementer",
            "implementation_author_session_label": "Implementation session",
            "candidate_git_sha": "b" * 40,
        }
        for review_number, record in ((6, plan), (6, implementation)):
            with self.subTest(review_type=record["type"]), patch(
                "scripts.validate_ai_workflow._has_exact_bootstrap_plan_review",
                return_value=True,
            ):
                reviews = (
                    historical_prefix
                    + f"\n\n## Review {review_number}\n\n"
                    + "\n".join(f"{key}={value}" for key, value in record.items())
                    + "\n"
                )
                self.assertEqual(
                    _plan_review_gate_errors(reviews, HISTORICAL_REVIEW_5_WORK), []
                )

    def test_later_legacy_plan_reviews_fail_closed_in_all_structural_forms(self) -> None:
        cases = {
            "prose-only": "A later prose-only review cannot provide authority.",
            "review-11-sentence": "This PLAN review is APPROVED.",
            "same-line": "- Review type: `PLAN`\n- Verdict: `APPROVED`",
            "markdown-bold": "- **Review type:** `PLAN`\n- **Verdict:** `APPROVED`",
            "multiline": "- Review type:\n  `PLAN`\n- Verdict:\n  `APPROVED`",
            "continuation-values": (
                "- Review type:\n  `PLAN`\n"
                "- Subject plan revision:\n  `4`\n"
                "- Verdict:\n  `APPROVED`"
            ),
            "malformed-multiline": (
                "- Review type:\n  `PL\n  AN`\n"
                "- Subject plan revision:\n  not-an-integer"
            ),
            "loose-approval": "PLAN APPROVED",
            "loose-split-approval": "PLAN\nAPPROVED",
            "unicode-whitespace": (
                "-\u00a0Review type:\u00a0\n\u00a0 `PLAN`\n"
                "-\u00a0Subject plan revision:\u00a0\n\u00a0 `4`"
            ),
            "zero-width-format-character": (
                "- Review\u200b type: `PLAN`\n- Verdict\u200b: `APPROVED`"
            ),
        }
        for name, record in cases.items():
            with self.subTest(name=name), temporary_repository() as root:
                append_review_block(root, record)
                errors = validate_repository(root)
                self.assertTrue(
                    any("missing canonical review type" in error for error in errors),
                    errors,
                )

    def test_bootstrap_fallback_rejects_a_later_multiline_legacy_record(self) -> None:
        historical_prefix = "# reviews\n\n" + "\n\n".join(
            f"## Review {number}\n\nhistorical" for number in range(1, 6)
        )
        later_record = (
            "## Review 6\n\n"
            "- Review type:\n  `PLAN`\n"
            "- Subject plan revision:\n  `4`\n"
            "- Verdict:\n  `APPROVED`\n"
        )
        with patch(
            "scripts.validate_ai_workflow._has_exact_bootstrap_plan_review",
            return_value=True,
        ):
            errors = _plan_review_gate_errors(
                historical_prefix + "\n\n" + later_record,
                HISTORICAL_REVIEW_5_WORK,
            )
        self.assertTrue(any("missing canonical review type" in error for error in errors), errors)

    def test_post_bootstrap_review_type_must_be_canonical_and_known(self) -> None:
        cases = {
            "missing": "review_id=ENG-016-PLAN-REVIEW-100\nverdict=APPROVED",
            "legacy-alias-only": (
                "review_id=ENG-016-PLAN-REVIEW-100\n"
                "review_type=PLAN\nverdict=APPROVED"
            ),
            "malformed": "type=plan",
            "unknown": "type=SECURITY",
        }
        for name, record in cases.items():
            with self.subTest(name=name), temporary_repository() as root:
                append_review_block(root, record)
                errors = validate_repository(root)
                self.assertTrue(errors)
                if name in {"missing", "legacy-alias-only"}:
                    self.assertTrue(
                        any("missing canonical review type" in error for error in errors),
                        errors,
                    )
                else:
                    self.assertTrue(
                        any("unknown canonical review type" in error for error in errors),
                        errors,
                    )

    def test_later_legacy_plan_bindings_never_extend_to_stale_or_future_work(self) -> None:
        for overrides in (
            {},
            {"Review ID": "ENG-016-PLAN-REVIEW-100"},
            {"Subject work item": "ENG-017"},
            {"Subject plan revision": "5"},
            {
                "Exact reviewed plan checkpoint SHA": "sha256:" + "0" * 64,
                "Independently recomputed normative SHA": "sha256:" + "0" * 64,
            },
        ):
            with self.subTest(overrides=overrides), temporary_repository() as root:
                append_review_block(root, legacy_plan_review_record(root, **overrides))
                errors = validate_repository(root)
                self.assertTrue(
                    any("missing canonical review type" in error for error in errors),
                    errors,
                )

        with temporary_repository() as root:
            active = root / "work" / "active"
            future = active / "AUTHCTX-001"
            (active / "ENG-016").rename(future)
            for name in ("WORK.md", "REVIEWS.md"):
                path = future / name
                path.write_text(
                    path.read_text(encoding="utf-8").replace("ENG-016", "AUTHCTX-001"),
                    encoding="utf-8",
                )
            reviews = future / "REVIEWS.md"
            reviews.write_text(
                reviews.read_text(encoding="utf-8").rstrip()
                + "\n\n## Review 2\n\n- Review type: `PLAN`\n- Verdict: `APPROVED`\n",
                encoding="utf-8",
            )
            errors = validate_repository(root)
            self.assertTrue(
                any("missing canonical review type" in error for error in errors),
                errors,
            )

    def test_canonical_review_narrative_can_mention_plan_and_approved(self) -> None:
        with temporary_repository() as root:
            candidate = "b" * 40
            bind_fixture_work(root, candidate=candidate)
            append_exact_implementation_review(root, candidate=candidate)
            append_finding(
                root,
                finding_id="F-NARRATIVE",
                severity="LOW",
                status="RESOLVED",
                summary="The approved PLAN wording is quoted only as finding context.",
            )
            self.assertEqual(validate_repository(root), [])

    def test_plan_author_identity_is_required_before_independence(self) -> None:
        cases = (
            ("Plan author actor label", "", "plan author actor label"),
            ("Plan author session label", "", "plan author session label"),
            ("Plan author actor label", " ", "whitespace-only"),
            ("Plan author session label", "   ", "whitespace-only"),
            ("Plan author actor label", "None", "null-like"),
            ("Plan author session label", "null", "null-like"),
        )
        for label, value, expected in cases:
            with self.subTest(label=label, value=value), temporary_repository() as root:
                replace_work_identity_field(root, label, value)
                errors = validate_repository(root)
                self.assertTrue(errors)
                self.assertTrue(
                    any(
                        "plan author" in error and expected in error
                        for error in errors
                    ),
                    errors,
                )

        with temporary_repository() as root:
            work = root / "work" / "active" / "ENG-016" / "WORK.md"
            text = work.read_text(encoding="utf-8")
            text = text.replace("Plan author actor label: `Fixture planner`\n\n", "")
            text = text.replace("Plan author session label: `Fixture planner session`\n\n", "")
            work.write_text(text, encoding="utf-8")
            replace_with_plan_review(root)
            errors = validate_repository(root)
            self.assertTrue(any("missing" in error and "plan author" in error for error in errors))

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
            self.assertTrue(any("reviewer actor equals author actor" in error for error in errors))
            self.assertTrue(any("reviewer session equals author session" in error for error in errors))

        with temporary_repository() as root:
            replace_with_plan_review(root)
            self.assertEqual(validate_repository(root), [])

    def test_invalid_state_and_equal_review_identity_are_rejected(self) -> None:
        with temporary_repository() as root:
            work = root / "work" / "active" / "ENG-016" / "WORK.md"
            text = work.read_text(encoding="utf-8").replace(
                "Status: `IMPLEMENTING`", "Status: `DONE`"
            )
            work.write_text(text, encoding="utf-8")
            append_exact_implementation_review(
                root,
                candidate="b" * 40,
                actor_label="Fixture implementer",
                implementation_author_actor_label="Fixture implementer",
            )
            errors = validate_repository(root)
            self.assertTrue(any("post-merge/DONE" in error for error in errors))
            self.assertTrue(any("reviewer actor equals" in error for error in errors))

    def test_finding_severity_status_decision_matrix(self) -> None:
        cases = (
            ("CRITICAL", "OPEN", True),
            ("CRITICAL", "RESOLVED", False),
            ("CRITICAL", "ACCEPTED_RESIDUAL_RISK", True),
            ("HIGH", "OPEN", True),
            ("HIGH", "RESOLVED", False),
            ("HIGH", "ACCEPTED_RESIDUAL_RISK", True),
            ("MEDIUM", "OPEN", True),
            ("MEDIUM", "RESOLVED", False),
            ("LOW", "OPEN", False),
            ("LOW", "RESOLVED", False),
        )
        for severity, status, blocked in cases:
            with self.subTest(severity=severity, status=status), temporary_repository() as root:
                append_finding(root, finding_id="F-MATRIX", severity=severity, status=status)
                errors = validate_repository(root)
                finding_errors = [error for error in errors if "F-MATRIX" in error]
                self.assertEqual(bool(finding_errors), blocked, errors)

        for severity, status in (("UNKNOWN", "OPEN"), ("MEDIUM", "UNKNOWN")):
            with self.subTest(severity=severity, status=status), temporary_repository() as root:
                append_finding(root, finding_id="F-UNKNOWN", severity=severity, status=status)
                errors = validate_repository(root)
                self.assertTrue(any("invalid finding" in error for error in errors), errors)

        with temporary_repository() as root:
            append_finding(
                root,
                finding_id="F-LOW-RISK",
                severity="LOW",
                status="ACCEPTED_RESIDUAL_RISK",
            )
            errors = validate_repository(root)
            self.assertTrue(any("missing residual-risk comment URL" in error for error in errors))

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
            self.assertTrue(any("missing canonical review type" in error for error in errors))

    def test_review_templates_round_trip_through_canonical_schema(self) -> None:
        template = (
            Path(__file__).resolve().parents[1]
            / "work"
            / "templates"
            / "REVIEWS.md"
        ).read_text(encoding="utf-8")
        self.assertNotIn("ENG-016", template)

        with temporary_repository() as root:
            reviews = root / "work" / "active" / "ENG-016" / "REVIEWS.md"
            reviews.write_text(
                "# reviews\n\n## Review 1\n\n"
                + populated_review_template(root, "PLAN")
                + "\n",
                encoding="utf-8",
            )
            self.assertEqual(validate_repository(root), [])

        with temporary_repository() as root:
            candidate = "b" * 40
            bind_fixture_work(root, candidate=candidate, status="READY_FOR_PR")
            append_review_block(
                root,
                populated_review_template(root, "IMPLEMENTATION", candidate=candidate),
                review_number="2",
            )
            self.assertEqual(validate_repository(root), [])

        with temporary_repository() as root:
            candidate = "b" * 40
            bind_fixture_work(root, candidate=candidate, status="READY_FOR_PR")
            append_review_block(
                root,
                populated_review_template(
                    root,
                    "IMPLEMENTATION",
                    candidate=candidate,
                    omitted="work_item_id",
                ),
                review_number="2",
            )
            errors = validate_repository(root)
            self.assertTrue(any("missing work_item_id" in error for error in errors), errors)

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

            foreign_url = "https://github.com/foreign/repository/issues/16#issuecomment-42"
            append_finding(
                root,
                finding_id="ENG-016-IMPL-FOREIGN",
                severity="MEDIUM",
                status="ACCEPTED_RESIDUAL_RISK",
                residual_url=foreign_url,
            )
            foreign_comment = dict(comment)
            foreign_comment["html_url"] = foreign_url
            foreign_comment["url"] = (
                "https://api.github.com/repos/foreign/repository/issues/comments/42"
            )
            foreign_comment["body"] = foreign_comment["body"].replace(
                "ENG-016-IMPL-TEST", "ENG-016-IMPL-FOREIGN"
            )
            errors = validate_repository(root, residual_comments=[foreign_comment])
            self.assertTrue(any("not on the canonical Issue" in error for error in errors))

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

    def test_active_cleanup_lifecycle_is_hermetic(self) -> None:
        with temporary_repository() as root:
            self.assertEqual(validate_repository(root), [])

            active_item = root / "work" / "active" / "ENG-016"
            (active_item / "WORK.md").unlink()
            (active_item / "REVIEWS.md").unlink()
            active_item.rmdir()

            self.assertEqual(validate_repository(root), [])
            command = subprocess.run(
                [
                    sys.executable,
                    str(Path(__file__).with_name("validate_ai_workflow.py")),
                    "--repo-root",
                    str(root),
                ],
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(command.returncode, 0, command.stderr)
            self.assertIn("AI workflow validation: PASS", command.stdout)

            (root / "work" / "active").rmdir()
            self.assertEqual(validate_repository(root), [])

    def test_unrelated_future_active_work_item_is_independent(self) -> None:
        with temporary_repository() as root:
            active = root / "work" / "active"
            future = active / "AUTHCTX-001"
            (active / "ENG-016").rename(future)
            for name in ("WORK.md", "REVIEWS.md"):
                path = future / name
                path.write_text(
                    path.read_text(encoding="utf-8").replace("ENG-016", "AUTHCTX-001"),
                    encoding="utf-8",
                )
            self.assertEqual(validate_repository(root), [])


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
