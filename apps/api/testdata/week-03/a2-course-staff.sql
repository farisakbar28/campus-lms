-- A2-C: Tenant A course_staff references Tenant B course_offering, with valid global user
-- Should fail on course_staff_tenant_id_course_offering_id_fkey
INSERT INTO course_staff (id, tenant_id, course_offering_id, user_id, role, source, permissions, active)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '44444444-4444-4444-4444-444444444444', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'instructor', 'local', '{}', true);