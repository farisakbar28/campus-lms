# Prompt Library

Templates matched to the owner's existing opencode combos: `plan`, `dev`,
`review`, `heavy`, `quick`, `spark`.

Copy, fill the `<...>` placeholders, send. Keep the context-loading lines —
they are what prevent the agent from inventing repo state.

| Combo | File | Use when |
|---|---|---|
| plan | `plan.md` | Starting a task; you want a plan before code |
| dev | `dev.md` | Executing an approved plan |
| review | `review.md` | Before merging; second-pass critique |
| heavy | `heavy.md` | Weekly report, complex debugging, architecture work |
| quick | `quick.md` | Small mechanical edits |
| spark | `spark.md` | Exploring options before deciding |
| — | `debug.md` | Something is broken and you need disciplined diagnosis |
| — | `teach.md` | You do not understand code that already exists |
