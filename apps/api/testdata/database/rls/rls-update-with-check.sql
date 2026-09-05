-- RLS UPDATE WITH CHECK: As Tenant A target its own membership A but change tenant_id to Tenant B
BEGIN;
SET LOCAL ROLE rls_verifier;
SET LOCAL app.tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

UPDATE memberships
SET tenant_id = 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
WHERE id = '55555555-5555-5555-5555-555555555555';

COMMIT;