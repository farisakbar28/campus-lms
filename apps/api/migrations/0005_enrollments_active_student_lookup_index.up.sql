-- 0005_enrollments_active_student_lookup_index.up.sql
-- Supports the tenant-scoped active student enrollment dashboard lookup.

CREATE INDEX enrollments_active_student_lookup_idx
    ON enrollments (tenant_id, student_user_id)
    WHERE status = 'active';
