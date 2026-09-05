-- An enrollment referencing a course offering in another tenant, with a
-- same-tenant student membership, should fail on the composite foreign key.
INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at, withdrawn_at, synced_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '44444444-4444-4444-4444-444444444444', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'ext-enroll-cross', 'active', NOW(), NULL, NOW());