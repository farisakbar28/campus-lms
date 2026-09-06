#!/usr/bin/env python3
"""Validate the machine-checkable part of the repository workflow contract.

This module deliberately validates repository artifacts only.  It does not
claim to prove human intent, reviewer independence, branch protection, or
live production state.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import struct
import subprocess
import sys
from pathlib import Path
from typing import Any, Iterable, Mapping


CONTRACT_VERSION = "campus-lms-work/v1"
PLAN_HASH_ALGORITHM = "campus-lms-plan-sha256-v1"
ISSUE_HASH_ALGORITHM = "campus-lms-issue-sha256-v1"

PLAN_BEGIN = "<!-- PLAN-NORMATIVE-BEGIN -->"
PLAN_END = "<!-- PLAN-NORMATIVE-END -->"
PHASE_BEGIN = "<!-- PHASE-NORMATIVE-BEGIN -->"
PHASE_END = "<!-- PHASE-NORMATIVE-END -->"

STATES = (
    "DRAFT",
    "PLAN_REVIEW",
    "APPROVED",
    "IMPLEMENTING",
    "IMPLEMENTATION_REVIEW",
    "READY_FOR_PR",
    "PR_VALIDATION",
    "READY_TO_MERGE",
)
RISKS = ("NORMAL", "HIGH", "CRITICAL")
FINDING_SEVERITIES = ("CRITICAL", "HIGH", "MEDIUM", "LOW")
FINDING_STATUSES = ("OPEN", "RESOLVED", "ACCEPTED_RESIDUAL_RISK")
REVIEW_VERDICTS = ("CHANGES_REQUIRED", "APPROVED")
SHA256_RE = re.compile(r"^sha256:[0-9a-f]{64}$")
GIT_SHA_RE = re.compile(r"^[0-9a-f]{40}$")
WORK_ID_RE = re.compile(r"^[A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*-[0-9]{3,}$")
WEB_ISSUE_URL_RE = re.compile(
    r"^https://github\.com/(?P<owner>[^/\s]+)/(?P<repo>[^/\s]+)"
    r"/issues/(?P<number>[0-9]+)$"
)
WEB_COMMENT_URL_RE = re.compile(
    r"^https://github\.com/(?P<owner>[^/\s]+)/(?P<repo>[^/\s]+)"
    r"/issues/(?P<number>[0-9]+)#issuecomment-(?P<comment>[0-9]+)$"
)

IMPLEMENTATION_REVIEW_FIELDS = (
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
FINDING_FIELDS = {
    "finding_id",
    "severity",
    "status",
    "summary",
    "residual_risk_comment_url",
}
REQUIRED_FINDING_FIELDS = {
    "finding_id",
    "severity",
    "status",
    "summary",
}
FORBIDDEN_WORK_DIRECTORY_NAMES = {
    "completed",
    "archive",
    "archives",
    "archived",
}

REQUIRED_DURABLE_FILES = (
    "docs/engineering/ai-workflow.md",
    "scripts/validate_ai_workflow.py",
    "scripts/test_validate_ai_workflow.py",
    "work/README.md",
    "work/templates/PHASE.md",
    "work/templates/REVIEWS.md",
    "work/templates/WORK.md",
)


class WorkflowValidationError(ValueError):
    """Raised for malformed workflow input supplied to a helper."""


def _normalise_lines(text: str) -> list[str]:
    """Apply the workflow LF and trailing-whitespace line normalization."""

    text = text.replace("\r\n", "\n").replace("\r", "\n")
    return [line.rstrip(" \t") for line in text.split("\n")]


def normalize_document(text: str) -> bytes:
    """Normalize a complete document using the plan-hash rules."""

    lines = _normalise_lines(text)
    while lines and lines[0] == "":
        lines.pop(0)
    while lines and lines[-1] == "":
        lines.pop()
    if not lines:
        return b""
    return ("\n".join(lines) + "\n").encode("utf-8")


def _marked_section(text: str, begin: str, end: str) -> bytes:
    lines = _normalise_lines(text)
    if lines.count(begin) != 1 or lines.count(end) != 1:
        raise WorkflowValidationError(
            f"expected exactly one {begin!r} and {end!r} marker"
        )
    begin_index = lines.index(begin)
    end_index = lines.index(end)
    if begin_index >= end_index:
        raise WorkflowValidationError("normative markers are out of order")
    content = lines[begin_index + 1 : end_index]
    while content and content[0] == "":
        content.pop(0)
    while content and content[-1] == "":
        content.pop()
    if not content:
        raise WorkflowValidationError("normative section is empty")
    return ("\n".join(content) + "\n").encode("utf-8")


def plan_hash_from_text(text: str) -> str:
    """Return the deterministic plan hash for a WORK.md document."""

    return "sha256:" + hashlib.sha256(
        _marked_section(text, PLAN_BEGIN, PLAN_END)
    ).hexdigest()


def phase_hash_from_text(text: str) -> str:
    """Return the deterministic phase hash for a PHASE.md document."""

    return "sha256:" + hashlib.sha256(
        _marked_section(text, PHASE_BEGIN, PHASE_END)
    ).hexdigest()


def issue_digest(issue: Mapping[str, Any]) -> str:
    """Compute the canonical title/body digest, ignoring comments."""

    title = issue.get("title")
    body = issue.get("body")
    if not isinstance(title, str) or not isinstance(body, str):
        raise WorkflowValidationError("issue JSON must contain string title/body")
    title_bytes = normalize_document(title)
    body_bytes = normalize_document(body)
    payload = (
        ISSUE_HASH_ALGORITHM.encode("ascii")
        + b"\0"
        + struct.pack(">Q", len(title_bytes))
        + title_bytes
        + struct.pack(">Q", len(body_bytes))
        + body_bytes
    )
    return "sha256:" + hashlib.sha256(payload).hexdigest()


def _read_utf8(path: Path) -> str:
    try:
        return path.read_bytes().decode("utf-8")
    except UnicodeDecodeError as exc:
        raise WorkflowValidationError(f"{path}: invalid UTF-8") from exc


def _field(text: str, label: str, *, multiline: bool = False) -> str | None:
    if multiline:
        pattern = (
            rf"(?m)^(?:-\s*)?{re.escape(label)}:[ \t]*\n"
            rf"[ \t]*`([^`\n]+)`[ \t]*\.?[ \t]*$"
        )
    else:
        pattern = rf"(?m)^(?:-\s*)?{re.escape(label)}:\s*`([^`\n]+)`\s*$"
    match = re.search(pattern, text)
    return match.group(1) if match else None


def _angle_field(text: str, label: str) -> str | None:
    match = re.search(
        rf"(?m)^(?:-\s*)?{re.escape(label)}:[ \t]*<([^<>\n]+)>[ \t]*$",
        text,
    )
    return match.group(1) if match else None


def _web_issue_parts(value: Any) -> tuple[str, str, int] | None:
    if not isinstance(value, str):
        return None
    match = WEB_ISSUE_URL_RE.fullmatch(value)
    if not match:
        return None
    return match.group("owner"), match.group("repo"), int(match.group("number"))


def _web_comment_parts(value: Any) -> tuple[str, str, int, int] | None:
    if not isinstance(value, str):
        return None
    match = WEB_COMMENT_URL_RE.fullmatch(value)
    if not match:
        return None
    return (
        match.group("owner"),
        match.group("repo"),
        int(match.group("number")),
        int(match.group("comment")),
    )


def _api_repository_url(owner: str, repo: str) -> str:
    return f"https://api.github.com/repos/{owner}/{repo}"


def _api_issue_url(owner: str, repo: str, number: int) -> str:
    return f"{_api_repository_url(owner, repo)}/issues/{number}"


def _key_value_fields(text: str) -> dict[str, list[str]]:
    fields: dict[str, list[str]] = {}
    for line in _normalise_lines(text):
        match = re.fullmatch(r"(?P<key>[a-z][a-z0-9_]*)=(?P<value>.*)", line)
        if match:
            fields.setdefault(match.group("key"), []).append(match.group("value"))
    return fields


def _review_blocks(reviews_text: str) -> list[str]:
    return re.split(r"(?m)(?=^##\s+Review\b)", reviews_text)


def _issue_number(text: str) -> int | None:
    match = re.search(
        r"(?m)^Issue:\s*`(?P<work>[A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*-[0-9]{3,})`"
        r"\s*/\s*GitHub issue\s*`#(?P<number>[0-9]+)`\s*$",
        text,
    )
    return int(match.group("number")) if match else None


def _work_id_from_heading(text: str) -> str | None:
    match = re.search(
        r"(?m)^#\s+(?P<work>[A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*-[0-9]{3,})\s+—",
        text,
    )
    return match.group("work") if match else None


def parse_work(text: str) -> dict[str, str | int | None]:
    """Extract the small mutable/header binding set from WORK.md."""

    values: dict[str, str | int | None] = {
        "work_id": _work_id_from_heading(text),
        "issue_number": _issue_number(text),
        "issue_url": _angle_field(text, "Issue URL"),
        "status": _field(text, "Status"),
        "issue_digest": _field(text, "Issue specification digest", multiline=True),
        "plan_revision": _field(text, "Plan revision"),
        "plan_hash": _field(text, "Plan hash"),
        "risk": _field(text, "Risk"),
        "release_required": _field(text, "Release required"),
        "approval_url": _field(text, "Human plan-approval comment"),
        "implementation_actor": _field(text, "Implementation author actor label"),
        "implementation_session": _field(
            text, "Implementation author session label"
        ),
        "candidate_sha": _field(text, "Candidate Git commit SHA"),
        "base_sha": _field(text, "Base Git revision", multiline=True),
    }
    return values


def _require(value: Any, name: str, errors: list[str]) -> str:
    if value is None or value == "":
        errors.append(f"missing {name}")
        return ""
    return str(value)


def _validate_work_document(
    path: Path,
    work_id_from_directory: str,
    reviews_text: str,
    residual_comments: Iterable[Mapping[str, Any]] = (),
) -> list[str]:
    errors: list[str] = []
    try:
        text = _read_utf8(path)
    except WorkflowValidationError as exc:
        return [str(exc)]

    values = parse_work(text)
    work_id = _require(values["work_id"], f"{path}: work-item heading", errors)
    if work_id and not WORK_ID_RE.fullmatch(work_id):
        errors.append(f"{path}: invalid work-item ID {work_id!r}")
    if work_id and re.fullmatch(r"W-[0-9]+", work_id):
        errors.append(f"{path}: generated W-#### work-item IDs are forbidden")
    if work_id and work_id != work_id_from_directory:
        errors.append(
            f"{path}: heading ID {work_id!r} does not match directory "
            f"{work_id_from_directory!r}"
        )
    if not WORK_ID_RE.fullmatch(work_id_from_directory):
        errors.append(f"{path}: invalid directory work-item ID")

    issue_number = values["issue_number"]
    if issue_number is None:
        errors.append(f"{path}: missing canonical GitHub issue mapping")
    issue_url = _require(values["issue_url"], f"{path}: Issue URL", errors)
    issue_parts = _web_issue_parts(issue_url)
    if issue_url and issue_parts is None:
        errors.append(f"{path}: invalid canonical Issue URL")
    elif issue_parts is not None and issue_number != issue_parts[2]:
        errors.append(f"{path}: Issue URL number does not match canonical mapping")

    status = _require(values["status"], f"{path}: Status", errors)
    if status and status not in STATES:
        errors.append(f"{path}: invalid state {status!r}")
    if status in {"DONE", "MERGED", "POST_MERGE", "RELEASED"}:
        errors.append(f"{path}: post-merge/DONE state is forbidden")

    issue_hash = _require(
        values["issue_digest"], f"{path}: issue specification digest", errors
    )
    if issue_hash and not SHA256_RE.fullmatch(issue_hash):
        errors.append(f"{path}: invalid issue specification digest")

    revision = _require(values["plan_revision"], f"{path}: plan revision", errors)
    if revision:
        try:
            if int(revision) < 1 or str(int(revision)) != revision:
                raise ValueError
        except ValueError:
            errors.append(f"{path}: plan revision must be a positive integer")

    plan_hash = _require(values["plan_hash"], f"{path}: plan hash", errors)
    if plan_hash and not SHA256_RE.fullmatch(plan_hash):
        errors.append(f"{path}: invalid plan hash")
    else:
        try:
            computed = plan_hash_from_text(text)
            if plan_hash and computed != plan_hash:
                errors.append(
                    f"{path}: plan hash mismatch (declared {plan_hash}, computed {computed})"
                )
        except WorkflowValidationError as exc:
            errors.append(f"{path}: {exc}")

    base_sha = _require(values["base_sha"], f"{path}: base Git revision", errors)
    if base_sha and not GIT_SHA_RE.fullmatch(base_sha):
        errors.append(f"{path}: invalid base Git revision")

    risk = _require(values["risk"], f"{path}: risk", errors)
    if risk and risk not in RISKS:
        errors.append(f"{path}: invalid risk {risk!r}")
    release_required = _require(
        values["release_required"], f"{path}: release required", errors
    )
    if release_required not in {"true", "false"}:
        errors.append(f"{path}: release required must be true or false")

    approval_url = _require(
        values["approval_url"], f"{path}: human plan-approval comment", errors
    )
    if approval_url != "NONE":
        approval_parts = _web_comment_parts(approval_url)
        if approval_parts is None:
            errors.append(f"{path}: invalid human plan-approval comment URL")
        elif issue_parts is None or approval_parts[:3] != issue_parts:
            errors.append(
                f"{path}: human plan-approval comment URL is not on the canonical Issue"
            )

    implementation_actor = _require(
        values["implementation_actor"],
        f"{path}: implementation author actor label",
        errors,
    )
    implementation_session = _require(
        values["implementation_session"],
        f"{path}: implementation author session label",
        errors,
    )
    candidate_sha = _require(
        values["candidate_sha"], f"{path}: candidate Git SHA", errors
    )
    if candidate_sha != "NONE" and not GIT_SHA_RE.fullmatch(candidate_sha):
        errors.append(f"{path}: invalid candidate Git SHA")

    review_ids = re.findall(r"(?m)^-?\s*Review ID:\s*`?([^`\n]+?)`?\s*$", reviews_text)
    if status not in {"DRAFT"} and not review_ids:
        errors.append(f"{path}: active work has no review record")

    approval_states = {
        "APPROVED",
        "IMPLEMENTING",
        "IMPLEMENTATION_REVIEW",
        "READY_FOR_PR",
        "PR_VALIDATION",
        "READY_TO_MERGE",
    }
    if status in approval_states and approval_url == "NONE":
        errors.append(f"{path}: {status} requires a plan-approval URL")
    if status in approval_states and not _has_approved_plan_review(reviews_text):
        errors.append(f"{path}: {status} requires an APPROVED plan review")

    implementation_states = {
        "IMPLEMENTING",
        "IMPLEMENTATION_REVIEW",
        "READY_FOR_PR",
        "PR_VALIDATION",
        "READY_TO_MERGE",
    }
    if status in implementation_states:
        if implementation_actor == "NONE" or implementation_session == "NONE":
            errors.append(f"{path}: {status} requires implementation author labels")
    if status in {
        "IMPLEMENTATION_REVIEW",
        "READY_FOR_PR",
        "PR_VALIDATION",
        "READY_TO_MERGE",
    } and candidate_sha == "NONE":
        errors.append(f"{path}: {status} requires a candidate Git SHA")

    work_bindings = {
        **values,
        "work_id": work_id,
        "issue_number": issue_number,
        "issue_url": issue_url,
    }
    if status in {"READY_FOR_PR", "PR_VALIDATION", "READY_TO_MERGE"}:
        if not _has_approved_implementation_review(reviews_text, work_bindings):
            errors.append(
                f"{path}: {status} requires an APPROVED implementation review "
                f"for {candidate_sha}"
            )

    errors.extend(
        _implementation_review_errors_for_document(
            path,
            reviews_text,
            work_bindings,
        )
    )
    errors.extend(
        _finding_gate_errors(
            path,
            reviews_text,
            work_bindings,
            residual_comments,
        )
    )
    errors.extend(_review_relationship_errors(path, reviews_text))
    return errors


def _is_implementation_review_block(block: str) -> bool:
    fields = _key_value_fields(block)
    if any(
        value == "IMPLEMENTATION"
        for key in ("type", "review_type")
        for value in fields.get(key, [])
    ):
        return True
    return bool(
        re.search(
            r"(?m)^-\s*Review type:\s*`?IMPLEMENTATION`?\s*$", block
        )
    )


def _implementation_review_errors(
    block: str,
    work: Mapping[str, str | int | None],
    *,
    require_current_binding: bool = False,
    require_approved: bool = False,
) -> list[str]:
    fields = _key_value_fields(block)
    errors: list[str] = []
    values: dict[str, str] = {}
    for key in IMPLEMENTATION_REVIEW_FIELDS:
        occurrences = fields.get(key, [])
        if not occurrences:
            errors.append(f"implementation review is missing {key}")
            continue
        if len(occurrences) != 1:
            errors.append(f"implementation review repeats {key}")
            continue
        value = occurrences[0]
        if not value.strip():
            errors.append(f"implementation review has an empty {key}")
            continue
        values[key] = value

    if values.get("type") != "IMPLEMENTATION":
        errors.append("implementation review type must be IMPLEMENTATION")

    revision = values.get("plan_revision")
    if revision and (
        not revision.isdigit() or int(revision) < 1 or str(int(revision)) != revision
    ):
        errors.append("implementation review has an invalid plan_revision")
    for key in ("plan_hash", "issue_digest"):
        value = values.get(key)
        if value and not SHA256_RE.fullmatch(value):
            errors.append(f"implementation review has an invalid {key}")
    candidate = values.get("candidate_git_sha")
    if candidate and not GIT_SHA_RE.fullmatch(candidate):
        errors.append("implementation review has an invalid candidate_git_sha")
    if values.get("verdict") and values["verdict"] not in REVIEW_VERDICTS:
        errors.append("implementation review has an invalid verdict")

    for key in (
        "actor_label",
        "session_label",
        "implementation_author_actor_label",
        "implementation_author_session_label",
        "fresh_session_attestation",
    ):
        if values.get(key) == "NONE":
            errors.append(f"implementation review has an invalid {key}")

    if (
        values.get("actor_label")
        and values.get("implementation_author_actor_label")
        and values["actor_label"] == values["implementation_author_actor_label"]
    ):
        errors.append("implementation review reviewer actor equals author actor")
    if (
        values.get("session_label")
        and values.get("implementation_author_session_label")
        and values["session_label"] == values["implementation_author_session_label"]
    ):
        errors.append("implementation review reviewer session equals author session")

    aliases = (
        ("review_type", "type"),
        ("reviewer_actor_label", "actor_label"),
        ("fresh_reviewer_session_label", "session_label"),
        ("implementation_candidate_sha", "candidate_git_sha"),
    )
    for alias, canonical in aliases:
        occurrences = fields.get(alias, [])
        if not occurrences:
            continue
        if len(occurrences) != 1:
            errors.append(f"implementation review repeats {alias}")
        elif canonical in values and occurrences[0] != values[canonical]:
            errors.append(
                f"implementation review {alias} does not match {canonical}"
            )

    if require_approved and values.get("verdict") != "APPROVED":
        errors.append("implementation review verdict is not APPROVED")

    if require_current_binding:
        expected = {
            "work_item_id": str(work.get("work_id")),
            "plan_revision": str(work.get("plan_revision")),
            "plan_hash": str(work.get("plan_hash")),
            "issue_digest": str(work.get("issue_digest")),
            "candidate_git_sha": str(work.get("candidate_sha")),
            "implementation_author_actor_label": str(work.get("implementation_actor")),
            "implementation_author_session_label": str(work.get("implementation_session")),
        }
        for key, expected_value in expected.items():
            if values.get(key) != expected_value:
                errors.append(f"implementation review {key} does not match current WORK.md")

    return errors


def _implementation_review_errors_for_document(
    path: Path,
    reviews_text: str,
    work: Mapping[str, str | int | None],
) -> list[str]:
    errors: list[str] = []
    for block in _review_blocks(reviews_text):
        if _is_implementation_review_block(block):
            errors.extend(
                f"{path}: {error}"
                for error in _implementation_review_errors(block, work)
            )
    return errors


def _has_approved_implementation_review(
    reviews_text: str, work: Mapping[str, str | int | None]
) -> bool:
    matches = 0
    for block in _review_blocks(reviews_text):
        if not _is_implementation_review_block(block):
            continue
        if not _implementation_review_errors(
            block,
            work,
            require_current_binding=True,
            require_approved=True,
        ):
            matches += 1
    return matches == 1


def _has_approved_plan_review(reviews_text: str) -> bool:
    blocks = _review_blocks(reviews_text)
    for block in blocks:
        if not re.search(r"(?m)^-\s*Review type:\s*`?PLAN`?\s*$", block):
            continue
        if re.search(r"(?m)^-\s*Verdict:\s*`?APPROVED`?\s*$", block):
            return True
        if re.search(r"(?m)^`APPROVED`\s*$", block):
            return True
    return False


def _review_relationship_errors(path: Path, reviews_text: str) -> list[str]:
    """Check the relational actor/session rules when review records expose them."""

    errors: list[str] = []
    blocks = re.split(r"(?m)(?=^##\s+Review\b)", reviews_text)
    for block in blocks:
        review_type = _block_field(block, "Review type")
        if review_type not in {"PLAN", "IMPLEMENTATION"}:
            continue
        actor = _block_field(block, "Actor")
        session = _block_field(block, "Session label")
        if review_type == "PLAN":
            author_actor = _block_field(block, "Plan author actor label")
            author_session = _block_field(block, "Plan author session label")
        else:
            author_actor = _block_field(block, "Implementation author actor label")
            author_session = _block_field(block, "Implementation author session label")
        if actor and author_actor and actor == author_actor:
            errors.append(f"{path}: {review_type} reviewer actor equals author actor")
        if session and author_session and session == author_session:
            errors.append(f"{path}: {review_type} reviewer session equals author session")
    return errors


def _block_field(block: str, label: str) -> str | None:
    match = re.search(
        rf"(?m)^-\s*{re.escape(label)}:\s*`?([^`\n]+?)`?\s*$", block
    )
    return match.group(1) if match else None


def _structured_finding_records(
    reviews_text: str,
) -> tuple[list[dict[str, str]], list[str]]:
    records: list[dict[str, str]] = []
    errors: list[str] = []

    def finish(current: dict[str, list[str]] | None) -> None:
        if current is None:
            return
        record: dict[str, str] = {}
        for key in REQUIRED_FINDING_FIELDS:
            occurrences = current.get(key, [])
            if not occurrences:
                errors.append(f"finding record is missing {key}")
            elif len(occurrences) != 1:
                errors.append(f"finding record repeats {key}")
            elif not occurrences[0].strip():
                errors.append(f"finding record has an empty {key}")
            else:
                record[key] = occurrences[0]
        if "residual_risk_comment_url" in current:
            occurrences = current["residual_risk_comment_url"]
            if len(occurrences) != 1:
                errors.append("finding record repeats residual_risk_comment_url")
            elif not occurrences[0].strip():
                errors.append("finding record has an empty residual_risk_comment_url")
            else:
                record["residual_risk_comment_url"] = occurrences[0]
        if record.get("finding_id"):
            records.append(record)

    for block in _review_blocks(reviews_text):
        current: dict[str, list[str]] | None = None
        for line in _normalise_lines(block):
            match = re.fullmatch(
                r"(?P<key>[a-z][a-z0-9_]*)=(?P<value>.*)", line
            )
            if not match:
                continue
            key = match.group("key")
            if key == "finding_id":
                finish(current)
                current = {"finding_id": [match.group("value")]}
            elif key in FINDING_FIELDS:
                if current is None:
                    errors.append(f"finding field {key} appears without finding_id")
                else:
                    current.setdefault(key, []).append(match.group("value"))
        finish(current)
    return records, errors


def _finding_gate_errors(
    path: Path,
    reviews_text: str,
    work: Mapping[str, str | int | None],
    residual_comments: Iterable[Mapping[str, Any]] = (),
) -> list[str]:
    """Check exact structured findings and require human residual-risk input."""

    records, errors = _structured_finding_records(reviews_text)
    latest: dict[str, dict[str, str]] = {}
    for record in records:
        finding_id = record.get("finding_id", "")
        if finding_id:
            latest[finding_id] = record

    residual_comments = tuple(residual_comments)
    for record in latest.values():
        finding_id = record.get("finding_id", "")
        severity = record.get("severity")
        status = record.get("status")
        if not finding_id or not severity or not status:
            continue
        if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9._-]*", finding_id):
            errors.append(f"{path}: malformed finding ID {finding_id!r}")
        if severity not in FINDING_SEVERITIES:
            errors.append(f"{path}: invalid finding severity {severity!r}")
            continue
        if status not in FINDING_STATUSES:
            errors.append(f"{path}: invalid finding status {status!r}")
            if severity in {"CRITICAL", "HIGH"}:
                errors.append(f"{path}: {severity} finding {finding_id} is not resolved")
            continue

        if severity in {"CRITICAL", "HIGH"}:
            if status == "ACCEPTED_RESIDUAL_RISK":
                errors.append(
                    f"{path}: {severity} finding {finding_id} cannot accept residual risk"
                )
            elif status != "RESOLVED":
                errors.append(f"{path}: {severity} finding {finding_id} is not resolved")

        if status == "ACCEPTED_RESIDUAL_RISK":
            comment_url = record.get("residual_risk_comment_url")
            if not comment_url:
                errors.append(
                    f"{path}: finding {finding_id} is missing residual-risk comment URL"
                )
                continue
            matching = [
                comment
                for comment in residual_comments
                if comment.get("html_url") == comment_url
            ]
            if len(matching) != 1:
                errors.append(
                    f"{path}: finding {finding_id} requires exactly one supplied residual-risk comment"
                )
            else:
                errors.extend(
                    f"{path}: {error}"
                    for error in validate_residual_comment(
                        matching[0], work, finding=record
                    )
                )
    return errors


def _validate_active_layout(
    root: Path, residual_comments: Iterable[Mapping[str, Any]] = ()
) -> list[str]:
    errors: list[str] = []
    work_root = root / "work"
    if work_root.exists():
        for path in work_root.rglob("*"):
            if path.is_dir() and path.name.lower() in FORBIDDEN_WORK_DIRECTORY_NAMES:
                relative = path.relative_to(work_root).as_posix()
                errors.append(f"forbidden work/{relative}/ archive shape exists")

    active_root = work_root / "active"
    if active_root.exists():
        for child in sorted(active_root.iterdir()):
            if not child.is_dir():
                errors.append(f"{child}: active work entries must be directories")
                continue
            files = {
                path.relative_to(child).as_posix()
                for path in child.rglob("*")
                if path.is_file()
            }
            expected = {"WORK.md", "REVIEWS.md"}
            if files != expected:
                errors.append(
                    f"{child}: active artifact set must be exactly WORK.md and REVIEWS.md"
                )
                continue
            reviews_text = _read_utf8(child / "REVIEWS.md")
            errors.extend(
                _validate_work_document(
                    child / "WORK.md", child.name, reviews_text, residual_comments
                )
            )

    phase_root = work_root / "phases" / "active"
    if phase_root.exists():
        for child in sorted(phase_root.iterdir()):
            if not child.is_dir():
                errors.append(f"{child}: active phase entries must be directories")
                continue
            files = {
                path.relative_to(child).as_posix()
                for path in child.rglob("*")
                if path.is_file()
            }
            expected = {"PHASE.md", "REVIEWS.md"}
            if files != expected:
                errors.append(
                    f"{child}: active phase artifact set must be exactly PHASE.md and REVIEWS.md"
                )
                continue
            errors.extend(_validate_phase_document(child / "PHASE.md", child.name))
    return errors


def _validate_phase_document(path: Path, phase_id_from_directory: str) -> list[str]:
    errors: list[str] = []
    try:
        text = _read_utf8(path)
    except WorkflowValidationError as exc:
        return [str(exc)]
    match = re.search(r"(?m)^Phase ID:\s*`([^`\n]+)`\s*$", text)
    if not match:
        errors.append(f"{path}: missing phase ID")
    elif match.group(1) != phase_id_from_directory:
        errors.append(f"{path}: phase ID does not match directory")
    revision = _field(text, "Plan revision") or _field(text, "Phase revision")
    if revision is None or not revision.isdigit() or int(revision) < 1:
        errors.append(f"{path}: invalid phase revision")
    declared = _field(text, "Plan hash") or _field(text, "Phase hash")
    if declared is None or not SHA256_RE.fullmatch(declared):
        errors.append(f"{path}: invalid phase hash")
    else:
        try:
            computed = phase_hash_from_text(text)
            if declared != computed:
                errors.append(f"{path}: phase hash mismatch")
        except WorkflowValidationError as exc:
            errors.append(f"{path}: {exc}")
    return errors


def _validate_durable_files(root: Path) -> list[str]:
    errors: list[str] = []
    for relative in REQUIRED_DURABLE_FILES:
        path = root / relative
        if not path.is_file():
            errors.append(f"missing durable workflow file: {relative}")
    if errors:
        return errors
    for relative in REQUIRED_DURABLE_FILES:
        path = root / relative
        try:
            text = _read_utf8(path)
        except WorkflowValidationError as exc:
            errors.append(str(exc))
            continue
        if relative != "scripts/validate_ai_workflow.py" and CONTRACT_VERSION not in text:
            errors.append(f"{relative}: missing contract version {CONTRACT_VERSION}")
    return errors


def _load_json(path: str) -> Mapping[str, Any]:
    raw = sys.stdin.read() if path == "-" else Path(path).read_text(encoding="utf-8")
    value = json.loads(raw)
    if not isinstance(value, dict):
        raise WorkflowValidationError(f"{path}: expected a JSON object")
    return value


def _parse_expected_work(root: Path) -> tuple[Path, dict[str, str | int | None]]:
    active_root = root / "work" / "active"
    candidates = [
        path / "WORK.md"
        for path in active_root.iterdir()
        if path.is_dir() and (path / "WORK.md").is_file()
    ] if active_root.exists() else []
    if len(candidates) != 1:
        raise WorkflowValidationError(
            "an explicit issue/comment input requires exactly one active WORK.md"
        )
    path = candidates[0]
    return path, parse_work(_read_utf8(path))


def validate_approval_comment(
    comment: Mapping[str, Any], work: Mapping[str, str | int | None]
) -> list[str]:
    errors: list[str] = []
    issue_parts = _web_issue_parts(work.get("issue_url"))
    expected_url = work.get("approval_url")
    approval_parts = _web_comment_parts(expected_url)
    if approval_parts is None:
        errors.append("WORK.md approval URL is not a canonical Issue comment URL")
    elif issue_parts is None or approval_parts[:3] != issue_parts:
        errors.append("WORK.md approval URL is not on the canonical Issue")

    url = comment.get("html_url")
    if not isinstance(url, str) or url != expected_url:
        errors.append("approval comment URL does not match WORK.md")
    comment_parts = _web_comment_parts(url)
    if comment_parts is None:
        errors.append("approval comment html_url is not a canonical Issue comment URL")
    elif issue_parts is None or comment_parts[:3] != issue_parts:
        errors.append("approval comment html_url is not on the canonical Issue")

    native_url = comment.get("url")
    if native_url is not None:
        if not isinstance(native_url, str):
            errors.append("approval comment API url is malformed")
        elif (
            issue_parts is None
            or comment_parts is None
            or native_url
            != _api_repository_url(issue_parts[0], issue_parts[1])
            + f"/issues/comments/{comment_parts[3]}"
        ):
            errors.append("approval comment API url is not the matching comment")
    comment_id = comment.get("id")
    if comment_id is not None:
        if isinstance(comment_id, bool) or not isinstance(comment_id, int) or comment_id < 1:
            errors.append("approval comment id is malformed")
        elif comment_parts is not None and comment_id != comment_parts[3]:
            errors.append("approval comment id does not match html_url")

    body = comment.get("body")
    expected_body = "\n".join(
        (
            "approval_type=PLAN",
            f"work_item_id={work.get('work_id')}",
            f"plan_revision={work.get('plan_revision')}",
            f"plan_hash={work.get('plan_hash')}",
            "decision=APPROVED",
        )
    )
    if not isinstance(body, str) or normalize_document(body).decode("utf-8") != (
        expected_body + "\n"
    ):
        errors.append("approval comment body does not exactly match PLAN schema")
    user = comment.get("user")
    if not isinstance(user, Mapping) or user.get("type") != "User":
        errors.append("approval comment author is not a GitHub User")
    elif user.get("is_bot") is True:
        errors.append("approval comment author is marked as a bot")
    author = comment.get("author")
    if isinstance(author, Mapping) and author.get("is_bot") is True:
        errors.append("approval comment author is marked as a bot")
    created = comment.get("created_at") or comment.get("createdAt")
    updated = comment.get("updated_at") or comment.get("updatedAt")
    if created is not None and updated is not None and created != updated:
        errors.append("approval comment is edited (created_at != updated_at)")
    return errors


def validate_residual_comment(
    comment: Mapping[str, Any],
    work: Mapping[str, str | int | None],
    *,
    finding: Mapping[str, str] | None = None,
) -> list[str]:
    errors: list[str] = []

    issue_parts = _web_issue_parts(work.get("issue_url"))
    html_url = comment.get("html_url")
    comment_parts = _web_comment_parts(html_url)
    if comment_parts is None:
        errors.append("residual-risk comment html_url is not a canonical Issue comment URL")
    elif issue_parts is None or comment_parts[:3] != issue_parts:
        errors.append("residual-risk comment html_url is not on the canonical Issue")

    comment_id = comment.get("id")
    if isinstance(comment_id, bool) or not isinstance(comment_id, int) or comment_id < 1:
        errors.append("residual-risk comment id is malformed")
    elif comment_parts is not None and comment_id != comment_parts[3]:
        errors.append("residual-risk comment id does not match html_url")

    native_url = comment.get("url")
    if not isinstance(native_url, str):
        errors.append("residual-risk comment API url is missing or malformed")
    elif (
        issue_parts is None
        or comment_parts is None
        or native_url
        != _api_repository_url(issue_parts[0], issue_parts[1])
        + f"/issues/comments/{comment_parts[3]}"
    ):
        errors.append("residual-risk comment API url is not the matching comment")

    body = comment.get("body")
    if not isinstance(body, str):
        errors.append("residual-risk comment body is missing")
        return errors
    fields: dict[str, str] = {}
    for line in normalize_document(body).decode("utf-8").splitlines():
        if "=" not in line:
            errors.append("residual-risk comment contains a non key=value line")
            return errors
        key, value = line.split("=", 1)
        if key in fields:
            errors.append(f"residual-risk comment repeats {key}")
            return errors
        fields[key] = value

    expected_keys = {
        "decision",
        "work_item_id",
        "finding_id",
        "plan_revision",
        "plan_hash",
        "candidate_git_sha",
        "rationale",
    }
    missing = expected_keys - fields.keys()
    unexpected = fields.keys() - expected_keys
    for key in sorted(missing):
        errors.append(f"residual-risk comment is missing {key}")
    for key in sorted(unexpected):
        errors.append(f"residual-risk comment has unsupported field {key}")

    expected = {
        "decision": "ACCEPTED_RESIDUAL_RISK",
        "work_item_id": str(work.get("work_id")),
        "plan_revision": str(work.get("plan_revision")),
        "plan_hash": str(work.get("plan_hash")),
    }
    for key, value in expected.items():
        if fields.get(key) != value:
            errors.append(f"residual-risk comment {key} does not match WORK.md")
    finding_id = fields.get("finding_id")
    if not finding_id:
        errors.append("residual-risk comment is missing finding_id")
    elif not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9._-]*", finding_id):
        errors.append("residual-risk comment has a malformed finding_id")
    if not fields.get("rationale"):
        errors.append("residual-risk comment is missing rationale")
    candidate = fields.get("candidate_git_sha")
    expected_candidate = work.get("candidate_sha")
    if candidate != expected_candidate:
        errors.append("residual-risk comment candidate_git_sha does not match WORK.md")
    if candidate != "NONE" and (candidate is None or not GIT_SHA_RE.fullmatch(candidate)):
        errors.append("residual-risk comment has an invalid candidate_git_sha")
    if finding is not None:
        if finding_id != finding.get("finding_id"):
            errors.append("residual-risk comment finding_id does not match the applicable finding")
        if finding.get("status") != "ACCEPTED_RESIDUAL_RISK":
            errors.append("residual-risk comment finding is not accepted residual risk")
        if finding.get("severity") in {"CRITICAL", "HIGH"}:
            errors.append("residual-risk comment cannot accept a HIGH or CRITICAL finding")
        expected_comment_url = finding.get("residual_risk_comment_url")
        if html_url != expected_comment_url:
            errors.append("residual-risk comment URL does not match the finding reference")

    user = comment.get("user")
    if not isinstance(user, Mapping) or user.get("type") != "User":
        errors.append("residual-risk comment author is not a GitHub User")
    elif user.get("is_bot") is True:
        errors.append("residual-risk comment author is marked as a bot")
    author = comment.get("author")
    if isinstance(author, Mapping) and author.get("is_bot") is True:
        errors.append("residual-risk comment author is marked as a bot")
    return errors


def _validate_supplied_issue(
    issue: Mapping[str, Any], work_path: Path, work: Mapping[str, str | int | None]
) -> list[str]:
    errors: list[str] = []
    try:
        computed = issue_digest(issue)
    except WorkflowValidationError as exc:
        errors.append(str(exc))
        computed = None
    if computed is not None and computed != work.get("issue_digest"):
        errors.append(
            f"{work_path}: supplied Issue title/body digest {computed} does not match WORK.md"
        )

    issue_parts = _web_issue_parts(work.get("issue_url"))
    if issue_parts is None:
        errors.append(f"{work_path}: WORK.md has no valid canonical Issue URL")

    number = issue.get("number")
    if isinstance(number, bool) or not isinstance(number, int):
        errors.append(f"{work_path}: supplied Issue number is missing or malformed")
    elif number != work.get("issue_number"):
        errors.append(f"{work_path}: supplied Issue number does not match WORK.md")

    if issue_parts is not None:
        owner, repo, issue_number = issue_parts
        expected_api_repo = _api_repository_url(owner, repo)
        expected_api_issue = _api_issue_url(owner, repo, issue_number)
        if issue.get("repository_url") != expected_api_repo:
            errors.append(
                f"{work_path}: supplied Issue repository identity does not match canonical mapping"
            )
        if issue.get("url") != expected_api_issue:
            errors.append(
                f"{work_path}: supplied Issue API url does not match canonical mapping"
            )
        if issue.get("html_url") != work.get("issue_url"):
            errors.append(
                f"{work_path}: supplied Issue html_url does not match canonical mapping"
            )
    return errors


def _active_pairs_at_revision(
    root: Path, revision: str
) -> tuple[dict[str, set[str]], list[str]]:
    result = subprocess.run(
        ["git", "ls-tree", "-r", "--name-only", revision, "--", "work/active"],
        cwd=root,
        capture_output=True,
        text=True,
        check=False,
    )
    if result.returncode != 0:
        return {}, ["git ls-tree failed while checking active cleanup target"]
    pairs: dict[str, set[str]] = {}
    errors: list[str] = []
    for path in result.stdout.splitlines():
        match = re.fullmatch(r"work/active/([^/]+)/(WORK|REVIEWS)\.md", path)
        if match is None:
            errors.append(f"candidate active tree contains an unpermitted path: {path}")
            continue
        pairs.setdefault(match.group(1), set()).add(match.group(2) + ".md")
    return pairs, errors


def validate_cleanup_diff(
    root: Path,
    base_sha: str,
    candidate_sha: str,
    head_sha: str,
    *,
    work_id: str | None = None,
) -> list[str]:
    """Require a candidate-to-head diff consisting only of active-file deletion."""

    errors: list[str] = []
    for label, value in (
        ("base", base_sha),
        ("candidate", candidate_sha),
        ("head", head_sha),
    ):
        if not GIT_SHA_RE.fullmatch(value):
            errors.append(f"{label} SHA is not a full 40-character hexadecimal SHA")
    if errors:
        return errors
    for ancestor, descendant, label in (
        (base_sha, candidate_sha, "base-to-candidate"),
        (candidate_sha, head_sha, "candidate-to-head"),
    ):
        result = subprocess.run(
            ["git", "merge-base", "--is-ancestor", ancestor, descendant],
            cwd=root,
            capture_output=True,
            text=True,
            check=False,
        )
        if result.returncode != 0:
            errors.append(f"{label} is not an ancestor relationship")
    result = subprocess.run(
        ["git", "diff", "--name-status", "--no-renames", candidate_sha, head_sha],
        cwd=root,
        capture_output=True,
        text=True,
        check=False,
    )
    if result.returncode != 0:
        return errors + ["git diff failed while checking cleanup-only changes"]
    deleted: dict[str, set[str]] = {}
    for line in result.stdout.splitlines():
        fields = line.split("\t")
        if len(fields) != 2 or fields[0] != "D":
            errors.append(f"cleanup diff contains a non-deletion change: {line}")
            continue
        match = re.fullmatch(r"work/active/([^/]+)/(WORK|REVIEWS)\.md", fields[1])
        if not match:
            errors.append(f"cleanup diff contains an unpermitted path: {fields[1]}")
            continue
        deleted.setdefault(match.group(1), set()).add(match.group(2) + ".md")
    if not deleted:
        errors.append("cleanup diff is empty")
    if len(deleted) != 1:
        errors.append("cleanup diff must delete exactly one reviewed work item")

    candidate_pairs, tree_errors = _active_pairs_at_revision(root, candidate_sha)
    errors.extend(tree_errors)
    for candidate_work_id, names in candidate_pairs.items():
        if names != {"WORK.md", "REVIEWS.md"}:
            errors.append(
                f"candidate active work item {candidate_work_id} must contain exactly WORK.md and REVIEWS.md"
            )

    if work_id is not None and not WORK_ID_RE.fullmatch(work_id):
        errors.append("cleanup reviewed work-item ID is malformed")
    if work_id is None:
        if len(candidate_pairs) != 1:
            errors.append(
                "cleanup requires an explicit reviewed work-item ID when the candidate has multiple active items"
            )
        elif candidate_pairs:
            work_id = next(iter(candidate_pairs))

    if work_id is not None:
        if work_id not in candidate_pairs:
            errors.append("cleanup reviewed work-item ID is not present in the candidate")
        if work_id not in deleted:
            errors.append("cleanup diff does not delete the reviewed work item")

    for deleted_work_id, names in deleted.items():
        if names != {"WORK.md", "REVIEWS.md"}:
            errors.append(
                f"cleanup diff for {deleted_work_id} must delete exactly WORK.md and REVIEWS.md"
            )
    if work_id is not None and set(deleted) != {work_id}:
        errors.append("cleanup diff contains a work item other than the reviewed work item")
    return errors


def validate_repository(
    root: Path,
    *,
    issue: Mapping[str, Any] | None = None,
    approval_comment: Mapping[str, Any] | None = None,
    residual_comments: Iterable[Mapping[str, Any]] = (),
    cleanup_shas: tuple[str, str, str] | None = None,
    cleanup_work_id: str | None = None,
) -> list[str]:
    residual_comments = tuple(residual_comments)
    errors = _validate_durable_files(root)
    errors.extend(_validate_active_layout(root, residual_comments))

    if issue is not None or approval_comment is not None or residual_comments:
        try:
            work_path, work = _parse_expected_work(root)
        except WorkflowValidationError as exc:
            errors.append(str(exc))
        else:
            if issue is not None:
                errors.extend(_validate_supplied_issue(issue, work_path, work))
            if approval_comment is not None:
                errors.extend(validate_approval_comment(approval_comment, work))
            reviews_text = _read_utf8(work_path.parent / "REVIEWS.md")
            finding_records, _ = _structured_finding_records(reviews_text)
            latest_findings = {
                record["finding_id"]: record
                for record in finding_records
                if record.get("finding_id")
            }
            for comment in residual_comments:
                body = comment.get("body")
                finding_id = None
                if isinstance(body, str):
                    for line in normalize_document(body).decode("utf-8").splitlines():
                        if line.startswith("finding_id="):
                            finding_id = line.split("=", 1)[1]
                            break
                finding = latest_findings.get(finding_id or "")
                if finding is None:
                    errors.append(
                        "residual-risk comment does not identify an applicable finding"
                    )
                errors.extend(
                    validate_residual_comment(comment, work, finding=finding)
                )

    if cleanup_shas is not None:
        errors.extend(
            validate_cleanup_diff(root, *cleanup_shas, work_id=cleanup_work_id)
        )
    return errors


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--repo-root",
        type=Path,
        default=Path(__file__).resolve().parents[1],
        help="repository root (defaults to the parent of scripts/)",
    )
    parser.add_argument(
        "--issue-json",
        help="fresh GitHub Issue JSON path, or '-' for stdin",
    )
    parser.add_argument(
        "--approval-json",
        help="fresh GitHub approval-comment JSON path, or '-' for stdin",
    )
    parser.add_argument(
        "--residual-risk-json",
        action="append",
        default=[],
        help="fresh residual-risk comment JSON path; may be repeated",
    )
    parser.add_argument("--base-sha")
    parser.add_argument("--candidate-sha")
    parser.add_argument("--head-sha")
    parser.add_argument(
        "--work-item-id",
        help="reviewed active work-item ID for cleanup-only SHA validation",
    )
    return parser


def main(argv: list[str] | None = None) -> int:
    args = _build_parser().parse_args(argv)
    root = args.repo_root.resolve()
    try:
        issue = _load_json(args.issue_json) if args.issue_json else None
        approval = _load_json(args.approval_json) if args.approval_json else None
        residual = [_load_json(path) for path in args.residual_risk_json]
    except (OSError, json.JSONDecodeError, WorkflowValidationError) as exc:
        print(f"AI workflow validation: FAIL\n- {exc}", file=sys.stderr)
        return 1

    sha_args = (args.base_sha, args.candidate_sha, args.head_sha)
    if any(value is not None for value in sha_args) and not all(
        value is not None for value in sha_args
    ):
        print(
            "AI workflow validation: FAIL\n- base, candidate, and head SHAs must be supplied together",
            file=sys.stderr,
        )
        return 1
    if args.work_item_id is not None and not all(value is not None for value in sha_args):
        print(
            "AI workflow validation: FAIL\n- --work-item-id requires cleanup SHAs",
            file=sys.stderr,
        )
        return 1
    cleanup = tuple(sha_args) if all(value is not None for value in sha_args) else None
    errors = validate_repository(
        root,
        issue=issue,
        approval_comment=approval,
        residual_comments=residual,
        cleanup_shas=cleanup,  # type: ignore[arg-type]
        cleanup_work_id=args.work_item_id,
    )
    if errors:
        print("AI workflow validation: FAIL", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print("AI workflow validation: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
