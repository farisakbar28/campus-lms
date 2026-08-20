-- apps/api/testdata/week-03/fail_term_range.sql
-- Proves that an academic term cannot start after it ends.

BEGIN;

INSERT INTO academic_terms (id, tenant_id, external_id, code, name, starts_at, ends_at, status)
VALUES (
  gen_random_uuid(),
  md5('seed-tenant-a')::uuid,
  'invalid-time-range',
  'INV',
  'Invalid Term',
  '2026-12-01 00:00:00+00',
  '2026-08-01 00:00:00+00', -- Ends before it starts
  'active'
);

ROLLBACK;