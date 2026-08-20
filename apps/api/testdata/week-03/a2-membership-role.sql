-- A2-B: Tenant A membership_role references Tenant B membership
-- Should fail on membership_roles_tenant_id_membership_id_fkey
INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_by, granted_at, revoked_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '66666666-6666-6666-6666-666666666666', 'lecturer', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NULL);