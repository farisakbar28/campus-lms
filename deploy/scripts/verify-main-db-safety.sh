#!/usr/bin/env bash
# Read-only source-state verification against the selected local backup
# manifest. The state validator remains pinned to the pre-0006, migration-5
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

if (( $# != 0 )); then
    die "verify-main-db-safety.sh accepts no positional arguments"
fi
if [[ -v RESTORE_DB || -v TARGET_DB ]]; then
    die "RESTORE_DB and TARGET_DB are not supported"
fi

runtime_config="$("${COMPOSE[@]}" exec -T postgres sh -ec 'printf "%s\\n%s\\n" "${POSTGRES_DB:-}" "${POSTGRES_USER:-}"')" ||
    die "could not discover PostgreSQL runtime configuration"
source_database="$(printf '%s\n' "$runtime_config" | sed -n '1p')"
source_user="$(printf '%s\n' "$runtime_config" | sed -n '2p')"
extra_config="$(printf '%s\n' "$runtime_config" | sed -n '3p')"
is_safe_identifier "$source_database" || die "POSTGRES_DB is empty or unsafe"
is_safe_identifier "$source_user" || die "POSTGRES_USER is empty or unsafe"
[ -z "$extra_config" ] || die "unexpected PostgreSQL runtime configuration output"

if [[ -v BACKUP_FILE ]]; then
    [ -n "$BACKUP_FILE" ] || die "BACKUP_FILE must not be empty"
    case "$BACKUP_FILE" in
        /*) ;;
        *) die "BACKUP_FILE must be an absolute path" ;;
    esac
    selected_backup="$BACKUP_FILE"
else
    shopt -s nullglob
    candidates=("$BACKUP_DIR"/*.dump)
    shopt -u nullglob
    [ "${#candidates[@]}" -gt 0 ] || die "no local backup exists in $BACKUP_DIR"
    mapfile -t sorted_candidates < <(printf '%s\n' "${candidates[@]}" | sort -r)
    selected_backup=""
    for candidate in "${sorted_candidates[@]}"; do
        if [ -f "$candidate" ] && [ -r "$candidate" ] && [ -f "${candidate}.state" ] && [ -r "${candidate}.state" ]; then
            selected_backup="$candidate"
            break
        fi
    done
    [ -n "$selected_backup" ] || die "no backup with a readable state manifest exists"
fi

[ -f "$selected_backup" ] || die "selected backup does not exist"
[ -r "$selected_backup" ] || die "selected backup is not readable"
state_manifest="${selected_backup}.state"
[ -f "$state_manifest" ] || die "selected backup state manifest does not exist"
[ -r "$state_manifest" ] || die "selected backup state manifest is not readable"

tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/campus-lms-source-safety.XXXXXX")"
source_state="$tmp_dir/source.state"
source_seed="$tmp_dir/source-seed.txt"
cleanup() {
    rm -rf -- "$tmp_dir"
}
trap cleanup EXIT

run_state() {
    "${COMPOSE[@]}" exec -T postgres \
        psql -X -q -v ON_ERROR_STOP=1 -A -t -P pager=off \
        -U "$source_user" -d "$source_database" < "$STATE_SQL"
}

run_seed_verifier() {
    "${COMPOSE[@]}" exec -T postgres \
        psql -X -v ON_ERROR_STOP=1 -P pager=off \
        -U "$source_user" -d "$source_database" < "$SEED_VERIFY_SQL"
}

source_exists="$("${COMPOSE[@]}" exec -T postgres \
    psql -X -q -v ON_ERROR_STOP=1 -A -t -P pager=off \
    -U "$source_user" -d postgres \
    -c "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '$source_database');")"
[ "$source_exists" = "t" ] || die "source database does not exist"

run_state > "$source_state"
if ! cmp -s "$state_manifest" "$source_state"; then
    echo "main_db_safety=fail" >&2
    diff -u "$state_manifest" "$source_state" >&2 || true
    die "source state differs from the stabilized backup manifest"
fi
cat "$source_state"

if ! run_seed_verifier > "$source_seed" 2>&1; then
    cat "$source_seed" >&2
    die "source Week 3 seed/invariant verification failed"
fi
cat "$source_seed"
echo "source_database=$source_database"
echo "source_exists=pass"
echo "source_manifest_comparison=pass"
echo "source_seed_verifier=pass"
echo "main_db_safety=pass"
