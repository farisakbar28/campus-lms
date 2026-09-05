-- 0006_auth_sessions_schema.up.sql
-- Global authentication session state; this table is intentionally not tenant-scoped.

CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_token_hash BYTEA NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    rotated_from UUID NULL,
    revoked_at TIMESTAMPTZ NULL,
    revoked_reason TEXT NULL,
    user_agent TEXT NULL,
    ip_address INET NULL,
    last_seen_at TIMESTAMPTZ NULL,
    CONSTRAINT auth_sessions_refresh_token_hash_key UNIQUE (refresh_token_hash),
    CONSTRAINT auth_sessions_id_user_id_key UNIQUE (id, user_id),
    CONSTRAINT auth_sessions_expires_after_issued_check
        CHECK (expires_at > issued_at),
    CONSTRAINT auth_sessions_revoked_after_issued_check
        CHECK (revoked_at IS NULL OR revoked_at >= issued_at),
    CONSTRAINT auth_sessions_last_seen_after_issued_check
        CHECK (last_seen_at IS NULL OR last_seen_at >= issued_at),
    CONSTRAINT auth_sessions_not_self_rotated_check
        CHECK (rotated_from IS NULL OR rotated_from <> id),
    CONSTRAINT auth_sessions_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT auth_sessions_rotated_from_user_id_fkey
        FOREIGN KEY (rotated_from, user_id) REFERENCES auth_sessions (id, user_id)
);

-- A predecessor can have at most one rotation child. NULL predecessors are
-- allowed for the initial session in each independent session lineage.
CREATE UNIQUE INDEX auth_sessions_rotated_from_unique_idx
    ON auth_sessions (rotated_from)
    WHERE rotated_from IS NOT NULL;

-- Supports lifecycle lookups for a user's sessions that are not revoked.
CREATE INDEX auth_sessions_active_user_idx
    ON auth_sessions (user_id)
    WHERE revoked_at IS NULL;

-- Supports a future expiry-cleanup job without adding an index for token data.
CREATE INDEX auth_sessions_expires_at_idx
    ON auth_sessions (expires_at);
