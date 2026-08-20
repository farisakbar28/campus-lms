-- A7 INVALID CASE: Tenant A course offering A + student_user_id = User B + User B only has membership in Tenant B
-- Should fail on enrollments_tenant_id_student_user_id_fkey
INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at, withdrawn_at, synced_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '33333333-3333-3333-3333-333333333333', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'ext-enroll-invalid', 'active', NOW(), NULL, NOW());