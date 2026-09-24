-- 0007_content_schema.up.sql
-- Tenant-scoped course content and file metadata.

CREATE TABLE modules (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    course_offering_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    position INTEGER NOT NULL CHECK (position >= 0),
    status TEXT NOT NULL CHECK (status IN ('draft', 'published', 'hidden')),
    available_from TIMESTAMPTZ NULL,
    available_until TIMESTAMPTZ NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (tenant_id, course_offering_id)
        REFERENCES course_offerings (tenant_id, id),
    CHECK (
        available_until IS NULL
        OR available_from IS NULL
        OR available_from < available_until
    )
);

CREATE TABLE lessons (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    module_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    position INTEGER NOT NULL CHECK (position >= 0),
    learning_mode TEXT NOT NULL CHECK (learning_mode IN ('asynchronous', 'synchronous', 'blended', 'onsite')),
    estimated_minutes INTEGER NOT NULL CHECK (estimated_minutes >= 0),
    available_from TIMESTAMPTZ NULL,
    available_until TIMESTAMPTZ NULL,
    status TEXT NOT NULL CHECK (status IN ('draft', 'published', 'hidden')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (tenant_id, module_id) REFERENCES modules (tenant_id, id),
    CHECK (
        available_until IS NULL
        OR available_from IS NULL
        OR available_from < available_until
    )
);

CREATE TABLE files (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    storage_key TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    checksum TEXT NOT NULL,
    malware_scan_status TEXT NOT NULL CHECK (malware_scan_status IN ('pending', 'clean', 'quarantined', 'failed')),
    uploaded_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE (tenant_id, id),
    UNIQUE (tenant_id, storage_key),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE TABLE materials (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    lesson_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL CHECK (type IN ('text', 'file', 'link', 'video', 'audio', 'embed', 'learning_package', 'external_tool')),
    file_id UUID NULL,
    external_url TEXT NULL,
    content TEXT NULL,
    position INTEGER NOT NULL CHECK (position >= 0),
    published BOOLEAN NOT NULL DEFAULT false,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (tenant_id, lesson_id) REFERENCES lessons (tenant_id, id),
    FOREIGN KEY (tenant_id, file_id) REFERENCES files (tenant_id, id),
    CHECK (
        (type = 'text' AND file_id IS NULL AND external_url IS NULL AND content IS NOT NULL)
        OR (type = 'file' AND file_id IS NOT NULL AND external_url IS NULL AND content IS NULL)
        OR (type IN ('link', 'video', 'audio', 'embed', 'learning_package', 'external_tool')
            AND file_id IS NULL AND external_url IS NOT NULL AND content IS NULL)
    )
);

CREATE INDEX modules_offering_position_idx
    ON modules (tenant_id, course_offering_id, position, id);

CREATE INDEX lessons_module_position_idx
    ON lessons (tenant_id, module_id, position, id);

CREATE INDEX materials_lesson_position_idx
    ON materials (tenant_id, lesson_id, position, id);

CREATE INDEX files_tenant_scan_status_idx
    ON files (tenant_id, malware_scan_status);

ALTER TABLE modules ENABLE ROW LEVEL SECURITY;
CREATE POLICY modules_tenant_isolation ON modules
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE lessons ENABLE ROW LEVEL SECURITY;
CREATE POLICY lessons_tenant_isolation ON lessons
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE files ENABLE ROW LEVEL SECURITY;
CREATE POLICY files_tenant_isolation ON files
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);

ALTER TABLE materials ENABLE ROW LEVEL SECURITY;
CREATE POLICY materials_tenant_isolation ON materials
    USING (tenant_id = current_setting('app.tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);
