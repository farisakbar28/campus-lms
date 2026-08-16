# -----------------------------------------------------------------------------
# TASK BRIEF — Week 2
# Agent: read agent/rules/30-docker-deploy.md before implementing.
# -----------------------------------------------------------------------------
#
# TARGETS
#   - multi-stage build (builder + minimal runtime)
#   - final image < 25 MB   (verify: make docker-size)
#   - runs as a NON-ROOT user (verify with docker inspect)
#   - HEALTHCHECK present
#   - builds for both linux/amd64 and linux/arm64
#   - layer caching: `go mod download` in its own layer before copying source
#
# HARD RULES
#   - never place a secret in ARG or ENV (it persists in image layers)
#   - CGO_ENABLED=0 for a static binary
#
# QUESTIONS THE HUMAN MUST BE ABLE TO ANSWER AFTERWARDS
#   - Why does CGO_ENABLED=0 matter for scratch/distroless images?
#   - distroless vs alpine: what is the debugging trade-off?
#   - Why does COPY order dominate rebuild speed?
#   - Why is a secret in a build ARG trivially extractable?
#
# EVIDENCE REQUIRED
#   docs/progress/evidence/week-02/image-size.txt
#   docs/progress/evidence/week-02/docker-inspect-user.txt
#   docs/progress/evidence/week-02/buildx-multiarch.txt

# TODO Week 2
