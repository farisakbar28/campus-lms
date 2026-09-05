-- A cross-tenant audit_log/course_offering reference with a valid actor user
-- should fail on the composite foreign-key constraint.
INSERT INTO audit_logs (id, tenant_id, actor_user_id, actor_role, action, entity_type, entity_id, course_offering_id, before_data, after_data, reason, ip_address, request_id, occurred_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'lecturer', 'UPDATE', 'course_offering', '44444444-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444444', '{}', '{}', 'test', '127.0.0.1', 'req-cross', NOW());