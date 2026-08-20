-- apps/api/testdata/week-03/baseline-main-db-precheck.sql
-- Fails non-zero if current DB is not structurally identical to v4 end state.

DO $$
DECLARE
    t_count int;
    rls_count int;
    a7_count int;
    term_check_count int;
    audit_total int;
    audit_select int;
    audit_insert int;
    audit_update int;
    audit_delete int;
    audit_all int;
    rls_verifier_exists boolean;
    fixture_count int;
BEGIN
    SELECT count(*) INTO t_count FROM pg_tables WHERE schemaname = 'public' AND tablename != 'schema_migrations';
    IF t_count != 11 THEN RAISE EXCEPTION 'FAIL: Expected 11 tables, found %', t_count; END IF;

    SELECT count(*) INTO rls_count FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = 'public' AND c.relrowsecurity = true;
    IF rls_count != 8 THEN RAISE EXCEPTION 'FAIL: Expected 8 RLS-enabled tables, found %', rls_count; END IF;

    SELECT count(*) INTO a7_count FROM pg_constraint c JOIN pg_class rel ON rel.oid = c.conrelid WHERE rel.relname = 'enrollments' AND c.conname = 'enrollments_tenant_id_student_user_id_fkey';
    IF a7_count != 1 THEN RAISE EXCEPTION 'FAIL: A7 FK missing'; END IF;

    SELECT count(*) INTO term_check_count FROM pg_constraint c JOIN pg_class rel ON rel.oid = c.conrelid WHERE rel.relname = 'academic_terms' AND c.conname = 'academic_terms_valid_time_range';
    IF term_check_count != 1 THEN RAISE EXCEPTION 'FAIL: academic_terms_valid_time_range missing'; END IF;

    SELECT count(*) INTO audit_total FROM pg_policy p JOIN pg_class c ON p.polrelid = c.oid WHERE c.relname = 'audit_logs';
    IF audit_total != 2 THEN RAISE EXCEPTION 'FAIL: Expected 2 audit policies, found %', audit_total; END IF;

    SELECT count(*) INTO audit_select FROM pg_policy p JOIN pg_class c ON p.polrelid = c.oid WHERE c.relname = 'audit_logs' AND p.polcmd = 'r';
    IF audit_select != 1 THEN RAISE EXCEPTION 'FAIL: Expected 1 SELECT audit policy, found %', audit_select; END IF;

    SELECT count(*) INTO audit_insert FROM pg_policy p JOIN pg_class c ON p.polrelid = c.oid WHERE c.relname = 'audit_logs' AND p.polcmd = 'a';
    IF audit_insert != 1 THEN RAISE EXCEPTION 'FAIL: Expected 1 INSERT audit policy, found %', audit_insert; END IF;

    SELECT count(*) INTO audit_update FROM pg_policy p JOIN pg_class c ON p.polrelid = c.oid WHERE c.relname = 'audit_logs' AND p.polcmd = 'w';
    IF audit_update != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 UPDATE audit policy, found %', audit_update; END IF;

    SELECT count(*) INTO audit_delete FROM pg_policy p JOIN pg_class c ON p.polrelid = c.oid WHERE c.relname = 'audit_logs' AND p.polcmd = 'd';
    IF audit_delete != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 DELETE audit policy, found %', audit_delete; END IF;
    
    SELECT count(*) INTO audit_all FROM pg_policy p JOIN pg_class c ON p.polrelid = c.oid WHERE c.relname = 'audit_logs' AND p.polcmd = '*';
    IF audit_all != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 ALL audit policy, found %', audit_all; END IF;

    SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = 'rls_verifier') INTO rls_verifier_exists;
    IF rls_verifier_exists THEN RAISE EXCEPTION 'FAIL: rls_verifier role exists'; END IF;

    SELECT count(*) INTO fixture_count FROM users WHERE email IN ('tenant1-admin@test.com', 'tenant1-user@test.com', 'tenant2-admin@test.com', 'global-admin@test.com');
    IF fixture_count != 0 THEN RAISE EXCEPTION 'FAIL: Fixture users exist'; END IF;

END $$;
