-- 0003_auth_membership_schema.up.sql
-- Tenant-scoped tables: memberships, membership_roles, audit_logs, course_staff, enrollments

CREATE TABLE memberships (
    id UUID,
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, user_id),
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE membership_roles (
    id UUID,
    tenant_id UUID NOT NULL,
    membership_id UUID NOT NULL,
    role TEXT NOT NULL,
    granted_by UUID,
    granted_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (tenant_id, membership_id) REFERENCES memberships (tenant_id, id),
    FOREIGN KEY (granted_by) REFERENCES users(id)
);

-- Prevent duplicate active roles (A1)
CREATE UNIQUE INDEX membership_roles_active_role_idx ON membership_roles (membership_id, role) WHERE revoked_at IS NULL;

CREATE TABLE audit_logs (
    id UUID,
    tenant_id UUID NOT NULL,
    actor_user_id UUID,
    actor_role TEXT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID,
    course_offering_id UUID,
    before_data JSONB,
    after_data JSONB,
    reason TEXT,
    ip_address TEXT,
    request_id TEXT,
    occurred_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (actor_user_id) REFERENCES users(id),
    FOREIGN KEY (tenant_id, course_offering_id) REFERENCES course_offerings (tenant_id, id)
);

CREATE TABLE course_staff (
    id UUID,
    tenant_id UUID NOT NULL,
    course_offering_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL,
    source TEXT,
    permissions JSONB,
    active BOOLEAN NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, id),
    UNIQUE (course_offering_id, user_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (tenant_id, course_offering_id) REFERENCES course_offerings (tenant_id, id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE enrollments (
    id UUID,
    tenant_id UUID NOT NULL,
    course_offering_id UUID NOT NULL,
    student_user_id UUID NOT NULL,
    external_id TEXT,
    status TEXT NOT NULL,
    enrolled_at TIMESTAMPTZ NOT NULL,
    withdrawn_at TIMESTAMPTZ NULL,
    synced_at TIMESTAMPTZ NULL,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, id),
    UNIQUE (course_offering_id, student_user_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (tenant_id, course_offering_id) REFERENCES course_offerings (tenant_id, id),
    FOREIGN KEY (tenant_id, student_user_id) REFERENCES memberships (tenant_id, user_id)
);

-- Enable RLS and create policies for tenant-scoped tables
ALTER TABLE memberships ENABLE ROW LEVEL SECURITY;
CREATE POLICY memberships_tenant_isolation ON memberships
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE membership_roles ENABLE ROW LEVEL SECURITY;
CREATE POLICY membership_roles_tenant_isolation ON membership_roles
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;
-- audit_logs is immutable: only SELECT and INSERT policies for normal tenant access
CREATE POLICY audit_logs_tenant_select ON audit_logs
    FOR SELECT
    USING (tenant_id = current_setting('app.tenant_id')::uuid);
CREATE POLICY audit_logs_tenant_insert ON audit_logs
    FOR INSERT
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE course_staff ENABLE ROW LEVEL SECURITY;
CREATE POLICY course_staff_tenant_isolation ON course_staff
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE enrollments ENABLE ROW LEVEL SECURITY;
CREATE POLICY enrollments_tenant_isolation ON enrollments
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);