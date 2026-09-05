#!/usr/bin/env bash
set -euo pipefail

export GOCACHE="${GOCACHE:-/tmp/campus-lms-go-build}"

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
api_dir="$repo_root/apps/api"
api_binary="$(mktemp)"
api_log="$(mktemp)"
pid=""

cleanup() {
	if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
		kill -TERM "$pid"
		wait "$pid" || true
	fi
	rm -f "$api_binary" "$api_log"
}
trap cleanup EXIT

(
	cd "$api_dir"
	go build -o "$api_binary" ./cmd/api
)

APP_ENV=development APP_PORT=18081 APP_LOG_LEVEL=info APP_SHUTDOWN_TIMEOUT=2s \
	"$api_binary" >"$api_log" 2>&1 &
pid="$!"
echo "PID: $pid"

for _ in {1..50}; do
	if curl --silent --fail --output /dev/null http://127.0.0.1:18081/readyz; then
		break
	fi
	sleep 0.1
done

echo "ACTION: sending SIGTERM to PID $pid"
kill -TERM "$pid"

set +e
wait "$pid"
exit_code="$?"
set -e
echo "EXIT CODE: $exit_code"
cat "$api_log"
pid=""

exit "$exit_code"
