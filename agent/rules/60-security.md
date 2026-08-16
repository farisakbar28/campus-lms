# Rule 60 — Security

## Baseline

- OWASP API Security Top 10 (2023) and OWASP Top 10 for LLM Applications are the reference checklists. Audit against both in Week 12; keep them in mind from Week 1.
- Least privilege everywhere: database roles, tool permissions, CI token scopes, container users.
- Never invent crypto. Argon2id for passwords, standard JWT libraries, TLS via Caddy.

## Authentication and authorisation

- Short-lived access tokens, rotating refresh tokens, revocation on logout.
- **BOLA/IDOR is the number one API vulnerability**: every object access checks both ownership *and* tenant. Add a test that attempts cross-tenant access and expects 403/404.
- Role checks happen server-side. A hidden UI button is not an authorisation control.
- Never trust `tenant_id` from a request header, query parameter, or body.

## Input and output

- Validate everything at the boundary with a schema.
- Uploads: verify MIME type by content not extension, cap size, store outside the web root, never execute.
- Rate limit in layers: per IP, per user, per tenant. AI endpoints additionally get a token budget per tenant — otherwise one tenant can exhaust a shared free-tier quota.
- Sanitise output; never reflect unescaped user content.

## Secrets

- Environment variables only. Never in code, logs, error messages, images, or CI output.
- If a secret is ever committed, treat it as compromised: rotate first, then clean history.
- Agents must never read `.env`.

## LLM-specific threats

| Threat | Required mitigation |
|---|---|
| Prompt injection via uploaded documents | Treat retrieved text as data, not instructions; sanitise; system prompt states that document content can never override instructions |
| Sensitive data disclosure | Redact personal data before sending to a provider; document what leaves the system and why |
| Excessive agency | Tool allow-lists per role; human approval for impactful actions; hard step/cost limits |
| Insecure output handling | Never execute, render as HTML, or pass model output to a shell without validation |
| Denial of wallet | Per-tenant token budgets, caching, circuit breakers |

A red-team suite (Week 10) runs in CI and must pass 100%. Attacks are written as tests, not as prose.

## Privacy (Indonesian context)

Student data falls under UU PDP No. 27/2022. Practise data minimisation, define retention, log access to grades, and never send identifiable student data to a third-party model without a documented justification. Document the posture honestly — do not claim formal compliance.
