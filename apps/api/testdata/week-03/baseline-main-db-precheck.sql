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
    IF t_count != 11 THEN RAISE EXCEPTION 'FAIL: Expected exactly 11 domain tables, found %', t_count; END IF;

    SELECT count(*) INTO t_count
    FROM pg_tables
    WHERE schemaname = 'public'
      AND tablename IN ('tenants', 'users', 'auth_identities', 'memberships', 'membership_roles', 'audit_logs', 'academic_terms', 'courses', 'course_offerings', 'course_staff', 'enrollments');
    IF t_count != 11 THEN RAISE EXCEPTION 'FAIL: Week 3 domain table set is incomplete'; END IF;

    SELECT count(*) INTO rls_count FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = 'public' AND c.relrowsecurity = true;
    IF rls_count != 8 THEN RAISE EXCEPTION 'FAIL: Expected 8 RLS-enabled tables, found %', rls_count; END IF;

    SELECT count(*) INTO rls_count
    FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE n.nspname = 'public'
      AND c.relrowsecurity = true
      AND c.relname IN ('memberships', 'membership_roles', 'audit_logs', 'academic_terms', 'courses', 'course_offerings', 'course_staff', 'enrollments');
    IF rls_count != 8 THEN RAISE EXCEPTION 'FAIL: RLS table set differs from authoritative Week 3 set'; END IF;

    SELECT count(*) INTO a7_count
    FROM pg_constraint c
    WHERE c.conname = 'enrollments_tenant_id_student_user_id_fkey'
      AND c.contype = 'f'
      AND c.conrelid = 'enrollments'::regclass
      AND c.confrelid = 'memberships'::regclass
      AND (SELECT array_agg(local_column.attname ORDER BY pair.ordinality)
           FROM unnest(c.conkey, c.confkey) WITH ORDINALITY AS pair(local_attnum, referenced_attnum, ordinality)
           JOIN pg_attribute local_column ON local_column.attrelid = c.conrelid AND local_column.attnum = pair.local_attnum) = ARRAY['tenant_id'::name, 'student_user_id'::name]
      AND (SELECT array_agg(referenced_column.attname ORDER BY pair.ordinality)
           FROM unnest(c.conkey, c.confkey) WITH ORDINALITY AS pair(local_attnum, referenced_attnum, ordinality)
           JOIN pg_attribute referenced_column ON referenced_column.attrelid = c.confrelid AND referenced_column.attnum = pair.referenced_attnum) = ARRAY['tenant_id'::name, 'user_id'::name];
    IF a7_count != 1 THEN RAISE EXCEPTION 'FAIL: A7 must map enrollments(tenant_id, student_user_id) to memberships(tenant_id, user_id)'; END IF;

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

    SELECT
        (SELECT count(*) FROM tenants WHERE id IN ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa','bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb')) +
        (SELECT count(*) FROM users WHERE id IN ('cccccccc-cccc-cccc-cccc-cccccccccccc','dddddddd-dddd-dddd-dddd-dddddddddddd')) +
        (SELECT count(*) FROM academic_terms WHERE id IN ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee','ffffffff-ffff-ffff-ffff-ffffffffffff')) +
        (SELECT count(*) FROM courses WHERE id IN ('11111111-1111-1111-1111-111111111111','22222222-2222-2222-2222-222222222222')) +
        (SELECT count(*) FROM course_offerings WHERE id IN ('33333333-3333-3333-3333-333333333333','44444444-4444-4444-4444-444444444444')) +
        (SELECT count(*) FROM memberships WHERE id IN ('55555555-5555-5555-5555-555555555555','66666666-6666-6666-6666-666666666666')) +
        (SELECT count(*) FROM membership_roles WHERE membership_id IN ('55555555-5555-5555-5555-555555555555','66666666-6666-6666-6666-666666666666')) +
        (SELECT count(*) FROM audit_logs WHERE id IN ('77777777-7777-7777-7777-777777777777','88888888-8888-8888-8888-888888888888')) +
        (SELECT count(*) FROM course_staff WHERE course_offering_id IN ('33333333-3333-3333-3333-333333333333','44444444-4444-4444-4444-444444444444')) +
        (SELECT count(*) FROM enrollments WHERE course_offering_id IN ('33333333-3333-3333-3333-333333333333','44444444-4444-4444-4444-444444444444'))
    INTO fixture_count;
    IF fixture_count != 0 THEN RAISE EXCEPTION 'FAIL: Deterministic verification fixtures remain: %', fixture_count; END IF;

END $$;
