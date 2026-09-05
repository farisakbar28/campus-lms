#!/usr/bin/env bash

# Verify migration 0006 against disposable local PostgreSQL databases.
# The Compose services supply credentials through their existing environment;
# this harness never reads or prints .env contents.

set -euo pipefail

MODE="${1:-}"
case "$MODE" in
    lifecycle|catalog|constraints|scope) ;;
    *)
        echo "usage: $0 {lifecycle|catalog|constraints|scope}" >&2
        exit 2
        ;;
esac

COMPOSE=(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml)
RUN_ID="$(date +%s%N)_$$_${RANDOM}"
DB_NAME="auth_session_${MODE}_${RUN_ID}"
FRESH_DB="auth_session_fresh_${RUN_ID}"
UPGRADE_DB="auth_session_upgrade_${RUN_ID}"
TMP_DIR="$(mktemp -d)"
CREATED_DATABASES=()

PRE_AUTH_SESSION_TABLES=(
    tenants
    users
    auth_identities
    academic_terms
    courses
    course_offerings
    memberships
    membership_roles
    audit_logs
    course_staff
    enrollments
)

TABLE_ARGS=()
for table_name in "${PRE_AUTH_SESSION_TABLES[@]}"; do
    TABLE_ARGS+=(--table="public.${table_name}")
done

cleanup() {
    local original_status="$1"
    local cleanup_status=0
    trap - EXIT

    for database_name in "${CREATED_DATABASES[@]}"; do
        if ! run_psql postgres -c "DROP DATABASE IF EXISTS \"${database_name}\";" >/dev/null; then
            cleanup_status=1
        fi
    done
    rm -rf -- "$TMP_DIR"

    if [ "$cleanup_status" -eq 0 ]; then
        echo "cleanup=success"
    else
        echo "cleanup=failure" >&2
    fi

    if [ "$original_status" -ne 0 ]; then
        exit "$original_status"
    fi
    exit "$cleanup_status"
}

on_exit() {
    local original_status="$?"
    cleanup "$original_status"
}

trap on_exit EXIT

run_psql() {
    local database_name="$1"
    shift
    "${COMPOSE[@]}" exec -T postgres sh -ec '
        set -eu
        database_name="$1"
        shift
        exec psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database_name" "$@"
    ' sh "$database_name" "$@"
}

run_migrate() {
    local database_name="$1"
    shift
    "${COMPOSE[@]}" exec -T api sh -ec '
        set -eu
        target_db="$1"
        shift
        test -n "$MIGRATE_DATABASE_URL" || { echo "MIGRATE_DATABASE_URL is not set" >&2; exit 1; }
        case "$target_db" in
            ""|[!a-z_]*|*[!a-z0-9_]*) echo "target database name is invalid" >&2; exit 1;;
        esac
        case "$MIGRATE_DATABASE_URL" in
            postgres://*|postgresql://*) ;;
            *) echo "MIGRATE_DATABASE_URL must use a PostgreSQL URL" >&2; exit 1;;
        esac
        scheme=${MIGRATE_DATABASE_URL%%://*}
        remainder=${MIGRATE_DATABASE_URL#*://}
        case "$remainder" in
            */*) ;;
            *) echo "MIGRATE_DATABASE_URL has no database component" >&2; exit 1;;
        esac
        authority=${remainder%%/*}
        path_and_query=${remainder#*/}
        source_database=${path_and_query%%\?*}
        suffix=${path_and_query#"$source_database"}
        case "$authority:$source_database" in
            :*|*:|*/*) echo "MIGRATE_DATABASE_URL has an invalid database component" >&2; exit 1;;
        esac
        migrate_url="${scheme}://${authority}/${target_db}${suffix}"
        migrate_path="/go/bin/migrate"
        test -x "$migrate_path" || { echo "migrate binary is missing" >&2; exit 1; }
        exec "$migrate_path" -path /src/migrations -database "$migrate_url" "$@"
    ' sh "$database_name" "$@"
}

dump_pre_auth_session_tables() {
    local database_name="$1"
    shift
    "${COMPOSE[@]}" exec -T postgres sh -ec '
        set -eu
        database_name="$1"
        shift
        exec pg_dump -U "$POSTGRES_USER" -d "$database_name" --no-owner --no-privileges --format=plain "$@"
    ' sh "$database_name" "${TABLE_ARGS[@]}" "$@"
}

create_database() {
    local database_name="$1"
    case "$database_name" in
        auth_session_*) ;;
        *) echo "unexpected database name" >&2; return 1;;
    esac
    run_psql postgres -c "CREATE DATABASE \"${database_name}\";" >/dev/null
    CREATED_DATABASES+=("$database_name")
}

assert_equal() {
    local label="$1"
    local expected="$2"
    local actual="$3"
    if [ "$actual" != "$expected" ]; then
        echo "${label}=FAIL expected=${expected@Q} actual=${actual@Q}" >&2
        return 1
    fi
    echo "${label}=pass value=${actual}"
}

expect_failure() {
    local label="$1"
    local pattern="$2"
    shift 2
    local output
    local status
    set +e
    output=$("$@" 2>&1)
    status=$?
    set -e
    echo "${label}_exit=${status}"
    echo "$output"
    if [ "$status" -eq 0 ] || ! grep -Eiq "$pattern" <<<"$output"; then
        echo "${label}=FAIL" >&2
        return 1
    fi
    echo "${label}=pass"
}

case "$MODE" in
    lifecycle)
        echo "fresh_db=${FRESH_DB}"
        echo "upgrade_db=${UPGRADE_DB}"
        create_database "$FRESH_DB"
        run_migrate "$FRESH_DB" up >/dev/null
        fresh_version=$(run_migrate "$FRESH_DB" version 2>&1)
        assert_equal fresh_0_to_6 6 "$fresh_version"
        assert_equal fresh_auth_sessions_table auth_sessions "$(run_psql "$FRESH_DB" -tA -c "SELECT to_regclass('public.auth_sessions');")"

        create_database "$UPGRADE_DB"
        run_migrate "$UPGRADE_DB" up 5 >/dev/null
        upgrade_before=$(run_migrate "$UPGRADE_DB" version 2>&1)
        assert_equal upgrade_before_0006 5 "$upgrade_before"
        run_psql "$UPGRADE_DB" -c "INSERT INTO users (id, email, display_name, status, created_at) VALUES ('11111111-1111-1111-1111-111111111111', 'auth-session-lifecycle@example.test', 'Auth Session Lifecycle', 'active', '2026-01-01T00:00:00Z');" >/dev/null
        run_psql "$UPGRADE_DB" -c "INSERT INTO tenants (id, slug, name, status, default_timezone, created_at) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'auth-session-lifecycle', 'Auth Session Lifecycle Tenant', 'active', 'UTC', '2026-01-01T00:00:00Z');" >/dev/null
        dump_pre_auth_session_tables "$UPGRADE_DB" > "$TMP_DIR/pre-auth-session-before.sql"
        sed -e '/^\\restrict /d' -e '/^\\unrestrict /d' "$TMP_DIR/pre-auth-session-before.sql" > "$TMP_DIR/pre-auth-session-before.normalized.sql"

        run_migrate "$UPGRADE_DB" up 1 >/dev/null
        upgrade_after=$(run_migrate "$UPGRADE_DB" version 2>&1)
        assert_equal upgrade_5_to_6 6 "$upgrade_after"
        run_migrate "$UPGRADE_DB" down 1 >/dev/null
        rollback_version=$(run_migrate "$UPGRADE_DB" version 2>&1)
        assert_equal rollback_6_to_5 5 "$rollback_version"
        assert_equal rollback_auth_sessions_absent t "$(run_psql "$UPGRADE_DB" -tA -c "SELECT to_regclass('public.auth_sessions') IS NULL;")"
        run_migrate "$UPGRADE_DB" up 1 >/dev/null
        reapply_version=$(run_migrate "$UPGRADE_DB" version 2>&1)
        assert_equal reapply_5_to_6 6 "$reapply_version"
        final_state=$(run_psql "$UPGRADE_DB" -tA -c "SELECT version || '|' || dirty FROM schema_migrations;")
        assert_equal final_version_dirty "6|false" "$final_state"

        dump_pre_auth_session_tables "$UPGRADE_DB" > "$TMP_DIR/pre-auth-session-after.sql"
        sed -e '/^\\restrict /d' -e '/^\\unrestrict /d' "$TMP_DIR/pre-auth-session-after.sql" > "$TMP_DIR/pre-auth-session-after.normalized.sql"
        if diff -u "$TMP_DIR/pre-auth-session-before.normalized.sql" "$TMP_DIR/pre-auth-session-after.normalized.sql"; then
            echo "pre_auth_session_tables_unchanged=pass"
        else
            echo "pre_auth_session_tables_unchanged=FAIL" >&2
            exit 1
        fi
        ;;
    catalog)
        echo "database=${DB_NAME}"
        create_database "$DB_NAME"
        run_migrate "$DB_NAME" up >/dev/null

        column_catalog=$(run_psql "$DB_NAME" -At -c "SELECT ordinal_position || '|' || column_name || '|' || data_type || '|' || is_nullable FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'auth_sessions' ORDER BY ordinal_position;")
        printf '%s\n' "$column_catalog"
        expected_columns=$(printf '%s\n' \
            '1|id|uuid|NO' \
            '2|user_id|uuid|NO' \
            '3|refresh_token_hash|bytea|NO' \
            '4|issued_at|timestamp with time zone|NO' \
            '5|expires_at|timestamp with time zone|NO' \
            '6|rotated_from|uuid|YES' \
            '7|revoked_at|timestamp with time zone|YES' \
            '8|revoked_reason|text|YES' \
            '9|user_agent|text|YES' \
            '10|ip_address|inet|YES' \
            '11|last_seen_at|timestamp with time zone|YES')
        assert_equal auth_sessions_columns "$expected_columns" "$column_catalog"

        constraint_catalog=$(run_psql "$DB_NAME" -At -c "SELECT conname || '|' || contype::text || '|' || convalidated::text || '|' || pg_get_constraintdef(oid) FROM pg_constraint WHERE conrelid = 'public.auth_sessions'::regclass ORDER BY conname;")
        printf '%s\n' "$constraint_catalog"
        for required_constraint in \
            'auth_sessions_id_user_id_key|u|true|UNIQUE (id, user_id)' \
            'auth_sessions_refresh_token_hash_key|u|true|UNIQUE (refresh_token_hash)' \
            'auth_sessions_user_id_fkey|f|true|FOREIGN KEY (user_id) REFERENCES users(id)' \
            'auth_sessions_rotated_from_user_id_fkey|f|true|FOREIGN KEY (rotated_from, user_id) REFERENCES auth_sessions(id, user_id)'; do
            if ! grep -Fq "$required_constraint" <<<"$constraint_catalog"; then
                echo "missing_constraint=${required_constraint}" >&2
                exit 1
            fi
        done
        if ! grep -Fq 'auth_sessions_not_self_rotated_check|c|true|CHECK' <<<"$constraint_catalog" || ! grep -Fq 'rotated_from <> id' <<<"$constraint_catalog"; then
            echo "missing_self_rotation_check" >&2
            exit 1
        fi
        if grep -Fq 'FOREIGN KEY (rotated_from) REFERENCES auth_sessions(id)' <<<"$constraint_catalog"; then
            echo "simple_rotated_from_fk=FAIL" >&2
            exit 1
        fi
        echo "simple_rotated_from_fk=pass absent"

        index_catalog=$(run_psql "$DB_NAME" -At -c "SELECT indexname || '|' || indexdef FROM pg_indexes WHERE schemaname = 'public' AND tablename = 'auth_sessions' ORDER BY indexname;")
        printf '%s\n' "$index_catalog"
        rotated_index=$(run_psql "$DB_NAME" -tA -c "SELECT indisunique || '|' || pg_get_expr(indpred, indrelid) FROM pg_index WHERE indexrelid = 'public.auth_sessions_rotated_from_unique_idx'::regclass;")
        assert_equal rotated_from_index "true|(rotated_from IS NOT NULL)" "$rotated_index"
        ;;
    constraints)
        echo "database=${DB_NAME}"
        create_database "$DB_NAME"
        run_migrate "$DB_NAME" up >/dev/null
        run_psql "$DB_NAME" -c "INSERT INTO users (id, email, display_name, status, created_at) VALUES ('11111111-1111-1111-1111-111111111111', 'auth-session-user-a@example.test', 'Auth Session User A', 'active', '2026-01-01T00:00:00Z'), ('22222222-2222-2222-2222-222222222222', 'auth-session-user-b@example.test', 'Auth Session User B', 'active', '2026-01-01T00:00:00Z');" >/dev/null
        run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at) VALUES ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', decode('00', 'hex'), '2026-01-01T00:00:00Z', '2026-02-01T00:00:00Z');" >/dev/null

        expect_failure cross_user_rotation 'foreign key' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at, rotated_from) VALUES ('aaaaaaaa-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', decode('0a', 'hex'), '2026-01-02T00:00:00Z', '2026-03-01T00:00:00Z', '33333333-3333-3333-3333-333333333333');"

        run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at) VALUES ('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', decode('01', 'hex'), '2026-01-02T00:00:00Z', '2026-03-01T00:00:00Z');" >/dev/null
        expect_failure self_rotation 'auth_sessions_not_self_rotated_check|check constraint' \
            run_psql "$DB_NAME" -c "UPDATE auth_sessions SET rotated_from = id WHERE id = '44444444-4444-4444-4444-444444444444';"

        run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at, rotated_from) VALUES ('55555555-5555-5555-5555-555555555555', '11111111-1111-1111-1111-111111111111', decode('02', 'hex'), '2026-01-02T00:00:00Z', '2026-03-01T00:00:00Z', '33333333-3333-3333-3333-333333333333');" >/dev/null
        run_psql "$DB_NAME" -c "UPDATE auth_sessions SET revoked_at = '2026-01-03T00:00:00Z', revoked_reason = 'rotation' WHERE id = '55555555-5555-5555-5555-555555555555';" >/dev/null
        retained_count=$(run_psql "$DB_NAME" -tA -c "SELECT count(*) FROM auth_sessions WHERE id IN ('33333333-3333-3333-3333-333333333333', '55555555-5555-5555-5555-555555555555');")
        assert_equal historical_sessions_retained 2 "$retained_count"

        expect_failure duplicate_refresh_token_hash 'duplicate key|unique constraint' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at) VALUES ('bbbbbbbb-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', decode('00', 'hex'), '2026-01-04T00:00:00Z', '2026-02-01T00:00:00Z');"
        expect_failure duplicate_rotated_from 'duplicate key|unique constraint' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at, rotated_from) VALUES ('cccccccc-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', decode('03', 'hex'), '2026-01-04T00:00:00Z', '2026-02-01T00:00:00Z', '33333333-3333-3333-3333-333333333333');"
        expect_failure expires_at_check 'auth_sessions_expires_after_issued_check|check constraint' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at) VALUES ('dddddddd-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', decode('04', 'hex'), '2026-02-01T00:00:00Z', '2026-01-01T00:00:00Z');"
        expect_failure revoked_at_check 'auth_sessions_revoked_after_issued_check|check constraint' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at, revoked_at) VALUES ('eeeeeeee-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', decode('05', 'hex'), '2026-02-01T00:00:00Z', '2026-03-01T00:00:00Z', '2026-01-31T00:00:00Z');"
        expect_failure last_seen_at_check 'auth_sessions_last_seen_after_issued_check|check constraint' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at, last_seen_at) VALUES ('ffffffff-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', decode('06', 'hex'), '2026-02-01T00:00:00Z', '2026-03-01T00:00:00Z', '2026-01-31T00:00:00Z');"
        expect_failure nonexistent_user 'foreign key' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at) VALUES ('00000000-1111-1111-1111-111111111111', '66666666-6666-6666-6666-666666666666', decode('07', 'hex'), '2026-02-01T00:00:00Z', '2026-03-01T00:00:00Z');"
        expect_failure nonexistent_predecessor 'foreign key' \
            run_psql "$DB_NAME" -c "INSERT INTO auth_sessions (id, user_id, refresh_token_hash, issued_at, expires_at, rotated_from) VALUES ('99999999-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', decode('08', 'hex'), '2026-02-01T00:00:00Z', '2026-03-01T00:00:00Z', '66666666-6666-6666-6666-666666666666');"
        ;;
    scope)
        echo "database=${DB_NAME}"
        create_database "$DB_NAME"
        run_migrate "$DB_NAME" up >/dev/null
        tenant_column_count=$(run_psql "$DB_NAME" -tA -c "SELECT count(*) FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'auth_sessions' AND column_name = 'tenant_id';")
        assert_equal auth_sessions_tenant_id_columns 0 "$tenant_column_count"
        rls_state=$(run_psql "$DB_NAME" -tA -c "SELECT relrowsecurity || '|' || relforcerowsecurity FROM pg_class WHERE oid = 'public.auth_sessions'::regclass;")
        assert_equal auth_sessions_rls "false|false" "$rls_state"
        pre_auth_session_tables=$(run_psql "$DB_NAME" -tA -c "SELECT string_agg(relname, ',' ORDER BY relname) FROM pg_class WHERE relnamespace = 'public'::regnamespace AND relkind = 'r' AND relname IN ('tenants', 'users', 'auth_identities', 'academic_terms', 'courses', 'course_offerings', 'memberships', 'membership_roles', 'audit_logs', 'course_staff', 'enrollments');")
        assert_equal pre_auth_session_tables_present academic_terms,audit_logs,auth_identities,course_offerings,course_staff,courses,enrollments,membership_roles,memberships,tenants,users "$pre_auth_session_tables"
        premature_tables=$(run_psql "$DB_NAME" -tA -c "SELECT COALESCE(string_agg(relname, ',' ORDER BY relname), '') FROM pg_class WHERE relnamespace = 'public'::regnamespace AND relkind = 'r' AND relname NOT IN ('schema_migrations', 'auth_sessions', 'tenants', 'users', 'auth_identities', 'academic_terms', 'courses', 'course_offerings', 'memberships', 'membership_roles', 'audit_logs', 'course_staff', 'enrollments');")
        assert_equal unexpected_non_auth_session_tables "" "$premature_tables"
        ;;
esac

echo "verification=${MODE} pass"
