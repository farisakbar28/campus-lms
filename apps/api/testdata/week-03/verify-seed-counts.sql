-- apps/api/testdata/week-03/verify-seed-counts.sql
-- Fails non-zero if ANY seeded invariant is violated.

DO $$
DECLARE
    v_tenants int;
    v_lecturers int;
    v_students int;
    v_users int;
    v_auth_identities int;
    v_memberships int;
    v_active_membership_roles int;
    v_active_lecturer_roles int;
    v_active_student_roles int;
    v_terms int;
    v_courses int;
    v_offerings int;
    v_course_staff int;
    v_enrollments int;

    v_offering_course_mismatch int;
    v_offering_term_mismatch int;
    v_invalid_term_range int;
    v_enrollment_offering_mismatch int;
    v_missing_student_membership int;
    v_inactive_student_membership int;
    v_missing_active_student_role int;
    v_duplicate_offering_student int;
    v_course_staff_offering_mismatch int;
    v_missing_lecturer_membership int;
    v_inactive_lecturer_membership int;
    v_missing_active_lecturer_role int;
    v_inactive_course_staff int;
    v_wrong_course_staff_role int;
    v_seed_course_staff_ids uuid[];
BEGIN
    SELECT array_agg(cs.id) INTO v_seed_course_staff_ids
    FROM course_staff cs
    JOIN course_offerings o ON o.id = cs.course_offering_id
    JOIN courses c ON c.id = o.course_id
    JOIN academic_terms t ON t.id = o.academic_term_id
    JOIN tenants tenant ON tenant.id = cs.tenant_id
    WHERE tenant.slug IN ('seed-campus-a', 'seed-campus-b', 'seed-campus-c')
      AND c.external_id LIKE 'seed-course-%'
      AND t.external_id LIKE 'seed-term-%'
      AND cs.id = md5('seed-staff-' || o.id)::uuid;

    SELECT count(*) INTO v_tenants FROM tenants WHERE slug IN ('seed-campus-a', 'seed-campus-b', 'seed-campus-c');
    IF v_tenants != 3 THEN RAISE EXCEPTION 'FAIL: Expected 3 tenants, got %', v_tenants; END IF;

    SELECT count(*) INTO v_lecturers FROM users WHERE email LIKE 'seed-lecturer-%';
    IF v_lecturers != 50 THEN RAISE EXCEPTION 'FAIL: Expected 50 lecturers, got %', v_lecturers; END IF;

    SELECT count(*) INTO v_students FROM users WHERE email LIKE 'seed-student-%';
    IF v_students != 2000 THEN RAISE EXCEPTION 'FAIL: Expected 2000 students, got %', v_students; END IF;

    SELECT count(*) INTO v_users FROM users WHERE email LIKE '%@seed.campus-lms.invalid';
    IF v_users != 2050 THEN RAISE EXCEPTION 'FAIL: Expected 2050 users, got %', v_users; END IF;

    SELECT count(*) INTO v_auth_identities FROM auth_identities WHERE provider = 'seed_provider';
    IF v_auth_identities != 2050 THEN RAISE EXCEPTION 'FAIL: Expected 2050 auth identities, got %', v_auth_identities; END IF;

    SELECT count(*) INTO v_memberships FROM memberships m JOIN users u ON m.user_id = u.id WHERE u.email LIKE '%@seed.campus-lms.invalid';
    IF v_memberships != 2050 THEN RAISE EXCEPTION 'FAIL: Expected 2050 memberships, got %', v_memberships; END IF;

    SELECT count(*) INTO v_active_membership_roles FROM membership_roles mr JOIN memberships m ON mr.membership_id = m.id JOIN users u ON m.user_id = u.id WHERE u.email LIKE '%@seed.campus-lms.invalid' AND mr.revoked_at IS NULL;
    IF v_active_membership_roles != 2050 THEN RAISE EXCEPTION 'FAIL: Expected 2050 active membership roles, got %', v_active_membership_roles; END IF;

    SELECT count(*) INTO v_active_lecturer_roles FROM membership_roles mr JOIN memberships m ON mr.membership_id = m.id JOIN users u ON m.user_id = u.id WHERE u.email LIKE 'seed-lecturer-%' AND mr.role = 'lecturer' AND mr.revoked_at IS NULL;
    IF v_active_lecturer_roles != 50 THEN RAISE EXCEPTION 'FAIL: Expected 50 active lecturer roles, got %', v_active_lecturer_roles; END IF;

    SELECT count(*) INTO v_active_student_roles FROM membership_roles mr JOIN memberships m ON mr.membership_id = m.id JOIN users u ON m.user_id = u.id WHERE u.email LIKE 'seed-student-%' AND mr.role = 'student' AND mr.revoked_at IS NULL;
    IF v_active_student_roles != 2000 THEN RAISE EXCEPTION 'FAIL: Expected 2000 active student roles, got %', v_active_student_roles; END IF;

    SELECT count(*) INTO v_terms FROM academic_terms WHERE external_id LIKE 'seed-term-%';
    IF v_terms != 6 THEN RAISE EXCEPTION 'FAIL: Expected 6 terms, got %', v_terms; END IF;

    SELECT count(*) INTO v_courses FROM courses WHERE external_id LIKE 'seed-course-%';
    IF v_courses != 200 THEN RAISE EXCEPTION 'FAIL: Expected 200 courses, got %', v_courses; END IF;

    SELECT count(*) INTO v_offerings FROM course_offerings WHERE external_id LIKE 'seed-offering-%';
    IF v_offerings != 400 THEN RAISE EXCEPTION 'FAIL: Expected 400 offerings, got %', v_offerings; END IF;

    SELECT count(*) INTO v_course_staff FROM course_staff WHERE id = ANY (v_seed_course_staff_ids);
    IF v_course_staff != 400 THEN RAISE EXCEPTION 'FAIL: Expected 400 course staff, got %', v_course_staff; END IF;

    SELECT count(*) INTO v_enrollments FROM enrollments WHERE external_id LIKE 'seed-enrollment-%';
    IF v_enrollments != 20000 THEN RAISE EXCEPTION 'FAIL: Expected 20000 enrollments, got %', v_enrollments; END IF;

    SELECT count(*) INTO v_offering_course_mismatch FROM course_offerings o JOIN courses c ON o.course_id = c.id WHERE o.external_id LIKE 'seed-offering-%' AND o.tenant_id != c.tenant_id;
    IF v_offering_course_mismatch != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 offering/course tenant mismatches, got %', v_offering_course_mismatch; END IF;

    SELECT count(*) INTO v_offering_term_mismatch FROM course_offerings o JOIN academic_terms t ON o.academic_term_id = t.id WHERE o.external_id LIKE 'seed-offering-%' AND o.tenant_id != t.tenant_id;
    IF v_offering_term_mismatch != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 offering/term tenant mismatches, got %', v_offering_term_mismatch; END IF;

    SELECT count(*) INTO v_invalid_term_range FROM academic_terms WHERE external_id LIKE 'seed-term-%' AND starts_at >= ends_at;
    IF v_invalid_term_range != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 invalid term ranges, got %', v_invalid_term_range; END IF;

    SELECT count(*) INTO v_enrollment_offering_mismatch FROM enrollments e JOIN course_offerings o ON e.course_offering_id = o.id WHERE e.external_id LIKE 'seed-enrollment-%' AND e.tenant_id != o.tenant_id;
    IF v_enrollment_offering_mismatch != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 enrollment/offering tenant mismatches, got %', v_enrollment_offering_mismatch; END IF;

    SELECT count(*) INTO v_missing_student_membership FROM enrollments e
        LEFT JOIN memberships m ON e.tenant_id = m.tenant_id AND e.student_user_id = m.user_id
        WHERE e.external_id LIKE 'seed-enrollment-%' AND m.id IS NULL;
    IF v_missing_student_membership != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 missing same-tenant student memberships, got %', v_missing_student_membership; END IF;

    SELECT count(*) INTO v_inactive_student_membership FROM enrollments e
        JOIN memberships m ON e.tenant_id = m.tenant_id AND e.student_user_id = m.user_id
        WHERE e.external_id LIKE 'seed-enrollment-%' AND m.status != 'active';
    IF v_inactive_student_membership != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 inactive student memberships, got %', v_inactive_student_membership; END IF;

    SELECT count(*) INTO v_missing_active_student_role FROM enrollments e
        JOIN memberships m ON e.tenant_id = m.tenant_id AND e.student_user_id = m.user_id
        LEFT JOIN membership_roles mr ON m.id = mr.membership_id AND mr.role = 'student' AND mr.revoked_at IS NULL
        WHERE e.external_id LIKE 'seed-enrollment-%' AND mr.id IS NULL;
    IF v_missing_active_student_role != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 missing ACTIVE student role, got %', v_missing_active_student_role; END IF;

    SELECT count(*) INTO v_duplicate_offering_student FROM (SELECT tenant_id, course_offering_id, student_user_id, count(*) FROM enrollments WHERE external_id LIKE 'seed-enrollment-%' GROUP BY 1,2,3 HAVING count(*) > 1) d;
    IF v_duplicate_offering_student != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 duplicate offering/student, got %', v_duplicate_offering_student; END IF;

    SELECT count(*) INTO v_course_staff_offering_mismatch FROM course_staff cs JOIN course_offerings o ON cs.course_offering_id = o.id WHERE cs.id = ANY (v_seed_course_staff_ids) AND cs.tenant_id != o.tenant_id;
    IF v_course_staff_offering_mismatch != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 course_staff/offering tenant mismatches, got %', v_course_staff_offering_mismatch; END IF;

    SELECT count(*) INTO v_missing_lecturer_membership FROM course_staff cs
        LEFT JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id
        WHERE cs.id = ANY (v_seed_course_staff_ids) AND m.id IS NULL;
    IF v_missing_lecturer_membership != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 missing same-tenant lecturer memberships, got %', v_missing_lecturer_membership; END IF;

    SELECT count(*) INTO v_inactive_lecturer_membership FROM course_staff cs
        JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id
        WHERE cs.id = ANY (v_seed_course_staff_ids) AND m.status != 'active';
    IF v_inactive_lecturer_membership != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 inactive lecturer memberships, got %', v_inactive_lecturer_membership; END IF;

    SELECT count(*) INTO v_missing_active_lecturer_role FROM course_staff cs
        JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id
        LEFT JOIN membership_roles mr ON m.id = mr.membership_id AND mr.role = 'lecturer' AND mr.revoked_at IS NULL
        WHERE cs.id = ANY (v_seed_course_staff_ids) AND mr.id IS NULL;
    IF v_missing_active_lecturer_role != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 missing ACTIVE lecturer role, got %', v_missing_active_lecturer_role; END IF;

    SELECT count(*) INTO v_inactive_course_staff FROM course_staff cs
        WHERE cs.id = ANY (v_seed_course_staff_ids) AND cs.active IS DISTINCT FROM true;
    IF v_inactive_course_staff != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 inactive seed course_staff rows, got %', v_inactive_course_staff; END IF;

    SELECT count(*) INTO v_wrong_course_staff_role FROM course_staff cs
        WHERE cs.id = ANY (v_seed_course_staff_ids) AND cs.role <> 'instructor';
    IF v_wrong_course_staff_role != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 seed course_staff rows with a non-instructor role, got %', v_wrong_course_staff_role; END IF;

END $$;

WITH seed_tenants AS (
    SELECT id FROM tenants WHERE slug IN ('seed-campus-a', 'seed-campus-b', 'seed-campus-c')
), seed_course_staff AS (
    SELECT cs.*
    FROM course_staff cs
    JOIN course_offerings o ON o.id = cs.course_offering_id
    JOIN courses c ON c.id = o.course_id
    JOIN academic_terms t ON t.id = o.academic_term_id
    JOIN seed_tenants st ON st.id = cs.tenant_id
    WHERE c.external_id LIKE 'seed-course-%'
      AND t.external_id LIKE 'seed-term-%'
      AND cs.id = md5('seed-staff-' || o.id)::uuid
), summary AS (
    SELECT
        (SELECT count(*) FROM seed_tenants) AS tenants,
        (SELECT count(*) FROM users WHERE email LIKE 'seed-lecturer-%') AS lecturers,
        (SELECT count(*) FROM users WHERE email LIKE 'seed-student-%') AS students,
        (SELECT count(*) FROM users WHERE email LIKE '%@seed.campus-lms.invalid') AS users,
        (SELECT count(*) FROM auth_identities ai JOIN users u ON u.id = ai.user_id WHERE ai.provider = 'seed_provider' AND u.email LIKE '%@seed.campus-lms.invalid') AS auth_identities,
        (SELECT count(*) FROM memberships m JOIN users u ON u.id = m.user_id JOIN seed_tenants st ON st.id = m.tenant_id WHERE u.email LIKE '%@seed.campus-lms.invalid') AS memberships,
        (SELECT count(*) FROM membership_roles mr JOIN memberships m ON m.id = mr.membership_id JOIN users u ON u.id = m.user_id JOIN seed_tenants st ON st.id = m.tenant_id WHERE u.email LIKE '%@seed.campus-lms.invalid' AND mr.revoked_at IS NULL) AS active_roles,
        (SELECT count(*) FROM membership_roles mr JOIN memberships m ON m.id = mr.membership_id JOIN users u ON u.id = m.user_id JOIN seed_tenants st ON st.id = m.tenant_id WHERE u.email LIKE 'seed-lecturer-%' AND mr.role = 'lecturer' AND mr.revoked_at IS NULL) AS lecturer_roles,
        (SELECT count(*) FROM membership_roles mr JOIN memberships m ON m.id = mr.membership_id JOIN users u ON u.id = m.user_id JOIN seed_tenants st ON st.id = m.tenant_id WHERE u.email LIKE 'seed-student-%' AND mr.role = 'student' AND mr.revoked_at IS NULL) AS student_roles,
        (SELECT count(*) FROM academic_terms t JOIN seed_tenants st ON st.id = t.tenant_id WHERE t.external_id LIKE 'seed-term-%') AS terms,
        (SELECT count(*) FROM courses c JOIN seed_tenants st ON st.id = c.tenant_id WHERE c.external_id LIKE 'seed-course-%') AS courses,
        (SELECT count(*) FROM course_offerings o JOIN seed_tenants st ON st.id = o.tenant_id WHERE o.external_id LIKE 'seed-offering-%') AS offerings,
        (SELECT count(*) FROM seed_course_staff) AS course_staff,
        (SELECT count(*) FROM enrollments e JOIN seed_tenants st ON st.id = e.tenant_id WHERE e.external_id LIKE 'seed-enrollment-%') AS enrollments,
        (SELECT count(*) FROM course_offerings o JOIN courses c ON c.id = o.course_id JOIN seed_tenants st ON st.id = o.tenant_id WHERE o.external_id LIKE 'seed-offering-%' AND o.tenant_id <> c.tenant_id) AS offering_course_mismatch,
        (SELECT count(*) FROM course_offerings o JOIN academic_terms t ON t.id = o.academic_term_id JOIN seed_tenants st ON st.id = o.tenant_id WHERE o.external_id LIKE 'seed-offering-%' AND o.tenant_id <> t.tenant_id) AS offering_term_mismatch,
        (SELECT count(*) FROM academic_terms t JOIN seed_tenants st ON st.id = t.tenant_id WHERE t.external_id LIKE 'seed-term-%' AND t.starts_at >= t.ends_at) AS invalid_term_range,
        (SELECT count(*) FROM enrollments e JOIN course_offerings o ON o.id = e.course_offering_id JOIN seed_tenants st ON st.id = e.tenant_id WHERE e.external_id LIKE 'seed-enrollment-%' AND e.tenant_id <> o.tenant_id) AS enrollment_offering_mismatch,
        (SELECT count(*) FROM enrollments e JOIN seed_tenants st ON st.id = e.tenant_id LEFT JOIN memberships m ON m.tenant_id = e.tenant_id AND m.user_id = e.student_user_id WHERE e.external_id LIKE 'seed-enrollment-%' AND m.id IS NULL) AS missing_student_membership,
        (SELECT count(*) FROM enrollments e JOIN seed_tenants st ON st.id = e.tenant_id JOIN memberships m ON m.tenant_id = e.tenant_id AND m.user_id = e.student_user_id WHERE e.external_id LIKE 'seed-enrollment-%' AND m.status <> 'active') AS inactive_student_membership,
        (SELECT count(*) FROM enrollments e JOIN seed_tenants st ON st.id = e.tenant_id JOIN memberships m ON m.tenant_id = e.tenant_id AND m.user_id = e.student_user_id LEFT JOIN membership_roles mr ON mr.membership_id = m.id AND mr.role = 'student' AND mr.revoked_at IS NULL WHERE e.external_id LIKE 'seed-enrollment-%' AND mr.id IS NULL) AS missing_active_student_role,
        (SELECT count(*) FROM (SELECT e.tenant_id, e.course_offering_id, e.student_user_id FROM enrollments e JOIN seed_tenants st ON st.id = e.tenant_id WHERE e.external_id LIKE 'seed-enrollment-%' GROUP BY e.tenant_id, e.course_offering_id, e.student_user_id HAVING count(*) > 1) duplicates) AS duplicate_offering_student,
        (SELECT count(*) FROM seed_course_staff cs JOIN course_offerings o ON o.id = cs.course_offering_id WHERE cs.tenant_id <> o.tenant_id) AS course_staff_offering_mismatch,
        (SELECT count(*) FROM seed_course_staff cs LEFT JOIN memberships m ON m.tenant_id = cs.tenant_id AND m.user_id = cs.user_id WHERE m.id IS NULL) AS missing_lecturer_membership,
        (SELECT count(*) FROM seed_course_staff cs JOIN memberships m ON m.tenant_id = cs.tenant_id AND m.user_id = cs.user_id WHERE m.status <> 'active') AS inactive_lecturer_membership,
        (SELECT count(*) FROM seed_course_staff cs JOIN memberships m ON m.tenant_id = cs.tenant_id AND m.user_id = cs.user_id LEFT JOIN membership_roles mr ON mr.membership_id = m.id AND mr.role = 'lecturer' AND mr.revoked_at IS NULL WHERE mr.id IS NULL) AS missing_active_lecturer_role,
        (SELECT count(*) FROM seed_course_staff WHERE active IS DISTINCT FROM true) AS inactive_course_staff,
        (SELECT count(*) FROM seed_course_staff WHERE role <> 'instructor') AS wrong_course_staff_role
)
SELECT line FROM summary CROSS JOIN LATERAL (VALUES
    ('tenants=' || tenants), ('lecturers=' || lecturers), ('students=' || students), ('users=' || users), ('auth_identities=' || auth_identities), ('memberships=' || memberships), ('active_roles=' || active_roles), ('lecturer_roles=' || lecturer_roles), ('student_roles=' || student_roles), ('terms=' || terms), ('courses=' || courses), ('offerings=' || offerings), ('course_staff=' || course_staff), ('enrollments=' || enrollments),
    ('violation_offering_course_mismatch=' || offering_course_mismatch), ('violation_offering_term_mismatch=' || offering_term_mismatch), ('violation_invalid_term_range=' || invalid_term_range), ('violation_enrollment_offering_mismatch=' || enrollment_offering_mismatch), ('violation_missing_student_membership=' || missing_student_membership), ('violation_inactive_student_membership=' || inactive_student_membership), ('violation_missing_active_student_role=' || missing_active_student_role), ('violation_duplicate_offering_student=' || duplicate_offering_student), ('violation_course_staff_offering_mismatch=' || course_staff_offering_mismatch), ('violation_missing_lecturer_membership=' || missing_lecturer_membership), ('violation_inactive_lecturer_membership=' || inactive_lecturer_membership), ('violation_missing_active_lecturer_role=' || missing_active_lecturer_role), ('violation_inactive_course_staff=' || inactive_course_staff), ('violation_wrong_course_staff_role=' || wrong_course_staff_role)
) AS output(line);
