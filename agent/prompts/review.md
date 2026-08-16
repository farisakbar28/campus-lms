# Prompt — REVIEW

```
Review the changes in <branch/diff/files> as a senior engineer would, before
this is merged.

Check, in this order:
1. Correctness — does it do what it claims? Any logic error?
2. Failure modes — what happens on timeout, empty input, concurrent access,
   dependency down, malformed data?
3. Security — OWASP API Top 10 (especially BOLA/IDOR and tenant isolation),
   secrets handling, input validation. If AI code: OWASP LLM Top 10.
4. Resource use — memory footprint (production VM has 1 GB), N+1 queries,
   unbounded loops, missing timeouts.
5. Tests — is the meaningful behaviour covered? What is deliberately untested?
6. Readability — will the owner understand this in three months?

For each finding: severity (blocker / should-fix / nitpick), the file and
line, why it matters, and the concrete fix.

End with the three questions an interviewer would most likely ask about this
change, so I can prepare for them.

Be direct. Do not soften findings. If the change is fine, say so briefly
rather than inventing issues.
```
