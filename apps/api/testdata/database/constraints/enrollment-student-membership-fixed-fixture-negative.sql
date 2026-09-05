-- An enrollment pairing a Tenant A offering with a user who belongs only to
-- Tenant B should fail on the composite student-membership foreign key.
INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at, withdrawn_at, synced_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '33333333-3333-3333-3333-333333333333', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'ext-enroll-invalid', 'active', NOW(), NULL, NOW());