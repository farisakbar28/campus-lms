#!/usr/bin/env bash
# apps/api/testdata/week-03/baseline-main-db-precheck.sh
# Validates that the existing main dev database matches the expected v4 structural state
# BEFORE baselining it to version 4 with golang-migrate force.

set -euo pipefail

# Run the strict structural precheck
echo "Running strict structural precheck..."

docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d campus_lms -f /src/testdata/week-03/baseline-main-db-precheck.sql

echo "Structural assertions PASS."

# Prove schema_migrations does not exist (or if it does, it must be exactly version 4 and clean)
echo "Checking schema_migrations state..."
MIGRATE_STATE=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d campus_lms -tA -c "SELECT to_regclass('schema_migrations');")

if [ "$MIGRATE_STATE" = "schema_migrations" ]; then
    echo "schema_migrations exists. Verifying state..."
    # Prefer reading version directly from DB to avoid CLI dependencies if schema is already tracked
    VERSION=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d campus_lms -tA -c "SELECT version FROM schema_migrations LIMIT 1;")
    DIRTY=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d campus_lms -tA -c "SELECT dirty FROM schema_migrations LIMIT 1;")
    
    if [ "$VERSION" = "4" ] && [ "$DIRTY" = "f" ]; then
        echo "Already baselined correctly to version 4."
        exit 0
    else
        echo "FAIL: schema_migrations exists but state is not 4 clean."
        echo "Version: $VERSION, Dirty: $DIRTY"
        exit 1
    fi
fi

echo "Precheck completed successfully. Database is ready for baseline force 4."
