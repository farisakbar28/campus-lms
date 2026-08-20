-- A1-3: Revoke the existing role by setting revoked_at, then insert that same role again
-- First, revoke the 'lecturer' role
UPDATE membership_roles
SET revoked_at = NOW()
WHERE membership_id = '55555555-5555-5555-5555-555555555555'
  AND role = 'lecturer'
  AND revoked_at IS NULL;

SELECT 'A1-3 SUCCESS: Role revoked via revoked_at' as result;

-- Now insert the same role again as a new active row
INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_by, granted_at, revoked_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', 'lecturer', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NULL);

SELECT 'A1-4 SUCCESS: Same role re-inserted after revocation' as result;
SELECT id, role, revoked_at FROM membership_roles WHERE membership_id = '55555555-5555-5555-5555-555555555555' ORDER BY granted_at;