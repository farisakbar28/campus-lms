#!/usr/bin/env bash
# apps/api/testdata/week-03/test-fresh-migrations.sh
# Proves golang-migrate works gracefully on a fresh disposable database.

set -euo pipefail

DB_NAME="campus_lms_week03_migrate_verify_$(date +%s)"

echo "Creating disposable database: $DB_NAME"
docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d postgres -c "CREATE DATABASE $DB_NAME;"

# Ensure cleanup happens on exit
trap 'echo "Cleaning up disposable database..."; docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"' EXIT

echo "Fetching root connection URL template..."
# Dynamically build the URL for the disposable DB without printing the secret
COMPOSE_CMD="docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml"
$COMPOSE_CMD exec -T api sh -ec '
    if [ ! -x /go/bin/migrate ]; then
        go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3
    fi
    
    test -n "$MIGRATE_DATABASE_URL" || exit 1
    
    # Replace the DB name in the URL safely using sed
    TEST_URL=$(echo "$MIGRATE_DATABASE_URL" | sed "s|/campus_lms?|/$1?|")
    
    echo "Running migrate up..."
    /go/bin/migrate -path /src/migrations -database "$TEST_URL" up
    
    echo "Verifying version..."
    VERSION=$(/go/bin/migrate -path /src/migrations -database "$TEST_URL" version 2>&1)
    echo "$VERSION"
    echo "$VERSION" | grep -q "^4$" || exit 1
    
    echo "Running migrate down 1..."
    /go/bin/migrate -path /src/migrations -database "$TEST_URL" down 1
    
    echo "Verifying version..."
    VERSION=$(/go/bin/migrate -path /src/migrations -database "$TEST_URL" version 2>&1)
    echo "$VERSION"
    echo "$VERSION" | grep -q "^3$" || exit 1
' sh "$DB_NAME"

echo "Verifying academic_terms CHECK constraint is absent..."
$COMPOSE_CMD exec -T postgres psql -U campus -d "$DB_NAME" -tA -c "
SELECT COUNT(*) FROM information_schema.check_constraints WHERE constraint_name = 'academic_terms_valid_time_range';
" | grep -q "^0$" || { echo "FAIL: Constraint still exists"; exit 1; }

echo "Running migrate up 1 to restore state..."
$COMPOSE_CMD exec -T api sh -ec '
    TEST_URL=$(echo "$MIGRATE_DATABASE_URL" | sed "s|/campus_lms?|/$1?|")
    /go/bin/migrate -path /src/migrations -database "$TEST_URL" up 1
    
    VERSION=$(/go/bin/migrate -path /src/migrations -database "$TEST_URL" version 2>&1)
    echo "$VERSION"
    echo "$VERSION" | grep -q "^4$" || exit 1
' sh "$DB_NAME"

echo "Verifying academic_terms CHECK constraint is restored..."
$COMPOSE_CMD exec -T postgres psql -U campus -d "$DB_NAME" -tA -c "
SELECT COUNT(*) FROM information_schema.check_constraints WHERE constraint_name = 'academic_terms_valid_time_range';
" | grep -q "^1$" || { echo "FAIL: Constraint not restored"; exit 1; }

# Run final structural check against the disposable DB
echo "Running structural verification against disposable DB..."
CHECK_OUTPUT=$($COMPOSE_CMD exec -T postgres psql -U campus -d "$DB_NAME" -tA < apps/api/testdata/week-03/main-db-final-check.sql)

echo "$CHECK_OUTPUT" | grep -q "^tables|11$" || { echo "FAIL: table count != 11"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^rls_enabled|8$" || { echo "FAIL: RLS not enabled on 8 tables"; exit 1; }

echo "SUCCESS: Fresh migration test completed perfectly."
