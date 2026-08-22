#!/usr/bin/env bash
# apps/api/testdata/week-03/test-fresh-migrations.sh
# Proves golang-migrate works gracefully on a fresh disposable database.

set -euo pipefail

COMPOSE=(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml)
DB_NAME="campus_lms_week03_migrate_verify_$(date +%s)"

run_migrate() {
    "${COMPOSE[@]}" exec -T api sh -ec '
        set -eu
        target_db="$1"
        shift
        test -n "$MIGRATE_DATABASE_URL" || { echo "MIGRATE_DATABASE_URL is not set" >&2; exit 1; }
        case "$target_db" in ""|[!a-z_]*|*[!a-z0-9_]*) echo "target database name is invalid" >&2; exit 1;; esac
        case "$MIGRATE_DATABASE_URL" in postgres://*|postgresql://*) ;; *) echo "MIGRATE_DATABASE_URL must use a PostgreSQL URL" >&2; exit 1;; esac
        scheme=${MIGRATE_DATABASE_URL%%://*}; remainder=${MIGRATE_DATABASE_URL#*://}
        case "$remainder" in */*) ;; *) echo "MIGRATE_DATABASE_URL has no database component" >&2; exit 1;; esac
        authority=${remainder%%/*}; path_and_query=${remainder#*/}; source_database=${path_and_query%%\?*}; suffix=${path_and_query#"$source_database"}
        case "$authority:$source_database" in :*|*:|*/*) echo "MIGRATE_DATABASE_URL has an invalid database component" >&2; exit 1;; esac
        migrate_url="${scheme}://${authority}/${target_db}${suffix}"
        test -x /go/bin/migrate || go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3
        exec /go/bin/migrate -path /src/migrations -database "$migrate_url" "$@"
    ' sh "$DB_NAME" "$@"
}

echo "disposable_db=${DB_NAME}"
"${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d postgres -c "CREATE DATABASE $DB_NAME;"
trap '"${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"' EXIT
echo "Running migrate up..."
run_migrate up
version=$(run_migrate version 2>&1); printf '%s\n' "$version"; [ "$version" = "4" ] || { echo "FAIL: expected version 4"; exit 1; }
echo "Running migrate down 1..."
run_migrate down 1
version=$(run_migrate version 2>&1); printf '%s\n' "$version"; [ "$version" = "3" ] || { echo "FAIL: expected version 3"; exit 1; }
echo "Verifying academic_terms CHECK constraint is absent..."
"${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$DB_NAME" -tA -c "SELECT count(*) FROM information_schema.check_constraints WHERE constraint_name = 'academic_terms_valid_time_range';" | grep -qx '0'
echo "Running migrate up 1 to restore state..."
run_migrate up 1
version=$(run_migrate version 2>&1); printf '%s\n' "$version"; [ "$version" = "4" ] || { echo "FAIL: expected version 4"; exit 1; }
echo "Verifying academic_terms CHECK constraint is restored..."
"${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$DB_NAME" -tA -c "SELECT count(*) FROM information_schema.check_constraints WHERE constraint_name = 'academic_terms_valid_time_range';" | grep -qx '1'
echo "Running strict structural verification against disposable DB..."
cat apps/api/testdata/week-03/baseline-main-db-precheck.sql | "${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$DB_NAME"
echo "structural_check=pass"
