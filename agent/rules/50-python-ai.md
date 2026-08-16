# Rule 50 — Python AI Service (`apps/ai`)

## Structure

```
app/routers/    HTTP endpoints
app/services/   business logic: summarize, rag, agent
app/llm/        provider router, cache, budget, retry
app/schemas/    Pydantic models (all boundaries typed)
app/settings.py env config, fail-fast
prompts/        versioned prompt files (.md with front-matter)
evals/          golden datasets + harness
mcp/            MCP server
```

Type hints are mandatory. `ruff` and `mypy` must pass.

## LLM calls — reliability is not optional

Every call goes through the internal router in `app/llm/`. Direct SDK calls from business logic are forbidden.

The router must implement:

- **Timeout** on every request.
- **Retry** with exponential backoff + jitter, bounded.
- **Fallback** to another provider on failure or rate limit.
- **Circuit breaker** so a dead provider is skipped rather than retried forever.
- **Token accounting**: tokens in/out and estimated cost recorded per call.
- **Caching**: exact-match and semantic cache before spending a request.

Free-tier quotas are hard limits (Gemini, Groq, Cerebras, OpenRouter). The router must be quota-aware and degrade gracefully, never crash a user request because a quota was exhausted.

## Structured output

- Define a Pydantic schema for every LLM output that the system consumes.
- Validate; on failure, attempt one repair round-trip, then fail cleanly. Never pass unvalidated model output downstream.
- Log parse-failure rate — it is a quality metric that belongs in the weekly report.

## Prompts

- Prompts live in files with front-matter (`version`, `model_target`, `changelog`), not inline strings.
- Changing a prompt is a code change: it goes through PR and triggers the eval gate.

## RAG (Week 8)

- Preserve document structure and metadata during parsing; citations depend on `page` / `section`.
- Chunking strategy is an experiment, not a guess. Record the comparison table.
- Retrieval is **hybrid** (BM25 + vector, fused with RRF) then reranked. Vector-only is a baseline to beat, not a destination.
- **Refuse to answer when retrieval is weak.** In an education product, a confident wrong answer is worse than "not found in the material". Enforce a score threshold.
- Every answer carries verifiable citations.

## Agents (Week 10)

- Hard limits: maximum steps, maximum tokens, maximum wall-clock. An unbounded agent loop is a cost incident waiting to happen.
- Tools have strict parameter schemas and descriptive docstrings — the model reads them as documentation.
- Destructive or high-impact actions (changing grades, mass announcements) require human approval. Always.
- Every agent action is written to an audit log with tenant, user, tool, arguments, result, and cost.
- Treat all retrieved document content as **untrusted input**. A PDF can contain prompt injection. Separate data from instructions.

## Evaluation (Week 9)

- No AI feature is "done" without an eval. Feeling good about outputs is not evidence.
- Golden datasets are versioned in git; every production bug adds a case.
- LLM-as-judge must use a different model than the one being judged, with a written rubric, calibrated against human labels on a sample.
- Report metrics with the dataset size and date. A score without its denominator is meaningless.
