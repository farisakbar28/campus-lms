#!/usr/bin/env bash
# apps/api/testdata/week-03/baseline-main-db.sh
# Performs the guarded one-time transition baselining a DB to migration version 4.
# Accepts TARGET_DB and TARGET_MIGRATE_URL environment variables for disposable DB testing.

set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DB="${TARGET_DB:-campus_lms}"

# 1. Run strict precheck
echo "Running strict structural precheck on ${TARGET_DB}..."
cat apps/api/testdata/week-03/baseline-main-db-precheck.sql | docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "${TARGET_DB}"
echo "Structural assertions PASS."

# 2. Check if already baselined (schema_migrations exists and is version 4)
echo "Checking schema_migrations state on ${TARGET_DB}..."
MIGRATE_STATE=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "${TARGET_DB}" -tA -c "SELECT to_regclass('schema_migrations');")

if [ "$MIGRATE_STATE" = "schema_migrations" ]; then
    VERSION=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "${TARGET_DB}" -tA -c "SELECT version FROM schema_migrations LIMIT 1;")
    DIRTY=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "${TARGET_DB}" -tA -c "SELECT dirty FROM schema_migrations LIMIT 1;")
    
    if [ "$VERSION" = "4" ] && [ "$DIRTY" = "f" ]; then
        echo "Database ${TARGET_DB} is already baselined correctly to version 4."
        exit 0
    else
        echo "FAIL: schema_migrations exists on ${TARGET_DB} but state is not 4 clean."
        echo "Version: $VERSION, Dirty: $DIRTY"
        exit 1
    fi
fi

echo "Applying force 4 metadata to ${TARGET_DB}..."

# 3. Force version 4
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T api sh -ec '
    MIGRATE_PATH="/go/bin/migrate"
    WANT_VER="github.com/golang-migrate/migrate/v4 v4.18.3"
    if [ -x "$MIGRATE_PATH" ]; then
        CUR_VER=$(go version -m "$MIGRATE_PATH" 2>/dev/null | grep "$WANT_VER" || true)
        if [ -z "$CUR_VER" ]; then
            rm -f "$MIGRATE_PATH"
        fi
    fi
    if [ ! -x "$MIGRATE_PATH" ]; then
        go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3
    fi
    URL="${TARGET_MIGRATE_URL:-$MIGRATE_DATABASE_URL}"
    test -n "$URL" || exit 1
    exec /go/bin/migrate -path /src/migrations -database "$URL" force 4
'

echo "Verify version..."
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T api sh -ec '
    URL="${TARGET_MIGRATE_URL:-$MIGRATE_DATABASE_URL}"
    exec /go/bin/migrate -path /src/migrations -database "$URL" version
'

echo "Verify no-op on up..."
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T api sh -ec '
    URL="${TARGET_MIGRATE_URL:-$MIGRATE_DATABASE_URL}"
    exec /go/bin/migrate -path /src/migrations -database "$URL" up
' 2>&1 | grep -q "no change" || { echo "Expected no change from migrate up"; exit 1; }

echo "Verify structural integrity remains intact..."
cat apps/api/testdata/week-03/baseline-main-db-precheck.sql | docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -v ON_ERROR_STOP=1 -U campus -d "${TARGET_DB}"

echo "Baseline successful on ${TARGET_DB}."
