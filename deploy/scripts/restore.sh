#!/usr/bin/env bash
# Restores a local custom-format archive into a fresh disposable database,
# verifies it against its paired state manifest, and removes that database
# afterwards. The state validator remains pinned to the pre-0006, migration-5
# baseline and is not validated for the repository's current schema;
# current-schema backup/restore reconciliation is future authorized engineering
# work.

set -euo pipefail
umask 077

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
COMPOSE=(
    docker compose
    -f "$REPO_ROOT/deploy/compose/docker-compose.yml"
    -f "$REPO_ROOT/deploy/compose/docker-compose.override.yml"
)
STATE_SQL="$REPO_ROOT/apps/api/testdata/database/backup-restore/backup-restore-state.sql"
SEED_VERIFY_SQL="$REPO_ROOT/apps/api/testdata/database/seed/verify-seed-counts.sql"
BACKUP_DIR="$REPO_ROOT/data/backups"

die() {
    echo "ERROR: $*" >&2
    exit 1
}

is_safe_identifier() {
    case "$1" in
        ''|[!a-z_]*|*[!a-z0-9_]*) return 1 ;;
        *) return 0 ;;
    esac
}

is_safe_target() {
    local target="$1"

    is_safe_identifier "$target" || return 1
    case "$target" in
        week03_restore_*) return 0 ;;
        *) return 1 ;;
    esac
}

require_state_line() {
    local state_file="$1"
    local expected="$2"

    grep -Fxq "$expected" "$state_file" || return 1
}

validate_state() {
    local state_file="$1"
    local expected_presence
    local expected_rls
    local table

    expected_presence="week3_table_presence=present:academic_terms,present:audit_logs,present:auth_identities,present:course_offerings,present:course_staff,present:courses,present:enrollments,present:membership_roles,present:memberships,present:tenants,present:users"
    expected_rls="rls_state=academic_terms:enabled,audit_logs:enabled,course_offerings:enabled,course_staff:enabled,courses:enabled,enrollments:enabled,membership_roles:enabled,memberships:enabled"

    require_state_line "$state_file" "state_format=week03-backup-restore-v1" || return 1
    require_state_line "$state_file" "week3_table_count=11" || return 1
    require_state_line "$state_file" "public_application_table_count=11" || return 1
    require_state_line "$state_file" "$expected_presence" || return 1
    require_state_line "$state_file" "schema_migrations=present" || return 1
    require_state_line "$state_file" "schema_migrations_rows=1" || return 1
    require_state_line "$state_file" "migration_version=5" || return 1
    require_state_line "$state_file" "migration_dirty=false" || return 1
    require_state_line "$state_file" "$expected_rls" || return 1
    require_state_line "$state_file" "rls_enabled_count=8" || return 1
    require_state_line "$state_file" "a7_fk_count=1" || return 1
    require_state_line "$state_file" "term_range_check_count=1" || return 1
    require_state_line "$state_file" "audit_policy_total=2" || return 1
    require_state_line "$state_file" "audit_policy_select=1" || return 1
    require_state_line "$state_file" "audit_policy_insert=1" || return 1
    require_state_line "$state_file" "audit_policy_update=0" || return 1
    require_state_line "$state_file" "audit_policy_delete=0" || return 1
    require_state_line "$state_file" "audit_policy_all=0" || return 1
    require_state_line "$state_file" "enrollments_active_student_lookup_index_count=1" || return 1
    grep -Eq '^enrollments_active_student_lookup_index_definition=CREATE INDEX ' "$state_file" || return 1

    for table in tenants users auth_identities memberships membership_roles audit_logs academic_terms courses course_offerings course_staff enrollments; do
        grep -Eq "^table_rows:${table}=[0-9]+$" "$state_file" || return 1
        grep -Eq "^table_fingerprint:${table}=[[:xdigit:]]{32}$" "$state_file" || return 1
    done
}

runtime_config="$("${COMPOSE[@]}" exec -T postgres sh -ec 'printf "%s\\n%s\\n" "${POSTGRES_DB:-}" "${POSTGRES_USER:-}"')" ||
    die "could not discover PostgreSQL runtime configuration"
source_database="$(printf '%s\n' "$runtime_config" | sed -n '1p')"
source_user="$(printf '%s\n' "$runtime_config" | sed -n '2p')"
extra_config="$(printf '%s\n' "$runtime_config" | sed -n '3p')"

is_safe_identifier "$source_database" || die "POSTGRES_DB is empty or not a safe local identifier"
is_safe_identifier "$source_user" || die "POSTGRES_USER is empty or not a safe local identifier"
[ -z "$extra_config" ] || die "unexpected PostgreSQL runtime configuration output"

if (( $# != 0 )); then
    die "restore.sh accepts no positional target or archive; use BACKUP_FILE=/absolute/path/file.dump"
fi
if [[ -v RESTORE_DB || -v TARGET_DB ]]; then
    die "RESTORE_DB and TARGET_DB overrides are not supported"
fi

run_state() {
    local database="$1"

    "${COMPOSE[@]}" exec -T postgres \
        psql -X -q -v ON_ERROR_STOP=1 -A -t -P pager=off \
        -U "$source_user" -d "$database" < "$STATE_SQL"
}

run_seed_verifier() {
    local database="$1"

    "${COMPOSE[@]}" exec -T postgres \
        psql -X -v ON_ERROR_STOP=1 -P pager=off \
        -U "$source_user" -d "$database" < "$SEED_VERIFY_SQL"
}

database_exists() {
    local database="$1"

    "${COMPOSE[@]}" exec -T postgres \
        psql -X -q -v ON_ERROR_STOP=1 -A -t -P pager=off \
        -U "$source_user" -d postgres \
        -c "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '$database');"
}

validate_archive_candidate() {
    local candidate="$1"
    local candidate_state="${candidate}.state"

    [ -f "$candidate" ] || return 1
    [ -r "$candidate" ] || return 1
    [ -s "$candidate" ] || return 1
    [ -f "$candidate_state" ] || return 1
    [ -r "$candidate_state" ] || return 1
    [ -s "$candidate_state" ] || return 1
    "${COMPOSE[@]}" exec -T postgres pg_restore --list < "$candidate" > /dev/null 2>/dev/null || return 1
    validate_state "$candidate_state" || return 1
}

select_archive() {
    local candidate
    local -a candidates=()
    local -a sorted_candidates=()

    if [[ -v BACKUP_FILE ]]; then
        [ -n "$BACKUP_FILE" ] || die "BACKUP_FILE must not be empty"
        case "$BACKUP_FILE" in
            /*) ;;
            *) die "BACKUP_FILE must be an absolute path" ;;
        esac
        validate_archive_candidate "$BACKUP_FILE" ||
            die "explicit backup is missing, unreadable, malformed, or not a valid custom archive"
        printf '%s\n' "$BACKUP_FILE"
        return 0
    fi

    shopt -s nullglob
    candidates=("$BACKUP_DIR"/*.dump)
    shopt -u nullglob
    [ "${#candidates[@]}" -gt 0 ] || die "no local .dump backup exists in $BACKUP_DIR"

    mapfile -t sorted_candidates < <(printf '%s\n' "${candidates[@]}" | sort -r)
    for candidate in "${sorted_candidates[@]}"; do
        if validate_archive_candidate "$candidate"; then
            printf '%s\n' "$candidate"
            return 0
        fi
        echo "candidate_rejected=$candidate" >&2
    done

    die "no valid local dump with a matching state manifest was found"
}

monotonic_ns() {
    python3 -c 'import time; print(time.monotonic_ns())'
}

BACKUP_FILE_SELECTED="$(select_archive)"
STATE_MANIFEST="${BACKUP_FILE_SELECTED}.state"

tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/campus-lms-restore.XXXXXX")"
source_before_restore="$tmp_dir/source-before-restore.state"
source_after_restore="$tmp_dir/source-after-restore.state"
restored_state="$tmp_dir/restored.state"
source_seed_before="$tmp_dir/source-seed-before.txt"
source_seed_after="$tmp_dir/source-seed-after.txt"
restored_seed="$tmp_dir/restored-seed.txt"
target_database=""
target_created=0

cleanup() {
    local original_status="$?"
    local cleanup_status=0

    trap - EXIT

    if (( target_created == 1 )); then
        if ! is_safe_target "$target_database" ||
            [ "$target_database" = "$source_database" ] ||
            [ "$target_database" = "campus_lms" ]; then
            echo "cleanup=failure" >&2
            echo "retained_target=$target_database" >&2
            cleanup_status=1
        elif "${COMPOSE[@]}" exec -T postgres \
            psql -X -q -v ON_ERROR_STOP=1 -U "$source_user" -d postgres \
            -c "DROP DATABASE \"$target_database\";" > /dev/null; then
            echo "cleanup=success"
        else
            echo "cleanup=failure" >&2
            echo "retained_target=$target_database" >&2
            cleanup_status=1
        fi
    else
        echo "cleanup=not_needed"
    fi

    rm -rf -- "$tmp_dir"

    if (( cleanup_status != 0 )); then
        exit 1
    fi
    exit "$original_status"
}
trap cleanup EXIT

echo "backup_file=$BACKUP_FILE_SELECTED"
echo "state_manifest=$STATE_MANIFEST"
echo "source_database=$source_database"
echo "source_user=$source_user"

run_state "$source_database" > "$source_before_restore" ||
    die "source state could not be read before restore"
validate_state "$source_before_restore" ||
    die "source state before restore is structurally inconsistent"
if ! cmp -s "$STATE_MANIFEST" "$source_before_restore"; then
    echo "source_manifest_comparison=fail" >&2
    diff -u "$STATE_MANIFEST" "$source_before_restore" >&2 || true
    die "source changed after the backup was created"
fi
echo "source_manifest_comparison=pass"

if ! run_seed_verifier "$source_database" > "$source_seed_before" 2>&1; then
    cat "$source_seed_before" >&2
    die "Week 3 source seed/invariant verification failed"
fi
cat "$source_seed_before"
echo "source_seed_verifier=pass"

source_exists="$(database_exists "$source_database")"
[ "$source_exists" = "t" ] || die "source database does not exist"
echo "source_exists=pass"

restore_timestamp="$(date -u +%Y%m%d%H%M%S)"
target_database="week03_restore_${restore_timestamp}_$$_${RANDOM}"
is_safe_target "$target_database" || die "generated restore target is unsafe"
[ "$target_database" != "$source_database" ] || die "restore target equals source database"
[ "$target_database" != "campus_lms" ] || die "restore target equals main database"
[ -n "$target_database" ] || die "restore target is empty"

target_exists="$(database_exists "$target_database")"
[ "$target_exists" = "f" ] || die "generated restore target already exists"
echo "restore_target=$target_database"
echo "restore_target_differs_from_source=pass"
echo "restore_target_absent_before_create=pass"

"${COMPOSE[@]}" exec -T postgres \
    psql -X -q -v ON_ERROR_STOP=1 -U "$source_user" -d postgres \
    -c "CREATE DATABASE \"$target_database\";" > /dev/null
target_created=1

target_tables="$("${COMPOSE[@]}" exec -T postgres \
    psql -X -q -v ON_ERROR_STOP=1 -A -t -P pager=off \
    -U "$source_user" -d "$target_database" \
    -c "SELECT count(*) FROM pg_tables WHERE schemaname = 'public';")"
target_schema_migrations="$("${COMPOSE[@]}" exec -T postgres \
    psql -X -q -v ON_ERROR_STOP=1 -A -t -P pager=off \
    -U "$source_user" -d "$target_database" \
    -c "SELECT CASE WHEN to_regclass('public.schema_migrations') IS NULL THEN 'absent' ELSE 'present' END;")"
[ "$target_tables" = "0" ] || die "restore target is not empty"
[ "$target_schema_migrations" = "absent" ] || die "restore target already has schema_migrations"
echo "target_initial_public_tables=$target_tables"
echo "target_initial_schema_migrations=$target_schema_migrations"
echo "target_initial_empty=pass"

restore_started="$(monotonic_ns)"
restore_status=0
if "${COMPOSE[@]}" exec -T postgres \
    pg_restore --exit-on-error --single-transaction --no-owner --no-acl \
    -U "$source_user" \
    --dbname "$target_database" < "$BACKUP_FILE_SELECTED"; then
    restore_status=0
else
    restore_status=$?
fi
restore_finished="$(monotonic_ns)"
echo "pg_restore_exit=$restore_status"
[ "$restore_status" -eq 0 ] || die "pg_restore failed"
restore_duration_ns=$((restore_finished - restore_started))
restore_duration_ms=$(( (restore_duration_ns + 999999) / 1000000 ))
echo "local_restore_duration_ns=$restore_duration_ns"
echo "local_restore_duration_ms=$restore_duration_ms"
echo "pg_restore=pass"

run_state "$target_database" > "$restored_state" || die "restored state snapshot failed"
validate_state "$restored_state" || die "restored state is structurally inconsistent"
echo "restored_state=begin"
cat "$restored_state"
echo "restored_state=end"

if ! cmp -s "$STATE_MANIFEST" "$restored_state"; then
    echo "source_restored_comparison=fail" >&2
    diff -u "$STATE_MANIFEST" "$restored_state" >&2 || true
    die "source and restored normalized states differ"
fi
echo "source_restored_comparison=pass"

if ! run_seed_verifier "$target_database" > "$restored_seed" 2>&1; then
    cat "$restored_seed" >&2
    die "restored Week 3 seed/invariant verification failed"
fi
cat "$restored_seed"
echo "restored_seed_verifier=pass"

run_state "$source_database" > "$source_after_restore" ||
    die "source state could not be read after restore"
validate_state "$source_after_restore" ||
    die "source state after restore is structurally inconsistent"
if ! cmp -s "$STATE_MANIFEST" "$source_after_restore"; then
    echo "main_db_safety=fail" >&2
    diff -u "$STATE_MANIFEST" "$source_after_restore" >&2 || true
    die "source/main database changed during restore verification"
fi
if ! run_seed_verifier "$source_database" > "$source_seed_after" 2>&1; then
    cat "$source_seed_after" >&2
    die "source/main Week 3 seed/invariant verification failed after restore"
fi
cat "$source_seed_after"
source_exists_after="$(database_exists "$source_database")"
[ "$source_exists_after" = "t" ] || die "source database disappeared during restore"
echo "source_exists_after_restore=pass"
echo "main_db_safety=pass"
echo "restore_drill=pass"
