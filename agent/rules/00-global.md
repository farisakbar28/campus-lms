# Rule 00 — Global Engineering Rules

Applies to every file in this repository.

## Language

- **English:** code, comments, docstrings, commit messages, README, ADRs, API docs, these rule files, prompts.
- **Bahasa Indonesia:** weekly reports, quizzes, journal entries, setup checklists, `agent/README.md`.

Rationale: the repo is a portfolio aimed at Indonesian *and* remote-international employers; English is the professional default for code. Learning material stays in the owner's first language to reduce cognitive load while absorbing new concepts.

## Definition of "done"

A change is done when **all** of these hold:

1. It compiles / imports cleanly.
2. It runs, and you have the output saved as evidence.
3. It has a test if it introduces behaviour.
4. Linters pass.
5. Its rationale is written down (code comment for local decisions, ADR draft for architectural ones).
6. The weekly report reflects it honestly, including what remains unverified.

## Change size

Prefer the smallest coherent change. The tree must remain runnable after every commit. If a change cannot be small, split it and say why in the report.

## Error handling

- Never swallow an error. Wrap with context: `fmt.Errorf("load config: %w", err)`.
- Never `panic` outside `main()` / initialisation.
- Errors crossing an API boundary get a stable error code and a message safe to show a user; internal detail goes to the log, not the response.

## Logging

- Structured only (`log/slog` in Go, `structlog`/stdlib JSON in Python).
- Every log line carries: `service`, `env`, and from Week 6 onward `trace_id`.
- Never log secrets, tokens, full request bodies, or personal data of students.
- Log levels mean something: `error` = a human must look; `warn` = degraded but handled; `info` = state change; `debug` = development only.

## Configuration

12-factor: everything from environment variables. No config file with real values in the repo. No default value for a security-sensitive variable — fail fast and loudly if it is missing.

## Comments

Explain **why**, not what. If a line needs a "what" comment, the code is unclear — fix the code. Every non-obvious trade-off gets one line, because the owner will be asked about it in an interview.

## Specification is human-owned

TASK BRIEF blocks, Definition-of-Done lists, and acceptance criteria are the
human's specification. Never edit them — not to clarify, not to generalise, not
to make a verification command pass. If the wording is wrong or a check hits a
false positive, raise it with the STOPPING format and let the human decide.

Verification commands prove behaviour. If a command is imprecise, propose a
better command; do not adjust the code or comments it inspects.

## Files you must never create

- `.env` with real values
- Backup copies like `main_old.go`, `handler.bak` — that is what git is for
- Large binaries or datasets (evidence files are plain text; keep them small)
