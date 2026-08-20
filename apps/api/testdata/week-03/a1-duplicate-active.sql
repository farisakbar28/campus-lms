-- A1-2: Attempt to insert the SAME active role again (should fail on membership_roles_active_role_idx)
INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_by, granted_at, revoked_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', 'lecturer', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NULL);