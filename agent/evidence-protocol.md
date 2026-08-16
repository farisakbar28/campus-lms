# Evidence Protocol — How to Prove What You Claim

> Audience: AI agents. This is the mechanism that makes agent-first learning safe.
> Violating this protocol is the most serious failure mode in this repo.

## 1. Why this exists

The human owner is learning from what you write. A fabricated benchmark or an
invented API signature does not just produce a bug — it produces a **false
belief** that gets repeated in a job interview. That damage is not visible in
the diff and not caught by tests.

Therefore: **claims without receipts are treated as defects.**

## 2. Evidence file format

Store every piece of evidence as a plain-text file:

```
docs/progress/evidence/week-<NN>/<slug>.txt
```

Each file starts with a mandatory header, then the **raw, unedited** output:

```
=== EVIDENCE ===
CLAIM:    API container image is under the 25 MB target
COMMAND:  docker images campus-lms-api:dev --format "{{.Size}}"
CWD:      /home/user/campus-lms
RUN AT:   2026-08-20T14:32:11+08:00
COMMIT:   a3f9c21
EXIT:     0
=== RAW OUTPUT ===
21.4MB
=== END ===
```

Rules:

- **Never edit the raw output.** Not to shorten, not to prettify. If it is long, keep the full file and quote an excerpt in the report.
- **Failures are evidence too.** A failing run gets saved with `EXIT: 1`. Do not hide it.
- One claim per file. Slugs must be descriptive: `image-size.txt`, `explain-analyze-course-list-before.txt`.
- Include the commit SHA so the claim is reproducible.

## 3. What requires evidence

| Claim type | Required evidence |
|---|---|
| Any latency / throughput number | k6 or `curl -w` output, full run |
| Any "test passes" statement | `go test` / `pytest` output including the summary line |
| Coverage percentage | `go tool cover -func` tail |
| Image size | `docker images` output |
| Query performance | `EXPLAIN (ANALYZE, BUFFERS)` before **and** after |
| Memory/CPU behaviour | `docker stats` snapshot or `free -h` |
| Security scan result | Trivy / gosec / govulncheck output |
| Eval score (RAG, agent) | Harness output with dataset size and date |
| "Service is healthy" | `curl` output with status code |
| Deployment success | Deploy script output + post-deploy health check |

## 4. What must NEVER be stated as fact

- A number you did not measure. Write `NOT MEASURED` — that is an acceptable, honest report line.
- "Typical" or "usually around" figures from memory.
- That a DoD item is met. You may write `DoD candidate: met, pending human verification`.
- That code works because it compiles. Compiling is not working.
- An API signature you did not verify against the pinned version in `go.mod` / `pyproject.toml`.
- Anything about production behaviour you observed only locally. Label it `LOCAL ONLY`.

## 5. Report structure (mandatory sections)

Every `docs/progress/week-NN.md` must contain, in this order:

1. **Ringkasan** — 3–5 sentences, plain language
2. **Dikerjakan Agent** — with file paths and commit SHAs
3. **Dikerjakan Manusia** — what the human did manually (ask them; do not assume)
4. **Keputusan yang Diambil** — decisions and trade-offs, linked to ADRs
5. **Angka & Bukti** — table of every measurement with a link to its evidence file
6. **Konsep yang Dipelajari** — the teaching section (see §6)
7. **Belum Terverifikasi** — mandatory; an empty section is a red flag
8. **Masalah & Cara Diselesaikan** — including dead ends, because that is where learning lives
9. **DoD Status** — proposed status only; the human ticks
10. **Untuk Minggu Depan** — carry-over items

Use `agent/templates/weekly-report.md` verbatim.

## 6. The teaching section

For each significant concept implemented that week, write:

```
### <Concept name>

**Apa:** one-sentence definition, no jargon-stacking.
**Kenapa dipakai di sini:** the concrete reason in THIS codebase, with file reference.
**Alternatif yang tidak dipilih:** and what it would have cost.
**Cara membuktikannya sendiri:** an exact command the human can run to see it.
**Pertanyaan interview terkait:** the question a hiring manager would actually ask.
```

The last two lines are what turn documentation into learning.

## 7. Quiz generation rules

Questions go in `docs/progress/quiz/week-NN.md`, using `agent/templates/quiz.md`.

- 8–12 questions per week.
- **Every question must be answerable from code that actually exists in this repo.** Verify by reading the file before writing the question.
- Mix: 40% "why did we do X" (reasoning), 30% "what happens if" (failure modes), 20% "where is X and what does it do" (orientation), 10% recall.
- Include the answer key in a collapsed section at the bottom, with the file path where the answer can be confirmed.
- Never ask about code you generated but did not run.

## 8. Self-audit before submitting a report

Run this checklist mentally, and state in the report that you did:

- [ ] Every number in this report links to an evidence file that exists
- [ ] I re-read the files I claim to have created
- [ ] I ran the tests I claim pass, in this session
- [ ] The "Belum Terverifikasi" section reflects genuine gaps
- [ ] I did not tick any DoD box
- [ ] I labelled FACT / INFERENCE / RECOMMENDATION where relevant
- [ ] Quiz questions reference real files, verified by reading them
