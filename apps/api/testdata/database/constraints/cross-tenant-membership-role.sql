-- A membership role referencing a membership in another tenant should fail
-- on the composite foreign-key constraint.
INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_by, granted_at, revoked_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '66666666-6666-6666-6666-666666666666', 'lecturer', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NULL);