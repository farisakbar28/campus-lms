-- 0001_tenant_identity_schema.up.sql
-- Global tables: tenants, users, auth_identities

CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    default_timezone TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    suspended_at TIMESTAMPTZ NULL
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE auth_identities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    provider TEXT NOT NULL,
    provider_subject TEXT NOT NULL,
    verified_at TIMESTAMPTZ NULL
);