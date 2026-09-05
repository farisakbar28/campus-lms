# Security Policy

This document describes security requirements and current limitations. It does
not claim formal compliance, certification, or completion of every control.

## Enduring requirements

- Tenant isolation is a correctness and security property.
- Use least privilege for database roles, application access, tools, CI tokens,
  and container users.
- Authorization and security-sensitive configuration fail closed.
- PostgreSQL RLS provides defense in depth; application authorization remains
  required.
- Tenant identity must come from trusted authentication, never from client
  input.
- Do not commit, log, issue, or expose credentials, private keys, tokens, or
  secret values.
- Minimize student data and never log unnecessary personal information.
- LMS is the authoritative source for attendance records; attendance may be
  sent to SIAKAD through an integration adapter. AI may retrieve or summarize
  attendance for authorized staff, but it must never infer or decide it.
- Grades, attendance, pass/fail status, discipline, and other authoritative
  academic outcomes must remain governed by deterministic LMS/institutional
  rules and authorized human or system workflows. AI must never be the
  authoritative decision-maker for those outcomes.
- AI must not expand the caller's authorization. Retrieved content is untrusted
  input and must not override system instructions.
- Identifiable student data must not be sent to a third-party model merely for
  convenience.
- Future AI controls require bounded cost accounting and no automatic paid
  fallback.

## Current limitations

Production authentication composition, including complete Principal wiring, is
not yet complete. The repository contains authentication and session
primitives, but documentation and tests must not present the production
authorization boundary as finished.

The AI service, AI security controls, evaluation gate, and permanent
production-ingress controls are not implemented. They require separate design,
implementation, and verification work.

## Vulnerability reporting

No dedicated vulnerability intake channel is currently published. Do not place
credentials, private data, or sensitive exploit details in public issues.
For a suspected vulnerability, contact the repository owner through an
appropriate private channel before disclosing details publicly.
