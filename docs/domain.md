# Domain Model — campus-lms

## 0. Purpose and domain boundary

campus-lms is a multi-tenant Learning Management System (LMS) for higher
education institutions. It is not the institution's authoritative academic
information system (SIAKAD).

The LMS manages learning after the institution has supplied the relevant
academic facts: course delivery, RPS (Rencana Pembelajaran Semester, the
semester learning plan) and learning outcomes, materials,
activities, assignments, examinations, discussions, attendance, progress,
assessment, feedback, gradebook, and learning-activity reporting.

Academic master data remains owned by SIAKAD and is exchanged through an
integration adapter. The adapter protects the canonical LMS model from any
vendor or protocol-specific representation.

### 0.1 Sources of truth

SIAKAD is authoritative for:

- academic terms and semesters;
- academic-program references;
- the master course catalogue;
- course-offering identity and section/class identity;
- instructor assignments supplied by the institution;
- the students assigned to an offering; and
- student enrollment and KRS (the institution's course-registration record).

These records are read-only from the ordinary LMS side. The LMS must not become
an alternate path for changing academic facts that belong to SIAKAD.

The LMS is authoritative for:

- course content and publication state;
- the RPS/learning-plan versions maintained in the LMS;
- LMS-managed CPMK and Sub-CPMK versions;
- modules, lessons, materials, and learning activities;
- assignments and submissions;
- quizzes, examinations, question banks, and attempts;
- discussion forums and announcements;
- attendance;
- activity completion and learning progress;
- rubrics, gradebook entries, feedback, and final-grade workflow; and
- the LMS audit trail.

Attendance is LMS-owned. Its source-of-truth decision is:

ATTENDANCE_SOURCE_OF_TRUTH=LMS

Final grades and attendance may be sent to SIAKAD through an outbound
integration adapter after the LMS rules for publication or finalization have
been satisfied.

### 0.2 Outside the core LMS domain

The following remain outside the responsibility of campus-lms:

- creating or changing faculties, departments, or academic programs;
- admissions;
- academic registration and KRS administration;
- tuition and payment processing;
- room scheduling and classroom management;
- complete institutional curriculum management;
- transcripts, GPA, graduation, judicium (the formal graduation decision), or
  diplomas;
- lecturer employment administration;
- one-to-one direct messaging;
- active SPADA or PDDikti implementation in the current core scope.

The LMS may receive reference data from these domains through an integration
adapter.

### 0.3 Terms

- Tenant: one institution using the SaaS platform.
- User: a global platform login identity.
- Membership: the relationship between a user and one tenant.
- Membership role: a tenant-scoped role granted through a membership.
- Course: a master course, such as IF101 — Algorithms and Programming.
- Course offering: one delivery of a course in a term and section.
- Course staff: an instructor or Teaching Assistant assigned to an offering.
- Module: a major topic or unit in a course structure.
- Lesson: a learning unit or session within a module.
- Learning activity: an activity that a student may or must complete.
- Enrollment: a student's participation in a course offering.
- Assessment: an activity that produces or may produce a score.
- Grade item: one component of a gradebook.
- Grade: a student's result for one grade item.
- Final grade: the result of the gradebook for one offering and enrollment.
- Activity completion: the state of completing a learning activity; it is not
  attendance.
- Attendance: a student's attendance record for one attendance session.

---

## 1. Actors and roles

There are two authorization levels:

1. tenant roles, granted through memberships and membership_roles; and
2. course-offering roles, granted through course_staff.

A user may have one global identity, memberships in several tenants, several
active roles in one tenant, and different roles in different tenants. Every
authorization decision remains tenant-scoped.

### 1.1 Super Admin

Super Admin is a global SaaS operator, not an academic administrator of a
campus.

May:

- create, activate, suspend, or deactivate tenants;
- view tenant metadata required for SaaS operations;
- manage global platform configuration and feature capability;
- view global integration health;
- perform explicitly audited break-glass support access;
- manage tenant-administrator accounts when operationally necessary;
- view platform_audit_logs; and
- perform platform administration that does not alter academic records.

May not:

- become an instructor merely by holding Super Admin access;
- assign student grades or alter submissions, attendance, or enrollment;
- change course content without tenant authorization;
- read tenant data without a legitimate, audited operational purpose;
- disable audit trails; or
- alter tenant data without an audit record.

Every break-glass access must record the reason, actor, target tenant,
timestamp, action, and audit trail.

### 1.2 Campus Administrator

Campus Administrator administers the LMS for one tenant. This role is not an
operator of SIAKAD and does not automatically receive academic authority.

May:

- manage tenant LMS configuration, branding, timezone, and policy;
- manage memberships and tenant roles;
- configure storage and file-upload policy;
- configure tenant integrations and attendance methods;
- configure default grading schemes;
- view synchronization status and permitted tenant audit records;
- support course lifecycle operations;
- archive or unarchive offerings when authorized; and
- manage tenant feature flags.

May not:

- create or change academic entities owned by SIAKAD;
- bypass SIAKAD to alter enrollment;
- change grades, submissions, or final-grade publication;
- impersonate a student or instructor outside an audited support path;
- delete audit logs; or
- hard-delete learning history.

### 1.3 Academic Operator

Academic Operator manages the boundary between SIAKAD and the LMS.

May:

- run and monitor inbound and outbound synchronization;
- inspect synchronization jobs and errors;
- retry failed synchronization;
- map external identifiers and resolve mapping conflicts;
- inspect terms, courses, offerings, instructor assignments, and enrollments;
- validate synchronization results;
- send a final grade that has already been published to SIAKAD;
- send attendance through the outbound adapter;
- run or monitor authorized course-copy/import operations; and
- archive an offering according to tenant operating policy.

May not:

- edit SIAKAD-owned academic facts directly in the LMS;
- add or remove a student to bypass SIAKAD;
- change an SIAKAD-supplied instructor assignment;
- change scores, submissions, quiz answers, or published final grades;
- change attendance without a valid correction workflow; or
- use integration mappings to move objects between tenants.

### 1.4 Instructor and Lead Instructor

A tenant lecturer membership does not grant access to every course. Course
authority comes from course_staff.

Course-offering roles are:

- lead_instructor;
- instructor; and
- teaching_assistant.

Lead Instructor may perform all Instructor actions and may also:

- manage the learning structure and the offering's RPS;
- manage CPMK and Sub-CPMK mappings;
- configure gradebook, grading scheme, and assessment weights;
- publish grades and finalize final grades;
- lock final grades;
- reopen a locked final grade with an audited reason;
- determine Teaching Assistant permissions;
- publish the course to students; and
- close the learning process for the offering.

Instructor may:

- create and edit modules and lessons;
- upload materials;
- create assignments, quizzes, examinations, rubrics, forums, and
  announcements;
- use question banks within the permitted scope;
- create attendance sessions;
- view enrolled students and submissions;
- grade and provide feedback;
- configure individual overrides and group activities;
- view learning progress; and
- view course reports.

Instructor may not:

- edit SIAKAD enrollment, course identity, academic term, or instructor
  assignment;
- move a student between offerings;
- cross a tenant boundary;
- hard-delete grades, submissions, or audit history; or
- publish or lock a final grade unless also acting as Lead Instructor.

### 1.5 Teaching Assistant

Teaching Assistant is a course-scoped role, not a tenant-wide administrative
role. Permissions may be narrowed by course configuration.

When authorized, a TA may:

- view the course roster and materials;
- upload or edit materials;
- help manage lessons and discussions;
- create or manage attendance sessions;
- view submissions and learning progress;
- provide feedback;
- grade as a draft score; and
- assist with rubric assessment.

A TA may not:

- publish or lock final grades;
- change grade schemes or grade weights;
- change enrollment, course ownership, or Lead Instructor assignment;
- change SIAKAD-owned data;
- send final grades to SIAKAD;
- change a locked final grade; or
- hard-delete academic history.

Every TA score stores graded_by, draft status, timestamp, and rubric detail
when applicable. The responsible instructor remains accountable for
publication.

### 1.6 Student

May:

- view offerings in which the student is enrolled;
- view published course overview and RPS;
- view available modules, lessons, and materials;
- complete activities;
- submit assignments and files;
- take quizzes and examinations within the permitted window;
- view published feedback and grades;
- participate in permitted discussion forums;
- use self check-in, QR, or PIN attendance where enabled;
- view personal attendance and progress; and
- receive announcements.

May not:

- access an un-enrolled offering;
- view another student's submission, grade, or attendance, except through an
  explicitly authorized future peer activity;
- change grades, enrollment, deadlines, or attendance sessions;
- edit a locked submission;
- attempt a quiz outside its permitted window;
- use another student's override;
- access unpublished content or draft grades.

---

## 2. Identity and tenant isolation

### 2.1 Global identity

users is the global identity table. Tenant roles are never stored as a global
property of a user.

The relationship is:

users → memberships → tenants

A user may be a lecturer in one tenant, a student in another, and hold several
active tenant roles in the same tenant through membership_roles.

### 2.2 Tenant isolation

Every learning-domain row has a deterministic relationship to exactly one
tenant. A query must validate tenant context together with object identity.

An object identifier alone is never an authorization boundary. RLS must prevent
objects from one tenant being read, referenced, or modified from another
tenant.

Tenant context is obtained from a trusted authenticated principal and is set
only for the transaction performing tenant-scoped database work. Missing or
invalid context fails closed.

### 2.3 Global tables

These are global because they are not owned by one tenant:

- tenants;
- users;
- auth_identities;
- auth_sessions;
- platform_admins; and
- platform_audit_logs.

All other learning-domain records are tenant-scoped, directly or through a
tenant-owned parent, and must remain subject to tenant isolation.

---

## 3. Tenant and identity entities

### tenants

Global attributes include id, slug, name, status, default_timezone,
created_at, and suspended_at.

### tenant_settings

Tenant-owned settings include branding, locale, upload policy, attendance
policy, grading defaults, and feature flags.

### users

Global user attributes include id, email, display_name, status, and
created_at.

### auth_identities

An authentication identity records provider, provider subject, user_id, and
verification state.

### memberships

A membership records id, tenant_id, user_id, status, and joined_at.

There is exactly one membership for a user in a tenant:

UNIQUE (tenant_id, user_id)

### membership_roles

Tenant roles are separate records, not a singular role column on memberships.
Each record contains id, tenant_id, membership_id, role, granted_by,
granted_at, and revoked_at.

Active role grants are unique per membership and role. Revocation records
revoked_at rather than deleting the grant, preserving audit history. The
canonical tenant roles are:

- tenant_admin;
- academic_operator;
- lecturer; and
- student.

teaching_assistant remains course-scoped through course_staff.

### external_identities

An external identity maps tenant_id, user_id, source, external_user_id, and
external_type for an integration.

### auth_sessions

Global authentication sessions support short-lived access and revocable
refresh sessions. A session records id, user_id, refresh_token_hash, issued_at,
expires_at, rotated_from, revoked_at, revoked_reason, user_agent, ip_address,
and last_seen_at.

Refresh tokens are stored as hashes. Rotation creates a new row and links it
through rotated_from. Reuse of a rotated token is treated as a possible
compromise and may revoke all sessions for the user.

---

## 4. SIAKAD-authoritative academic references

These entities are tenant-scoped but externally authoritative and read-only
from the ordinary LMS side.

### academic_terms

Records id, tenant_id, external_id, code, name, starts_at, ends_at, status,
and synced_at. The owner is SIAKAD.

### academic_program_refs

Provides program references for filtering and reporting. SIAKAD owns the
external_id, code, name, and status.

### courses

The course master records id, tenant_id, external_id, code, name, credits,
academic_program_ref_id, status, and synced_at. A course is not a semester
class.

### course_offerings

An offering instantiates one course in one academic term and section. It
contains id, tenant_id, external_id, course_id, academic_term_id,
external_section_code, display_name, lms_status, published_at, closed_at,
archived_at, course_plan_version_id, and created_at.

SIAKAD owns course identity, term, and external section identity. The LMS owns
LMS status, publication, content, learning plan, and archive state.

The lifecycle is:

draft → published → active → closed → archived

An archived offering is read-only by default.

### course_staff

Records id, tenant_id, course_offering_id, user_id, role, source, permissions,
and active. The role is lead_instructor, instructor, or
teaching_assistant.

Lead Instructor and Instructor assignments normally originate in SIAKAD. A TA
may be supplied by an integration or assigned locally under tenant policy.

There is at most one staff record for a user in an offering:

UNIQUE (course_offering_id, user_id)

### enrollments

Records id, tenant_id, course_offering_id, student_user_id, external_id, status,
enrolled_at, withdrawn_at, and synced_at. SIAKAD owns enrollment.

There is at most one enrollment for a student in an offering:

UNIQUE (course_offering_id, student_user_id)

The enrollment also has a composite foreign-key relationship to the student's
membership in the same tenant. Service authorization still verifies that the
membership and enrollment are active.

Use statuses such as active, withdrawn, completed, and inactive. If a
synchronization no longer returns an enrollment that has learning activity,
the record is not hard-deleted.

---

## 5. RPS and learning outcomes

CPL (institutional learning outcomes) is an external reference and is
read-only. CPMK (course learning outcomes), Sub-CPMK (sub-course learning
outcomes), and the LMS learning plan are LMS-managed and versioned.

### learning_outcomes

Records id, tenant_id, course_id, type, code, title, description, external_id,
source, version, and active. type is CPL, CPMK, or SUB_CPMK.

### outcome_mappings

Maps CPL to CPMK and CPMK to Sub-CPMK through tenant_id, parent_outcome_id,
child_outcome_id, and an optional weight.

### course_plan_versions

An RPS/learning-plan version records id, tenant_id, course_id, version_number,
title, description, learning_guidance, status, created_by, published_by,
published_at, and created_at. Status is draft, published, or retired.

An active offering retains the version it used. A material change creates a
new version rather than overwriting historical content.

### course_plan_outcomes

Associates a plan version with the outcome versions it uses, preserving the
correct learning-outcome snapshot for historical offerings.

---

## 6. Learning structure

The instructional structure is:

Course
└── Course Offering
    ├── Course Plan / RPS
    ├── Module
    │   └── Lesson
    │       ├── Material
    │       └── Learning Activity
    └── Gradebook

### modules

Records id, tenant_id, course_offering_id, title, description, position,
status, available_from, available_until, and created_by. Status is draft,
published, or hidden.

A module may represent a week, topic, chapter, unit, or another learning
block. The domain does not require one module to equal one week.

### lessons

Records id, tenant_id, module_id, title, description, position, learning_mode,
estimated_minutes, available_from, available_until, and status.

learning_mode may be asynchronous, synchronous, blended, or onsite.

### lesson_outcomes

Maps a lesson to CPMK or Sub-CPMK.

---

## 7. Materials and files

### files

Tenant-owned file metadata includes id, tenant_id, storage_key,
original_filename, mime_type, size_bytes, checksum, uploaded_by,
malware_scan_status, and created_at. Binary objects are stored outside the
domain row, and every file has tenant ownership.

### materials

A material belongs to a lesson and records id, tenant_id, lesson_id, title,
description, type, file_id, external_url, content, position, published, and
created_by.

Supported material types include text, file, link, video, audio, embed,
learning_package, and external_tool. PDF is one possible format, not the only
permitted format.

---

## 8. Learning activities

### learning_activities

The base activity records id, tenant_id, lesson_id, type, title, description,
position, available_from, due_at, cutoff_at, published, completion_required,
and created_by.

Activity types include assignment, quiz, discussion, attendance, and
external_tool.

### activity_outcomes

Maps activities to CPMK or Sub-CPMK, optionally with a weight. One assessment
may measure more than one outcome.

---

## 9. Assignments and submissions

### assignments

Records id, tenant_id, learning_activity_id, instructions, submission_mode,
submission_types, max_attempts, max_score, group_mode, allow_resubmission,
require_submission_statement, and created_by.

Submission types may be text, file, or text_and_file. group_mode is individual
or group.

### assignment_overrides

An override targets exactly one enrollment or one group. It records id,
tenant_id, assignment_id, enrollment_id or group_id, availability and cutoff
times, max_attempts, reason, and created_by.

The invariant is:

enrollment_id XOR group_id

An override changes only the targeted student's or group's window.

### submissions

A logical submission records id, tenant_id, assignment_id, enrollment_id or
group_id, status, current_version, first_submitted_at, last_submitted_at, and
locked_at.

Status is draft, submitted, returned, resubmission_allowed, or locked.

### submission_versions

Every submission or resubmission creates an immutable version with id,
tenant_id, submission_id, attempt_number, text_content, submitted_at, and
submitted_by. Earlier versions are never overwritten.

### submission_files

Maps a submission version to tenant-owned files.

---

## 10. Group learning

### course_groups

Records id, tenant_id, course_offering_id, name, description, and created_by.

### course_group_members

Records tenant_id, group_id, enrollment_id, and joined_at. There is at most one
membership for an enrollment in a group:

UNIQUE (group_id, enrollment_id)

Groups may support group assignments, discussions, and collaborative
activities.

---

## 11. Quizzes and examinations

Quizzes and examinations use the same domain engine with different
configuration.

### question_banks

Records id, tenant_id, an optional course_id or course_offering_id, name,
description, visibility, and created_by. A bank may be offering-specific or
reusable within the same course. Cross-tenant sharing is prohibited unless a
future library explicitly authorizes it.

### questions

Records id, tenant_id, question_bank_id, type, created_by, and status.
Supported types include multiple_choice_single, multiple_choice_multiple,
true_false, short_answer, essay, and numeric. The model remains extensible.

### question_versions

Records id, tenant_id, question_id, version_number, prompt, explanation,
default_points, configuration, and created_at. A question version is
immutable.

### question_options

For selectable questions, records id, tenant_id, question_version_id, content,
position, is_correct, and score_fraction. Correct-answer metadata must not be
sent to a student before the feedback policy permits it.

### quizzes

Records id, tenant_id, learning_activity_id, mode, time_limit_seconds,
attempt_limit, shuffle_questions, shuffle_answers, grade_method, max_score,
feedback_release_policy, and created_by. mode is quiz or exam.

### quiz_question_rules

Defines fixed or random questions, category or pool selection, points, and
order. Random selection is resolved when an attempt starts and stored with the
attempt.

### quiz_overrides

Provides authorized student or group overrides for opening time, closing time,
attempt limit, and time limit.

### quiz_attempts

Records id, tenant_id, quiz_id, enrollment_id, attempt_number, started_at,
expires_at, submitted_at, status, and score.

Status is in_progress, submitted, auto_submitted, graded, or invalidated.

### quiz_attempt_questions

Stores the question-version snapshot assigned to an attempt. This preserves
randomization and examination history.

### quiz_responses

Records tenant_id, quiz_attempt_question_id, response, answered_at, auto_score,
manual_score, feedback, and graded_by. Essay and other manual responses may
await instructor or TA grading.

---

## 12. Rubrics

### rubrics

Records id, tenant_id, course_offering_id, name, description, status, and
created_by.

### rubric_criteria

Records tenant_id, id, rubric_id, title, description, weight, and position.

### rubric_levels

Records tenant_id, id, rubric_criterion_id, label, description, and score.

A rubric used for grading retains a version or snapshot suitable for
historical review.

---

## 13. Gradebook

The domain does not assume a 1–100 scale. A user interface may display that
range by default, but the model uses raw score, maximum score, weight,
normalized score, and a configurable grade scheme.

### grade_categories

Records id, tenant_id, course_offering_id, name, weight, and position.

### grade_items

A grade item records id, tenant_id, course_offering_id, grade_category_id,
name, max_score, weight, grading_type, rubric_id, and published.

Its source is represented by nullable typed references:

- assignment_id;
- quiz_id; or
- attendance_ref.

At most one source reference may be present. A manual item has all source
references null. This preserves referential integrity without an unbounded
polymorphic foreign key.

### grades

Records id, tenant_id, grade_item_id, enrollment_id, raw_score,
normalized_score, status, feedback, graded_by, graded_at, published_by, and
published_at.

Status is draft, published, or overridden. There is at most one grade for a
grade item and enrollment:

UNIQUE (grade_item_id, enrollment_id)

A TA may create or change only draft scores.

### grade_schemes

Records id, tenant_id, name, type, configuration, and active. A scheme may
produce labels such as A, AB, B, BC, C, D, or E, but no scheme is assumed to
apply to every institution.

### final_grades

Records id, tenant_id, course_offering_id, enrollment_id, numeric_result,
grade_label, status, calculated_at, published_by, published_at, locked_by,
locked_at, sync_status, and synced_at.

Status is draft, published, or locked. There is at most one final grade for an
offering and enrollment:

UNIQUE (course_offering_id, enrollment_id)

Only Lead Instructor may transition draft to published or published to locked.

---

## 14. Grade change and audit

Every grade or final-grade change produces an audit record. At minimum the
record includes actor, role, tenant, offering, student, grade item, previous
value, new value, reason for an override or reopen, timestamp, and request or
correlation ID.

A locked final grade cannot be edited directly. Correction follows:

locked
→ Lead Instructor requests reopen with a reason
→ audit
→ published or draft
→ correction
→ republish
→ lock
→ outbound synchronization

Academic Operator cannot alter a score as a shortcut for a synchronization
failure.

---

## 15. Attendance

Attendance is separate from activity completion and is owned by the LMS.

### attendance_sessions

Records id, tenant_id, course_offering_id, optional lesson_id or
learning_activity_id, title, method, opens_at, closes_at, late_after,
created_by, and status.

Methods include manual, self_check_in, pin, qr, activity_completion, and
external.

### attendance_credentials

QR and PIN values are check-in credentials, not attendance records. A
credential records id, tenant_id, attendance_session_id, type, credential_hash,
valid_from, valid_until, and revoked_at.

Credentials must be time-bounded, session-specific, non-reusable after the
session, rotatable or revocable, and unusable across tenants or courses. QR
tokens should be sufficiently random and short-lived. A PIN must not be stored
in plaintext where hashing is possible.

### attendance_records

Records id, tenant_id, attendance_session_id, enrollment_id, status,
check_in_method, checked_in_at, recorded_by, override_reason, and updated_at.

Status is present, late, absent, or excused. There is at most one record for a
student in a session:

UNIQUE (attendance_session_id, enrollment_id)

Manual correction requires an actor and reason in the audit trail.

---

## 16. Activity completion and progress

### activity_completions

Records id, tenant_id, learning_activity_id, enrollment_id, status, progress,
completed_at, and completion_source.

Status is not_started, in_progress, or completed. There is at most one
completion record for an activity and enrollment:

UNIQUE (learning_activity_id, enrollment_id)

Completion may contribute to attendance only when the course explicitly
enables that rule; the attendance and completion records remain separate.

---

## 17. Discussion forums

The core domain has no one-to-one direct message or private-chat feature.

### discussion_forums

Records id, tenant_id, learning_activity_id, mode, group_mode,
students_can_create_threads, available_from, and available_until.

### discussion_threads

Records id, tenant_id, forum_id, optional group_id, title, created_by,
created_at, locked_at, and pinned_at.

### discussion_posts

Records id, tenant_id, thread_id, optional parent_post_id, author_user_id,
content, created_at, edited_at, and deleted_at.

Moderation uses soft deletion. Content relevant to audit or moderation is not
immediately destroyed.

---

## 18. Announcements

### announcements

Records id, tenant_id, course_offering_id, title, content, published_by,
published_at, and expires_at. The default audience is active enrollments and
course staff.

### announcement_receipts

Optional read tracking records announcement_id, tenant_id, user_id, and
read_at.

---

## 19. Course copy and import

An old offering may provide instructional structure for a new offering.

May be copied:

- an RPS as a new draft version;
- local CPMK and Sub-CPMK structure;
- modules and lessons;
- materials;
- learning activities;
- assignment and quiz configuration;
- permitted question banks;
- rubrics;
- forum configuration; and
- gradebook structure.

Must not be copied:

- enrollment;
- attendance records or credentials;
- submissions;
- quiz attempts or responses;
- student grades or final grades;
- activity completion;
- student discussion posts; or
- audit logs.

### course_copy_jobs

Records id, tenant_id, source_course_offering_id,
target_course_offering_id, options, status, requested_by, started_at,
completed_at, and error.

Copy creates new entity identifiers. The target must not reference mutable
records belonging to the old offering.

---

## 20. Integration boundary

Integrations are adapters at the edge of the canonical LMS domain:

SIAKAD
   ↓
Integration Adapter
   ↓
Canonical LMS Domain

The core model must not depend on the vendor or protocol details of a
particular SIAKAD.

### integrations

Records id, tenant_id, type, provider, status, configuration_reference, and
last_success_at. Secrets are referenced through secret-management facilities,
never stored as plaintext configuration.

Types may include siakad, spada, pddikti, scorm, xapi, lti, and other.

### external_mappings

Records id, tenant_id, integration_id, entity_type, internal_id, external_id,
external_version, and last_synced_at.

Both directions are unique in one integration context:

UNIQUE (integration_id, entity_type, external_id)
UNIQUE (integration_id, entity_type, internal_id)

### sync_jobs

Records id, tenant_id, integration_id, direction, entity_type, status,
started_at, completed_at, cursor, and summary. direction is inbound or
outbound.

### sync_errors

Records id, tenant_id, sync_job_id, external_id, entity_type, error_code,
error_message, retryable, resolved_at, and resolved_by.

Inbound SIAKAD synchronization may provide terms, program references, users
and external identities, courses, offerings, course staff, enrollments, and
CPL references. It must be idempotent: repeating the same payload does not
create duplicates.

Outbound synchronization may provide final grades and attendance. Only data
that satisfies publication or finalization rules may be sent.

### Future adapters

SPADA and PDDikti are future integrations. SCORM, xAPI, and LTI are integration
boundaries rather than mandatory first-release protocol implementations.
External API changes must not force the core learning model to mirror a vendor
payload.

### learning_packages, external_tools, and learning_events

learning_packages describes packages such as SCORM through tenant-owned
metadata, standard, version, file, creator, and status.

external_tools describes a tenant integration name, protocol,
configuration_reference, and status. Credentials remain outside plaintext
domain configuration.

learning_events provides a tenant-scoped boundary for event type, user,
offering, activity, occurrence time, and payload. High-volume retention may be
handled separately from transactional records.

---

## 21. Audit logs

### audit_logs

audit_logs is an immutable tenant-scoped record containing id, tenant_id,
actor_user_id, actor_role, action, entity_type, entity_id, optional
course_offering_id, before_data, after_data, reason, ip_address, request_id,
and occurred_at.

Audit is required for sensitive actions, including role changes, staff
assignment, publication or unpublication, deadline changes, overrides,
grading, grade publication, final-grade lock or reopen, attendance correction,
archive or unarchive, integration conflict resolution, break-glass access, and
tenant or integration configuration changes.

Normal application users cannot edit audit records.

### platform_audit_logs

The global platform record covers tenant creation or suspension, Super Admin
changes, break-glass actions, global configuration, and cross-tenant support.

---

## 22. Business invariants

### 22.1 Multi-tenancy

1. Data from one tenant is never visible to another.
2. Tenant foreign keys refer to entities in the same tenant.
3. A course cannot contain another tenant's materials.
4. An enrollment requires a membership in the same tenant.
5. Integration mappings cannot cross tenants.
6. Background jobs carry tenant_id.
7. File access validates tenant ownership.
8. Tenant data cache keys include tenant scope.

### 22.2 Identity and authorization

1. Users do not store tenant roles.
2. Access requires an active membership.
3. One user may have many tenant roles through active membership_roles.
4. A lecturer membership alone does not grant course authority.
5. Course staff must be assigned to the offering.
6. A student must have an active enrollment.
7. Authorization checks tenant, role, course scope, object scope, and lifecycle.

### 22.3 SIAKAD

1. Academic master facts are not edited manually in ordinary LMS workflows.
2. Synchronization is idempotent.
3. External identifiers are unique in their integration context.
4. A failed synchronization does not delete existing data.
5. Failures create sync_errors.
6. Operators do not alter grades to hide an integration failure.
7. Conflict resolution is audited.
8. Enrollment changes do not erase historical learning data.

### 22.4 Course offering

1. A course and course offering are different entities.
2. An offering has one academic term.
3. Student content is visible only after the offering or activity is published.
4. Archived offerings are read-only by default.
5. Reopening an archive is audited.
6. A later offering cannot mutate an earlier offering's history.

### 22.5 Enrollment

1. A student has at most one enrollment in an offering.
2. A student accesses only enrolled offerings.
3. Normal operations do not create or alter SIAKAD enrollment locally.
4. Withdrawal preserves submissions, grades, and history.
5. Student access follows enrollment status and tenant policy.

### 22.6 Time and deadlines

Business timestamps are timezone-aware. Storage may use canonical UTC with the
relevant tenant or course timezone retained. Server-local time is never the
meaning of a deadline.

Assignments have available_from, due_at, and cutoff_at:

- due_at is the academic deadline;
- cutoff_at is the hard stop; and
- a late submission is after due_at and before cutoff_at.

The tenant may set cutoff_at equal to due_at for a strict policy or later for
late submissions. Individual extensions use assignment_overrides.

### 22.7 Submissions

1. A submission belongs to a valid enrollment or group.
2. Submission outside the allowed window is rejected unless overridden.
3. Every resubmission creates a new immutable version.
4. Earlier versions are never overwritten.
5. Staff cannot silently replace a student's file.
6. A locked submission cannot be edited by the student.
7. Group submissions require valid group membership.

### 22.8 Quizzes and examinations

1. An attempt cannot start outside its window without an override.
2. Attempt limits are enforced.
3. Timers are server-side.
4. Question versions are frozen for an attempt.
5. Random selection is stored at attempt start.
6. Correct answers are withheld until permitted by feedback policy.
7. Submitted attempts cannot be modified by students.
8. Manual regrading is audited.
9. Invalidated attempts have a reason.

### 22.9 Grades

1. TA grading is draft-only.
2. Instructors may grade but cannot publish final grades unless Lead Instructor.
3. Only Lead Instructor may publish or lock final grades.
4. Scores are not assumed to be 1–100.
5. Scores exceeding grade-item rules require an explicit override.
6. Grade changes record actor and timestamp.
7. Final-grade changes are audited.
8. Locked final grades cannot be edited directly.
9. Outbound synchronization sends only permitted published/finalized grades.
10. Synchronization failure does not change internal academic values.

### 22.10 Attendance

1. One student has one record per attendance session.
2. QR/PIN works only for the correct session and window.
3. Credentials are not reusable across sessions.
4. Corrections record actor and reason.
5. Completion equals attendance only when an explicit course rule says so.
6. Archiving does not remove attendance.
7. Outbound attendance synchronization is idempotent.

### 22.11 Discussions

1. A user sees only accessible course forums.
2. Group forums are visible only to the right group and authorized staff.
3. Post edits and deletion follow moderation policy.
4. Soft-deleted content may be retained for audit and moderation.
5. Forums do not create private direct-message access.

### 22.12 Course copy

1. Copy takes instructional structure only.
2. Student-generated data is excluded.
3. Target entities receive new identifiers.
4. Grades, submissions, attendance, and attempts are excluded.
5. Copy operations are audited.
6. Core copying remains within one tenant.

### 22.13 Deletion and retention

The normal application path does not hard-delete submissions, submission
versions, quiz attempts, quiz responses, grades, final grades, attendance
records, audit logs, relevant synchronization history, or published course-plan
history.

Use archive, soft delete, or a separately governed retention/anonymization
workflow when privacy policy requires it.

---

## 23. RLS classification

The following are global and do not use ordinary tenant RLS:

- tenants;
- users;
- auth_identities;
- auth_sessions;
- platform_admins; and
- platform_audit_logs.

Access remains restricted by platform policy.

The following are tenant-scoped and require RLS:

- tenant_settings, memberships, membership_roles, external_identities;
- academic_terms, academic_program_refs, courses, course_offerings,
  course_staff, enrollments;
- learning_outcomes, outcome_mappings, course_plan_versions,
  course_plan_outcomes;
- modules, lessons, lesson_outcomes, materials, files;
- learning_activities, activity_outcomes;
- assignments, assignment_overrides, submissions, submission_versions,
  submission_files;
- course_groups, course_group_members;
- question_banks, questions, question_versions, question_options, quizzes,
  quiz_question_rules, quiz_overrides, quiz_attempts,
  quiz_attempt_questions, quiz_responses;
- rubrics, rubric_criteria, rubric_levels;
- grade_categories, grade_items, grades, grade_schemes, final_grades;
- attendance_sessions, attendance_credentials, attendance_records;
- activity_completions;
- discussion_forums, discussion_threads, discussion_posts;
- announcements, announcement_receipts;
- course_copy_jobs;
- integrations, external_mappings, sync_jobs, sync_errors;
- learning_packages, external_tools, learning_events; and
- audit_logs.

Every tenant-scoped table stores tenant_id explicitly, including child tables
where it can also be derived through a parent. Composite foreign keys use
(tenant_id, parent_id) and make cross-tenant references structurally
impossible. A child table's tenant_id is therefore redundant by design: it
supports simple, safe RLS and referential checks.

---

## 24. Authorization hierarchy

Authorization follows this order:

Authenticated user
  ↓
Active membership
  ↓
Correct tenant
  ↓
Tenant-role permission
  ↓
Course enrollment or course staff
  ↓
Object permission
  ↓
Lifecycle and state permit the action

For grading:

Authenticated user
  ↓
Active tenant membership
  ↓
Authorized course_staff record
  ↓
Instructor, Lead Instructor, or permitted TA role
  ↓
Grade item belongs to the same offering
  ↓
Grade is not final-locked
  ↓
Action permitted

Holding a role is never sufficient by itself.

---

## 25. Lifecycle references

Course offering:

draft → published → active → closed → archived

Assignment submission:

draft → submitted → returned → resubmission_allowed → submitted → locked

Quiz attempt:

in_progress → submitted or auto_submitted → graded

An attempt may be invalidated with a recorded reason.

Grade:

draft → published

Final grade:

draft → published → locked

Correction uses the audited reopen flow.

---

## 26. ERD and ownership summary

The central relationships are:

users → memberships → tenants

tenants → academic_terms
tenants → courses → course_offerings
course_offerings → course_staff and enrollments
course_offerings → course_plan_versions → modules → lessons
lessons → materials and learning_activities
learning_activities → assignments, quizzes, forums, attendance
assignments → submissions → submission_versions
quizzes → quiz_attempts → quiz_attempt_questions → quiz_responses
course_offerings → grade_categories → grade_items → grades
course_offerings → final_grades
course_offerings → attendance_sessions → attendance_records

Inbound and outbound integration is:

SIAKAD → integration adapter → canonical LMS
canonical LMS → integration adapter → SIAKAD

Global platform records:

- tenants;
- users;
- authentication identities and sessions; and
- platform audit.

SIAKAD-authoritative records:

- academic terms;
- program references;
- course master;
- offering identity;
- instructor assignments;
- student enrollment; and
- CPL references.

LMS-authoritative records:

- publication and offering lifecycle;
- RPS versions and CPMK/Sub-CPMK;
- modules, lessons, materials, and activities;
- assignments and submissions;
- quizzes, examinations, question banks, and rubrics;
- forums and announcements;
- attendance and activity completion;
- gradebook, grades, and final grades; and
- audit trail.

Outbound records:

- published/finalized final grades; and
- attendance.

---

## 27. Database and security requirements

Database constraints should enforce, where applicable:

- unique membership per tenant and user;
- unique active membership role per membership and role;
- unique staff per offering and user;
- unique enrollment per offering and student;
- unique grade per grade item and enrollment;
- unique final grade per offering and enrollment;
- unique attendance record per session and enrollment;
- unique activity completion per activity and enrollment;
- unique external mapping in both directions; and
- composite tenant-consistency foreign keys.

Check constraints should cover valid scores, time ranges, non-negative values,
exactly-one submission owner, and permissible states where the rule can be
represented safely. Service code must enforce rules that require authorization
or multi-step state transitions.

Minimum security requirements are:

- deny-by-default authorization;
- tenant RLS plus object-level authorization;
- trusted transaction-local tenant context;
- signed and short-lived file access;
- secrets in a secret manager or reference;
- hashed attendance credentials;
- immutable audit trails;
- idempotent integration;
- concurrency protection for grading where needed;
- server-side examination timers;
- no correct-answer leakage;
- malware-scanning boundaries for files;
- archive or soft-delete for academic history; and
- request/correlation identifiers for sensitive operations.

Super Admin is not an automatic application-level bypass of tenant RLS.
Break-glass access is an explicit privileged path and is audited.

---

## 28. Reporting and analytics

Student-level reporting may show:

- incomplete activities;
- missing submissions;
- unattempted quizzes;
- personal attendance;
- published grades; and
- personal course progress.

Instructor-level reporting may show:

- active enrollment;
- missing submissions;
- score distribution;
- activity completion;
- attendance;
- student progress;
- outcome-to-assessment mapping; and
- gradebook information.

Tenant-level reporting may show:

- active offerings;
- usage and adoption;
- synchronization health;
- completion and assessment activity;
- attendance exports; and
- final-grade synchronization status.

Reporting authorization must prevent access to an unrelated course or tenant.

---

## 29. Extensibility boundaries

Future capability may include SPADA, PDDikti, SCORM, xAPI, LTI, video
conferencing, plagiarism checking, online proctoring, object storage,
notification, and analytics-warehouse adapters.

These boundaries must not cause core entities to copy one provider's API
schema. The LMS domain remains coherent without any future adapter.

---

## 30. Locked domain decisions

1. campus-lms is an LMS, not SIAKAD.
2. SIAKAD is authoritative for institutional academic facts.
3. Course and course offering are distinct entities.
4. The LMS is designed for campus-grade isolation and history.
5. Grading is configurable and not hard-coded to 1–100.
6. Academic Operator and Teaching Assistant are distinct authority scopes.
7. A TA can provide draft scores only.
8. Assessments support assignments, quizzes/examinations, question banks,
   randomization, attempt rules, overrides, group work, rubrics,
   formative/summative use, outcome mapping, and weighted gradebooks.
9. Attendance supports manual, self check-in, QR/PIN, activity-completion, and
   external methods.
10. Attendance and activity completion are separate domains.
11. Core interaction uses announcements and discussion forums; direct messaging
    is outside the core.
12. Course copying excludes student-generated and academic-history data.
13. SPADA, PDDikti, SCORM, xAPI, and LTI are integration boundaries.
14. users is global; memberships and membership_roles are tenant-scoped.
15. A user may hold several active roles in one tenant.
16. Tenant consistency is enforced with explicit tenant_id and composite
    foreign keys, not only service checks.
17. Every tenant-scoped table stores tenant_id explicitly.
18. Attendance is LMS-owned and may be exported through an adapter.
19. Attendance may be a gradebook component through an explicit typed reference.
20. Only Lead Instructor may publish or lock final grades.
21. Course history, submissions, grades, attendance, and audit records are not
    hard-deleted through normal application operations.
22. Authentication sessions use hashed, revocable refresh-token records.
23. The AI layer, if enabled later, is optional and must not become a
    dependency of the core learning domain.
