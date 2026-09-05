#!/usr/bin/env bash
# Disposable verification for the production migrator image. The explicit local
# mode is test-only: docker-compose.prod.yml never sets it, so production always
# enforces a TLS PostgreSQL URL before invoking golang-migrate.

set -euo pipefail

IMAGE="${MIGRATOR_IMAGE:-campus-lms-migrator:phase4b}"
NETWORK="campus-lms-phase4b-migration"
POSTGRES="campus-lms-phase4b-postgres"
HARNESS="$(mktemp -d /tmp/campus-lms-migrator-test.XXXXXX)"
FAKE_MIGRATE="$HARNESS/migrate"
TLS_ERROR='ERROR: production migration requires exactly one supported sslmode'

cleanup() {
	docker rm -f "$POSTGRES" >/dev/null 2>&1 || true
	docker network rm "$NETWORK" >/dev/null 2>&1 || true
	rm -rf -- "$HARNESS"
}
trap cleanup EXIT

cat > "$FAKE_MIGRATE" <<'EOF'
#!/bin/sh
printf '%s\n' fake-migrate-invoked
EOF
chmod 0555 "$FAKE_MIGRATE"

run_tls_case() {
	local name="$1"
	local database_url="$2"
	local expected="$3"
	local output
	local status

	set +e
	output="$(docker run --rm \
		-v "$FAKE_MIGRATE:/usr/local/bin/migrate:ro" \
		-e "MIGRATE_DATABASE_URL=$database_url" \
		"$IMAGE" 2>&1)"
	status=$?
	set -e
	if [ "$expected" = "accept" ]; then
		test "$status" -eq 0
		grep -Fqx 'fake-migrate-invoked' <<<"$output"
	else
		test "$status" -ne 0
		grep -Fqx "$TLS_ERROR" <<<"$output"
	fi
	if grep -Fq 'db.example' <<<"$output" || grep -Fq 'user' <<<"$output"; then
		echo "tls_case=$name secret_safe=false" >&2
		exit 1
	fi
	printf 'tls_case=%s result=%sed\n' "$name" "$expected"
}

run_tls_case missing_sslmode 'postgres://user@db.example/campus_lms?connect_timeout=1' reject
run_tls_case disable 'postgres://user@db.example/campus_lms?sslmode=disable' reject
run_tls_case prefer 'postgres://user@db.example/campus_lms?sslmode=prefer' reject
run_tls_case unknown 'postgres://user@db.example/campus_lms?sslmode=unknown' reject
run_tls_case duplicate 'postgres://user@db.example/campus_lms?sslmode=require&sslmode=verify-full' reject
run_tls_case require 'postgres://user@db.example/campus_lms?sslmode=require' accept
run_tls_case verify_ca 'postgres://user@db.example/campus_lms?sslmode=verify-ca' accept
run_tls_case verify_full 'postgres://user@db.example/campus_lms?sslmode=verify-full' accept

cleanup
docker network create "$NETWORK" >/dev/null
docker run -d --name "$POSTGRES" --network "$NETWORK" --network-alias postgres \
	-e POSTGRES_USER=campus \
	-e POSTGRES_DB=campus_lms \
	-e POSTGRES_HOST_AUTH_METHOD=trust \
	postgres:16.14-alpine3.23 >/dev/null

ready=0
for attempt in $(seq 1 45); do
	if docker exec "$POSTGRES" psql -U campus -d campus_lms -c 'SELECT 1;' >/dev/null 2>&1; then
		ready=1
		break
	fi
	sleep 1
done
test "$ready" = 1
database_url='postgres://campus@postgres:5432/campus_lms?sslmode=disable'

docker run --rm --network "$NETWORK" \
	-e MIGRATOR_LOCAL_TEST_MODE=1 \
	-e "MIGRATE_DATABASE_URL=$database_url" \
	"$IMAGE"
first_version="$(docker exec "$POSTGRES" psql -U campus -d campus_lms -Atc "SELECT version || '|' || dirty FROM schema_migrations;")"
test "$first_version" = "6|false"
printf 'first_up_version=%s\n' "$first_version"

second_output="$(docker run --rm --network "$NETWORK" \
	-e MIGRATOR_LOCAL_TEST_MODE=1 \
	-e "MIGRATE_DATABASE_URL=$database_url" \
	"$IMAGE" 2>&1)"
grep -Fqx 'no change' <<<"$second_output"
echo 'second_up_output=no change'
second_version="$(docker exec "$POSTGRES" psql -U campus -d campus_lms -Atc "SELECT version || '|' || dirty FROM schema_migrations;")"
test "$second_version" = "6|false"
printf 'second_up_version=%s\n' "$second_version"

set +e
docker run --rm --network "$NETWORK" \
	-e MIGRATOR_LOCAL_TEST_MODE=1 \
	-e 'MIGRATE_DATABASE_URL=postgres://campus@invalid.invalid/campus_lms?sslmode=disable' \
	"$IMAGE" >/dev/null 2>&1
failure_status=$?
set -e
test "$failure_status" -ne 0
printf 'local_migration_failure_exit=%s\n' "$failure_status"

if grep -Eq '(^|[[:space:]])down([[:space:]]|$)' apps/api/migrate-entrypoint.sh; then
	echo 'production_migration_down_command=present' >&2
	exit 1
fi
echo 'production_migration_down_command=absent'
echo 'database_target=disposable_local_postgres'
echo 'migrator_tls_and_migration_verification=pass'
