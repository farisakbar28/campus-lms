#!/usr/bin/env bash
# apps/api/testdata/week-03/baseline-main-db.sh
# TARGET_DB is a non-secret PostgreSQL identifier; the migration URL stays inside api.

set -euo pipefail

COMPOSE=(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml)
TARGET_DB="${TARGET_DB:-campus_lms}"

case "$TARGET_DB" in
    ''|[!a-z_]*|*[!a-z0-9_]*) echo "FAIL: TARGET_DB must be a lowercase PostgreSQL identifier."; exit 1 ;;
esac

run_migrate() {
    "${COMPOSE[@]}" exec -T api sh -ec '
        set -eu
        target_db="$1"
        shift
        test -n "$MIGRATE_DATABASE_URL" || { echo "MIGRATE_DATABASE_URL is not set" >&2; exit 1; }
        case "$target_db" in ""|[!a-z_]*|*[!a-z0-9_]*) echo "target database name is invalid" >&2; exit 1;; esac
        case "$MIGRATE_DATABASE_URL" in postgres://*|postgresql://*) ;; *) echo "MIGRATE_DATABASE_URL must use a PostgreSQL URL" >&2; exit 1;; esac
        scheme=${MIGRATE_DATABASE_URL%%://*}
        remainder=${MIGRATE_DATABASE_URL#*://}
        case "$remainder" in */*) ;; *) echo "MIGRATE_DATABASE_URL has no database component" >&2; exit 1;; esac
        authority=${remainder%%/*}
        path_and_query=${remainder#*/}
        source_database=${path_and_query%%\?*}
        suffix=${path_and_query#"$source_database"}
        case "$authority:$source_database" in :*|*:|*/*) echo "MIGRATE_DATABASE_URL has an invalid database component" >&2; exit 1;; esac
        migrate_url="${scheme}://${authority}/${target_db}${suffix}"
        migrate_path="/go/bin/migrate"
        want_version="github.com/golang-migrate/migrate/v4 v4.18.3"
        if [ -x "$migrate_path" ] && ! go version -m "$migrate_path" 2>/dev/null | grep -q "$want_version"; then rm -f "$migrate_path"; fi
        if [ ! -x "$migrate_path" ]; then go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3; fi
        exec "$migrate_path" -path /src/migrations -database "$migrate_url" "$@"
    ' sh "$TARGET_DB" "$@"
}

run_precheck() {
    cat apps/api/testdata/week-03/baseline-main-db-precheck.sql |
        "${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$TARGET_DB"
}

echo "target=${TARGET_DB}"
echo "Running strict structural precheck..."
run_precheck
echo "structural_precheck=pass"

migration_relation=$("${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$TARGET_DB" -tA -c "SELECT to_regclass('schema_migrations');")
if [ "$migration_relation" = "schema_migrations" ]; then
    version=$("${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$TARGET_DB" -tA -c "SELECT version FROM schema_migrations LIMIT 1;")
    dirty=$("${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$TARGET_DB" -tA -c "SELECT dirty FROM schema_migrations LIMIT 1;")
    if [ "$version" != "4" ] || [ "$dirty" != "f" ]; then echo "FAIL: schema_migrations exists but is not version 4 and clean."; exit 1; fi
    echo "force_skipped=true"
else
    echo "baseline_force=4"
    run_migrate force 4
    version=$("${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$TARGET_DB" -tA -c "SELECT version FROM schema_migrations LIMIT 1;")
    dirty=$("${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$TARGET_DB" -tA -c "SELECT dirty FROM schema_migrations LIMIT 1;")
    if [ "$version" != "4" ] || [ "$dirty" != "f" ]; then echo "FAIL: force did not create a clean version 4 schema_migrations state."; exit 1; fi
    echo "force_skipped=false"
fi

echo "version=${version}"
echo "dirty=false"
up_output=$(run_migrate up 2>&1)
printf '%s\n' "$up_output"
case "$up_output" in *"no change"*) echo "migrate_up=no_change";; *) echo "FAIL: expected migrate up to report no change."; exit 1;; esac
echo "Running strict structural postcheck..."
run_precheck
echo "structural_check=pass"
echo "structural_postcheck=pass"
