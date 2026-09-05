-- A course offering referencing a course in another tenant, with a valid
-- same-tenant academic term, should fail on the composite foreign key.
INSERT INTO course_offerings (id, tenant_id, external_id, course_id, academic_term_id, external_section_code, display_name, lms_status, published_at, closed_at, archived_at, created_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ext-offering-cross', '22222222-2222-2222-2222-222222222222', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'X', 'Cross-tenant', 'draft', NULL, NULL, NULL, NOW());