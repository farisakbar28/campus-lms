-- apps/api/testdata/week-03/seed.sql
-- Deterministic SQL seeder for Week 3.

BEGIN;

-- 1. Idempotency cleanup: delete seed-owned rows in safe dependency order
DELETE FROM audit_logs WHERE actor_role = 'seed_marker';
DELETE FROM enrollments WHERE id IN (SELECT id FROM enrollments WHERE external_id LIKE 'seed-enrollment-%');
DELETE FROM course_staff WHERE source = 'seed';
DELETE FROM course_offerings WHERE external_id LIKE 'seed-offering-%';
DELETE FROM courses WHERE external_id LIKE 'seed-course-%';
DELETE FROM academic_terms WHERE external_id LIKE 'seed-term-%';
DELETE FROM membership_roles WHERE id IN (SELECT mr.id FROM membership_roles mr JOIN memberships m ON mr.membership_id = m.id WHERE m.status = 'seed_active');
DELETE FROM memberships WHERE status = 'seed_active';
DELETE FROM auth_identities WHERE provider = 'seed_provider';
DELETE FROM users WHERE email LIKE '%@seed.campus-lms.invalid';
DELETE FROM tenants WHERE slug LIKE 'seed-%';

-- 2. Create deterministic tenants (3 tenants)
INSERT INTO tenants (id, slug, name, status, default_timezone, created_at)
VALUES
  (md5('seed-tenant-a')::uuid, 'seed-campus-a', 'Seed Campus A', 'active', 'Asia/Jakarta', '2026-08-01 00:00:00+00'),
  (md5('seed-tenant-b')::uuid, 'seed-campus-b', 'Seed Campus B', 'active', 'Asia/Jakarta', '2026-08-01 00:00:00+00'),
  (md5('seed-tenant-c')::uuid, 'seed-campus-c', 'Seed Campus C', 'active', 'Asia/Jakarta', '2026-08-01 00:00:00+00');

-- 3. Create global users (2050 total: 50 lecturers, 2000 students)
WITH lecturer_nums AS (SELECT generate_series(1, 50) AS n),
     student_nums AS (SELECT generate_series(1, 2000) AS n)
INSERT INTO users (id, email, display_name, status, created_at)
SELECT md5('seed-lecturer-' || n)::uuid, 'seed-lecturer-' || n || '@seed.campus-lms.invalid', 'Seed Lecturer ' || n, 'active', '2026-08-01 00:00:00+00'::timestamptz
FROM lecturer_nums
UNION ALL
SELECT md5('seed-student-' || n)::uuid, 'seed-student-' || n || '@seed.campus-lms.invalid', 'Seed Student ' || n, 'active', '2026-08-01 00:00:00+00'::timestamptz
FROM student_nums;

-- 4. Create auth identities
INSERT INTO auth_identities (id, user_id, provider, provider_subject, verified_at)
SELECT md5('seed-auth-' || id)::uuid, id, 'seed_provider', 'seed-sub-' || id, '2026-08-01 00:00:00+00'::timestamptz
FROM users WHERE email LIKE '%@seed.campus-lms.invalid';

-- 5. Create memberships (assigning deterministic tenants)
-- Lecturer tenants (distribute 50 lecturers across 3 tenants)
WITH tenants AS (
  SELECT id as tenant_id, row_number() over() - 1 as t_idx FROM tenants WHERE slug LIKE 'seed-%' ORDER BY slug
),
lecturer_members AS (
  SELECT u.id as user_id, (CAST(substring(u.email from 'seed-lecturer-([0-9]+)') AS int) - 1) % 3 as t_idx
  FROM users u WHERE u.email LIKE 'seed-lecturer-%'
)
INSERT INTO memberships (id, tenant_id, user_id, status, joined_at)
SELECT md5('seed-membership-' || l.user_id)::uuid, t.tenant_id, l.user_id, 'seed_active', '2026-08-01 00:00:00+00'::timestamptz
FROM lecturer_members l JOIN tenants t ON l.t_idx = t.t_idx;

-- Student tenants (distribute 2000 students across 3 tenants)
WITH tenants AS (
  SELECT id as tenant_id, row_number() over() - 1 as t_idx FROM tenants WHERE slug LIKE 'seed-%' ORDER BY slug
),
student_members AS (
  SELECT u.id as user_id, (CAST(substring(u.email from 'seed-student-([0-9]+)') AS int) - 1) % 3 as t_idx
  FROM users u WHERE u.email LIKE 'seed-student-%'
)
INSERT INTO memberships (id, tenant_id, user_id, status, joined_at)
SELECT md5('seed-membership-' || s.user_id)::uuid, t.tenant_id, s.user_id, 'seed_active', '2026-08-01 00:00:00+00'::timestamptz
FROM student_members s JOIN tenants t ON s.t_idx = t.t_idx;

-- 6. Create active membership roles
INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_at)
SELECT md5('seed-role-' || m.id)::uuid, m.tenant_id, m.id, 
  CASE WHEN u.email LIKE 'seed-lecturer-%' THEN 'lecturer' ELSE 'student' END,
  '2026-08-01 00:00:00+00'::timestamptz
FROM memberships m JOIN users u ON m.user_id = u.id
WHERE m.status = 'seed_active';

-- 7. Create academic terms (6 total: 2 per tenant)
WITH term_nums AS (SELECT generate_series(1, 2) AS n)
INSERT INTO academic_terms (id, tenant_id, external_id, code, name, starts_at, ends_at, status)
SELECT md5('seed-term-' || t.slug || '-' || n.n)::uuid, t.id, 'seed-term-' || t.slug || '-' || n.n, 'T' || n.n, 'Term ' || n.n,
  '2026-08-01 00:00:00+00'::timestamptz, '2026-12-01 00:00:00+00'::timestamptz, 'active'
FROM tenants t CROSS JOIN term_nums n WHERE t.slug LIKE 'seed-%';

-- 8. Create courses (200 courses across 3 tenants: ~67 per tenant)
WITH course_nums AS (SELECT generate_series(1, 200) AS n),
tenants AS (
  SELECT id as tenant_id, row_number() over() - 1 as t_idx FROM tenants WHERE slug LIKE 'seed-%' ORDER BY slug
)
INSERT INTO courses (id, tenant_id, external_id, code, name, credits, status)
SELECT md5('seed-course-' || c.n)::uuid, t.tenant_id, 'seed-course-' || c.n, 'C' || c.n, 'Seed Course ' || c.n, 3, 'active'
FROM course_nums c JOIN tenants t ON (c.n - 1) % 3 = t.t_idx;

-- 9. Create course offerings (400 offerings: 2 per course, one in each term)
WITH term_assign AS (
  SELECT id as term_id, tenant_id, row_number() over(partition by tenant_id ORDER BY code) as t_idx
  FROM academic_terms WHERE external_id LIKE 'seed-term-%'
)
INSERT INTO course_offerings (id, tenant_id, external_id, course_id, academic_term_id, lms_status, created_at)
SELECT md5('seed-offering-' || c.external_id || '-' || t.t_idx)::uuid, c.tenant_id, 'seed-offering-' || c.external_id || '-' || t.t_idx, c.id, t.term_id, 'published', '2026-08-01 00:00:00+00'::timestamptz
FROM courses c JOIN term_assign t ON c.tenant_id = t.tenant_id
WHERE c.external_id LIKE 'seed-course-%';

-- 10. Create course staff (assign 1 lecturer per offering, from same tenant)
WITH lecturers AS (
  SELECT m.user_id, m.tenant_id, row_number() over(partition by m.tenant_id ORDER BY m.user_id) as l_idx
  FROM memberships m JOIN users u ON m.user_id = u.id WHERE u.email LIKE 'seed-lecturer-%'
),
offerings AS (
  SELECT id as offering_id, tenant_id, row_number() over(partition by tenant_id ORDER BY external_id) as o_idx
  FROM course_offerings WHERE external_id LIKE 'seed-offering-%'
)
INSERT INTO course_staff (id, tenant_id, course_offering_id, user_id, role, source, active)
SELECT md5('seed-staff-' || o.offering_id)::uuid, o.tenant_id, o.offering_id, 
  -- pick lecturer by cycling through available lecturers for this tenant
  (SELECT l.user_id FROM lecturers l WHERE l.tenant_id = o.tenant_id AND l.l_idx = ((o.o_idx - 1) % (SELECT count(*) FROM lecturers WHERE tenant_id = o.tenant_id)) + 1),
  'instructor', 'seed', true
FROM offerings o;

-- 11. Create enrollments (20,000 enrollments: exactly 50 per offering, picking from same-tenant students)
WITH students AS (
  SELECT m.user_id, m.tenant_id, row_number() over(partition by m.tenant_id ORDER BY m.user_id) as s_idx
  FROM memberships m JOIN users u ON m.user_id = u.id WHERE u.email LIKE 'seed-student-%'
),
offerings AS (
  SELECT id as offering_id, tenant_id, row_number() over(partition by tenant_id ORDER BY external_id) as o_idx
  FROM course_offerings WHERE external_id LIKE 'seed-offering-%'
),
enrollment_matrix AS (
  -- For each offering, generate 50 rows
  SELECT o.offering_id, o.tenant_id, o.o_idx, g.n as slot
  FROM offerings o CROSS JOIN generate_series(1, 50) g(n)
)
INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at)
SELECT 
  md5('seed-enrollment-' || em.offering_id || '-' || em.slot)::uuid, 
  em.tenant_id, 
  em.offering_id, 
  -- map each slot to a distinct student in this tenant, shifting by offering index to mix them up
  (SELECT s.user_id FROM students s WHERE s.tenant_id = em.tenant_id AND s.s_idx = ((em.o_idx + em.slot - 2) % (SELECT count(*) FROM students WHERE tenant_id = em.tenant_id)) + 1),
  'seed-enrollment-' || em.offering_id || '-' || em.slot,
  'active', 
  '2026-08-01 00:00:00+00'::timestamptz
FROM enrollment_matrix em;

COMMIT;