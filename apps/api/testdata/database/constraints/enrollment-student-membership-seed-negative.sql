-- Proves that an enrollment cannot bypass the A7 tenant consistency invariant.

BEGIN;

-- Attempt to enroll a student from Campus B into a course offering in Campus A
INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at)
SELECT 
  'ffffffff-ffff-ffff-ffff-ffffffffffff'::uuid,
  o.tenant_id,
  o.id,
  (
    SELECT m.user_id 
    FROM memberships m
    JOIN membership_roles mr ON m.id = mr.membership_id
    WHERE m.tenant_id = md5('seed-tenant-b')::uuid 
      AND m.status = 'active'
      AND mr.role = 'student'
      AND mr.revoked_at IS NULL
      AND NOT EXISTS (
          SELECT 1 FROM memberships m2 WHERE m2.user_id = m.user_id AND m2.tenant_id = md5('seed-tenant-a')::uuid
      )
    ORDER BY m.user_id 
    LIMIT 1
  ),
  'invalid-cross-tenant',
  'active',
  '2026-08-01 00:00:00+00'
FROM course_offerings o 
WHERE o.tenant_id = md5('seed-tenant-a')::uuid
ORDER BY o.id
LIMIT 1;

ROLLBACK;