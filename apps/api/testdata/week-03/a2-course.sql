-- A2-A: Tenant A course_offering references Tenant B course, with valid Tenant A academic term
-- Should fail on course_offerings_tenant_id_course_id_fkey
INSERT INTO course_offerings (id, tenant_id, external_id, course_id, academic_term_id, external_section_code, display_name, lms_status, published_at, closed_at, archived_at, created_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ext-offering-cross', '22222222-2222-2222-2222-222222222222', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'X', 'Cross-tenant', 'draft', NULL, NULL, NULL, NOW());