-- A cross-tenant course_staff/course_offering reference with a valid global
-- user should fail on the composite foreign-key constraint.
INSERT INTO course_staff (id, tenant_id, course_offering_id, user_id, role, source, permissions, active)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '44444444-4444-4444-4444-444444444444', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'instructor', 'local', '{}', true);