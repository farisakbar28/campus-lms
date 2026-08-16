# Prompt — HEAVY

Use for weekly reports, deep debugging, and architectural work.

## Weekly report

```
Compile docs/progress/week-<NN>.md for this week.

Sources of truth — read them, do not rely on memory:
- git log for this week (list the actual commits and SHAs)
- the files actually changed
- every evidence file in docs/progress/evidence/week-<NN>/
- the roadmap section for MINGGU <NN>

Follow agent/templates/weekly-report.md exactly. Requirements:
- Every number links to an existing evidence file. If a number has no
  evidence, write NOT MEASURED instead.
- The "Belum Terverifikasi" section must be genuine. An empty one will be
  rejected.
- Ask me what I did manually this week before filling "Dikerjakan Manusia".
- Do not tick any DoD checkbox; propose status only.
- Run the self-audit in agent/evidence-protocol.md §8 and state that you did.

Then generate docs/progress/quiz/week-<NN>.md following agent/templates/quiz.md.
Every question must be answerable from code that exists in this repo — verify
by opening the file before writing the question.
```

## Deep debugging

```
Problem: <symptom, with exact error text>
Already tried: <what you did and what happened>

Work through this systematically:
1. State 3-5 hypotheses ranked by likelihood, with your reasoning.
2. For each, the cheapest command that would confirm or eliminate it.
3. Run them in order. Show me the raw output at each step.
4. Conclude with the actual root cause — distinguish clearly between what
   you PROVED and what you are INFERRING.
5. Fix, verify, and save evidence.
6. Tell me what would have caught this earlier (a test? an alert? a type?).

Do not guess at a fix before the diagnosis is complete.
```
