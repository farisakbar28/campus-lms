-- apps/api/testdata/week-03/verify-seed-counts.sql
-- Strictly validates exact counts and invariants for the Week 3 seed dataset.

WITH counts AS (
  SELECT
    (SELECT count(*) FROM tenants WHERE slug LIKE 'seed-%') AS tenants,
    (SELECT count(*) FROM users WHERE email LIKE 'seed-lecturer-%') AS lecturers,
    (SELECT count(*) FROM users WHERE email LIKE 'seed-student-%') AS students,
    (SELECT count(*) FROM users WHERE email LIKE '%@seed.campus-lms.invalid') AS users_total,
    (SELECT count(*) FROM auth_identities WHERE provider = 'seed_provider') AS auth_identities,
    (SELECT count(*) FROM memberships WHERE status = 'seed_active') AS memberships,
    (SELECT count(*) FROM membership_roles mr JOIN memberships m ON mr.membership_id = m.id WHERE m.status = 'seed_active') AS roles_total,
    (SELECT count(*) FROM membership_roles mr JOIN memberships m ON mr.membership_id = m.id WHERE m.status = 'seed_active' AND mr.role = 'lecturer') AS roles_lecturer,
    (SELECT count(*) FROM membership_roles mr JOIN memberships m ON mr.membership_id = m.id WHERE m.status = 'seed_active' AND mr.role = 'student') AS roles_student,
    (SELECT count(*) FROM academic_terms WHERE external_id LIKE 'seed-term-%') AS terms,
    (SELECT count(*) FROM courses WHERE external_id LIKE 'seed-course-%') AS courses,
    (SELECT count(*) FROM course_offerings WHERE external_id LIKE 'seed-offering-%') AS offerings,
    (SELECT count(*) FROM course_staff WHERE source = 'seed') AS course_staff_total,
    (SELECT count(*) FROM enrollments WHERE external_id LIKE 'seed-enrollment-%') AS enrollments,

    -- Consistency violations
    (SELECT count(*) FROM course_offerings o JOIN courses c ON o.course_id = c.id WHERE o.external_id LIKE 'seed-offering-%' AND o.tenant_id != c.tenant_id) AS cross_tenant_course,
    (SELECT count(*) FROM course_offerings o JOIN academic_terms t ON o.academic_term_id = t.id WHERE o.external_id LIKE 'seed-offering-%' AND o.tenant_id != t.tenant_id) AS cross_tenant_term,
    (SELECT count(*) FROM academic_terms WHERE external_id LIKE 'seed-term-%' AND starts_at >= ends_at) AS invalid_term_range,
    (SELECT count(*) FROM enrollments e JOIN course_offerings o ON e.course_offering_id = o.id WHERE e.external_id LIKE 'seed-enrollment-%' AND e.tenant_id != o.tenant_id) AS cross_tenant_enrollment,
    (SELECT count(*) FROM enrollments e WHERE e.external_id LIKE 'seed-enrollment-%' AND NOT EXISTS (SELECT 1 FROM memberships m WHERE m.tenant_id = e.tenant_id AND m.user_id = e.student_user_id)) AS missing_membership,
    (SELECT count(*) FROM enrollments e JOIN memberships m ON e.tenant_id = m.tenant_id AND e.student_user_id = m.user_id WHERE e.external_id LIKE 'seed-enrollment-%' AND m.status != 'seed_active') AS inactive_membership,
    (SELECT count(*) FROM enrollments e JOIN memberships m ON e.tenant_id = m.tenant_id AND e.student_user_id = m.user_id LEFT JOIN membership_roles mr ON mr.membership_id = m.id AND mr.role = 'student' WHERE e.external_id LIKE 'seed-enrollment-%' AND mr.id IS NULL) AS missing_student_role,
    (SELECT count(*) FROM (SELECT course_offering_id, student_user_id FROM enrollments WHERE external_id LIKE 'seed-enrollment-%' GROUP BY course_offering_id, student_user_id HAVING count(*) > 1) d) AS duplicate_enrollment,
    (SELECT count(*) FROM course_staff cs JOIN course_offerings o ON cs.course_offering_id = o.id WHERE cs.source = 'seed' AND cs.tenant_id != o.tenant_id) AS cross_tenant_staff,
    (SELECT count(*) FROM course_staff cs WHERE cs.source = 'seed' AND NOT EXISTS (SELECT 1 FROM memberships m WHERE m.tenant_id = cs.tenant_id AND m.user_id = cs.user_id)) AS missing_staff_membership,
    (SELECT count(*) FROM course_staff cs JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id WHERE cs.source = 'seed' AND m.status != 'seed_active') AS inactive_staff_membership,
    (SELECT count(*) FROM course_staff cs JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id LEFT JOIN membership_roles mr ON mr.membership_id = m.id AND mr.role = 'lecturer' WHERE cs.source = 'seed' AND mr.id IS NULL) AS missing_lecturer_role
)
SELECT
  'tenants=' || tenants ||
  ' lecturers=' || lecturers ||
  ' students=' || students ||
  ' users=' || users_total ||
  ' auth_identities=' || auth_identities ||
  ' memberships=' || memberships ||
  ' roles_total=' || roles_total ||
  ' roles_lecturer=' || roles_lecturer ||
  ' roles_student=' || roles_student ||
  ' terms=' || terms ||
  ' courses=' || courses ||
  ' offerings=' || offerings ||
  ' course_staff=' || course_staff_total ||
  ' enrollments=' || enrollments ||
  ' v_ct_course=' || cross_tenant_course ||
  ' v_ct_term=' || cross_tenant_term ||
  ' v_term_range=' || invalid_term_range ||
  ' v_ct_enrollment=' || cross_tenant_enrollment ||
  ' v_miss_mem=' || missing_membership ||
  ' v_inact_mem=' || inactive_membership ||
  ' v_miss_role=' || missing_student_role ||
  ' v_dup_enroll=' || duplicate_enrollment ||
  ' v_ct_staff=' || cross_tenant_staff ||
  ' v_miss_staff_mem=' || missing_staff_membership ||
  ' v_inact_staff_mem=' || inactive_staff_membership ||
  ' v_miss_lec_role=' || missing_lecturer_role
FROM counts;