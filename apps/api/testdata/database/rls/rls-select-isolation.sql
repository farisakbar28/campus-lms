-- RLS SELECT isolation test
BEGIN;
SET LOCAL ROLE rls_verifier;
SET LOCAL app.tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- Tenant A should see audit log A (1 row)
SELECT count(*) as tenant_a_sees_own FROM audit_logs WHERE id = '77777777-7777-7777-7777-777777777777';

-- Tenant A should NOT see audit log B (0 rows)
SELECT count(*) as tenant_a_sees_other FROM audit_logs WHERE id = '88888888-8888-8888-8888-888888888888';

COMMIT;