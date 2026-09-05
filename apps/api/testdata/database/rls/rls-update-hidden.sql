-- RLS UPDATE USING: As Tenant A target the existing Tenant B membership (should update 0 rows)
BEGIN;
SET LOCAL ROLE rls_verifier;
SET LOCAL app.tenant_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

UPDATE memberships
SET status = 'suspended'
WHERE id = '66666666-6666-6666-6666-666666666666';

COMMIT;