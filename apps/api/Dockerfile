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

FROM --platform=$BUILDPLATFORM golang:1.23.6-alpine3.21 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Copy the module definition before application source so dependency downloads
# remain cached when only Go source files change.
COPY go.mod go.sum* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# distroless intentionally has neither a shell nor an HTTP client. BusyBox is
# copied solely to run the Docker health probe without adding a full OS layer.
FROM busybox:1.37.0-musl AS healthcheck

# This target is retained only for the Week 2 size comparison. It deliberately
# keeps the Go toolchain in the runtime image; production uses the final stage.
FROM golang:1.23.6-alpine3.21 AS alpine-runtime

COPY --from=builder /out/api /api

USER 65532:65532

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["wget", "--spider", "--quiet", "http://127.0.0.1:8080/healthz"]

ENTRYPOINT ["/api"]

FROM gcr.io/distroless/static-debian12:nonroot AS busybox-runtime

COPY --from=builder /out/api /api
COPY --from=healthcheck /bin/busybox /busybox

USER nonroot:nonroot

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/busybox", "wget", "--spider", "--quiet", "http://127.0.0.1:8080/healthz"]

ENTRYPOINT ["/api"]

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/api /api

USER nonroot:nonroot

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/api", "-healthcheck"]

ENTRYPOINT ["/api"]
