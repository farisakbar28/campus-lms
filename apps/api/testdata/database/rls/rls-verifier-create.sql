-- Create rls_verifier role
CREATE ROLE rls_verifier WITH NOLOGIN NOSUPERUSER NOBYPASSRLS;

-- Grant schema usage
GRANT USAGE ON SCHEMA public TO rls_verifier;

-- Grant minimal table privileges for verification
GRANT SELECT, UPDATE ON memberships TO rls_verifier;
GRANT SELECT, INSERT ON audit_logs TO rls_verifier;

-- Verify role properties
SELECT rolname, rolcanlogin, rolsuper, rolbypassrls
FROM pg_roles
WHERE rolname = 'rls_verifier';

-- Verify rls_verifier is NOT owner of audit_logs or memberships
SELECT c.relname, pg_get_userbyid(c.relowner) as owner
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'public'
  AND c.relname IN ('audit_logs', 'memberships');