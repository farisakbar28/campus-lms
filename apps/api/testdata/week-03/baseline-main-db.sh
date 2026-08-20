#!/usr/bin/env bash
# apps/api/testdata/week-03/baseline-main-db.sh
# Performs the guarded one-time transition baselining the existing main dev DB to migration version 4.

set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 1. Run strict precheck
bash "$DIR/baseline-main-db-precheck.sh"

# 2. Check if already baselined (schema_migrations exists and is version 4)
MIGRATE_STATE=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d campus_lms -tA -c "SELECT to_regclass('schema_migrations');")
if [ "$MIGRATE_STATE" = "schema_migrations" ]; then
    echo "Database is already baselined."
    exit 0
fi

echo "Applying force 4 metadata..."

# 3. Force version 4
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T api sh -ec '
    if [ ! -x /go/bin/migrate ]; then
        go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3
    fi
    test -n "$MIGRATE_DATABASE_URL" || exit 1
    exec /go/bin/migrate -path /src/migrations -database "$MIGRATE_DATABASE_URL" force 4
'

echo "Verify version..."
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T api sh -ec '
    exec /go/bin/migrate -path /src/migrations -database "$MIGRATE_DATABASE_URL" version
'

echo "Verify no-op on up..."
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T api sh -ec '
    exec /go/bin/migrate -path /src/migrations -database "$MIGRATE_DATABASE_URL" up
' 2>&1 | grep -q "no change" || { echo "Expected no change from migrate up"; exit 1; }

echo "Verify structural integrity remains intact..."
bash "$DIR/baseline-main-db-precheck.sh"

echo "Baseline successful."