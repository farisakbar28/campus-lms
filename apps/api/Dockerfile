# Builds static API and migration-runner artifacts.
# The API runtime stays non-root; build ARG and ENV must never contain secrets
# because image layers retain their values.

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

# The migration runner is a separate artifact so the production API runtime
# does not gain a database-administration binary or a shell. The version is
# pinned here because production deploys must never install tools dynamically.
FROM builder AS migrator-builder

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go install -tags postgres \
    github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

# distroless intentionally has neither a shell nor an HTTP client. BusyBox is
# copied solely to run the Docker health probe without adding a full OS layer.
FROM busybox:1.37.0-musl AS healthcheck

# This comparison target retains the Go toolchain; production uses the final stage.
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

FROM gcr.io/distroless/static-debian12:nonroot AS api-runtime

COPY --from=builder /out/api /api

USER nonroot:nonroot

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/api", "-healthcheck"]

ENTRYPOINT ["/api"]

# This image is used only by the one-shot production Compose migrator service.
# Alpine supplies the shell needed to expand MIGRATE_DATABASE_URL at runtime;
# the API image above remains distroless and shell-free.
FROM alpine:3.21 AS migrator

RUN apk add --no-cache ca-certificates \
    && addgroup -S migrator \
    && adduser -S -G migrator migrator

COPY --from=migrator-builder /go/bin/migrate /usr/local/bin/migrate
COPY migrations /migrations
COPY migrate-entrypoint.sh /usr/local/bin/production-migrate

RUN chmod 0555 /usr/local/bin/migrate /usr/local/bin/production-migrate \
    && chown -R migrator:migrator /migrations

USER migrator:migrator

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/migrate", "-version"]

ENTRYPOINT ["/usr/local/bin/production-migrate"]

# Preserve the historical no-target build contract used by `make docker-build`:
# it must continue to produce the API image, while the migrator is selected
# explicitly with `--target migrator`.
FROM api-runtime AS final
