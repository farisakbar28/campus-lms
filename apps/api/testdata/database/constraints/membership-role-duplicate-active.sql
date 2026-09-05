-- Attempt to insert the same active role again; the active-role uniqueness
-- constraint should reject the statement.
INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_by, granted_at, revoked_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', 'lecturer', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NULL);