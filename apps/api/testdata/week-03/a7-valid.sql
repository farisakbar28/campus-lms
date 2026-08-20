-- A7 VALID CASE: Tenant A course offering A + student_user_id = User A + membership (Tenant A, User A) exists
INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at, withdrawn_at, synced_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '33333333-3333-3333-3333-333333333333', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'ext-enroll-valid', 'active', NOW(), NULL, NOW());

SELECT 'A7-1 SUCCESS: Valid enrollment inserted' as result;
SELECT id, tenant_id, course_offering_id, student_user_id FROM enrollments WHERE student_user_id = 'cccccccc-cccc-cccc-cccc-cccccccccccc';