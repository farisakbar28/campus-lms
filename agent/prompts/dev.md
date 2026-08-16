# Prompt — DEV

```
Execute the approved plan.

Constraints (from AGENTS.md and agent/rules/):
- Smallest coherent change; the tree stays runnable.
- English in code, comments, and commit messages.
- Explain WHY for every non-obvious decision, in a brief comment.
- Do not touch .env, do not run cloud write commands, do not tick DoD boxes.

After implementing, you MUST:
1. Run it. Capture the raw output.
2. Save evidence to docs/progress/evidence/week-<NN>/<slug>.txt using the
   header format in agent/evidence-protocol.md.
3. Report honestly:
   - what you implemented (file paths)
   - what you verified, with the command used
   - what you could NOT verify, and why
   - any assumption you had to make
4. Propose the commit message (Conventional Commits).

If a command fails, show me the failure verbatim. Do not summarise it into
something more comfortable, and do not silently work around it.
```
