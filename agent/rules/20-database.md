# Rule 20 — Database (PostgreSQL, pgvector)

## Environments

| Environment | Where | Purpose |
|---|---|---|
| Local dev/test | Postgres 16 in Docker | Where the owner learns operations: migrations, indexes, EXPLAIN, backup, restore |
| Production | Neon (free tier) | Because the Azure B1s VM has only 1 GB RAM |

Keep them schema-identical. Any drift is a bug.

## Migrations

- Versioned, forward-only, one concern per migration.
- **Never edit a migration that has been applied anywhere.** Write a new one.
- Every migration must be reversible in principle; write the down-migration or explain why it is impossible.
- Schema changes that touch running code use **expand-contract**: add new → backfill → switch reads → switch writes → drop old. Never a breaking change in a single deploy.

## Multi-tenancy (shared schema + RLS)

- Every tenant-scoped table has `tenant_id` NOT NULL with an index.
- `ENABLE ROW LEVEL SECURITY` plus a policy using `current_setting('app.tenant_id')`.
- Policies need both `USING` (read) and `WITH CHECK` (write) — omitting `WITH CHECK` lets a tenant write rows it cannot read.
- There must be a test that **fails if RLS is disabled**. Isolation is a correctness property, not a configuration detail.
- Beware connection pooling: `SET LOCAL` inside a transaction, never `SET` on a pooled connection (it leaks to the next request — a real and dangerous bug).

## Indexes

- Add an index because a query needs it, with `EXPLAIN` evidence before and after — never speculatively.
- Composite index column order follows the leftmost-prefix rule; state the reasoning in the migration comment.
- Consider partial indexes for soft-deleted or status-filtered rows.
- Remember the cost: every index slows writes and consumes disk.

## Vectors (from Week 8)

- `pgvector` with HNSW for retrieval; record the `m` / `ef_construction` / `ef_search` values used and why.
- Vector tables are tenant-scoped and RLS-protected too. A retrieval leak across tenants is a data breach, not a bug.
- Store chunk metadata (`material_id`, `page`, `section`) — citations depend on it.

## Data handling

- `timestamptz` always, never naive timestamps.
- Money/grades: `numeric`, never floating point.
- Student data is personal data (UU PDP context): collect the minimum, never log it, never send it to an LLM provider without a documented reason.

## Backup

- `pg_dump -Fc`, timestamped, uploaded off-machine, 7-day retention.
- **A backup that has never been restored is not a backup.** The restore drill is a deliverable, and its duration (RTO) is a number that goes in the report.
