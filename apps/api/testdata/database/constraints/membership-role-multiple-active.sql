-- Insert two distinct active roles for one membership.
-- Membership A: 55555555-5555-5555-5555-555555555555 (Tenant A + User A)
INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_by, granted_at, revoked_at)
VALUES
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', 'lecturer', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NULL),
  (gen_random_uuid(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555', 'student', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NULL);

SELECT 'A1-1 SUCCESS: Two distinct active roles inserted' as result;
SELECT id, role, revoked_at FROM membership_roles WHERE membership_id = '55555555-5555-5555-5555-555555555555' AND revoked_at IS NULL;