#!/usr/bin/env bash
# -----------------------------------------------------------------------------
# TASK BRIEF — Week 4
# Agent: read agent/rules/20-database.md.
# -----------------------------------------------------------------------------
#
# REQUIREMENTS
#   - pg_dump in custom format (-Fc), not plain SQL
#   - filename contains a UTC timestamp
#   - upload to object storage (Supabase Storage / MinIO)
#   - 7-day retention; delete older dumps
#   - Telegram notification ON FAILURE (silent backup failure is the most
#     common way people lose data)
#   - non-zero exit on failure so cron/systemd records it
#
# Golden rule: a backup that has never been restored is not a backup.
set -euo pipefail

echo "TODO Week 4" >&2
exit 1
