# ai-svc — Python FastAPI

**Not started. Begins in Week 7.**

Do not start early. This service depends on the job queue (Week 7), pgvector
(Week 8), and observability (Week 6). Building it now means rewriting it later.

Agent: read `agent/rules/50-python-ai.md` before touching this directory.

## Planned contents

| Week | Contents |
|---|---|
| 7 | FastAPI skeleton, multi-provider LLM router, material auto-summary |
| 8 | RAG pipeline: parsing, chunking, embedding, hybrid search, reranking |
| 9 | Eval harness (ragas + LLM-as-judge), Langfuse integration |
| 10 | LangGraph agent, MCP server, guardrails, red-team suite |

## Planned structure

```
apps/ai/
  app/
    routers/      HTTP endpoints
    services/     rag, summarize, agent
    llm/          provider router, cache, budget, retry, circuit breaker
    schemas/      Pydantic models
    settings.py   env config, fail-fast
  prompts/        versioned prompts (.md with front-matter)
  evals/          golden datasets + harness
  mcp/            MCP server
  tests/
  pyproject.toml
  Dockerfile
```
