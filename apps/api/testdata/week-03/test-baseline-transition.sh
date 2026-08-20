#!/usr/bin/env bash
# apps/api/testdata/week-03/test-baseline-transition.sh

set -euo pipefail

DISPOSABLE_DB="test_baseline_$(date +%s)"
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Trap for cleanup
cleanup() {
    echo "Cleaning up disposable database ${DISPOSABLE_DB}..."
    docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d postgres -c "DROP DATABASE IF EXISTS ${DISPOSABLE_DB};" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "Creating disposable database ${DISPOSABLE_DB}..."
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d postgres -c "CREATE DATABASE ${DISPOSABLE_DB};"

echo "Applying migrations manually..."
cat apps/api/migrations/0001_tenant_identity_schema.up.sql | docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d "${DISPOSABLE_DB}"
cat apps/api/migrations/0002_academic_core_schema.up.sql | docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d "${DISPOSABLE_DB}"
cat apps/api/migrations/0003_auth_membership_schema.up.sql | docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d "${DISPOSABLE_DB}"
cat apps/api/migrations/0004_academic_term_time_range_check.up.sql | docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d "${DISPOSABLE_DB}"

# Test schema_migrations absent
echo "Verifying schema_migrations absent..."
MIGRATE_STATE=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "${DISPOSABLE_DB}" -tA -c "SELECT to_regclass('schema_migrations');")
if [ "$MIGRATE_STATE" = "schema_migrations" ]; then
    echo "FAIL: schema_migrations unexpectedly exists before baseline."
    exit 1
fi

echo "Running guarded baseline..."
TARGET_DB="${DISPOSABLE_DB}" TARGET_MIGRATE_URL="postgres://campus:campus@postgres:5432/${DISPOSABLE_DB}?sslmode=disable" bash "$DIR/baseline-main-db.sh"

echo "Disposable baseline transition test PASS."
