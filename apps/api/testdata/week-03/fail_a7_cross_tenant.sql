-- apps/api/testdata/week-03/fail_a7_cross_tenant.sql
-- Proves that an enrollment cannot bypass the A7 tenant consistency invariant.

BEGIN;

-- Attempt to enroll a student from Campus B into a course offering in Campus A
INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at)
SELECT 
  gen_random_uuid(),
  o.tenant_id,
  o.id,
  (SELECT user_id FROM memberships WHERE tenant_id = md5('seed-tenant-b')::uuid LIMIT 1),
  'invalid-cross-tenant',
  'active',
  '2026-08-01 00:00:00+00'
FROM course_offerings o 
WHERE o.tenant_id = md5('seed-tenant-a')::uuid
LIMIT 1;

ROLLBACK;