-- Deterministic UUIDs
-- Tenant A: aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
-- Tenant B: bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb
-- User A: cccccccc-cccc-cccc-cccc-cccccccccccc
-- User B: dddddddd-dddd-dddd-dddd-dddddddddddd
-- Academic Term A: eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee
-- Academic Term B: ffffffff-ffff-ffff-ffff-ffffffffffff
-- Course A: 11111111-1111-1111-1111-111111111111
-- Course B: 22222222-2222-2222-2222-222222222222
-- Course Offering A: 33333333-3333-3333-3333-333333333333
-- Course Offering B: 44444444-4444-4444-4444-444444444444
-- Membership A: 55555555-5555-5555-5555-555555555555
-- Membership B: 66666666-6666-6666-6666-666666666666
-- Audit Log A: 77777777-7777-7777-7777-777777777777
-- Audit Log B: 88888888-8888-8888-8888-888888888888

-- Insert tenants
INSERT INTO tenants (id, slug, name, status, default_timezone, created_at, suspended_at)
VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'tenant-a', 'Tenant A', 'active', 'Asia/Jakarta', NOW(), NULL),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'tenant-b', 'Tenant B', 'active', 'Asia/Jakarta', NOW(), NULL);

-- Insert users
INSERT INTO users (id, email, display_name, status, created_at)
VALUES
  ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'usera@example.com', 'User A', 'active', NOW()),
  ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'userb@example.com', 'User B', 'active', NOW());

-- Insert academic terms (valid time range: starts_at < ends_at)
INSERT INTO academic_terms (id, tenant_id, external_id, code, name, starts_at, ends_at, status, synced_at)
VALUES
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ext-term-a', '2026-1', '2026 Ganjil', '2026-08-01 00:00:00+07', '2026-12-31 23:59:59+07', 'active', NOW()),
  ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'ext-term-b', '2026-1', '2026 Ganjil', '2026-08-01 00:00:00+07', '2026-12-31 23:59:59+07', 'active', NOW());

-- Insert courses
INSERT INTO courses (id, tenant_id, external_id, code, name, credits, status, synced_at)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ext-course-a', 'CS101', 'Intro to CS', 3, 'active', NOW()),
  ('22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'ext-course-b', 'CS101', 'Intro to CS', 3, 'active', NOW());

-- Insert course offerings
INSERT INTO course_offerings (id, tenant_id, external_id, course_id, academic_term_id, external_section_code, display_name, lms_status, published_at, closed_at, archived_at, created_at)
VALUES
  ('33333333-3333-3333-3333-333333333333', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ext-offering-a', '11111111-1111-1111-1111-111111111111', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'A', 'CS101-A-2026-1', 'published', NOW(), NULL, NULL, NOW()),
  ('44444444-4444-4444-4444-444444444444', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'ext-offering-b', '22222222-2222-2222-2222-222222222222', 'ffffffff-ffff-ffff-ffff-ffffffffffff', 'A', 'CS101-A-2026-1', 'published', NOW(), NULL, NULL, NOW());

-- Insert memberships
INSERT INTO memberships (id, tenant_id, user_id, status, joined_at)
VALUES
  ('55555555-5555-5555-5555-555555555555', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'active', NOW()),
  ('66666666-6666-6666-6666-666666666666', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'active', NOW());

-- Insert audit logs
INSERT INTO audit_logs (id, tenant_id, actor_user_id, actor_role, action, entity_type, entity_id, course_offering_id, before_data, after_data, reason, ip_address, request_id, occurred_at)
VALUES
  ('77777777-7777-7777-7777-777777777777', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'lecturer', 'CREATE', 'course_offering', '33333333-3333-3333-3333-333333333333', '33333333-3333-3333-3333-333333333333', '{}', '{"name":"CS101-A"}', 'initial', '127.0.0.1', 'req-a', NOW()),
  ('88888888-8888-8888-8888-888888888888', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'lecturer', 'CREATE', 'course_offering', '44444444-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444444', '{}', '{"name":"CS101-A"}', 'initial', '127.0.0.1', 'req-b', NOW());