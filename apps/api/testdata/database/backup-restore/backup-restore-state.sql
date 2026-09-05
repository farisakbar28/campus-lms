-- Read-only normalized state snapshot for the local backup drill.
-- The output is compared byte-for-byte between source and restored databases.

SELECT 'state_format=week03-backup-restore-v1';

WITH expected(tablename) AS (
    VALUES
        ('tenants'),
        ('users'),
        ('auth_identities'),
        ('memberships'),
        ('membership_roles'),
        ('audit_logs'),
        ('academic_terms'),
        ('courses'),
        ('course_offerings'),
        ('course_staff'),
        ('enrollments')
)
SELECT 'week3_table_count=' || count(*)::text
FROM pg_tables AS actual
JOIN expected ON expected.tablename = actual.tablename
WHERE actual.schemaname = 'public';

SELECT 'public_application_table_count=' || count(*)::text
FROM pg_tables
WHERE schemaname = 'public'
  AND tablename <> 'schema_migrations';

WITH expected(tablename) AS (
    VALUES
        ('tenants'),
        ('users'),
        ('auth_identities'),
        ('memberships'),
        ('membership_roles'),
        ('audit_logs'),
        ('academic_terms'),
        ('courses'),
        ('course_offerings'),
        ('course_staff'),
        ('enrollments')
)
SELECT 'week3_table_presence=' || string_agg(
    CASE WHEN actual.tablename IS NULL
         THEN 'missing:' || expected.tablename
         ELSE 'present:' || expected.tablename
    END,
    ',' ORDER BY expected.tablename
)
FROM expected
LEFT JOIN pg_tables AS actual
  ON actual.schemaname = 'public'
 AND actual.tablename = expected.tablename;

SELECT 'schema_migrations=' || CASE
    WHEN to_regclass('public.schema_migrations') IS NULL THEN 'absent'
    ELSE 'present'
END;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NULL THEN
        RAISE EXCEPTION 'schema_migrations is missing';
    END IF;
END
$$;

SELECT 'schema_migrations_rows=' || count(*)::text
FROM schema_migrations;

SELECT 'migration_version=' || COALESCE(
    string_agg(version::text, ',' ORDER BY version),
    '<missing>'
)
FROM schema_migrations;

SELECT 'migration_dirty=' || COALESCE(
    string_agg(CASE WHEN dirty THEN 'true' ELSE 'false' END, ',' ORDER BY version),
    '<missing>'
)
FROM schema_migrations;

WITH expected(tablename) AS (
    VALUES
        ('memberships'),
        ('membership_roles'),
        ('audit_logs'),
        ('academic_terms'),
        ('courses'),
        ('course_offerings'),
        ('course_staff'),
        ('enrollments')
), actual AS (
    SELECT c.relname AS tablename, c.relrowsecurity
    FROM pg_class AS c
    JOIN pg_namespace AS n ON n.oid = c.relnamespace
    WHERE n.nspname = 'public'
)
SELECT 'rls_state=' || string_agg(
    expected.tablename || ':' || CASE
        WHEN actual.relrowsecurity THEN 'enabled'
        ELSE 'disabled'
    END,
    ',' ORDER BY expected.tablename
)
FROM expected
LEFT JOIN actual ON actual.tablename = expected.tablename;

SELECT 'rls_enabled_count=' || count(*)::text
FROM pg_class AS c
JOIN pg_namespace AS n ON n.oid = c.relnamespace
WHERE n.nspname = 'public'
  AND c.relname IN (
      'memberships', 'membership_roles', 'audit_logs', 'academic_terms',
      'courses', 'course_offerings', 'course_staff', 'enrollments'
  )
  AND c.relrowsecurity = true;

SELECT 'a7_fk_count=' || count(*)::text
FROM pg_constraint AS c
WHERE c.conname = 'enrollments_tenant_id_student_user_id_fkey'
  AND c.contype = 'f'
  AND c.conrelid = 'public.enrollments'::regclass
  AND c.confrelid = 'public.memberships'::regclass
  AND (
      SELECT array_agg(local_column.attname ORDER BY pair.ordinality)
      FROM unnest(c.conkey, c.confkey) WITH ORDINALITY AS pair(local_attnum, referenced_attnum, ordinality)
      JOIN pg_attribute AS local_column
        ON local_column.attrelid = c.conrelid
       AND local_column.attnum = pair.local_attnum
  ) = ARRAY['tenant_id'::name, 'student_user_id'::name]
  AND (
      SELECT array_agg(referenced_column.attname ORDER BY pair.ordinality)
      FROM unnest(c.conkey, c.confkey) WITH ORDINALITY AS pair(local_attnum, referenced_attnum, ordinality)
      JOIN pg_attribute AS referenced_column
        ON referenced_column.attrelid = c.confrelid
       AND referenced_column.attnum = pair.referenced_attnum
  ) = ARRAY['tenant_id'::name, 'user_id'::name];

SELECT 'term_range_check_count=' || count(*)::text
FROM pg_constraint AS c
WHERE c.conname = 'academic_terms_valid_time_range'
  AND c.contype = 'c'
  AND c.conrelid = 'public.academic_terms'::regclass;

SELECT 'audit_policy_total=' || count(*)::text
FROM pg_policy
WHERE polrelid = 'public.audit_logs'::regclass;

SELECT 'audit_policy_select=' || count(*)::text
FROM pg_policy
WHERE polrelid = 'public.audit_logs'::regclass
  AND polcmd = 'r';

SELECT 'audit_policy_insert=' || count(*)::text
FROM pg_policy
WHERE polrelid = 'public.audit_logs'::regclass
  AND polcmd = 'a';

SELECT 'audit_policy_update=' || count(*)::text
FROM pg_policy
WHERE polrelid = 'public.audit_logs'::regclass
  AND polcmd = 'w';

SELECT 'audit_policy_delete=' || count(*)::text
FROM pg_policy
WHERE polrelid = 'public.audit_logs'::regclass
  AND polcmd = 'd';

SELECT 'audit_policy_all=' || count(*)::text
FROM pg_policy
WHERE polrelid = 'public.audit_logs'::regclass
  AND polcmd = '*';

SELECT 'enrollments_active_student_lookup_index_count=' || count(*)::text
FROM pg_class
WHERE oid = to_regclass('public.enrollments_active_student_lookup_idx');

SELECT 'enrollments_active_student_lookup_index_definition=' || COALESCE(
    (
        SELECT pg_get_indexdef(indexrelid)
        FROM pg_index
        WHERE indexrelid = to_regclass('public.enrollments_active_student_lookup_idx')
    ),
    '<missing>'
);

SELECT 'table_rows:tenants=' || count(*)::text FROM public.tenants;
SELECT 'table_fingerprint:tenants=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.tenants AS t;

SELECT 'table_rows:users=' || count(*)::text FROM public.users;
SELECT 'table_fingerprint:users=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.users AS t;

SELECT 'table_rows:auth_identities=' || count(*)::text FROM public.auth_identities;
SELECT 'table_fingerprint:auth_identities=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.auth_identities AS t;

SELECT 'table_rows:memberships=' || count(*)::text FROM public.memberships;
SELECT 'table_fingerprint:memberships=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.memberships AS t;

SELECT 'table_rows:membership_roles=' || count(*)::text FROM public.membership_roles;
SELECT 'table_fingerprint:membership_roles=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.membership_roles AS t;

SELECT 'table_rows:audit_logs=' || count(*)::text FROM public.audit_logs;
SELECT 'table_fingerprint:audit_logs=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.audit_logs AS t;

SELECT 'table_rows:academic_terms=' || count(*)::text FROM public.academic_terms;
SELECT 'table_fingerprint:academic_terms=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.academic_terms AS t;

SELECT 'table_rows:courses=' || count(*)::text FROM public.courses;
SELECT 'table_fingerprint:courses=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.courses AS t;

SELECT 'table_rows:course_offerings=' || count(*)::text FROM public.course_offerings;
SELECT 'table_fingerprint:course_offerings=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.course_offerings AS t;

SELECT 'table_rows:course_staff=' || count(*)::text FROM public.course_staff;
SELECT 'table_fingerprint:course_staff=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.course_staff AS t;

SELECT 'table_rows:enrollments=' || count(*)::text FROM public.enrollments;
SELECT 'table_fingerprint:enrollments=' || md5(COALESCE(string_agg(md5(row_to_json(t)::text), '' ORDER BY t.id), '')) FROM public.enrollments AS t;
