#!/usr/bin/env bash
# Local-only Caddy-facing readiness proof. It uses a temporary local certificate
# trusted explicitly by curl; this is not a substitute for Cloudflare Origin CA
# validation and deploy.sh itself never disables certificate verification.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CADDY_CONTAINER="campus-lms-phase4b-caddy-readiness"
TLS_DIR=""

cleanup() {
	docker rm -f "$CADDY_CONTAINER" >/dev/null 2>&1 || true
	make -C "$REPO_ROOT" down >/dev/null 2>&1 || true
	if [ -n "$TLS_DIR" ]; then rm -rf -- "$TLS_DIR"; fi
}
trap cleanup EXIT

cleanup
TLS_DIR="$(mktemp -d /tmp/campus-lms-caddy-tls.XXXXXX)"
openssl req -x509 -newkey rsa:2048 -nodes -days 1 \
	-subj '/CN=localhost' \
	-addext 'subjectAltName=DNS:localhost' \
	-keyout "$TLS_DIR/key.pem" \
	-out "$TLS_DIR/cert.pem" >/dev/null 2>&1
chmod 0600 "$TLS_DIR/key.pem"
make -C "$REPO_ROOT" up >/dev/null 2>&1

for attempt in $(seq 1 60); do
	if curl --fail --silent --show-error http://127.0.0.1:8080/readyz >/dev/null 2>&1; then
		break
	fi
	if [ "$attempt" = 60 ]; then
		echo 'internal_api_readiness=fail' >&2
		exit 1
	fi
	sleep 1
done

docker run -d --name "$CADDY_CONTAINER" --network campus-lms_default \
	-p 127.0.0.1:8443:443 \
	-e CADDY_HOSTNAME=localhost \
	-e CADDY_ORIGIN_CERT_FILE=/etc/caddy/origin/cert.pem \
	-e CADDY_ORIGIN_KEY_FILE=/etc/caddy/origin/key.pem \
	-v "$REPO_ROOT/deploy/caddy/Caddyfile:/etc/caddy/Caddyfile:ro" \
	-v "$TLS_DIR:/etc/caddy/origin:ro" \
	caddy:2.10.2-alpine \
	caddy run --config /etc/caddy/Caddyfile --adapter caddyfile >/dev/null

for attempt in $(seq 1 30); do
	if curl --fail --silent --show-error --cacert "$TLS_DIR/cert.pem" https://localhost:8443/readyz >/dev/null 2>&1; then
		echo 'internal_api_readiness=pass'
		echo 'caddy_proxy_readiness=pass'
		echo 'caddy_facing_mode=temporary_explicit_ca_only'
		echo 'phase4d_external_tls=not_executed'
		exit 0
	fi
	sleep 1
done

echo 'caddy_proxy_readiness=fail' >&2
exit 1
