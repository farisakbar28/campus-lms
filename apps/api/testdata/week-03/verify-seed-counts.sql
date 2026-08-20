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
BEGIN
    SELECT count(*) INTO v_tenants FROM tenants WHERE slug LIKE 'seed-%';
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

    SELECT count(*) INTO v_course_staff FROM course_staff WHERE source = 'seed';
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

    SELECT count(*) INTO v_course_staff_offering_mismatch FROM course_staff cs JOIN course_offerings o ON cs.course_offering_id = o.id WHERE cs.source = 'seed' AND cs.tenant_id != o.tenant_id;
    IF v_course_staff_offering_mismatch != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 course_staff/offering tenant mismatches, got %', v_course_staff_offering_mismatch; END IF;

    SELECT count(*) INTO v_missing_lecturer_membership FROM course_staff cs
        LEFT JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id
        WHERE cs.source = 'seed' AND m.id IS NULL;
    IF v_missing_lecturer_membership != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 missing same-tenant lecturer memberships, got %', v_missing_lecturer_membership; END IF;

    SELECT count(*) INTO v_inactive_lecturer_membership FROM course_staff cs
        JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id
        WHERE cs.source = 'seed' AND m.status != 'active';
    IF v_inactive_lecturer_membership != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 inactive lecturer memberships, got %', v_inactive_lecturer_membership; END IF;

    SELECT count(*) INTO v_missing_active_lecturer_role FROM course_staff cs
        JOIN memberships m ON cs.tenant_id = m.tenant_id AND cs.user_id = m.user_id
        LEFT JOIN membership_roles mr ON m.id = mr.membership_id AND mr.role = 'lecturer' AND mr.revoked_at IS NULL
        WHERE cs.source = 'seed' AND mr.id IS NULL;
    IF v_missing_active_lecturer_role != 0 THEN RAISE EXCEPTION 'FAIL: Expected 0 missing ACTIVE lecturer role, got %', v_missing_active_lecturer_role; END IF;

END $$;

SELECT 'Seed data verification passed: All invariants met.' AS result;