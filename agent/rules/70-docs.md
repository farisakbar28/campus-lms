# Rule 70 — Documentation and Reporting

Documentation here is not decoration; it is the learning mechanism and the portfolio surface.

## Weekly report (`docs/progress/week-NN.md`)

Use `agent/templates/weekly-report.md` exactly. Mandatory sections, in order:
Ringkasan · Dikerjakan Agent · Dikerjakan Manusia · Keputusan · Angka & Bukti · Konsep yang Dipelajari · **Belum Terverifikasi** · Masalah & Cara Diselesaikan · DoD Status · Untuk Minggu Depan.

Writing rules:

- Plain Indonesian in the narrative sections. No marketing language, no "successfully implemented a robust solution".
- Every number links to its evidence file.
- Ask the human what they did manually; do not infer it and do not leave the section empty.
- Dead ends and mistakes belong in the report. They are the most educational part and the best interview material.
- Never tick a DoD checkbox — propose status only.

## Quiz (`docs/progress/quiz/week-NN.md`)

Use `agent/templates/quiz.md`. Every question must be answerable from code that exists — verify by reading the file first. Include an answer key with file references at the bottom.

## ADRs (`docs/adr/`)

You draft **Context** and **Options**. The human writes **Decision** and **Consequences** — they must own and be able to defend the trade-off. Never mark an ADR "Accepted" yourself.

## Experiment notes (`docs/notes/`)

Tables with real measurements. Include the failed configurations — knowing what does not work is half the value, and it demonstrates method.

## README (`README.md`, English)

The README is what a recruiter reads first. It must state the problem, the architecture, the **numbers**, the trade-offs, and the demo links. Keep the status table current. When the project matures in Week 12, rewrite the top section to lead with outcomes, not setup instructions.

## Runbooks (`docs/runbook/`)

Written to be followed at 3 a.m. by someone stressed: symptom → detection → diagnosis commands → mitigation → permanent fix → prevention. Every scenario must have been rehearsed at least once, with its MTTR recorded.

## Tone

Direct and honest. State uncertainty as uncertainty. If something is fragile, say it is fragile — the owner needs to know where the bodies are buried before an interviewer finds them.
