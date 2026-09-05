-- RLS Negative Control: Prove RLS being disabled exposes cross-tenant data
BEGIN;
-- As owner/admin: disable RLS
ALTER TABLE audit_logs DISABLE ROW LEVEL SECURITY;

-- Now as rls_verifier/Tenant A
SET LOCAL ROLE rls_verifier;
SET LOCAL app.tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

-- Query the known Tenant B audit log - should be visible now
SELECT count(*) as tenant_b_visible_when_rls_disabled FROM audit_logs WHERE id = '88888888-8888-8888-8888-888888888888';

ROLLBACK;

-- After rollback, verify RLS is re-enabled
SELECT relname, relrowsecurity
FROM pg_class
WHERE relname = 'audit_logs';