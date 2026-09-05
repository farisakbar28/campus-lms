#!/usr/bin/env bash
# Local safety tests for deploy.sh. Docker, curl, and stat are replaced with
# fakes, so no service, registry, secret file, or cloud resource is contacted.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_SCRIPT="$SCRIPT_DIR/deploy.sh"
HARNESS="$(mktemp -d /tmp/campus-lms-deploy-safety.XXXXXX)"
ORIGIN_CERT_FILE="$HARNESS/origin-cert.pem"
ORIGIN_KEY_FILE="$HARNESS/origin-key.pem"
ORIGIN_CA_ROOT_FILE="$HARNESS/origin-ca.pem"
PREVIOUS_API='ghcr.io/example/campus-lms-api@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
PREVIOUS_MIGRATOR='ghcr.io/example/campus-lms-migrator@sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'
PREVIOUS_CADDY='caddy@sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc'
CANDIDATE_API='ghcr.io/example/campus-lms-api@sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd'
CANDIDATE_MIGRATOR='ghcr.io/example/campus-lms-migrator@sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee'
CANDIDATE_CADDY='caddy@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff'

cleanup() { rm -rf -- "$HARNESS"; }
trap cleanup EXIT

bash -n "$DEPLOY_SCRIPT"
test "$PREVIOUS_API" != "$CANDIDATE_API"
test "$PREVIOUS_MIGRATOR" != "$CANDIDATE_MIGRATOR"
test "$PREVIOUS_CADDY" != "$CANDIDATE_CADDY"

mkdir -p "$HARNESS/bin" "$HARNESS/state"
api_env="$HARNESS/api.env"
migrator_env="$HARNESS/migrator.env"
touch "$api_env" "$migrator_env"
chmod 0600 "$api_env" "$migrator_env"
printf '%s\n' 'test certificate placeholder' > "$ORIGIN_CERT_FILE"
printf '%s\n' 'test private key placeholder' > "$ORIGIN_KEY_FILE"
printf '%s\n' 'test CA placeholder' > "$ORIGIN_CA_ROOT_FILE"
chmod 0600 "$ORIGIN_KEY_FILE"

cat > "$HARNESS/bin/stat" <<'EOF'
#!/usr/bin/env bash
case "${2:-}" in
%u) printf '%s\n' "${FAKE_STAT_UID:-0}" ;;
%a) printf '%s\n' "${FAKE_STAT_MODE:-600}" ;;
*) exit 1 ;;
esac
EOF
cat > "$HARNESS/bin/docker" <<'EOF'
#!/usr/bin/env bash
printf 'api_image=%s migrator_image=%s caddy_image=%s args=%s\n' \
  "${CAMPUS_LMS_API_IMAGE:?}" "${CAMPUS_LMS_MIGRATOR_IMAGE:?}" \
  "${CAMPUS_LMS_CADDY_IMAGE:?}" "$*" >> "${FAKE_DOCKER_LOG:?}"
case " $* " in
*" up "*)
  printf '%s\n%s\n%s\n' "$CAMPUS_LMS_API_IMAGE" "$CAMPUS_LMS_MIGRATOR_IMAGE" \
    "$CAMPUS_LMS_CADDY_IMAGE" > "${FAKE_ACTIVE_IMAGES:?}"
  ;;
esac
EOF
cat > "$HARNESS/bin/curl" <<'EOF'
#!/usr/bin/env bash
url="${!#}"
mapfile -t active < "${FAKE_ACTIVE_IMAGES:?}"
printf 'url=%s active_api=%s active_migrator=%s active_caddy=%s\n' \
  "$url" "${active[0]:-none}" "${active[1]:-none}" "${active[2]:-none}" >> "${FAKE_CURL_LOG:?}"
case "${FAKE_CURL_MODE:?}" in
success) exit 0 ;;
rollback)
  case "$url" in
  http://127.0.0.1:8080/readyz) exit 0 ;;
  https://*) [ "${active[2]:-}" != "${FAKE_CANDIDATE_CADDY:?}" ] ;;
  *) exit 1 ;;
  esac
  ;;
*) exit 2 ;;
esac
EOF
chmod 0755 "$HARNESS/bin/stat" "$HARNESS/bin/docker" "$HARNESS/bin/curl"

run_deploy() {
	PATH="$HARNESS/bin:$PATH" \
	FAKE_DOCKER_LOG="$HARNESS/docker.log" \
	FAKE_ACTIVE_IMAGES="$HARNESS/active-images" \
	FAKE_CURL_LOG="$HARNESS/curl.log" \
	FAKE_CANDIDATE_CADDY="$CANDIDATE_CADDY" \
	CAMPUS_LMS_API_IMAGE="$PREVIOUS_API" \
	CAMPUS_LMS_MIGRATOR_IMAGE="$PREVIOUS_MIGRATOR" \
	CAMPUS_LMS_CADDY_IMAGE="$PREVIOUS_CADDY" \
	CAMPUS_LMS_TARGET_API_IMAGE="$CANDIDATE_API" \
	CAMPUS_LMS_TARGET_MIGRATOR_IMAGE="$CANDIDATE_MIGRATOR" \
	CAMPUS_LMS_TARGET_CADDY_IMAGE="$CANDIDATE_CADDY" \
	CAMPUS_LMS_HOSTNAME='api.example.invalid' \
	CAMPUS_LMS_API_ENV_FILE="$api_env" \
	CAMPUS_LMS_MIGRATOR_ENV_FILE="$migrator_env" \
	CAMPUS_LMS_API_MEMORY_LIMIT=256m \
	CAMPUS_LMS_MIGRATOR_MEMORY_LIMIT=128m \
	CAMPUS_LMS_CADDY_MEMORY_LIMIT=64m \
	CAMPUS_LMS_CADDY_ORIGIN_CERT_FILE="$ORIGIN_CERT_FILE" \
	CAMPUS_LMS_CADDY_ORIGIN_KEY_FILE="$ORIGIN_KEY_FILE" \
	CAMPUS_LMS_CADDY_ORIGIN_CA_ROOT_FILE="$ORIGIN_CA_ROOT_FILE" \
	CAMPUS_LMS_STATE_DIR="$HARNESS/state" \
	CAMPUS_LMS_READY_TIMEOUT_SECONDS=10 \
	CAMPUS_LMS_READY_INTERVAL_SECONDS=1 \
	bash "$DEPLOY_SCRIPT"
}

set +e
owner_output="$(PATH="$HARNESS/bin:$PATH" FAKE_STAT_UID=1000 CAMPUS_LMS_API_IMAGE="$PREVIOUS_API" CAMPUS_LMS_MIGRATOR_IMAGE="$PREVIOUS_MIGRATOR" CAMPUS_LMS_CADDY_IMAGE="$PREVIOUS_CADDY" CAMPUS_LMS_HOSTNAME='api.example.invalid' CAMPUS_LMS_API_ENV_FILE="$api_env" CAMPUS_LMS_MIGRATOR_ENV_FILE="$migrator_env" CAMPUS_LMS_API_MEMORY_LIMIT=256m CAMPUS_LMS_MIGRATOR_MEMORY_LIMIT=128m CAMPUS_LMS_CADDY_MEMORY_LIMIT=64m CAMPUS_LMS_CADDY_ORIGIN_CERT_FILE="$ORIGIN_CERT_FILE" CAMPUS_LMS_CADDY_ORIGIN_KEY_FILE="$ORIGIN_KEY_FILE" CAMPUS_LMS_CADDY_ORIGIN_CA_ROOT_FILE="$ORIGIN_CA_ROOT_FILE" bash "$DEPLOY_SCRIPT" 2>&1)"
owner_status=$?
set -e
test "$owner_status" -ne 0
grep -Fq 'must be owned by UID 0' <<<"$owner_output"
echo 'secret_owner_rejected=pass'

printf '%s\n' "api_image=$PREVIOUS_API" "migrator_image=$PREVIOUS_MIGRATOR" \
	"caddy_image=$PREVIOUS_CADDY" 'hostname=api.example.invalid' > "$HARNESS/state/deployment-state"
chmod 0600 "$HARNESS/state/deployment-state"
printf '%s\n%s\n%s\n' "$PREVIOUS_API" "$PREVIOUS_MIGRATOR" "$PREVIOUS_CADDY" > "$HARNESS/active-images"

set +e
rollback_output="$(FAKE_CURL_MODE=rollback run_deploy 2>&1)"
rollback_status=$?
set -e
test "$rollback_status" -ne 0
grep -Fq 'candidate failed readiness; application rollback completed' <<<"$rollback_output"
grep -Fq "api_image=$CANDIDATE_API migrator_image=$CANDIDATE_MIGRATOR caddy_image=$CANDIDATE_CADDY" "$HARNESS/docker.log"
grep -F 'args=' "$HARNESS/docker.log" | grep -F 'run --rm --no-deps migrator' | grep -F "migrator_image=$CANDIDATE_MIGRATOR" >/dev/null
grep -F 'args=' "$HARNESS/docker.log" | grep -F 'pull api caddy' | grep -F "api_image=$PREVIOUS_API" | grep -F "caddy_image=$PREVIOUS_CADDY" >/dev/null
grep -F 'args=' "$HARNESS/docker.log" | grep -F 'up -d --no-deps api caddy' | grep -F "api_image=$PREVIOUS_API" | grep -F "caddy_image=$PREVIOUS_CADDY" >/dev/null
grep -F "url=https://api.example.invalid:8443/readyz active_api=$CANDIDATE_API active_migrator=$CANDIDATE_MIGRATOR active_caddy=$CANDIDATE_CADDY" "$HARNESS/curl.log" >/dev/null
grep -F "url=https://api.example.invalid:8443/readyz active_api=$PREVIOUS_API active_migrator=$PREVIOUS_MIGRATOR active_caddy=$PREVIOUS_CADDY" "$HARNESS/curl.log" >/dev/null
grep -Fq "api_image=$PREVIOUS_API" "$HARNESS/state/deployment-state"
grep -Fq "migrator_image=$PREVIOUS_MIGRATOR" "$HARNESS/state/deployment-state"
grep -Fq "caddy_image=$PREVIOUS_CADDY" "$HARNESS/state/deployment-state"
if grep -Eq '(^|[[:space:]])down([[:space:]]|$)' "$HARNESS/docker.log"; then
	echo 'automatic_down_command=true' >&2
	exit 1
fi
echo "rollback_candidate_api=$CANDIDATE_API"
echo "rollback_previous_api=$PREVIOUS_API"
echo "rollback_candidate_caddy=$CANDIDATE_CADDY"
echo "rollback_previous_caddy=$PREVIOUS_CADDY"
echo "rollback_candidate_migrator=$CANDIDATE_MIGRATOR"
echo 'rollback_digest_discrimination=pass'
echo 'rollback_previous_state_intact=pass'
echo 'rollback_caddy_readiness_after_previous_digest=pass'
echo 'automatic_down_command=false'

: > "$HARNESS/docker.log"
: > "$HARNESS/curl.log"
printf '%s\n' "api_image=$PREVIOUS_API" "migrator_image=$PREVIOUS_MIGRATOR" \
	"caddy_image=$PREVIOUS_CADDY" 'hostname=api.example.invalid' > "$HARNESS/state/deployment-state"
printf '%s\n%s\n%s\n' "$PREVIOUS_API" "$PREVIOUS_MIGRATOR" "$PREVIOUS_CADDY" > "$HARNESS/active-images"
success_output="$(FAKE_CURL_MODE=success run_deploy 2>&1)"
grep -Fq 'deployment=pass' <<<"$success_output"
grep -Fq "api_image=$CANDIDATE_API" "$HARNESS/state/deployment-state"
grep -Fq "migrator_image=$CANDIDATE_MIGRATOR" "$HARNESS/state/deployment-state"
grep -Fq "caddy_image=$CANDIDATE_CADDY" "$HARNESS/state/deployment-state"
echo 'success_state_contains_candidate_digests=pass'
echo 'deploy_script_safety=pass'
