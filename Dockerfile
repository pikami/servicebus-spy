# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM oven/bun:1-alpine AS web-builder

WORKDIR /app/web

COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile

COPY web/ ./
RUN bun run build


FROM --platform=$BUILDPLATFORM alpine:3.21 AS certs

RUN apk add --no-cache ca-certificates


FROM alpine:3.21

WORKDIR /app

COPY servicebus-spy .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=web-builder /app/web/dist ./web/dist

ENV SERVICEBUS_SPY_WEB_DIST=/app/web/dist

ENTRYPOINT ["/app/servicebus-spy"]
