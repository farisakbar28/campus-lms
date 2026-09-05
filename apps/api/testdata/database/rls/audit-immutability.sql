-- Grant UPDATE/DELETE table privileges for immutability test
GRANT UPDATE, DELETE ON audit_logs TO rls_verifier;

-- Attempt UPDATE own audit log as Tenant A
BEGIN;
SET LOCAL ROLE rls_verifier;
SET LOCAL app.tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

UPDATE audit_logs
SET reason = 'modified'
WHERE id = '77777777-7777-7777-7777-777777777777';

COMMIT;

-- Attempt DELETE own audit log as Tenant A
BEGIN;
SET LOCAL ROLE rls_verifier;
SET LOCAL app.tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

DELETE FROM audit_logs
WHERE id = '77777777-7777-7777-7777-777777777777';

COMMIT;

-- Verify audit log A still exists unchanged
SELECT id, reason FROM audit_logs WHERE id = '77777777-7777-7777-7777-777777777777';