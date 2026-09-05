#!/usr/bin/env bash
# Reproducible authoritative parser check for the production Compose definition.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
COMPOSE_FILE="$REPO_ROOT/deploy/compose/docker-compose.prod.yml"
TEMPORARY="$(mktemp -d /tmp/campus-lms-prod-compose.XXXXXX)"
API_ENV_FILE="$TEMPORARY/api.env"
MIGRATOR_ENV_FILE="$TEMPORARY/migrator.env"
ORIGIN_CERT_FILE="$TEMPORARY/origin-cert.pem"
ORIGIN_KEY_FILE="$TEMPORARY/origin-key.pem"

cleanup() { rm -rf -- "$TEMPORARY"; }
trap cleanup EXIT
printf '%s\n' 'API_ENV_SENTINEL=api-only' > "$API_ENV_FILE"
printf '%s\n' 'MIGRATOR_ENV_SENTINEL=migrator-only' > "$MIGRATOR_ENV_FILE"
printf '%s\n' 'test certificate placeholder' > "$ORIGIN_CERT_FILE"
printf '%s\n' 'test private key placeholder' > "$ORIGIN_KEY_FILE"
chmod 0600 "$API_ENV_FILE" "$MIGRATOR_ENV_FILE" "$ORIGIN_KEY_FILE"

export CAMPUS_LMS_API_IMAGE='ghcr.io/example/campus-lms-api@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
export CAMPUS_LMS_MIGRATOR_IMAGE='ghcr.io/example/campus-lms-migrator@sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'
export CAMPUS_LMS_CADDY_IMAGE='caddy:2.10.2-alpine@sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc'
export CAMPUS_LMS_HOSTNAME='api.example.invalid'
export CAMPUS_LMS_API_ENV_FILE="$API_ENV_FILE"
export CAMPUS_LMS_MIGRATOR_ENV_FILE="$MIGRATOR_ENV_FILE"
export CAMPUS_LMS_CADDY_ORIGIN_CERT_FILE="$ORIGIN_CERT_FILE"
export CAMPUS_LMS_CADDY_ORIGIN_KEY_FILE="$ORIGIN_KEY_FILE"
export CAMPUS_LMS_API_MEMORY_LIMIT=256m
export CAMPUS_LMS_MIGRATOR_MEMORY_LIMIT=128m
export CAMPUS_LMS_CADDY_MEMORY_LIMIT=64m

docker compose -p campus-lms-phase4b -f "$COMPOSE_FILE" config --quiet
parsed="$(docker compose -p campus-lms-phase4b -f "$COMPOSE_FILE" config --format json)"
services="$(jq -r '.services | keys | join(",")' <<<"$parsed")"
api_ports="$(jq -c '.services.api.ports | map({host_ip, target, published, protocol})' <<<"$parsed")"
caddy_ports="$(jq -c '.services.caddy.ports | map({host_ip, target, published, protocol})' <<<"$parsed")"
api_environment="$(jq -c '.services.api.environment' <<<"$parsed")"
migrator_environment="$(jq -c '.services.migrator.environment' <<<"$parsed")"
logging="$(jq -c '[.services | to_entries[] | {service:.key,logging:.value.logging}]' <<<"$parsed")"

test "$services" = 'api,caddy,migrator'
test "$api_ports" = '[{"host_ip":"127.0.0.1","target":8080,"published":"8080","protocol":"tcp"}]'
test "$caddy_ports" = '[{"host_ip":"127.0.0.1","target":443,"published":"8443","protocol":"tcp"}]'
grep -Fq '"API_ENV_SENTINEL":"api-only"' <<<"$api_environment"
if grep -Fq '"MIGRATOR_ENV_SENTINEL":"migrator-only"' <<<"$api_environment"; then exit 1; fi
grep -Fq '"MIGRATOR_ENV_SENTINEL":"migrator-only"' <<<"$migrator_environment"
if grep -Fq '"API_ENV_SENTINEL":"api-only"' <<<"$migrator_environment"; then exit 1; fi
test "$logging" = '[{"service":"api","logging":{"driver":"json-file","options":{"max-file":"5","max-size":"10m"}}},{"service":"caddy","logging":{"driver":"json-file","options":{"max-file":"5","max-size":"10m"}}},{"service":"migrator","logging":{"driver":"json-file","options":{"max-file":"5","max-size":"10m"}}}]'

printf 'services=%s\n' "$services"
printf 'api_ports=%s\n' "$api_ports"
printf 'caddy_ports=%s\n' "$caddy_ports"
printf 'api_env_boundary=%s\n' "$api_environment"
printf 'migrator_env_boundary=%s\n' "$migrator_environment"
printf 'logging=%s\n' "$logging"
echo 'excluded_services=postgres,redis,frontend,minio,observability,ai,worker'
echo 'production_compose_verification=pass'
