-- Cleanup verifier role
DROP OWNED BY rls_verifier;
DROP ROLE rls_verifier;

-- Verify role is gone
SELECT count(*) as role_count
FROM pg_roles
WHERE rolname = 'rls_verifier';