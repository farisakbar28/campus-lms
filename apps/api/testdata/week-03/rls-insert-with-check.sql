-- RLS INSERT WITH CHECK: As Tenant A attempt to INSERT audit log with tenant_id Tenant B
BEGIN;
SET LOCAL ROLE rls_verifier;
SET LOCAL app.tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

INSERT INTO audit_logs (id, tenant_id, actor_user_id, actor_role, action, entity_type, entity_id, course_offering_id, before_data, after_data, reason, ip_address, request_id, occurred_at)
VALUES
  (gen_random_uuid(), 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'lecturer', 'CREATE', 'test', gen_random_uuid(), '33333333-3333-3333-3333-333333333333', '{}', '{}', 'test', '127.0.0.1', 'req-rls', NOW());

COMMIT;