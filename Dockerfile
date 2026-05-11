# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY main.go app.go ./
COPY internal ./internal

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    set -eux; \
    GOARM_VALUE=""; \
    if [ "${TARGETARCH}" = "arm" ] && [ "${TARGETVARIANT}" = "v7" ]; then \
      GOARM_VALUE="7"; \
    fi; \
    export CGO_ENABLED=0; \
    export GOOS="${TARGETOS:-linux}"; \
    export GOARCH="${TARGETARCH:-amd64}"; \
    if [ -n "${GOARM_VALUE}" ]; then export GOARM="${GOARM_VALUE}"; fi; \
    go build -trimpath -ldflags="-s -w" -o /out/visitor-ip-api-go .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/visitor-ip-api-go /app/visitor-ip-api-go

EXPOSE 8466

ENV LISTEN_ADDR=:8466

ENTRYPOINT ["/app/visitor-ip-api-go"]
