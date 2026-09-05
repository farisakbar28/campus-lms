#!/usr/bin/env bash
# Reproducible build and provenance inspection for the dedicated migrator.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
IMAGE="${MIGRATOR_IMAGE:-campus-lms-migrator:phase4b}"
CONTAINER="campus-lms-phase4b-migrator-inspect"
TEMPORARY="$(mktemp -d /tmp/campus-lms-migrator-build.XXXXXX)"

cleanup() {
	docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
	rm -rf -- "$TEMPORARY"
}
trap cleanup EXIT

docker build --quiet --target migrator -t "$IMAGE" -f "$REPO_ROOT/apps/api/Dockerfile" "$REPO_ROOT/apps/api" >/dev/null
container_id="$(docker create --name "$CONTAINER" "$IMAGE")"
docker cp "$container_id:/usr/local/bin/migrate" "$TEMPORARY/migrate"
metadata="$(go version -m "$TEMPORARY/migrate")"

grep -Fq $'\tpath\tgithub.com/golang-migrate/migrate/v4/cmd/migrate' <<<"$metadata"
grep -Fq $'\tmod\tgithub.com/golang-migrate/migrate/v4\tv4.18.3' <<<"$metadata"

printf 'image_id=%s\n' "$(docker image inspect --format '{{.Id}}' "$IMAGE")"
printf 'architecture=%s\n' "$(docker image inspect --format '{{.Architecture}}' "$IMAGE")"
printf 'user=%s\n' "$(docker image inspect --format '{{.Config.User}}' "$IMAGE")"
printf 'entrypoint=%s\n' "$(docker image inspect --format '{{json .Config.Entrypoint}}' "$IMAGE")"
printf 'healthcheck=%s\n' "$(docker image inspect --format '{{json .Config.Healthcheck.Test}}' "$IMAGE")"
printf '%s\n' "$metadata"
echo 'migrator_build_and_version_verification=pass'
