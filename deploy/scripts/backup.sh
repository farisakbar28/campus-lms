#!/usr/bin/env bash
# Creates a local custom-format PostgreSQL archive and paired normalized state
# manifest. Off-machine copy, retention, and failure notification are not
# implemented here. The state validator remains pinned to the pre-0006,
# migration-5 baseline and is not validated for the repository's current
# schema; current-schema backup/restore reconciliation is future authorized
# engineering work.

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

require_state_line() {
    local state_file="$1"
    local expected="$2"

    grep -Fxq "$expected" "$state_file" || die "state precondition failed: $expected"
}

validate_state() {
    local state_file="$1"
    local expected_presence
    local expected_rls
    local table

    expected_presence="week3_table_presence=present:academic_terms,present:audit_logs,present:auth_identities,present:course_offerings,present:course_staff,present:courses,present:enrollments,present:membership_roles,present:memberships,present:tenants,present:users"
    expected_rls="rls_state=academic_terms:enabled,audit_logs:enabled,course_offerings:enabled,course_staff:enabled,courses:enabled,enrollments:enabled,membership_roles:enabled,memberships:enabled"

    require_state_line "$state_file" "state_format=week03-backup-restore-v1"
    require_state_line "$state_file" "week3_table_count=11"
    require_state_line "$state_file" "public_application_table_count=11"
    require_state_line "$state_file" "$expected_presence"
    require_state_line "$state_file" "schema_migrations=present"
    require_state_line "$state_file" "schema_migrations_rows=1"
    require_state_line "$state_file" "migration_version=5"
    require_state_line "$state_file" "migration_dirty=false"
    require_state_line "$state_file" "$expected_rls"
    require_state_line "$state_file" "rls_enabled_count=8"
    require_state_line "$state_file" "a7_fk_count=1"
    require_state_line "$state_file" "term_range_check_count=1"
    require_state_line "$state_file" "audit_policy_total=2"
    require_state_line "$state_file" "audit_policy_select=1"
    require_state_line "$state_file" "audit_policy_insert=1"
    require_state_line "$state_file" "audit_policy_update=0"
    require_state_line "$state_file" "audit_policy_delete=0"
    require_state_line "$state_file" "audit_policy_all=0"
    require_state_line "$state_file" "enrollments_active_student_lookup_index_count=1"
    grep -Eq '^enrollments_active_student_lookup_index_definition=CREATE INDEX ' "$state_file" ||
        die "state precondition failed: migration-0005 index definition is missing"

    for table in tenants users auth_identities memberships membership_roles audit_logs academic_terms courses course_offerings course_staff enrollments; do
        grep -Eq "^table_rows:${table}=[0-9]+$" "$state_file" ||
            die "state precondition failed: row count missing for $table"
        grep -Eq "^table_fingerprint:${table}=[[:xdigit:]]{32}$" "$state_file" ||
            die "state precondition failed: fingerprint missing for $table"
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

run_state() {
    local database="$1"

    "${COMPOSE[@]}" exec -T postgres \
        psql -X -q -v ON_ERROR_STOP=1 -A -t -P pager=off \
        -U "$source_user" -d "$database" < "$STATE_SQL"
}

mkdir -p "$BACKUP_DIR"
tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/campus-lms-backup.XXXXXX")"
timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
tmp_dump="$(mktemp "$BACKUP_DIR/.campus_lms_${timestamp}.XXXXXX.dump.partial")"
tmp_state_before="$tmp_dir/source-before.state"
tmp_state_after="$tmp_dir/source-after.state"
artifact="$BACKUP_DIR/campus_lms_${timestamp}_$$_${RANDOM}.dump"
state_manifest="${artifact}.state"

cleanup() {
    rm -f -- "$tmp_dump" "$tmp_state_before" "$tmp_state_after"
    rm -rf -- "$tmp_dir"
}
trap cleanup EXIT

echo "source_database=$source_database"
echo "source_user=$source_user"
echo "source_state_before=begin"
run_state "$source_database" > "$tmp_state_before" || die "source state snapshot before pg_dump failed"
validate_state "$tmp_state_before"
cat "$tmp_state_before"
echo "source_state_before=end"

echo "backup_format=custom"
echo "pg_dump=begin"
if ! "${COMPOSE[@]}" exec -T postgres \
    pg_dump -Fc -U "$source_user" -d "$source_database" > "$tmp_dump"; then
    die "pg_dump failed; partial archive was discarded"
fi
echo "pg_dump=pass"

echo "source_state_after=begin"
run_state "$source_database" > "$tmp_state_after" || die "source state snapshot after pg_dump failed"
validate_state "$tmp_state_after"
cat "$tmp_state_after"
echo "source_state_after=end"

if ! cmp -s "$tmp_state_before" "$tmp_state_after"; then
    echo "source_stability=fail" >&2
    diff -u "$tmp_state_before" "$tmp_state_after" >&2 || true
    die "source state changed during pg_dump"
fi
echo "source_stability=pass"

test -s "$tmp_dump" || die "pg_dump produced an empty archive"
echo "archive_nonempty=pass"

echo "archive_validation=begin"
"${COMPOSE[@]}" exec -T postgres pg_restore --list < "$tmp_dump" > /dev/null
echo "archive_validation=pass"

mv -- "$tmp_dump" "$artifact"
mv -- "$tmp_state_before" "$state_manifest"

echo "artifact_path=$artifact"
echo "state_manifest=$state_manifest"
echo "artifact_archive_validation=begin"
"${COMPOSE[@]}" exec -T postgres pg_restore --list < "$artifact" > /dev/null
echo "artifact_archive_validation=pass"
echo "backup=pass"
