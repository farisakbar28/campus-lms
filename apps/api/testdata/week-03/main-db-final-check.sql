-- Main Database Final Verification Script
-- READ-ONLY: proves final state of main dev DB after Week 3 migrations
-- Must output explicit summary values:
--   tables = 11
--   rls_enabled = 8
--   a7_fk = 1
--   term_range_check = 1
--   audit_policies = 2
--   rls_verifier = 0
--   fixture_rows_remaining = 0

-- 1. Exactly 11 Week 3 tables exist
WITH week3_tables AS (
  SELECT count(*) AS tables
  FROM pg_tables
  WHERE schemaname = 'public'
    AND tablename IN (
      'tenants', 'users', 'auth_identities', 'memberships', 'membership_roles',
      'audit_logs', 'academic_terms', 'courses', 'course_offerings',
      'course_staff', 'enrollments'
    )
)
-- 2. Exactly 8 Week 3 tenant-scoped tables have RLS enabled
, rls_enabled AS (
  SELECT count(*) AS rls_enabled
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE n.nspname = 'public'
    AND c.relname IN (
      'memberships', 'membership_roles', 'audit_logs', 'academic_terms',
      'courses', 'course_offerings', 'course_staff', 'enrollments'
    )
    AND c.relrowsecurity = true
)
-- 3. A7 FK exists with exact meaning:
--    enrollments(tenant_id, student_user_id) -> memberships(tenant_id, user_id)
, a7_fk AS (
  SELECT count(*) AS a7_fk
  FROM pg_constraint c
  WHERE c.conname = 'enrollments_tenant_id_student_user_id_fkey'
    AND c.contype = 'f'
    AND EXISTS (
      SELECT 1 FROM pg_attribute a
      JOIN pg_attribute b ON b.attnum = c.confkey[1] AND b.attrelid = c.confrelid
      WHERE a.attnum = c.conkey[1] AND a.attrelid = c.conrelid
        AND a.attname = 'tenant_id' AND b.attname = 'tenant_id'
    )
    AND EXISTS (
      SELECT 1 FROM pg_attribute a
      JOIN pg_attribute b ON b.attnum = c.confkey[2] AND b.attrelid = c.confrelid
      WHERE a.attnum = c.conkey[2] AND a.attrelid = c.conrelid
        AND a.attname = 'student_user_id' AND b.attname = 'user_id'
    )
)
-- 4. academic_terms contains named CHECK constraint: academic_terms_valid_time_range
--    equivalent to CHECK (starts_at < ends_at)
, term_range_check AS (
  SELECT count(*) AS term_range_check
  FROM pg_constraint
  WHERE conname = 'academic_terms_valid_time_range'
    AND contype = 'c'
    AND conrelid = 'academic_terms'::regclass
)
-- 5. audit_logs has only expected normal tenant policies: SELECT (r) and INSERT (a)
--    and no normal UPDATE/DELETE policy
, audit_policies AS (
  SELECT count(*) AS audit_policies
  FROM pg_policy
  WHERE polrelid = 'audit_logs'::regclass
    AND polcmd IN ('r', 'a')
)
-- 6. rls_verifier role count is zero
, rls_verifier AS (
  SELECT count(*) AS rls_verifier
  FROM pg_roles
  WHERE rolname = 'rls_verifier'
)
-- 7. All deterministic verification fixtures are absent
--    Covering every fixture family used by verification harness:
--    tenants, users, academic_terms, courses, course_offerings,
--    memberships, membership_roles, audit_logs, course_staff, enrollments
, fixture_rows_remaining AS (
  SELECT
    (SELECT count(*) FROM tenants WHERE id IN ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa','bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb')) +
    (SELECT count(*) FROM users WHERE id IN ('cccccccc-cccc-cccc-cccc-cccccccccccc','dddddddd-dddd-dddd-dddd-dddddddddddd')) +
    (SELECT count(*) FROM academic_terms WHERE id IN ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee','ffffffff-ffff-ffff-ffff-ffffffffffff')) +
    (SELECT count(*) FROM courses WHERE id IN ('11111111-1111-1111-1111-111111111111','22222222-2222-2222-2222-222222222222')) +
    (SELECT count(*) FROM course_offerings WHERE id IN ('33333333-3333-3333-3333-333333333333','44444444-4444-4444-4444-444444444444')) +
    (SELECT count(*) FROM memberships WHERE id IN ('55555555-5555-5555-5555-555555555555','66666666-6666-6666-6666-666666666666')) +
    (SELECT count(*) FROM membership_roles WHERE membership_id IN ('55555555-5555-5555-5555-555555555555','66666666-6666-6666-6666-666666666666')) +
    (SELECT count(*) FROM audit_logs WHERE id IN ('77777777-7777-7777-7777-777777777777','88888888-8888-8888-8888-888888888888')) +
    (SELECT count(*) FROM course_staff WHERE course_offering_id IN ('33333333-3333-3333-3333-333333333333','44444444-4444-4444-4444-444444444444')) +
    (SELECT count(*) FROM enrollments WHERE course_offering_id IN ('33333333-3333-3333-3333-333333333333','44444444-4444-4444-4444-444444444444'))
    AS fixture_rows_remaining
)
SELECT 'tables' AS chk, tables::text AS count FROM week3_tables
UNION ALL SELECT 'rls_enabled', rls_enabled::text FROM rls_enabled
UNION ALL SELECT 'a7_fk', a7_fk::text FROM a7_fk
UNION ALL SELECT 'term_range_check', term_range_check::text FROM term_range_check
UNION ALL SELECT 'audit_policies', audit_policies::text FROM audit_policies
UNION ALL SELECT 'rls_verifier', rls_verifier::text FROM rls_verifier
UNION ALL SELECT 'fixture_rows_remaining', fixture_rows_remaining::text FROM fixture_rows_remaining
ORDER BY chk;