# Agent Policy — Permissions, Hard Stops, Escalation

> Audience: AI agents. Read together with `AGENTS.md`.

## 1. Permission matrix

| Action | Permission | Notes |
|---|---|---|
| Read any file in repo | ✅ Free | |
| Create/modify source code | ✅ Free | Follow `agent/rules/` for that area |
| Create/modify infra config (Dockerfile, Compose, CI, Caddyfile) | ✅ Free | |
| Create/modify migrations | ⚠️ Conditional | New migration files: free. **Editing an already-applied migration: forbidden** |
| Run build / test / lint / format | ✅ Free | Capture output as evidence |
| Run local containers (`docker compose up`) | ✅ Free | Respect memory limits; laptop may have 8 GB |
| Run load tests (k6) | ⚠️ Conditional | Announce first; can saturate the laptop |
| Install project dependencies (`go get`, `pip install`, `npm i`) | ⚠️ Conditional | State why the dependency is needed and its licence. Prefer stdlib |
| Install system packages (`apt install`) | ❌ Ask first | Changes the machine, not the repo |
| Modify `.wslconfig`, shell rc files, global git config | ❌ Ask first | Outside repo scope |
| Read `.env` | ❌ Never | Use `.env.example`. If you need a variable, add it there as a placeholder |
| Write `.env` | ❌ Never | Human-only |
| Any cloud write operation (`az`, `terraform apply`, Neon/Cloudflare API) | ❌ Never | Real money, real data. Draft the command, human executes |
| `git commit` | ✅ Free | Conventional Commits. Small, logical commits |
| `git push` | ⚠️ Conditional | Feature branches: free. `main`: only via PR |
| `git push --force`, history rewrite, branch deletion | ❌ Never | |
| Tick a DoD checkbox | ❌ Never | Propose; the human ticks |
| Edit a signed-off weekly report | ❌ Never | Append a correction note in the next week's report instead |
| Write ADR **Context** and **Options** | ✅ Free | Drafting is helpful |
| Write ADR **Decision** and **Consequences** | ❌ Never | The human owns trade-offs; they must be able to defend them |

## 2. Hard stops — halt and ask

Stop immediately and ask the human when:

1. The task requires a credential you do not have (and must not have).
2. The task requires a cloud write operation.
3. The roadmap says one thing and the repo state says another.
4. Two rules in `agent/rules/` conflict.
5. The requested change would break a documented DoD from a previous week.
6. You are about to add a third-party dependency for something the standard library handles.
7. You cannot verify your work locally — say so *before* implementing, not after.
8. You notice a security issue (leaked secret, IDOR, injection) — report it, do not silently "fix and move on".

## 3. Escalation format

When stopping, use exactly this shape so the human can decide quickly:

```
STOPPING — need a decision

Task:        <what was asked>
Blocker:     <what prevents proceeding>
Rule:        <which policy/rule triggered this>
Options:
  A) <option> — cost: <...>  risk: <...>
  B) <option> — cost: <...>  risk: <...>
Recommendation: <A or B, and why — labelled as RECOMMENDATION, not FACT>
Need from you: <the exact one thing you need>
```

## 4. Dependency policy

Before adding any dependency, answer in the report:

1. What problem does it solve that the standard library cannot?
2. How actively is it maintained (last release date)?
3. What is its licence?
4. What is the removal cost if it turns out to be wrong?

Reject a dependency if the honest answer to (1) is "convenience". This project deliberately favours the standard library — it teaches more and reduces supply-chain surface, which matters because CI runs Trivy and `govulncheck`.

## 5. Resource discipline

The development laptop is an **AMD Ryzen 5 7430U, 8–16 GB RAM, 256 GB SSD**. The production VM is an **Azure B1s: 1 vCPU, 1 GB RAM**.

Consequences you must respect:

- Never assume memory headroom. Set container memory limits.
- Never start the full observability stack alongside the core stack without warning the human.
- Prefer batch/offline processing (embeddings) over always-on services.
- Keep Docker images small; run `docker system prune` guidance in the report if disk pressure appears.
- Before suggesting a tool, check whether it fits 1 GB of RAM in production.
