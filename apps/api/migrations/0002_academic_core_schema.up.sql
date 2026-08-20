-- 0002_academic_core_schema.up.sql
-- Tenant-scoped tables: academic_terms, courses, course_offerings

CREATE TABLE academic_terms (
    id UUID,
    tenant_id UUID NOT NULL,
    external_id TEXT,
    code TEXT,
    name TEXT,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    synced_at TIMESTAMPTZ NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE TABLE courses (
    id UUID,
    tenant_id UUID NOT NULL,
    external_id TEXT,
    code TEXT,
    name TEXT,
    credits INTEGER NOT NULL,
    status TEXT NOT NULL,
    synced_at TIMESTAMPTZ NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE TABLE course_offerings (
    id UUID,
    tenant_id UUID NOT NULL,
    external_id TEXT,
    course_id UUID NOT NULL,
    academic_term_id UUID NOT NULL,
    external_section_code TEXT,
    display_name TEXT,
    lms_status TEXT NOT NULL,
    published_at TIMESTAMPTZ NULL,
    closed_at TIMESTAMPTZ NULL,
    archived_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (tenant_id, course_id) REFERENCES courses (tenant_id, id),
    FOREIGN KEY (tenant_id, academic_term_id) REFERENCES academic_terms (tenant_id, id)
);

-- Enable RLS and create policies for tenant-scoped tables
ALTER TABLE academic_terms ENABLE ROW LEVEL SECURITY;
CREATE POLICY academic_terms_tenant_isolation ON academic_terms
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE courses ENABLE ROW LEVEL SECURITY;
CREATE POLICY courses_tenant_isolation ON courses
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE course_offerings ENABLE ROW LEVEL SECURITY;
CREATE POLICY course_offerings_tenant_isolation ON course_offerings
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);