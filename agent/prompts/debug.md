# Prompt — DEBUG (disciplined diagnosis)

```
Symptom: <exact behaviour, exact error text, when it started>
Environment: <local / VM / CI>
Last known good: <commit or time>

Do NOT propose a fix yet.

Step 1 — Reproduce. Give me the exact command. Confirm you can see it too.
Step 2 — Narrow. Is it code, config, data, network, or resource? Prove which.
Step 3 — Hypotheses. List them ranked, with the evidence for and against.
Step 4 — Test one at a time. Show raw output after each.
Step 5 — Root cause. Label clearly: PROVEN vs INFERRED.
Step 6 — Fix, verify, save evidence.
Step 7 — Prevention: what test, alert, or type would have caught this?

If you find yourself guessing, stop and say so.
```
