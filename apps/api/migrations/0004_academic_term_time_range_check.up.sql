-- 0004_academic_term_time_range_check.up.sql
-- Add CHECK constraint: starts_at < ends_at for academic_terms

ALTER TABLE academic_terms
    ADD CONSTRAINT academic_terms_valid_time_range
    CHECK (starts_at < ends_at);