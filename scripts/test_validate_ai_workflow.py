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
    GIT_SHA_RE,
    issue_digest,
    normalize_document,
    plan_hash_from_text,
    validate_approval_comment,
    validate_cleanup_diff,
    validate_repository,
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
            self.assertEqual(validate_cleanup_diff(root, base, candidate, head), [])
            git_commit(root, "unpermitted", {"README.md": "changed\n"})
            bad_head = subprocess.check_output(
                ["git", "rev-parse", "HEAD"], cwd=root, text=True
            ).strip()
            self.assertTrue(validate_cleanup_diff(root, base, candidate, bad_head))


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

Issue specification digest:
`{issue_digest(issue)}`

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
        "## Review 1\n\n- Review ID: `ENG-016-PLAN-REVIEW-001`\n- Review type: `PLAN`\n- Verdict: `APPROVED`\n",
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
