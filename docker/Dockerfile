# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22
ARG ALPINE_VERSION=3.20

FROM golang:${GO_VERSION}-alpine AS builder

ARG APP_VERSION=dev
ARG APP_COMMIT_HASH=0
ARG APP_BUILD_TIME=0

WORKDIR /src
COPY . .

RUN go build \
    -ldflags="-w -s \
    -X 'main.version=${APP_VERSION}' \
    -X 'main.commit=${APP_COMMIT_HASH}' \
    -X 'main.date=${APP_BUILD_TIME}'" \
    -o /app \
    ./cmd

FROM alpine:${ALPINE_VERSION} AS app

COPY --from=builder /app /app

CMD ["/app"]
