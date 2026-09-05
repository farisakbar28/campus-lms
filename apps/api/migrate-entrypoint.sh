#!/bin/sh
set -eu

: "${MIGRATE_DATABASE_URL:?MIGRATE_DATABASE_URL is required}"

if [ "${MIGRATOR_LOCAL_TEST_MODE:-}" != "1" ]; then
	query=""
	case "$MIGRATE_DATABASE_URL" in
	*\?*)
		query="${MIGRATE_DATABASE_URL#*\?}"
		query="${query%%\#*}"
		;;
	esac

	sslmode_count=0
	sslmode_value=""
	old_ifs="$IFS"
	IFS='&'
	set -f
	set -- $query
	set +f
	IFS="$old_ifs"
	for parameter in "$@"; do
		case "$parameter" in
		sslmode=*)
			sslmode_count=$((sslmode_count + 1))
			sslmode_value="${parameter#sslmode=}"
			;;
		esac
	done

	if [ "$sslmode_count" -ne 1 ]; then
		echo "ERROR: production migration requires exactly one supported sslmode" >&2
		exit 1
	fi

	case "$sslmode_value" in
	require|verify-ca|verify-full) ;;
	*)
		echo "ERROR: production migration requires exactly one supported sslmode" >&2
		exit 1
		;;
	esac
fi

exec /usr/local/bin/migrate \
    -path /migrations \
    -database "$MIGRATE_DATABASE_URL" \
    up
