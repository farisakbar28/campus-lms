#!/usr/bin/env bash
# -----------------------------------------------------------------------------
# TASK BRIEF — Week 4
# Agent: read agent/rules/30-docker-deploy.md. This script runs on the VM;
# the agent writes it but the HUMAN runs it the first time.
# -----------------------------------------------------------------------------
#
# REQUIREMENTS
#   - pull image from GHCR by DIGEST, never by :latest tag
#   - run database migrations (expand-contract pattern, Week 5)
#   - docker compose up -d with the production file
#   - health check: poll /readyz until 200, 60s timeout
#   - AUTOMATIC ROLLBACK after 3 failures: revert to the previous digest
#   - print total deploy duration (this becomes the "lead time" metric)
#
# DoD: deploy completes with zero failed requests (prove with a curl loop
# during deploy), and rollback finishes in under 2 minutes.
set -euo pipefail

echo "TODO Week 4" >&2
exit 1
