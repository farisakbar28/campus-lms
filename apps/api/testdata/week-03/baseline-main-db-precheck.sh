#!/usr/bin/env bash
# apps/api/testdata/week-03/baseline-main-db-precheck.sql equivalent runner
# Validates that the existing main dev database matches the expected v4 structural state
# BEFORE baselining it to version 4 with golang-migrate force.

set -euo pipefail

# Run the final check logic to assert invariants
echo "Running strict structural precheck..."

# We expect exact values from the final verification
CHECK_OUTPUT=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d campus_lms -tA < apps/api/testdata/week-03/main-db-final-check.sql)

# Check all assertions
echo "$CHECK_OUTPUT" | grep -q "^tables|11$" || { echo "FAIL: table count != 11"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^rls_enabled|8$" || { echo "FAIL: RLS not enabled on 8 tables"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^a7_fk|1$" || { echo "FAIL: A7 FK missing"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^term_range_check|1$" || { echo "FAIL: academic_term check missing"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^audit_policy_total|2$" || { echo "FAIL: audit policy count != 2"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^audit_policy_select|1$" || { echo "FAIL: audit policy SELECT != 1"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^audit_policy_insert|1$" || { echo "FAIL: audit policy INSERT != 1"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^audit_policy_update|0$" || { echo "FAIL: audit policy UPDATE != 0"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^audit_policy_delete|0$" || { echo "FAIL: audit policy DELETE != 0"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^audit_policy_all|0$" || { echo "FAIL: audit policy ALL != 0"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^rls_verifier|0$" || { echo "FAIL: rls_verifier role exists"; exit 1; }
echo "$CHECK_OUTPUT" | grep -q "^fixture_rows_remaining|0$" || { echo "FAIL: fixture users exist"; exit 1; }

echo "Structural assertions PASS."

# Prove schema_migrations does not exist (or if it does, it must be exactly version 4 and clean)
echo "Checking schema_migrations state..."
MIGRATE_STATE=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T postgres psql -U campus -d campus_lms -tA -c "SELECT to_regclass('schema_migrations');")

if [ "$MIGRATE_STATE" = "schema_migrations" ]; then
    echo "schema_migrations exists. Verifying state..."
    VERSION_OUT=$(docker compose -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec -T api sh -ec '
        test -n "$MIGRATE_DATABASE_URL" || exit 1
        exec /go/bin/migrate -path /src/migrations -database "$MIGRATE_DATABASE_URL" version
    ' 2>&1 || true)
    
    if echo "$VERSION_OUT" | grep -q "^4$"; then
        echo "Already baselined correctly to version 4."
        exit 0
    else
        echo "FAIL: schema_migrations exists but version is not 4 clean."
        echo "Version output: $VERSION_OUT"
        exit 1
    fi
fi

echo "Precheck completed successfully. Database is ready for baseline force 4."