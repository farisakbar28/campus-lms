#!/usr/bin/env bash
#
# Operator-run deployment for immutable API, migrator, and Caddy images.
#
# The caller supplies immutable image references and separate VM-local API and
# migrator environment files. This script never prints those files, performs no
# cloud operation, and rolls back application images only. Database schema
# rollback is intentionally absent: production migrations must be
# expand-contract compatible before this script is used.

set -euo pipefail
umask 077

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
COMPOSE_FILE="$REPO_ROOT/deploy/compose/docker-compose.prod.yml"
COMPOSE_PROJECT="${CAMPUS_LMS_COMPOSE_PROJECT:-campus-lms-prod}"
STATE_DIR="${CAMPUS_LMS_STATE_DIR:-/var/lib/campus-lms}"
STATE_FILE="${CAMPUS_LMS_STATE_FILE:-$STATE_DIR/deployment-state}"
INTERNAL_API_READINESS_URL="${CAMPUS_LMS_INTERNAL_API_READINESS_URL:-http://127.0.0.1:8080/readyz}"
CADDY_READINESS_ADDRESS="${CAMPUS_LMS_CADDY_READINESS_ADDRESS:-127.0.0.1}"
CADDY_READINESS_PORT="${CAMPUS_LMS_CADDY_READINESS_PORT:-8443}"
READY_TIMEOUT_SECONDS="${CAMPUS_LMS_READY_TIMEOUT_SECONDS:-60}"
READY_INTERVAL_SECONDS="${CAMPUS_LMS_READY_INTERVAL_SECONDS:-2}"
ROLLBACK_FAILURE_THRESHOLD=3

die() {
	echo "ERROR: $*" >&2
	exit 1
}

require_nonempty() {
	local name="$1"
	local value="${!name-}"
	[ -n "$value" ] || die "$name is required"
}

require_digest_image() {
	local name="$1"
	local value="${!name-}"
	require_nonempty "$name"
	[[ "$value" =~ ^[^[:space:]]+@sha256:[0-9a-fA-F]{64}$ ]] ||
		die "$name must be an immutable image digest"
	[[ "$value" != *:latest@* ]] || die "$name must not use latest"
}

require_positive_integer() {
	local name="$1"
	local value="${!name-}"
	[[ "$value" =~ ^[1-9][0-9]*$ ]] || die "$name must be a positive integer"
}

require_secret_env_file() {
	local name="$1"
	local path="${!name-}"

	require_nonempty "$name"
	[ -f "$path" ] || die "$name is missing"
	[ -r "$path" ] || die "$name is not readable; production deployment requires root privileges"
	[ "$(stat -c '%u' "$path" 2>/dev/null || true)" = "0" ] ||
		die "$name must be owned by UID 0"
	[ "$(stat -c '%a' "$path" 2>/dev/null || true)" = "600" ] ||
		die "$name must have mode 0600"
}

require_root_secret_file() {
	local name="$1"
	local path="${!name-}"

	require_nonempty "$name"
	[ -f "$path" ] || die "$name is missing"
	[ -r "$path" ] || die "$name is not readable; production deployment requires root privileges"
	[ "$(stat -c '%u' "$path" 2>/dev/null || true)" = "0" ] ||
		die "$name must be owned by UID 0"
	[ "$(stat -c '%a' "$path" 2>/dev/null || true)" = "600" ] ||
		die "$name must have mode 0600"
}

require_readable_file() {
	local name="$1"
	local path="${!name-}"

	require_nonempty "$name"
	[ -f "$path" ] || die "$name is missing"
	[ -r "$path" ] || die "$name is not readable"
}

compose() {
	CAMPUS_LMS_API_IMAGE="$CAMPUS_LMS_API_IMAGE" \
	CAMPUS_LMS_MIGRATOR_IMAGE="$CAMPUS_LMS_MIGRATOR_IMAGE" \
	CAMPUS_LMS_CADDY_IMAGE="$CAMPUS_LMS_CADDY_IMAGE" \
	CAMPUS_LMS_HOSTNAME="$CAMPUS_LMS_HOSTNAME" \
	CAMPUS_LMS_API_ENV_FILE="$CAMPUS_LMS_API_ENV_FILE" \
	CAMPUS_LMS_MIGRATOR_ENV_FILE="$CAMPUS_LMS_MIGRATOR_ENV_FILE" \
	CAMPUS_LMS_API_MEMORY_LIMIT="$CAMPUS_LMS_API_MEMORY_LIMIT" \
	CAMPUS_LMS_MIGRATOR_MEMORY_LIMIT="$CAMPUS_LMS_MIGRATOR_MEMORY_LIMIT" \
	CAMPUS_LMS_CADDY_MEMORY_LIMIT="$CAMPUS_LMS_CADDY_MEMORY_LIMIT" \
	CAMPUS_LMS_CADDY_ORIGIN_CERT_FILE="$CAMPUS_LMS_CADDY_ORIGIN_CERT_FILE" \
	CAMPUS_LMS_CADDY_ORIGIN_KEY_FILE="$CAMPUS_LMS_CADDY_ORIGIN_KEY_FILE" \
	docker compose -p "$COMPOSE_PROJECT" -f "$COMPOSE_FILE" "$@"
}

validate_inputs() {
	require_digest_image CAMPUS_LMS_API_IMAGE
	require_digest_image CAMPUS_LMS_MIGRATOR_IMAGE
	require_digest_image CAMPUS_LMS_CADDY_IMAGE
	require_nonempty CAMPUS_LMS_HOSTNAME
	[[ "$CAMPUS_LMS_HOSTNAME" =~ ^[A-Za-z0-9.-]+$ ]] || die "CAMPUS_LMS_HOSTNAME is invalid"
	require_secret_env_file CAMPUS_LMS_API_ENV_FILE
	require_secret_env_file CAMPUS_LMS_MIGRATOR_ENV_FILE
	require_root_secret_file CAMPUS_LMS_CADDY_ORIGIN_KEY_FILE
	require_readable_file CAMPUS_LMS_CADDY_ORIGIN_CERT_FILE
	require_readable_file CAMPUS_LMS_CADDY_ORIGIN_CA_ROOT_FILE
	require_nonempty CAMPUS_LMS_API_MEMORY_LIMIT
	require_nonempty CAMPUS_LMS_MIGRATOR_MEMORY_LIMIT
	require_nonempty CAMPUS_LMS_CADDY_MEMORY_LIMIT
	require_nonempty CADDY_READINESS_ADDRESS
	require_positive_integer CADDY_READINESS_PORT
	if [ -n "${CAMPUS_LMS_CADDY_READINESS_HOST:-}" ]; then
		[[ "$CAMPUS_LMS_CADDY_READINESS_HOST" =~ ^[A-Za-z0-9.-]+$ ]] ||
			die "CAMPUS_LMS_CADDY_READINESS_HOST is invalid"
	fi
	require_positive_integer READY_TIMEOUT_SECONDS
	require_positive_integer READY_INTERVAL_SECONDS
	command -v docker >/dev/null 2>&1 || die "docker is required"
	command -v curl >/dev/null 2>&1 || die "curl is required for readiness checks"
	[ -f "$COMPOSE_FILE" ] || die "production Compose file is missing"
	compose config --quiet || die "production Compose configuration is invalid"
}

write_state() {
	local target="$1"
	local temporary

	mkdir -p "$STATE_DIR"
	temporary="$(mktemp "$STATE_DIR/.deployment-state.XXXXXX")"
	{
		printf 'api_image=%s\n' "$CAMPUS_LMS_API_IMAGE"
		printf 'migrator_image=%s\n' "$CAMPUS_LMS_MIGRATOR_IMAGE"
		printf 'caddy_image=%s\n' "$CAMPUS_LMS_CADDY_IMAGE"
		printf 'hostname=%s\n' "$CAMPUS_LMS_HOSTNAME"
	} > "$temporary"
	chmod 0600 "$temporary"
	mv -f "$temporary" "$target"
}

load_previous_state() {
	local key value
	local found_api=0
	local found_migrator=0
	local found_caddy=0
	local found_hostname=0

	[ -f "$STATE_FILE" ] || return 1
	while IFS='=' read -r key value; do
		case "$key" in
		api_image)
			CAMPUS_LMS_API_IMAGE="$value"
			found_api=1
			;;
		migrator_image)
			CAMPUS_LMS_MIGRATOR_IMAGE="$value"
			found_migrator=1
			;;
		caddy_image)
			CAMPUS_LMS_CADDY_IMAGE="$value"
			found_caddy=1
			;;
		hostname)
			CAMPUS_LMS_HOSTNAME="$value"
			found_hostname=1
			;;
		'')
			;;
		*)
			return 1
			;;
		esac
	done < "$STATE_FILE"

	[ "$found_api" -eq 1 ] && [ "$found_migrator" -eq 1 ] &&
		[ "$found_caddy" -eq 1 ] && [ "$found_hostname" -eq 1 ] || return 1
	require_digest_image CAMPUS_LMS_API_IMAGE
	require_digest_image CAMPUS_LMS_MIGRATOR_IMAGE
	require_digest_image CAMPUS_LMS_CADDY_IMAGE
	[[ "$CAMPUS_LMS_HOSTNAME" =~ ^[A-Za-z0-9.-]+$ ]] || return 1
}

check_internal_api_readiness() {
	curl --fail --silent --show-error --max-time 2 --output /dev/null "$INTERNAL_API_READINESS_URL"
}

check_caddy_readiness() {
	local host="${CAMPUS_LMS_CADDY_READINESS_HOST:-$CAMPUS_LMS_HOSTNAME}"

	curl --fail --silent --show-error --max-time 2 --output /dev/null \
		--cacert "$CAMPUS_LMS_CADDY_ORIGIN_CA_ROOT_FILE" \
		--resolve "${host}:${CADDY_READINESS_PORT}:${CADDY_READINESS_ADDRESS}" \
		"https://${host}:${CADDY_READINESS_PORT}/readyz"
}

wait_for_readiness() {
	local deadline=$((SECONDS + READY_TIMEOUT_SECONDS))
	local failures=0
	local api_ready=0
	local caddy_ready=0

	while (( SECONDS < deadline )); do
		api_ready=0
		caddy_ready=0
		if check_internal_api_readiness; then
			api_ready=1
			echo "api_readiness=pass"
		else
			echo "api_readiness=fail"
		fi
		if check_caddy_readiness; then
			caddy_ready=1
			echo "caddy_readiness=pass"
		else
			echo "caddy_readiness=fail"
		fi
		if (( api_ready == 1 && caddy_ready == 1 )); then
			echo "readiness=pass"
			return 0
		fi
		failures=$((failures + 1))
		if (( failures >= ROLLBACK_FAILURE_THRESHOLD )); then
			echo "readiness=fail"
			return 1
		fi
		sleep "$READY_INTERVAL_SECONDS"
	done

	echo "readiness=timeout"
	return 1
}

rollback_application() {
	if ! load_previous_state; then
		echo "rollback=unavailable_previous_state" >&2
		return 1
	fi

	if ! compose pull api caddy; then
		echo "rollback=pull_failed" >&2
		return 1
	fi
	if ! compose up -d --no-deps api caddy; then
		echo "rollback=update_failed" >&2
		return 1
	fi
	if wait_for_readiness; then
		echo "rollback=pass"
		return 0
	fi

	echo "rollback=fail" >&2
	return 1
}

main() {
	local started_at
	local duration_ms
	local previous_state_available=0
	local candidate_api_image="${CAMPUS_LMS_API_IMAGE-}"
	local candidate_migrator_image="${CAMPUS_LMS_MIGRATOR_IMAGE-}"
	local candidate_caddy_image="${CAMPUS_LMS_CADDY_IMAGE-}"
	local candidate_hostname="${CAMPUS_LMS_HOSTNAME-}"

	validate_inputs
	started_at="$(date +%s%N)"

	if load_previous_state; then
		previous_state_available=1
	fi
	# Restore candidate values after reading the previous state for validation.
	CAMPUS_LMS_API_IMAGE="${CAMPUS_LMS_TARGET_API_IMAGE:-$candidate_api_image}"
	CAMPUS_LMS_MIGRATOR_IMAGE="${CAMPUS_LMS_TARGET_MIGRATOR_IMAGE:-$candidate_migrator_image}"
	CAMPUS_LMS_CADDY_IMAGE="${CAMPUS_LMS_TARGET_CADDY_IMAGE:-$candidate_caddy_image}"
	CAMPUS_LMS_HOSTNAME="$candidate_hostname"
	validate_inputs

	if (( previous_state_available == 1 )); then
		cp -- "$STATE_FILE" "${STATE_FILE}.previous"
		chmod 0600 "${STATE_FILE}.previous"
	fi

	compose pull api migrator caddy
	compose run --rm --no-deps migrator
	compose up -d --no-deps api caddy

	if ! wait_for_readiness; then
		echo "candidate_readiness=fail" >&2
		if (( previous_state_available == 1 )) && rollback_application; then
			die "candidate failed readiness; application rollback completed"
		fi
		die "candidate failed readiness and application rollback was unavailable or failed"
	fi

	write_state "$STATE_FILE"
	duration_ms=$(( ($(date +%s%N) - started_at + 999999) / 1000000 ))
	echo "candidate_readiness=pass"
	echo "deploy_duration_ms=$duration_ms"
	echo "deployment=pass"
}

main "$@"
