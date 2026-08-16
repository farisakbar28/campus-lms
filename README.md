# campus-lms

Multi-tenant SaaS Learning Management System for higher education — built as the
practical vehicle for a 12-week engineering curriculum covering infrastructure,
production operations, and AI engineering.

**Status:** Week 0 — scaffolding complete, implementation starting.

---

## Why this repository looks the way it does

Two goals run in parallel here, and they are equally important:

1. **Ship a production-grade system.** Real deployment, real CI/CD, real
   observability, real AI features — not a tutorial project.
2. **Make the owner genuinely understand it.** Well enough to defend every
   decision in a technical interview.

The second goal constrains the first. That is why this repo contains an
unusual amount of process: an agent operating contract (`AGENTS.md`), an
evidence protocol, weekly reports, and verification quizzes. Most of the
implementation is done by AI agents; the learning happens through
**evidence-based verification** rather than through typing.

If you are an AI agent working in this repository, start with `AGENTS.md`.
It is not optional reading.

---

## Architecture (target)

```
[ Cloudflare Pages ]  Next.js frontend
          |
[ Caddy ] -> [ api-go ] -> [ Neon Postgres + pgvector ]
                |
             [ ai-svc ] (FastAPI: RAG, LangGraph agent, MCP server)
                |
             [ redis ]  [ worker-go ]
```

Deliberately distributed: the production VM (Azure B1s) has **1 GB of RAM**, so
the database, frontend, and heavy observability components live elsewhere on
free tiers. That constraint is documented as an architecture decision, not
hidden — see `docs/adr/`.

**Total infrastructure budget: $0.** No credit card anywhere in the stack.

---

## Tech stack

| Layer | Choice |
|---|---|
| API | Go (chi), modular monolith |
| AI service | Python, FastAPI, LangGraph |
| Frontend | Next.js, TypeScript, Tailwind |
| Database | PostgreSQL 16 + pgvector (Docker locally, Neon in production) |
| Cache / queue | Redis, Go job queue |
| Reverse proxy | Caddy (automatic TLS) |
| CI/CD | GitHub Actions, GHCR |
| Observability | OpenTelemetry, Prometheus, Grafana, Loki, Tempo |
| LLM observability | Langfuse (self-hosted) |
| Evaluation | ragas, promptfoo |
| Infrastructure | Azure VM, Cloudflare, Terraform, k3s (learning) |

---

## Repository layout

```
apps/
  api/            Go API — auth, RBAC, multi-tenancy, LMS domain
  web/            Next.js frontend
  ai/             Python AI service — RAG, agents, MCP
deploy/
  compose/        docker compose (dev and prod)
  caddy/          reverse proxy and TLS
  scripts/        deploy, backup, restore
agent/            agent operating rules, prompts, templates, checklists
docs/
  roadmap.md      the 12-week curriculum (source of truth for what to build when)
  adr/            architecture decision records
  setup/          human-only setup checklists
  progress/       weekly reports, evidence, quizzes (the learning record)
  runbook/        incident procedures
  notes/          experiment notes with measurements
  journal/        daily journal
```

---

## Getting started

```bash
cp .env.example .env      # fill in values (never committed)
make help                 # list all commands
make todo                 # outstanding work, labelled by week
make agent-context        # what an agent must read before starting
```

Once Week 2 is complete:

```bash
make up                   # start the core stack
make health               # verify health endpoints
make logs                 # follow logs
```

Prerequisites: Docker + Compose v2, Go 1.23+, Node 20+, Python 3.12+.

---

## The learning system

**New here? Read `WORKFLOW.md` first** — one page, explains the whole loop.

| Command | Purpose |
|---|---|
| `make week-init W=01` | Scaffold the week's report, quiz, and evidence directory |
| `make evidence W=01 SLUG=x CMD="..."` | Run a command and save its raw output as evidence |
| `make gate` | Show the checklist that must pass before starting the next week |
| `make journal` | Create today's journal entry |

Every measurable claim in this repository traces to a file in
`docs/progress/evidence/`. Numbers without evidence are treated as defects.
See `agent/evidence-protocol.md`.

---

## Progress

| Week | Focus | Status |
|---|---|---|
| 0 | Cloud account hardening, agent setup | 🟡 in progress |
| 1 | Linux, networking, Git, repo foundation | ⬜ |
| 2 | Docker and Compose | ⬜ |
| 3 | PostgreSQL, multi-tenancy with RLS | ⬜ |
| 4 | Azure deployment, TLS, hardening | ⬜ |
| 5 | CI/CD and testing | ⬜ |
| 6 | Observability and performance | ⬜ |
| 7 | LLM features in production | ⬜ |
| 8 | Production RAG | ⬜ |
| 9 | Evaluation and LLM observability | ⬜ |
| 10 | Agents, MCP, guardrails | ⬜ |
| 11 | Kubernetes and IaC | ⬜ |
| 12 | Security, documentation, portfolio | ⬜ |

---

## TODO — rewrite this section in Week 12

By Week 12 this README is the first thing a recruiter reads. Replace the top
with outcomes rather than setup:

- [ ] The real problem this solves, in plain language
- [ ] Architecture diagram (real, not ASCII)
- [ ] **Numbers**: p95 latency, load test results, coverage, RAG eval scores, cost per request
- [ ] Key technical decisions and their trade-offs (link to ADRs)
- [ ] Live demo, 5-minute video, read-only Grafana dashboard, eval report
- [ ] One-command local setup, verified by someone else in under 10 minutes

---

## License

MIT — see `LICENSE`. *(TODO Week 1: add the LICENSE file.)*
