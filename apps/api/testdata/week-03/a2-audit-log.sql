-- A2-D: Tenant A audit_log references Tenant B course_offering, with valid actor user
-- Should fail on audit_logs_tenant_id_course_offering_id_fkey
INSERT INTO audit_logs (id, tenant_id, actor_user_id, actor_role, action, entity_type, entity_id, course_offering_id, before_data, after_data, reason, ip_address, request_id, occurred_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'lecturer', 'UPDATE', 'course_offering', '44444444-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444444', '{}', '{}', 'test', '127.0.0.1', 'req-cross', NOW());