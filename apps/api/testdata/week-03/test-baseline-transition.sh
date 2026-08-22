#!/usr/bin/env bash
# apps/api/testdata/week-03/test-baseline-transition.sh

set -euo pipefail

COMPOSE=(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml)
DISPOSABLE_DB="week03_baseline_$(date +%s%N)_$$_${RANDOM}"
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cleanup() {
    local original_status="$1"

    if "${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d postgres -c "DROP DATABASE IF EXISTS ${DISPOSABLE_DB};" >/dev/null; then
        echo "cleanup=success"
        return "$original_status"
    fi

    echo "cleanup=failure" >&2
    if [ "$original_status" -eq 0 ]; then
        return 1
    fi
    return "$original_status"
}

on_exit() {
    local original_status="$?"
    trap - EXIT
    cleanup "$original_status"
    exit "$?"
}

echo "disposable_db=${DISPOSABLE_DB}"
"${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d postgres -c "CREATE DATABASE ${DISPOSABLE_DB};"
trap on_exit EXIT
echo "Applying committed migrations manually with ON_ERROR_STOP in a transaction..."
for migration in apps/api/migrations/000{1_tenant_identity_schema,2_academic_core_schema,3_auth_membership_schema,4_academic_term_time_range_check}.up.sql; do
    "${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -1 -U campus -d "$DISPOSABLE_DB" < "$migration"
done
migration_relation=$("${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$DISPOSABLE_DB" -tA -c "SELECT to_regclass('schema_migrations');")
if [ -n "$migration_relation" ]; then echo "FAIL: schema_migrations exists before baseline."; exit 1; fi
echo "pre_schema_migrations=absent"
TARGET_DB="$DISPOSABLE_DB" bash "$DIR/baseline-main-db.sh"
post_state=$("${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$DISPOSABLE_DB" -tA -c "SELECT to_regclass('schema_migrations');")
if [ "$post_state" != "schema_migrations" ]; then echo "FAIL: schema_migrations is absent after baseline."; exit 1; fi
echo "post_schema_migrations=present"
"${COMPOSE[@]}" exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "$DISPOSABLE_DB" -tA -c "SELECT 'version=' || version FROM schema_migrations UNION ALL SELECT 'dirty=' || CASE WHEN dirty THEN 'true' ELSE 'false' END FROM schema_migrations;"
